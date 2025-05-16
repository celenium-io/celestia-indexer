// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/legacy"
	coreClient "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	tmTypes "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
)

// MsgCreateClient defines a message to create an IBC client
func MsgCreateClient(ctx *context.Context, status types.Status, data types.PackedBytes, m *coreClient.MsgCreateClient) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreateClient
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	if err != nil || status == types.StatusFailed {
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

// MsgUpdateClientV6 defines a sdk.Msg to update an IBC client state using the given header
func MsgUpdateClientV6(ctx *context.Context, status types.Status, data types.PackedBytes, m *legacy.MsgUpdateClient) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUpdateClient
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	if err != nil || status == types.StatusFailed {
		return msgType, addresses, err
	}

	if data == nil {
		return msgType, addresses, nil
	}

	if m.Header != nil {
		var header tmTypes.Header
		if err := header.Unmarshal(m.Header.Value); err != nil {
			return msgType, addresses, err
		}
		data["Header"] = header
	}
	return msgType, addresses, err
}

// MsgUpdateClient defines a sdk.Msg to update an IBC client state using the given header
func MsgUpdateClient(ctx *context.Context, status types.Status, data types.PackedBytes, m *coreClient.MsgUpdateClient) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUpdateClient
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	if err != nil || status == types.StatusFailed {
		return msgType, addresses, err
	}

	if data == nil {
		return msgType, addresses, nil
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
func MsgSubmitMisbehaviour(ctx *context.Context, m *legacy.MsgSubmitMisbehaviour) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSubmitMisbehaviour
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
