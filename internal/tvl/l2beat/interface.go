// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package l2beat

import (
	"context"
	"github.com/celenium-io/celestia-indexer/internal/storage"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IApi interface {
	TVL(ctx context.Context, rollupName string, timeframe storage.TvlTimeframe) (result TVLResponse, err error)
}
