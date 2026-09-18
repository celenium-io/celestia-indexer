// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celestiaorg/celestia-app/v10/pkg/appconsts"
	fibreTypes "github.com/celestiaorg/celestia-app/v10/x/fibre/types"
	valaddrTypes "github.com/celestiaorg/celestia-app/v10/x/valaddr/types"
	"github.com/celestiaorg/go-square/v4/inclusion"
	"github.com/celestiaorg/go-square/v4/share"
	"github.com/cometbft/cometbft/crypto/merkle"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

const (
	fibreSigner    = "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l"
	fibreAuthority = "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6"
	fibreValidator = "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x"

	// fibreTxId is deliberately different from the message id the decoder
	// assigns (it counts from 1) so a swap of the two shows up in assertions.
	fibreTxId = 42

	// fibreChunk is the paid upload step: a promise's blob_size is always a
	// multiple of it (4096 rows x 64 byte minimum row size).
	fibreChunk = appconsts.PFBFibreChunkSize
)

// fibreNamespace is a valid non-reserved blob namespace: version byte 0 plus a
// 28 byte id.
var fibreNamespace = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22}

// sizedMsg is the subset of the proto interface decode.Message needs to size a message.
type sizedMsg interface {
	cosmosTypes.Msg
	Size() int
}

// decodeFibreMsg runs a message through the decoder with a fresh context and
// returns the decoded message plus the addresses the handler registered.
func decodeFibreMsg(t *testing.T, msg sizedMsg) (decode.DecodedMsg, *context.Context, *storage.Block, time.Time) {
	t.Helper()
	return decodeFibreMsgWithStatus(t, msg, storageTypes.StatusSuccess)
}

func decodeFibreMsgWithStatus(
	t *testing.T, msg sizedMsg, status storageTypes.Status,
) (decode.DecodedMsg, *context.Context, *storage.Block, time.Time) {
	t.Helper()

	block, now := testsuite.EmptyBlock()
	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   now,
	}

	dm, err := decode.Message(decodeCtx, msg, 0, status, fibreTxId)
	require.NoError(t, err)
	// Fibre payload never reaches the data square, so it adds nothing to the
	// square size the block accounts for.
	require.EqualValues(t, 0, dm.BlobsSize)

	return dm, decodeCtx, decodeCtx.Block, now
}

// requireMsgAddress asserts the handler linked exactly one address to the message
// under the expected role.
func requireMsgAddress(t *testing.T, ctx *context.Context, address string, addrType storageTypes.MsgAddressType) {
	t.Helper()

	addresses := ctx.AddressMessages.Values()
	require.Len(t, addresses, 1)
	require.Equal(t, addrType, addresses[0].Type)
	require.Equal(t, address, addresses[0].Address.Address)
	require.EqualValues(t, 1, addresses[0].MsgId)
}

func TestFibre_MsgDepositToEscrow(t *testing.T) {
	msg := &fibreTypes.MsgDepositToEscrow{
		Signer: fibreSigner,
		Amount: cosmosTypes.NewInt64Coin("utia", 1_000_000),
	}

	dm, ctx, block, now := decodeFibreMsg(t, msg)

	require.Equal(t, storage.Message{
		Id:       1,
		Height:   block.Height,
		Time:     now,
		Position: 0,
		Type:     storageTypes.MsgDepositToEscrow,
		TxId:     fibreTxId,
		Data:     mustMsgToMap(t, msg),
		Size:     msg.Size(),
	}, dm.Msg)
	requireMsgAddress(t, ctx, fibreSigner, storageTypes.MsgAddressTypeSigner)
}

