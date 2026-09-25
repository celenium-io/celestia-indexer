// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/hex"
	"strings"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	celestials "github.com/celenium-io/celestial-module/pkg/storage"
	"github.com/uptrace/bun"
)

// Search -
type Search struct {
	db bun.IDB
}

// NewSearch -
func NewSearch(db bun.IDB) storage.ISearch {
	return &Search{
		db: db,
	}
}

func (s *Search) Search(ctx context.Context, query []byte) (results []storage.SearchResult, err error) {
	blockQuery := s.db.NewSelect().
		Model((*storage.Block)(nil)).
		ColumnExpr("id, ? as value, 'block' as type", hex.EncodeToString(query)).
		Where("hash = ?", query).
		WhereOr("data_hash = ?", query)
	txQuery := s.db.NewSelect().
		Model((*storage.Tx)(nil)).
		ColumnExpr("id, encode(hash, 'hex') as value, 'tx' as type").
		Where("hash = ?", query)

	union := blockQuery.UnionAll(txQuery)

	err = s.db.NewSelect().
		TableExpr("(?) as search", union).
		Limit(10).
		Offset(0).
		Scan(ctx, &results)

	return
}

func (s *Search) SearchText(ctx context.Context, text string) (results []storage.SearchResult, err error) {
	text = strings.ToUpper(text)
	text = "%" + text + "%"
	validatorQuery := s.db.NewSelect().
		Model((*storage.Validator)(nil)).
		ColumnExpr("id, moniker as value, 'validator' as type").
		Where("moniker ILIKE ?", text)
	rollupQuery := s.db.NewSelect().
		Model((*storage.Rollup)(nil)).
		ColumnExpr("id, name as value, 'rollup' as type").
		Where("name ILIKE ?", text)
	namespaceQuery := s.db.NewSelect().
		Model((*storage.Namespace)(nil)).
		ColumnExpr("id, encode(namespace_id, 'hex') as value, 'namespace' as type").
		Where("encode(namespace_id, 'hex') ILIKE ?", text)
	celestialsQuery := s.db.NewSelect().
		Model((*celestials.Celestial)(nil)).
		ColumnExpr("address_id as id, id as value, 'celestial' as type").
		Where("id ILIKE ?", text)

	union := rollupQuery.
		UnionAll(namespaceQuery).
		UnionAll(validatorQuery).
		UnionAll(celestialsQuery)

	err = s.db.NewSelect().
		TableExpr("(?) as search", union).
		Limit(10).
		Offset(0).
		Scan(ctx, &results)

	return
}
