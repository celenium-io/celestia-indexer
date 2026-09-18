// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/celestiaorg/celestia-app/v10/app"
	"github.com/celestiaorg/celestia-app/v10/app/encoding"
	fibreTypes "github.com/celestiaorg/celestia-app/v10/x/fibre/types"
	"github.com/stretchr/testify/require"
)

// TestAppState_UnmarshalFibre pins the genesis shape against the codec that
// actually produces it, so a change in how the app serializes the fibre params
// fails here instead of silently yielding empty constants.
func TestAppState_UnmarshalFibre(t *testing.T) {
	cfg := encoding.MakeConfig(app.ModuleEncodingRegisters...)

	genesis := fibreTypes.DefaultGenesis()
	genesis.Params.PaymentPromiseHeightWindow = 2500
	genesis.Params.WithdrawalDelay = 36 * time.Hour

	raw, err := cfg.Codec.MarshalJSON(genesis)
	require.NoError(t, err)

	var appState AppState
	require.NoError(t, json.Unmarshal([]byte(`{"fibre":`+string(raw)+`}`), &appState))
	require.NotNil(t, appState.Fibre)

	// Durations arrive as protobuf duration strings, uint64 as decimal strings.
	require.Equal(t, "129600s", appState.Fibre.Params.WithdrawalDelay)
	require.Equal(t, "3600s", appState.Fibre.Params.PaymentPromiseTimeout)
	require.Equal(t, "2500", appState.Fibre.Params.PaymentPromiseHeightWindow)
	require.Equal(t, "14400s", appState.Fibre.Params.ShardRetention)
	require.Equal(t, "2199023255552", appState.Fibre.Params.FullStakeStorageBudget)

	parsed, err := time.ParseDuration(appState.Fibre.Params.WithdrawalDelay)
	require.NoError(t, err)
	require.Equal(t, genesis.Params.WithdrawalDelay, parsed)
}

// TestAppState_FibreMissing covers genesis of a chain that predates the module.
func TestAppState_FibreMissing(t *testing.T) {
	var appState AppState
	require.NoError(t, json.Unmarshal([]byte(`{"minfee":{"network_min_gas_price":"0.000001"}}`), &appState))
	require.Nil(t, appState.Fibre)
}
