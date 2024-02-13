// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/hex"
	"strings"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

// Search -
type Search struct {
	db *database.Bun
}

// NewSearch -
func NewSearch(db *database.Bun) *Search {
	return &Search{
		db: db,
	}
}

func (s *Search) Search(ctx context.Context, query []byte) (results []storage.SearchResult, err error) {
	blockQuery := s.db.DB().NewSelect().
		Model((*storage.Block)(nil)).
		ColumnExpr("id, ? as value, 'block' as type", hex.EncodeToString(query)).
		Where("hash = ?", query).
		WhereOr("data_hash = ?", query)
	txQuery := s.db.DB().NewSelect().
		Model((*storage.Tx)(nil)).
		ColumnExpr("id, encode(hash, 'hex') as value, 'tx' as type").
		Where("hash = ?", query)

	union := blockQuery.UnionAll(txQuery)

	err = s.db.DB().NewSelect().
		TableExpr("(?) as search", union).
		Limit(10).
		Offset(0).
		Scan(ctx, &results)

	return
}

func (s *Search) SearchText(ctx context.Context, text string) (results []storage.SearchResult, err error) {
	text = strings.ToUpper(text)
	text = "%" + text + "%"
	validatorQuery := s.db.DB().NewSelect().
		Model((*storage.Validator)(nil)).
		ColumnExpr("id, moniker as value, 'validator' as type").
		Where("moniker ILIKE ?", text)
	rollupQuery := s.db.DB().NewSelect().
		Model((*storage.Rollup)(nil)).
		ColumnExpr("id, name as value, 'rollup' as type").
		Where("name ILIKE ?", text)

	union := rollupQuery.UnionAll(validatorQuery)

	err = s.db.DB().NewSelect().
		TableExpr("(?) as search", union).
		Limit(10).
		Offset(0).
		Scan(ctx, &results)

	return
}
