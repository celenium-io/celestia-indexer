// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Proposal -
type Proposal struct {
	*postgres.Table[*storage.Proposal]
}

// NewProposal -
func NewProposal(db *database.Bun) *Proposal {
	return &Proposal{
		Table: postgres.NewTable[*storage.Proposal](db),
	}
}

func (p *Proposal) ByProposer(ctx context.Context, id uint64, limit, offset int) (proposals []storage.Proposal, err error) {
	query := p.DB().NewSelect().Model(&proposals).
		Where("proposer_id = ?", id).
		Order("id desc")

	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}
	err = query.Scan(ctx)
	return
}
