// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"fmt"
	"strings"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/cosmos/cosmos-sdk/codec"
	cosmosDistrTypesV1Beta1 "github.com/cosmos/cosmos-sdk/x/distribution/types"
	cosmosGovTypesV1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	cosmosGovTypesV1Beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	paramsV1Beta "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	ibcTypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
)

// MsgSubmitProposalV1
func MsgSubmitProposalV1(ctx *context.Context, codec codec.Codec, status storageTypes.Status, msg *cosmosGovTypesV1.MsgSubmitProposal) (storageTypes.MsgType, []storage.AddressWithType, []any, *storage.Proposal, error) {
	msgType := storageTypes.MsgSubmitProposal
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeProposer, address: msg.Proposer},
	}, ctx.Block.Height)
	if err != nil {
		return msgType, addresses, nil, nil, err
	}

	if status != storageTypes.StatusSuccess {
		return msgType, addresses, nil, nil, nil
	}

	prpsl := &storage.Proposal{
		Height: ctx.Block.Height,
		Proposer: &storage.Address{
			Address: msg.Proposer,
		},
		CreatedAt: ctx.Block.Time,
		Status:    storageTypes.ProposalStatusInactive,
		Type:      storageTypes.ProposalTypeText,
		Title:     "Proposal with messages",
	}

	var sb strings.Builder
	if _, err := sb.WriteString("Proposal contains messages:\r\n"); err != nil {
		return msgType, addresses, nil, nil, errors.Wrap(err, "building proposal description from messages")
	}
	for i := range msg.Messages {
		if _, err := sb.WriteString(fmt.Sprintf("%d. %s\r\n", i+1, msg.Messages[i].TypeUrl)); err != nil {
			return msgType, addresses, nil, nil, errors.Wrap(err, "building proposal description from messages")
		}
	}
	prpsl.Description = sb.String()
	return msgType, addresses, nil, prpsl, nil
}

// MsgSubmitProposalV1Beta
func MsgSubmitProposalV1Beta(ctx *context.Context, codec codec.Codec, status storageTypes.Status, msg *cosmosGovTypesV1Beta1.MsgSubmitProposal) (storageTypes.MsgType, []storage.AddressWithType, any, *storage.Proposal, error) {
	msgType := storageTypes.MsgSubmitProposal
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeProposer, address: msg.Proposer},
	}, ctx.Block.Height)
	if err != nil {
		return msgType, addresses, nil, nil, err
	}
	if status != storageTypes.StatusSuccess {
		return msgType, addresses, nil, nil, nil
	}

	prpsl := &storage.Proposal{
		Height: ctx.Block.Height,
		Proposer: &storage.Address{
			Address: msg.Proposer,
		},
		CreatedAt: ctx.Block.Time,
		Status:    storageTypes.ProposalStatusInactive,
	}

	switch msg.Content.TypeUrl {
	case "/cosmos.gov.v1beta1.TextProposal":
		var proposal cosmosGovTypesV1Beta1.TextProposal
		if err := proposal.Unmarshal(msg.Content.Value); err != nil {
			return msgType, addresses, nil, nil, errors.Wrap(err, "unmarshalling text proposal for submit proposal content")
		}
		prpsl.Title = proposal.Title
		prpsl.Description = proposal.Description
		prpsl.Type = storageTypes.ProposalTypeText
		return msgType, addresses, proposal, prpsl, nil
	case "/cosmos.params.v1beta1.ParameterChangeProposal":
		var proposal paramsV1Beta.ParameterChangeProposal
		if err := proposal.Unmarshal(msg.Content.Value); err != nil {
			return msgType, addresses, nil, nil, errors.Wrap(err, "unmarshalling param change proposal for submit proposal content")
		}
		prpsl.Title = proposal.Title
		prpsl.Description = proposal.Description
		prpsl.Type = storageTypes.ProposalTypeParamChanged
		prpsl.Changes, err = json.Marshal(proposal.Changes)
		if err != nil {
			return msgType, addresses, nil, nil, errors.Wrap(err, "marshalling changes proposal for submit proposal content")
		}

		return msgType, addresses, proposal, prpsl, nil
	case "/ibc.core.client.v1.ClientUpdateProposal":
		var proposal ibcTypes.ClientUpdateProposal //nolint
		if err := proposal.Unmarshal(msg.Content.Value); err != nil {
			return msgType, addresses, nil, nil, errors.Wrap(err, "unmarshalling client update proposal for submit proposal content")
		}
		prpsl.Title = proposal.Title
		prpsl.Description = proposal.Description
		prpsl.Type = storageTypes.ProposalTypeClientUpdate
		prpsl.Changes, err = json.Marshal(map[string]any{
			"SubjectClientId":    proposal.SubjectClientId,
			"SubstituteClientId": proposal.SubstituteClientId,
		})
		if err != nil {
			return msgType, addresses, nil, nil, errors.Wrap(err, "marshalling changes proposal for submit proposal content")
		}
		return msgType, addresses, proposal, prpsl, nil

	case "/cosmos.distribution.v1beta1.CommunityPoolSpendProposal":
		var proposal cosmosDistrTypesV1Beta1.CommunityPoolSpendProposal //nolint
		if err := proposal.Unmarshal(msg.Content.Value); err != nil {
			return msgType, addresses, nil, nil, errors.Wrap(err, "unmarshalling community pool spend proposal for submit proposal content")
		}
		prpsl.Title = proposal.Title
		prpsl.Description = proposal.Description
		prpsl.Type = storageTypes.ProposalTypeCommunityPoolSpend
		prpsl.Changes, err = json.Marshal(map[string]any{
			"Recipient": proposal.Recipient,
			"Amount":    proposal.Amount,
		})
		if err != nil {
			return msgType, addresses, nil, nil, errors.Wrap(err, "marshalling changes proposal for submit proposal content")
		}
		return msgType, addresses, proposal, prpsl, nil

	default:
		return msgType, addresses, nil, nil, errors.Errorf("unknown content type in submit proposal: %s", msg.Content.TypeUrl)
	}
}

// MsgExecLegacyContent is used to wrap the legacy content field into a message.
// This ensures backwards compatibility with v1beta1.MsgSubmitProposal.
func MsgExecLegacyContent(ctx *context.Context, m *cosmosGovTypesV1.MsgExecLegacyContent) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgExecLegacyContent
	addresses, err := createAddresses(
		ctx,
		addressesData{
			{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
		}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgVote defines a message to cast a vote.
func MsgVote(ctx *context.Context, voterAddress string) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgVote
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeVoter, address: voterAddress},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgVoteWeighted defines a message to cast a vote.
func MsgVoteWeighted(ctx *context.Context, voterAddress string) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgVoteWeighted
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeVoter, address: voterAddress},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgDeposit defines a message to submit a deposit to an existing proposal.
func MsgDeposit(ctx *context.Context, depositorAddress string) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgDeposit
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeDepositor, address: depositorAddress},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
