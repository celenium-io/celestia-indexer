// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	cosmosGovTypesV1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// MsgSubmitProposal defines a sdk.Msg type that supports submitting arbitrary
// proposal Content.
func MsgSubmitProposal(ctx *context.Context, proposerAddress string) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSubmitProposal
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeProposer, address: proposerAddress},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgExecLegacyContent is used to wrap the legacy content field into a message.
// This ensures backwards compatibility with v1beta1.MsgSubmitProposal.
func MsgExecLegacyContent(ctx *context.Context, m *cosmosGovTypesV1.MsgExecLegacyContent) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgExecLegacyContent
	addresses, err := createAddresses(
		addressesData{
			{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
		}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgVote defines a message to cast a vote.
func MsgVote(ctx *context.Context, voterAddress string) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgVote
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeVoter, address: voterAddress},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgVoteWeighted defines a message to cast a vote.
func MsgVoteWeighted(ctx *context.Context, voterAddress string) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgVoteWeighted
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeVoter, address: voterAddress},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgDeposit defines a message to submit a deposit to an existing proposal.
func MsgDeposit(ctx *context.Context, depositorAddress string) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgDeposit
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDepositor, address: depositorAddress},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
