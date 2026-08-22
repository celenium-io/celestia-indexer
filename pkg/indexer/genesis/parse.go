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
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/pkg/errors"
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
	grants        []*storage.Grant

	bondedTokensPool *storage.Address
	feeCollector     *storage.Address
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
		grants:        make([]*storage.Grant, 0),
	}
}

func (module *Module) parse(genesis types.GenesisOutput) (parsedData, error) {
	if genesis.AppState.Staking.Exported {
		return parsedData{}, errors.New("genesis was produced by state export (app_state.staking.exported=true): validators and delegations are not read from gen_txs in this case and parsing them is not supported")
	}

	data := newParsedData()
	block := storage.Block{
		Time:    genesis.GenesisTime,
		Height:  pkgTypes.Level(genesis.InitialHeight - 1),
		AppHash: []byte(genesis.AppHash),
		ChainId: genesis.ChainID,
		Txs:     make([]storage.Tx, len(genesis.AppState.Genutil.GenTxs)),
		Stats: storage.BlockStats{
			Time:          genesis.GenesisTime,
			Height:        pkgTypes.Level(genesis.InitialHeight - 1),
			TxCount:       int64(len(genesis.AppState.Genutil.GenTxs)),
			EventsCount:   0,
			Fee:           storageTypes.NumericZero(),
			SupplyChange:  storageTypes.NumericZero(),
			InflationRate: storageTypes.NumericZero(),
		},
		MessageTypes: storageTypes.NewMsgTypeBits(),
	}
	module.parseDenomMetadata(genesis.AppState.Bank.DenomMetadata, &data)

	decodeCtx := context.NewContext()
	decodeCtx.Block = &block
	currencyBase := getBaseCurrency(data.denomMetadata)

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

		tx := &block.Txs[index]
		tx.Height = block.Height
		tx.Time = block.Time
		tx.Position = int64(index)
		tx.TimeoutHeight = txWithTimeoutHeight.GetTimeoutHeight()
		tx.MessagesCount = int64(len(txDecoded.GetMsgs()))
		tx.Fee = storageTypes.NumericZero()
		tx.Status = storageTypes.StatusSuccess
		tx.Memo = memoTx.GetMemo()
		tx.MessageTypes = storageTypes.NewMsgTypeBitMask()
		tx.Messages = make([]storage.Message, len(txDecoded.GetMsgs()))
		tx.Events = nil

		if err := tx.SetId(); err != nil {
			return data, err
		}

		for msgIndex, msg := range txDecoded.GetMsgs() {
			decoded, err := decode.Message(decodeCtx, msg, msgIndex, storageTypes.StatusSuccess, tx.Id)
			if err != nil {
				return data, errors.Wrap(err, "decode genesis message")
			}

			tx.Messages[msgIndex] = decoded.Msg
			tx.MessageTypes.SetByMsgType(decoded.Msg.Type)
			block.MessageTypes.SetByMsgType(decoded.Msg.Type)
			tx.BlobsSize += decoded.BlobsSize
		}

		if t, ok := txDecoded.(cosmosTypes.FeeTx); ok {
			value, err := decode.DecodeFee(t.GetFee())
			if err != nil {
				return data, errors.Wrap(err, "fee decoding")
			}
			tx.Fee = storageTypes.NewNumeric(value)
			block.Stats.Fee = block.Stats.Fee.Add(tx.Fee)
			if data.feeCollector != nil {
				data.feeCollector.AddSpendableBalance(currencyBase, tx.Fee)
			}
		}

		if t, ok := txDecoded.(signing.Tx); ok {
			signers, err := t.GetSigners()
			if payer := t.FeePayer(); len(payer) > 0 {
				addr, err := addEmptyAddress(payer, currencyBase, block.Height, &data)
				if err != nil {
					return data, errors.Wrap(err, "add signer")
				}
				addr.AddSpendableBalance(currencyBase, tx.Fee.Neg())
			}

			if err != nil {
				pubKeys, err := t.GetPubKeys()
				if err != nil {
					return data, errors.Wrap(err, "get pub keys")
				}
				for _, pk := range pubKeys {
					if _, err := addEmptyAddress(pk.Address().Bytes(), currencyBase, block.Height, &data); err != nil {
						return data, errors.Wrap(err, "add signer")
					}
				}
			} else {
				for i := range signers {
					if _, err := addEmptyAddress(signers[i], currencyBase, block.Height, &data); err != nil {
						return data, errors.Wrap(err, "add signer")
					}
				}
			}
		}

	}

	for _, addr := range decodeCtx.Addresses.Values() {
		if a, ok := data.addresses[addr.String()]; ok {
			for i := range addr.Balances {
				for j := range a.Balances {
					if a.Balances[j].Currency == addr.Balances[i].Currency {
						a.Balances[j].Spendable = a.Balances[j].Spendable.Add(addr.Balances[i].Spendable)
						a.Balances[j].Delegated = a.Balances[j].Delegated.Add(addr.Balances[i].Delegated)
						a.Balances[j].Unbonding = a.Balances[j].Unbonding.Add(addr.Balances[i].Unbonding)
						break
					}
				}
			}
		} else {
			data.addresses[addr.String()] = addr
		}
	}

	if err := module.parseConstants(genesis.AppState, genesis.ConsensusParams, &data); err != nil {
		return data, errors.Wrap(err, "parse constants")
	}

	module.parseTotalSupply(genesis.AppState.Bank.Supply, &block)

	data.validators = decodeCtx.Validators.Values()
	data.stakingLogs = decodeCtx.StakingLogs

	for _, delegation := range decodeCtx.Delegations.All() {
		data.delegations = append(data.delegations, *delegation)
		if data.bondedTokensPool != nil && len(data.bondedTokensPool.Balances) > 0 {
			data.bondedTokensPool.Balances[0].Spendable = data.bondedTokensPool.Balances[0].Spendable.Add(delegation.Amount)
		}
		if addr, ok := data.addresses[delegation.Address.Address]; ok && len(addr.Balances) > 0 {
			addr.Balances[0].Spendable = addr.Balances[0].Spendable.Sub(delegation.Amount)
		}
	}

	if err := module.parseFeeGrants(genesis.AppState.Feegrant.FeeGrants, block, &data); err != nil {
		return data, errors.Wrap(err, "parse genesis fee grants")
	}

	if err := module.parseAccounts(genesis.AppState.Auth.Accounts, block, &data); err != nil {
		return data, errors.Wrap(err, "parse genesis accounts")
	}
	if err := module.parseBalances(genesis.AppState.Bank.Balances, block.Height, &data); err != nil {
		return data, errors.Wrap(err, "parse genesis account balances")
	}

	data.block = block
	return data, nil
}

