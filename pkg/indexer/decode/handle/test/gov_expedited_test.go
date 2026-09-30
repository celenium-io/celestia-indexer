// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	json "github.com/bytedance/sonic"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	cosmosGovTypesV1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	cosmosGovTypesV1Beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	paramsV1Beta "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

const govAuthority = "celestia10d07y265gmmuvt4z0w9aw880jnsr700jtgz4v7"

func newDecodeContext() *context.Context {
	block, _ := testsuite.EmptyBlock()
	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}
	return decodeCtx
}

func expeditedGovParams() cosmosGovTypesV1.Params {
	votingPeriod := 7 * 24 * time.Hour
	expeditedVotingPeriod := 24 * time.Hour
	return cosmosGovTypesV1.Params{
		MinDeposit:            types.NewCoins(types.NewCoin("utia", math.NewInt(10_000_000_000))),
		ExpeditedMinDeposit:   types.NewCoins(types.NewCoin("utia", math.NewInt(50_000_000_000))),
		MaxDepositPeriod:      &votingPeriod,
		VotingPeriod:          &votingPeriod,
		ExpeditedVotingPeriod: &expeditedVotingPeriod,
		Quorum:                "0.334000000000000000",
		Threshold:             "0.500000000000000000",
		ExpeditedThreshold:    "0.667000000000000000",
		VetoThreshold:         "0.334000000000000000",
	}
}

func TestDecodeMsg_MsgSubmitProposalV1_Expedited(t *testing.T) {
	for _, expedited := range []bool{true, false} {
		m, ok := createMsgSubmitProposalV1().(*cosmosGovTypesV1.MsgSubmitProposal)
		require.True(t, ok)
		m.Expedited = expedited

		dm, err := decode.Message(newDecodeContext(), m, 0, storageTypes.StatusSuccess, 0)
		require.NoError(t, err)
		require.NotNil(t, dm.Msg.Proposal)
		require.NotNil(t, dm.Msg.Proposal.Expedited)
		require.Equal(t, expedited, *dm.Msg.Proposal.Expedited)
		require.Equal(t, expedited, dm.Msg.Proposal.IsExpedited())
	}
}

// v1beta1 has no expedited flag: nil keeps the stored value on upsert
func TestDecodeMsg_MsgSubmitProposalV1Beta1_ExpeditedIsNil(t *testing.T) {
	dm, err := decode.Message(newDecodeContext(), createMsgSubmitProposalV1Beta1(), 0, storageTypes.StatusSuccess, 0)
	require.NoError(t, err)
	require.NotNil(t, dm.Msg.Proposal)
	require.Nil(t, dm.Msg.Proposal.Expedited)
	require.False(t, dm.Msg.Proposal.IsExpedited())
}

func proposalChanges(t *testing.T, proposal *storage.Proposal) map[string]string {
	t.Helper()
	require.NotNil(t, proposal)
	var changes []paramsV1Beta.ParamChange
	require.NoError(t, json.Unmarshal(proposal.Changes, &changes))
	result := make(map[string]string, len(changes))
	for _, c := range changes {
		result[c.Subspace+"_"+c.Key] = c.Value
	}
	return result
}

func submitProposalWith(t *testing.T, typeURL string, msg interface{ Marshal() ([]byte, error) }) (*context.Context, *storage.Proposal) {
	t.Helper()
	value, err := msg.Marshal()
	require.NoError(t, err)

	m := &cosmosGovTypesV1.MsgSubmitProposal{
		Messages:       []*codecTypes.Any{{TypeUrl: typeURL, Value: value}},
		InitialDeposit: make([]types.Coin, 0),
		Proposer:       govAuthority,
	}

	ctx := newDecodeContext()
	dm, err := decode.Message(ctx, m, 0, storageTypes.StatusSuccess, 0)
	require.NoError(t, err)
	require.NotNil(t, dm.Msg.Proposal)
	return ctx, dm.Msg.Proposal
}

// gov params change on chain through a proposal; constants are applied only when it passes
func TestDecodeMsg_MsgSubmitProposalV1_GovUpdateParams(t *testing.T) {
	ctx, proposal := submitProposalWith(t, "/cosmos.gov.v1.MsgUpdateParams", &cosmosGovTypesV1.MsgUpdateParams{
		Authority: govAuthority,
		Params:    expeditedGovParams(),
	})

	require.Equal(t, storageTypes.ProposalTypeParamChanged, proposal.Type)
	require.EqualValues(t, 0, ctx.Constants.Len())

	changes := proposalChanges(t, proposal)
	require.Equal(t, "10000000000utia", changes["gov_min_deposit"])
	require.Equal(t, "604800000000000", changes["gov_voting_period"])
	require.Equal(t, "0.500000000000000000", changes["gov_threshold"])
	require.Equal(t, "86400000000000", changes["gov_expedited_voting_period"])
	require.Equal(t, "0.667000000000000000", changes["gov_expedited_threshold"])
	require.Equal(t, "50000000000utia", changes["gov_expedited_min_deposit"])
}

