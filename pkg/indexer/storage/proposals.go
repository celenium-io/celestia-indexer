// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/pkg/errors"
)

func (module *Module) saveProposals(
	ctx context.Context,
	tx storage.Transaction,
	proposals []*storage.Proposal,
	votes []*storage.Vote,
	addrToId map[string]uint64,
) error {
	if len(votes) > 0 {
		for i := range votes {
			if votes[i].Voter != nil {
				voterId, ok := addrToId[votes[i].Voter.Address]
				if !ok {
					return errors.Errorf("unknown voter address: %s", votes[i].Voter.Address)
				}
				votes[i].VoterId = voterId
			}
		}

		if err := tx.SaveVotes(ctx, votes...); err != nil {
			return errors.Wrap(err, "save votes")
		}
	}

	for i := range proposals {
		if proposals[i].Proposer != nil {
			proposerId, ok := addrToId[proposals[i].Proposer.Address]
			if !ok {
				return errors.Errorf("unknown proposer address for proposal: %s", proposals[i].Proposer.Address)
			}
			proposals[i].ProposerId = proposerId
		}

		if !proposals[i].CreatedAt.IsZero() {
			constant, err := module.constants.Get(ctx, types.ModuleNameGov, "max_deposit_period")
			if err != nil {
				return errors.Wrap(err, "can't find max_deposit_period constant")
			}
			maxDepositPeriod, err := time.ParseDuration(constant.Value)
			if err != nil {
				return errors.Wrap(err, "can't parse max_deposit_period value")
			}
			proposals[i].DepositTime = proposals[i].CreatedAt.Add(maxDepositPeriod)
		}
	}

	return tx.SaveProposals(ctx, proposals...)
}
