// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icaTypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	"github.com/stretchr/testify/require"
)

// An unregistered nested type followed by a registered one must not pull
// the nested handlers into the next top-level message's events.
func Test_handleRecvPacket_NestedDoesNotOverrunTopLevel(t *testing.T) {
	const (
		delegator = "celestia1lkrd86urrmhmsvgzfygjsguv3cgv0036hrj0m9"
		validator = "celestiavaloper1qe8uuf5x69c526h4nzxwv4ltftr73v7q5qhs58"
	)

	ctx := context.NewContext()
	ctx.Block = &storage.Block{Height: 100, Time: time.Now().UTC()}

	events := []storage.Event{
		{Type: types.EventTypeMessage, Data: map[string]string{"action": "/ibc.core.channel.v1.MsgRecvPacket", "msg_index": "0"}},
		{Type: types.EventTypeRecvPacket, Data: map[string]string{"packet_dst_port": "icahost", "msg_index": "0"}},
		{Type: types.EventTypeMessage, Data: map[string]string{"module": "ibc_channel", "msg_index": "0"}},
		// nested MsgSetWithdrawAddress (no event handler) without a closing module event
		{Type: types.EventTypeSetWithdrawAddress, Data: map[string]string{"withdraw_address": delegator, "msg_index": "0"}},
		// nested MsgDelegate
		{Type: types.EventTypeDelegate, Data: map[string]string{"validator": validator, "delegator": delegator, "amount": "1utia", "new_shares": "1", "msg_index": "0"}},
		{Type: types.EventTypeIcs27Packet, Data: map[string]string{"module": "interchainaccounts", "success": "true", "msg_index": "0"}},
		{Type: types.EventTypeWriteAcknowledgement, Data: map[string]string{"packet_dst_port": "icahost", "msg_index": "0"}},
		// next top-level MsgDelegate
		{Type: types.EventTypeMessage, Data: map[string]string{"action": "/cosmos.staking.v1beta1.MsgDelegate", "module": "staking", "sender": delegator, "msg_index": "1"}},
		{Type: types.EventTypeDelegate, Data: map[string]string{"validator": validator, "delegator": delegator, "amount": "999utia", "new_shares": "999", "msg_index": "1"}},
	}

	msg := &storage.Message{
		Id:     1,
		Type:   types.MsgRecvPacket,
		Height: 100,
		Data: map[string]any{
			"Packet": map[string]any{
				"DestinationPort": "icahost",
				"Data": map[string]any{
					"Type": icaTypes.EXECUTE_TX,
					"Data": []cosmosTypes.Msg{
						&cosmosDistributionTypes.MsgSetWithdrawAddress{
							DelegatorAddress: delegator,
							WithdrawAddress:  delegator,
						},
						&cosmosStakingTypes.MsgDelegate{
							DelegatorAddress: delegator,
							ValidatorAddress: validator,
							Amount:           cosmosTypes.NewCoin("utia", math.NewInt(1)),
						},
					},
				},
			},
		},
	}

	c := NewCursor(events)
	require.NoError(t, handleRecvPacket(ctx, c, msg))

	next, ok := c.Peek()
	require.True(t, ok)
	require.Equal(t, "/cosmos.staking.v1beta1.MsgDelegate", next.Data["action"])

	// nested messages aren't stored: their addresses must point to the RecvPacket
	require.Positive(t, ctx.AddressMessages.Len())
	for _, ma := range ctx.AddressMessages.All() {
		require.Equal(t, msg.Id, ma.MsgId)
	}
	require.Empty(t, ctx.Messages)

	for _, d := range ctx.Delegations.All() {
		require.NotEqual(t, "999", d.Amount.String(), "top-level delegation consumed by nested handler")
	}
}
