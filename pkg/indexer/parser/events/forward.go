// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"encoding/json"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/pkg/errors"
)

func handleForward(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/celestia.forwarding.v1.MsgForward" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processForward(ctx, events, msg, idx)
}

func processForward(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	var forwarding = storage.Forwarding{
		Height: msg.Height,
		Time:   msg.Time,
	}

	var tokens = make([]map[string]string, 0)

	for ; len(events) > *idx; *idx += 1 {
		switch events[*idx].Type {
		case types.EventTypeCelestiaforwardingv1EventTokenForwarded:
			forwarded, err := decode.NewEventTokenForwarded(events[*idx].Data)
			if err != nil {
				return errors.Wrap(err, "decoding token forwarded event")
			}

			token := map[string]string{
				"denom":  forwarded.Denom,
				"amount": forwarded.Amount,
			}
			if forwarded.Error != "" {
				token["error"] = forwarded.Error
			}
			tokens = append(tokens, token)

		case types.EventTypeCelestiaforwardingv1EventForwardingComplete:
			complete, err := decode.NewEventForwardingComplete(events[*idx].Data)
			if err != nil {
				return errors.Wrap(err, "decoding forwarding complete event")
			}
			forwarding.Address = &storage.Address{
				Address:      complete.ForwardAddress,
				Height:       msg.Height,
				LastHeight:   msg.Height,
				IsForwarding: true,
				Balance:      storage.EmptyBalance(),
			}
			if err := ctx.AddAddress(forwarding.Address); err != nil {
				return errors.Wrap(err, "adding forwarding address to context")
			}

			forwarding.SuccessCount = complete.SuccessfulCount
			forwarding.FailedCount = complete.FailedCount
			forwarding.DestDomain = complete.DestinationDomain
			forwarding.DestRecipient = complete.DestinationRecipient

			transfers, err := json.Marshal(tokens)
			if err != nil {
				return errors.Wrap(err, "json marshalling transfers")
			}
			forwarding.Transfers = transfers
		default:
			if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "" {
				msg.Forwarding = &forwarding
				return nil
			}
		}
	}

	msg.Forwarding = &forwarding
	return nil
}
