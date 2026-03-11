// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"strings"

	json "github.com/bytedance/sonic"
	"github.com/celenium-io/celestia-indexer/internal/math"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	icaTypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	transferTypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	coreChannel "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// MsgChannelOpenInit defines an sdk.Msg to initialize a channel handshake. It
// is called by a relayer on Chain A.
func MsgChannelOpenInit(ctx *context.Context, msgId uint64, m *coreChannel.MsgChannelOpenInit) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgChannelOpenInit
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgChannelOpenTry defines a msg sent by a Relayer to try to open a channel
// on Chain B. The version field within the Channel field has been deprecated. Its
// value will be ignored by core IBC.
func MsgChannelOpenTry(ctx *context.Context, msgId uint64, m *coreChannel.MsgChannelOpenTry) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgChannelOpenTry
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgChannelOpenAck defines a msg sent by a Relayer to Chain A to acknowledge
// the change of channel state to TRYOPEN on Chain B.
func MsgChannelOpenAck(ctx *context.Context, msgId uint64, m *coreChannel.MsgChannelOpenAck) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgChannelOpenAck
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgChannelOpenConfirm defines a msg sent by a Relayer to Chain B to
// acknowledge the change of channel state to OPEN on Chain A.
func MsgChannelOpenConfirm(ctx *context.Context, msgId uint64, m *coreChannel.MsgChannelOpenConfirm) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgChannelOpenConfirm
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgChannelCloseInit defines a msg sent by a Relayer to Chain A
// to close a channel with Chain B.
func MsgChannelCloseInit(ctx *context.Context, msgId uint64, m *coreChannel.MsgChannelCloseInit) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgChannelCloseInit
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgChannelCloseConfirm defines a msg sent by a Relayer to Chain B
// to acknowledge the change of channel state to CLOSED on Chain A.
func MsgChannelCloseConfirm(ctx *context.Context, msgId uint64, m *coreChannel.MsgChannelCloseConfirm) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgChannelCloseConfirm
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgRecvPacket receives an incoming IBC packet
func MsgRecvPacket(ctx *context.Context, status storageTypes.Status, codec codec.Codec, data storageTypes.PackedBytes, txId, msgId uint64, m *coreChannel.MsgRecvPacket) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgRecvPacket
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	if err != nil || status == storageTypes.StatusFailed {
		return msgType, err
	}

	packetMap, ok := data["Packet"].(map[string]any)
	if !ok {
		return msgType, errors.Wrap(err, "Packet is not map")
	}

	switch m.Packet.DestinationPort {
	case "icahost":
		var packet icaTypes.InterchainAccountPacketData
		if err := json.Unmarshal(m.Packet.Data, &packet); err != nil {
			return msgType, errors.Wrap(err, "InterchainAccountPacketData")
		}

		packetMapData := map[string]any{
			"Type": packet.Type,
			"Memo": packet.Memo,
		}

		var tx icaTypes.CosmosTx
		if err := codec.Unmarshal(packet.Data, &tx); err != nil {
			if err := codec.UnmarshalJSON(packet.Data, &tx); err != nil {
				packetMapData["Data"] = []cosmosTypes.Msg{}
				packetMap["Data"] = packetMapData
				return msgType, nil
			}
		}

		msgs := make([]cosmosTypes.Msg, len(tx.Messages))
		for i, rawMsg := range tx.Messages {
			var msg cosmosTypes.Msg
			if err := codec.UnpackAny(rawMsg, &msg); err != nil {
				return msgType, errors.Wrap(err, "cosmosTypes.Msg")
			}
			if grant, ok := msg.(*authz.MsgGrant); ok {
				grant.Grant.Authorization = nil // TODO: make more beautiful
			}
			msgs[i] = msg
		}
		packetMapData["Data"] = msgs
		packetMap["Data"] = packetMapData
		return msgType, nil

	case "transfer":
		var packet transferTypes.FungibleTokenPacketData
		if err := json.Unmarshal(m.Packet.Data, &packet); err != nil {
			return msgType, errors.Wrap(err, "FungibleTokenPacketData")
		}
		packetMap["Data"] = packet

		transfer := &storage.IbcTransfer{
			Amount:    decimal.RequireFromString(packet.Amount),
			Memo:      packet.Memo,
			ChannelId: m.Packet.DestinationChannel,
			Port:      m.Packet.DestinationPort,
			Sequence:  m.Packet.Sequence,
			Denom:     packet.Denom,
			Height:    ctx.Block.Height,
			Time:      ctx.Block.Time,
			TxId:      txId,
		}

		partsDenom := strings.Split(packet.Denom, "/")
		if len(partsDenom) == 3 {
			transfer.Denom = partsDenom[2]
		}

		if m.Packet.TimeoutHeight.RevisionHeight > 0 {
			transfer.HeightTimeout = m.Packet.TimeoutHeight.RevisionHeight
		}
		if m.Packet.TimeoutTimestamp > 0 {
			ts := math.TimeFromNano(m.Packet.TimeoutTimestamp)
			transfer.Timeout = &ts
		}

		channel := &storage.IbcChannel{
			Id:             m.Packet.DestinationChannel,
			TransfersCount: 1,
			Status:         storageTypes.IbcChannelStatusInitialization,
		}
		prefix, hash, err := pkgTypes.Address(packet.Receiver).Decode()
		if err != nil {
			return msgType, nil
		}
		if prefix == pkgTypes.AddressPrefixCelestia {
			transfer.Receiver = &storage.Address{
				Address:    packet.Receiver,
				Balance:    storage.EmptyBalance(),
				Height:     ctx.Block.Height,
				LastHeight: ctx.Block.Height,
				Hash:       hash,
			}
			ctx.AddAddressMessage(&storage.MsgAddress{
				MsgId:   msgId,
				Type:    storageTypes.MsgAddressTypeReceiver,
				Address: transfer.Receiver,
			})
			channel.Received = channel.Received.Add(transfer.Amount)
		} else {
			transfer.ReceiverAddress = &packet.Receiver
		}
		prefix, hash, err = pkgTypes.Address(packet.Sender).Decode()
		if err != nil {
			return msgType, nil
		}
		if prefix == pkgTypes.AddressPrefixCelestia {
			transfer.Sender = &storage.Address{
				Address:    packet.Sender,
				Balance:    storage.EmptyBalance(),
				Height:     ctx.Block.Height,
				LastHeight: ctx.Block.Height,
				Hash:       hash,
			}
			ctx.AddAddressMessage(&storage.MsgAddress{
				MsgId:   msgId,
				Type:    storageTypes.MsgAddressTypeSender,
				Address: transfer.Sender,
			})
			channel.Sent = channel.Sent.Add(transfer.Amount)
		} else {
			transfer.SenderAddress = &packet.Sender
		}
		ctx.AddIbcChannel(channel)
		ctx.AddIbcTransfer(transfer)
		return msgType, nil
	default:
		return msgType, errors.Errorf("unknown destination port: %s", m.Packet.DestinationPort)
	}
}

