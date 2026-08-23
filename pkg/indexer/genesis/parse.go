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

func (module *Module) parse(genesis types.GenesisOutput) (*context.Context, error) {
	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
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

	var (
		bondedTokensPool *storage.Address
		feeCollector     *storage.Address
	)

	module.parseDenomMetadata(decodeCtx, genesis.AppState.Bank.DenomMetadata)
	currencyBase := getBaseCurrency(decodeCtx.DenomMetadata)

	if err := module.parseAccounts(decodeCtx, currencyBase, genesis.ModuleAccs); err != nil {
		return decodeCtx, errors.Wrap(err, "parse genesis accounts")
	}

	for acc := range decodeCtx.Addresses.AllValues() {
		switch acc.Name {
		case "bonded_tokens_pool":
			bondedTokensPool = acc
		case "fee_collector":
			feeCollector = acc
		}
	}

	for index, genTx := range genesis.AppState.Genutil.GenTxs {
		txDecoded, err := decode.JsonTx(genTx)
		if err != nil {
			return decodeCtx, errors.Wrapf(err, "failed to decode GenTx '%s'", genTx)
		}

		memoTx, ok := txDecoded.(cosmosTypes.TxWithMemo)
		if !ok {
			return decodeCtx, errors.Wrapf(err, "expected TxWithMemo, got %T", genTx)
		}
		txWithTimeoutHeight, ok := txDecoded.(cosmosTypes.TxWithTimeoutHeight)
		if !ok {
			return decodeCtx, errors.Wrapf(err, "expected TxWithTimeoutHeight, got %T", genTx)
		}

		tx := &decodeCtx.Block.Txs[index]
		tx.Height = decodeCtx.Block.Height
		tx.Time = decodeCtx.Block.Time
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
			return decodeCtx, err
		}

		for msgIndex, msg := range txDecoded.GetMsgs() {
			decoded, err := decode.Message(decodeCtx, msg, msgIndex, storageTypes.StatusSuccess, tx.Id)
			if err != nil {
				return decodeCtx, errors.Wrap(err, "decode genesis message")
			}

			tx.Messages[msgIndex] = decoded.Msg
			tx.MessageTypes.SetByMsgType(decoded.Msg.Type)
			decodeCtx.Block.MessageTypes.SetByMsgType(decoded.Msg.Type)
			tx.BlobsSize += decoded.BlobsSize
		}

		if t, ok := txDecoded.(cosmosTypes.FeeTx); ok {
			value, err := decode.DecodeFee(t.GetFee())
			if err != nil {
				return decodeCtx, errors.Wrap(err, "fee decoding")
			}
			tx.Fee = storageTypes.NewNumeric(value)
			decodeCtx.Block.Stats.Fee = decodeCtx.Block.Stats.Fee.Add(tx.Fee)
			if feeCollector != nil {
				feeCollector.AddSpendableBalance(currencyBase, tx.Fee)
			}
		}

		if t, ok := txDecoded.(signing.Tx); ok {
			signers, err := t.GetSigners()
			if payer := t.FeePayer(); len(payer) > 0 {
				addr, err := addEmptyAddress(decodeCtx, payer, currencyBase)
				if err != nil {
					return decodeCtx, errors.Wrap(err, "add signer")
				}
				addr.AddSpendableBalance(currencyBase, tx.Fee.Neg())
			}

			if err != nil {
				pubKeys, err := t.GetPubKeys()
				if err != nil {
					return decodeCtx, errors.Wrap(err, "get pub keys")
				}
				for i := range pubKeys {
					signerAddress, err := addEmptyAddress(decodeCtx, pubKeys[i].Address().Bytes(), currencyBase)
					if err != nil {
						return decodeCtx, errors.Wrap(err, "add signer")
					}
					tx.Signers = append(tx.Signers, *signerAddress)
				}
			} else {
				for i := range signers {
					signerAddress, err := addEmptyAddress(decodeCtx, signers[i], currencyBase)
					if err != nil {
						return decodeCtx, errors.Wrap(err, "add signer")
					}
					tx.Signers = append(tx.Signers, *signerAddress)
				}
			}
		}
	}

	if err := module.parseConstants(decodeCtx, genesis.AppState, genesis.ConsensusParams); err != nil {
		return decodeCtx, errors.Wrap(err, "parse constants")
	}

	module.parseTotalSupply(genesis.AppState.Bank.Supply, decodeCtx.Block)

	for _, delegation := range decodeCtx.Delegations.All() {
		if bondedTokensPool != nil && len(bondedTokensPool.Balances) > 0 {
			bondedTokensPool.Balances[0].Spendable = bondedTokensPool.Balances[0].Spendable.Add(delegation.Amount)
		}
		if addr, ok := decodeCtx.Addresses.Get(delegation.Address.Address); ok && len(addr.Balances) > 0 {
			addr.Balances[0].Spendable = addr.Balances[0].Spendable.Sub(delegation.Amount)
		}
	}

	if err := module.parseFeeGrants(decodeCtx, genesis.AppState.Feegrant.FeeGrants); err != nil {
		return decodeCtx, errors.Wrap(err, "parse genesis fee grants")
	}
	if err := module.parseAuthzGrants(decodeCtx, genesis.AppState.Authz.Authorization); err != nil {
		return decodeCtx, errors.Wrap(err, "parse genesis authz grants")
	}

	if err := module.parseAccounts(decodeCtx, currencyBase, genesis.AppState.Auth.Accounts); err != nil {
		return decodeCtx, errors.Wrap(err, "parse genesis accounts")
	}
	if err := module.parseBalances(decodeCtx, genesis.AppState.Bank.Balances); err != nil {
		return decodeCtx, errors.Wrap(err, "parse genesis account balances")
	}

	return decodeCtx, nil
}

