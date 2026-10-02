// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/stretchr/testify/require"
)

// Real tx: 0A5524EE31AA746CD136BB7421F503B2D5D27F01ECFF65BE5C4DB9C23D2AD9B6 at height 10955535
// Event flow: message(action) → coin events → EventSendRemoteTransfer → hyperlane core events → EventTokenForwarded

func Test_handleForward(t *testing.T) {
	ts := time.Now()

	t.Run("nil event cursor", func(t *testing.T) {
		ctx := context.NewContext()
		msg := &storage.Message{
			Type:   types.MsgForward,
			Height: 100,
			Time:   ts,
		}
		err := handleForward(ctx, nil, msg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "nil event cursor")
	})

	t.Run("nil message", func(t *testing.T) {
		ctx := context.NewContext()
		c := NewCursor(nil)
		err := handleForward(ctx, c, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "nil message")
	})

	t.Run("unexpected action", func(t *testing.T) {
		ctx := context.NewContext()
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
		c := NewCursor(events)
		err := handleForward(ctx, c, msg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unexpected event action")
	})

	// Mirrors real tx 0A5524EE...: EventSendRemoteTransfer precedes EventTokenForwarded.
	t.Run("success: send remote transfer then token forwarded", func(t *testing.T) {
		ctx := context.NewContext()
		msg := &storage.Message{
			Id:     42,
			TxId:   7,
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

		c := NewCursor(events)
		err := handleForward(ctx, c, msg)
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
		require.EqualValues(t, 7, fwd.TxId)
		require.EqualValues(t, 42, fwd.MsgId)
	})

	// Intermediate non-forwarding events (coin_spent, coin_received, etc.) are skipped.
	t.Run("skips intermediate events without action key", func(t *testing.T) {
		ctx := context.NewContext()
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
					"sender":             "\"celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60\"",
					"token_id":           "\"0x726f757465725f61707000000000000000000000000000020000000000000024\"",
				},
			},
			{
				Height: 100,
				Type:   types.EventTypeHyperlanecorev1EventDispatch,
				Data:   testDispatchEvent(),
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

		c := NewCursor(events)
		err := handleForward(ctx, c, msg)
		require.NoError(t, err)
		require.Len(t, ctx.Forwardings, 1)
		require.Equal(t, uint64(1), ctx.Forwardings[0].DestDomain)
		require.Equal(t, "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", ctx.Forwardings[0].Address.Address)
	})

	t.Run("stops at next action event", func(t *testing.T) {
		ctx := context.NewContext()
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

		c := NewCursor(events)
		err := handleForward(ctx, c, msg)
		require.NoError(t, err)
		require.Len(t, ctx.Forwardings, 1)

		next, ok := c.Peek()
		require.True(t, ok, "cursor should stop at the next action event, not consume it")
		require.Equal(t, "/cosmos.bank.v1beta1.MsgSend", next.Data["action"])
	})

	t.Run("multiple messages in sequence", func(t *testing.T) {
		ctx := context.NewContext()
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
					"sender":             "\"celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60\"",
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
					"sender":             "\"celestia1ccqy2wlzf2zndn4vspmuksw5frqq0ufsgw4gmt\"",
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

		c := NewCursor(events)
		for i := range msgs {
			err := handleForward(ctx, c, msgs[i])
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

		c := NewCursor(events)
		err := handleForward(ctx, c, msg)
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

		c := NewCursor(events)
		err := handleForward(ctx, c, msg)
		require.NoError(t, err)
		require.Empty(t, ctx.Forwardings)
	})

	// Real tx 89CDA5CF at height 10727934 — pre-v8 format (before PR #6920), absent on mainnet and mocha-5.
	// EventTokenForwarded carries success/error but no token_id, so the parser fails instead of losing the forward.
	t.Run("pre-v8 format: missing token_id returns error", func(t *testing.T) {
		ctx := context.NewContext()
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
				Data: map[string]string{
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
				Data: map[string]string{
					"dest_domain":      "2147483647",
					"dest_recipient":   "\"0x00000000000000000000000050cB97b8613003DB9B278Bb89d3ab3C377F99727\"",
					"forward_addr":     "\"celestia184xycj78zfyhxd07rmg0wpc70r8scjm7jwj403\"",
					"tokens_failed":    "0",
					"tokens_forwarded": "1",
				},
			},
		}

		c := NewCursor(events)
		err := handleForward(ctx, c, msg)
		require.ErrorContains(t, err, "token_id is missing")
	})
}

// Real dispatch of the warp route 0x726f…0000 to domain 56 (see remote_transfer_test.go).
func testDispatchEvent() map[string]string {
	return map[string]string{
		"sender":            "0x726f757465725f61707000000000000000000000000000010000000000000000",
		"message":           "0x03000000e86d696c6b726f757465725f61707000000000000000000000000000010000000000000000000000380000000000000000000000007b4bf9feccff207ef2cb7101ceb15b8516021acd000000000000000000000000e0f9f661f106d6da1974fdc12a904e936834b3f8000000000000000000000000000000000000000000000000000000000c939ac0",
		"recipient":         "0x0000000000000000000000007b4bf9feccff207ef2cb7101ceb15b8516021acd",
		"destination":       "56",
		"origin_mailbox_id": "0x68797065726c616e650000000000000000000000000000000000000000000000",
	}
}

const testDispatchMessageId = "dcdb3f985ecd20c313c58c0f6b2a0d7ea980349134ee4813f6bd53cfe5bf0a1e"

func Test_processForward_HyperlaneTransfer(t *testing.T) {
	const (
		forwardAddr = "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60"
		tokenId     = "0x726f757465725f61707000000000000000000000000000010000000000000000" //nolint:gosec
	)
	ts := time.Now()

	sendEvent := storage.Event{
		Type: types.EventTypeHyperlanewarpv1EventSendRemoteTransfer,
		Data: map[string]string{
			"amount":             "1000000utia",
			"sender":             forwardAddr,
			"token_id":           tokenId,
			"recipient":          "0x000000000000000000000000e0f9f661f106d6da1974fdc12a904e936834b3f8",
			"destination_domain": "56",
		},
	}
	dispatchEvent := storage.Event{
		Type: types.EventTypeHyperlanecorev1EventDispatch,
		Data: testDispatchEvent(),
	}
	gasEvent := storage.Event{
		Type: types.EventTypeHyperlanecorepostDispatchv1EventGasPayment,
		Data: map[string]string{
			"igp_id":      "0x726f757465725f706f73745f6469737061746368000000040000000000000000",
			"payment":     "4379232utia",
			"gas_amount":  "64000",
			"message_id":  "0x" + testDispatchMessageId,
			"destination": "56",
		},
	}
	forwardedEvent := storage.Event{
		Type: types.EventTypeCelestiaforwardingv1EventTokenForwarded,
		Data: map[string]string{
			"forward_addr": forwardAddr,
			"token_id":     tokenId,
			"denom":        "utia",
			"amount":       "1000000",
			"message_id":   "0x" + testDispatchMessageId,
		},
	}
	actionEvent := storage.Event{
		Type: "message",
		Data: map[string]string{"action": "/celestia.forwarding.v1.MsgForward"},
	}
	nextEvent := storage.Event{
		Type: "message",
		Data: map[string]string{"action": "/cosmos.bank.v1beta1.MsgSend"},
	}
	newMsg := func() *storage.Message {
		return &storage.Message{Id: 42, TxId: 7, Type: types.MsgForward, Height: 1036866, Time: ts}
	}

	t.Run("forward creates outgoing hyperlane transfer", func(t *testing.T) {
		ctx := context.NewContext()
		c := NewCursor([]storage.Event{actionEvent, sendEvent, dispatchEvent, gasEvent, forwardedEvent, nextEvent})
		require.NoError(t, handleForward(ctx, c, newMsg()))

		require.Len(t, ctx.Forwardings, 1)
		require.Len(t, ctx.HlTransfers, 1)

		transfer := ctx.HlTransfers[0]
		require.Equal(t, types.HLTransferTypeSend, transfer.Type)
		require.EqualValues(t, 1036866, transfer.Height)
		require.EqualValues(t, 7, transfer.TxId)
		require.NotNil(t, transfer.Address)
		require.Equal(t, forwardAddr, transfer.Address.Address)
		require.Equal(t, "e0f9f661f106d6da1974fdc12a904e936834b3f8", transfer.CounterpartyAddress)
		require.EqualValues(t, 56, transfer.Counterparty)
		require.Equal(t, "1000000", transfer.Amount.String())
		require.Equal(t, "utia", transfer.Denom)
		require.Equal(t, testDispatchMessageId, hex.EncodeToString(transfer.MessageId))

		require.NotNil(t, transfer.GasPayment)
		require.Equal(t, "4379232", transfer.GasPayment.Amount.String())
		require.Equal(t, "64000", transfer.GasPayment.GasAmount.String())

		// the forwarding points at the same hyperlane message
		require.Equal(t, "0x"+hex.EncodeToString(transfer.MessageId), ctx.Forwardings[0].MessageId)

		// mailbox and token counters
		require.Equal(t, 1, ctx.HlMailboxes.Len())
		require.EqualValues(t, 1, ctx.HlMailboxes.Values()[0].SentMessages)
		require.Equal(t, 1, ctx.HlTokens.Len())
		require.EqualValues(t, 1, ctx.HlTokens.Values()[0].SentTransfers)
		require.Equal(t, "1000000", ctx.HlTokens.Values()[0].Sent.String())

		// the forwarding address is flagged even though the transfer registered it first
		address, ok := ctx.Addresses.Get(transfer.Address.String())
		require.True(t, ok)
		require.True(t, address.IsForwarding)

		next, ok := c.Peek()
		require.True(t, ok)
		require.Equal(t, "/cosmos.bank.v1beta1.MsgSend", next.Data["action"])
	})

	// Without EventTokenForwarded the token is unknown: nothing is saved instead of failing in storage.
	t.Run("no token forwarded event before next message", func(t *testing.T) {
		ctx := context.NewContext()
		c := NewCursor([]storage.Event{actionEvent, sendEvent, dispatchEvent, gasEvent, nextEvent})
		require.NoError(t, handleForward(ctx, c, newMsg()))

		require.Empty(t, ctx.Forwardings)
		require.Empty(t, ctx.HlTransfers)

		next, ok := c.Peek()
		require.True(t, ok)
		require.Equal(t, "/cosmos.bank.v1beta1.MsgSend", next.Data["action"])
	})
}
