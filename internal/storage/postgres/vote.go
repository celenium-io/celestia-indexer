// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Vote -
type Vote struct {
	*postgres.Table[*storage.Vote]
}

// NewVote -
func NewVote(db *database.Bun) *Vote {
	return &Vote{
		Table: postgres.NewTable[*storage.Vote](db),
	}
}

// ByProposalId -
func (v *Vote) ByProposalId(ctx context.Context, proposalId uint64, fltrs storage.VoteFilters) (votes []storage.Vote, err error) {
	subQuery := v.DB().NewSelect().
		Model((*storage.Vote)(nil)).
		Where("proposal_id = ?", proposalId)

	subQuery = limitScope(subQuery, fltrs.Limit)
	if fltrs.Offset > 0 {
		subQuery = subQuery.Offset(fltrs.Offset)
	}
	if len(fltrs.Option) > 0 {
		subQuery = subQuery.Where("option IN (?)", bun.In(fltrs.Option))
	}
	if fltrs.VoterType == types.VoterTypeValidator {
		subQuery = subQuery.Where("validator_id != 0")
	}
	if fltrs.VoterType == types.VoterTypeAddress {
		subQuery = subQuery.Where("validator_id = 0")
	}
	if fltrs.AddressId != nil {
		subQuery.Where("voter_id = ?", fltrs.AddressId)
	}
	if fltrs.ValidatorId != nil {
		subQuery.Where("validator_id = ?", fltrs.ValidatorId)
	}

	query := v.DB().NewSelect().
		TableExpr("(?) as votes", subQuery).
		ColumnExpr("votes.*").
		ColumnExpr("validator.id as validator__id").
		ColumnExpr("validator.cons_address as validator__cons_address").
		ColumnExpr("validator.moniker as validator__moniker").
		ColumnExpr("address.address as voter__address").
		ColumnExpr("celestial.id as voter__celestials__id").
		ColumnExpr("celestial.image_url as voter__celestials__image_url").
		OrderExpr("votes.time desc").
		Join("left join validator on validator.id = votes.validator_id").
		Join("left join address on address.id = votes.voter_id").
		Join("left join celestial on celestial.address_id = votes.voter_id and celestial.status = 'PRIMARY'")
	err = query.Scan(ctx, &votes)

	return
}

// ByVoterId -
func (v *Vote) ByVoterId(ctx context.Context, voterId uint64, fltrs storage.VoteFilters) (votes []storage.Vote, err error) {
	subQuery := v.DB().NewSelect().
		Model((*storage.Vote)(nil)).
		Where("voter_id = ?", voterId)

	subQuery = limitScope(subQuery, fltrs.Limit)
	if fltrs.Offset > 0 {
		subQuery = subQuery.Offset(fltrs.Offset)
	}

	query := v.DB().NewSelect().
		TableExpr("(?) as votes", subQuery).
		ColumnExpr("votes.*").
		ColumnExpr("validator.id as validator__id").
		ColumnExpr("validator.cons_address as validator__cons_address").
		ColumnExpr("validator.moniker as validator__moniker").
		ColumnExpr("proposal.id as proposal__id, proposal.status as proposal__status, proposal.title as proposal__title, proposal.description as proposal__description").
		OrderExpr("votes.time desc").
		Join("left join validator on validator.id = votes.validator_id").
		Join("left join proposal on proposal.id = votes.proposal_id")
	err = query.Scan(ctx, &votes)

	return
}

// ByValidatorId -
func (v *Vote) ByValidatorId(ctx context.Context, validatorId uint64, fltrs storage.VoteFilters) (votes []storage.Vote, err error) {
	subQuery := v.DB().NewSelect().
		Model((*storage.Vote)(nil)).
		Where("validator_id = ?", validatorId)

	subQuery = limitScope(subQuery, fltrs.Limit)
	if fltrs.Offset > 0 {
		subQuery = subQuery.Offset(fltrs.Offset)
	}

	query := v.DB().NewSelect().
		TableExpr("(?) as votes", subQuery).
		ColumnExpr("votes.*").
		ColumnExpr("address.address as voter__address").
		ColumnExpr("celestial.id as voter__celestials__id").
		ColumnExpr("celestial.image_url as voter__celestials__image_url").
		OrderExpr("votes.time desc").
		Join("left join address on address.id = votes.voter_id").
		Join("left join celestial on celestial.address_id = votes.voter_id and celestial.status = 'PRIMARY'")
	err = query.Scan(ctx, &votes)

	return
}
