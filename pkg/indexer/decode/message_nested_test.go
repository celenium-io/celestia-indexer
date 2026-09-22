// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package decode

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

const (
	nestedTestHeight    = 100
	nestedTestDelegator = "celestia1lkrd86urrmhmsvgzfygjsguv3cgv0036hrj0m9"
	nestedTestValidator = "celestiavaloper1qe8uuf5x69c526h4nzxwv4ltftr73v7q5qhs58"
)

func newNestedTestContext() *context.Context {
	ctx := context.NewContext()
	ctx.Block = &storage.Block{Height: nestedTestHeight, Time: time.Now().UTC()}
	return ctx
}

func newNestedTestDelegate() *cosmosStakingTypes.MsgDelegate {
	return &cosmosStakingTypes.MsgDelegate{
		DelegatorAddress: nestedTestDelegator,
		ValidatorAddress: nestedTestValidator,
		Amount:           cosmosTypes.NewCoin("utia", math.NewInt(1)),
	}
}

func TestMessage_IdFromBlockHeightAndPosition(t *testing.T) {
	ctx := newNestedTestContext()

	first, err := Message(ctx, createMsgUnknown(), 0, storageTypes.StatusSuccess, 0)
	require.NoError(t, err)
	second, err := Message(ctx, createMsgUnknown(), 1, storageTypes.StatusSuccess, 0)
	require.NoError(t, err)

	require.Equal(t, uint64(nestedTestHeight)<<24, first.Msg.Id)
	require.Equal(t, uint64(nestedTestHeight)<<24|1, second.Msg.Id)
	require.EqualValues(t, nestedTestHeight, first.Msg.Height)
}

func TestNestedMessage_UsesParentId(t *testing.T) {
	ctx := newNestedTestContext()

	parent, err := Message(ctx, createMsgUnknown(), 0, storageTypes.StatusSuccess, 0)
	require.NoError(t, err)

	nested, err := NestedMessage(ctx, newNestedTestDelegate(), 0, storageTypes.StatusSuccess, 0, parent.Msg.Id)
	require.NoError(t, err)
	require.Equal(t, parent.Msg.Id, nested.Msg.Id)
	require.Equal(t, storageTypes.MsgDelegate, nested.Msg.Type)
	require.EqualValues(t, nestedTestHeight, nested.Msg.Height)

	require.Equal(t, 2, ctx.AddressMessages.Len())
	for _, ma := range ctx.AddressMessages.All() {
		require.Equal(t, parent.Msg.Id, ma.MsgId)
	}
}

func TestNestedMessage_DoesNotAdvanceMsgCounter(t *testing.T) {
	ctx := newNestedTestContext()

	parent, err := Message(ctx, createMsgUnknown(), 0, storageTypes.StatusSuccess, 0)
	require.NoError(t, err)
	_, err = NestedMessage(ctx, newNestedTestDelegate(), 0, storageTypes.StatusSuccess, 0, parent.Msg.Id)
	require.NoError(t, err)
	next, err := Message(ctx, createMsgUnknown(), 1, storageTypes.StatusSuccess, 0)
	require.NoError(t, err)

	require.Equal(t, parent.Msg.Id+1, next.Msg.Id)
}

// Identical nested messages under one parent collapse into a single
// msg_address row instead of violating its (address, msg, type) key.
func TestNestedMessage_DuplicateAddressesDeduplicated(t *testing.T) {
	ctx := newNestedTestContext()

	for i := range 2 {
		_, err := NestedMessage(ctx, newNestedTestDelegate(), i, storageTypes.StatusSuccess, 0, 42)
		require.NoError(t, err)
	}
	require.Equal(t, 2, ctx.AddressMessages.Len())
}
