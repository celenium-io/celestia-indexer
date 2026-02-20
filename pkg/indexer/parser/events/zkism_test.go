// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/stretchr/testify/require"
)

func makeZkISMContext() *context.Context {
	ctx := context.NewContext()
	ctx.Block = &storage.Block{
		Height: 1_500_000,
		Time:   time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC),
	}
	return ctx
}

// ──────────────────────────────────────────────────────────
// handleCreateZkISM
// ──────────────────────────────────────────────────────────

func Test_handleCreateZkISM_NilIndex(t *testing.T) {
	ctx := makeZkISMContext()
	msg := &storage.Message{Type: types.MsgCreateInterchainSecurityModule}
	err := handleCreateZkISM(ctx, nil, msg, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "nil event index")
}

func Test_handleCreateZkISM_NilMessage(t *testing.T) {
	ctx := makeZkISMContext()
	idx := 0
	err := handleCreateZkISM(ctx, nil, nil, &idx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "nil message")
}

func Test_handleCreateZkISM_WrongAction(t *testing.T) {
	ctx := makeZkISMContext()
	idx := 0
	msg := &storage.Message{Type: types.MsgCreateInterchainSecurityModule}
	events := []storage.Event{
		{
			Height: 1_500_000,
			Type:   "message",
			Data: map[string]any{
				"action": "/cosmos.bank.v1beta1.MsgSend",
			},
		},
	}
	err := handleCreateZkISM(ctx, events, msg, &idx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected event action")
}

func Test_handleCreateZkISM_Success(t *testing.T) {
	stateRoot := testsuite.RandomBytes(32)
	merkleTree := testsuite.RandomBytes(32)
	groth16VKey := testsuite.RandomBytes(64)
	stateTransitionVKey := testsuite.RandomBytes(32)
	stateMembershipVKey := testsuite.RandomBytes(32)
	state := testsuite.RandomBytes(64)

	toHex := func(b []byte) string { return "0x" + hex.EncodeToString(b) }

	ctx := makeZkISMContext()
	idx := 0
	msg := &storage.Message{Type: types.MsgCreateInterchainSecurityModule}
	events := []storage.Event{
		{
			Height: 1_500_000,
			Type:   "message",
			Data: map[string]any{
				"action": "/celestia.zkism.v1.MsgCreateInterchainSecurityModule",
			},
		},
		{
			Height: 1_500_000,
			Type:   types.EventTypeCelestiazkismv1EventCreateInterchainSecurityModule,
			Data: map[string]any{
				"id":                    "42",
				"creator":               "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
				"state":                 toHex(state),
				"state_root":            toHex(stateRoot),
				"merkle_tree_address":   toHex(merkleTree),
				"groth16_vkey":          toHex(groth16VKey),
				"state_transition_vkey": toHex(stateTransitionVKey),
				"state_membership_vkey": toHex(stateMembershipVKey),
			},
		},
	}

	err := handleCreateZkISM(ctx, events, msg, &idx)
	require.NoError(t, err)
	require.NotNil(t, msg.ZkISM)

	ism := msg.ZkISM
	require.EqualValues(t, 42, ism.ExternalId)
	require.Equal(t, ctx.Block.Height, ism.Height)
	require.Equal(t, ctx.Block.Time, ism.Time)
	require.Equal(t, stateRoot, ism.StateRoot)
	require.Equal(t, merkleTree, ism.MerkleTreeAddress)
	require.Equal(t, groth16VKey, ism.Groth16VKey)
	require.Equal(t, stateTransitionVKey, ism.StateTransitionVKey)
	require.Equal(t, stateMembershipVKey, ism.StateMembershipVKey)
	require.NotNil(t, ism.Creator)
	require.Equal(t, "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", ism.Creator.Address)

	// The ISM must be stored in the context keyed by external id
	stored, ok := ctx.ZkISMs.Get(42)
	require.True(t, ok)
	require.Equal(t, ism, stored)
}

func Test_handleCreateZkISM_StopsAtNextAction(t *testing.T) {
	state := testsuite.RandomBytes(32)
	stateRoot := testsuite.RandomBytes(32)
	merkleTree := testsuite.RandomBytes(32)
	groth16VKey := testsuite.RandomBytes(32)
	stateTransVKey := testsuite.RandomBytes(32)
	stateMembVKey := testsuite.RandomBytes(32)

	toHex := func(b []byte) string { return "0x" + hex.EncodeToString(b) }

	ctx := makeZkISMContext()
	idx := testsuite.Ptr(0)
	events := []storage.Event{
		{
			Type: "message",
			Data: map[string]any{"action": "/celestia.zkism.v1.MsgCreateInterchainSecurityModule"},
		},
		{
			Type: types.EventTypeCelestiazkismv1EventCreateInterchainSecurityModule,
			Data: map[string]any{
				"id":                    "1",
				"creator":               "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
				"state":                 toHex(state),
				"state_root":            toHex(stateRoot),
				"merkle_tree_address":   toHex(merkleTree),
				"groth16_vkey":          toHex(groth16VKey),
				"state_transition_vkey": toHex(stateTransVKey),
				"state_membership_vkey": toHex(stateMembVKey),
			},
		},
		{
			// second message — should NOT be consumed
			Type: "message",
			Data: map[string]any{"action": "/celestia.zkism.v1.MsgUpdateInterchainSecurityModule"},
		},
	}

	msg := &storage.Message{Type: types.MsgCreateInterchainSecurityModule}
	err := handleCreateZkISM(ctx, events, msg, idx)
	require.NoError(t, err)
	require.NotNil(t, msg.ZkISM)
	require.Equal(t, 2, *idx, "index must stop before the next action event")
}

// ──────────────────────────────────────────────────────────
// handleUpdateZkISM
// ──────────────────────────────────────────────────────────

func Test_handleUpdateZkISM_NilIndex(t *testing.T) {
	ctx := makeZkISMContext()
	msg := &storage.Message{Type: types.MsgUpdateInterchainSecurityModule}
	err := handleUpdateZkISM(ctx, nil, msg, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "nil event index")
}

func Test_handleUpdateZkISM_NilMessage(t *testing.T) {
	ctx := makeZkISMContext()
	idx := 0
	err := handleUpdateZkISM(ctx, nil, nil, &idx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "nil message")
}

func Test_handleUpdateZkISM_WrongAction(t *testing.T) {
	ctx := makeZkISMContext()
	idx := 0
	msg := &storage.Message{Type: types.MsgUpdateInterchainSecurityModule}
	events := []storage.Event{
		{
			Type: "message",
			Data: map[string]any{"action": "/cosmos.bank.v1beta1.MsgSend"},
		},
	}
	err := handleUpdateZkISM(ctx, events, msg, &idx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected event action")
}

func Test_handleUpdateZkISM_Success(t *testing.T) {
	newState := testsuite.RandomBytes(64)
	newStateRoot := testsuite.RandomBytes(32)
	toHex := func(b []byte) string { return "0x" + hex.EncodeToString(b) }

	ctx := makeZkISMContext()
	idx := 0
	msg := &storage.Message{Type: types.MsgUpdateInterchainSecurityModule}
	events := []storage.Event{
		{
			Type: "message",
			Data: map[string]any{"action": "/celestia.zkism.v1.MsgUpdateInterchainSecurityModule"},
		},
		{
			Type: types.EventTypeCelestiazkismv1EventUpdateInterchainSecurityModule,
			Data: map[string]any{
				"id":             "7",
				"signer":         "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
				"new_state":      toHex(newState),
				"new_state_root": toHex(newStateRoot),
			},
		},
	}

	err := handleUpdateZkISM(ctx, events, msg, &idx)
	require.NoError(t, err)
	require.NotNil(t, msg.ZkISMUpdate)

	upd := msg.ZkISMUpdate
	require.EqualValues(t, 7, upd.ZkISMExternalId)
	require.Equal(t, ctx.Block.Height, upd.Height)
	require.Equal(t, ctx.Block.Time, upd.Time)
	require.Equal(t, newState, upd.NewState)
	require.Equal(t, newStateRoot, upd.NewStateRoot)
	require.NotNil(t, upd.Signer)
	require.Equal(t, "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj", upd.Signer.Address)
}

func Test_handleUpdateZkISM_UpdatesContextState(t *testing.T) {
	// If the ISM was created in the same block it must be updated in the context.
	oldState := testsuite.RandomBytes(64)
	oldRoot := testsuite.RandomBytes(32)
	newState := testsuite.RandomBytes(64)
	newRoot := testsuite.RandomBytes(32)
	toHex := func(b []byte) string { return "0x" + hex.EncodeToString(b) }

	ctx := makeZkISMContext()
	// Simulate ISM created earlier in this block
	ctx.ZkISMs.Set(7, &storage.ZkISM{ExternalId: 7, State: oldState, StateRoot: oldRoot})

	idx := 0
	msg := &storage.Message{Type: types.MsgUpdateInterchainSecurityModule}
	events := []storage.Event{
		{
			Type: "message",
			Data: map[string]any{"action": "/celestia.zkism.v1.MsgUpdateInterchainSecurityModule"},
		},
		{
			Type: types.EventTypeCelestiazkismv1EventUpdateInterchainSecurityModule,
			Data: map[string]any{
				"id":             "7",
				"signer":         "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
				"new_state":      toHex(newState),
				"new_state_root": toHex(newRoot),
			},
		},
	}

	err := handleUpdateZkISM(ctx, events, msg, &idx)
	require.NoError(t, err)

	stored, ok := ctx.ZkISMs.Get(7)
	require.True(t, ok)
	require.Equal(t, newState, stored.State)
	require.Equal(t, newRoot, stored.StateRoot)
}

// ──────────────────────────────────────────────────────────
// handleSubmitZkISMMessages
// ──────────────────────────────────────────────────────────

func Test_handleSubmitZkISMMessages_NilIndex(t *testing.T) {
	ctx := makeZkISMContext()
	msg := &storage.Message{Type: types.MsgSubmitMessages}
	err := handleSubmitZkISMMessages(ctx, nil, msg, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "nil event index")
}

func Test_handleSubmitZkISMMessages_NilMessage(t *testing.T) {
	ctx := makeZkISMContext()
	idx := 0
	err := handleSubmitZkISMMessages(ctx, nil, nil, &idx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "nil message")
}

func Test_handleSubmitZkISMMessages_WrongAction(t *testing.T) {
	ctx := makeZkISMContext()
	idx := 0
	msg := &storage.Message{Type: types.MsgSubmitMessages}
	events := []storage.Event{
		{
			Type: "message",
			Data: map[string]any{"action": "/cosmos.bank.v1beta1.MsgSend"},
		},
	}
	err := handleSubmitZkISMMessages(ctx, events, msg, &idx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected event action")
}

func Test_handleSubmitZkISMMessages_SingleMessage(t *testing.T) {
	stateRoot := testsuite.RandomBytes(32)
	msgId := testsuite.RandomBytes(32)
	toHex := func(b []byte) string { return "0x" + hex.EncodeToString(b) }

	ctx := makeZkISMContext()
	idx := 0
	msg := &storage.Message{Type: types.MsgSubmitMessages}
	events := []storage.Event{
		{
			Type: "message",
			Data: map[string]any{"action": "/celestia.zkism.v1.MsgSubmitMessages"},
		},
		{
			Type: types.EventTypeCelestiazkismv1EventSubmitMessages,
			Data: map[string]any{
				"id":          "3",
				"signer":      "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
				"state_root":  toHex(stateRoot),
				"message_ids": toHex(msgId),
			},
		},
	}

	err := handleSubmitZkISMMessages(ctx, events, msg, &idx)
	require.NoError(t, err)
	require.Len(t, msg.ZkISMMessages, 1)

	m := msg.ZkISMMessages[0]
	require.EqualValues(t, 3, m.ZkISMExternalId)
	require.Equal(t, ctx.Block.Height, m.Height)
	require.Equal(t, ctx.Block.Time, m.Time)
	require.Equal(t, stateRoot, m.StateRoot)
	require.Equal(t, msgId, m.MessageId)
	require.NotNil(t, m.Signer)
	require.Equal(t, "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", m.Signer.Address)
}

func Test_handleSubmitZkISMMessages_MultipleMessages(t *testing.T) {
	stateRoot := testsuite.RandomBytes(32)
	msgId1 := testsuite.RandomBytes(32)
	msgId2 := testsuite.RandomBytes(32)
	msgId3 := testsuite.RandomBytes(32)
	toHex := func(b []byte) string { return "0x" + hex.EncodeToString(b) }

	ctx := makeZkISMContext()
	idx := 0
	msg := &storage.Message{Type: types.MsgSubmitMessages}
	events := []storage.Event{
		{
			Type: "message",
			Data: map[string]any{"action": "/celestia.zkism.v1.MsgSubmitMessages"},
		},
		{
			Type: types.EventTypeCelestiazkismv1EventSubmitMessages,
			Data: map[string]any{
				"id":          "5",
				"signer":      "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
				"state_root":  toHex(stateRoot),
				"message_ids": toHex(msgId1) + "," + toHex(msgId2) + "," + toHex(msgId3),
			},
		},
	}

	err := handleSubmitZkISMMessages(ctx, events, msg, &idx)
	require.NoError(t, err)
	require.Len(t, msg.ZkISMMessages, 3)

	for _, m := range msg.ZkISMMessages {
		require.EqualValues(t, 5, m.ZkISMExternalId)
		require.Equal(t, stateRoot, m.StateRoot)
		require.NotNil(t, m.Signer)
	}
	require.Equal(t, msgId1, msg.ZkISMMessages[0].MessageId)
	require.Equal(t, msgId2, msg.ZkISMMessages[1].MessageId)
	require.Equal(t, msgId3, msg.ZkISMMessages[2].MessageId)
}

func Test_handleSubmitZkISMMessages_SequentialMessages(t *testing.T) {
	stateRoot1 := testsuite.RandomBytes(32)
	stateRoot2 := testsuite.RandomBytes(32)
	msgId1 := testsuite.RandomBytes(32)
	msgId2 := testsuite.RandomBytes(32)
	toHex := func(b []byte) string { return "0x" + hex.EncodeToString(b) }

	ctx := makeZkISMContext()
	idx := testsuite.Ptr(0)
	events := []storage.Event{
		{
			Type: "message",
			Data: map[string]any{"action": "/celestia.zkism.v1.MsgSubmitMessages"},
		},
		{
			Type: types.EventTypeCelestiazkismv1EventSubmitMessages,
			Data: map[string]any{
				"id":          "1",
				"signer":      "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
				"state_root":  toHex(stateRoot1),
				"message_ids": toHex(msgId1),
			},
		},
		{
			Type: "message",
			Data: map[string]any{"action": "/celestia.zkism.v1.MsgSubmitMessages"},
		},
		{
			Type: types.EventTypeCelestiazkismv1EventSubmitMessages,
			Data: map[string]any{
				"id":          "2",
				"signer":      "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj",
				"state_root":  toHex(stateRoot2),
				"message_ids": toHex(msgId2),
			},
		},
	}

	msgs := []*storage.Message{
		{Type: types.MsgSubmitMessages},
		{Type: types.MsgSubmitMessages},
	}

	for i := range msgs {
		err := handleSubmitZkISMMessages(ctx, events, msgs[i], idx)
		require.NoError(t, err)
	}

	require.Len(t, msgs[0].ZkISMMessages, 1)
	require.EqualValues(t, 1, msgs[0].ZkISMMessages[0].ZkISMExternalId)
	require.Equal(t, msgId1, msgs[0].ZkISMMessages[0].MessageId)
	require.Equal(t, stateRoot1, msgs[0].ZkISMMessages[0].StateRoot)

	require.Len(t, msgs[1].ZkISMMessages, 1)
	require.EqualValues(t, 2, msgs[1].ZkISMMessages[0].ZkISMExternalId)
	require.Equal(t, msgId2, msgs[1].ZkISMMessages[0].MessageId)
	require.Equal(t, stateRoot2, msgs[1].ZkISMMessages[0].StateRoot)
}
