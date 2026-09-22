// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/pkg/errors"
)

func handleAcknowledgement(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	action := decoder.StringFromMap(event.Data, "action")
	isValidMsg := action == "/ibc.core.channel.v1.MsgAcknowledgement"
	if !isValidMsg {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processAcknowledgement(ctx, c, msg)
}

func processAcknowledgement(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	event, ok := c.Peek()
	if !ok || event.Type == storageTypes.EventTypeMessage {
		ctx.RemoveIbcTransfer(msg.Id)
		return nil
	}
	packet, err := decoder.Map(msg.Data, "Packet")
	if err != nil {
		return err
	}

	port, err := (storageTypes.PackedBytes)(packet).GetString("SourcePort")
	if err != nil {
		return err
	}

	switch port {
	case "icahost":
		mapData, err := decoder.Map(packet, "Data")
		if err != nil {
			return errors.Wrap(err, "get data map")
		}

		msgs, err := decoder.MessagesFromMap(mapData, "Data")
		if err != nil {
			return errors.Wrap(err, "get messages from data map")
		}

		// nested handlers must not reach the next top-level message's events
		nested := c.Sub("action")
		for i := range msgs {
			decodedMsg, err := decode.NestedMessage(ctx, msgs[i], i, storageTypes.StatusSuccess, 0, msg.Id)
			if err != nil {
				return errors.Wrap(err, "decode message in Acknowledgement")
			}
			if err := handle(ctx, nested, &decodedMsg.Msg, ibcEventHandlers, "module"); err != nil {
				return errors.Wrap(err, "handle IBC msg event")
			}
		}

	case "transfer":
		transfer := ctx.IbcTransferByMsg(msg.Id)
		if transfer == nil {
			return nil
		}
		if err := ctx.AddAddress(transfer.Sender); err != nil {
			return err
		}
		if err := ctx.AddAddress(transfer.Receiver); err != nil {
			return err
		}

		var hasFtp, failed bool
		// MsgEvents also yields the tx's last event, where the error ack sits in relayer batches
		for event := range c.MsgEvents("action") {
			switch event.Type {
			case storageTypes.EventTypeAcknowledgePacket:
				ack, err := decode.NewAcknowledgementPacket(event.Data)
				if err != nil {
					return errors.Wrap(err, "ack packet")
				}
				transfer.ConnectionId = ack.PacketConnection
			case storageTypes.EventTypeFungibleTokenPacket:
				hasFtp = true
				if decode.NewFungibleTokenPacket(event.Data).Error != "" {
					failed = true
				}
			}
		}

		if !hasFtp || failed {
			ctx.RemoveIbcTransfer(msg.Id)
			return nil
		}
		ctx.AddIbcChannelTransfer(transfer)
	}

	return nil
}
