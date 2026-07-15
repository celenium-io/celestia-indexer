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

func handleRecvPacket(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	action := decoder.StringFromMap(event.Data, "action")
	isValidMsg := action == "/ibc.core.channel.v1.MsgRecvPacket"
	if !isValidMsg {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processRecvPacket(ctx, c, msg)
}

func processRecvPacket(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	event, ok := c.Peek()
	if !ok || event.Type == storageTypes.EventTypeMessage {
		ctx.RemoveLastIbcTransfer()
		return nil
	}

	transfer := ctx.GetLastIbcTransfer()
	var chanId string
	if event.Type == storageTypes.EventTypeRecvPacket && transfer != nil {
		rp, err := decode.NewRecvPacket(event.Data)
		if err != nil {
			return err
		}
		transfer.ConnectionId = rp.Connection
		chanId = rp.DstChannel
	}

	c.Skip(2)

	event, ok = c.Peek()
	if !ok {
		return nil
	}

	if event.Type == storageTypes.EventTypeIbccallbackerrorIcs27Packet {
		c.Next()
		if _, ok := c.Peek(); !ok {
			ctx.RemoveLastIbcTransfer()
			ctx.DeleteIbcChannel(chanId)
			return nil
		}
		event, _ = c.Peek()
	}

	if event.Type == storageTypes.EventTypeWriteAcknowledgement {
		c.Skip(2)
		ctx.RemoveLastIbcTransfer()
		ctx.DeleteIbcChannel(chanId)
		return nil
	}

	packet, err := decoder.Map(msg.Data, "Packet")
	if err != nil {
		return err
	}

	port := (storageTypes.PackedBytes)(packet).GetStringOrDefault("DestinationPort")

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

		for i := range msgs {
			decodedMsg, err := decode.Message(ctx, msgs[i], i, storageTypes.StatusSuccess, 0)
			if err != nil {
				return errors.Wrap(err, "decode message in RecvPacket")
			}

			if err := handle(ctx, c, &decodedMsg.Msg, ibcEventHandlers, "module"); err != nil {
				return errors.Wrap(err, "handle IBC msg event")
			}
		}

	case "transfer":
		current, _ := c.Peek()
		action := decoder.StringFromMap(current.Data, "action")

		if transfer == nil {
			return nil
		}
		if err := ctx.AddAddress(transfer.Sender); err != nil {
			return err
		}
		if err := ctx.AddAddress(transfer.Receiver); err != nil {
			return err
		}

		hasFtp := false
		for action == "" {
			if len(c.Remaining()) <= 1 {
				break
			}
			c.Next()
			next, _ := c.Peek()
			action = decoder.StringFromMap(next.Data, "action")

			if next.Type == storageTypes.EventTypeFungibleTokenPacket {
				hasFtp = true
				ftp := decode.NewFungibleTokenPacket(next.Data)
				if ftp.Error != "" {
					ctx.RemoveLastIbcTransfer()
					ctx.DeleteIbcChannel(chanId)
				}
			}
		}

		if !hasFtp {
			ctx.RemoveLastIbcTransfer()
			ctx.DeleteIbcChannel(chanId)
		}
	}

	return nil
}
