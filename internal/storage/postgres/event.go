// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Event -
type Event struct {
	*postgres.Table[*storage.Event]
}

// NewEvent -
func NewEvent(db *database.Bun) *Event {
	return &Event{
		Table: postgres.NewTable[*storage.Event](db),
	}
}

// ByTxId -
func (e *Event) ByTxId(ctx context.Context, txId uint64, limit, offset int) (events []storage.Event, err error) {
	query := e.DB().NewSelect().Model(&events).
		Where("tx_id = ?", txId)
	query = limitScope(query, limit)

	if offset > 0 {
		query = query.Offset(offset)
	}

	err = query.Scan(ctx)
	return
}

// ByBlock -
func (e *Event) ByBlock(ctx context.Context, height pkgTypes.Level, limit, offset int) (events []storage.Event, err error) {
	query := e.DB().NewSelect().Model(&events).
		Where("height = ?", height).
		Where("tx_id IS NULL")

	query = limitScope(query, limit)

	if offset > 0 {
		query = query.Offset(offset)
	}
	err = query.Scan(ctx)
	return
}
