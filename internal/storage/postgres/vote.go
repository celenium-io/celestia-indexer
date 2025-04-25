// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"github.com/celenium-io/celestia-indexer/internal/storage"
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
		Relation("Voter").
		Relation("Validator").
		Order("time desc")

	query = limitScope(query, fltrs.Limit)
	if fltrs.Offset > 0 {
		query = query.Offset(fltrs.Offset)
	}

	if fltrs.Option != "" {
		query = query.Where("option = ?", fltrs.Option)
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
		Relation("Voter").
		Relation("Validator").
		Order("time desc")

	query = limitScope(query, fltrs.Limit)
	if fltrs.Offset > 0 {
		query = query.Offset(fltrs.Offset)
	}

	err = query.Scan(ctx)
	return
}

// ByValidatorId -
func (v *Vote) ByValidatorId(ctx context.Context, validatorId uint64, fltrs storage.VoteFilters) (votes []storage.Vote, err error) {
	query := v.DB().NewSelect().
		Model(&votes).
		Where("validator_id = ?", validatorId).
		Relation("Voter").
		Relation("Validator").
		Order("time desc")

	query = limitScope(query, fltrs.Limit)
	if fltrs.Offset > 0 {
		query = query.Offset(fltrs.Offset)
	}

	err = query.Scan(ctx)
	return
}