func addEmptyAddress(pk []byte, denom string, height pkgTypes.Level, data *parsedData) (*storage.Address, error) {
	address, err := pkgTypes.NewAddressFromBytes(pk)
	if err != nil {
		return nil, errors.Wrap(err, "NewAddressFromBytes")
	}

	readableAddress := address.String()
	if addr, ok := data.addresses[readableAddress]; ok {
		return addr, nil
	}
	data.addresses[readableAddress] = &storage.Address{
		Address:    readableAddress,
		Hash:       pk,
		Height:     height,
		LastHeight: height,
		Balances: []storage.Balance{
			{
				Spendable: storageTypes.NumericZero(),
				Delegated: storageTypes.NumericZero(),
				Unbonding: storageTypes.NumericZero(),
				Currency:  denom,
			},
		},
	}
	return data.addresses[readableAddress], nil
}

func (module *Module) parseTotalSupply(supply []types.Supply, block *storage.Block) {
	if len(supply) == 0 {
		return
	}

	if totalSupply, err := storageTypes.NumericFromString(supply[0].Amount); err == nil {
		block.Stats.SupplyChange = totalSupply
	}
}

func getBaseCurrency(denomMetadata []storage.DenomMetadata) string {
	currencyBase := currency.DefaultCurrency
	if len(denomMetadata) > 0 {
		currencyBase = denomMetadata[0].Base
	}
	return currencyBase
}

func (module *Module) parseAccounts(accounts []types.Account, block storage.Block, data *parsedData) error {
	currencyBase := getBaseCurrency(data.denomMetadata)

	for i := range accounts {
		address := storage.Address{
			Height:     block.Height,
			LastHeight: block.Height,
			Balances: []storage.Balance{
				{
					Spendable: storageTypes.NumericZero(),
					Delegated: storageTypes.NumericZero(),
					Unbonding: storageTypes.NumericZero(),
					Currency:  currencyBase,
				},
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

			switch address.Name {
			case "bonded_tokens_pool":
				data.bondedTokensPool = &address
			case "fee_collector":
				data.feeCollector = &address
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

		case strings.Contains(accounts[i].Type, "PermanentLockedAccount"):
			readableAddress = accounts[i].BaseVestingAccount.BaseAccount.Address
			if err := parseVesting(accounts[i], block, readableAddress, storageTypes.VestingTypePermanent, data); err != nil {
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

func addAddress(data *parsedData, block storage.Block, address string) error {
	if _, ok := data.addresses[address]; ok {
		return nil
	}
	currencyBase := getBaseCurrency(data.denomMetadata)

	addr := &storage.Address{
		Address:    address,
		Height:     block.Height,
		LastHeight: block.Height,
		Balances: []storage.Balance{
			{
				Spendable: storageTypes.NumericZero(),
				Delegated: storageTypes.NumericZero(),
				Unbonding: storageTypes.NumericZero(),
				Currency:  currencyBase,
			},
		},
	}
	_, hash, err := pkgTypes.Address(address).Decode()
	if err != nil {
		return err
	}
	addr.Hash = hash
	data.addresses[address] = addr
	return nil
}

func (module *Module) parseBalances(balances []types.Balances, height pkgTypes.Level, data *parsedData) error {
	for i := range balances {
		if len(balances[i].Coins) == 0 {
			continue
		}

		addr, ok := data.addresses[balances[i].Address]
		if !ok {
			_, hash, err := pkgTypes.Address(balances[i].Address).Decode()
			if err != nil {
				return err
			}
			addr = &storage.Address{
				Hash:       hash,
				Address:    balances[i].Address,
				Height:     height,
				LastHeight: height,
				Balances:   make([]storage.Balance, 0, len(balances[i].Coins)),
			}
			data.addresses[addr.String()] = addr
		}

		for _, coin := range balances[i].Coins {
			value, err := storageTypes.NumericFromString(coin.Amount)
			if err != nil {
				continue
			}
			addr.AddSpendableBalance(coin.Denom, value)
		}
	}

	return nil
}

func getAmountFromOriginalVesting(vestings []types.Coins) (storageTypes.Numeric, error) {
	amount := storageTypes.NumericZero().Copy()

	for i := range vestings {
		val, err := storageTypes.NumericFromString(vestings[i].Amount)
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