func TestFibre_MsgRequestWithdrawal(t *testing.T) {
	msg := &fibreTypes.MsgRequestWithdrawal{
		Signer: fibreSigner,
		Amount: cosmosTypes.NewInt64Coin("utia", 500_000),
	}

	dm, ctx, block, now := decodeFibreMsg(t, msg)

	require.Equal(t, storage.Message{
		Id:       1,
		Height:   block.Height,
		Time:     now,
		Position: 0,
		Type:     storageTypes.MsgRequestWithdrawal,
		TxId:     fibreTxId,
		Data:     mustMsgToMap(t, msg),
		Size:     msg.Size(),
	}, dm.Msg)
	requireMsgAddress(t, ctx, fibreSigner, storageTypes.MsgAddressTypeSigner)
}

func TestFibre_MsgPayForFibre(t *testing.T) {
	msg := newMsgPayForFibre(t, fibreChunk)

	dm, ctx, block, now := decodeFibreMsg(t, msg)

	require.Equal(t, storage.Message{
		Id:       1,
		Height:   block.Height,
		Time:     now,
		Position: 0,
		Type:     storageTypes.MsgPayForFibre,
		TxId:     fibreTxId,
		Data:     mustMsgToMap(t, msg),
		Size:     msg.Size(),
	}, dm.Msg)
	// The escrow owner is a pubkey inside the promise, so only the submitter is indexed.
	requireMsgAddress(t, ctx, fibreSigner, storageTypes.MsgAddressTypeSigner)

	// The promise namespace is tracked apart from PFB traffic: pff_count and
	// fibre_size move, pfb_count and size do not.
	namespaces := ctx.Namespaces.Values()
	require.Len(t, namespaces, 1)
	ns := namespaces[0]
	require.EqualValues(t, 0, ns.Version)
	require.Equal(t, fibreNamespace[1:], ns.NamespaceID)
	require.EqualValues(t, 1, ns.PffCount)
	require.EqualValues(t, fibreChunk, ns.FibreSize)
	require.EqualValues(t, 1, ns.BlobsCount)
	require.EqualValues(t, 0, ns.PfbCount)
	require.EqualValues(t, 0, ns.Size)
	require.False(t, ns.Reserved)
	require.Equal(t, block.Height, ns.FirstHeight)
	require.Equal(t, block.Height, ns.LastHeight)
	require.Equal(t, now, ns.LastMessageTime)

	nsMsgs := ctx.NamespaceMessages.Values()
	require.Len(t, nsMsgs, 1)
	require.EqualValues(t, 1, nsMsgs[0].MsgId)
	require.EqualValues(t, fibreTxId, nsMsgs[0].TxId)
	require.EqualValues(t, fibreChunk, nsMsgs[0].Size)
	require.Equal(t, ns, nsMsgs[0].Namespace)

	require.Len(t, dm.BlobLogs, 1)
	blob := dm.BlobLogs[0]
	require.Equal(t, storageTypes.BlobSourceFibre, blob.Source)
	require.Equal(t, int(share.ShareVersionTwo), blob.ShareVersion)
	require.EqualValues(t, fibreChunk, blob.Size)
	require.EqualValues(t, 1, blob.MsgId)
	require.EqualValues(t, fibreTxId, blob.TxId)
	require.Equal(t, fibreSigner, blob.Signer.Address)
	require.Equal(t, ns, blob.Namespace)
	require.Equal(t, block.Height, blob.Height)
	require.Equal(t, now, blob.Time)
	// Fibre blobs carry no payload in the square, so there is nothing to sniff.
	require.Empty(t, blob.ContentType)

	// The commitment must be the one celestia-app puts in the square for this
	// message, not something recomputed differently on our side.
	require.Equal(t, systemBlobCommitment(t, msg), blob.Commitment)

	// One 256 KiB chunk: 650_000 + 45_000, charged at 1 utia per gas.
	requireFibrePayment(t, blob, fibreChunk)
}

