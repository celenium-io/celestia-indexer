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
		msg.IbcTransfer = nil
		msg.IbcChannel = nil
		return nil
	}
	if events[*idx].Type == storageTypes.EventTypeRecvPacket && msg.IbcTransfer != nil {
		rp, err := decode.NewRecvPacket(events[*idx].Data)
		if err != nil {
			return err
		}
		msg.IbcTransfer.ConnectionId = rp.Connection
	}

	*idx += 2

	if len(events)-1 < *idx {
		return nil
	}

	if events[*idx].Type == storageTypes.EventTypeWriteAcknowledgement {
		*idx += 2
		msg.IbcTransfer = nil
		msg.IbcChannel = nil
		return nil
	}

	packet, err := decoder.Map(msg.Data, "Packet")
	if err != nil {
		return err
	}

	port := decoder.StringFromMap(packet, "DestinationPort")

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
			decodedMsg, err := decode.Message(ctx, msgs[i], i, storageTypes.StatusSuccess)
			if err != nil {
				return errors.Wrap(err, "decode message in RecvPacket")
			}

			if err := handle(ctx, events, &decodedMsg.Msg, idx, ibcEventHandlers, "module"); err != nil {
				return errors.Wrap(err, "handle IBC msg event")
			}
			msg.Addresses = append(msg.Addresses, decodedMsg.Addresses...)
		}

	case "transfer":
		var action = decoder.StringFromMap(events[*idx].Data, "action")

		if msg.IbcTransfer == nil {
			return nil
		}

		if err := ctx.AddAddress(msg.IbcTransfer.Sender); err != nil {
			return err
		}
		if err := ctx.AddAddress(msg.IbcTransfer.Receiver); err != nil {
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
					msg.IbcTransfer = nil
					msg.IbcChannel = nil
				}
			}
		}
		if !hasFtp {
			msg.IbcTransfer = nil
			msg.IbcChannel = nil
		}
	}

	return nil
}
