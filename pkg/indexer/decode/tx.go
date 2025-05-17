// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/celestiaorg/celestia-app/v4/app"
	"github.com/celestiaorg/celestia-app/v4/app/encoding"
	appBlobTypes "github.com/celestiaorg/celestia-app/v4/x/blob/types"
	blobTypes "github.com/cometbft/cometbft/proto/tendermint/types"
	tmTypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type DecodedTx struct {
	AuthInfo      tx.AuthInfo
	TimeoutHeight uint64
	Memo          string
	Messages      []cosmosTypes.Msg
	Fee           decimal.Decimal
	Signers       map[types.Address][]byte
	Blobs         []*blobTypes.Blob
}

var (
	cfg, txDecoder = createDecoder()
)

func Tx(b types.BlockData, index int) (d DecodedTx, err error) {
	raw := b.Block.Txs[index]
	if bTx, isBlob := tmTypes.UnmarshalBlobTx(raw); isBlob {
		raw = bTx.Tx
		d.Blobs = bTx.Blobs
	}

	d.AuthInfo, d.Fee, err = decodeAuthInfo(raw)
	if err != nil {
		return
	}

	d.TimeoutHeight, d.Memo, d.Messages, err = decodeCosmosTx(txDecoder, raw)
	if err != nil {
		return
	}
	d.Signers = make(map[types.Address][]byte)

	for i := range d.Messages {
		if pfb, ok := d.Messages[i].(*appBlobTypes.MsgPayForBlobs); ok {
			address := types.Address(pfb.Signer)
			_, hash, err := address.Decode()
			if err != nil {
				return d, errors.Wrap(err, "decode PFB signer")
			}
			d.Signers[address] = hash
		}
	}

	for _, signer := range d.AuthInfo.GetSignerInfos() {
		publickKey := signer.GetPublicKey()
		if publickKey == nil {
			continue
		}
		var pk secp256k1.PubKey
		if err := cfg.Codec.Unmarshal(publickKey.Value, &pk); err != nil {
			return d, errors.Wrap(err, "signer decoding")
		}
		address, err := types.NewAddressFromBytes(pk.Bytes())
		if err != nil {
			return d, err
		}
		d.Signers[address] = pk.Bytes()
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

func decodeAuthInfo(raw tmTypes.Tx) (tx.AuthInfo, decimal.Decimal, error) {
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

	fee, ok := getFeeInDenom(amount, currency.Utia)
	if !ok {
		if fee, ok = getFeeInDenom(amount, currency.Tia); !ok {
			// TODO stop indexer if tx is not in failed status
			return decimal.Zero, errors.New("couldn't find fee amount in utia or in tia denom")
		}
	}

	return fee, nil
}

func getFeeInDenom(amount cosmosTypes.Coins, denom string) (decimal.Decimal, bool) {
	ok, utiaCoin := amount.Find(denom)
	if !ok {
		return decimal.Zero, false
	}

	switch denom {
	case currency.Utia:
		fee := decimal.NewFromBigInt(utiaCoin.Amount.BigInt(), 0)
		return fee, true
	case currency.Tia:
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
