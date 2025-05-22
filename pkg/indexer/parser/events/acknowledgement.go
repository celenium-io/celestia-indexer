// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package events

import (
	"strings"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
)

func handleAcknowledgement(_ *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	action := decoder.StringFromMap(events[*idx].Data, "action")
	isValidMsg := action == "/ibc.core.channel.v1.MsgAcknowledgement"
	if !isValidMsg {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processAcknowledgement(events, msg, idx)
}

func processAcknowledgement(events []storage.Event, msg *storage.Message, idx *int) error {
	if len(events)-1 < *idx || events[*idx].Type == storageTypes.EventTypeMessage {
		return nil
	}
	msg.IbcTransfer = &storage.IbcTransfer{
		Height: msg.Height,
		Time:   msg.Time,
	}

	ackEvent := events[*idx]
	if ackEvent.Type != storageTypes.EventTypeAcknowledgePacket {
		return errors.Errorf("invalid event type: %s", ackEvent.Type)
	}

	recvPacket, err := decode.NewRecvPacket(ackEvent.Data)
	if err != nil {
		return errors.Wrap(err, "parse acknowledge packet event")
	}

	msg.IbcChannel = &storage.IbcChannel{
		Id:             recvPacket.SrcChannel,
		TransfersCount: 1,
		Status:         storageTypes.IbcChannelStatusInitialization,
	}

	msg.IbcTransfer.ChannelId = recvPacket.SrcChannel
	msg.IbcTransfer.Port = recvPacket.SrcPort
	msg.IbcTransfer.ConnectionId = recvPacket.Connection
	msg.IbcTransfer.Sequence = recvPacket.Sequence
	if recvPacket.TimeoutHeight > 0 {
		msg.IbcTransfer.HeightTimeout = recvPacket.TimeoutHeight
	}
	if !recvPacket.Timeout.IsZero() {
		msg.IbcTransfer.Timeout = &recvPacket.Timeout
	}

	*idx += 2

	fundEvent := events[*idx]

	switch fundEvent.Type {
	case storageTypes.EventTypeWriteAcknowledgement:
		*idx += 2
		msg.IbcTransfer = nil
		return nil
	case storageTypes.EventTypeFungibleTokenPacket:
		ftp := decode.NewFungibleTokenPacket(fundEvent.Data)

		parts := strings.Split(ftp.Denom, "/")
		msg.IbcTransfer.Denom = parts[len(parts)-1]

		msg.IbcTransfer.Amount = ftp.Amount
		msg.IbcTransfer.Memo = ftp.Memo

		if strings.HasPrefix(ftp.Receiver, types.AddressPrefixCelestia) {
			msg.IbcTransfer.Receiver = &storage.Address{
				Address: ftp.Receiver,
			}
			msg.IbcChannel.Received = msg.IbcChannel.Received.Add(ftp.Amount)
		} else {
			msg.IbcTransfer.ReceiverAddress = &ftp.Receiver
		}
		if strings.HasPrefix(ftp.Sender, types.AddressPrefixCelestia) {
			msg.IbcTransfer.Sender = &storage.Address{
				Address: ftp.Sender,
			}
			msg.IbcChannel.Sent = msg.IbcChannel.Sent.Add(ftp.Amount)
		} else {
			msg.IbcTransfer.SenderAddress = &ftp.Sender
		}

		*idx += 2
		return nil
	default:
		return errors.Errorf("unexpected event for MsgAcknowledgement: %s", fundEvent.Type)
	}
}
