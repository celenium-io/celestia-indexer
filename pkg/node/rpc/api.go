// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/pool"
	"github.com/celenium-io/celestia-indexer/pkg/node/types"
	jxpkg "github.com/go-faster/jx"
	kgzip "github.com/klauspost/compress/gzip"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/dipdup-io/go-lib/config"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

var decoderPool = pool.New(
	func() *jxpkg.Decoder { return jxpkg.Decode(nil, 512*1024) }, // 512 KB — large strings handled via StrAppend+base64BufPool
)

// gzipPool reuses klauspost gzip readers to avoid per-request allocations.
// Factory returns a zero-value Reader; Reset initialises all fields before use.
var gzipPool = pool.New(func() *kgzip.Reader { return &kgzip.Reader{} })

// bodyBufPool reuses large byte buffers for reading compressed and decompressed
// response bodies entirely into memory before parsing. This eliminates the
// per-byte syscall pattern that occurs when the flate decompressor reads
// directly from the network socket via bufio.ReadByte.
var bodyBufPool = pool.New(func() *bytes.Buffer {
	return bytes.NewBuffer(make([]byte, 0, 16*1024*1024)) // 16 MB initial
})

// bodyBufMaxCap is the maximum buffer capacity to return to bodyBufPool.
// Buffers that grew beyond this threshold are discarded so the pool does not
// permanently hold memory from unusually large blocks.
const bodyBufMaxCap = 64 * 1024 * 1024 // 64 MB

func releaseBodyBuf(buf *bytes.Buffer) {
	if buf.Cap() > bodyBufMaxCap {
		return // let GC collect oversized buffers
	}
	buf.Reset()
	bodyBufPool.Put(buf)
}

func acquireGzipReader(r io.Reader) (*kgzip.Reader, error) {
	gz := gzipPool.Get()
	if err := gz.Reset(r); err != nil {
		gzipPool.Put(gz)
		return nil, err
	}
	return gz, nil
}

func releaseGzipReader(gz *kgzip.Reader) {
	_ = gz.Close()
	gzipPool.Put(gz)
}

// openBody reads the entire response body into memory, then (if gzip-encoded)
// decompresses it fully in-memory before returning a bytes.Reader to the caller.
// Keeping both steps in RAM eliminates the syscall-per-bufio-fill overhead that
// dominates CPU when the flate decompressor reads directly from the network.
// The returned cleanup function must be called when the reader is no longer needed.
// openBody reads the entire response body and optionally decompresses gzip.
// Returns the reader, compressed byte count, uncompressed byte count, and cleanup func.
func (api *API) openBody(resp *http.Response) (io.Reader, int64, int64, func()) {
	compBuf := bodyBufPool.Get()
	compBuf.Reset()
	if _, err := io.Copy(compBuf, resp.Body); err != nil {
		releaseBodyBuf(compBuf)
		api.log.Warn().Err(err).
			Str("url", resp.Request.URL.String()).
			Msg("reading response body failed, streaming raw")
		return resp.Body, 0, 0, nil
	}

	compressedBytes := int64(compBuf.Len())

	if resp.Header.Get("Content-Encoding") != "gzip" {
		return bytes.NewReader(compBuf.Bytes()), compressedBytes, compressedBytes, func() { releaseBodyBuf(compBuf) }
	}

	gz, err := acquireGzipReader(bytes.NewReader(compBuf.Bytes()))
	if err != nil {
		api.log.Warn().Err(err).
			Str("url", resp.Request.URL.String()).
			Msg("gzip reader init failed, reading raw body")
		return bytes.NewReader(compBuf.Bytes()), compressedBytes, compressedBytes, func() { releaseBodyBuf(compBuf) }
	}

	decompBuf := bodyBufPool.Get()
	decompBuf.Reset()
	if _, err := io.Copy(decompBuf, gz); err != nil {
		releaseGzipReader(gz)
		releaseBodyBuf(compBuf)
		releaseBodyBuf(decompBuf)
		api.log.Warn().Err(err).
			Str("url", resp.Request.URL.String()).
			Msg("gzip decompression failed, streaming raw")
		return resp.Body, 0, 0, nil
	}
	releaseGzipReader(gz)

	uncompressedBytes := int64(decompBuf.Len())
	releaseBodyBuf(compBuf)

	return bytes.NewReader(decompBuf.Bytes()), compressedBytes, uncompressedBytes, func() { releaseBodyBuf(decompBuf) }
}

