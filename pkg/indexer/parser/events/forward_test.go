// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/stretchr/testify/require"
)

func Test_handleForward(t *testing.T) {
	ts := time.Now()

	t.Run("nil event index", func(t *testing.T) {
		ctx := context.NewContext()
		msg := &storage.Message{
			Type:   types.MsgForward,
			Height: 100,
			Time:   ts,
		}
		err := handleForward(ctx, nil, msg, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "nil event index")
	})

	t.Run("nil message", func(t *testing.T) {
		ctx := context.NewContext()
		idx := 0
		err := handleForward(ctx, nil, nil, &idx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "nil message")
	})

	t.Run("unexpected action", func(t *testing.T) {
		ctx := context.NewContext()
		idx := 0
		msg := &storage.Message{
			Type:   types.MsgForward,
			Height: 100,
			Time:   ts,
		}
		events := []storage.Event{
			{
				Height: 100,
				Type:   "message",
				Data: map[string]any{
					"action": "/cosmos.bank.v1beta1.MsgSend",
				},
			},
		}
		err := handleForward(ctx, events, msg, &idx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unexpected event action")
	})

	t.Run("success with token forwarded and complete", func(t *testing.T) {
		ctx := context.NewContext()
		idx := 0
		msg := &storage.Message{
			Type:   types.MsgForward,
			Height: 100,
			Time:   ts,
		}
		events := []storage.Event{
			{
				Height: 100,
				Type:   "message",
				Data: map[string]any{
					"action": "/celestia.forwarding.v1.MsgForward",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
				Data: map[string]any{
					"forward_address": "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
					"denom":           "utia",
					"amount":          "1000",
					"message_id":      "msg-1",
					"success":         "true",
					"error":           "",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
				Data: map[string]any{
					"forward_address": "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
					"denom":           "uatom",
					"amount":          "500",
					"message_id":      "msg-2",
					"success":         "false",
					"error":           "insufficient funds",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventForwardingComplete,
				Data: map[string]any{
					"forward_address":       "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
					"destination_domain":    "1",
					"destination_recipient": "0101",
					"successful_count":      "1",
					"failed_count":          "1",
				},
			},
		}

		err := handleForward(ctx, events, msg, &idx)
		require.NoError(t, err)
		require.NotNil(t, msg.Forwarding)
		require.Equal(t, uint64(1), msg.Forwarding.SuccessCount)
		require.Equal(t, uint64(1), msg.Forwarding.FailedCount)
		require.Equal(t, uint64(1), msg.Forwarding.DestDomain)
		require.NotNil(t, msg.Forwarding.Address)
		require.Equal(t, "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", msg.Forwarding.Address.Address)
		require.True(t, msg.Forwarding.Address.IsForwarding)
		require.NotNil(t, msg.Forwarding.Transfers)
	})

	t.Run("stops at next action event", func(t *testing.T) {
		ctx := context.NewContext()
		idx := 0
		msg := &storage.Message{
			Type:   types.MsgForward,
			Height: 100,
			Time:   ts,
		}
		events := []storage.Event{
			{
				Height: 100,
				Type:   "message",
				Data: map[string]any{
					"action": "/celestia.forwarding.v1.MsgForward",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
				Data: map[string]any{
					"forward_address": "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
					"denom":           "utia",
					"amount":          "1000",
					"message_id":      "msg-1",
					"success":         "true",
					"error":           "",
				},
			},
			{
				Height: 100,
				Type:   "message",
				Data: map[string]any{
					"action": "/cosmos.bank.v1beta1.MsgSend",
				},
			},
		}

		err := handleForward(ctx, events, msg, &idx)
		require.NoError(t, err)
		require.NotNil(t, msg.Forwarding)
		require.Equal(t, 2, idx, "index should stop at next action event")
	})

	t.Run("multiple messages in sequence", func(t *testing.T) {
		ctx := context.NewContext()
		idx := testsuite.Ptr(0)
		events := []storage.Event{
			{
				Height: 100,
				Type:   "message",
				Data: map[string]any{
					"action": "/celestia.forwarding.v1.MsgForward",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
				Data: map[string]any{
					"forward_address": "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
					"denom":           "utia",
					"amount":          "1000",
					"message_id":      "msg-1",
					"success":         "true",
					"error":           "",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventForwardingComplete,
				Data: map[string]any{
					"forward_address":       "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
					"destination_domain":    "1",
					"destination_recipient": "AAEC",
					"successful_count":      "1",
					"failed_count":          "0",
				},
			},
			{
				Height: 100,
				Type:   "message",
				Data: map[string]any{
					"action": "/celestia.forwarding.v1.MsgForward",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
				Data: map[string]any{
					"forward_address": "celestia1ccqy2wlzf2zndn4vspmuksw5frqq0ufsgw4gmt",
					"denom":           "uatom",
					"amount":          "2000",
					"message_id":      "msg-2",
					"success":         "true",
					"error":           "",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventForwardingComplete,
				Data: map[string]any{
					"forward_address":       "celestia1ccqy2wlzf2zndn4vspmuksw5frqq0ufsgw4gmt",
					"destination_domain":    "2",
					"destination_recipient": "010203",
					"successful_count":      "1",
					"failed_count":          "0",
				},
			},
		}

		msgs := []*storage.Message{
			{
				Type:   types.MsgForward,
				Height: 100,
				Time:   ts,
			},
			{
				Type:   types.MsgForward,
				Height: 100,
				Time:   ts,
			},
		}

		for i := range msgs {
			err := handleForward(ctx, events, msgs[i], idx)
			require.NoError(t, err)
			require.NotNil(t, msgs[i].Forwarding)
		}

		require.Equal(t, "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", msgs[0].Forwarding.Address.Address)
		require.Equal(t, "celestia1ccqy2wlzf2zndn4vspmuksw5frqq0ufsgw4gmt", msgs[1].Forwarding.Address.Address)
		require.Equal(t, uint64(1), msgs[0].Forwarding.DestDomain)
		require.Equal(t, uint64(2), msgs[1].Forwarding.DestDomain)
	})

	t.Run("no events after action", func(t *testing.T) {
		ctx := context.NewContext()
		idx := 0
		msg := &storage.Message{
			Type:   types.MsgForward,
			Height: 100,
			Time:   ts,
		}
		events := []storage.Event{
			{
				Height: 100,
				Type:   "message",
				Data: map[string]any{
					"action": "/celestia.forwarding.v1.MsgForward",
				},
			},
		}

		err := handleForward(ctx, events, msg, &idx)
		require.NoError(t, err)
		require.NotNil(t, msg.Forwarding)
		require.Nil(t, msg.Forwarding.Address)
	})

	t.Run("token with error included in transfers", func(t *testing.T) {
		ctx := context.NewContext()
		idx := 0
		msg := &storage.Message{
			Type:   types.MsgForward,
			Height: 100,
			Time:   ts,
		}
		events := []storage.Event{
			{
				Height: 100,
				Type:   "message",
				Data: map[string]any{
					"action": "/celestia.forwarding.v1.MsgForward",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
				Data: map[string]any{
					"forward_address": "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
					"denom":           "utia",
					"amount":          "1000",
					"message_id":      "msg-1",
					"success":         "false",
					"error":           "some error",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventForwardingComplete,
				Data: map[string]any{
					"forward_address":       "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
					"destination_domain":    "1",
					"destination_recipient": "AAEC",
					"successful_count":      "0",
					"failed_count":          "1",
				},
			},
		}

		err := handleForward(ctx, events, msg, &idx)
		require.NoError(t, err)
		require.NotNil(t, msg.Forwarding)
		require.Contains(t, string(msg.Forwarding.Transfers), `"error":"some error"`)
		require.Contains(t, string(msg.Forwarding.Transfers), `"denom":"utia"`)
		require.Contains(t, string(msg.Forwarding.Transfers), `"amount":"1000"`)
	})
}
