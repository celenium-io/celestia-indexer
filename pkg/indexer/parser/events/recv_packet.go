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

func handleRecvPacket(_ *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	action := decoder.StringFromMap(events[*idx].Data, "action")
	isValidMsg := action == "/ibc.core.channel.v1.MsgRecvPacket"
	if !isValidMsg {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processRecvPacket(events, msg, idx)
}

func processRecvPacket(events []storage.Event, msg *storage.Message, idx *int) error {
	if len(events)-1 < *idx || events[*idx].Type == storageTypes.EventTypeMessage {
		return nil
	}
	msg.IbcTransfer = &storage.IbcTransfer{
		Height: msg.Height,
		Time:   msg.Time,
	}

	recvPacketEvent := events[*idx]
	if recvPacketEvent.Type != storageTypes.EventTypeRecvPacket {
		return errors.Errorf("invalid event type: %s", recvPacketEvent.Type)
	}

	recvPacket, err := decode.NewRecvPacket(recvPacketEvent.Data)
	if err != nil {
		return errors.Wrap(err, "parse recv packet event")
	}

	msg.IbcChannel = &storage.IbcChannel{
		Id:             recvPacket.DstChannel,
		TransfersCount: 1,
		Status:         storageTypes.IbcChannelStatusInitialization,
	}

	msg.IbcTransfer.ChannelId = recvPacket.DstChannel
	msg.IbcTransfer.Port = recvPacket.DstPort
	msg.IbcTransfer.ConnectionId = recvPacket.Connection
	msg.IbcTransfer.Sequence = recvPacket.Sequence
	if recvPacket.TimeoutHeight > 0 {
		msg.IbcTransfer.HeightTimeout = recvPacket.TimeoutHeight
	}
	if !recvPacket.Timeout.IsZero() {
		msg.IbcTransfer.Timeout = &recvPacket.Timeout
	}

	*idx += 3
	coinReceivedEvent := events[*idx]
	if coinReceivedEvent.Type == storageTypes.EventTypeMessage {
		*idx += 1
		msg.IbcTransfer = nil
		msg.IbcChannel = nil
		return nil
	}
	if coinReceivedEvent.Type != storageTypes.EventTypeCoinReceived {
		return errors.Errorf("invalid event type: %s", coinReceivedEvent.Type)
	}

	received, err := decode.NewCoinReceived(coinReceivedEvent.Data)
	if err != nil {
		return errors.Wrap(err, "parse coinr received in recv packet")
	}

	*idx += 3
	fundEvent := events[*idx]
	if fundEvent.Type != storageTypes.EventTypeFungibleTokenPacket {
		return nil
	}

	ftp := decode.NewFungibleTokenPacket(fundEvent.Data)
	if received.Amount.Amount.String() == ftp.Amount.String() {
		msg.IbcTransfer.Denom = received.Amount.GetDenom()
	} else {
		parts := strings.Split(ftp.Denom, "/")
		msg.IbcTransfer.Denom = parts[len(parts)-1]
	}

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

	action := decoder.StringFromMap(events[*idx].Data, "action")
	for action == "" && len(events)-1 > *idx {
		*idx++
		action = decoder.StringFromMap(events[*idx].Data, "action")
	}

	if action == "" {
		*idx += 1
	}
	return nil
}
