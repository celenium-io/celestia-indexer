// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rpc

import (
	"encoding/base64"
	"encoding/hex"
	stdjson "encoding/json"
	"testing"
	"time"

	jxpkg "github.com/go-faster/jx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	nodeTypes "github.com/celenium-io/celestia-indexer/pkg/node/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	cmtTypes "github.com/cometbft/cometbft/types"
)

// ── test helper ──────────────────────────────────────────────────────────────

// jdec creates a fresh jx decoder over a JSON string literal.
// Uses the built-in pool so resources are reused.
func jdec(s string) *jxpkg.Decoder {
	d := jxpkg.GetDecoder()
	d.ResetBytes([]byte(s))
	return d
}

func mustDecodeHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	require.NoError(t, err)
	return b
}

func mustDecodeB64(t *testing.T, s string) []byte {
	t.Helper()
	b, err := base64.StdEncoding.DecodeString(s)
	require.NoError(t, err)
	return b
}

// ── jxInternStr ──────────────────────────────────────────────────────────────

func TestJxInternStr_Known(t *testing.T) {
	for _, s := range []string{"amount", "sender", "transfer", "celestia.blob.v1.EventPayForBlobs", "message_id"} {
		t.Run(s, func(t *testing.T) {
			d := jdec(`"` + s + `"`)
			defer jxpkg.PutDecoder(d)

			got, err := jxInternStr(d)
			require.NoError(t, err)
			require.Equal(t, s, got)
		})
	}
}

func TestJxInternStr_Unknown(t *testing.T) {
	d := jdec(`"definitely_not_in_the_table_xyz"`)
	defer jxpkg.PutDecoder(d)

	got, err := jxInternStr(d)
	require.NoError(t, err)
	require.Equal(t, "definitely_not_in_the_table_xyz", got)
}

func TestJxInternStr_Empty(t *testing.T) {
	d := jdec(`""`)
	defer jxpkg.PutDecoder(d)

	got, err := jxInternStr(d)
	require.NoError(t, err)
	require.Equal(t, "", got)
}

// ── jxHex ────────────────────────────────────────────────────────────────────

func TestJxHex_Valid(t *testing.T) {
	d := jdec(`"DEADBEEF"`)
	defer jxpkg.PutDecoder(d)

	got, err := jxHex(d)
	require.NoError(t, err)
	require.Equal(t, []byte{0xDE, 0xAD, 0xBE, 0xEF}, got)
}

func TestJxHex_Empty(t *testing.T) {
	d := jdec(`""`)
	defer jxpkg.PutDecoder(d)

	got, err := jxHex(d)
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestJxHex_Invalid(t *testing.T) {
	d := jdec(`"ZZZZ"`)
	defer jxpkg.PutDecoder(d)

	_, err := jxHex(d)
	require.Error(t, err)
}

func TestJxHex_LongerHash(t *testing.T) {
	d := jdec(`"ABCDEF012345ABCDEF012345ABCDEF012345ABCD"`)
	defer jxpkg.PutDecoder(d)

	got, err := jxHex(d)
	require.NoError(t, err)
	require.Equal(t, mustDecodeHex(t, "ABCDEF012345ABCDEF012345ABCDEF012345ABCD"), got)
}

// ── jxInt64 ──────────────────────────────────────────────────────────────────

func TestJxInt64(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`"0"`, 0},
		{`"1234567"`, 1234567},
		{`"-99"`, -99},
		{`"9223372036854775807"`, 9223372036854775807}, // math.MaxInt64
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			d := jdec(tt.input)
			defer jxpkg.PutDecoder(d)

			got, err := jxInt64(d)
			require.NoError(t, err)
			require.Equal(t, tt.expected, got)
		})
	}
}

func TestJxInt64_Invalid(t *testing.T) {
	d := jdec(`"not_a_number"`)
	defer jxpkg.PutDecoder(d)

	_, err := jxInt64(d)
	require.Error(t, err)
}

// ── jxUint64 ─────────────────────────────────────────────────────────────────

func TestJxUint64(t *testing.T) {
	tests := []struct {
		input    string
		expected uint64
	}{
		{`"0"`, 0},
		{`"128"`, 128},
		{`"18446744073709551615"`, 18446744073709551615}, // math.MaxUint64
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			d := jdec(tt.input)
			defer jxpkg.PutDecoder(d)

			got, err := jxUint64(d)
			require.NoError(t, err)
			require.Equal(t, tt.expected, got)
		})
	}
}

// ── jxDuration ───────────────────────────────────────────────────────────────

func TestJxDuration(t *testing.T) {
	// 172800000000000 ns = 48 hours
	d := jdec(`"172800000000000"`)
	defer jxpkg.PutDecoder(d)

	got, err := jxDuration(d)
	require.NoError(t, err)
	require.Equal(t, 48*time.Hour, got)
}

func TestJxDuration_Zero(t *testing.T) {
	d := jdec(`"0"`)
	defer jxpkg.PutDecoder(d)

	got, err := jxDuration(d)
	require.NoError(t, err)
	require.Equal(t, time.Duration(0), got)
}

// ── jxTime ───────────────────────────────────────────────────────────────────

