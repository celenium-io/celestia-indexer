// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package genesis

import (
	"strconv"
	"strings"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/node/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
)

func (module *Module) parseConstants(appState types.AppState, consensus pkgTypes.ConsensusParams, data *parsedData) error {
	// consensus
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameConsensus,
		Name:   "block_max_bytes",
		Value:  strconv.FormatInt(consensus.Block.MaxBytes, 10),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameConsensus,
		Name:   "block_max_gas",
		Value:  strconv.FormatInt(consensus.Block.MaxGas, 10),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameConsensus,
		Name:   "evidence_max_age_num_blocks",
		Value:  strconv.FormatInt(consensus.Evidence.MaxAgeNumBlocks, 10),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameConsensus,
		Name:   "evidence_max_bytes",
		Value:  strconv.FormatInt(consensus.Evidence.MaxBytes, 10),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameConsensus,
		Name:   "evidence_max_age_duration",
		Value:  strconv.FormatInt(consensus.Evidence.MaxAgeDuration.Nanoseconds(), 10),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameConsensus,
		Name:   "validator_pub_key_types",
		Value:  strings.Join(consensus.Validator.PubKeyTypes, ", "),
	})

	// auth
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameAuth,
		Name:   "max_memo_characters",
		Value:  appState.Auth.Params.MaxMemoCharacters,
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameAuth,
		Name:   "tx_sig_limit",
		Value:  appState.Auth.Params.TxSigLimit,
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameAuth,
		Name:   "tx_size_cost_per_byte",
		Value:  appState.Auth.Params.TxSizeCostPerByte,
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameAuth,
		Name:   "sig_verify_cost_ed25519",
		Value:  appState.Auth.Params.SigVerifyCostEd25519,
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameAuth,
		Name:   "sig_verify_cost_secp256k1",
		Value:  appState.Auth.Params.SigVerifyCostSecp256K1,
	})

	// blob
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameBlob,
		Name:   "gas_per_blob_byte",
		Value:  strconv.FormatInt(int64(appState.Blob.Params.GasPerBlobByte), 10),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameBlob,
		Name:   "gov_max_square_size",
		Value:  appState.Blob.Params.GovMaxSquareSize,
	})

	// crisis
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameCrisis,
		Name:   "constant_fee",
		Value:  appState.Crisis.ConstantFee.String(),
	})

	// distribution
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameDistribution,
		Name:   "community_tax",
		Value:  appState.Distribution.Params.CommunityTax,
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameDistribution,
		Name:   "base_proposer_reward",
		Value:  appState.Distribution.Params.BaseProposerReward,
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameDistribution,
		Name:   "bonus_proposer_reward",
		Value:  appState.Distribution.Params.BonusProposerReward,
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameDistribution,
		Name:   "withdraw_addr_enabled",
		Value:  strconv.FormatBool(appState.Distribution.Params.WithdrawAddrEnabled),
	})

	// gov
	if len(appState.Gov.DepositParams.MinDeposit) > 0 {
		data.constants = append(data.constants, storage.Constant{
			Module: storageTypes.ModuleNameGov,
			Name:   "min_deposit",
			Value:  appState.Gov.DepositParams.MinDeposit[0].String(),
		})
	}

	maxDepositPeriod, err := time.ParseDuration(appState.Gov.DepositParams.MaxDepositPeriod)
	if err != nil {
		return errors.Wrap(err, "max deposit period")
	}

	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGov,
		Name:   "max_deposit_period",
		Value:  strconv.FormatInt(maxDepositPeriod.Nanoseconds(), 10),
	})

	votingPower, err := time.ParseDuration(appState.Gov.VotingParams.VotingPeriod)
	if err != nil {
		return errors.Wrap(err, "voting power")
	}

	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGov,
		Name:   "voting_period",
		Value:  strconv.FormatInt(votingPower.Nanoseconds(), 10),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGov,
		Name:   "quorum",
		Value:  appState.Gov.TallyParams.Quorum,
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGov,
		Name:   "threshold",
		Value:  appState.Gov.TallyParams.Threshold,
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGov,
		Name:   "veto_threshold",
		Value:  appState.Gov.TallyParams.VetoThreshold,
	})

	// slashing
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameSlashing,
		Name:   "signed_blocks_window",
		Value:  appState.Slashing.Params.SignedBlocksWindow,
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameSlashing,
		Name:   "min_signed_per_window",
		Value:  appState.Slashing.Params.MinSignedPerWindow,
	})
	downtimeJailDuration, err := time.ParseDuration(appState.Slashing.Params.DowntimeJailDuration)
	if err != nil {
		return errors.Wrap(err, "DowntimeJailDuration")
	}
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameSlashing,
		Name:   "downtime_jail_duration",
		Value:  strconv.FormatInt(downtimeJailDuration.Nanoseconds(), 10),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameSlashing,
		Name:   "slash_fraction_double_sign",
		Value:  appState.Slashing.Params.SlashFractionDoubleSign,
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameSlashing,
		Name:   "slash_fraction_downtime",
		Value:  appState.Slashing.Params.SlashFractionDowntime,
	})

	// staking
	unbondingTime, err := time.ParseDuration(appState.Staking.Params.UnbondingTime)
	if err != nil {
		return errors.Wrap(err, "unbonding time")
	}

	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameStaking,
		Name:   "unbonding_time",
		Value:  strconv.FormatInt(unbondingTime.Nanoseconds(), 10),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameStaking,
		Name:   "max_validators",
		Value:  strconv.FormatInt(int64(appState.Staking.Params.MaxValidators), 10),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameStaking,
		Name:   "max_entries",
		Value:  strconv.FormatInt(int64(appState.Staking.Params.MaxEntries), 10),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameStaking,
		Name:   "historical_entries",
		Value:  strconv.FormatInt(int64(appState.Staking.Params.HistoricalEntries), 10),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameStaking,
		Name:   "bond_denom",
		Value:  appState.Staking.Params.BondDenom,
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameStaking,
		Name:   "min_commission_rate",
		Value:  appState.Staking.Params.MinCommissionRate,
	})

	return nil
}
