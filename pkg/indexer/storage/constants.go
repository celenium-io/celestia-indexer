// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/bytedance/sonic"
	sdkSync "github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	paramsV1Beta "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
)

func saveConstantUpdates(
	ctx context.Context,
	tx storage.BlockTx,
	consts *sdkSync.Map[string, *storage.Constant],
) error {
	if consts.Len() == 0 {
		return nil
	}

	newConstants := make([]storage.Constant, 0, consts.Len())
	for value := range consts.AllValues() {
		newConstants = append(newConstants, *value)
	}
	return tx.SaveConstants(ctx, newConstants...)
}

func saveProposalConstantUpdates(
	ctx context.Context,
	tx storage.BlockTx,
	proposalsRepo storage.IProposal,
	proposals *sdkSync.Map[uint64, *storage.Proposal],
) error {
	consts := make(map[string]storage.Constant)
	for p := range appliedProposals(proposals.AllValues()) {
		proposal, err := proposalsRepo.GetByID(ctx, p.Id)
		if err != nil {
			return errors.Wrapf(err, "receiving proposal %d", p.Id)
		}
		if proposal.Type != types.ProposalTypeParamChanged {
			continue
		}
		if err := parseParamChanges(consts, proposal.Changes); err != nil {
			return errors.Wrapf(err, "parse proposal changes %d", p.Id)
		}
	}
	if len(consts) == 0 {
		return nil
	}
	return tx.SaveConstants(ctx, slices.Collect(maps.Values(consts))...)
}

func parseParamChanges(consts map[string]storage.Constant, changesBytes []byte) error {
	if len(changesBytes) == 0 {
		return nil
	}
	var changes []paramsV1Beta.ParamChange
	if err := sonic.Unmarshal(changesBytes, &changes); err != nil {
		return err
	}
	for i := range changes {
		module, err := types.ParseModuleName(changes[i].Subspace)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("%s_%s", changes[i].Subspace, changes[i].Key)
		consts[key] = storage.Constant{
			Module: module,
			Name:   changes[i].Key,
			Value:  changes[i].Value,
		}
	}
	return nil
}