func TestJxTime(t *testing.T) {
	const ts = "2024-06-15T12:34:56.789000000Z"
	d := jdec(`"` + ts + `"`)
	defer jxpkg.PutDecoder(d)

	got, err := jxTime(d)
	require.NoError(t, err)
	want, _ := time.Parse(time.RFC3339Nano, ts)
	require.True(t, got.Equal(want))
}

func TestJxTime_WithOffset(t *testing.T) {
	const ts = "2024-01-01T00:00:00+03:00"
	d := jdec(`"` + ts + `"`)
	defer jxpkg.PutDecoder(d)

	got, err := jxTime(d)
	require.NoError(t, err)
	want, _ := time.Parse(time.RFC3339Nano, ts)
	require.True(t, got.Equal(want))
}

func TestJxTime_Invalid(t *testing.T) {
	d := jdec(`"not-a-time"`)
	defer jxpkg.PutDecoder(d)

	_, err := jxTime(d)
	require.Error(t, err)
}

// ── jxEventAttribute ─────────────────────────────────────────────────────────

func TestJxEventAttribute_AllFields(t *testing.T) {
	d := jdec(`{"key":"amount","value":"1000utia"}`)
	defer jxpkg.PutDecoder(d)

	got, err := jxEventAttribute(d)
	require.NoError(t, err)
	require.Equal(t, "amount", got.Key)
	require.Equal(t, "1000utia", got.Value)
}

func TestJxEventAttribute_KnownKeyInterned(t *testing.T) {
	d := jdec(`{"key":"sender","value":"celestia1abc"}`)
	defer jxpkg.PutDecoder(d)

	got, err := jxEventAttribute(d)
	require.NoError(t, err)
	require.Equal(t, "sender", got.Key)
	require.Equal(t, "celestia1abc", got.Value)
}

func TestJxEventAttribute_UnknownFieldsSkipped(t *testing.T) {
	d := jdec(`{"key":"fee","value":"100utia","index":true,"extra":42}`)
	defer jxpkg.PutDecoder(d)

	got, err := jxEventAttribute(d)
	require.NoError(t, err)
	require.Equal(t, "fee", got.Key)
	require.Equal(t, "100utia", got.Value)
}

func TestJxEventAttribute_OnlyKey(t *testing.T) {
	d := jdec(`{"key":"validator"}`)
	defer jxpkg.PutDecoder(d)

	got, err := jxEventAttribute(d)
	require.NoError(t, err)
	require.Equal(t, "validator", got.Key)
	require.Equal(t, "", got.Value)
}

// ── jxEvent ──────────────────────────────────────────────────────────────────

func TestJxEvent_WithAttributes(t *testing.T) {
	d := jdec(`{
		"type": "transfer",
		"attributes": [
			{"key": "sender",   "value": "celestia1aaa"},
			{"key": "receiver", "value": "celestia1bbb"},
			{"key": "amount",   "value": "500utia"}
		]
	}`)
	defer jxpkg.PutDecoder(d)

	got, err := jxEvent(d)
	require.NoError(t, err)
	require.Equal(t, "transfer", got.Type)
	require.Len(t, got.Attributes, 3)
	require.Equal(t, "sender", got.Attributes[0].Key)
	require.Equal(t, "celestia1aaa", got.Attributes[0].Value)
	require.Equal(t, "receiver", got.Attributes[1].Key)
	require.Equal(t, "amount", got.Attributes[2].Key)
	require.Equal(t, "500utia", got.Attributes[2].Value)
}

func TestJxEvent_EmptyAttributes(t *testing.T) {
	d := jdec(`{"type":"message","attributes":[]}`)
	defer jxpkg.PutDecoder(d)

	got, err := jxEvent(d)
	require.NoError(t, err)
	require.Equal(t, "message", got.Type)
	require.Empty(t, got.Attributes)
}

func TestJxEvent_UnknownFieldsSkipped(t *testing.T) {
	d := jdec(`{"type":"tx","unknown_field":"ignore","attributes":[]}`)
	defer jxpkg.PutDecoder(d)

	got, err := jxEvent(d)
	require.NoError(t, err)
	require.Equal(t, "tx", got.Type)
}

// ── jxTxResult ───────────────────────────────────────────────────────────────