// TestFibre_MsgPayForFibreMultipleChunks pins the per chunk part of the gas
// formula: a promise is charged per started 256 KiB chunk.
func TestFibre_MsgPayForFibreMultipleChunks(t *testing.T) {
	for _, chunks := range []uint32{1, 2, 8} {
		size := fibreChunk * chunks
		msg := newMsgPayForFibre(t, size)

		dm, _, _, _ := decodeFibreMsg(t, msg)

		require.Len(t, dm.BlobLogs, 1)
		blob := dm.BlobLogs[0]
		require.EqualValues(t, size, blob.Size)
		require.Equal(t, systemBlobCommitment(t, msg), blob.Commitment)
		requireFibrePayment(t, blob, size)

		expected := decimal.NewFromInt(int64(appconsts.PFBFibreGasFixedCost + appconsts.PFBFibreGasPerChunk*uint64(chunks)))
		require.Truef(t, expected.Equal(blob.GasConsumed), "%d chunks: want %s, got %s", chunks, expected, blob.GasConsumed)
	}
}

// TestFibre_MsgPayForFibreFailed checks a reverted PFF leaves no trace beyond
// the message itself and its signer.
func TestFibre_MsgPayForFibreFailed(t *testing.T) {
	msg := newMsgPayForFibre(t, fibreChunk)

	dm, ctx, _, _ := decodeFibreMsgWithStatus(t, msg, storageTypes.StatusFailed)

	require.Equal(t, storageTypes.MsgPayForFibre, dm.Msg.Type)
	require.Empty(t, dm.BlobLogs)
	require.Empty(t, ctx.Namespaces.Values())
	require.Empty(t, ctx.NamespaceMessages.Values())
	// The signer is still linked: a failed message is part of its history.
	requireMsgAddress(t, ctx, fibreSigner, storageTypes.MsgAddressTypeSigner)
}

func TestFibre_MsgPaymentPromiseTimeout(t *testing.T) {
	msg := &fibreTypes.MsgPaymentPromiseTimeout{
		Signer:         fibreSigner,
		PaymentPromise: testPaymentPromise(t, fibreChunk),
	}

	dm, ctx, block, now := decodeFibreMsg(t, msg)

	require.Equal(t, storage.Message{
		Id:       1,
		Height:   block.Height,
		Time:     now,
		Position: 0,
		Type:     storageTypes.MsgPaymentPromiseTimeout,
		TxId:     fibreTxId,
		Data:     mustMsgToMap(t, msg),
		Size:     msg.Size(),
	}, dm.Msg)
	requireMsgAddress(t, ctx, fibreSigner, storageTypes.MsgAddressTypeSigner)

	// A timed out promise also charges the escrow, but nothing was uploaded, so
	// no blob is recorded. The settlement itself is not indexed either.
	require.Empty(t, dm.BlobLogs)
	require.Empty(t, ctx.Namespaces.Values())
	require.Empty(t, ctx.NamespaceMessages.Values())
}

func TestFibre_MsgUpdateFibreParams(t *testing.T) {
	params := fibreTypes.DefaultParams()
	params.WithdrawalDelay = 48 * time.Hour
	params.PaymentPromiseTimeout = 2 * time.Hour
	params.PaymentPromiseHeightWindow = 2000
	params.ShardRetention = 6 * time.Hour
	params.FullStakeStorageBudget = 1 << 40
	require.NoError(t, params.Validate())

	msg := &fibreTypes.MsgUpdateFibreParams{
		Authority: fibreAuthority,
		Params:    params,
	}

	dm, ctx, block, now := decodeFibreMsg(t, msg)

	require.Equal(t, storage.Message{
		Id:       1,
		Height:   block.Height,
		Time:     now,
		Position: 0,
		Type:     storageTypes.MsgUpdateFibreParams,
		TxId:     fibreTxId,
		Data:     mustMsgToMap(t, msg),
		Size:     msg.Size(),
	}, dm.Msg)
	requireMsgAddress(t, ctx, fibreAuthority, storageTypes.MsgAddressTypeAuthority)

	// Durations are stored in nanoseconds, like the other duration constants.
	require.Equal(t, map[string]string{
		"withdrawal_delay":              "172800000000000",
		"payment_promise_timeout":       "7200000000000",
		"payment_promise_height_window": "2000",
		"shard_retention":               "21600000000000",
		"full_stake_storage_budget":     "1099511627776",
	}, fibreConstants(t, ctx))
}