// MsgTimeout receives a timed-out packet
func MsgTimeout(ctx *context.Context, msgId uint64, m *coreChannel.MsgTimeout) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgTimeout
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgTimeoutOnClose timed-out packet upon counterparty channel closure
func MsgTimeoutOnClose(ctx *context.Context, msgId uint64, m *coreChannel.MsgTimeoutOnClose) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgTimeoutOnClose
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgAcknowledgement receives incoming IBC acknowledgement
func MsgAcknowledgement(ctx *context.Context, status storageTypes.Status, codec codec.Codec, data storageTypes.PackedBytes, txId, msgId uint64, m *coreChannel.MsgAcknowledgement) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgAcknowledgement
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	if err != nil || status == storageTypes.StatusFailed {
		return msgType, err
	}

	packetMap, ok := data["Packet"].(map[string]any)
	if !ok {
		return msgType, errors.Wrap(err, "Packet is not map")
	}

	switch m.Packet.SourcePort {
	case "icahost":
		var packet icaTypes.InterchainAccountPacketData
		if err := json.Unmarshal(m.Packet.Data, &packet); err != nil {
			return msgType, errors.Wrap(err, "InterchainAccountPacketData")
		}

		packetMapData := map[string]any{
			"Type": packet.Type,
			"Memo": packet.Memo,
		}

		var tx icaTypes.CosmosTx
		if err := codec.Unmarshal(packet.Data, &tx); err != nil {
			if err := codec.UnmarshalJSON(packet.Data, &tx); err != nil {
				packetMapData["Data"] = []cosmosTypes.Msg{}
				packetMap["Data"] = packetMapData
				return msgType, nil
			}
		}

		msgs := make([]cosmosTypes.Msg, len(tx.Messages))
		for i, rawMsg := range tx.Messages {
			var msg cosmosTypes.Msg
			if err := codec.UnpackAny(rawMsg, &msg); err != nil {
				return msgType, errors.Wrap(err, "cosmosTypes.Msg")
			}
			msgs[i] = msg
		}
		packetMapData["Data"] = msgs
		packetMap["Data"] = packetMapData
		return msgType, nil

	case "transfer":
		var packet transferTypes.FungibleTokenPacketData
		if err := json.Unmarshal(m.Packet.Data, &packet); err != nil {
			return msgType, errors.Wrap(err, "FungibleTokenPacketData")
		}
		packetMap["Data"] = packet

		transfer := &storage.IbcTransfer{
			Amount:    decimal.RequireFromString(packet.Amount),
			Memo:      packet.Memo,
			ChannelId: m.Packet.SourceChannel,
			Port:      m.Packet.SourcePort,
			Sequence:  m.Packet.Sequence,
			Denom:     packet.Denom,
			Height:    ctx.Block.Height,
			Time:      ctx.Block.Time,
			TxId:      txId,
		}

		partsDenom := strings.Split(packet.Denom, "/")
		if len(partsDenom) == 3 {
			transfer.Denom = partsDenom[2]
		}

		if m.Packet.TimeoutHeight.RevisionHeight > 0 {
			transfer.HeightTimeout = m.Packet.TimeoutHeight.RevisionHeight
		}
		if m.Packet.TimeoutTimestamp > 0 {
			ts := math.TimeFromNano(m.Packet.TimeoutTimestamp)
			transfer.Timeout = &ts
		}

		channel := &storage.IbcChannel{
			Id:             m.Packet.SourceChannel,
			TransfersCount: 1,
			Status:         storageTypes.IbcChannelStatusInitialization,
		}
		prefix, hash, err := pkgTypes.Address(packet.Receiver).Decode()
		if err != nil {
			return msgType, nil
		}
		if prefix == pkgTypes.AddressPrefixCelestia {
			transfer.Receiver = &storage.Address{
				Address:    packet.Receiver,
				Balance:    storage.EmptyBalance(),
				Height:     ctx.Block.Height,
				LastHeight: ctx.Block.Height,
				Hash:       hash,
			}
			ctx.AddAddressMessage(&storage.MsgAddress{
				MsgId:   msgId,
				Type:    storageTypes.MsgAddressTypeReceiver,
				Address: transfer.Receiver,
			})
			channel.Received = channel.Received.Add(transfer.Amount)
		} else {
			transfer.ReceiverAddress = &packet.Receiver
		}
		prefix, hash, err = pkgTypes.Address(packet.Sender).Decode()
		if err != nil {
			return msgType, nil
		}
		if prefix == pkgTypes.AddressPrefixCelestia {
			transfer.Sender = &storage.Address{
				Address:    packet.Sender,
				Balance:    storage.EmptyBalance(),
				Height:     ctx.Block.Height,
				LastHeight: ctx.Block.Height,
				Hash:       hash,
			}
			ctx.AddAddressMessage(&storage.MsgAddress{
				MsgId:   msgId,
				Type:    storageTypes.MsgAddressTypeSender,
				Address: transfer.Sender,
			})
			channel.Sent = channel.Sent.Add(transfer.Amount)
		} else {
			transfer.SenderAddress = &packet.Sender
		}

		ctx.AddIbcChannel(channel)
		ctx.AddIbcTransfer(transfer)
		return msgType, nil
	default:
		return msgType, errors.Errorf("unknown source port: %s", m.Packet.SourcePort)
	}
}

func MsgUpdateParamsChannel(ctx *context.Context, msgId uint64, m *coreChannel.MsgUpdateParams) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateParams
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height, msgId)
	return msgType, err
}
