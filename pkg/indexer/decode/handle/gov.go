// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"fmt"
	"strconv"
	"strings"

	consensusv1 "cosmossdk.io/api/cosmos/consensus/v1"
	slashingv1beta1 "cosmossdk.io/api/cosmos/slashing/v1beta1"
	"cosmossdk.io/math"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	blobTypes "github.com/celestiaorg/celestia-app/v6/x/blob/types"
	"github.com/cosmos/cosmos-sdk/codec"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	cosmosGovTypesV1Beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	paramsV1Beta "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	ibcTypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
)

// MsgSubmitProposalV1
func MsgSubmitProposalV1(ctx *context.Context, codec codec.Codec, status storageTypes.Status, msg *v1.MsgSubmitProposal) (storageTypes.MsgType, []storage.AddressWithType, []any, *storage.Proposal, error) {
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
		CreatedAt:   ctx.Block.Time,
		Status:      storageTypes.ProposalStatusInactive,
		Type:        storageTypes.ProposalTypeText,
		Title:       msg.Title,
		Description: msg.Summary,
		Metadata:    msg.Metadata,
	}

	if prpsl.Title == "" {
		prpsl.Title = "Proposal with messages"
	}

	changes := make(map[string]any, 0)
	var sb strings.Builder
	if _, err := sb.WriteString("Proposal contains messages:\r\n"); err != nil {
		return msgType, addresses, nil, nil, errors.Wrap(err, "building proposal description from messages")
	}
	for i := range msg.Messages {
		switch msg.Messages[i].TypeUrl {
		case "/cosmos.slashing.v1beta.MsgUpdateParams":
			var params slashingv1beta1.MsgUpdateParams
			if err := codec.Unmarshal(msg.Messages[i].Value, &params); err != nil {
				return msgType, addresses, nil, nil, errors.Wrap(err, "unmarshalling proposal with slashing.v1beta1.MsgUpdateParams")
			}
			if p := params.GetParams(); p != nil {
				paramChanges := make(map[string]any)
				prpsl.Type = storageTypes.ProposalTypeParamChanged

				var slashFractionDoubleSign math.LegacyDec
				if err := slashFractionDoubleSign.Unmarshal(p.GetSlashFractionDoubleSign()); err != nil {
					return msgType, addresses, nil, nil, errors.Wrap(err, "slash_fraction_double_sign")
				}
				ctx.AddConstant(storageTypes.ModuleNameSlashing, "slash_fraction_double_sign", slashFractionDoubleSign.String())
				paramChanges["slash_fraction_double_sign"] = slashFractionDoubleSign.String()

				var slashFractionDowntime math.LegacyDec
				if err := slashFractionDoubleSign.Unmarshal(p.GetSlashFractionDowntime()); err != nil {
					return msgType, addresses, nil, nil, errors.Wrap(err, "slash_fraction_downtime")
				}
				ctx.AddConstant(storageTypes.ModuleNameSlashing, "slash_fraction_downtime", slashFractionDowntime.String())
				paramChanges["slash_fraction_downtime"] = slashFractionDowntime.String()

				downtimeJailDuration := strconv.FormatInt(int64(p.GetDowntimeJailDuration().GetNanos()), 10)
				ctx.AddConstant(storageTypes.ModuleNameSlashing, "downtime_jail_duration", downtimeJailDuration)
				paramChanges["downtime_jail_duration"] = downtimeJailDuration

				minSignedPerWindow := string(p.GetMinSignedPerWindow())
				ctx.AddConstant(storageTypes.ModuleNameSlashing, "min_signed_per_window", minSignedPerWindow)
				paramChanges["min_signed_per_window"] = minSignedPerWindow

				signedBlocksWindow := strconv.FormatInt(p.GetSignedBlocksWindow(), 10)
				ctx.AddConstant(storageTypes.ModuleNameSlashing, "signed_blocks_window", signedBlocksWindow)
				paramChanges["signed_blocks_window"] = signedBlocksWindow

				changes[storageTypes.ModuleNameSlashing.String()] = paramChanges
			}

		case "/cosmos.distribution.v1beta1.MsgUpdateParams":
			var params distributionTypes.MsgUpdateParams
			if err := codec.Unmarshal(msg.Messages[i].Value, &params); err != nil {
				return msgType, addresses, nil, nil, errors.Wrap(err, "unmarshalling proposal with cosmos.distribution.v1beta1.MsgUpdateParams")
			}

			paramChanges := make(map[string]any)
			prpsl.Type = storageTypes.ProposalTypeParamChanged

			communityTax := params.Params.CommunityTax.String()
			ctx.AddConstant(storageTypes.ModuleNameDistribution, "community_tax", communityTax)
			paramChanges["community_tax"] = communityTax

			baseProposerReward := params.Params.BaseProposerReward.String() //nolint
			ctx.AddConstant(storageTypes.ModuleNameDistribution, "base_proposer_reward", baseProposerReward)
			paramChanges["base_proposer_reward"] = baseProposerReward

			bonusProposerReward := params.Params.BonusProposerReward.String() //nolint
			ctx.AddConstant(storageTypes.ModuleNameDistribution, "bonus_proposer_reward", bonusProposerReward)
			paramChanges["bonus_proposer_reward"] = bonusProposerReward

			withdrawAddrEnabled := strconv.FormatBool(params.Params.WithdrawAddrEnabled)
			ctx.AddConstant(storageTypes.ModuleNameDistribution, "withdraw_addr_enabled", withdrawAddrEnabled)
			paramChanges["withdraw_addr_enabled"] = withdrawAddrEnabled

			changes[storageTypes.ModuleNameDistribution.String()] = paramChanges

		case "/cosmos.gov.v1.MsgUpdateParams":
			var params v1.MsgUpdateParams
			if err := codec.Unmarshal(msg.Messages[i].Value, &params); err != nil {
				return msgType, addresses, nil, nil, errors.Wrap(err, "unmarshalling proposal with cosmos.gov.v1.MsgUpdateParams")
			}
			paramChanges := make(map[string]any)
			prpsl.Type = storageTypes.ProposalTypeParamChanged

			if minDeposits := params.Params.GetMinDeposit(); len(minDeposits) > 0 {
				ctx.AddConstant(storageTypes.ModuleNameGov, "min_deposit", minDeposits[0].Amount.String())
				paramChanges["min_deposit"] = minDeposits[0].Amount.String()
			}

			if maxDepositPeriod := params.Params.GetMaxDepositPeriod(); maxDepositPeriod != nil {
				value := strconv.FormatInt(maxDepositPeriod.Nanoseconds(), 10)
				ctx.AddConstant(storageTypes.ModuleNameGov, "max_deposit_period", value)
				paramChanges["max_deposit_period"] = value
			}

			if votingPeriod := params.Params.GetVotingPeriod(); votingPeriod != nil {
				value := strconv.FormatInt(votingPeriod.Nanoseconds(), 10)
				ctx.AddConstant(storageTypes.ModuleNameGov, "voting_period", value)
				paramChanges["voting_period"] = value
			}

			quorum := params.Params.GetQuorum()
			ctx.AddConstant(storageTypes.ModuleNameGov, "quorum", quorum)
			paramChanges["quorum"] = quorum

			threshold := params.Params.GetThreshold()
			ctx.AddConstant(storageTypes.ModuleNameGov, "threshold", threshold)
			paramChanges["threshold"] = threshold

			vetoThreshold := params.Params.GetVetoThreshold()
			ctx.AddConstant(storageTypes.ModuleNameGov, "veto_threshold", vetoThreshold)
			paramChanges["veto_threshold"] = vetoThreshold

			changes[storageTypes.ModuleNameGov.String()] = paramChanges

		case "/celestia.blob.v1.MsgUpdateBlobParams":
			var params blobTypes.MsgUpdateBlobParams
			if err := codec.Unmarshal(msg.Messages[i].Value, &params); err != nil {
				return msgType, addresses, nil, nil, errors.Wrap(err, "unmarshalling proposal with cosmos.gov.v1.MsgUpdateParams")
			}
			paramChanges := make(map[string]any)
			prpsl.Type = storageTypes.ProposalTypeParamChanged

			gasBlobPerByte := strconv.FormatUint(uint64(params.Params.GetGasPerBlobByte()), 10)
			ctx.AddConstant(storageTypes.ModuleNameBlob, "gas_per_blob_byte", gasBlobPerByte)
			paramChanges["gas_per_blob_byte"] = gasBlobPerByte

			maxSquareSize := strconv.FormatUint(params.Params.GetGovMaxSquareSize(), 10)
			ctx.AddConstant(storageTypes.ModuleNameBlob, "gov_max_square_size", maxSquareSize)
			paramChanges["gov_max_square_size"] = maxSquareSize

			changes[storageTypes.ModuleNameBlob.String()] = paramChanges

		case "/cosmos.consensus.v1.MsgUpdateParams":
			var params consensusv1.MsgUpdateParams
			if err := codec.Unmarshal(msg.Messages[i].Value, &params); err != nil {
				return msgType, addresses, nil, nil, errors.Wrap(err, "unmarshalling proposal with cosmos.gov.v1.MsgUpdateParams")
			}

			paramChanges := make(map[string]any)
			prpsl.Type = storageTypes.ProposalTypeParamChanged

			if block := params.GetBlock(); block != nil {
				maxBytes := strconv.FormatInt(block.GetMaxBytes(), 10)
				ctx.AddConstant(storageTypes.ModuleNameConsensus, "block_max_bytes", maxBytes)
				paramChanges["block_max_bytes"] = maxBytes

				maxGas := strconv.FormatInt(block.GetMaxGas(), 10)
				ctx.AddConstant(storageTypes.ModuleNameConsensus, "block_max_gas", maxGas)
				paramChanges["block_max_gas"] = maxGas
			}

			if evidence := params.GetEvidence(); evidence != nil {
				maxAgeNumBlocks := strconv.FormatInt(evidence.GetMaxAgeNumBlocks(), 10)
				ctx.AddConstant(storageTypes.ModuleNameConsensus, "evidence_max_age_num_blocks", maxAgeNumBlocks)
				paramChanges["block_max_bytes"] = maxAgeNumBlocks

				maxBytes := strconv.FormatInt(evidence.GetMaxBytes(), 10)
				ctx.AddConstant(storageTypes.ModuleNameConsensus, "evidence_max_bytes", maxBytes)
				paramChanges["evidence_max_bytes"] = maxBytes

				if age := evidence.GetMaxAgeDuration(); age != nil {
					value := strconv.FormatInt(int64(age.GetNanos()), 10)
					ctx.AddConstant(storageTypes.ModuleNameConsensus, "evidence_max_age_duration", value)
					paramChanges["evidence_max_age_duration"] = value
				}
			}

			changes[storageTypes.ModuleNameConsensus.String()] = paramChanges
		}

		if _, err := sb.WriteString(fmt.Sprintf("%d. %s\r\n", i+1, msg.Messages[i].TypeUrl)); err != nil {
			return msgType, addresses, nil, nil, errors.Wrap(err, "building proposal description from messages")
		}
	}
	if prpsl.Description == "" {
		prpsl.Description = sb.String()
	}
	if len(changes) > 0 {
		prpsl.Changes, err = json.Marshal(changes)
		if err != nil {
			return msgType, addresses, nil, nil, errors.Wrap(err, "marshalling changes proposal v1")
		}
	}
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
		var proposal distributionTypes.CommunityPoolSpendProposal //nolint
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
func MsgExecLegacyContent(ctx *context.Context, m *v1.MsgExecLegacyContent) (storageTypes.MsgType, []storage.AddressWithType, error) {
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

func MsgUpdateParamsGov(ctx *context.Context, m *v1.MsgUpdateParams) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUpdateParams
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
