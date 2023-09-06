package parser

import (
	"testing"

	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	nodeTypes "github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEvents_EmptyEventsResults(t *testing.T) {
	block := types.BlockData{
		ResultBlockResults: nodeTypes.ResultBlockResults{
			TxsResults: make([]*nodeTypes.ResponseDeliverTx, 0),
		},
	}

	resultEvents := parseEvents(block, make([]nodeTypes.Event, 0))

	assert.Empty(t, resultEvents)
}

func TestParseEvents_SuccessTx(t *testing.T) {
	events := []nodeTypes.Event{
		{
			Type: "coin_spent",
			Attributes: []nodeTypes.EventAttribute{
				{
					Key:   "c3BlbmRlcg==",
					Value: "Y2VsZXN0aWExdjY5bnB6NncwN3h0NGhkdWU5eGR3a3V4eHZ2ZDZlYTl5MjZlcXI=",
					Index: true,
				},
				{
					Key:   "YW1vdW50",
					Value: "NzAwMDB1dGlh",
					Index: true,
				},
			},
		},
	}

	txRes := nodeTypes.ResponseDeliverTx{
		Code:      0,
		Data:      []byte{},
		Log:       "[]",
		Info:      "info",
		GasWanted: 12000,
		GasUsed:   1000,
		Events:    events,
		Codespace: "celestia-explorer",
	}
	block, now := createBlock(txRes, 1)

	resultEvents := parseEvents(block, events)

	assert.Len(t, resultEvents, 1)

	e := resultEvents[0]
	assert.Equal(t, block.Height, e.Height)
	assert.Equal(t, now, e.Time)
	assert.Equal(t, uint64(0), e.Position)
	assert.Equal(t, storageTypes.EventTypeCoinSpent, e.Type)
	assert.Nil(t, e.TxId)

	attrs := map[string]any{
		"spender": "celestia1v69npz6w07xt4hdue9xdwkuxxvvd6ea9y26eqr",
		"amount":  "70000utia",
	}
	assert.Equal(t, attrs, e.Data)
}

func Test_decodeEventAttribute(t *testing.T) {
	tests := []struct {
		name string
		data string
		want string
	}{
		{
			name: "test 1",
			data: "Y2VsZXN0aWExczQ1NXJoenh3Yzh3YzlrcXBoeHV0NzUyNHVtMDY3YzhwZGNjamo=",
			want: "celestia1s455rhzxwc8wc9kqphxut7524um067c8pdccjj",
		}, {
			name: "test 2",
			data: "c3BlbmRlcg==",
			want: "spender",
		}, {
			name: "test 3",
			data: "YW1vdW50",
			want: "amount",
		}, {
			name: "test 4",
			data: "NzAwMDB1dGlh",
			want: "70000utia",
		}, {
			name: "test 5",
			data: "bW9kdWxl",
			want: "module",
		}, {
			name: "test 6",
			data: "cmVjZWl2ZXI=",
			want: "receiver",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeEventAttribute(tt.data)
			require.Equal(t, tt.want, got)
		})
	}
}