func TestJxTxResult_AllFields(t *testing.T) {
	// log is a JSON string (the most common case in CometBFT responses).
	// RawAppend returns the raw JSON token including the surrounding quotes.
	const input = `{
		"code": 0,
		"log": "tx executed successfully",
		"gas_wanted": "200000",
		"gas_used": "82000",
		"events": [
			{
				"type": "tx",
				"attributes": [
					{"key": "fee",    "value": "500utia"},
					{"key": "fee_payer", "value": "celestia1abc"}
				]
			}
		],
		"codespace": ""
	}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	got, err := jxTxResult(d)
	require.NoError(t, err)
	require.Equal(t, uint32(0), got.Code)
	// RawAppend stores the raw JSON token; for a JSON string this includes the quotes.
	require.Equal(t, stdjson.RawMessage(`"tx executed successfully"`), got.Log)
	require.Equal(t, int64(200000), got.GasWanted)
	require.Equal(t, int64(82000), got.GasUsed)
	require.Len(t, got.Events, 1)
	require.Equal(t, "tx", got.Events[0].Type)
	require.Len(t, got.Events[0].Attributes, 2)
	require.Equal(t, "", got.Codespace)
}

func TestJxTxResult_LogAsObject(t *testing.T) {
	// log can also be a JSON object in some ABCI responses.
	d := jdec(`{"code":0,"log":{"msg":"ok","data":42},"gas_wanted":"100","gas_used":"50","events":[]}`)
	defer jxpkg.PutDecoder(d)

	got, err := jxTxResult(d)
	require.NoError(t, err)
	require.Equal(t, stdjson.RawMessage(`{"msg":"ok","data":42}`), got.Log)
}

func TestJxTxResult_NonZeroCode(t *testing.T) {
	d := jdec(`{"code":12,"codespace":"sdk","log":"\"tx failed\"","gas_wanted":"100","gas_used":"50","events":[]}`)
	defer jxpkg.PutDecoder(d)

	got, err := jxTxResult(d)
	require.NoError(t, err)
	require.Equal(t, uint32(12), got.Code)
	require.Equal(t, "sdk", got.Codespace)
	require.True(t, got.IsFailed())
}

func TestJxTxResult_UnknownFieldsSkipped(t *testing.T) {
	d := jdec(`{"code":0,"gas_wanted":"100","gas_used":"50","events":[],"new_field":"ignored"}`)
	defer jxpkg.PutDecoder(d)

	_, err := jxTxResult(d)
	require.NoError(t, err)
}

// ── jxConsensusParamsBlock ───────────────────────────────────────────────────

func TestJxConsensusParamsBlock(t *testing.T) {
	d := jdec(`{"max_bytes":"22020096","max_gas":"-1"}`)
	defer jxpkg.PutDecoder(d)

	got, err := jxConsensusParamsBlock(d)
	require.NoError(t, err)
	require.Equal(t, int64(22020096), got.MaxBytes)
	require.Equal(t, int64(-1), got.MaxGas)
}

func TestJxConsensusParamsBlock_UnknownFieldsSkipped(t *testing.T) {
	d := jdec(`{"max_bytes":"1024","max_gas":"0","future_param":"42"}`)
	defer jxpkg.PutDecoder(d)

	got, err := jxConsensusParamsBlock(d)
	require.NoError(t, err)
	require.Equal(t, int64(1024), got.MaxBytes)
}

// ── jxConsensusParamsEvidence ────────────────────────────────────────────────

func TestJxConsensusParamsEvidence(t *testing.T) {
	// 172800000000000 ns = 48 hours
	d := jdec(`{"max_age_num_blocks":"100000","max_age_duration":"172800000000000","max_bytes":"1048576"}`)
	defer jxpkg.PutDecoder(d)

	got, err := jxConsensusParamsEvidence(d)
	require.NoError(t, err)
	require.Equal(t, int64(100000), got.MaxAgeNumBlocks)
	require.Equal(t, 48*time.Hour, got.MaxAgeDuration)
	require.Equal(t, int64(1048576), got.MaxBytes)
}

// ── jxConsensusParamsValidator ───────────────────────────────────────────────

func TestJxConsensusParamsValidator(t *testing.T) {
	d := jdec(`{"pub_key_types":["ed25519","secp256k1"]}`)
	defer jxpkg.PutDecoder(d)

	got, err := jxConsensusParamsValidator(d)
	require.NoError(t, err)
	require.Equal(t, []string{"ed25519", "secp256k1"}, got.PubKeyTypes)
}

func TestJxConsensusParamsValidator_Empty(t *testing.T) {
	d := jdec(`{"pub_key_types":[]}`)
	defer jxpkg.PutDecoder(d)

	got, err := jxConsensusParamsValidator(d)
	require.NoError(t, err)
	require.Empty(t, got.PubKeyTypes)
}

// ── jxConsensusParams ────────────────────────────────────────────────────────

func TestJxConsensusParams_Full(t *testing.T) {
	const input = `{
		"block":     {"max_bytes":"22020096","max_gas":"-1"},
		"evidence":  {"max_age_num_blocks":"100000","max_age_duration":"172800000000000","max_bytes":"1048576"},
		"validator": {"pub_key_types":["ed25519"]}
	}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	got, err := jxConsensusParams(d)
	require.NoError(t, err)
	require.NotNil(t, got.Block)
	require.Equal(t, int64(22020096), got.Block.MaxBytes)
	require.NotNil(t, got.Evidence)
	require.Equal(t, int64(100000), got.Evidence.MaxAgeNumBlocks)
	require.NotNil(t, got.Validator)
	require.Equal(t, []string{"ed25519"}, got.Validator.PubKeyTypes)
}

func TestJxConsensusParams_NullSections(t *testing.T) {
	d := jdec(`{"block":null,"evidence":null,"validator":null}`)
	defer jxpkg.PutDecoder(d)

	got, err := jxConsensusParams(d)
	require.NoError(t, err)
	require.Nil(t, got.Block)
	require.Nil(t, got.Evidence)
	require.Nil(t, got.Validator)
}

