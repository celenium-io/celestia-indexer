// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	testsuite "github.com/dipdup-io/celestia-indexer/internal/test_suite"
	"github.com/stretchr/testify/require"
)

func TestNewCoinSpent(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]any
		wantBody CoinSpent
		wantErr  bool
	}{
		{
			name: "test 1",
			m: map[string]any{
				"spender": "spender",
				"amount":  "1utia",
			},
			wantBody: CoinSpent{
				Spender: "spender",
				Amount:  testsuite.Ptr(types.NewCoin("utia", types.OneInt())),
			},
		}, {
			name: "test 2",
			m: map[string]any{
				"invalid": "invalid",
				"amount":  "1utia",
			},
			wantErr:  true,
			wantBody: CoinSpent{},
		}, {
			name: "test 3",
			m: map[string]any{
				"spender": "spender",
				"amount":  "invalid",
			},
			wantErr: true,
			wantBody: CoinSpent{
				Spender: "spender",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, err := NewCoinSpent(tt.m)
			require.True(t, (err != nil) == tt.wantErr)
			require.Equal(t, tt.wantBody, gotBody)
		})
	}
}

func TestNewCoinReceived(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]any
		wantBody CoinReceived
		wantErr  bool
	}{
		{
			name: "test 1",
			m: map[string]any{
				"receiver": "receiver",
				"amount":   "42utia",
			},
			wantBody: CoinReceived{
				Receiver: "receiver",
				Amount:   testsuite.Ptr(types.NewCoin("utia", types.NewInt(42))),
			},
		}, {
			name: "test 2",
			m: map[string]any{
				"invalid": "invalid",
				"amount":  "13utia",
			},
			wantErr:  true,
			wantBody: CoinReceived{},
		}, {
			name: "test 3",
			m: map[string]any{
				"receiver": "receiver",
				"amount":   "invalid",
			},
			wantErr: true,
			wantBody: CoinReceived{
				Receiver: "receiver",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, err := NewCoinReceived(tt.m)
			require.True(t, (err != nil) == tt.wantErr)
			require.Equal(t, tt.wantBody, gotBody)
		})
	}
}
