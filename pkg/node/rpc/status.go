// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rpc

import (
	"context"

	"github.com/celenium-io/celestia-indexer/pkg/node/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
)

const pathStatus = "status"

type syncInfoMinimal struct {
	LatestBlockHeight pkgTypes.Level `json:"latest_block_height,string"`
}
type statusMinimal struct {
	SyncInfo syncInfoMinimal `json:"sync_info"`
}

func (api *API) CurrentHead(ctx context.Context) (pkgTypes.Level, error) {
	var sr types.Response[statusMinimal]
	if err := api.get(ctx, pathStatus, nil, &sr); err != nil {
		return 0, errors.Wrap(err, "api.get")
	}

	if sr.Error != nil {
		return 0, errors.Wrapf(types.ErrRequest, "current head request %d error: %s", sr.Id, sr.Error.Error())
	}

	return sr.Result.SyncInfo.LatestBlockHeight, nil
}
