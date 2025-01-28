// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package api

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/celestials"
	"github.com/goccy/go-json"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/opus-domini/fast-shot/constant/mime"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

type Api struct {
	client      fastshot.ClientHttpMethods
	timeout     time.Duration
	rateLimiter *rate.Limiter
}

func New(baseUrl string, opts ...ApiOption) Api {
	api := Api{
		client:      fastshot.NewClient(baseUrl).Build(),
		timeout:     time.Second * 10,
		rateLimiter: rate.NewLimiter(rate.Every(time.Second/time.Duration(5)), 5),
	}

	for i := range opts {
		opts[i](&api)
	}

	return api
}

func (api Api) Changes(ctx context.Context, chainId string, opts ...celestials.ChangeOption) (changes celestials.Changes, err error) {
	var opt = celestials.ChangeOptions{
		ChainId: chainId,
	}
	for i := range opts {
		opts[i](&opt)
	}

	response, err := api.client.POST("api/resolver/changes").
		Context().Set(ctx).
		Body().AsJSON(opt).
		Header().AddContentType(mime.JSON).
		Send()
	if err != nil {
		return changes, err
	}

	if response.Status().IsError() {
		body, err := response.Body().AsString()
		if err != nil {
			return changes, err
		}
		return changes, errors.Errorf("status=%d text=%s", response.Status().Code(), body)
	}

	err = json.NewDecoder(response.Raw().Body).Decode(&changes)
	return
}
