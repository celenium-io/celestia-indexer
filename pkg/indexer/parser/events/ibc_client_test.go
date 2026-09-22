// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	cmtTypes "github.com/cometbft/cometbft/proto/tendermint/types"
	clientTypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	solomachine "github.com/cosmos/ibc-go/v8/modules/light-clients/06-solomachine"
	tendermint "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	"github.com/stretchr/testify/require"
)

const (
	actionCreateClient = "/ibc.core.client.v1.MsgCreateClient"
	actionUpdateClient = "/ibc.core.client.v1.MsgUpdateClient"
)

var clientTime = time.Date(2026, 9, 22, 0, 0, 0, 0, time.UTC)

func clientEvent(typ types.EventType, data map[string]string) storage.Event {
	return storage.Event{Height: 100, Type: typ, Data: data}
}

func ibcClientModuleEvent() storage.Event {
	return clientEvent(types.EventTypeMessage, map[string]string{"module": "ibc_client"})
}

func createClientEvents(id, typ, height string) []storage.Event {
	return []storage.Event{
		msgEvent(actionCreateClient),
		clientEvent(types.EventTypeCreateClient, map[string]string{
			"client_id":        id,
			"client_type":      typ,
			"consensus_height": height,
		}),
		ibcClientModuleEvent(),
	}
}

func updateClientEvents(id, typ, height string) []storage.Event {
	return []storage.Event{
		msgEvent(actionUpdateClient),
		clientEvent(types.EventTypeUpdateClient, map[string]string{
			"client_id":         id,
			"client_type":       typ,
			"consensus_height":  height,
			"consensus_heights": height,
		}),
		ibcClientModuleEvent(),
	}
}

func misbehaviourEvents(id, typ string) []storage.Event {
	return []storage.Event{
		msgEvent(actionUpdateClient),
		clientEvent(types.EventTypeClientMisbehaviour, map[string]string{
			"client_id":   id,
			"client_type": typ,
		}),
		ibcClientModuleEvent(),
	}
}

func clientMsg(typ types.MsgType, data map[string]any) *storage.Message {
	return &storage.Message{Type: typ, Height: 100, Time: clientTime, TxId: 1, Data: data}
}

// runClientTx feeds msgs through Handle over one shared cursor, as parseTxs does.
func runClientTx(t *testing.T, ctx *context.Context, evts []storage.Event, msgs ...*storage.Message) {
	t.Helper()
	c := NewCursor(evts)
	for _, msg := range msgs {
		require.NoError(t, Handle(ctx, c, msg), "msg %s", msg.Type)
	}
	_, ok := c.Peek()
	require.False(t, ok, "all events must be consumed")
}

func tmClientState() tendermint.ClientState {
	return tendermint.ClientState{
		ChainId:         "osmosis-1",
		TrustLevel:      tendermint.Fraction{Numerator: 1, Denominator: 3},
		TrustingPeriod:  time.Hour,
		UnbondingPeriod: 2 * time.Hour,
		MaxClockDrift:   time.Second,
		LatestHeight:    clientTypes.NewHeight(1, 100),
	}
}

func tmHeader() tendermint.Header {
	return tendermint.Header{
		SignedHeader: &cmtTypes.SignedHeader{
			Header: &cmtTypes.Header{ChainID: "osmosis-1"},
		},
	}
}

func TestCreateClient_Tendermint(t *testing.T) {
	ctx := context.NewContext()
	msg := clientMsg(types.MsgCreateClient, map[string]any{
		"Signer":      "celestia1signer",
		"ClientState": tmClientState(),
	})
	runClientTx(t, ctx, createClientEvents("07-tendermint-1", "07-tendermint", "1-100"), msg)

	client, ok := ctx.IbcClients.Get("07-tendermint-1")
	require.True(t, ok)
	require.Equal(t, "07-tendermint", client.Type)
	require.Equal(t, time.Hour, client.TrustingPeriod)
	require.EqualValues(t, 100, client.LatestRevisionHeight)
	require.EqualValues(t, 1, client.LatestRevisionNumber)
	require.EqualValues(t, 3, client.TrustLevelDenominator)
	require.Equal(t, "celestia1signer", client.Creator.Address)
}

