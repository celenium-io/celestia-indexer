package hyperlane

import (
	"bufio"
	"bytes"
	"context"
	"gopkg.in/yaml.v3"
	"time"

	fastshot "github.com/opus-domini/fast-shot"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

type Api struct {
	client    fastshot.ClientHttpMethods
	rateLimit *rate.Limiter
	timeout   time.Duration
}

func NewApi(baseUrl string, opts ...ApiOption) Api {
	api := Api{
		client:    fastshot.NewClient(baseUrl).Build(),
		rateLimit: rate.NewLimiter(rate.Every(time.Second/time.Duration(10)), 10),
		timeout:   time.Second * 30,
	}

	for i := range opts {
		opts[i](&api)
	}

	return api
}

func (api Api) ChainMetadata(ctx context.Context) (map[uint64]ChainMetadata, error) {
	if api.rateLimit != nil {
		if err := api.rateLimit.Wait(ctx); err != nil {
			return nil, err
		}
	}

	requestCtx, cancel := context.WithTimeout(ctx, api.timeout)
	defer cancel()

	resp, err := api.client.GET("/hyperlane-xyz/hyperlane-registry/main/chains/metadata.yaml").
		Context().Set(requestCtx).
		Send()
	if err != nil {
		return nil, err
	}

	if resp.Status().IsError() {
		return nil, errors.Errorf("invalid status: %d", resp.Status().Code())
	}

	if resp.Status().IsError() {
		return nil, errors.Errorf("invalid status: %d", resp.Status().Code())
	}

	scanner := bufio.NewScanner(resp.Raw().Body)
	var buf bytes.Buffer
	skipped := false
	for scanner.Scan() {
		if !skipped {
			skipped = true
			continue
		}
		buf.Write(scanner.Bytes())
		buf.WriteByte('\n')
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var raw map[string]ChainMetadata
	if err := yaml.Unmarshal(buf.Bytes(), &raw); err != nil {
		return nil, err
	}

	result := make(map[uint64]ChainMetadata)
	for _, data := range raw {
		result[data.DomainId] = data
	}

	return result, nil
}