const (
	celeniumUserAgent = "Celenium Indexer"
)

type API struct {
	client      *http.Client
	cfg         config.DataSource
	rps         int
	rateLimit   *rate.Limiter
	log         zerolog.Logger
	disableGzip bool
}

func NewAPI(cfg config.DataSource, opts ...ApiOption) API {
	rps := cfg.RequestsPerSecond
	if cfg.RequestsPerSecond < 1 || cfg.RequestsPerSecond > 100 {
		rps = 10
	}
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 10
	}

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	// Disable stdlib's transparent gzip decompression so we can use the
	// faster klauspost/compress implementation instead.
	t.DisableCompression = true
	t.ReadBufferSize = 256 * 1024

	a := API{
		client: &http.Client{
			Transport: t,
			Timeout:   time.Second * time.Duration(timeout),
		},
		cfg:       cfg,
		rps:       rps,
		rateLimit: rate.NewLimiter(rate.Every(time.Second/time.Duration(rps)), rps),
		log:       log.With().Str("module", "node rpc").Logger(),
	}
	for _, opt := range opts {
		opt(&a)
	}
	return a
}

func (api *API) getStream(ctx context.Context, path string, args map[string]string, fn func(*jxpkg.Decoder) error) error {
	u, err := url.Parse(api.cfg.URL)
	if err != nil {
		return err
	}
	u.Path, err = url.JoinPath(u.Path, path)
	if err != nil {
		return err
	}

	values := u.Query()
	for key, value := range args {
		values.Add(key, value)
	}
	u.RawQuery = values.Encode()

	if api.rateLimit != nil {
		if err := api.rateLimit.Wait(ctx); err != nil {
			return err
		}
	}

	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", celeniumUserAgent)
	if !api.disableGzip {
		req.Header.Set("Accept-Encoding", "gzip")
	}

	response, err := api.client.Do(req) //nolint:gosec
	if err != nil {
		return err
	}
	defer closeWithLogError(response.Body, api.log)

	api.log.Trace().
		Int64("ms", time.Since(start).Milliseconds()).
		Str("url", u.String()).
		Msg("request")

	if response.StatusCode != http.StatusOK {
		return errors.Errorf("invalid status: %d", response.StatusCode)
	}

	streamBody, _, _, cleanup := api.openBody(response)
	if cleanup != nil {
		defer cleanup()
	}
	d := decoderPool.Get()
	d.Reset(streamBody)
	defer decoderPool.Put(d)
	return fn(d)
}

func (api *API) postStream(ctx context.Context, requests []types.Request, fn func(*jxpkg.Decoder) error) error {
	u, err := url.Parse(api.cfg.URL)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(requests); err != nil {
		return errors.Wrap(err, "invalid bulk post request")
	}

	if api.rateLimit != nil {
		if err := api.rateLimit.Wait(ctx); err != nil {
			return err
		}
	}

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", celeniumUserAgent)
	if !api.disableGzip {
		req.Header.Set("Accept-Encoding", "gzip")
	}

	response, err := api.client.Do(req) //nolint:gosec
	if err != nil {
		return err
	}
	defer closeWithLogError(response.Body, api.log)

	if response.StatusCode != http.StatusOK {
		return errors.Errorf("invalid status: %d", response.StatusCode)
	}

	streamBody, compressedBytes, uncompressedBytes, cleanup := api.openBody(response)
	if cleanup != nil {
		defer cleanup()
	}

	elapsed := time.Since(start)
	mbps := float64(uncompressedBytes) / elapsed.Seconds() / (1024 * 1024)
	api.log.Debug().
		Int64("ms", elapsed.Milliseconds()).
		Int64("compressed_kb", compressedBytes/1024).
		Int64("uncompressed_kb", uncompressedBytes/1024).
		Float64("mbps", mbps).
		Str("url", u.String()).
		Msg("post request")

	d := decoderPool.Get()
	d.Reset(streamBody)
	defer decoderPool.Put(d)

	return fn(d)
}

func closeWithLogError(stream io.ReadCloser, log zerolog.Logger) {
	if _, err := io.Copy(io.Discard, stream); err != nil {
		log.Err(err).Msg("api copy GET body response to discard")
	}
	if err := stream.Close(); err != nil {
		log.Err(err).Msg("api close GET body request")
	}
}
