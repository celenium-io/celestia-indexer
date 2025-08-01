// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rpc

import (
	"context"

	"github.com/celenium-io/celestia-indexer/pkg/types"
)

func (api *API) Head(ctx context.Context) (types.ResultBlock, error) {
	return api.Block(ctx, 0)
}
