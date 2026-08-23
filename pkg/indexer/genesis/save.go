// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package genesis

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/pkg/errors"
)

var errCantFindAddressInContext = errors.New("can't find address in context")

func (module *Module) save(ctx context.Context, decodeCtx *decodeContext.Context) error {
	start := time.Now()
	module.Log.Info().Uint64("height", uint64(decodeCtx.Block.Height)).Msg("saving block...")
	tx, err := postgres.BeginTransaction(ctx, module.storage.Transactable)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	newConstants := make([]storage.Constant, 0, decodeCtx.Constants.Len())
	for constant := range decodeCtx.Constants.AllValues() {
		newConstants = append(newConstants, *constant)
	}
	if err := tx.SaveConstants(ctx, newConstants...); err != nil {
		return tx.HandleError(ctx, err)
	}

	for i := range decodeCtx.DenomMetadata {
		if err := tx.Add(ctx, &decodeCtx.DenomMetadata[i]); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if err := tx.Add(ctx, decodeCtx.Block); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Add(ctx, &decodeCtx.Block.Stats); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.SaveTransactions(ctx, decodeCtx.Block.Txs...); err != nil {
		return tx.HandleError(ctx, err)
	}

	var (
		messages   = make([]*storage.Message, 0, 1000)
		events     = make([]any, len(decodeCtx.Block.Events))
		namespaces = make(map[string]*storage.Namespace, 0)
	)

	for i := range decodeCtx.Block.Events {
		events[i] = &decodeCtx.Block.Events[i]
	}

	for i := range decodeCtx.Block.Txs {
		for j := range decodeCtx.Block.Txs[i].Messages {
			messages = append(messages, &decodeCtx.Block.Txs[i].Messages[j])

			for k := range decodeCtx.Block.Txs[i].Messages[j].Namespace {
				key := decodeCtx.Block.Txs[i].Messages[j].Namespace[k].String()
				if _, ok := namespaces[key]; !ok {
					decodeCtx.Block.Txs[i].Messages[j].Namespace[k].PfbCount = 1
					namespaces[key] = &decodeCtx.Block.Txs[i].Messages[j].Namespace[k]
				}
			}
		}

		for j := range decodeCtx.Block.Txs[i].Events {
			decodeCtx.Block.Txs[i].Events[j].TxId = &decodeCtx.Block.Txs[i].Id
			events = append(events, &decodeCtx.Block.Txs[i].Events[j])
		}
	}

	var totalAccounts int64
	if decodeCtx.Addresses.Len() > 0 {
		entities := decodeCtx.Addresses.Values()

		totalAccounts, err = tx.SaveAddresses(ctx, entities...)
		if err != nil {
			return tx.HandleError(ctx, err)
		}

		balances := make([]storage.Balance, 0, len(entities))
		for i := range entities {
			for j := range entities[i].Balances {
				entities[i].Balances[j].Id = entities[i].Id
			}
			balances = append(balances, entities[i].Balances...)
		}
		if err := tx.SaveBalances(ctx, balances...); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	var totalNamespaces int64
	if len(namespaces) > 0 {
		entities := make([]*storage.Namespace, 0, len(namespaces))
		for key := range namespaces {
			entities = append(entities, namespaces[key])
		}

		totalNamespaces, err = tx.SaveNamespaces(ctx, entities...)
		if err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	validators := decodeCtx.Validators.Values()
	totalValidators, err := tx.SaveValidators(ctx, validators...)
	if err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.SaveMessages(ctx, messages...); err != nil {
		return tx.HandleError(ctx, err)
	}

	var msgVals []storage.MsgValidator
	for i := range messages {
		for _, val := range messages[i].Validators {
			for j := range validators {
				if validators[j].Address == val {
					msgVals = append(msgVals, storage.MsgValidator{
						Height:      0,
						Time:        decodeCtx.Block.Time,
						MsgId:       messages[i].Id,
						ValidatorId: validators[j].Id,
					})
					break
				}
			}
		}
	}

	if err := tx.SaveMsgValidator(ctx, msgVals...); err != nil {
		return tx.HandleError(ctx, err)
	}

	for i := range decodeCtx.StakingLogs {
		if address, ok := decodeCtx.Addresses.Get(decodeCtx.StakingLogs[i].Address.Address); ok {
			decodeCtx.StakingLogs[i].AddressId = &address.Id
		} else {
			return tx.HandleError(ctx, errors.Wrap(errCantFindAddressInContext, decodeCtx.StakingLogs[i].Address.Address))
		}

		for j := range validators {
			if validators[j].Address == decodeCtx.StakingLogs[i].Validator.Address {
				decodeCtx.StakingLogs[i].ValidatorId = validators[j].Id
				break
			}
		}
	}

	if err := tx.SaveStakingLogs(ctx, decodeCtx.StakingLogs...); err != nil {
		return tx.HandleError(ctx, err)
	}

	for i := range decodeCtx.VestingAccounts {
		if decodeCtx.VestingAccounts[i] == nil || decodeCtx.VestingAccounts[i].Address == nil {
			return tx.HandleError(ctx, errors.New("nil pointer for vesting"))
		}
		if address, ok := decodeCtx.Addresses.Get(decodeCtx.VestingAccounts[i].Address.Address); ok {
			decodeCtx.VestingAccounts[i].AddressId = address.Id
		} else {
			return tx.HandleError(ctx, errors.Wrap(errCantFindAddressInContext, decodeCtx.VestingAccounts[i].Address.Address))
		}
	}

	if err := tx.SaveVestingAccounts(ctx, decodeCtx.VestingAccounts...); err != nil {
		return tx.HandleError(ctx, err)
	}

	periods := make([]storage.VestingPeriod, 0)
	for i := range decodeCtx.VestingAccounts {
		for j := range decodeCtx.VestingAccounts[i].VestingPeriods {
			decodeCtx.VestingAccounts[i].VestingPeriods[j].VestingAccountId = decodeCtx.VestingAccounts[i].Id
		}
		periods = append(periods, decodeCtx.VestingAccounts[i].VestingPeriods...)
	}
	if err := tx.SaveVestingPeriods(ctx, periods...); err != nil {
		return tx.HandleError(ctx, err)
	}

	delegations := make([]storage.Delegation, 0, decodeCtx.Delegations.Len())
	for delegation := range decodeCtx.Delegations.AllValues() {
		if address, ok := decodeCtx.Addresses.Get(delegation.Address.Address); ok {
			delegation.AddressId = address.Id
		} else {
			return tx.HandleError(ctx, errors.Wrap(errCantFindAddressInContext, delegation.Address.Address))
		}

		for j := range validators {
			if validators[j].Address == delegation.Validator.Address {
				delegation.ValidatorId = validators[j].Id
				break
			}
		}

		delegations = append(delegations, *delegation)
	}
	if err := tx.SaveDelegations(ctx, delegations...); err != nil {
		return tx.HandleError(ctx, err)
	}

	if len(events) > 0 {
		if err := tx.BulkSave(ctx, events); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	var namespaceMsgs []*storage.NamespaceMessage
	for i := range messages {
		for j := range messages[i].Namespace {
			if messages[i].Namespace[j].Id == 0 { // in case of duplication of writing to one namespace inside one messages
				continue
			}
			namespaceMsgs = append(namespaceMsgs, &storage.NamespaceMessage{
				MsgId:       messages[i].Id,
				NamespaceId: messages[i].Namespace[j].Id,
				Time:        messages[i].Time,
				Height:      messages[i].Height,
				TxId:        messages[i].TxId,
				Size:        uint64(messages[i].Namespace[j].Size),
			})
		}
	}
	if err := tx.SaveNamespaceMessage(ctx, namespaceMsgs...); err != nil {
		return tx.HandleError(ctx, err)
	}

	var signers []storage.Signer
	for i := range decodeCtx.Block.Txs {
		for _, address := range decodeCtx.Block.Txs[i].Signers {
			signer, ok := decodeCtx.Addresses.Get(address.Address)
			if !ok {
				return tx.HandleError(ctx, errors.Wrap(errCantFindAddressInContext, address.Address))
			}
			signers = append(signers, storage.Signer{
				TxId:      decodeCtx.Block.Txs[i].Id,
				AddressId: signer.Id,
			})
		}
	}

	if err := tx.SaveSigners(ctx, signers...); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Add(ctx, &storage.State{
		Name:            module.indexerName,
		LastHeight:      decodeCtx.Block.Height,
		LastTime:        decodeCtx.Block.Time,
		LastHash:        decodeCtx.Block.Hash,
		ChainId:         decodeCtx.Block.ChainId,
		TotalTx:         decodeCtx.Block.Stats.TxCount,
		TotalSupply:     decodeCtx.Block.Stats.SupplyChange,
		TotalFee:        decodeCtx.Block.Stats.Fee,
		TotalBlobsSize:  decodeCtx.Block.Stats.BlobsSize,
		TotalAccounts:   totalAccounts,
		TotalNamespaces: totalNamespaces,
		TotalValidators: totalValidators,
	}); err != nil {
		return tx.HandleError(ctx, err)
	}

	grants := make([]*storage.Grant, 0, decodeCtx.Grants.Len())
	for grant := range decodeCtx.Grants.AllValues() {
		if address, ok := decodeCtx.Addresses.Get(grant.Grantee.Address); ok {
			grant.GranteeId = address.Id
		} else {
			return tx.HandleError(ctx, errors.Wrap(errCantFindAddressInContext, grant.Grantee.Address))
		}
		if address, ok := decodeCtx.Addresses.Get(grant.Granter.Address); ok {
			grant.GranterId = address.Id
		} else {
			return tx.HandleError(ctx, errors.Wrap(errCantFindAddressInContext, grant.Granter.Address))
		}
		grants = append(grants, grant)
	}

	if err := tx.SaveGrants(ctx, grants...); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Flush(ctx); err != nil {
		return tx.HandleError(ctx, err)
	}
	module.Log.Info().
		Uint64("height", decodeCtx.Block.Id).
		Int64("block_ns_size", decodeCtx.Block.Stats.BlobsSize).
		Str("block_fee", decodeCtx.Block.Stats.Fee.String()).
		Int64("ms", time.Since(start).Milliseconds()).
		Msg("block saved")
	return nil
}