func addEmptyAddress(ctx *context.Context, pk []byte, denom string) (*storage.Address, error) {
	address, err := pkgTypes.NewAddressFromBytes(pk)
	if err != nil {
		return nil, errors.Wrap(err, "NewAddressFromBytes")
	}

	readableAddress := address.String()
	if addr, ok := ctx.Addresses.Get(readableAddress); ok {
		return addr, nil
	}

	addr := &storage.Address{
		Address:    readableAddress,
		Hash:       pk,
		Height:     ctx.Block.Height,
		LastHeight: ctx.Block.Height,
		Balances: []storage.Balance{
			{
				Spendable: storageTypes.NumericZero(),
				Delegated: storageTypes.NumericZero(),
				Unbonding: storageTypes.NumericZero(),
				Currency:  denom,
			},
		},
	}
	err = ctx.AddAddress(addr)
	return addr, err
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

func (module *Module) parseAccounts(ctx *context.Context, denom string, accounts []types.Account) error {
	for i := range accounts {
		address := storage.Address{
			Height:     ctx.Block.Height,
			LastHeight: ctx.Block.Height,
			Balances: []storage.Balance{
				{
					Spendable: storageTypes.NumericZero(),
					Delegated: storageTypes.NumericZero(),
					Unbonding: storageTypes.NumericZero(),
					Currency:  denom,
				},
			},
		}

		var readableAddress string

		switch {
		case strings.Contains(accounts[i].Type, "PeriodicVestingAccount"):
			readableAddress = accounts[i].BaseVestingAccount.BaseAccount.Address
			if err := parseVesting(ctx, accounts[i], readableAddress, storageTypes.VestingTypePeriodic); err != nil {
				return err
			}

		case strings.Contains(accounts[i].Type, "ModuleAccount"):
			readableAddress = accounts[i].BaseAccount.Address
			address.Name = accounts[i].Name

		case strings.Contains(accounts[i].Type, "BaseAccount"):
			readableAddress = accounts[i].Address

		case strings.Contains(accounts[i].Type, "ContinuousVestingAccount"):
			readableAddress = accounts[i].BaseVestingAccount.BaseAccount.Address
			if err := parseVesting(ctx, accounts[i], readableAddress, storageTypes.VestingTypeContinuous); err != nil {
				return err
			}

		case strings.Contains(accounts[i].Type, "DelayedVestingAccount"):
			readableAddress = accounts[i].BaseVestingAccount.BaseAccount.Address
			if err := parseVesting(ctx, accounts[i], readableAddress, storageTypes.VestingTypeDelayed); err != nil {
				return err
			}

		case strings.Contains(accounts[i].Type, "PermanentLockedAccount"):
			readableAddress = accounts[i].BaseVestingAccount.BaseAccount.Address
			if err := parseVesting(ctx, accounts[i], readableAddress, storageTypes.VestingTypePermanent); err != nil {
				return err
			}

		default:
			return errors.Errorf("unknown account type: %s", accounts[i].Type)
		}

		address.Address = readableAddress
		if err := ctx.AddAddress(&address); err != nil {
			return err
		}
	}
	return nil
}

func (module *Module) parseBalances(ctx *context.Context, balances []types.Balances) error {
	for i := range balances {
		if len(balances[i].Coins) == 0 {
			continue
		}

		addr, ok := ctx.Addresses.Get(balances[i].Address)
		if !ok {
			addr = &storage.Address{
				Address:    balances[i].Address,
				Height:     ctx.Block.Height,
				LastHeight: ctx.Block.Height,
				Balances:   make([]storage.Balance, 0, len(balances[i].Coins)),
			}
		}

		for _, coin := range balances[i].Coins {
			value, err := storageTypes.NumericFromString(coin.Amount)
			if err != nil {
				continue
			}
			addr.AddSpendableBalance(coin.Denom, value)
		}

		if !ok {
			if err := ctx.AddAddress(addr); err != nil {
				return err
			}
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

func parseVesting(ctx *context.Context, acc types.Account, address string, typ storageTypes.VestingType) error {
	amount, err := getAmountFromOriginalVesting(acc.BaseVestingAccount.OriginalVesting)
	if err != nil {
		return err
	}

	v := storage.VestingAccount{
		Height: ctx.Block.Height,
		Time:   ctx.Block.Time,
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
	ctx.AddVestingAccount(&v)
	return nil
}