// uint32 params are JSON numbers in ProtoJSON
func TestDecodeMsg_MsgSubmitProposalV1_StakingUpdateParams(t *testing.T) {
	params := stakingTypes.DefaultParams()
	params.MaxValidators = 105
	params.BondDenom = "utia"

	_, proposal := submitProposalWith(t, "/cosmos.staking.v1beta1.MsgUpdateParams", &stakingTypes.MsgUpdateParams{
		Authority: govAuthority,
		Params:    params,
	})

	require.Equal(t, storageTypes.ProposalTypeParamChanged, proposal.Type)
	changes := proposalChanges(t, proposal)
	require.Equal(t, "105", changes["staking_max_validators"])
	require.Equal(t, "7", changes["staking_max_entries"])
	require.Equal(t, "10000", changes["staking_historical_entries"])
	require.Equal(t, "utia", changes["staking_bond_denom"])
	require.Equal(t, "1814400000000000", changes["staking_unbonding_time"])
	require.Equal(t, "0.000000000000000000", changes["staking_min_commission_rate"])
}

// a message without params must not drop the proposal
func TestDecodeMsg_MsgSubmitProposalV1_MessageWithoutParams(t *testing.T) {
	_, proposal := submitProposalWith(t, "/cosmos.gov.v1.MsgExecLegacyContent", &cosmosGovTypesV1.MsgExecLegacyContent{
		Authority: govAuthority,
	})

	require.Equal(t, storageTypes.ProposalTypeText, proposal.Type)
	require.Empty(t, proposal.Changes)
	require.Equal(t, "Proposal contains messages:\r\n1. /cosmos.gov.v1.MsgExecLegacyContent\r\n", proposal.Description)
}

// v1beta1 changes are stored normalized: they are applied as constants when the proposal passes
func TestDecodeMsg_MsgSubmitProposalV1Beta1_NormalizedChanges(t *testing.T) {
	content := paramsV1Beta.NewParameterChangeProposal("title", "description", []paramsV1Beta.ParamChange{
		paramsV1Beta.NewParamChange("staking", "MaxValidators", "103"),
		paramsV1Beta.NewParamChange("baseapp", "BlockParams", `{"max_bytes":"100","max_gas":"-1"}`),
		paramsV1Beta.NewParamChange("gov", "votingparams", `{"voting_period":"604800000000000"}`),
	})
	value, err := content.Marshal()
	require.NoError(t, err)

	m := &cosmosGovTypesV1Beta1.MsgSubmitProposal{
		Content:        &codecTypes.Any{TypeUrl: "/cosmos.params.v1beta1.ParameterChangeProposal", Value: value},
		InitialDeposit: make(types.Coins, 0),
		Proposer:       govAuthority,
	}

	ctx := newDecodeContext()
	dm, err := decode.Message(ctx, m, 0, storageTypes.StatusSuccess, 0)
	require.NoError(t, err)

	require.Equal(t, storageTypes.ProposalTypeParamChanged, dm.Msg.Proposal.Type)
	require.EqualValues(t, 0, ctx.Constants.Len())
	require.Equal(t, map[string]string{
		"staking_max_validators":    "103",
		"consensus_block_max_bytes": "100",
		"consensus_block_max_gas":   "-1",
		"gov_voting_period":         "604800000000000",
	}, proposalChanges(t, dm.Msg.Proposal))
}

// every message is listed in the description, even when it has no constants to keep
func TestDecodeMsg_MsgSubmitProposalV1_DescriptionListsAllMessages(t *testing.T) {
	send := &bankTypes.MsgSend{
		FromAddress: govAuthority,
		ToAddress:   govAuthority,
		Amount:      types.NewCoins(types.NewCoin("utia", math.NewInt(1))),
	}
	sendValue, err := send.Marshal()
	require.NoError(t, err)

	m := &cosmosGovTypesV1.MsgSubmitProposal{
		Messages: []*codecTypes.Any{
			// bank is not a module with constants
			{TypeUrl: "/cosmos.bank.v1beta1.MsgSend", Value: sendValue},
			{TypeUrl: "/cosmos.bank.v1beta1.MsgSend", Value: sendValue},
		},
		InitialDeposit: make([]types.Coin, 0),
		Proposer:       govAuthority,
	}

	dm, err := decode.Message(newDecodeContext(), m, 0, storageTypes.StatusSuccess, 0)
	require.NoError(t, err)
	require.NotNil(t, dm.Msg.Proposal)
	require.Equal(t, storageTypes.ProposalTypeText, dm.Msg.Proposal.Type)
	require.Equal(t,
		"Proposal contains messages:\r\n1. /cosmos.bank.v1beta1.MsgSend\r\n2. /cosmos.bank.v1beta1.MsgSend\r\n",
		dm.Msg.Proposal.Description,
	)
}
