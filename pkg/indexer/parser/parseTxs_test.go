// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package parser

import (
	"encoding/base64"
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTxs_EmptyTxsResults(t *testing.T) {
	block := types.BlockData{
		ResultBlockResults: types.ResultBlockResults{
			TxsResults: make([]*types.ResponseDeliverTx, 0),
		},
	}

	p := NewModule(config.Indexer{})
	decodeCtx := context.NewContext()
	resultTxs, err := p.parseTxs(decodeCtx, block)

	assert.NoError(t, err)
	assert.Empty(t, resultTxs)
}

func mustDecodeBase64(s string) string {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return string(data)
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
	block, now := testsuite.CreateTestBlockWithAppVersion(txRes, 3, 4)

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height:       1000,
		Time:         now,
		MessageTypes: storageTypes.NewMsgTypeBitMask(),
	}
	p := NewModule(config.Indexer{})
	resultTxs, err := p.parseTxs(decodeCtx, block)

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

	p := NewModule(config.Indexer{})
	resultTxs, err := p.parseTxs(decodeCtx, block)

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

	p := NewModule(config.Indexer{})
	resultTxs, err := p.parseTxs(decodeCtx, block)

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

func TestParseTxs_PayForBlob(t *testing.T) {
	txRes := types.ResponseDeliverTx{
		Code:      0,
		Data:      []byte{},
		Log:       "[{\"msg_index\":0,\"events\":[{\"type\":\"celestia.blob.v1.EventPayForBlobs\",\"attributes\":[{\"key\":\"blob_sizes\",\"value\":\"[2]\"},{\"key\":\"namespaces\",\"value\":\"[\\\"AAAAAAAAAAAAAAAAAAAAAAAAAEJpDCBNOWAP3dM=\\\"]\"},{\"key\":\"signer\",\"value\":\"\\\"celestia1j52ntqu7l734fjpa9lvylmtekaq0xqzhc22l0w\\\"\"}]},{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"/celestia.blob.v1.MsgPayForBlobs\"}]}]}]",
		Info:      "info",
		GasWanted: 79796,
		GasUsed:   65177,
		Events: []types.Event{
			{
				Type: "coin_spent",
				Attributes: []types.EventAttribute{
					{
						Key:   "spender",
						Value: "celestia1j52ntqu7l734fjpa9lvylmtekaq0xqzhc22l0w",
						Index: true,
					},
					{
						Key:   "amount",
						Value: "7980utia",
						Index: true,
					},
				},
			},
		},
		Codespace: "celestia-explorer",
	}
	raw, err := base64.StdEncoding.DecodeString("CoQCCp8BCpwBCiAvY2VsZXN0aWEuYmxvYi52MS5Nc2dQYXlGb3JCbG9icxJ4Ci9jZWxlc3RpYTFqNTJudHF1N2w3MzRmanBhOWx2eWxtdGVrYXEweHF6aGMyMmwwdxIdAAAAAAAAAAAAAAAAAAAAAAAAAEJpDCBNOWAP3dMaAQIiICF4PtPB1eUbDxdy5XvDx/gdk1BrBlLAYrHn5cYesAeRQgEAEh4KCBIECgIIARgBEhIKDAoEdXRpYRIENzk4MBC07wQaQOjRPPhYMdn12jdebWXpDJaDIRwmBsJ85ke8a8nwb18CLMcXzovh7/dvZm/FH1Cxe4x8NDQjY4Ethm73qhPb/pQSIgocAAAAAAAAAAAAAAAAAAAAAAAAQmkMIE05YA/d0xICZ20aBEJMT0I=")
	require.NoError(t, err)
	block, now := testsuite.CreateBlockWithTxs(txRes, raw, 1)

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height:       1000,
		Time:         now,
		MessageTypes: storageTypes.NewMsgTypeBitMask(),
	}

	p := NewModule(config.Indexer{})
	resultTxs, err := p.parseTxs(decodeCtx, block)

	require.NoError(t, err)
	require.Len(t, resultTxs, 1)

	tx := resultTxs[0]
	require.Equal(t, now, tx.Time)
	require.Equal(t, storageTypes.StatusSuccess, tx.Status)
	require.Equal(t, "", tx.Error)
	require.EqualValues(t, 79796, tx.GasWanted)
	require.EqualValues(t, 65177, tx.GasUsed)
	require.Equal(t, "celestia-explorer", tx.Codespace)
	require.Len(t, tx.Signers, 1)
}
