// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	coreClient "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	tmTypes "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
)

// MsgCreateClient defines a message to create an IBC client
func MsgCreateClient(ctx *context.Context, status storageTypes.Status, data storageTypes.PackedBytes, msgId uint64, m *coreClient.MsgCreateClient) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgCreateClient
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	if err != nil || status == storageTypes.StatusFailed {
		return msgType, err
	}

	if data == nil {
		return msgType, nil
	}

	if m.ClientState != nil {
		var clientState tmTypes.ClientState
		if err := clientState.Unmarshal(m.ClientState.Value); err != nil {
			return msgType, err
		}
		data["ClientState"] = clientState
	}

	if m.ConsensusState != nil {
		var consensusState tmTypes.ConsensusState
		if err := consensusState.Unmarshal(m.ConsensusState.Value); err != nil {
			return msgType, err
		}
		data["ConsensusState"] = consensusState
	}

	return msgType, nil
}

// MsgUpdateClient defines a sdk.Msg to update an IBC client state using the given header
func MsgUpdateClient(ctx *context.Context, status storageTypes.Status, data storageTypes.PackedBytes, msgId uint64, m *coreClient.MsgUpdateClient) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateClient
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	if err != nil || status == storageTypes.StatusFailed {
		return msgType, err
	}

	if data == nil {
		return msgType, nil
	}

	if m.ClientMessage != nil {
		var header tmTypes.Header
		if err := header.Unmarshal(m.ClientMessage.Value); err != nil {
			return msgType, err
		}
		data["Header"] = header
		delete(data, "ClientMessage")
	}

	return msgType, err
}

// MsgUpgradeClient defines a sdk.Msg to upgrade an IBC client to a new client state
func MsgUpgradeClient(ctx *context.Context, msgId uint64, m *coreClient.MsgUpgradeClient) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpgradeClient
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgSubmitMisbehaviour defines a sdk.Msg type that submits Evidence for light client misbehavior
func MsgSubmitMisbehaviour(ctx *context.Context, msgId uint64, m *coreClient.MsgSubmitMisbehaviour) (storageTypes.MsgType, error) { //nolint
	msgType := storageTypes.MsgSubmitMisbehaviour
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgRecoverClient defines the message used to recover a frozen or expired client.
func MsgRecoverClient(ctx *context.Context, msgId uint64, m *coreClient.MsgRecoverClient) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgRecoverClient
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgIBCSoftwareUpgrade defines the message used to schedule an upgrade of an IBC client using a v1 governance proposal
func MsgIBCSoftwareUpgrade(ctx *context.Context, msgId uint64, m *coreClient.MsgIBCSoftwareUpgrade) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgIBCSoftwareUpgrade
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgUpdateParams defines the sdk.Msg type to update the client parameters.
func MsgUpdateParams(ctx *context.Context, msgId uint64, m *coreClient.MsgUpdateParams) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateParams
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)
	return msgType, err
}
