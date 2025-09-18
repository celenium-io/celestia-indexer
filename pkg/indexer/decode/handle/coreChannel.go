// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"strings"

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
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// MsgChannelOpenInit defines an sdk.Msg to initialize a channel handshake. It
// is called by a relayer on Chain A.
func MsgChannelOpenInit(ctx *context.Context, m *coreChannel.MsgChannelOpenInit) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgChannelOpenInit
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgChannelOpenTry defines a msg sent by a Relayer to try to open a channel
// on Chain B. The version field within the Channel field has been deprecated. Its
// value will be ignored by core IBC.
func MsgChannelOpenTry(ctx *context.Context, m *coreChannel.MsgChannelOpenTry) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgChannelOpenTry
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgChannelOpenAck defines a msg sent by a Relayer to Chain A to acknowledge
// the change of channel state to TRYOPEN on Chain B.
func MsgChannelOpenAck(ctx *context.Context, m *coreChannel.MsgChannelOpenAck) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgChannelOpenAck
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgChannelOpenConfirm defines a msg sent by a Relayer to Chain B to
// acknowledge the change of channel state to OPEN on Chain A.
func MsgChannelOpenConfirm(ctx *context.Context, m *coreChannel.MsgChannelOpenConfirm) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgChannelOpenConfirm
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgChannelCloseInit defines a msg sent by a Relayer to Chain A
// to close a channel with Chain B.
func MsgChannelCloseInit(ctx *context.Context, m *coreChannel.MsgChannelCloseInit) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgChannelCloseInit
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgChannelCloseConfirm defines a msg sent by a Relayer to Chain B
// to acknowledge the change of channel state to CLOSED on Chain A.
func MsgChannelCloseConfirm(ctx *context.Context, m *coreChannel.MsgChannelCloseConfirm) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgChannelCloseConfirm
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgRecvPacket receives an incoming IBC packet
func MsgRecvPacket(ctx *context.Context, codec codec.Codec, data storageTypes.PackedBytes, m *coreChannel.MsgRecvPacket) (storageTypes.MsgType, []storage.AddressWithType, *storage.IbcTransfer, *storage.IbcChannel, error) {
	msgType := storageTypes.MsgRecvPacket
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	if err != nil {
		return msgType, addresses, nil, nil, err
	}

	packetMap, ok := data["Packet"].(map[string]any)
	if !ok {
		return msgType, addresses, nil, nil, errors.Wrap(err, "Packet is not map")
	}

	switch m.Packet.DestinationPort {
	case "icahost":
		var packet icaTypes.InterchainAccountPacketData
		if err := json.Unmarshal(m.Packet.Data, &packet); err != nil {
			return msgType, addresses, nil, nil, errors.Wrap(err, "InterchainAccountPacketData")
		}

		packetMapData := map[string]any{
			"Type": packet.Type,
			"Memo": packet.Memo,
		}

		var tx icaTypes.CosmosTx
		if err := codec.Unmarshal(packet.Data, &tx); err != nil {
			if err := codec.UnmarshalJSON(packet.Data, &tx); err != nil {
				return msgType, addresses, nil, nil, errors.Wrap(err, "icaTypes.CosmosTx")
			}
		}

		msgs := make([]cosmosTypes.Msg, len(tx.Messages))
		for i, rawMsg := range tx.Messages {
			var msg cosmosTypes.Msg
			if err := codec.UnpackAny(rawMsg, &msg); err != nil {
				return msgType, addresses, nil, nil, errors.Wrap(err, "cosmosTypes.Msg")
			}
			if grant, ok := msg.(*authz.MsgGrant); ok {
				grant.Grant.Authorization = nil // TODO: make more beautiful
			}
			msgs[i] = msg
		}
		packetMapData["Data"] = msgs
		packetMap["Data"] = packetMapData
		return msgType, addresses, nil, nil, nil

	case "transfer":
		var packet transferTypes.FungibleTokenPacketData
		if err := json.Unmarshal(m.Packet.Data, &packet); err != nil {
			return msgType, addresses, nil, nil, errors.Wrap(err, "FungibleTokenPacketData")
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
		prefix, _, err := pkgTypes.Address(packet.Receiver).Decode()
		if err != nil {
			return msgType, addresses, nil, nil, nil
		}
		if prefix == pkgTypes.AddressPrefixCelestia {
			transfer.Receiver = &storage.Address{
				Address: packet.Receiver,
				Balance: storage.EmptyBalance(),
			}
			addresses = append(addresses, storage.AddressWithType{
				Address: *transfer.Receiver,
				Type:    storageTypes.MsgAddressTypeReceiver,
			})
			channel.Received = channel.Received.Add(transfer.Amount)
		} else {
			transfer.ReceiverAddress = &packet.Receiver
		}
		prefix, _, err = pkgTypes.Address(packet.Sender).Decode()
		if err != nil {
			return msgType, addresses, nil, nil, nil
		}
		if prefix == pkgTypes.AddressPrefixCelestia {
			transfer.Sender = &storage.Address{
				Address: packet.Sender,
				Balance: storage.EmptyBalance(),
			}
			addresses = append(addresses, storage.AddressWithType{
				Address: *transfer.Sender,
				Type:    storageTypes.MsgAddressTypeSender,
			})
			channel.Sent = channel.Sent.Add(transfer.Amount)
		} else {
			transfer.SenderAddress = &packet.Sender
		}

		return msgType, addresses, transfer, channel, nil
	default:
		return msgType, addresses, nil, nil, errors.Errorf("unknown destination port: %s", m.Packet.DestinationPort)
	}
}

// MsgTimeout receives a timed-out packet
func MsgTimeout(ctx *context.Context, m *coreChannel.MsgTimeout) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgTimeout
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgTimeoutOnClose timed-out packet upon counterparty channel closure
func MsgTimeoutOnClose(ctx *context.Context, m *coreChannel.MsgTimeoutOnClose) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgTimeoutOnClose
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgAcknowledgement receives incoming IBC acknowledgement
func MsgAcknowledgement(ctx *context.Context, codec codec.Codec, data storageTypes.PackedBytes, m *coreChannel.MsgAcknowledgement) (storageTypes.MsgType, []storage.AddressWithType, *storage.IbcTransfer, *storage.IbcChannel, error) {
	msgType := storageTypes.MsgAcknowledgement
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)

	packetMap, ok := data["Packet"].(map[string]any)
	if !ok {
		return msgType, addresses, nil, nil, errors.Wrap(err, "Packet is not map")
	}

	switch m.Packet.SourcePort {
	case "icahost":
		var packet icaTypes.InterchainAccountPacketData
		if err := json.Unmarshal(m.Packet.Data, &packet); err != nil {
			return msgType, addresses, nil, nil, errors.Wrap(err, "InterchainAccountPacketData")
		}

		packetMapData := map[string]any{
			"Type": packet.Type,
			"Memo": packet.Memo,
		}

		var tx icaTypes.CosmosTx
		if err := codec.Unmarshal(packet.Data, &tx); err != nil {
			if err := codec.UnmarshalJSON(packet.Data, &tx); err != nil {
				return msgType, addresses, nil, nil, errors.Wrap(err, "icaTypes.CosmosTx")
			}
		}

		msgs := make([]cosmosTypes.Msg, len(tx.Messages))
		for i, rawMsg := range tx.Messages {
			var msg cosmosTypes.Msg
			if err := codec.UnpackAny(rawMsg, &msg); err != nil {
				return msgType, addresses, nil, nil, errors.Wrap(err, "cosmosTypes.Msg")
			}
			msgs[i] = msg
		}
		packetMapData["Data"] = msgs
		packetMap["Data"] = packetMapData
		return msgType, addresses, nil, nil, nil

	case "transfer":
		var packet transferTypes.FungibleTokenPacketData
		if err := json.Unmarshal(m.Packet.Data, &packet); err != nil {
			return msgType, addresses, nil, nil, errors.Wrap(err, "FungibleTokenPacketData")
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
		prefix, _, err := pkgTypes.Address(packet.Receiver).Decode()
		if err != nil {
			return msgType, addresses, nil, nil, nil
		}
		if prefix == pkgTypes.AddressPrefixCelestia {
			transfer.Receiver = &storage.Address{
				Address: packet.Receiver,
				Balance: storage.EmptyBalance(),
			}
			addresses = append(addresses, storage.AddressWithType{
				Address: *transfer.Receiver,
				Type:    storageTypes.MsgAddressTypeReceiver,
			})
			channel.Received = channel.Received.Add(transfer.Amount)
		} else {
			transfer.ReceiverAddress = &packet.Receiver
		}
		prefix, _, err = pkgTypes.Address(packet.Sender).Decode()
		if err != nil {
			return msgType, addresses, nil, nil, nil
		}
		if prefix == pkgTypes.AddressPrefixCelestia {
			transfer.Sender = &storage.Address{
				Address: packet.Sender,
				Balance: storage.EmptyBalance(),
			}
			addresses = append(addresses, storage.AddressWithType{
				Address: *transfer.Sender,
				Type:    storageTypes.MsgAddressTypeSender,
			})
			channel.Sent = channel.Sent.Add(transfer.Amount)
		} else {
			transfer.SenderAddress = &packet.Sender
		}

		return msgType, addresses, transfer, channel, nil
	default:
		return msgType, addresses, nil, nil, errors.Errorf("unknown source port: %s", m.Packet.SourcePort)
	}
}

func MsgUpdateParamsChannel(ctx *context.Context, m *coreChannel.MsgUpdateParams) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUpdateParams
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
