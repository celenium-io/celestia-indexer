// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"encoding/base64"
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestParseTxs_EmptyTxsResults(t *testing.T) {
	block := types.BlockData{
		ResultBlockResults: types.ResultBlockResults{
			TxsResults: make([]*types.ResponseDeliverTx, 0),
		},
	}

	decodeCtx := context.NewContext()
	resultTxs, err := parseTxs(decodeCtx, block)

	assert.NoError(t, err)
	assert.Empty(t, resultTxs)
}

func mustDecodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func TestParseTxs_SuccessTx(t *testing.T) {
	txRes := types.ResponseDeliverTx{
		Code:      0,
		Data:      []byte{},
		Log:       "[]",
		Info:      "info",
		GasWanted: 12000,
		GasUsed:   1000,
		Codespace: "celestia-explorer",
		Events: []types.Event{
			{
				Type: "message",
				Attributes: []types.EventAttribute{
					{
						Key:   mustDecodeBase64("YWN0aW9u"),
						Value: mustDecodeBase64("L2Nvc21vcy5zdGFraW5nLnYxYmV0YTEuTXNnQmVnaW5SZWRlbGVnYXRl"),
						Index: true,
					},
				},
			},
			{
				Type: "coin_spent",
				Attributes: []types.EventAttribute{
					{
						Key:   mustDecodeBase64("c3BlbmRlcg=="),
						Value: mustDecodeBase64("Y2VsZXN0aWExanY2NXMzZ3JxZjZ2NmpsM2RwNHQ2Yzl0OXJrOTljZDhrNDR2bmo="),
						Index: true,
					}, {
						Key:   mustDecodeBase64("YW1vdW50"),
						Value: mustDecodeBase64("NDI3NzMxdXRpYQ=="),
						Index: true,
					},
				},
			},
			{
				Type: "coin_received",
				Attributes: []types.EventAttribute{
					{
						Key:   mustDecodeBase64("cmVjZWl2ZXI="),
						Value: mustDecodeBase64("Y2VsZXN0aWExMjUzdmRsZGxmeGx3eXBuZGgzZnp6cXQ5ZG4wcjV3MGRldHY1dWg="),
						Index: true,
					}, {
						Key:   mustDecodeBase64("YW1vdW50"),
						Value: mustDecodeBase64("NDI3NzMxdXRpYQ=="),
						Index: true,
					},
				},
			},
			{
				Type: "transfer",
				Attributes: []types.EventAttribute{
					{
						Key:   mustDecodeBase64("cmVjaXBpZW50"),
						Value: mustDecodeBase64("Y2VsZXN0aWExMjUzdmRsZGxmeGx3eXBuZGgzZnp6cXQ5ZG4wcjV3MGRldHY1dWg="),
						Index: true,
					},
					{
						Key:   mustDecodeBase64("c2VuZGVy"),
						Value: mustDecodeBase64("Y2VsZXN0aWExanY2NXMzZ3JxZjZ2NmpsM2RwNHQ2Yzl0OXJrOTljZDhrNDR2bmo="),
						Index: true,
					},
					{
						Key:   mustDecodeBase64("YW1vdW50"),
						Value: mustDecodeBase64("NDI3NzMxdXRpYQ=="),
						Index: true,
					},
				},
			},
			{
				Type: "message",
				Attributes: []types.EventAttribute{
					{
						Key:   mustDecodeBase64("c2VuZGVy"),
						Value: mustDecodeBase64("Y2VsZXN0aWExanY2NXMzZ3JxZjZ2NmpsM2RwNHQ2Yzl0OXJrOTljZDhrNDR2bmo="),
						Index: true,
					},
				},
			},
			{
				Type: "withdraw_rewards",
				Attributes: []types.EventAttribute{
					{
						Key:   mustDecodeBase64("YW1vdW50"),
						Value: mustDecodeBase64("NDI3NzMxdXRpYQ=="),
						Index: true,
					},
					{
						Key:   mustDecodeBase64("dmFsaWRhdG9y"),
						Value: mustDecodeBase64("Y2VsZXN0aWF2YWxvcGVyMXY1aHJxbHY4ZHFnenZ5MHB3enF6ZzBneHk4OTlybTRrbHp4bTA3"),
						Index: true,
					},
					{
						Key:   mustDecodeBase64("ZGVsZWdhdG9y"),
						Value: mustDecodeBase64("Y2VsZXN0aWExMjUzdmRsZGxmeGx3eXBuZGgzZnp6cXQ5ZG4wcjV3MGRldHY1dWg="),
						Index: true,
					},
				},
			},
			{
				Type: "coin_spent",
				Attributes: []types.EventAttribute{
					{
						Key:   mustDecodeBase64("c3BlbmRlcg=="),
						Value: mustDecodeBase64("Y2VsZXN0aWExanY2NXMzZ3JxZjZ2NmpsM2RwNHQ2Yzl0OXJrOTljZDhrNDR2bmo="),
						Index: true,
					},
					{
						Key:   mustDecodeBase64("YW1vdW50"),
						Value: mustDecodeBase64("NHV0aWE="),
						Index: true,
					},
				},
			},
			{
				Type: "coin_received",
				Attributes: []types.EventAttribute{
					{
						Key:   mustDecodeBase64("cmVjZWl2ZXI="),
						Value: mustDecodeBase64("Y2VsZXN0aWExMjUzdmRsZGxmeGx3eXBuZGgzZnp6cXQ5ZG4wcjV3MGRldHY1dWg="),
						Index: true,
					},
					{
						Key:   mustDecodeBase64("YW1vdW50"),
						Value: mustDecodeBase64("NHV0aWE="),
						Index: true,
					},
				},
			},
			{
				Type: "transfer",
				Attributes: []types.EventAttribute{
					{
						Key:   mustDecodeBase64("cmVjaXBpZW50"),
						Value: mustDecodeBase64("Y2VsZXN0aWExMjUzdmRsZGxmeGx3eXBuZGgzZnp6cXQ5ZG4wcjV3MGRldHY1dWg="),
						Index: true,
					},
					{
						Key:   mustDecodeBase64("c2VuZGVy"),
						Value: mustDecodeBase64("Y2VsZXN0aWExanY2NXMzZ3JxZjZ2NmpsM2RwNHQ2Yzl0OXJrOTljZDhrNDR2bmo="),
						Index: true,
					},
					{
						Key:   mustDecodeBase64("YW1vdW50"),
						Value: mustDecodeBase64("NHV0aWE="),
						Index: true,
					},
				},
			},
			{
				Type: "message",
				Attributes: []types.EventAttribute{
					{
						Key:   mustDecodeBase64("c2VuZGVy"),
						Value: mustDecodeBase64("Y2VsZXN0aWExanY2NXMzZ3JxZjZ2NmpsM2RwNHQ2Yzl0OXJrOTljZDhrNDR2bmo="),
						Index: true,
					},
				},
			},
			{
				Type: "withdraw_rewards",
				Attributes: []types.EventAttribute{
					{
						Key:   mustDecodeBase64("YW1vdW50"),
						Value: mustDecodeBase64("NHV0aWE="),
						Index: true,
					},
					{
						Key:   mustDecodeBase64("dmFsaWRhdG9y"),
						Value: mustDecodeBase64("Y2VsZXN0aWF2YWxvcGVyMXU4MjVzcmxkaGV2N3Q0d25kM2hwbGhycGhhaGpmazdmZjN3ZmRy"),
						Index: true,
					},
					{
						Key:   mustDecodeBase64("ZGVsZWdhdG9y"),
						Value: mustDecodeBase64("Y2VsZXN0aWExMjUzdmRsZGxmeGx3eXBuZGgzZnp6cXQ5ZG4wcjV3MGRldHY1dWg="),
						Index: true,
					},
				},
			},
			{
				Type: "redelegate",
				Attributes: []types.EventAttribute{
					{
						Key:   mustDecodeBase64("c291cmNlX3ZhbGlkYXRvcg=="),
						Value: mustDecodeBase64("Y2VsZXN0aWF2YWxvcGVyMXY1aHJxbHY4ZHFnenZ5MHB3enF6ZzBneHk4OTlybTRrbHp4bTA3"),
						Index: true,
					},
					{
						Key:   mustDecodeBase64("ZGVzdGluYXRpb25fdmFsaWRhdG9y"),
						Value: mustDecodeBase64("Y2VsZXN0aWF2YWxvcGVyMXU4MjVzcmxkaGV2N3Q0d25kM2hwbGhycGhhaGpmazdmZjN3ZmRy"),
						Index: true,
					},
					{
						Key:   mustDecodeBase64("YW1vdW50"),
						Value: mustDecodeBase64("MjYwMjAwMDB1dGlh"),
						Index: true,
					},
					{
						Key:   mustDecodeBase64("Y29tcGxldGlvbl90aW1l"),
						Value: mustDecodeBase64("MjAyNC0wMy0xN1QyMjoyMjoyM1o="),
						Index: true,
					},
				},
			},
			{
				Type: "message",
				Attributes: []types.EventAttribute{
					{
						Key:   mustDecodeBase64("bW9kdWxl"),
						Value: mustDecodeBase64("c3Rha2luZw=="),
						Index: true,
					},
					{
						Key:   mustDecodeBase64("c2VuZGVy"),
						Value: mustDecodeBase64("Y2VsZXN0aWExMjUzdmRsZGxmeGx3eXBuZGgzZnp6cXQ5ZG4wcjV3MGRldHY1dWg="),
						Index: true,
					},
				},
			},
		},
	}
	block, now := testsuite.CreateTestBlock(txRes, 3)

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height:       1000,
		Time:         now,
		MessageTypes: storageTypes.NewMsgTypeBitMask(),
	}
	resultTxs, err := parseTxs(decodeCtx, block)

	assert.NoError(t, err)
	assert.Len(t, resultTxs, 3)

	f := resultTxs[0]
	assert.Equal(t, now, f.Time)
	assert.Equal(t, storageTypes.StatusSuccess, f.Status)
	assert.Equal(t, "", f.Error)
	assert.Equal(t, int64(12000), f.GasWanted)
	assert.Equal(t, int64(1000), f.GasUsed)
	assert.Equal(t, "celestia-explorer", f.Codespace)
}

