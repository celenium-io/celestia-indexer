// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	blobTypes "github.com/celestiaorg/celestia-app/v10/x/blob/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types"
	consensusTypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	cosmosGovTypesV1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/stretchr/testify/require"
)

// Mainnet proposal #7 "Set min_signed_per_window to .001": durations are stored in nanoseconds.
func TestDecodeMsg_MainnetProposal7_SlashingParams(t *testing.T) {
	_, proposal := submitProposalWith(t, "/cosmos.slashing.v1beta1.MsgUpdateParams", &slashingTypes.MsgUpdateParams{
		Authority: govAuthority,
		Params: slashingTypes.Params{
			SignedBlocksWindow:      10000,
			MinSignedPerWindow:      math.LegacyMustNewDecFromStr("0.001"),
			DowntimeJailDuration:    60 * time.Second,
			SlashFractionDoubleSign: math.LegacyMustNewDecFromStr("0.02"),
			SlashFractionDowntime:   math.LegacyZeroDec(),
		},
	})

	require.Equal(t, storageTypes.ProposalTypeParamChanged, proposal.Type)
	require.Equal(t, map[string]string{
		"slashing_signed_blocks_window":       "10000",
		"slashing_min_signed_per_window":      "0.001000000000000000",
		"slashing_downtime_jail_duration":     "60000000000",
		"slashing_slash_fraction_double_sign": "0.020000000000000000",
		"slashing_slash_fraction_downtime":    "0.000000000000000000",
	}, proposalChanges(t, proposal))
}

// Mainnet proposal #8 "Increase max square size to 256 and block size to 32 MiB": blob and consensus params in one proposal.
func TestDecodeMsg_MainnetProposal8_BlobAndConsensusParams(t *testing.T) {
	blobValue, err := (&blobTypes.MsgUpdateBlobParams{
		Authority: govAuthority,
		Params: blobTypes.Params{
			GasPerBlobByte:   8,
			GovMaxSquareSize: 256,
		},
	}).Marshal()
	require.NoError(t, err)

	consensusValue, err := (&consensusTypes.MsgUpdateParams{
		Authority: govAuthority,
		Block: &cmtproto.BlockParams{
			MaxBytes: 33554432,
			MaxGas:   -1,
		},
		Evidence: &cmtproto.EvidenceParams{
			MaxAgeNumBlocks: 242640,
			MaxAgeDuration:  1213200 * time.Second,
			MaxBytes:        1048576,
		},
		Validator: &cmtproto.ValidatorParams{
			PubKeyTypes: []string{"ed25519"},
		},
	}).Marshal()
	require.NoError(t, err)

	m := &cosmosGovTypesV1.MsgSubmitProposal{
		Messages: []*codecTypes.Any{
			{TypeUrl: "/celestia.blob.v1.MsgUpdateBlobParams", Value: blobValue},
			{TypeUrl: "/cosmos.consensus.v1.MsgUpdateParams", Value: consensusValue},
		},
		InitialDeposit: make([]types.Coin, 0),
		Proposer:       govAuthority,
		Title:          "Increase max square size to 256 and block size to 32 MiB",
		Summary:        "summary",
	}

	dm, err := decode.Message(newDecodeContext(), m, 0, storageTypes.StatusSuccess, 0)
	require.NoError(t, err)
	require.NotNil(t, dm.Msg.Proposal)
	require.Equal(t, storageTypes.ProposalTypeParamChanged, dm.Msg.Proposal.Type)
	require.Equal(t, map[string]string{
		"blob_gas_per_blob_byte":                "8",
		"blob_gov_max_square_size":              "256",
		"consensus_block_max_bytes":             "33554432",
		"consensus_block_max_gas":               "-1",
		"consensus_evidence_max_age_num_blocks": "242640",
		"consensus_evidence_max_bytes":          "1048576",
		"consensus_evidence_max_age_duration":   "1213200000000000",
	}, proposalChanges(t, dm.Msg.Proposal))
}
