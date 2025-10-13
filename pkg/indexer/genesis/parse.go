// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package genesis

import (
	"strings"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/node/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type parsedData struct {
	block         storage.Block
	addresses     map[string]*storage.Address
	validators    []*storage.Validator
	stakingLogs   []storage.StakingLog
	delegations   []storage.Delegation
	constants     []storage.Constant
	denomMetadata []storage.DenomMetadata
	vestings      []*storage.VestingAccount

	bondedTokensPool *storage.Address
}

func newParsedData() parsedData {
	return parsedData{
		addresses:     make(map[string]*storage.Address),
		validators:    make([]*storage.Validator, 0),
		stakingLogs:   make([]storage.StakingLog, 0),
		delegations:   make([]storage.Delegation, 0),
		constants:     make([]storage.Constant, 0),
		denomMetadata: make([]storage.DenomMetadata, 0),
		vestings:      make([]*storage.VestingAccount, 0),
	}
}

func (module *Module) parse(genesis types.GenesisOutput) (parsedData, error) {
	data := newParsedData()
	block := storage.Block{
		Time:    genesis.GenesisTime,
		Height:  pkgTypes.Level(genesis.InitialHeight - 1),
		AppHash: []byte(genesis.AppHash),
		ChainId: genesis.ChainID,
		Txs:     make([]storage.Tx, 0),
		Stats: storage.BlockStats{
			Time:          genesis.GenesisTime,
			Height:        pkgTypes.Level(genesis.InitialHeight - 1),
			TxCount:       int64(len(genesis.AppState.Genutil.GenTxs)),
			EventsCount:   0,
			Fee:           decimal.Zero,
			SupplyChange:  decimal.Zero,
			InflationRate: decimal.Zero,
		},
		MessageTypes: storageTypes.NewMsgTypeBits(),
	}

	decodeCtx := context.NewContext()
	decodeCtx.Block = &block

	if err := module.parseAccounts(genesis.ModuleAccs, block, &data); err != nil {
		return data, errors.Wrap(err, "parse genesis accounts")
	}

	for index, genTx := range genesis.AppState.Genutil.GenTxs {
		txDecoded, err := decode.JsonTx(genTx)
		if err != nil {
			return data, errors.Wrapf(err, "failed to decode GenTx '%s'", genTx)
		}

		memoTx, ok := txDecoded.(cosmosTypes.TxWithMemo)
		if !ok {
			return data, errors.Wrapf(err, "expected TxWithMemo, got %T", genTx)
		}
		txWithTimeoutHeight, ok := txDecoded.(cosmosTypes.TxWithTimeoutHeight)
		if !ok {
			return data, errors.Wrapf(err, "expected TxWithTimeoutHeight, got %T", genTx)
		}

		tx := storage.Tx{
			Height:        block.Height,
			Time:          block.Time,
			Position:      int64(index),
			TimeoutHeight: txWithTimeoutHeight.GetTimeoutHeight(),
			MessagesCount: int64(len(txDecoded.GetMsgs())),
			Fee:           decimal.Zero,
			Status:        storageTypes.StatusSuccess,
			Memo:          memoTx.GetMemo(),
			MessageTypes:  storageTypes.NewMsgTypeBitMask(),

			Messages: make([]storage.Message, len(txDecoded.GetMsgs())),
			Events:   nil,
		}

		for msgIndex, msg := range txDecoded.GetMsgs() {
			decoded, err := decode.Message(decodeCtx, msg, msgIndex, storageTypes.StatusSuccess)
			if err != nil {
				return data, errors.Wrap(err, "decode genesis message")
			}

			tx.Messages[msgIndex] = decoded.Msg
			tx.MessageTypes.SetByMsgType(decoded.Msg.Type)
			block.MessageTypes.SetByMsgType(decoded.Msg.Type)
			tx.BlobsSize += decoded.BlobsSize
		}

		block.Txs = append(block.Txs, tx)
	}

	for _, addr := range decodeCtx.Addresses.Values() {
		data.addresses[addr.String()] = addr
	}

	module.parseDenomMetadata(genesis.AppState.Bank.DenomMetadata, &data)
	if err := module.parseConstants(genesis.AppState, genesis.ConsensusParams, &data); err != nil {
		return data, errors.Wrap(err, "parse constants")
	}

	module.parseTotalSupply(genesis.AppState.Bank.Supply, &block)

	if err := module.parseAccounts(genesis.AppState.Auth.Accounts, block, &data); err != nil {
		return data, errors.Wrap(err, "parse genesis accounts")
	}
	if err := module.parseBalances(genesis.AppState.Bank.Balances, block.Height, &data); err != nil {
		return data, errors.Wrap(err, "parse genesis account balances")
	}

	data.validators = decodeCtx.Validators.Values()
	data.stakingLogs = decodeCtx.StakingLogs

	_ = decodeCtx.Delegations.Range(func(_ string, value *storage.Delegation) (error, bool) {
		data.delegations = append(data.delegations, *value)
		if data.bondedTokensPool != nil {
			data.bondedTokensPool.Balance.Spendable = data.bondedTokensPool.Balance.Spendable.Add(value.Amount)
		}
		if addr, ok := data.addresses[value.Address.Address]; ok {
			addr.Balance.Spendable = addr.Balance.Spendable.Sub(value.Amount)
		}
		return nil, false
	})

	data.block = block
	return data, nil
}

func (module *Module) parseTotalSupply(supply []types.Supply, block *storage.Block) {
	if len(supply) == 0 {
		return
	}

	if totalSupply, err := decimal.NewFromString(supply[0].Amount); err == nil {
		block.Stats.SupplyChange = totalSupply
	}
}

func (module *Module) parseAccounts(accounts []types.Account, block storage.Block, data *parsedData) error {
	currencyBase := currency.DefaultCurrency
	if len(data.denomMetadata) > 0 {
		currencyBase = data.denomMetadata[0].Base
	}

	for i := range accounts {
		address := storage.Address{
			Height:     block.Height,
			LastHeight: block.Height,
			Balance: storage.Balance{
				Spendable: decimal.Zero,
				Delegated: decimal.Zero,
				Unbonding: decimal.Zero,
				Currency:  currencyBase,
			},
		}

		var readableAddress string

		switch {
		case strings.Contains(accounts[i].Type, "PeriodicVestingAccount"):
			readableAddress = accounts[i].BaseVestingAccount.BaseAccount.Address
			if err := parseVesting(accounts[i], block, readableAddress, storageTypes.VestingTypePeriodic, data); err != nil {
				return err
			}

		case strings.Contains(accounts[i].Type, "ModuleAccount"):
			readableAddress = accounts[i].BaseAccount.Address
			address.Name = accounts[i].Name

			if address.Name == "bonded_tokens_pool" {
				data.bondedTokensPool = &address
			}

		case strings.Contains(accounts[i].Type, "BaseAccount"):
			readableAddress = accounts[i].Address

		case strings.Contains(accounts[i].Type, "ContinuousVestingAccount"):
			readableAddress = accounts[i].BaseVestingAccount.BaseAccount.Address
			if err := parseVesting(accounts[i], block, readableAddress, storageTypes.VestingTypeContinuous, data); err != nil {
				return err
			}

		case strings.Contains(accounts[i].Type, "DelayedVestingAccount"):
			readableAddress = accounts[i].BaseVestingAccount.BaseAccount.Address
			if err := parseVesting(accounts[i], block, readableAddress, storageTypes.VestingTypeDelayed, data); err != nil {
				return err
			}

		default:
			return errors.Errorf("unknown account type: %s", accounts[i].Type)
		}

		if _, ok := data.addresses[readableAddress]; !ok {
			_, hash, err := pkgTypes.Address(readableAddress).Decode()
			if err != nil {
				return err
			}
			address.Hash = hash
			address.Address = readableAddress
			data.addresses[address.String()] = &address
		}
	}
	return nil
}

func (module *Module) parseBalances(balances []types.Balances, height pkgTypes.Level, data *parsedData) error {
	for i := range balances {
		if len(balances[i].Coins) == 0 {
			continue
		}

		_, hash, err := pkgTypes.Address(balances[i].Address).Decode()
		if err != nil {
			return err
		}
		address := storage.Address{
			Hash:       hash,
			Address:    balances[i].Address,
			Height:     height,
			LastHeight: height,
			Balance: storage.Balance{
				Spendable: decimal.Zero,
				Delegated: decimal.Zero,
				Unbonding: decimal.Zero,
				Currency:  balances[i].Coins[0].Denom,
			},
		}
		if balance, err := decimal.NewFromString(balances[i].Coins[0].Amount); err == nil {
			address.Balance.Spendable = address.Balance.Spendable.Add(balance)
		}

		if addr, ok := data.addresses[address.String()]; ok {
			addr.Balance.Spendable = addr.Balance.Spendable.Add(address.Balance.Spendable)
		} else {
			data.addresses[address.String()] = &address
		}
	}

	return nil
}

func getAmountFromOriginalVesting(vestings []types.Coins) (decimal.Decimal, error) {
	var amount = decimal.Zero.Copy()

	for i := range vestings {
		val, err := decimal.NewFromString(vestings[i].Amount)
		if err != nil {
			return amount, err
		}
		amount = amount.Add(val)
	}

	return amount, nil
}

func parseVesting(acc types.Account, block storage.Block, address string, typ storageTypes.VestingType, data *parsedData) error {
	amount, err := getAmountFromOriginalVesting(acc.BaseVestingAccount.OriginalVesting)
	if err != nil {
		return err
	}

	v := storage.VestingAccount{
		Height: block.Height,
		Time:   block.Time,
		Address: &storage.Address{
			Address: address,
		},
		Type:           typ,
		Amount:         amount,
		VestingPeriods: make([]storage.VestingPeriod, 0),
	}

	if acc.BaseVestingAccount.EndTime > 0 {
		t := time.Unix(acc.BaseVestingAccount.EndTime, 0).UTC()
		v.EndTime = &t
	}

	var periodTime = v.Time
	if acc.StartTime != nil {
		t := time.Unix(*acc.StartTime, 0).UTC()
		v.StartTime = &t
		periodTime = t
	}

	for i := range acc.VestingPeriods {
		period := storage.VestingPeriod{
			Height: v.Height,
		}
		amount, err := getAmountFromOriginalVesting(acc.VestingPeriods[i].Amount)
		if err != nil {
			return err
		}
		period.Amount = amount
		periodTime = periodTime.Add(time.Second * time.Duration(acc.VestingPeriods[i].Length))
		period.Time = periodTime
		v.VestingPeriods = append(v.VestingPeriods, period)
	}

	data.vestings = append(data.vestings, &v)
	return nil
}