func TestJxConsensusParams_UnknownFieldsSkipped(t *testing.T) {
	d := jdec(`{"block":{"max_bytes":"512","max_gas":"0"},"abci":{}}`)
	defer jxpkg.PutDecoder(d)

	got, err := jxConsensusParams(d)
	require.NoError(t, err)
	require.NotNil(t, got.Block)
}

// ── jxResultBlockResults ─────────────────────────────────────────────────────

func TestJxResultBlockResults_NullTxsResults(t *testing.T) {
	d := jdec(`{"height":"8880520","txs_results":null,"finalize_block_events":[]}`)
	defer jxpkg.PutDecoder(d)

	got, err := jxResultBlockResults(d)
	require.NoError(t, err)
	require.Equal(t, pkgTypes.Level(8880520), got.Height)
	require.Nil(t, got.TxsResults)
	require.Empty(t, got.FinalizeBlockEvents)
}

func TestJxResultBlockResults_WithTxsResults(t *testing.T) {
	const input = `{
		"height": "1234",
		"txs_results": [
			{"code":0,"log":"\"ok\"","gas_wanted":"100000","gas_used":"40000","events":[]},
			{"code":5,"log":"\"fail\"","gas_wanted":"50000","gas_used":"50000","events":[],"codespace":"sdk"}
		],
		"finalize_block_events": [
			{"type":"coin_received","attributes":[{"key":"receiver","value":"addr1"},{"key":"amount","value":"10uatom"}]}
		]
	}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	got, err := jxResultBlockResults(d)
	require.NoError(t, err)
	require.Equal(t, pkgTypes.Level(1234), got.Height)
	require.Len(t, got.TxsResults, 2)
	require.Equal(t, uint32(0), got.TxsResults[0].Code)
	require.Equal(t, uint32(5), got.TxsResults[1].Code)
	require.Equal(t, "sdk", got.TxsResults[1].Codespace)
	require.Len(t, got.FinalizeBlockEvents, 1)
	require.Equal(t, "coin_received", got.FinalizeBlockEvents[0].Type)
}

func TestJxResultBlockResults_WithConsensusParams(t *testing.T) {
	const input = `{
		"height": "999",
		"txs_results": null,
		"finalize_block_events": [],
		"consensus_param_updates": {
			"block":     {"max_bytes":"22020096","max_gas":"-1"},
			"validator": {"pub_key_types":["ed25519"]}
		}
	}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	got, err := jxResultBlockResults(d)
	require.NoError(t, err)
	require.NotNil(t, got.ConsensusParamUpdates)
	require.NotNil(t, got.ConsensusParamUpdates.Block)
	require.Equal(t, int64(22020096), got.ConsensusParamUpdates.Block.MaxBytes)
}

func TestJxResultBlockResults_NullConsensusParams(t *testing.T) {
	d := jdec(`{"height":"1","txs_results":null,"finalize_block_events":[],"consensus_param_updates":null}`)
	defer jxpkg.PutDecoder(d)

	got, err := jxResultBlockResults(d)
	require.NoError(t, err)
	require.Nil(t, got.ConsensusParamUpdates)
}

// ── jxHeader ─────────────────────────────────────────────────────────────────

const testHashHex = "ABCDEF1234567890ABCDEF1234567890ABCDEF12"

func testHashBytes(t *testing.T) []byte {
	t.Helper()
	return mustDecodeHex(t, testHashHex)
}

