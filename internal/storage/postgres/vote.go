// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
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
	query := v.DB().NewSelect().
		Model(&votes).
		Where("proposal_id = ?", proposalId).
		ColumnExpr("vote.*").
		ColumnExpr("validator.id as validator__id").
		ColumnExpr("validator.cons_address as validator__cons_address").
		ColumnExpr("validator.moniker as validator__moniker").
		ColumnExpr("address.address as voter__address").
		ColumnExpr("celestial.id as voter__celestials__id").
		ColumnExpr("celestial.image_url as voter__celestials__image_url").
		Join("left join validator on validator.id = vote.validator_id").
		Join("left join address on address.id = vote.voter_id").
		Join("left join celestial on celestial.address_id = vote.voter_id and celestial.status = 'PRIMARY'")

	query = limitScope(query, fltrs.Limit)
	if fltrs.Offset > 0 {
		query = query.Offset(fltrs.Offset)
	}

	if fltrs.Option != "" {
		query = query.Where("option = ?", fltrs.Option)
	}

	if fltrs.VoterType == types.VoterTypeValidator {
		query = query.Where("validator_id = ?", 0)
	}
	if fltrs.VoterType == types.VoterTypeAddress {
		query = query.Where("validator_id != ?", 0)
	}

	query = query.Order("time desc")
	err = query.Scan(ctx)

	return
}

// ByVoterId -
func (v *Vote) ByVoterId(ctx context.Context, voterId uint64, fltrs storage.VoteFilters) (votes []storage.Vote, err error) {
	query := v.DB().NewSelect().
		Model(&votes).
		Where("voter_id = ?", voterId).
		ColumnExpr("vote.*").
		ColumnExpr("validator.id as validator__id").
		ColumnExpr("validator.cons_address as validator__cons_address").
		ColumnExpr("validator.moniker as validator__moniker").
		ColumnExpr("address.address as voter__address").
		ColumnExpr("celestial.id as voter__celestials__id").
		ColumnExpr("celestial.image_url as voter__celestials__image_url").
		Join("left join validator on validator.id = vote.validator_id").
		Join("left join address on address.id = vote.voter_id").
		Join("left join celestial on celestial.address_id = vote.voter_id and celestial.status = 'PRIMARY'")

	query = limitScope(query, fltrs.Limit)
	if fltrs.Offset > 0 {
		query = query.Offset(fltrs.Offset)
	}

	query = query.Order("time desc")
	err = query.Scan(ctx)

	return
}

// ByValidatorId -
func (v *Vote) ByValidatorId(ctx context.Context, validatorId uint64, fltrs storage.VoteFilters) (votes []storage.Vote, err error) {
	query := v.DB().NewSelect().
		Model(&votes).
		Where("validator_id = ?", validatorId).
		ColumnExpr("vote.*").
		ColumnExpr("validator.id as validator__id").
		ColumnExpr("validator.cons_address as validator__cons_address").
		ColumnExpr("validator.moniker as validator__moniker").
		ColumnExpr("address.address as voter__address").
		ColumnExpr("celestial.id as voter__celestials__id").
		ColumnExpr("celestial.image_url as voter__celestials__image_url").
		Join("left join validator on validator.id = vote.validator_id").
		Join("left join address on address.id = vote.voter_id").
		Join("left join celestial on celestial.address_id = vote.voter_id and celestial.status = 'PRIMARY'")

	query = limitScope(query, fltrs.Limit)
	if fltrs.Offset > 0 {
		query = query.Offset(fltrs.Offset)
	}

	query = query.Order("time desc")
	err = query.Scan(ctx)

	return
}
