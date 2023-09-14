package decode

import (
	"github.com/celestiaorg/celestia-app/app"
	"github.com/celestiaorg/celestia-app/app/encoding"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
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
		return tx.AuthInfo{}, decimal.Decimal{}, errors.Wrap(e, "unmarshaling tx error")
	}

	var authInfo tx.AuthInfo
	if e := cfg.Codec.Unmarshal(txRaw.AuthInfoBytes, &authInfo); e != nil {
		return tx.AuthInfo{}, decimal.Decimal{}, errors.Wrap(e, "decoding tx auth_info error")
	}
	amount := authInfo.GetFee().GetAmount()
	if len(amount) > 1 {
		// TODO stop indexer if tx is not in failed status
		return tx.AuthInfo{}, decimal.Decimal{}, errors.Errorf("found fee in %d currencies", len(amount))
	}
	ok, utiaCoin := amount.Find("utia")
	if !ok {
		// TODO stop indexer if tx is not in failed status
		return tx.AuthInfo{}, decimal.Decimal{}, errors.New("while getting fee amount in utia")
	}
	fee := decimal.NewFromBigInt(utiaCoin.Amount.BigInt(), 0)
	return authInfo, fee, nil
}

func createDecoder() (encoding.Config, cosmosTypes.TxDecoder) {
	cfg := encoding.MakeConfig(app.ModuleEncodingRegisters...)
	return cfg, cfg.TxConfig.TxDecoder()
}

func JsonTx(raw []byte) (cosmosTypes.Tx, error) {
	cfg := encoding.MakeConfig(app.ModuleEncodingRegisters...)
	return cfg.TxConfig.TxJSONDecoder()(raw)
}
