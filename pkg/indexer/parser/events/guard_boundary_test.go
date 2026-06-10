// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/stretchr/testify/require"
)

// msgEvent builds a single "message" event carrying the dispatch action.
func msgEvent(action string) storage.Event {
	return storage.Event{
		Height: 100,
		Type:   types.EventTypeMessage,
		Data:   map[string]string{"action": action},
	}
}

type guardHandlerFn func(*context.Context, []storage.Event, *storage.Message, *int) error

// Test_eventHandlers_truncatedEvents_returnErrorNotPanic verifies that every
// event handler which advances *idx by a fixed offset guards the subsequent
// events[*idx] access. When the node returns fewer events than the walker
// expects, the handler must return an error instead of panicking with an
// index-out-of-range. Each case feeds only the leading "message" event, so the
// guard right after the increment is the code path under test.
func Test_eventHandlers_truncatedEvents_returnErrorNotPanic(t *testing.T) {
	tests := []struct {
		name    string
		handler guardHandlerFn
		events  []storage.Event
		msg     *storage.Message
	}{
		{
			name:    "channel_close",
			handler: handleChannelClose,
			events:  []storage.Event{msgEvent("/ibc.core.channel.v1.MsgChannelCloseConfirm")},
			msg:     &storage.Message{Type: types.MsgUnknown},
		},
		{
			name:    "channel_open_confirm",
			handler: handleChannelOpenConfirm,
			events:  []storage.Event{msgEvent("/ibc.core.channel.v1.MsgChannelOpenConfirm")},
			msg:     &storage.Message{Type: types.MsgUnknown},
		},
		{
			name:    "channel_open_init",
			handler: handleChannelOpenInit,
			events:  []storage.Event{msgEvent("/ibc.core.channel.v1.MsgChannelOpenInit")},
			msg:     &storage.Message{Type: types.MsgUnknown},
		},
		{
			name:    "create_client",
			handler: handleCreateClient,
			events:  []storage.Event{msgEvent("/ibc.core.client.v1.MsgCreateClient")},
			msg:     &storage.Message{Type: types.MsgUnknown},
		},
		{
			name:    "update_client",
			handler: handleUpdateClient,
			events:  []storage.Event{msgEvent("/ibc.core.client.v1.MsgUpdateClient")},
			msg:     &storage.Message{Type: types.MsgUnknown},
		},
		{
			name:    "send",
			handler: handleSend,
			events:  []storage.Event{msgEvent("/cosmos.bank.v1beta1.MsgSend")},
			msg:     &storage.Message{Type: types.MsgUnknown},
		},
		{
			name:    "vote",
			handler: handleVote,
			events:  []storage.Event{msgEvent("/cosmos.gov.v1.MsgVote")},
			msg:     &storage.Message{Type: types.MsgUnknown},
		},
		{
			name:    "submit_proposal first guard",
			handler: handleSubmitProposal,
			events:  []storage.Event{msgEvent("/cosmos.gov.v1.MsgSubmitProposal")},
			msg:     &storage.Message{Type: types.MsgUnknown, Proposal: &storage.Proposal{}},
		},
		{
			name:    "deposit",
			handler: handleDeposit,
			events:  []storage.Event{msgEvent("/cosmos.gov.v1.MsgDeposit")},
			msg:     &storage.Message{Type: types.MsgUnknown},
		},
		{
			// module is empty on the first event, so the handler advances *idx by
			// one and reads the next (missing) event — the inner guard must fire.
			name:    "unjail inner guard",
			handler: handleUnjail,
			events:  []storage.Event{msgEvent("/cosmos.slashing.v1beta1.MsgUnjail")},
			msg:     &storage.Message{Type: types.MsgUnknown},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := 0
			require.NotPanics(t, func() {
				err := tt.handler(context.NewContext(), tt.events, tt.msg, &idx)
				require.Error(t, err)
			})
		})
	}
}

// Test_submitProposal_secondGuard_returnErrorNotPanic exercises the second guard
// in processSubmitProposal: after the proposal id is parsed, *idx jumps by 5.
// The events slice here is long enough to pass the first guard and the type
// check, but too short for the +5 jump, so the second guard must return an
// error rather than panic.
func Test_submitProposal_secondGuard_returnErrorNotPanic(t *testing.T) {
	events := []storage.Event{
		msgEvent("/cosmos.gov.v1.MsgSubmitProposal"),
		{
			Height: 100,
			Type:   types.EventTypeSubmitProposal,
			Data:   map[string]string{"proposal_id": "1"},
		},
	}
	msg := &storage.Message{Type: types.MsgUnknown, Proposal: &storage.Proposal{}}

	idx := 0
	require.NotPanics(t, func() {
		err := handleSubmitProposal(context.NewContext(), events, msg, &idx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "after parsing proposal id")
	})
}