// TestFibre_DefaultParamsMatchModuleDefaults keeps the constants the v10 upgrade
// seeds in sync with the app defaults the keeper falls back to.
func TestFibre_DefaultParamsMatchModuleDefaults(t *testing.T) {
	ctx := context.NewContext()
	ctx.Block = &storage.Block{}
	ctx.AddFibreParams(fibreTypes.DefaultParams())

	require.Equal(t, map[string]string{
		"withdrawal_delay":              "86400000000000",
		"payment_promise_timeout":       "3600000000000",
		"payment_promise_height_window": "1000",
		"shard_retention":               "14400000000000",
		"full_stake_storage_budget":     "2199023255552",
	}, fibreConstants(t, ctx))
}

func TestValaddr_MsgSetFibreProviderInfo(t *testing.T) {
	// The signer is a validator operator address, not an account.
	msg := &valaddrTypes.MsgSetFibreProviderInfo{
		Signer: fibreValidator,
		Host:   "fibre.example.com:443",
	}

	dm, ctx, block, now := decodeFibreMsg(t, msg)

	require.Equal(t, storage.Message{
		Id:         1,
		Height:     block.Height,
		Time:       now,
		Position:   0,
		Type:       storageTypes.MsgSetFibreProviderInfo,
		TxId:       fibreTxId,
		Data:       mustMsgToMap(t, msg),
		Size:       msg.Size(),
		Validators: []string{fibreValidator},
	}, dm.Msg)
	requireMsgAddress(t, ctx, fibreValidator, storageTypes.MsgAddressTypeValidator)

	validators := ctx.Validators.Values()
	require.Len(t, validators, 1)
	require.Equal(t, fibreValidator, validators[0].Address)
	require.NotNil(t, validators[0].FibreHost)
	require.Equal(t, "fibre.example.com:443", *validators[0].FibreHost)
	require.NotNil(t, validators[0].FibreHostHeight)
	require.Equal(t, block.Height, *validators[0].FibreHostHeight)
	require.EqualValues(t, 1, validators[0].MessagesCount)

	// The handler builds on storage.EmptyValidator(), so an upsert of this
	// row never blanks out a validator's existing profile fields: they carry
	// the "do not modify" sentinel rather than an empty string.
	require.Equal(t, storage.DoNotModify, validators[0].Moniker)
	require.Equal(t, storage.DoNotModify, validators[0].Website)
	require.Equal(t, storage.DoNotModify, validators[0].Identity)
	require.Equal(t, storage.DoNotModify, validators[0].Contacts)
	require.Equal(t, storage.DoNotModify, validators[0].Details)

	// Numeric fields must be a real zero decimal, not a nil-backed zero
	// value, since they get summed on later messages for this validator.
	require.True(t, validators[0].Rate.IsZero())
	require.True(t, validators[0].MaxRate.IsZero())
	require.True(t, validators[0].MaxChangeRate.IsZero())
	require.True(t, validators[0].MinSelfDelegation.IsZero())
	require.True(t, validators[0].Stake.IsZero())
	require.True(t, validators[0].Rewards.IsZero())
	require.True(t, validators[0].Commissions.IsZero())
}

// TestValaddr_MsgSetFibreProviderInfoUpdates covers the "or update" half of the
// message: a validator may move its fibre server, and the newest host wins.
func TestValaddr_MsgSetFibreProviderInfoUpdates(t *testing.T) {
	block, now := testsuite.EmptyBlock()
	ctx := context.NewContext()
	ctx.Block = &storage.Block{Height: block.Height, Time: now}

	for _, host := range []string{"old.example.com:443", "new.example.com:8443"} {
		_, err := decode.Message(ctx, &valaddrTypes.MsgSetFibreProviderInfo{
			Signer: fibreValidator,
			Host:   host,
		}, 0, storageTypes.StatusSuccess, fibreTxId)
		require.NoError(t, err)
	}

	validators := ctx.Validators.Values()
	require.Len(t, validators, 1)
	require.NotNil(t, validators[0].FibreHost)
	require.Equal(t, "new.example.com:8443", *validators[0].FibreHost)
	require.EqualValues(t, 2, validators[0].MessagesCount)
}

