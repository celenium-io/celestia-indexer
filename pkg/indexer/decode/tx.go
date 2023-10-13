// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"github.com/celestiaorg/celestia-app/app"
	"github.com/celestiaorg/celestia-app/app/encoding"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/dipdup-io/celestia-indexer/internal/consts"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	tmTypes "github.com/tendermint/tendermint/types"
)

type DecodedTx struct {
	AuthInfo      tx.AuthInfo
	TimeoutHeight uint64
	Memo          string
	Messages      []cosmosTypes.Msg
	Fee           decimal.Decimal
	Signers       map[string]struct{}
}

var (
	cfg, decoder = createDecoder()
)

func Tx(b types.BlockData, index int) (d DecodedTx, err error) {
	raw := b.Block.Txs[index]
	if bTx, isBlob := tmTypes.UnmarshalBlobTx(raw); isBlob {
		raw = bTx.Tx
	}

	d.AuthInfo, d.Fee, err = decodeAuthInfo(cfg, raw)
	if err != nil {
		return
	}

	d.TimeoutHeight, d.Memo, d.Messages, err = decodeCosmosTx(decoder, raw)
	if err != nil {
		return
	}

	d.Signers = make(map[string]struct{})
	for i := range d.Messages {
		for _, signer := range d.Messages[i].GetSigners() {
			d.Signers[signer.String()] = struct{}{}
		}
	}

	return
}

func decodeCosmosTx(decoder cosmosTypes.TxDecoder, raw tmTypes.Tx) (timeoutHeight uint64, memo string, messages []cosmosTypes.Msg, err error) {
	txDecoded, err := decoder(raw)
	if err != nil {
		err = errors.Wrap(err, "decoding tx error")
		return
	}

	if t, ok := txDecoded.(cosmosTypes.TxWithTimeoutHeight); ok {
		timeoutHeight = t.GetTimeoutHeight()
	}
	if t, ok := txDecoded.(cosmosTypes.TxWithMemo); ok {
		memo = t.GetMemo()
	}

	messages = txDecoded.GetMsgs()
	return
}

func decodeAuthInfo(cfg encoding.Config, raw tmTypes.Tx) (tx.AuthInfo, decimal.Decimal, error) {
	var txRaw tx.TxRaw
	if e := cfg.Codec.Unmarshal(raw, &txRaw); e != nil {
		return tx.AuthInfo{}, decimal.Decimal{}, errors.Wrap(e, "unmarshalling tx error")
	}

	var authInfo tx.AuthInfo
	if e := cfg.Codec.Unmarshal(txRaw.AuthInfoBytes, &authInfo); e != nil {
		return tx.AuthInfo{}, decimal.Decimal{}, errors.Wrap(e, "decoding tx auth_info error")
	}

	fee, err := decodeFee(authInfo)
	if err != nil {
		return authInfo, decimal.Zero, err
	}

	return authInfo, fee, nil
}

func decodeFee(authInfo tx.AuthInfo) (decimal.Decimal, error) {
	amount := authInfo.GetFee().GetAmount()

	if amount == nil {
		return decimal.Zero, nil
	}

	if len(amount) > 1 {
		// TODO stop indexer if tx is not in failed status
		return decimal.Zero, errors.Errorf("found fee in %d currencies", len(amount))
	}

	fee, ok := getFeeInDenom(amount, consts.Utia)
	if !ok {
		if fee, ok = getFeeInDenom(amount, consts.Tia); !ok {
			// TODO stop indexer if tx is not in failed status
			return decimal.Zero, errors.New("couldn't find fee amount in utia or in tia denom")
		}
	}

	return fee, nil
}

func getFeeInDenom(amount cosmosTypes.Coins, denom consts.Denom) (decimal.Decimal, bool) {
	ok, utiaCoin := amount.Find(string(denom))
	if !ok {
		return decimal.Zero, false
	}

	switch denom {
	case consts.Utia:
		fee := decimal.NewFromBigInt(utiaCoin.Amount.BigInt(), 0)
		return fee, true
	case consts.Tia:
		fee := decimal.NewFromBigInt(utiaCoin.Amount.BigInt(), 6)
		return fee, true
	default:
		return decimal.Zero, false
	}
}

func createDecoder() (encoding.Config, cosmosTypes.TxDecoder) {
	cfg := encoding.MakeConfig(app.ModuleEncodingRegisters...)
	return cfg, cfg.TxConfig.TxDecoder()
}

func JsonTx(raw []byte) (cosmosTypes.Tx, error) {
	return cfg.TxConfig.TxJSONDecoder()(raw)
}
