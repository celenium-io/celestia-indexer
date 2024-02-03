// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package rpc

import (
	"context"
	"github.com/celenium-io/celestia-indexer/pkg/types"
)

func (api *API) Head(ctx context.Context) (types.ResultBlock, error) {
	return api.Block(ctx, 0)
}
