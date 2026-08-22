// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/stretchr/testify/require"
)

const (
	tallyQuorumLegacy = "0.334"
	tallyQuorumUnused = "0.999"
)

func TestGov_GetDepositParams(t *testing.T) {
	t.Run("legacy format", func(t *testing.T) {
		gov := Gov{
			DepositParams: &DepositParams{
				MinDeposit:       []Coins{{Denom: currency.Utia, Amount: "100"}},
				MaxDepositPeriod: "172800s",
			},
			Params: &GovParams{
				MinDeposit:       []Coins{{Denom: currency.Utia, Amount: "999"}},
				MaxDepositPeriod: "999s",
			},
		}
		got, err := gov.GetDepositParams()
		require.NoError(t, err)
		require.Equal(t, "172800s", got.MaxDepositPeriod)
		require.Equal(t, []Coins{{Denom: currency.Utia, Amount: "100"}}, got.MinDeposit)
	})

	t.Run("new params format", func(t *testing.T) {
		gov := Gov{
			Params: &GovParams{
				MinDeposit:       []Coins{{Denom: currency.Utia, Amount: "500"}},
				MaxDepositPeriod: "86400s",
			},
		}
		got, err := gov.GetDepositParams()
		require.NoError(t, err)
		require.Equal(t, "86400s", got.MaxDepositPeriod)
		require.Equal(t, []Coins{{Denom: currency.Utia, Amount: "500"}}, got.MinDeposit)
	})

	t.Run("no params", func(t *testing.T) {
		gov := Gov{}
		got, err := gov.GetDepositParams()
		require.Error(t, err)
		require.Zero(t, got)
	})
}

func TestGov_GetVotingParams(t *testing.T) {
	t.Run("legacy format", func(t *testing.T) {
		gov := Gov{
			VotingParams: &VotingParams{VotingPeriod: "172800s"},
			Params:       &GovParams{VotingPeriod: "999s"},
		}
		got, err := gov.GetVotingParams()
		require.NoError(t, err)
		require.Equal(t, "172800s", got.VotingPeriod)
	})

	t.Run("new params format", func(t *testing.T) {
		gov := Gov{
			Params: &GovParams{VotingPeriod: "86400s"},
		}
		got, err := gov.GetVotingParams()
		require.NoError(t, err)
		require.Equal(t, "86400s", got.VotingPeriod)
	})

	t.Run("no params", func(t *testing.T) {
		gov := Gov{}
		got, err := gov.GetVotingParams()
		require.Error(t, err)
		require.Zero(t, got)
	})
}

func TestGov_GetTallyParams(t *testing.T) {
	t.Run("legacy format", func(t *testing.T) {
		gov := Gov{
			TallyParams: &TallyParams{
				Quorum:        tallyQuorumLegacy,
				Threshold:     "0.5",
				VetoThreshold: tallyQuorumLegacy,
			},
			Params: &GovParams{
				Quorum:        tallyQuorumUnused,
				Threshold:     tallyQuorumUnused,
				VetoThreshold: tallyQuorumUnused,
			},
		}
		got, err := gov.GetTallyParams()
		require.NoError(t, err)
		require.Equal(t, tallyQuorumLegacy, got.Quorum)
		require.Equal(t, "0.5", got.Threshold)
		require.Equal(t, tallyQuorumLegacy, got.VetoThreshold)
	})

	t.Run("new params format", func(t *testing.T) {
		gov := Gov{
			Params: &GovParams{
				Quorum:        tallyQuorumLegacy,
				Threshold:     "0.5",
				VetoThreshold: tallyQuorumLegacy,
			},
		}
		got, err := gov.GetTallyParams()
		require.NoError(t, err)
		require.Equal(t, tallyQuorumLegacy, got.Quorum)
		require.Equal(t, "0.5", got.Threshold)
		require.Equal(t, tallyQuorumLegacy, got.VetoThreshold)
	})

	t.Run("no params", func(t *testing.T) {
		gov := Gov{}
		got, err := gov.GetTallyParams()
		require.Error(t, err)
		require.Zero(t, got)
	})
}
