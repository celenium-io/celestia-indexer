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

func handleRecvPacket(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	action := decoder.StringFromMap(events[*idx].Data, "action")
	isValidMsg := action == "/ibc.core.channel.v1.MsgRecvPacket"
	if !isValidMsg {
		return errors.Errorf("unexpected event action %s for message type %s (idx=%d, event=%s)", action, msg.Type.String(), *idx, events[*idx].Type)
	}
	*idx += 1
	return processRecvPacket(ctx, events, msg, idx)
}

func processRecvPacket(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if len(events)-1 < *idx || events[*idx].Type == storageTypes.EventTypeMessage {
		ctx.RemoveLastIbcTransfer()
		return nil
	}

	transfer := ctx.GetLastIbcTransfer()
	var chanId string
	if events[*idx].Type == storageTypes.EventTypeRecvPacket && transfer != nil {
		rp, err := decode.NewRecvPacket(events[*idx].Data)
		if err != nil {
			return err
		}
		transfer.ConnectionId = rp.Connection
		chanId = rp.DstChannel
	}

	*idx += 2

	if len(events)-1 < *idx {
		return nil
	}

	if events[*idx].Type == storageTypes.EventTypeIbccallbackerrorIcs27Packet {
		*idx += 1
	}

	if events[*idx].Type == storageTypes.EventTypeWriteAcknowledgement {
		*idx += 2
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

			if err := handle(ctx, events, &decodedMsg.Msg, idx, ibcEventHandlers, "module"); err != nil {
				return errors.Wrap(err, "handle IBC msg event")
			}
		}

	case "transfer":
		var action = decoder.StringFromMap(events[*idx].Data, "action")

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
		for action == "" && len(events)-1 > *idx {
			*idx += 1
			action = decoder.StringFromMap(events[*idx].Data, "action")

			if events[*idx].Type == storageTypes.EventTypeFungibleTokenPacket {
				hasFtp = true
				ftp := decode.NewFungibleTokenPacket(events[*idx].Data)
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
