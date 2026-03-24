// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rpc

import (
	"context"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	jxpkg "github.com/go-faster/jx"

	"github.com/pkg/errors"
)

const pathStatus = "status"

func (api *API) CurrentHead(ctx context.Context) (pkgTypes.Level, error) {
	var level pkgTypes.Level
	err := api.getStream(ctx, pathStatus, nil, func(d *jxpkg.Decoder) error {
		return jxResponse(d, func(d *jxpkg.Decoder) error {
			var err error
			level, err = jxStatusMinimal(d)
			return err
		})
	})
	return level, errors.Wrap(err, "CurrentHead")
}
