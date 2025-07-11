// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package blob

import (
	context "context"
	"strconv"
	"strings"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/node/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/go-lib/config"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/opus-domini/fast-shot/constant/mime"
	"github.com/pkg/errors"
)

type Celenium struct {
	client fastshot.ClientHttpMethods
}

func NewCelenium(datasource config.DataSource) *Celenium {
	timeout := time.Second * 30
	if datasource.Timeout > 0 {
		timeout = time.Duration(datasource.Timeout) * time.Second
	}
	return &Celenium{
		client: fastshot.NewClient(datasource.URL).
			Config().SetTimeout(timeout).
			Build(),
	}
}

func (c *Celenium) Blobs(ctx context.Context, height pkgTypes.Level, hash ...string) ([]types.Blob, error) {
	body := map[string]string{
		"height": strconv.FormatInt(int64(height), 10),
	}
	if len(hash) > 0 {
		body["namespaces"] = strings.Join(hash, ",")
	}
	response, err := c.client.
		GET("/v1/blob").
		Context().Set(ctx).
		Header().AddContentType(mime.JSON).
		Query().AddParams(body).
		Send()
	if err != nil {
		return nil, err
	}
	if response.Status().IsError() {
		str, err := response.Body().AsString()
		if err != nil {
			return nil, errors.Wrap(err, "reading error message")
		}
		return nil, errors.New(str)
	}
	var blobs []types.Blob
	err = response.Body().AsJSON(&blobs)
	return blobs, err
}

func (c *Celenium) Blob(ctx context.Context, height pkgTypes.Level, hash, commitment string) (types.Blob, error) {
	body := map[string]any{
		"commitment": commitment,
		"hash":       hash,
		"height":     height,
	}
	response, err := c.client.
		POST("/v1/blob").
		Context().Set(ctx).
		Header().AddContentType(mime.JSON).
		Body().AsJSON(body).
		Send()
	if err != nil {
		return types.Blob{}, err
	}
	if response.Status().IsError() {
		str, err := response.Body().AsString()
		if err != nil {
			return types.Blob{}, errors.Wrap(err, "reading error message")
		}
		return types.Blob{}, errors.New(str)
	}
	var blob types.Blob
	err = response.Body().AsJSON(&blob)
	return blob, err
}
