package rpc

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/dipdup-net/go-lib/config"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
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

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = rps
	t.MaxConnsPerHost = rps
	t.MaxIdleConnsPerHost = rps

	return API{
		client: &http.Client{
			Transport: t,
		},
		cfg:       cfg,
		rps:       rps,
		rateLimit: rate.NewLimiter(rate.Every(time.Second/time.Duration(rps)), rps),
		log:       log.With().Str("module", "node rpc").Logger(),
	}
}

// get -
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	response, err := api.client.Do(req)
	if err != nil {
		return err
	}

	defer func(b io.ReadCloser) {
		if err := b.Close(); err != nil {
			log.Err(err).Msg("api close GET body request")
		}
	}(response.Body)

	err = json.NewDecoder(response.Body).Decode(output)
	return err
}
