// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package genesis

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/node/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/stretchr/testify/require"
)

const notADuration = "not-a-duration"

func mustNanos(t *testing.T, s string) string {
	t.Helper()
	d, err := time.ParseDuration(s)
	require.NoError(t, err)
	return strconv.FormatInt(d.Nanoseconds(), 10)
}

func testAppState() types.AppState {
	return types.AppState{
		Auth: types.Auth{
			Params: types.AuthParams{
				MaxMemoCharacters:      "256",
				TxSigLimit:             "7",
				TxSizeCostPerByte:      "10",
				SigVerifyCostEd25519:   "590",
				SigVerifyCostSecp256K1: "1000",
			},
		},
		Blob: types.BlobState{
			Params: types.BlobParams{
				GasPerBlobByte:   8,
				GovMaxSquareSize: "64",
			},
		},
		Crisis: types.Crisis{
			ConstantFee: types.Coins{Denom: "utia", Amount: "1000"},
		},
		Distribution: types.Distribution{
			Params: types.DistributionParams{
				CommunityTax:        "0.02",
				BaseProposerReward:  "0",
				BonusProposerReward: "0",
				WithdrawAddrEnabled: true,
			},
		},
		Gov: types.Gov{
			Params: &types.GovParams{
				MinDeposit:       []types.Coins{{Denom: "utia", Amount: "10000000"}},
				MaxDepositPeriod: "172800s",
				VotingPeriod:     "115200s",
				Quorum:           "0.334",
				Threshold:        "0.5",
				VetoThreshold:    "0.334",
			},
		},
		Slashing: types.Slashing{
			Params: types.SlashingParams{
				SignedBlocksWindow:      "5000",
				MinSignedPerWindow:      "0.1",
				DowntimeJailDuration:    "600s",
				SlashFractionDoubleSign: "0.05",
				SlashFractionDowntime:   "0.01",
			},
		},
		Staking: types.Staking{
			Params: types.StakingParams{
				UnbondingTime:     "1814400s",
				MaxValidators:     100,
				MaxEntries:        7,
				HistoricalEntries: 10000,
				BondDenom:         "utia",
				MinCommissionRate: "0.05",
			},
		},
		MinFee: types.MinFee{
			NetworkMinGasPrice: "0.000001",
		},
	}
}

func testConsensusParams() pkgTypes.ConsensusParams {
	return pkgTypes.ConsensusParams{
		Block: &pkgTypes.BlockParams{
			MaxBytes: 21000000,
			MaxGas:   -1,
		},
		Evidence: &pkgTypes.EvidenceParams{
			MaxAgeNumBlocks: 100000,
			MaxAgeDuration:  48 * time.Hour,
			MaxBytes:        1048576,
		},
		Validator: &pkgTypes.ValidatorParams{
			PubKeyTypes: []string{"ed25519", "secp256k1"},
		},
	}
}

// constantsMap materializes ctx.Constants into "module.name" -> value for comparison.
func constantsMap(ctx *decodeContext.Context) map[string]string {
	got := make(map[string]string, ctx.Constants.Len())
	for _, c := range ctx.Constants.Values() {
		got[fmt.Sprintf("%s.%s", c.Module, c.Name)] = c.Value
	}
	return got
}

func TestParseConstants(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})
	ctx := decodeContext.NewContext()
	appState := testAppState()
	consensus := testConsensusParams()

	err := module.parseConstants(ctx, appState, consensus)
	require.NoError(t, err)

	want := map[string]string{
		"consensus.block_max_bytes":             "21000000",
		"consensus.block_max_gas":               "-1",
		"consensus.evidence_max_age_num_blocks": "100000",
		"consensus.evidence_max_bytes":          "1048576",
		"consensus.evidence_max_age_duration":   strconv.FormatInt((48 * time.Hour).Nanoseconds(), 10),
		"consensus.validator_pub_key_types":     "ed25519, secp256k1",

		"auth.max_memo_characters":       "256",
		"auth.tx_sig_limit":              "7",
		"auth.tx_size_cost_per_byte":     "10",
		"auth.sig_verify_cost_ed25519":   "590",
		"auth.sig_verify_cost_secp256k1": "1000",

		"blob.gas_per_blob_byte":   "8",
		"blob.gov_max_square_size": "64",

		"crisis.constant_fee": "1000utia",

		"distribution.community_tax":         "0.02",
		"distribution.base_proposer_reward":  "0",
		"distribution.bonus_proposer_reward": "0",
		"distribution.withdraw_addr_enabled": "true",

		"gov.min_deposit":        "10000000utia",
		"gov.max_deposit_period": mustNanos(t, "172800s"),
		"gov.voting_period":      mustNanos(t, "115200s"),
		"gov.quorum":             "0.334",
		"gov.threshold":          "0.5",
		"gov.veto_threshold":     "0.334",

		"slashing.signed_blocks_window":       "5000",
		"slashing.min_signed_per_window":      "0.1",
		"slashing.downtime_jail_duration":     mustNanos(t, "600s"),
		"slashing.slash_fraction_double_sign": "0.05",
		"slashing.slash_fraction_downtime":    "0.01",

		"staking.unbonding_time":      mustNanos(t, "1814400s"),
		"staking.max_validators":      "100",
		"staking.max_entries":         "7",
		"staking.historical_entries":  "10000",
		"staking.bond_denom":          "utia",
		"staking.min_commission_rate": "0.05",

		"minfee.network_min_gas_price": "0.000001",
	}

	require.Equal(t, want, constantsMap(ctx))
}

func TestParseConstants_ZeroMinDepositOmitted(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})
	ctx := decodeContext.NewContext()
	appState := testAppState()
	appState.Gov.Params.MinDeposit = nil

	err := module.parseConstants(ctx, appState, testConsensusParams())
	require.NoError(t, err)
	require.NotContains(t, constantsMap(ctx), "gov.min_deposit")
}

func TestParseConstants_InvalidDuration(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*types.AppState)
	}{
		{
			name: "max_deposit_period",
			mutate: func(a *types.AppState) {
				a.Gov.Params.MaxDepositPeriod = notADuration
			},
		},
		{
			name: "voting_period",
			mutate: func(a *types.AppState) {
				a.Gov.Params.VotingPeriod = notADuration
			},
		},
		{
			name: "downtime_jail_duration",
			mutate: func(a *types.AppState) {
				a.Slashing.Params.DowntimeJailDuration = notADuration
			},
		},
		{
			name: "unbonding_time",
			mutate: func(a *types.AppState) {
				a.Staking.Params.UnbondingTime = notADuration
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			module := NewModule(postgres.Storage{}, config.Indexer{})
			ctx := decodeContext.NewContext()
			appState := testAppState()
			tc.mutate(&appState)

			err := module.parseConstants(ctx, appState, testConsensusParams())
			require.Error(t, err)
		})
	}
}

func TestParseConstants_MinFeeNewParamsFormatTakesPriority(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})
	ctx := decodeContext.NewContext()
	appState := testAppState()
	appState.MinFee.Params.NetworkMinGasPrice = "0.00002"

	err := module.parseConstants(ctx, appState, testConsensusParams())
	require.NoError(t, err)
	require.Equal(t, "0.00002", constantsMap(ctx)["minfee.network_min_gas_price"])
}

func TestParseConstants_MissingGovParams(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})
	ctx := decodeContext.NewContext()
	appState := testAppState()
	appState.Gov.Params = nil

	err := module.parseConstants(ctx, appState, testConsensusParams())
	require.Error(t, err)
}
