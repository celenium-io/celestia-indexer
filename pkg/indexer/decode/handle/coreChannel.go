// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	stdMath "math"
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
	"github.com/rs/zerolog/log"
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
func MsgRecvPacket(
	ctx *context.Context,
	status storageTypes.Status,
	codec codec.Codec,
	data storageTypes.PackedBytes,
	txId, msgId uint64,
	m *coreChannel.MsgRecvPacket,

) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgRecvPacket
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	if err != nil || status == storageTypes.StatusFailed {
		return msgType, err
	}

	packetMap, ok := data["Packet"].(map[string]any)
	if !ok {
		return msgType, errors.Errorf("Packet is not map: %T", data["Packet"])
	}

	err = handlePacketData(ctx, codec, packetMap, m.Packet, m.Packet.DestinationPort, m.Packet.DestinationChannel, txId, msgId)
	return msgType, err
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
		return msgType, errors.Errorf("Packet is not map: %T", data["Packet"])
	}

	err = handlePacketData(ctx, codec, packetMap, m.Packet, m.Packet.SourcePort, m.Packet.SourceChannel, txId, msgId)
	return msgType, err
}

func MsgUpdateParamsChannel(ctx *context.Context, msgId uint64, m *coreChannel.MsgUpdateParams) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateParams
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// handlePacketData decodes relayer-supplied packet payload. Malformed data or unknown ports are
// accepted on-chain with an error ack, so they are logged and skipped instead of halting the indexer.
func handlePacketData(
	ctx *context.Context,
	codec codec.Codec,
	packetMap map[string]any,
	packet coreChannel.Packet,
	port, channel string,
	txId, msgId uint64,
) error {
	switch port {
	case "icahost":
		handleIcaPacketData(ctx, codec, packetMap, packet, port, channel, msgId)
		return nil
	case "transfer":
		return handleTransferPacketData(ctx, packetMap, packet, port, channel, txId, msgId)
	default:
		logSkippedPacket(ctx, packet, port, channel, msgId, errors.New("unknown port"))
		return nil
	}
}

func handleIcaPacketData(
	ctx *context.Context,
	codec codec.Codec,
	packetMap map[string]any,
	packet coreChannel.Packet,
	port, channel string,
	msgId uint64,
) {
	var data icaTypes.InterchainAccountPacketData
	if err := json.Unmarshal(packet.Data, &data); err != nil {
		logSkippedPacket(ctx, packet, port, channel, msgId, errors.Wrap(err, "InterchainAccountPacketData"))
		return
	}

	packetMapData := map[string]any{
		"Type": data.Type,
		"Memo": data.Memo,
		"Data": []cosmosTypes.Msg{},
	}
	packetMap["Data"] = packetMapData

	var tx icaTypes.CosmosTx
	if err := codec.Unmarshal(data.Data, &tx); err != nil {
		if err := codec.UnmarshalJSON(data.Data, &tx); err != nil {
			return
		}
	}

	msgs := make([]cosmosTypes.Msg, len(tx.Messages))
	for i, rawMsg := range tx.Messages {
		var msg cosmosTypes.Msg
		if err := codec.UnpackAny(rawMsg, &msg); err != nil {
			logSkippedPacket(ctx, packet, port, channel, msgId, errors.Wrap(err, "cosmosTypes.Msg"))
			return
		}
		if grant, ok := msg.(*authz.MsgGrant); ok {
			grant.Grant.Authorization = nil // TODO: make more beautiful
		}
		msgs[i] = msg
	}
	packetMapData["Data"] = msgs
}

func handleTransferPacketData(
	ctx *context.Context,
	packetMap map[string]any,
	packet coreChannel.Packet,
	port, channel string,
	txId, msgId uint64,
) error {
	var data transferTypes.FungibleTokenPacketData
	if err := json.Unmarshal(packet.Data, &data); err != nil {
		logSkippedPacket(ctx, packet, port, channel, msgId, errors.Wrap(err, "FungibleTokenPacketData"))
		return nil
	}
	packetMap["Data"] = data

	amount, err := storageTypes.NumericFromString(data.Amount)
	if err != nil {
		logSkippedPacket(ctx, packet, port, channel, msgId, errors.Wrap(err, "parse transfer amount"))
		return nil
	}
	transfer := &storage.IbcTransfer{
		Amount:    amount,
		Memo:      data.Memo,
		ChannelId: channel,
		Port:      port,
		Sequence:  packet.Sequence,
		Denom:     data.Denom,
		Height:    ctx.Block.Height,
		Time:      ctx.Block.Time,
		TxId:      txId,
	}

	partsDenom := strings.Split(data.Denom, "/")
	if len(partsDenom) == 3 {
		transfer.Denom = partsDenom[2]
	}

	// height_timeout is bigint; values above MaxInt64 (e.g. MaxUint64 from Astria) mean no height timeout
	if packet.TimeoutHeight.RevisionHeight > 0 && packet.TimeoutHeight.RevisionHeight <= stdMath.MaxInt64 {
		transfer.HeightTimeout = packet.TimeoutHeight.RevisionHeight
	}
	if packet.TimeoutTimestamp > 0 {
		ts := math.TimeFromNano(packet.TimeoutTimestamp)
		transfer.Timeout = &ts
	}

	transfer.Receiver, transfer.ReceiverAddress, err = ibcTransferParty(ctx, data.Receiver, msgId, storageTypes.MsgAddressTypeReceiver)
	if err != nil {
		return errors.Wrap(err, "receiver")
	}
	transfer.Sender, transfer.SenderAddress, err = ibcTransferParty(ctx, data.Sender, msgId, storageTypes.MsgAddressTypeSender)
	if err != nil {
		return errors.Wrap(err, "sender")
	}
	ctx.AddIbcTransfer(msgId, transfer)
	return nil
}

func logSkippedPacket(ctx *context.Context, packet coreChannel.Packet, port, channel string, msgId uint64, err error) {
	log.Warn().
		Err(err).
		Int64("height", int64(ctx.Block.Height)).
		Uint64("msg_id", msgId).
		Str("port", port).
		Str("channel", channel).
		Uint64("sequence", packet.Sequence).
		Msg("skip invalid IBC packet data")
}

// ibcTransferParty indexes a celestia address; any other one, including non-bech32 (e.g. 0x...), is kept as a raw string
func ibcTransferParty(ctx *context.Context, address string, msgId uint64, typ storageTypes.MsgAddressType) (*storage.Address, *string, error) {
	prefix, hash, err := pkgTypes.Address(address).Decode()
	if err != nil || prefix != pkgTypes.AddressPrefixCelestia {
		return nil, &address, nil
	}
	addr := &storage.Address{
		Address:    address,
		Balances:   []storage.Balance{storage.EmptyBalance()},
		Height:     ctx.Block.Height,
		LastHeight: ctx.Block.Height,
		Hash:       hash,
	}
	if err := ctx.AddAddress(addr); err != nil {
		return nil, nil, errors.Wrap(err, "AddAddress")
	}
	ctx.AddAddressMessage(&storage.MsgAddress{
		MsgId:   msgId,
		Type:    typ,
		Address: addr,
	})
	return addr, nil, nil
}
