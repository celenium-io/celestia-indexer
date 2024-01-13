// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEvents_EmptyEventsResults(t *testing.T) {
	block := types.BlockData{
		ResultBlockResults: types.ResultBlockResults{
			TxsResults: make([]*types.ResponseDeliverTx, 0),
		},
	}

	resultEvents := parseEvents(block, make([]types.Event, 0))

	assert.Empty(t, resultEvents)
}

func TestParseEvents_SuccessTx(t *testing.T) {
	raw := `[{
		"type": "coin_spent",
		"attributes": [
			{
				"key": "c3BlbmRlcg==",
				"value": "Y2VsZXN0aWExcDMzMHN0YXB1c3lrZnNzNDdxcmhxbHVram5jdmd5emY2Z2R1ZnM=",
				"index": true
			},
			{
				"key": "YW1vdW50",
				"value": "NDA0OTR1dGlh",
				"index": true
			}
		]
	}]`
	var events []types.Event
	err := json.Unmarshal([]byte(raw), &events)
	require.NoError(t, err)

	txRes := types.ResponseDeliverTx{
		Code:      0,
		Data:      []byte{},
		Log:       "[]",
		Info:      "info",
		GasWanted: 12000,
		GasUsed:   1000,
		Events:    events,
		Codespace: "celestia-explorer",
	}
	block, now := testsuite.CreateTestBlock(txRes, 1)

	resultEvents := parseEvents(block, events)

	assert.Len(t, resultEvents, 1)

	e := resultEvents[0]
	assert.Equal(t, block.Height, e.Height)
	assert.Equal(t, now, e.Time)
	assert.Equal(t, int64(0), e.Position)
	assert.Equal(t, storageTypes.EventTypeCoinSpent, e.Type)
	assert.Nil(t, e.TxId)

	attrs := map[string]any{
		"spender": "celestia1p330stapusykfss47qrhqlukjncvgyzf6gdufs",
		"amount":  "40494utia",
	}
	assert.Equal(t, attrs, e.Data)
}

func BenchmarkParseEvent(b *testing.B) {
	block := types.BlockData{
		ResultBlock: types.ResultBlock{
			Block: &types.Block{
				Header: types.Header{
					Time: time.Now(),
				},
			},
		},
		ResultBlockResults: types.ResultBlockResults{
			Height: 100,
		},
	}
	raw := `{
		"type": "coin_spent",
		"attributes": [
			{
				"key": "c3BlbmRlcg==",
				"value": "Y2VsZXN0aWExcDMzMHN0YXB1c3lrZnNzNDdxcmhxbHVram5jdmd5emY2Z2R1ZnM=",
				"index": true
			},
			{
				"key": "YW1vdW50",
				"value": "NDA0OTR1dGlh",
				"index": true
			}
		]
	}`
	var event types.Event
	err := json.Unmarshal([]byte(raw), &event)
	require.NoError(b, err)

	resultEvent := storage.Event{}
	b.Run("parse event", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			parseEvent(block, event, 10, &resultEvent)
		}
	})
}
