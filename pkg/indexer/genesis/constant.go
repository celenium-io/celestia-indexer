// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package genesis

import (
	"strconv"
	"strings"
	"time"

	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/node/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
)

func (module *Module) parseConstants(ctx *decodeContext.Context, appState types.AppState, consensus pkgTypes.ConsensusParams) error {
	// consensus
	ctx.AddConstant(
		storageTypes.ModuleNameConsensus,
		"block_max_bytes",
		strconv.FormatInt(consensus.Block.MaxBytes, 10),
	)
	ctx.AddConstant(
		storageTypes.ModuleNameConsensus,
		"block_max_gas",
		strconv.FormatInt(consensus.Block.MaxGas, 10),
	)
	ctx.AddConstant(
		storageTypes.ModuleNameConsensus,
		"evidence_max_age_num_blocks",
		strconv.FormatInt(consensus.Evidence.MaxAgeNumBlocks, 10),
	)
	ctx.AddConstant(
		storageTypes.ModuleNameConsensus,
		"evidence_max_bytes",
		strconv.FormatInt(consensus.Evidence.MaxBytes, 10),
	)
	ctx.AddConstant(
		storageTypes.ModuleNameConsensus,
		"evidence_max_age_duration",
		strconv.FormatInt(consensus.Evidence.MaxAgeDuration.Nanoseconds(), 10),
	)
	ctx.AddConstant(
		storageTypes.ModuleNameConsensus,
		"validator_pub_key_types",
		strings.Join(consensus.Validator.PubKeyTypes, ", "),
	)

	// auth
	ctx.AddConstant(
		storageTypes.ModuleNameAuth,
		"max_memo_characters",
		appState.Auth.Params.MaxMemoCharacters,
	)
	ctx.AddConstant(
		storageTypes.ModuleNameAuth,
		"tx_sig_limit",
		appState.Auth.Params.TxSigLimit,
	)
	ctx.AddConstant(
		storageTypes.ModuleNameAuth,
		"tx_size_cost_per_byte",
		appState.Auth.Params.TxSizeCostPerByte,
	)
	ctx.AddConstant(
		storageTypes.ModuleNameAuth,
		"sig_verify_cost_ed25519",
		appState.Auth.Params.SigVerifyCostEd25519,
	)
	ctx.AddConstant(
		storageTypes.ModuleNameAuth,
		"sig_verify_cost_secp256k1",
		appState.Auth.Params.SigVerifyCostSecp256K1,
	)

	// blob
	ctx.AddConstant(
		storageTypes.ModuleNameBlob,
		"gas_per_blob_byte",
		strconv.FormatInt(int64(appState.Blob.Params.GasPerBlobByte), 10),
	)
	ctx.AddConstant(
		storageTypes.ModuleNameBlob,
		"gov_max_square_size",
		appState.Blob.Params.GovMaxSquareSize,
	)

	// crisis
	ctx.AddConstant(
		storageTypes.ModuleNameCrisis,
		"constant_fee",
		appState.Crisis.ConstantFee.String(),
	)

	// distribution
	ctx.AddConstant(
		storageTypes.ModuleNameDistribution,
		"community_tax",
		appState.Distribution.Params.CommunityTax,
	)
	ctx.AddConstant(
		storageTypes.ModuleNameDistribution,
		"base_proposer_reward",
		appState.Distribution.Params.BaseProposerReward,
	)
	ctx.AddConstant(
		storageTypes.ModuleNameDistribution,
		"bonus_proposer_reward",
		appState.Distribution.Params.BonusProposerReward,
	)
	ctx.AddConstant(
		storageTypes.ModuleNameDistribution,
		"withdraw_addr_enabled",
		strconv.FormatBool(appState.Distribution.Params.WithdrawAddrEnabled),
	)

	// gov
	govDepositParams, err := appState.Gov.GetDepositParams()
	if err != nil {
		return err
	}

	if len(govDepositParams.MinDeposit) > 0 {
		ctx.AddConstant(
			storageTypes.ModuleNameGov,
			"min_deposit",
			govDepositParams.MinDeposit[0].String(),
		)
	}

	maxDepositPeriod, err := time.ParseDuration(govDepositParams.MaxDepositPeriod)
	if err != nil {
		return errors.Wrap(err, "max deposit period")
	}

	ctx.AddConstant(
		storageTypes.ModuleNameGov,
		"max_deposit_period",
		strconv.FormatInt(maxDepositPeriod.Nanoseconds(), 10),
	)

	votingParams, err := appState.Gov.GetVotingParams()
	if err != nil {
		return err
	}
	votingPeriod, err := time.ParseDuration(votingParams.VotingPeriod)
	if err != nil {
		return errors.Wrap(err, "voting period")
	}

	ctx.AddConstant(
		storageTypes.ModuleNameGov,
		"voting_period",
		strconv.FormatInt(votingPeriod.Nanoseconds(), 10),
	)

	tallyParams, err := appState.Gov.GetTallyParams()
	if err != nil {
		return err
	}
	ctx.AddConstant(
		storageTypes.ModuleNameGov,
		"quorum",
		tallyParams.Quorum,
	)
	ctx.AddConstant(
		storageTypes.ModuleNameGov,
		"threshold",
		tallyParams.Threshold,
	)
	ctx.AddConstant(
		storageTypes.ModuleNameGov,
		"veto_threshold",
		tallyParams.VetoThreshold,
	)

	// slashing
	ctx.AddConstant(
		storageTypes.ModuleNameSlashing,
		"signed_blocks_window",
		appState.Slashing.Params.SignedBlocksWindow,
	)
	ctx.AddConstant(
		storageTypes.ModuleNameSlashing,
		"min_signed_per_window",
		appState.Slashing.Params.MinSignedPerWindow,
	)
	downtimeJailDuration, err := time.ParseDuration(appState.Slashing.Params.DowntimeJailDuration)
	if err != nil {
		return errors.Wrap(err, "DowntimeJailDuration")
	}
	ctx.AddConstant(
		storageTypes.ModuleNameSlashing,
		"downtime_jail_duration",
		strconv.FormatInt(downtimeJailDuration.Nanoseconds(), 10),
	)
	ctx.AddConstant(
		storageTypes.ModuleNameSlashing,
		"slash_fraction_double_sign",
		appState.Slashing.Params.SlashFractionDoubleSign,
	)
	ctx.AddConstant(
		storageTypes.ModuleNameSlashing,
		"slash_fraction_downtime",
		appState.Slashing.Params.SlashFractionDowntime,
	)

	// staking
	unbondingTime, err := time.ParseDuration(appState.Staking.Params.UnbondingTime)
	if err != nil {
		return errors.Wrap(err, "unbonding time")
	}

	ctx.AddConstant(
		storageTypes.ModuleNameStaking,
		"unbonding_time",
		strconv.FormatInt(unbondingTime.Nanoseconds(), 10),
	)
	ctx.AddConstant(
		storageTypes.ModuleNameStaking,
		"max_validators",
		strconv.FormatInt(int64(appState.Staking.Params.MaxValidators), 10),
	)
	ctx.AddConstant(
		storageTypes.ModuleNameStaking,
		"max_entries",
		strconv.FormatInt(int64(appState.Staking.Params.MaxEntries), 10),
	)
	ctx.AddConstant(
		storageTypes.ModuleNameStaking,
		"historical_entries",
		strconv.FormatInt(int64(appState.Staking.Params.HistoricalEntries), 10),
	)
	ctx.AddConstant(
		storageTypes.ModuleNameStaking,
		"bond_denom",
		appState.Staking.Params.BondDenom,
	)
	ctx.AddConstant(
		storageTypes.ModuleNameStaking,
		"min_commission_rate",
		appState.Staking.Params.MinCommissionRate,
	)

	// minfee
	ctx.AddConstant(
		storageTypes.ModuleNameMinfee,
		"network_min_gas_price",
		appState.MinFee.GetNetworkMinGasPrice(),
	)

	return nil
}