func TestValaddr_MsgSetFibreProviderInfoFailed(t *testing.T) {
	msg := &valaddrTypes.MsgSetFibreProviderInfo{
		Signer: fibreValidator,
		Host:   "fibre.example.com:443",
	}

	dm, ctx, _, _ := decodeFibreMsgWithStatus(t, msg, storageTypes.StatusFailed)

	require.Equal(t, storageTypes.MsgSetFibreProviderInfo, dm.Msg.Type)
	require.Empty(t, dm.Msg.Validators)
	require.Empty(t, ctx.Validators.Values())
	requireMsgAddress(t, ctx, fibreValidator, storageTypes.MsgAddressTypeValidator)
}

func fibreConstants(t *testing.T, ctx *context.Context) map[string]string {
	t.Helper()

	constants := make(map[string]string)
	for _, c := range ctx.Constants.Values() {
		require.Equal(t, storageTypes.ModuleNameFibre, c.Module)
		constants[c.Name] = c.Value
	}
	return constants
}

// requireFibrePayment checks the blob against the module's own gas formula and
// the 1 utia per gas rate the keeper settles at.
func requireFibrePayment(t *testing.T, blob *storage.BlobLog, blobSize uint32) {
	t.Helper()

	gas := decimal.NewFromUint64(fibreTypes.EstimateGasForPayForFibre(blobSize))
	require.Truef(t, gas.Equal(blob.GasConsumed), "gas: want %s, got %s", gas, blob.GasConsumed)

	payment := fibreTypes.PaymentAmount(blobSize)
	require.Equal(t, payment.Amount.String(), blob.Fee.String())
}

// systemBlobCommitment computes the commitment through celestia-app's own
// SystemBlob(), which is what ends up in the data square.
func systemBlobCommitment(t *testing.T, msg *fibreTypes.MsgPayForFibre) string {
	t.Helper()

	blob, err := msg.SystemBlob()
	require.NoError(t, err)

	commitment, err := inclusion.CreateCommitment(blob, merkle.HashFromByteSlices, appconsts.SubtreeRootThreshold)
	require.NoError(t, err)

	return base64.StdEncoding.EncodeToString(commitment)
}

func newMsgPayForFibre(t *testing.T, blobSize uint32) *fibreTypes.MsgPayForFibre {
	t.Helper()

	msg := &fibreTypes.MsgPayForFibre{
		Signer:         fibreSigner,
		PaymentPromise: testPaymentPromise(t, blobSize),
		ValidatorSignatures: [][]byte{
			testsuite.RandomBytes(64),
			testsuite.RandomBytes(64),
		},
	}
	require.NoError(t, msg.ValidateBasic())
	return msg
}

func testPaymentPromise(t *testing.T, blobSize uint32) fibreTypes.PaymentPromise {
	t.Helper()

	promise := fibreTypes.PaymentPromise{
		ChainId:           "mocha-4",
		Height:            100,
		Namespace:         fibreNamespace,
		BlobSize:          blobSize,
		BlobVersion:       fibreTypes.BlobVersionZero,
		Commitment:        testsuite.RandomBytes(32),
		CreationTimestamp: time.Unix(1757000000, 0).UTC(),
		SignerPublicKey:   secp256k1.PubKey{Key: testsuite.RandomBytes(33)},
		Signature:         testsuite.RandomBytes(64),
	}
	// The fixture must be a promise the chain would actually accept.
	require.NoError(t, promise.ValidateBasic())
	return promise
}
