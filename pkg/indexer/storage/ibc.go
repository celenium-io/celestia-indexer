// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"slices"
	"time"

	json "github.com/bytedance/sonic"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	sdkSync "github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
)

func saveIbcClients(
	ctx context.Context,
	tx storage.Transaction,
	clients []*storage.IbcClient,
	addrToId map[string]uint64,
) (int64, error) {
	if len(clients) == 0 {
		return 0, nil
	}

	for i := range clients {
		if clients[i].Creator != nil {
			if addrId, ok := addrToId[clients[i].Creator.Address]; ok {
				clients[i].CreatorId = addrId
			} else {
				return 0, errors.Wrap(errCantFindAddress, clients[i].Creator.Address)
			}
		}
	}

	return tx.SaveIbcClients(ctx, clients...)
}

type clientRecovery struct {
	SubjectClientId    string
	SubstituteClientId string
}

// parseClientRecoveries reads client_update proposal changes: a list of subject/substitute pairs
func parseClientRecoveries(changes []byte) ([]clientRecovery, error) {
	if len(changes) == 0 {
		return nil, nil
	}
	var list []clientRecovery
	if err := json.Unmarshal(changes, &list); err != nil {
		return nil, err
	}
	// other shapes, such as param changes of a mixed proposal, carry no subject
	return slices.DeleteFunc(list, func(r clientRecovery) bool {
		return r.SubjectClientId == ""
	}), nil
}

// recoverIbcClients applies governance client recoveries: the substitute comes from the applied proposal
func (module *Module) recoverIbcClients(
	ctx context.Context,
	tx storage.Transaction,
	subjects []string,
	proposals *sdkSync.Map[uint64, *storage.Proposal],
	blockTime time.Time,
) error {
	if len(subjects) == 0 {
		return nil
	}

	substitutes := make(map[string]string)
	for _, p := range proposals.Values() {
		if p.Status != types.ProposalStatusApplied {
			continue
		}
		proposal, err := tx.Proposal(ctx, p.Id)
		if err != nil {
			return errors.Wrapf(err, "receiving proposal %d", p.Id)
		}
		if proposal.Type != types.ProposalTypeClientUpdate {
			continue
		}
		recoveries, err := parseClientRecoveries(proposal.Changes)
		if err != nil {
			// not worth halting the indexer: subjects of this proposal are still unfrozen below
			module.Log.Warn().Err(err).Uint64("proposal_id", p.Id).Msg("can't parse client recoveries")
			continue
		}
		for _, r := range recoveries {
			substitutes[r.SubjectClientId] = r.SubstituteClientId
		}
	}

	for _, subject := range subjects {
		substitute, ok := substitutes[subject]
		if !ok {
			// e.g. proposal submitted before recoveries were stored: at least unfreeze
			module.Log.Warn().Str("client_id", subject).Msg("substitute of recovered IBC client not found")
		}
		if err := tx.RecoverIbcClient(ctx, subject, substitute, blockTime); err != nil {
			return errors.Wrapf(err, "recovering IBC client %s", subject)
		}
	}
	return nil
}

func saveIbcChannels(
	ctx context.Context,
	tx storage.Transaction,
	channels []*storage.IbcChannel,
	addrToId map[string]uint64,
) error {
	if len(channels) == 0 {
		return nil
	}

	for i := range channels {
		if channels[i].ConnectionId != "" {
			conn, err := tx.IbcConnection(ctx, channels[i].ConnectionId)
			if err != nil {
				return errors.Wrap(err, "receiving connection for channel")
			}
			channels[i].ClientId = conn.ClientId
		}

		if channels[i].Creator != nil {
			if addrId, ok := addrToId[channels[i].Creator.Address]; ok {
				channels[i].CreatorId = addrId
			} else {
				return errors.Wrap(errCantFindAddress, channels[i].Creator.Address)
			}
		}
	}

	return tx.SaveIbcChannels(ctx, channels...)
}

func saveIbcTransfers(
	ctx context.Context,
	tx storage.Transaction,
	transfers []*storage.IbcTransfer,
	addrToId map[string]uint64,
) error {
	if len(transfers) == 0 {
		return nil
	}

	for i := range transfers {
		if transfers[i].Sender != nil {
			if addrId, ok := addrToId[transfers[i].Sender.Address]; ok {
				transfers[i].SenderId = &addrId
			} else {
				return errors.Wrap(errCantFindAddress, transfers[i].Sender.Address)
			}
		}
		if transfers[i].Receiver != nil {
			if addrId, ok := addrToId[transfers[i].Receiver.Address]; ok {
				transfers[i].ReceiverId = &addrId
			} else {
				return errors.Wrap(errCantFindAddress, transfers[i].Receiver.Address)
			}
		}
	}

	return tx.SaveIbcTransfers(ctx, transfers...)
}