func TestJxHeader_AllFields(t *testing.T) {
	input := `{
		"chain_id": "celestia",
		"height":   "1234567",
		"time":     "2024-06-15T12:34:56.000000000Z",
		"last_block_id": {"hash": "` + testHashHex + `","parts":{"total":1,"hash":""}},
		"last_commit_hash":   "` + testHashHex + `",
		"data_hash":          "` + testHashHex + `",
		"validators_hash":    "` + testHashHex + `",
		"next_validators_hash":"` + testHashHex + `",
		"consensus_hash":     "` + testHashHex + `",
		"app_hash":           "` + testHashHex + `",
		"last_results_hash":  "` + testHashHex + `",
		"evidence_hash":      "` + testHashHex + `",
		"proposer_address":   "` + testHashHex + `",
		"version": {"block":"11","app":"2"}
	}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	var h pkgTypes.Header
	err := jxHeader(d, &h)
	require.NoError(t, err)

	wantHash := testHashBytes(t)
	wantTime, _ := time.Parse(time.RFC3339Nano, "2024-06-15T12:34:56.000000000Z")

	require.Equal(t, "celestia", h.ChainID)
	require.Equal(t, int64(1234567), h.Height)
	require.True(t, h.Time.Equal(wantTime))
	require.Equal(t, wantHash, []byte(h.LastBlockID.Hash))
	require.Equal(t, wantHash, []byte(h.LastCommitHash))
	require.Equal(t, wantHash, []byte(h.DataHash))
	require.Equal(t, wantHash, []byte(h.ValidatorsHash))
	require.Equal(t, wantHash, []byte(h.NextValidatorsHash))
	require.Equal(t, wantHash, []byte(h.ConsensusHash))
	require.Equal(t, wantHash, []byte(h.AppHash))
	require.Equal(t, wantHash, []byte(h.LastResultsHash))
	require.Equal(t, wantHash, []byte(h.EvidenceHash))
	require.Equal(t, wantHash, []byte(h.ProposerAddress))
}

func TestJxHeader_ChainIDInterned(t *testing.T) {
	// "celestia" is in knownEventStrings if we added it;
	// either way the value must be correct.
	d := jdec(`{"chain_id":"celestia","height":"1","time":"2024-01-01T00:00:00Z"}`)
	defer jxpkg.PutDecoder(d)

	var h pkgTypes.Header
	err := jxHeader(d, &h)
	require.NoError(t, err)
	require.Equal(t, "celestia", h.ChainID)
}

func TestJxHeader_EmptyHashes(t *testing.T) {
	d := jdec(`{"chain_id":"test","height":"1","time":"2024-01-01T00:00:00Z","last_block_id":{"hash":""}}`)
	defer jxpkg.PutDecoder(d)

	var h pkgTypes.Header
	err := jxHeader(d, &h)
	require.NoError(t, err)
	require.Nil(t, []byte(h.LastBlockID.Hash))
}

// ── jxData ───────────────────────────────────────────────────────────────────

func TestJxData_WithTxs(t *testing.T) {
	// AQIDBA== = base64({1,2,3,4});  BQYHCAk= = base64({5,6,7,8,9})
	d := jdec(`{"txs":["AQIDBA==","BQYHCAk="],"square_size":"128"}`)
	defer jxpkg.PutDecoder(d)

	var data pkgTypes.Data
	err := jxData(d, &data)
	require.NoError(t, err)
	require.Equal(t, uint64(128), data.SquareSize)
	require.Len(t, data.Txs, 2)
	require.Equal(t, mustDecodeB64(t, "AQIDBA=="), data.Txs[0])
	require.Equal(t, mustDecodeB64(t, "BQYHCAk="), data.Txs[1])
}

func TestJxData_EmptyTxs(t *testing.T) {
	d := jdec(`{"txs":[],"square_size":"64"}`)
	defer jxpkg.PutDecoder(d)

	var data pkgTypes.Data
	err := jxData(d, &data)
	require.NoError(t, err)
	require.Equal(t, uint64(64), data.SquareSize)
	require.Empty(t, data.Txs)
}

func TestJxData_UnknownFieldsSkipped(t *testing.T) {
	d := jdec(`{"txs":[],"square_size":"32","hash":"ABCD"}`)
	defer jxpkg.PutDecoder(d)

	var data pkgTypes.Data
	err := jxData(d, &data)
	require.NoError(t, err)
}

func TestJxData_InvalidBase64(t *testing.T) {
	d := jdec(`{"txs":["not-valid-base64!!!"],"square_size":"1"}`)
	defer jxpkg.PutDecoder(d)

	var data pkgTypes.Data
	err := jxData(d, &data)
	require.Error(t, err)
}

// ── jxCommitSig ──────────────────────────────────────────────────────────────

func TestJxCommitSig_AllFields(t *testing.T) {
	input := `{
		"block_id_flag": 2,
		"validator_address": "` + testHashHex + `",
		"timestamp": "2024-06-15T12:34:56.000000000Z",
		"signature": "abc123ignored=="
	}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	got, err := jxCommitSig(d)
	require.NoError(t, err)
	require.Equal(t, cmtTypes.BlockIDFlag(2), got.BlockIDFlag)
	require.Equal(t, testHashBytes(t), []byte(got.ValidatorAddress))
	wantTime, _ := time.Parse(time.RFC3339Nano, "2024-06-15T12:34:56.000000000Z")
	require.True(t, got.Timestamp.Equal(wantTime))
}

func TestJxCommitSig_AbsentValidator(t *testing.T) {
	// BlockIDFlagAbsent = 1, validator_address empty
	d := jdec(`{"block_id_flag":1,"validator_address":"","timestamp":"2024-01-01T00:00:00Z"}`)
	defer jxpkg.PutDecoder(d)

	got, err := jxCommitSig(d)
	require.NoError(t, err)
	require.Equal(t, cmtTypes.BlockIDFlag(1), got.BlockIDFlag)
	require.Nil(t, []byte(got.ValidatorAddress))
}

// ── jxCommit ─────────────────────────────────────────────────────────────────

