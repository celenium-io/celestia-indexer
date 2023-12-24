// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package binance

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/dipdup-net/go-lib/config"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type API struct {
	client    *http.Client
	cfg       config.DataSource
	rateLimit *rate.Limiter
	log       zerolog.Logger
}

func NewAPI(cfg config.DataSource) API {
	rps := cfg.RequestsPerSecond
	if cfg.RequestsPerSecond < 1 || cfg.RequestsPerSecond > 100 {
		rps = 10
	}

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = rps
	t.MaxConnsPerHost = rps
	t.MaxIdleConnsPerHost = rps

	return API{
		client: &http.Client{
			Transport: t,
		},
		cfg:       cfg,
		rateLimit: rate.NewLimiter(rate.Every(time.Second/time.Duration(rps)), rps),
		log:       log.With().Str("module", "binance_api").Logger(),
	}
}

func (api *API) get(ctx context.Context, path string, args map[string]string, output any) error {
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

	response, err := api.client.Do(req)
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

	err = json.NewDecoder(response.Body).Decode(output)
	return err
}

func closeWithLogError(stream io.ReadCloser, log zerolog.Logger) {
	if _, err := io.Copy(io.Discard, stream); err != nil {
		log.Err(err).Msg("binance api copy GET body response to discard")
	}
	if err := stream.Close(); err != nil {
		log.Err(err).Msg("binance api close GET body request")
	}
}
