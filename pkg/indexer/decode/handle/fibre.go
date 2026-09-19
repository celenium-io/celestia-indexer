// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"encoding/base64"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/celestiaorg/celestia-app/v10/pkg/appconsts"
	fibreTypes "github.com/celestiaorg/celestia-app/v10/x/fibre/types"
	"github.com/celestiaorg/go-square/v4/inclusion"
	"github.com/celestiaorg/go-square/v4/share"
	"github.com/cometbft/cometbft/crypto/merkle"
	"github.com/cosmos/btcutil/bech32"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// MsgDepositToEscrow tops up the signer's fibre escrow account, which pays for
// blobs served over fibre.
func MsgDepositToEscrow(
	ctx *context.Context, msgId uint64, m *fibreTypes.MsgDepositToEscrow,
) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgDepositToEscrow
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgRequestWithdrawal starts a withdrawal from the signer's escrow account.
// The payout happens later, after fibre.WithdrawalDelay.
func MsgRequestWithdrawal(
	ctx *context.Context, msgId uint64, m *fibreTypes.MsgRequestWithdrawal,
) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgRequestWithdrawal
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgPayForFibre settles a validator-signed payment promise against the escrow
// account of the promise signer. Only the submitter is indexed: the escrow
// owner is identified by a secp256k1 pubkey inside the promise, not a bech32
// address, and the promise itself stays in the raw message data.
func MsgPayForFibre(
	ctx *context.Context, status storageTypes.Status, txId, msgId uint64, m *fibreTypes.MsgPayForFibre,
) (storageTypes.MsgType, []*storage.BlobLog, error) {
	msgType := storageTypes.MsgPayForFibre
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	if err != nil {
		return msgType, nil, err
	}

	shareNamespace, err := share.NewNamespaceFromBytes(m.PaymentPromise.GetNamespace())
	if err != nil {
		return storageTypes.MsgUnknown, nil, errors.Wrap(err, "NewNamespaceFromBytes")
	}

	size := int64(m.PaymentPromise.GetBlobSize())

	namespace := &storage.Namespace{
		FirstHeight:     ctx.Block.Height,
		Version:         shareNamespace.Version(),
		NamespaceID:     shareNamespace.ID(),
		PffCount:        1,
		Reserved:        shareNamespace.IsReserved(),
		LastHeight:      ctx.Block.Height,
		LastMessageTime: ctx.Block.Time,
	}

	if status == storageTypes.StatusSuccess {
		namespace.BlobsCount = 1
		namespace.FibreSize = size

		ns := ctx.AddNamespace(namespace)

		owner, err := pkgTypes.NewAddressFromBytes(m.PaymentPromise.SignerPublicKey.Address().Bytes())
		if err != nil {
			return msgType, nil, err
		}

		signer := &storage.Address{
			Address:    owner.String(),
			Height:     ctx.Block.Height,
			LastHeight: ctx.Block.Height,
		}
		if err := ctx.AddAddress(signer); err != nil {
			return msgType, nil, err
		}

		_, signerBytes, err := bech32.DecodeToBase256(m.Signer)
		if err != nil {
			return msgType, nil, err
		}

		b, err := share.NewV2Blob(
			shareNamespace,
			m.PaymentPromise.GetBlobVersion(),
			m.PaymentPromise.GetCommitment(),
			signerBytes)
		if err != nil {
			return msgType, nil, errors.Wrap(err, "can't create blob")
		}

		commitment, err := inclusion.CreateCommitment(b, merkle.HashFromByteSlices, appconsts.SubtreeRootThreshold)
		if err != nil {
			return msgType, nil, errors.Wrap(err, "can't compute commitment")
		}

		estimatedGas := fibreTypes.EstimateGasForPayForFibre(m.PaymentPromise.GetBlobSize())
		gas := decimal.NewFromUint64(estimatedGas)

		blob := &storage.BlobLog{
			Commitment:   base64.StdEncoding.EncodeToString(commitment),
			Size:         size,
			Namespace:    ns,
			Height:       ctx.Block.Height,
			Time:         ctx.Block.Time,
			Signer:       signer,
			ShareVersion: int(share.ShareVersionTwo),
			MsgId:        msgId,
			TxId:         txId,
			Source:       storageTypes.BlobSourceFibre,
			GasConsumed:  gas,
			Fee:          storageTypes.NewNumeric(gas), // 1 utia per gas
		}

		ctx.AddNamespaceMessage(&storage.NamespaceMessage{
			MsgId:     msgId,
			TxId:      txId,
			Height:    ctx.Block.Height,
			Time:      ctx.Block.Time,
			Namespace: ns,
			Size:      uint64(size),
			Source:    blob.Source,
		})
		return msgType, []*storage.BlobLog{blob}, nil
	}

	return msgType, nil, nil
}

// MsgPaymentPromiseTimeout settles a promise that was never claimed within
// fibre.PaymentPromiseTimeout.
func MsgPaymentPromiseTimeout(
	ctx *context.Context, msgId uint64, m *fibreTypes.MsgPaymentPromiseTimeout,
) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgPaymentPromiseTimeout
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgUpdateFibreParams updates the fibre module params via governance.
func MsgUpdateFibreParams(
	ctx *context.Context, msgId uint64, m *fibreTypes.MsgUpdateFibreParams,
) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateFibreParams
	ctx.AddFibreParams(m.Params)

	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height, msgId)
	return msgType, err
}