func TestParseTxs_FailedTx(t *testing.T) {
	txRes := types.ResponseDeliverTx{
		Code:      1,
		Data:      []byte{},
		Log:       "something weird happened",
		Info:      "info",
		GasWanted: 12000,
		GasUsed:   1000,
		Events:    nil,
		Codespace: "celestia-explorer",
	}
	block, now := testsuite.CreateTestBlock(txRes, 1)

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height:       1000,
		Time:         now,
		MessageTypes: storageTypes.NewMsgTypeBitMask(),
	}
	resultTxs, err := parseTxs(decodeCtx, block)

	assert.NoError(t, err)
	assert.Len(t, resultTxs, 1)

	f := resultTxs[0]
	assert.Equal(t, now, f.Time)
	assert.Equal(t, storageTypes.StatusFailed, f.Status)
	assert.Equal(t, "something weird happened", f.Error)
	assert.Equal(t, int64(12000), f.GasWanted)
	assert.Equal(t, int64(1000), f.GasUsed)
	assert.Equal(t, "celestia-explorer", f.Codespace)
}

func TestParseTxs_FailedTxWithNonstandardErrorCode(t *testing.T) {
	txRes := types.ResponseDeliverTx{
		Code:      300,
		Data:      []byte{},
		Log:       "something unusual happened",
		Info:      "info",
		GasWanted: 12000,
		GasUsed:   1000,
		Events:    nil,
		Codespace: "celestia-explorer",
	}
	block, now := testsuite.CreateTestBlock(txRes, 1)

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height:       1000,
		Time:         now,
		MessageTypes: storageTypes.NewMsgTypeBitMask(),
	}
	resultTxs, err := parseTxs(decodeCtx, block)

	assert.NoError(t, err)
	assert.Len(t, resultTxs, 1)

	f := resultTxs[0]
	assert.Equal(t, now, f.Time)
	assert.Equal(t, storageTypes.StatusFailed, f.Status)
	assert.Equal(t, "something unusual happened", f.Error)
	assert.Equal(t, int64(12000), f.GasWanted)
	assert.Equal(t, int64(1000), f.GasUsed)
	assert.Equal(t, "celestia-explorer", f.Codespace)
}
