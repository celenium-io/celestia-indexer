// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package parser

import (
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/stretchr/testify/require"
)

const (
	testFeeCollector = "celestia17xpfvakm2amg962yls6f84z3kell8c5lpnjs3s"
	testDistribution = "celestia1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8k44vnj"
	testSender       = "celestia1p330stapusykfss47qrhqlukjncvgyzf6gdufs"
	testSender2      = "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6"
)

func coinSpentEvent(spender, amount string) types.Event {
	return types.Event{
		Type: "coin_spent",
		Attributes: []types.EventAttribute{
			{Key: "spender", Value: spender},
			{Key: "amount", Value: amount},
		},
	}
}

func coinReceivedEvent(receiver, amount string) types.Event {
	return types.Event{
		Type: "coin_received",
		Attributes: []types.EventAttribute{
			{Key: "receiver", Value: receiver},
			{Key: "amount", Value: amount},
		},
	}
}

// feeEvents is what the ante handler emits: the signer pays the fee to the fee collector.
func feeEvents(sender, amount string) []types.Event {
	return []types.Event{
		coinSpentEvent(sender, amount),
		coinReceivedEvent(testFeeCollector, amount),
	}
}

func TestGetFirstTxEvent(t *testing.T) {
	first := feeEvents(testSender, "10utia")
	second := feeEvents(testSender2, "20utia")

	tests := []struct {
		name    string
		results []types.ResponseDeliverTx
		want    *types.Event
	}{
		{
			name:    "no txs",
			results: nil,
			want:    nil,
		}, {
			name:    "single tx",
			results: []types.ResponseDeliverTx{{Events: first}},
			want:    &first[0],
		}, {
			name:    "single tx without events",
			results: []types.ResponseDeliverTx{{Code: 1}},
			want:    nil,
		}, {
			name:    "several txs",
			results: []types.ResponseDeliverTx{{Events: first}, {Events: second}},
			want:    &first[0],
		}, {
			name:    "events only in the last tx",
			results: []types.ResponseDeliverTx{{Code: 1}, {Code: 1}, {Events: second}},
			want:    &second[0],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFirstTxEvent(tt.results)
			if tt.want == nil {
				require.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			require.True(t, tt.want.Compare(*got))
		})
	}
}

// Some nodes return tx events once more inside finalize_block_events, between begin and end block events.
func TestParseBlockEvents_DuplicatedTxEvents(t *testing.T) {
	begin := []types.Event{
		coinReceivedEvent(testFeeCollector, "100utia"),
	}
	end := []types.Event{
		coinSpentEvent(testFeeCollector, "130utia"),
		coinReceivedEvent(testDistribution, "130utia"),
	}

	tests := []struct {
		name string
		txs  []types.ResponseDeliverTx
		// duplicated tells if the node repeats tx events in block events
		duplicated bool
		want       map[string]string
	}{
		{
			name:       "single tx",
			txs:        []types.ResponseDeliverTx{{Events: feeEvents(testSender, "30utia")}},
			duplicated: true,
			want: map[string]string{
				testSender:       "-30",
				testFeeCollector: "0",
				testDistribution: "130",
			},
		}, {
			name: "several txs",
			txs: []types.ResponseDeliverTx{
				{Events: feeEvents(testSender, "10utia")},
				{Events: feeEvents(testSender2, "20utia")},
			},
			duplicated: true,
			want: map[string]string{
				testSender:       "-10",
				testSender2:      "-20",
				testFeeCollector: "0",
				testDistribution: "130",
			},
		}, {
			name: "events only in the last tx",
			txs: []types.ResponseDeliverTx{
				{Code: 1},
				{Events: feeEvents(testSender, "30utia")},
			},
			duplicated: true,
			want: map[string]string{
				testSender:       "-30",
				testFeeCollector: "0",
				testDistribution: "130",
			},
		}, {
			name:       "single tx without duplicates",
			txs:        []types.ResponseDeliverTx{{Events: feeEvents(testSender, "30utia")}},
			duplicated: false,
			want: map[string]string{
				testSender:       "-30",
				testFeeCollector: "0",
				testDistribution: "130",
			},
		}, {
			name:       "no txs",
			txs:        nil,
			duplicated: false,
			want: map[string]string{
				testFeeCollector: "-30",
				testDistribution: "130",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blockEvents := append([]types.Event{}, begin...)
			if tt.duplicated {
				for i := range tt.txs {
					blockEvents = append(blockEvents, tt.txs[i].Events...)
				}
			}
			blockEvents = append(blockEvents, end...)

			block := &types.BlockData{
				ResultBlock: types.ResultBlock{
					Block: &types.Block{
						Header: types.Header{Time: time.Now()},
					},
				},
				ResultBlockResults: types.ResultBlockResults{
					TxsResults:          tt.txs,
					FinalizeBlockEvents: blockEvents,
				},
			}

			// the same order as in Module.parse: tx events first, then block events
			ctx := context.NewContext()
			for i := range tt.txs {
				_, err := parseEvents(ctx, block, tt.txs[i].Events, uint64(i+1))
				require.NoError(t, err)
				ctx.TxEventsCount += len(tt.txs[i].Events)
			}

			result, err := parseBlockEvents(ctx, block, block.FinalizeBlockEvents, getFirstTxEvent(block.TxsResults))
			require.NoError(t, err)
			// duplicates are skipped by handlers, but still returned
			require.Len(t, result, len(blockEvents))

			require.Equal(t, len(tt.want), ctx.Addresses.Len())
			for address, want := range tt.want {
				addr, ok := ctx.Addresses.Get(address)
				require.True(t, ok, address)
				require.Len(t, addr.Balances, 1, address)
				require.Equal(t, "utia", addr.Balances[0].Currency, address)
				require.Equal(t, want, addr.Balances[0].Spendable.String(), address)
			}
		})
	}
}