func TestCreateClient_Solomachine(t *testing.T) {
	ctx := context.NewContext()
	msg := clientMsg(types.MsgCreateClient, map[string]any{
		"Signer":      "celestia1signer",
		"ClientState": solomachine.ClientState{Sequence: 7},
	})
	runClientTx(t, ctx, createClientEvents("06-solomachine-1", "06-solomachine", "0-7"), msg)

	client, ok := ctx.IbcClients.Get("06-solomachine-1")
	require.True(t, ok)
	require.Equal(t, "06-solomachine", client.Type)
	require.EqualValues(t, 7, client.LatestRevisionHeight)
	require.Zero(t, client.TrustingPeriod)
	require.Empty(t, client.ChainId)
	require.Equal(t, clientTime, client.CreatedAt)
}

// The cursor must stop at the next message's action, whatever the client type.
func TestCreateClient_FollowedByUpdateInSameTx(t *testing.T) {
	tests := []struct {
		name  string
		id    string
		typ   string
		state any
		hdr   any
	}{
		{"tendermint", "07-tendermint-1", "07-tendermint", tmClientState(), tmHeader()},
		{"solomachine", "06-solomachine-1", "06-solomachine", solomachine.ClientState{Sequence: 1}, solomachine.Header{Timestamp: 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.NewContext()
			evts := append(createClientEvents(tt.id, tt.typ, "1-100"), updateClientEvents(tt.id, tt.typ, "1-101")...)
			runClientTx(t, ctx, evts,
				clientMsg(types.MsgCreateClient, map[string]any{"ClientState": tt.state}),
				clientMsg(types.MsgUpdateClient, map[string]any{"Header": tt.hdr}),
			)
			client, ok := ctx.IbcClients.Get(tt.id)
			require.True(t, ok)
			require.EqualValues(t, 101, client.LatestRevisionHeight)
		})
	}
}

func TestUpdateClient(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		evts       []storage.Event
		data       map[string]any
		wantHeight uint64
		wantChain  string
		wantFrozen uint64
	}{
		{
			name:       "tendermint header, update",
			id:         "07-tendermint-1",
			evts:       updateClientEvents("07-tendermint-1", "07-tendermint", "1-150"),
			data:       map[string]any{"Header": tmHeader()},
			wantHeight: 150,
			wantChain:  "osmosis-1",
		}, {
			name:       "solomachine header, update",
			id:         "06-solomachine-1",
			evts:       updateClientEvents("06-solomachine-1", "06-solomachine", "0-8"),
			data:       map[string]any{"Header": solomachine.Header{Timestamp: 5}},
			wantHeight: 8,
		}, {
			name:       "tendermint header, misbehaviour",
			id:         "07-tendermint-1",
			evts:       misbehaviourEvents("07-tendermint-1", "07-tendermint"),
			data:       map[string]any{"Header": tmHeader()},
			wantFrozen: 1,
		}, {
			name:       "tendermint misbehaviour, misbehaviour",
			id:         "07-tendermint-1",
			evts:       misbehaviourEvents("07-tendermint-1", "07-tendermint"),
			data:       map[string]any{"Misbehaviour": tendermint.Misbehaviour{}},
			wantFrozen: 1,
		}, {
			name:       "solomachine misbehaviour, misbehaviour",
			id:         "06-solomachine-1",
			evts:       misbehaviourEvents("06-solomachine-1", "06-solomachine"),
			data:       map[string]any{"Misbehaviour": solomachine.Misbehaviour{}},
			wantFrozen: 1,
		}, {
			// invalid tm misbehaviour: UpdateState returns no heights, update_client has empty consensus_height
			name: "tendermint misbehaviour, update without heights",
			id:   "07-tendermint-1",
			evts: updateClientEvents("07-tendermint-1", "07-tendermint", ""),
			data: map[string]any{"Misbehaviour": tendermint.Misbehaviour{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.NewContext()
			runClientTx(t, ctx, tt.evts, clientMsg(types.MsgUpdateClient, tt.data))

			client, ok := ctx.IbcClients.Get(tt.id)
			require.True(t, ok)
			require.Equal(t, tt.wantHeight, client.LatestRevisionHeight)
			require.Equal(t, tt.wantChain, client.ChainId)
			require.Equal(t, tt.wantFrozen, client.FrozenRevisionHeight)
			require.Equal(t, clientTime, client.UpdatedAt)
		})
	}
}

