// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package hyperlane

import (
	"context"
	"io"

	"github.com/celenium-io/celestia-indexer/pkg/node/hyperlane"
)

//go:generate mockgen -source=$GOFILE -destination=mock.go -package=hyperlane -typed
type IChainStore interface {
	io.Closer

	Start(ctx context.Context)
	Get(domainId uint64) (hyperlane.ChainMetadata, bool)
	Set(metadata map[uint64]hyperlane.ChainMetadata)
	All() map[uint64]hyperlane.ChainMetadata
}
