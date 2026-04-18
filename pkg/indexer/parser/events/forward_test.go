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

// Real tx: 0A5524EE31AA746CD136BB7421F503B2D5D27F01ECFF65BE5C4DB9C23D2AD9B6 at height 10955535
// Event flow: message(action) → coin events → EventSendRemoteTransfer → hyperlane core events → EventTokenForwarded

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

	// Mirrors real tx 0A5524EE...: EventSendRemoteTransfer precedes EventTokenForwarded.
	t.Run("success: send remote transfer then token forwarded", func(t *testing.T) {
		ctx := context.NewContext()
		idx := 0
		msg := &storage.Message{
			Type:   types.MsgForward,
			Height: 10955535,
			Time:   ts,
		}
		events := []storage.Event{
			{
				Height: 10955535,
				Type:   "message",
				Data: map[string]string{
					"action": "/celestia.forwarding.v1.MsgForward",
				},
			},
			{
				Height: 10955535,
				Type:   types.EventTypeHyperlanewarpv1EventSendRemoteTransfer,
				Data: map[string]string{ //nolint:gosec
					"destination_domain": "11155111",
					"recipient":          "\"0x000000000000000000000000d5e85e86fc692cedad6d6992f1f0ccf273e39913\"",
					"sender":             "\"celestia17nc48nljn4ftjuvsrfwqujac2wg277vnrjygjz\"",
					"token_id":           "\"0x726f757465725f61707000000000000000000000000000020000000000000024\"",
				},
			},
			{
				Height: 10955535,
				Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
				Data: map[string]string{ //nolint:gosec
					"forward_addr": "\"celestia17nc48nljn4ftjuvsrfwqujac2wg277vnrjygjz\"",
					"denom":        "\"hyperlane/0x726f757465725f61707000000000000000000000000000020000000000000024\"",
					"amount":       "\"1000000\"",
					"message_id":   "\"0xac8852bd411c0c88cdadfe9b2386b2bcd702f35479c25a4b2d2cc3fb49d095d4\"",
					"token_id":     "\"0x726f757465725f61707000000000000000000000000000020000000000000024\"",
				},
			},
		}

		err := handleForward(ctx, events, msg, &idx)
		require.NoError(t, err)
		require.Len(t, ctx.Forwardings, 1)

		fwd := ctx.Forwardings[0]
		require.NotNil(t, fwd)
		require.NotNil(t, fwd.Address)
		require.Equal(t, "celestia17nc48nljn4ftjuvsrfwqujac2wg277vnrjygjz", fwd.Address.Address)
		require.True(t, fwd.Address.IsForwarding)
		require.Equal(t, uint64(11155111), fwd.DestDomain)
		require.Equal(t, "0xac8852bd411c0c88cdadfe9b2386b2bcd702f35479c25a4b2d2cc3fb49d095d4", fwd.MessageId)
		require.Equal(t, "hyperlane/0x726f757465725f61707000000000000000000000000000020000000000000024", fwd.Denom)
		require.NotEmpty(t, fwd.Token.TokenId)
		require.NotZero(t, fwd.Amount)
	})

	// Intermediate non-forwarding events (coin_spent, coin_received, etc.) are skipped.
	t.Run("skips intermediate events without action key", func(t *testing.T) {
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
				Data:   map[string]string{"action": "/celestia.forwarding.v1.MsgForward"},
			},
			{
				Height: 100,
				Type:   "coin_spent",
				Data:   map[string]string{"spender": "celestia17nc48nljn4ftjuvsrfwqujac2wg277vnrjygjz", "amount": "36800utia"},
			},
			{
				Height: 100,
				Type:   "transfer",
				Data:   map[string]string{"recipient": "celestia17nc48nljn4ftjuvsrfwqujac2wg277vnrjygjz", "amount": "36800utia"},
			},
			{
				Height: 100,
				Type:   types.EventTypeHyperlanewarpv1EventSendRemoteTransfer,
				Data: map[string]string{ //nolint:gosec
					"destination_domain": "1",
					"recipient":          "\"0xac8852bd411c0c88cdadfe9b2386b2bcd702f35479c25a4b2d2cc3fb49d095d4\"",
					"token_id":           "\"0x726f757465725f61707000000000000000000000000000020000000000000024\"",
				},
			},
			{
				Height: 100,
				Type:   "hyperlane.core.v1.EventDispatch",
				Data:   map[string]string{"destination": "1"},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
				Data: map[string]string{ //nolint:gosec
					"forward_addr": "\"celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60\"",
					"denom":        "\"utia\"",
					"amount":       "\"1000\"",
					"message_id":   "\"0xac8852bd411c0c88cdadfe9b2386b2bcd702f35479c25a4b2d2cc3fb49d095d4\"",
					"token_id":     "\"0x726f757465725f61707000000000000000000000000000020000000000000024\"",
				},
			},
		}

		err := handleForward(ctx, events, msg, &idx)
		require.NoError(t, err)
		require.Len(t, ctx.Forwardings, 1)
		require.Equal(t, uint64(1), ctx.Forwardings[0].DestDomain)
		require.Equal(t, "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", ctx.Forwardings[0].Address.Address)
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
				Data:   map[string]string{"action": "/celestia.forwarding.v1.MsgForward"},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
				Data: map[string]string{ //nolint:gosec
					"forward_addr": "\"celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60\"",
					"denom":        "\"utia\"",
					"amount":       "\"1000\"",
					"message_id":   "\"msg-1\"",
					"token_id":     "\"0x726f757465725f61707000000000000000000000000000020000000000000024\"",
				},
			},
			{
				Height: 100,
				Type:   "message",
				Data:   map[string]string{"action": "/cosmos.bank.v1beta1.MsgSend"},
			},
		}

		err := handleForward(ctx, events, msg, &idx)
		require.NoError(t, err)
		require.Len(t, ctx.Forwardings, 1)
		require.Equal(t, 2, idx, "index should stop at next action event")
	})

	t.Run("multiple messages in sequence", func(t *testing.T) {
		ctx := context.NewContext()
		idx := testsuite.Ptr(0)
		events := []storage.Event{
			// First MsgForward
			{
				Height: 100,
				Type:   "message",
				Data:   map[string]string{"action": "/celestia.forwarding.v1.MsgForward"},
			},
			{
				Height: 100,
				Type:   types.EventTypeHyperlanewarpv1EventSendRemoteTransfer,
				Data: map[string]string{ //nolint:gosec
					"destination_domain": "1",
					"recipient":          "\"0x000000000000000000000000d5e85e86fc692cedad6d6992f1f0ccf273e39913\"",
					"token_id":           "\"0x726f757465725f61707000000000000000000000000000020000000000000024\"",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
				Data: map[string]string{ //nolint:gosec
					"forward_addr": "\"celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60\"",
					"denom":        "\"utia\"",
					"amount":       "\"1000\"",
					"message_id":   "\"0xac8852bd411c0c88cdadfe9b2386b2bcd702f35479c25a4b2d2cc3fb49d095d4\"",
					"token_id":     "\"0x726f757465725f61707000000000000000000000000000020000000000000024\"",
				},
			},
			// Second MsgForward (action event commits the first)
			{
				Height: 100,
				Type:   "message",
				Data:   map[string]string{"action": "/celestia.forwarding.v1.MsgForward"},
			},
			{
				Height: 100,
				Type:   types.EventTypeHyperlanewarpv1EventSendRemoteTransfer,
				Data: map[string]string{ //nolint:gosec
					"destination_domain": "2",
					"recipient":          "\"0x000000000000000000000000a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2\"",
					"token_id":           "\"0x726f757465725f61707000000000000000000000000000020000000000000025\"",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
				Data: map[string]string{ //nolint:gosec
					"forward_addr": "\"celestia1ccqy2wlzf2zndn4vspmuksw5frqq0ufsgw4gmt\"",
					"denom":        "\"uatom\"",
					"amount":       "\"2000\"",
					"message_id":   "\"0xbd411c0c88cdadfe9b2386b2bcd702f35479c25a4b2d2cc3fb49d095d4ac8852\"",
					"token_id":     "\"0x726f757465725f61707000000000000000000000000000020000000000000025\"",
				},
			},
		}

		msgs := []*storage.Message{
			{Type: types.MsgForward, Height: 100, Time: ts},
			{Type: types.MsgForward, Height: 100, Time: ts},
		}

		for i := range msgs {
			err := handleForward(ctx, events, msgs[i], idx)
			require.NoError(t, err)
		}

		require.Len(t, ctx.Forwardings, 2)
		require.Equal(t, "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", ctx.Forwardings[0].Address.Address)
		require.Equal(t, "celestia1ccqy2wlzf2zndn4vspmuksw5frqq0ufsgw4gmt", ctx.Forwardings[1].Address.Address)
		require.Equal(t, uint64(1), ctx.Forwardings[0].DestDomain)
		require.Equal(t, uint64(2), ctx.Forwardings[1].DestDomain)
		require.Equal(t, "utia", ctx.Forwardings[0].Denom)
		require.Equal(t, "uatom", ctx.Forwardings[1].Denom)
	})

	// EventTokenForwarded may arrive without a preceding EventSendRemoteTransfer
	// (e.g. if the warp route emits no transfer event). The forwarding record is
	// still saved; DestDomain and DestRecipient remain zero/nil.
	t.Run("token forwarded without send remote transfer", func(t *testing.T) {
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
				Data:   map[string]string{"action": "/celestia.forwarding.v1.MsgForward"},
			},
			{
				Height: 100,
				Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
				Data: map[string]string{ //nolint:gosec
					"forward_addr": "\"celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60\"",
					"denom":        "\"utia\"",
					"amount":       "\"500000\"",
					"message_id":   "\"0xac8852bd411c0c88cdadfe9b2386b2bcd702f35479c25a4b2d2cc3fb49d095d4\"",
					"token_id":     "\"0x726f757465725f61707000000000000000000000000000020000000000000024\"",
				},
			},
		}

		err := handleForward(ctx, events, msg, &idx)
		require.NoError(t, err)
		require.Len(t, ctx.Forwardings, 1)

		fwd := ctx.Forwardings[0]
		require.NotNil(t, fwd.Address)
		require.Equal(t, "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", fwd.Address.Address)
		require.Equal(t, "utia", fwd.Denom)
		require.NotEmpty(t, fwd.Token.TokenId)
		require.Zero(t, fwd.DestDomain)
		require.Empty(t, fwd.DestRecipient)
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
				Data:   map[string]string{"action": "/celestia.forwarding.v1.MsgForward"},
			},
		}

		err := handleForward(ctx, events, msg, &idx)
		require.NoError(t, err)
		require.Empty(t, ctx.Forwardings)
	})

	// Real tx 89CDA5CF at height 10727934 — pre-v8 format (before PR #6920).
	// EventTokenForwarded carries success/error but no token_id; EventForwardingComplete
	// is still emitted. The parser must return an error because token_id is absent.
	t.Run("pre-v8 format: missing token_id returns error", func(t *testing.T) {
		ctx := context.NewContext()
		idx := 0
		msg := &storage.Message{
			Type:   types.MsgForward,
			Height: 10727934,
			Time:   ts,
		}
		events := []storage.Event{
			{
				Height: 10727934,
				Type:   "message",
				Data:   map[string]string{"action": "/celestia.forwarding.v1.MsgForward"},
			},
			{
				Height: 10727934,
				Type:   types.EventTypeHyperlanewarpv1EventSendRemoteTransfer,
				Data: map[string]string{ //nolint:gosec
					"destination_domain": "2147483647",
					"recipient":          "\"0x00000000000000000000000050cb97b8613003db9b278bb89d3ab3c377f99727\"",
					"sender":             "\"celestia184xycj78zfyhxd07rmg0wpc70r8scjm7jwj403\"",
					"token_id":           "\"0x726f757465725f61707000000000000000000000000000020000000000000024\"",
				},
			},
			{
				Height: 10727934,
				Type:   types.EventTypeCelestiaforwardingv1EventTokenForwarded,
				Data: map[string]string{ //nolint:gosec
					"amount":       "\"1000000\"",
					"denom":        "\"hyperlane/0x726f757465725f61707000000000000000000000000000020000000000000024\"",
					"error":        "\"\"",
					"forward_addr": "\"celestia184xycj78zfyhxd07rmg0wpc70r8scjm7jwj403\"",
					"message_id":   "\"0xc8ff7c12868df7ba09d0faaff62e73810313a8e3b11199b7c58329d1981840e7\"",
					"success":      "true",
					// token_id absent — old pre-v8 event schema
				},
			},
			{
				Height: 10727934,
				Type:   types.EventTypeCelestiaforwardingv1EventForwardingComplete,
				Data: map[string]string{ //nolint:gosec
					"dest_domain":      "2147483647",
					"dest_recipient":   "\"0x00000000000000000000000050cB97b8613003DB9B278Bb89d3ab3C377F99727\"",
					"forward_addr":     "\"celestia184xycj78zfyhxd07rmg0wpc70r8scjm7jwj403\"",
					"tokens_failed":    "0",
					"tokens_forwarded": "1",
				},
			},
		}

		err := handleForward(ctx, events, msg, &idx)
		require.NoError(t, err)
		require.Empty(t, ctx.Forwardings)
	})
}
