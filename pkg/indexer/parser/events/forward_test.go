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
				Data: map[string]string{
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
				Data: map[string]string{
					"action": "/celestia.forwarding.v1.MsgForward",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
				Data: map[string]string{
					"forward_addr": "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
					"denom":        "utia",
					"amount":       "1000",
					"message_id":   "msg-1",
					"success":      "true",
					"error":        "",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
				Data: map[string]string{
					"forward_addr": "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
					"denom":        "uatom",
					"amount":       "500",
					"message_id":   "msg-2",
					"success":      "false",
					"error":        "insufficient funds",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventForwardingComplete,
				Data: map[string]string{
					"forward_addr":     "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
					"dest_domain":      "1",
					"dest_recipient":   "0101",
					"tokens_forwarded": "1",
					"tokens_failed":    "1",
				},
			},
		}

		err := handleForward(ctx, events, msg, &idx)
		require.NoError(t, err)
		require.Len(t, ctx.Forwardings, 1)
		require.NotNil(t, ctx.Forwardings[0])
		require.Equal(t, uint64(1), ctx.Forwardings[0].SuccessCount)
		require.Equal(t, uint64(1), ctx.Forwardings[0].FailedCount)
		require.Equal(t, uint64(1), ctx.Forwardings[0].DestDomain)
		require.NotNil(t, ctx.Forwardings[0].Address)
		require.Equal(t, "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", ctx.Forwardings[0].Address.Address)
		require.True(t, ctx.Forwardings[0].Address.IsForwarding)
		require.NotNil(t, ctx.Forwardings[0].Transfers)
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
				Data: map[string]string{
					"action": "/celestia.forwarding.v1.MsgForward",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
				Data: map[string]string{
					"forward_addr": "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
					"denom":        "utia",
					"amount":       "1000",
					"message_id":   "msg-1",
					"success":      "true",
					"error":        "",
				},
			},
			{
				Height: 100,
				Type:   "message",
				Data: map[string]string{
					"action": "/cosmos.bank.v1beta1.MsgSend",
				},
			},
		}

		err := handleForward(ctx, events, msg, &idx)
		require.NoError(t, err)
		require.Len(t, ctx.Forwardings, 1)
		require.NotNil(t, ctx.Forwardings[0])
		require.Equal(t, 2, idx, "index should stop at next action event")
	})

	t.Run("multiple messages in sequence", func(t *testing.T) {
		ctx := context.NewContext()
		idx := testsuite.Ptr(0)
		events := []storage.Event{
			{
				Height: 100,
				Type:   "message",
				Data: map[string]string{
					"action": "/celestia.forwarding.v1.MsgForward",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
				Data: map[string]string{
					"forward_addr": "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
					"denom":        "utia",
					"amount":       "1000",
					"message_id":   "msg-1",
					"success":      "true",
					"error":        "",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventForwardingComplete,
				Data: map[string]string{
					"forward_addr":     "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
					"dest_domain":      "1",
					"dest_recipient":   "AAEC",
					"tokens_forwarded": "1",
					"tokens_failed":    "0",
				},
			},
			{
				Height: 100,
				Type:   "message",
				Data: map[string]string{
					"action": "/celestia.forwarding.v1.MsgForward",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
				Data: map[string]string{
					"forward_addr": "celestia1ccqy2wlzf2zndn4vspmuksw5frqq0ufsgw4gmt",
					"denom":        "uatom",
					"amount":       "2000",
					"message_id":   "msg-2",
					"success":      "true",
					"error":        "",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventForwardingComplete,
				Data: map[string]string{
					"forward_addr":     "celestia1ccqy2wlzf2zndn4vspmuksw5frqq0ufsgw4gmt",
					"dest_domain":      "2",
					"dest_recipient":   "010203",
					"tokens_forwarded": "1",
					"tokens_failed":    "0",
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
		}

		require.Len(t, ctx.Forwardings, 2)
		require.NotNil(t, ctx.Forwardings[0])
		require.Equal(t, "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", ctx.Forwardings[0].Address.Address)
		require.Equal(t, "celestia1ccqy2wlzf2zndn4vspmuksw5frqq0ufsgw4gmt", ctx.Forwardings[1].Address.Address)
		require.Equal(t, uint64(1), ctx.Forwardings[0].DestDomain)
		require.Equal(t, uint64(2), ctx.Forwardings[1].DestDomain)
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
				Data: map[string]string{
					"action": "/celestia.forwarding.v1.MsgForward",
				},
			},
		}

		err := handleForward(ctx, events, msg, &idx)
		require.NoError(t, err)
		require.Len(t, ctx.Forwardings, 1)
		require.NotNil(t, ctx.Forwardings[0])
		require.Nil(t, ctx.Forwardings[0].Address)
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
				Data: map[string]string{
					"action": "/celestia.forwarding.v1.MsgForward",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
				Data: map[string]string{
					"forward_addr": "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
					"denom":        "utia",
					"amount":       "1000",
					"message_id":   "msg-1",
					"success":      "false",
					"error":        "some error",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventForwardingComplete,
				Data: map[string]string{
					"forward_addr":     "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60",
					"dest_domain":      "1",
					"dest_recipient":   "AAEC",
					"tokens_forwarded": "0",
					"tokens_failed":    "1",
				},
			},
		}

		err := handleForward(ctx, events, msg, &idx)
		require.NoError(t, err)
		require.Len(t, ctx.Forwardings, 1)
		require.NotNil(t, ctx.Forwardings[0])
		require.Contains(t, string(ctx.Forwardings[0].Transfers), `"error":"some error"`)
		require.Contains(t, string(ctx.Forwardings[0].Transfers), `"denom":"utia"`)
		require.Contains(t, string(ctx.Forwardings[0].Transfers), `"amount":"1000"`)
	})
}