func TestJxCommit_WithSignatures(t *testing.T) {
	input := `{
		"height": "999",
		"round": 0,
		"block_id": {"hash":"` + testHashHex + `"},
		"signatures": [
			{"block_id_flag":2,"validator_address":"` + testHashHex + `","timestamp":"2024-01-01T00:00:00Z"},
			{"block_id_flag":1,"validator_address":"","timestamp":"2024-01-01T00:00:00Z"}
		]
	}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	got, err := jxCommit(d)
	require.NoError(t, err)
	require.Equal(t, int64(999), got.Height)
	require.Len(t, got.Signatures, 2)
	require.Equal(t, cmtTypes.BlockIDFlag(2), got.Signatures[0].BlockIDFlag)
	require.Equal(t, cmtTypes.BlockIDFlag(1), got.Signatures[1].BlockIDFlag)
}

func TestJxCommit_EmptySignatures(t *testing.T) {
	d := jdec(`{"height":"1","signatures":[]}`)
	defer jxpkg.PutDecoder(d)

	got, err := jxCommit(d)
	require.NoError(t, err)
	require.Equal(t, int64(1), got.Height)
	require.Empty(t, got.Signatures)
}

// ── jxResultBlock ─────────────────────────────────────────────────────────────

func TestJxResultBlock_Full(t *testing.T) {
	input := `{
		"block_id": {"hash":"` + testHashHex + `","parts":{"total":8,"hash":""}},
		"block": {
			"header": {
				"chain_id": "celestia",
				"height":   "42",
				"time":     "2024-06-15T00:00:00Z",
				"last_block_id": {"hash":"` + testHashHex + `"},
				"last_commit_hash":    "` + testHashHex + `",
				"data_hash":           "` + testHashHex + `",
				"validators_hash":     "` + testHashHex + `",
				"next_validators_hash":"` + testHashHex + `",
				"consensus_hash":      "` + testHashHex + `",
				"app_hash":            "` + testHashHex + `",
				"last_results_hash":   "` + testHashHex + `",
				"evidence_hash":       "` + testHashHex + `",
				"proposer_address":    "` + testHashHex + `"
			},
			"data": {
				"txs": ["AQIDBA=="],
				"square_size": "32"
			},
			"last_commit": {
				"height": "41",
				"signatures": [
					{"block_id_flag":2,"validator_address":"` + testHashHex + `","timestamp":"2024-06-15T00:00:00Z"}
				]
			},
			"evidence": {}
		}
	}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	got, err := jxResultBlock(d)
	require.NoError(t, err)

	wantHash := testHashBytes(t)
	require.Equal(t, wantHash, []byte(got.BlockID.Hash))
	require.NotNil(t, got.Block)
	require.Equal(t, "celestia", got.Block.ChainID)
	require.Equal(t, int64(42), got.Block.Height)
	require.Equal(t, uint64(32), got.Block.SquareSize)
	require.Len(t, got.Block.Txs, 1)
	require.Equal(t, mustDecodeB64(t, "AQIDBA=="), got.Block.Txs[0])
	require.NotNil(t, got.Block.LastCommit)
	require.Equal(t, int64(41), got.Block.LastCommit.Height)
	require.Len(t, got.Block.LastCommit.Signatures, 1)
}

