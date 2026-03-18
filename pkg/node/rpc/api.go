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
	func() *jxpkg.Decoder { return jxpkg.Decode(nil, 4*1024*1024) }, // 4 МБ
)

// gzipPool reuses klauspost gzip readers to avoid per-request allocations.
// klauspost/compress/gzip is significantly faster than stdlib compress/flate for
// decompression, using SIMD and other optimisations.
// Factory returns a zero-value Reader; Reset initialises all fields before use.
var gzipPool = pool.New(func() *kgzip.Reader { return &kgzip.Reader{} })

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

// openBody returns an io.Reader for the response body. If the server sent a
// gzip-compressed response, the returned reader is a pooled klauspost gzip
// reader; callers must call the returned cleanup function when done.
// The underlying resp.Body is always closed separately via closeWithLogError.
func (api *API) openBody(resp *http.Response) (io.Reader, func()) {
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := acquireGzipReader(resp.Body)
		if err != nil {
			api.log.Warn().Err(err).
				Str("url", resp.Request.URL.String()).
				Msg("gzip reader init failed, reading raw body")
			return resp.Body, nil
		}
		return gz, func() { releaseGzipReader(gz) }
	}
	return resp.Body, nil
}

const (
	celeniumUserAgent = "Celenium Indexer"
)

type API struct {
	client    *http.Client
	cfg       config.DataSource
	rps       int
	rateLimit *rate.Limiter
	log       zerolog.Logger
}

func NewAPI(cfg config.DataSource) API {
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

	return API{
		client: &http.Client{
			Transport: t,
			Timeout:   time.Second * time.Duration(timeout),
		},
		cfg:       cfg,
		rps:       rps,
		rateLimit: rate.NewLimiter(rate.Every(time.Second/time.Duration(rps)), rps),
		log:       log.With().Str("module", "node rpc").Logger(),
	}
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
	req.Header.Set("Accept-Encoding", "gzip")

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

	streamBody, cleanup := api.openBody(response)
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
	req.Header.Set("Accept-Encoding", "gzip")

	response, err := api.client.Do(req) //nolint:gosec
	if err != nil {
		return err
	}
	defer closeWithLogError(response.Body, api.log)

	api.log.Trace().
		Int64("ms", time.Since(start).Milliseconds()).
		Str("url", u.String()).
		Msg("post request")

	if response.StatusCode != http.StatusOK {
		return errors.Errorf("invalid status: %d", response.StatusCode)
	}

	streamBody, cleanup := api.openBody(response)
	if cleanup != nil {
		defer cleanup()
	}

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
