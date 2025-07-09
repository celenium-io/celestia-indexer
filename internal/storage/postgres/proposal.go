// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
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

func (p *Proposal) ListWithFilters(ctx context.Context, filters storage.ListProposalFilters) (proposals []storage.Proposal, err error) {
	query := p.DB().NewSelect().Model((*storage.Proposal)(nil))
	query = limitScope(query, filters.Limit)
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}
	query = sortScope(query, "id", filters.Sort)
	if filters.ProposerId > 0 {
		query = query.Where("proposer_id = ?", filters.ProposerId)
	}

	if len(filters.Status) > 0 {
		query = query.Where("status IN (?)", bun.In(filters.Status))
	}
	if len(filters.Type) > 0 {
		query = query.Where("type IN (?)", bun.In(filters.Type))
	}

	err = p.DB().NewSelect().
		ColumnExpr("proposals.*").
		ColumnExpr("proposer.address as proposer__address").
		ColumnExpr("celestial.id as proposer__celestials__id, celestial.image_url as proposer__celestials__image_url").
		TableExpr("(?) as proposals", query).
		Join("left join address as proposer ON proposals.proposer_id = proposer.id").
		Join("left join celestial on celestial.address_id = proposals.proposer_id and celestial.status = 'PRIMARY'").
		Scan(ctx, &proposals)

	return
}

func (p *Proposal) ById(ctx context.Context, id uint64) (proposal storage.Proposal, err error) {
	err = p.DB().NewSelect().
		Model(&proposal).
		ColumnExpr("proposal.*").
		ColumnExpr("proposer.address as proposer__address").
		ColumnExpr("celestial.id as proposer__celestials__id, celestial.image_url as proposer__celestials__image_url").
		Where("proposal.id = ?", id).
		Join("left join address as proposer ON proposal.proposer_id = proposer.id").
		Join("left join celestial on celestial.address_id = proposal.proposer_id and celestial.status = 'PRIMARY'").
		Scan(ctx)
	return
}