func TestJxResultBlock_NoTxs(t *testing.T) {
	input := `{
		"block_id": {"hash":"` + testHashHex + `"},
		"block": {
			"header": {"chain_id":"test","height":"1","time":"2024-01-01T00:00:00Z"},
			"data":   {"txs":[],"square_size":"1"},
			"last_commit": {"height":"0","signatures":[]}
		}
	}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	got, err := jxResultBlock(d)
	require.NoError(t, err)
	require.NotNil(t, got.Block)
	require.Empty(t, got.Block.Txs)
}

// ── jxBatchResponse ───────────────────────────────────────────────────────────

const minimalBlockJSON = `{
	"block_id": {"hash":"DEADBEEF"},
	"block": {
		"header": {"chain_id":"celestia","height":"100","time":"2024-01-01T00:00:00Z"},
		"data":   {"txs":[],"square_size":"32"},
		"last_commit": {"height":"99","signatures":[]}
	}
}`

const minimalResultsJSON = `{
	"height": "100",
	"txs_results": null,
	"finalize_block_events": [
		{"type":"transfer","attributes":[{"key":"amount","value":"100utia"}]}
	]
}`

func batchJSON(block, results string) string {
	return `[{"jsonrpc":"2.0","id":-1,"result":` + block + `},` +
		`{"jsonrpc":"2.0","id":-1,"result":` + results + `}]`
}

func TestJxBatchResponse_Normal(t *testing.T) {
	batch := batchJSON(minimalBlockJSON, minimalResultsJSON)
	d := jdec(batch)
	defer jxpkg.PutDecoder(d)

	var called int
	err := jxBatchResponse(d, func(bd pkgTypes.BlockData) error {
		called++
		require.Equal(t, pkgTypes.Level(100), bd.Height)
		require.Equal(t, int64(100), bd.Block.Height)
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, 1, called, "fn should be called exactly once per block+results pair")
}

func TestJxBatchResponse_MultiplePairs(t *testing.T) {
	batch := `[` +
		`{"jsonrpc":"2.0","id":-1,"result":` + minimalBlockJSON + `},` +
		`{"jsonrpc":"2.0","id":-1,"result":` + minimalResultsJSON + `},` +
		`{"jsonrpc":"2.0","id":-1,"result":` + minimalBlockJSON + `},` +
		`{"jsonrpc":"2.0","id":-1,"result":` + minimalResultsJSON + `}` +
		`]`
	d := jdec(batch)
	defer jxpkg.PutDecoder(d)

	var called int
	err := jxBatchResponse(d, func(bd pkgTypes.BlockData) error {
		called++
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, 2, called)
}

func TestJxBatchResponse_NullError(t *testing.T) {
	// error:null is equivalent to no error — must not fail.
	batch := `[` +
		`{"jsonrpc":"2.0","id":-1,"result":` + minimalBlockJSON + `,"error":null},` +
		`{"jsonrpc":"2.0","id":-1,"result":` + minimalResultsJSON + `,"error":null}` +
		`]`
	d := jdec(batch)
	defer jxpkg.PutDecoder(d)

	var called int
	err := jxBatchResponse(d, func(bd pkgTypes.BlockData) error {
		called++
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, 1, called)
}

func TestJxBatchResponse_RPCErrorBeforeResult(t *testing.T) {
	// error field appears before result in JSON object — error must be propagated
	// when "result" key is subsequently encountered.
	batch := `[` +
		`{"jsonrpc":"2.0","id":-1,"error":{"code":-32603,"message":"internal error"},"result":` + minimalBlockJSON + `},` +
		`{"jsonrpc":"2.0","id":-1,"result":` + minimalResultsJSON + `}` +
		`]`
	d := jdec(batch)
	defer jxpkg.PutDecoder(d)

	err := jxBatchResponse(d, func(bd pkgTypes.BlockData) error {
		return nil
	})
	require.Error(t, err)
	require.ErrorIs(t, err, nodeTypes.ErrRequest)
	require.Contains(t, err.Error(), "internal error")
}

func TestJxBatchResponse_FnError(t *testing.T) {
	// fn returning an error must abort processing.
	batch := batchJSON(minimalBlockJSON, minimalResultsJSON)
	d := jdec(batch)
	defer jxpkg.PutDecoder(d)

	var fnErr = nodeTypes.ErrRequest
	err := jxBatchResponse(d, func(bd pkgTypes.BlockData) error {
		return fnErr
	})
	require.ErrorIs(t, err, nodeTypes.ErrRequest)
}

func TestJxBatchResponse_RPCErrorOnBlockResults(t *testing.T) {
	batch := `[` +
		`{"jsonrpc":"2.0","id":-1,"result":` + minimalBlockJSON + `},` +
		`{"jsonrpc":"2.0","id":-1,"error":{"code":-32603,"message":"could not find results for height #10825944"}}` +
		`]`
	d := jdec(batch)
	defer jxpkg.PutDecoder(d)

	err := jxBatchResponse(d, func(bd pkgTypes.BlockData) error {
		return nil
	})
	require.Error(t, err)
	require.ErrorIs(t, err, nodeTypes.ErrRequest)
	require.Contains(t, err.Error(), "could not find results for height #10825944")
}

// ── jxResponse ───────────────────────────────────────────────────────────────

func TestJxResponse_ResultField(t *testing.T) {
	// fn should be called with the decoder positioned at "result" value.
	input := `{"jsonrpc":"2.0","id":1,"result":{"answer":42}}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	var called bool
	err := jxResponse(d, func(d *jxpkg.Decoder) error {
		called = true
		return d.Skip()
	})
	require.NoError(t, err)
	require.True(t, called, "fn must be called when 'result' key is present")
}

func TestJxResponse_ResultDecoded(t *testing.T) {
	// fn receives a decoder positioned at the result value and can decode it.
	input := `{"jsonrpc":"2.0","id":-1,"result":{"height":"12345"}}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	var got string
	err := jxResponse(d, func(d *jxpkg.Decoder) error {
		return d.ObjBytes(func(d *jxpkg.Decoder, key []byte) error {
			if string(key) == "height" {
				var err error
				got, err = d.Str()
				return err
			}
			return d.Skip()
		})
	})
	require.NoError(t, err)
	require.Equal(t, "12345", got)
}

func TestJxResponse_ErrorFieldNonNull(t *testing.T) {
	// Non-null "error" object must return ErrRequest containing the message.
	input := `{"jsonrpc":"2.0","id":1,"error":{"code":-32600,"message":"something went wrong"}}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	err := jxResponse(d, func(d *jxpkg.Decoder) error {
		return nil
	})
	require.Error(t, err)
	require.ErrorIs(t, err, nodeTypes.ErrRequest)
	require.Contains(t, err.Error(), "something went wrong")
}

func TestJxResponse_ErrorFieldNull(t *testing.T) {
	// null "error" must be treated as no error; fn is still called when result follows.
	input := `{"jsonrpc":"2.0","id":1,"error":null,"result":42}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	var called bool
	err := jxResponse(d, func(d *jxpkg.Decoder) error {
		called = true
		return d.Skip()
	})
	require.NoError(t, err)
	require.True(t, called)
}

func TestJxResponse_UnknownFieldsSkipped(t *testing.T) {
	// "jsonrpc" and "id" fields must be silently skipped.
	input := `{"jsonrpc":"2.0","id":-1,"result":true}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	var called bool
	err := jxResponse(d, func(d *jxpkg.Decoder) error {
		called = true
		return d.Skip()
	})
	require.NoError(t, err)
	require.True(t, called)
}

