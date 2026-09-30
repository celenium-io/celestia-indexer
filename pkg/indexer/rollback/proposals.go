// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rollback

import (
	"context"
	"strconv"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/pkg/errors"
)

func rollbackProposals(
	ctx context.Context,
	tx storage.GovTx,
	constants storage.IConstant,
	proposalsRepo storage.IProposal,
	deletedEvents []storage.Event,
) error {
	activeProposals := make([]*storage.Proposal, 0)
	for i := range deletedEvents {
		if deletedEvents[i].Type != storageTypes.EventTypeActiveProposal {
			continue
		}

		status, err := decode.NewProposalStatus(deletedEvents[i].Data)
		if err != nil {
			return err
		}

		if status.Id > 0 && status.Result != "" {
			if status.Result == "expedited_proposal_rejected" {
				votingPeriodConst, err := constants.Get(ctx, storageTypes.ModuleNameGov, "expedited_voting_period")
				if err != nil {
					return err
				}
				votingPeriod, err := strconv.ParseInt(votingPeriodConst.Value, 10, 64)
				if err != nil {
					return err
				}
				proposal, err := proposalsRepo.GetByID(ctx, status.Id)
				if err != nil {
					return err
				}
				if proposal.ActivationTime == nil {
					return errors.Errorf("nil activation time for proposal: %d", status.Id)
				}

				endTime := proposal.ActivationTime.Add(time.Duration(votingPeriod))
				activeProposals = append(activeProposals, &storage.Proposal{
					Id:        status.Id,
					Expedited: new(true),
					EndTime:   &endTime,
				})
			}
		}
	}

	if len(activeProposals) == 0 {
		return nil
	}
	_, err := tx.SaveProposals(ctx, activeProposals...)
	return err
}
