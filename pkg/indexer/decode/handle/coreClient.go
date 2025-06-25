// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	coreClient "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	tmTypes "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
)

// MsgCreateClient defines a message to create an IBC client
func MsgCreateClient(ctx *context.Context, status storageTypes.Status, data storageTypes.PackedBytes, m *coreClient.MsgCreateClient) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreateClient
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	if err != nil || status == storageTypes.StatusFailed {
		return msgType, addresses, err
	}

	if data == nil {
		return msgType, addresses, nil
	}

	if m.ClientState != nil {
		var clientState tmTypes.ClientState
		if err := clientState.Unmarshal(m.ClientState.Value); err != nil {
			return msgType, addresses, err
		}
		data["ClientState"] = clientState
	}

	if m.ConsensusState != nil {
		var consensusState tmTypes.ConsensusState
		if err := consensusState.Unmarshal(m.ConsensusState.Value); err != nil {
			return msgType, addresses, err
		}
		data["ConsensusState"] = consensusState
	}

	return msgType, addresses, nil
}

// MsgUpdateClient defines a sdk.Msg to update an IBC client state using the given header
func MsgUpdateClient(ctx *context.Context, status storageTypes.Status, data storageTypes.PackedBytes, m *coreClient.MsgUpdateClient) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUpdateClient
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	if err != nil || status == storageTypes.StatusFailed {
		return msgType, addresses, err
	}

	if data == nil {
		return msgType, addresses, nil
	}

	if m.ClientMessage != nil {
		var header tmTypes.Header
		if err := header.Unmarshal(m.ClientMessage.Value); err != nil {
			return msgType, addresses, err
		}
		data["Header"] = header
		delete(data, "ClientMessage")
	}

	return msgType, addresses, err
}

// MsgUpgradeClient defines a sdk.Msg to upgrade an IBC client to a new client state
func MsgUpgradeClient(ctx *context.Context, m *coreClient.MsgUpgradeClient) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUpgradeClient
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgSubmitMisbehaviour defines a sdk.Msg type that submits Evidence for light client misbehavior
func MsgSubmitMisbehaviour(ctx *context.Context, m *coreClient.MsgSubmitMisbehaviour) (storageTypes.MsgType, []storage.AddressWithType, error) { //nolint
	msgType := storageTypes.MsgSubmitMisbehaviour
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgRecoverClient defines the message used to recover a frozen or expired client.
func MsgRecoverClient(ctx *context.Context, m *coreClient.MsgRecoverClient) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgRecoverClient
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgIBCSoftwareUpgrade defines the message used to schedule an upgrade of an IBC client using a v1 governance proposal
func MsgIBCSoftwareUpgrade(ctx *context.Context, m *coreClient.MsgIBCSoftwareUpgrade) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgIBCSoftwareUpgrade
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgUpdateParams defines the sdk.Msg type to update the client parameters.
func MsgUpdateParams(ctx *context.Context, m *coreClient.MsgUpdateParams) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUpdateParams
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