func TestJxResponse_ErrorBeforeResult(t *testing.T) {
	// "error" appearing before "result" must abort; fn must NOT be called.
	input := `{"error":{"code":-1,"message":"node error"},"result":{}}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	var called bool
	err := jxResponse(d, func(d *jxpkg.Decoder) error {
		called = true
		return d.Skip()
	})
	require.Error(t, err)
	require.ErrorIs(t, err, nodeTypes.ErrRequest)
	require.False(t, called)
}

func TestJxResponse_FnError(t *testing.T) {
	// An error returned by fn must be propagated.
	input := `{"result":{}}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	sentinel := errors.New("inner error")
	err := jxResponse(d, func(d *jxpkg.Decoder) error {
		_ = d.Skip()
		return sentinel
	})
	require.ErrorIs(t, err, sentinel)
}

// ── jxStatusMinimal ──────────────────────────────────────────────────────────

func TestJxStatusMinimal_Normal(t *testing.T) {
	input := `{
		"node_info": {"protocol_version": {"p2p": "8"}},
		"sync_info": {
			"latest_block_hash": "ABCDEF",
			"latest_block_height": "987654",
			"latest_block_time": "2024-01-01T00:00:00Z",
			"catching_up": false
		},
		"validator_info": {}
	}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	level, err := jxStatusMinimal(d)
	require.NoError(t, err)
	require.Equal(t, pkgTypes.Level(987654), level)
}

func TestJxStatusMinimal_OnlySyncInfo(t *testing.T) {
	// Minimal response containing only sync_info.
	input := `{"sync_info":{"latest_block_height":"42"}}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	level, err := jxStatusMinimal(d)
	require.NoError(t, err)
	require.Equal(t, pkgTypes.Level(42), level)
}

func TestJxStatusMinimal_NoSyncInfo(t *testing.T) {
	// Missing sync_info — level stays zero, no error.
	input := `{"node_info":{},"validator_info":{}}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	level, err := jxStatusMinimal(d)
	require.NoError(t, err)
	require.Equal(t, pkgTypes.Level(0), level)
}

func TestJxStatusMinimal_ExtraFieldsInSyncInfo(t *testing.T) {
	// Fields other than latest_block_height inside sync_info are skipped.
	input := `{"sync_info":{"catching_up":false,"latest_block_height":"100","some_hash":"DEADBEEF"}}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	level, err := jxStatusMinimal(d)
	require.NoError(t, err)
	require.Equal(t, pkgTypes.Level(100), level)
}

func TestJxStatusMinimal_InvalidHeight(t *testing.T) {
	// Non-numeric height string must return a parse error.
	input := `{"sync_info":{"latest_block_height":"not_a_number"}}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	_, err := jxStatusMinimal(d)
	require.Error(t, err)
}

// ── jxGenesisChunk ───────────────────────────────────────────────────────────

func TestJxGenesisChunk_AllFields(t *testing.T) {
	payload := []byte(`{"genesis_time":"2023-01-01T00:00:00Z"}`)
	encoded := base64.StdEncoding.EncodeToString(payload)
	input := `{"chunk":"2","total":"5","data":"` + encoded + `"}`

	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	var gc GenesisChunk
	err := jxGenesisChunk(d, &gc)
	require.NoError(t, err)
	require.Equal(t, int64(2), gc.Chunk)
	require.Equal(t, int64(5), gc.Total)
	require.Equal(t, payload, gc.Data)
}

func TestJxGenesisChunk_FirstChunk(t *testing.T) {
	// chunk=0, total=1 — single-chunk genesis.
	payload := []byte(`{}`)
	encoded := base64.StdEncoding.EncodeToString(payload)
	input := `{"chunk":"0","total":"1","data":"` + encoded + `"}`

	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	var gc GenesisChunk
	err := jxGenesisChunk(d, &gc)
	require.NoError(t, err)
	require.Equal(t, int64(0), gc.Chunk)
	require.Equal(t, int64(1), gc.Total)
	require.Equal(t, payload, gc.Data)
}

func TestJxGenesisChunk_InvalidBase64(t *testing.T) {
	input := `{"chunk":"0","total":"1","data":"!!!not-valid-base64!!!"}`
	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	var gc GenesisChunk
	err := jxGenesisChunk(d, &gc)
	require.Error(t, err)
	require.Contains(t, err.Error(), "genesis chunk base64 decode")
}

func TestJxGenesisChunk_UnknownFieldsSkipped(t *testing.T) {
	payload := []byte(`[]`)
	encoded := base64.StdEncoding.EncodeToString(payload)
	input := `{"chunk":"1","total":"3","extra_field":"ignored","data":"` + encoded + `"}`

	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	var gc GenesisChunk
	err := jxGenesisChunk(d, &gc)
	require.NoError(t, err)
	require.Equal(t, int64(1), gc.Chunk)
	require.Equal(t, int64(3), gc.Total)
	require.Equal(t, payload, gc.Data)
}

func TestJxGenesisChunk_EmptyData(t *testing.T) {
	// Empty base64 string decodes to an empty (nil) byte slice.
	encoded := base64.StdEncoding.EncodeToString([]byte{})
	input := `{"chunk":"0","total":"1","data":"` + encoded + `"}`

	d := jdec(input)
	defer jxpkg.PutDecoder(d)

	var gc GenesisChunk
	err := jxGenesisChunk(d, &gc)
	require.NoError(t, err)
	require.Empty(t, gc.Data)
}