func TestUpdateClient_UnexpectedEvent(t *testing.T) {
	evts := []storage.Event{
		msgEvent(actionUpdateClient),
		clientEvent(types.EventTypeCreateClient, map[string]string{"client_id": "07-tendermint-1"}),
		ibcClientModuleEvent(),
	}
	err := Handle(context.NewContext(), NewCursor(evts), clientMsg(types.MsgUpdateClient, map[string]any{}))
	require.Error(t, err)
}

// Relayers often send update + misbehaviour in one block: the freeze must survive the merge.
func TestUpdateClient_UpdateThenMisbehaviourInOneBlock(t *testing.T) {
	ctx := context.NewContext()
	id := "07-tendermint-1"
	evts := append(updateClientEvents(id, "07-tendermint", "1-150"), misbehaviourEvents(id, "07-tendermint")...)
	runClientTx(t, ctx, evts,
		clientMsg(types.MsgUpdateClient, map[string]any{"Header": tmHeader()}),
		clientMsg(types.MsgUpdateClient, map[string]any{"Misbehaviour": tendermint.Misbehaviour{}}),
	)

	client, ok := ctx.IbcClients.Get(id)
	require.True(t, ok)
	require.EqualValues(t, 150, client.LatestRevisionHeight)
	require.EqualValues(t, 1, client.FrozenRevisionHeight)
	require.Equal(t, "osmosis-1", client.ChainId)
}

const actionUpgradeClient = "/ibc.core.client.v1.MsgUpgradeClient"

func upgradeClientEvents(id, height string) []storage.Event {
	return []storage.Event{
		msgEvent(actionUpgradeClient),
		clientEvent(types.EventTypeUpgradeClient, map[string]string{
			"client_id":        id,
			"client_type":      "07-tendermint",
			"consensus_height": height,
		}),
		ibcClientModuleEvent(),
	}
}

func TestUpgradeClient(t *testing.T) {
	ctx := context.NewContext()
	id := "07-tendermint-1"
	upgraded := tendermint.ClientState{
		ChainId:         "osmosis-2",
		UnbondingPeriod: 3 * time.Hour,
		LatestHeight:    clientTypes.NewHeight(2, 1),
	}
	runClientTx(t, ctx, upgradeClientEvents(id, "2-1"),
		clientMsg(types.MsgUpgradeClient, map[string]any{"ClientState": upgraded}),
	)

	client, ok := ctx.IbcClients.Get(id)
	require.True(t, ok)
	require.EqualValues(t, 2, client.LatestRevisionNumber)
	require.EqualValues(t, 1, client.LatestRevisionHeight)
	require.Equal(t, "osmosis-2", client.ChainId)
	require.Equal(t, 3*time.Hour, client.UnbondingPeriod)
	require.Equal(t, clientTime, client.UpdatedAt)
}

// update on the old revision, then upgrade in the same block: the new revision wins
func TestUpgradeClient_AfterUpdateInSameBlock(t *testing.T) {
	ctx := context.NewContext()
	id := "07-tendermint-1"
	evts := append(updateClientEvents(id, "07-tendermint", "1-1000"), upgradeClientEvents(id, "2-1")...)
	runClientTx(t, ctx, evts,
		clientMsg(types.MsgUpdateClient, map[string]any{"Header": tmHeader()}),
		clientMsg(types.MsgUpgradeClient, map[string]any{"ClientState": tendermint.ClientState{
			ChainId:         "osmosis-2",
			UnbondingPeriod: 3 * time.Hour,
		}}),
	)

	client, ok := ctx.IbcClients.Get(id)
	require.True(t, ok)
	require.EqualValues(t, 2, client.LatestRevisionNumber)
	require.EqualValues(t, 1, client.LatestRevisionHeight)
	require.Equal(t, "osmosis-2", client.ChainId)
	require.Equal(t, 3*time.Hour, client.UnbondingPeriod)
}

func TestUpgradeClient_UnexpectedEvent(t *testing.T) {
	evts := []storage.Event{
		msgEvent(actionUpgradeClient),
		clientEvent(types.EventTypeUpdateClient, map[string]string{"client_id": "07-tendermint-1"}),
		ibcClientModuleEvent(),
	}
	err := Handle(context.NewContext(), NewCursor(evts), clientMsg(types.MsgUpgradeClient, map[string]any{}))
	require.Error(t, err)
}
