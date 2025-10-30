// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package decode

import (
	"encoding/base64"
	"testing"

	"cosmossdk.io/math"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	hyperlaneWarp "github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	"github.com/celenium-io/celestia-indexer/internal/currency"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	nodeTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/cosmos/cosmos-sdk/types"
	cosmosGovTypesV1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeTx_TxWithMemo(t *testing.T) {
	deliverTx := nodeTypes.ResponseDeliverTx{
		Code:      0,
		Data:      []byte{18, 45, 10, 43, 47, 99, 111, 115, 109, 111, 115, 46, 115, 116, 97, 107, 105, 110, 103, 46, 118, 49, 98, 101, 116, 97, 49, 46, 77, 115, 103, 68, 101, 108, 101, 103, 97, 116, 101, 82, 101, 115, 112, 111, 110, 115, 101},
		Log:       `[{"msg_index":0,"events":[{"type":"coin_received","attributes":[{"key":"receiver","value":"celestia1h2kqw44hdq5dwlcvsw8f2l49lkehtf9wp95kth"},{"key":"amount","value":"1562utia"}]}]}]`,
		Info:      "",
		GasWanted: 200000,
		GasUsed:   170049,
		Events:    []nodeTypes.Event{},
		Codespace: "",
	}
	txData := []byte{10, 252, 1, 10, 225, 1, 10, 42, 47, 99, 111, 115, 109, 111, 115, 46, 115, 116, 97, 107, 105, 110, 103, 46, 118, 49, 98, 101, 116, 97, 49, 46, 77, 115, 103, 66, 101, 103, 105, 110, 82, 101, 100, 101, 108, 101, 103, 97, 116, 101, 18, 178, 1, 10, 47, 99, 101, 108, 101, 115, 116, 105, 97, 49, 100, 97, 118, 122, 52, 48, 107, 97, 116, 57, 51, 116, 52, 57, 108, 106, 114, 107, 109, 107, 108, 53, 117, 113, 104, 113, 113, 52, 53, 101, 48, 116, 101, 100, 103, 102, 56, 97, 18, 54, 99, 101, 108, 101, 115, 116, 105, 97, 118, 97, 108, 111, 112, 101, 114, 49, 114, 102, 108, 117, 116, 107, 51, 101, 117, 119, 56, 100, 99, 119, 97, 101, 104, 120, 119, 117, 103, 99, 109, 57, 112, 101, 119, 107, 100, 110, 53, 54, 120, 106, 108, 104, 50, 54, 26, 54, 99, 101, 108, 101, 115, 116, 105, 97, 118, 97, 108, 111, 112, 101, 114, 49, 100, 97, 118, 122, 52, 48, 107, 97, 116, 57, 51, 116, 52, 57, 108, 106, 114, 107, 109, 107, 108, 53, 117, 113, 104, 113, 113, 52, 53, 101, 48, 116, 117, 106, 50, 115, 51, 109, 34, 15, 10, 4, 117, 116, 105, 97, 18, 7, 49, 48, 48, 48, 48, 48, 48, 18, 22, 116, 101, 115, 116, 32, 117, 105, 32, 114, 101, 100, 101, 108, 101, 103, 97, 116, 101, 32, 116, 120, 32, 18, 103, 10, 80, 10, 70, 10, 31, 47, 99, 111, 115, 109, 111, 115, 46, 99, 114, 121, 112, 116, 111, 46, 115, 101, 99, 112, 50, 53, 54, 107, 49, 46, 80, 117, 98, 75, 101, 121, 18, 35, 10, 33, 2, 205, 82, 66, 173, 172, 164, 110, 151, 162, 183, 151, 111, 80, 96, 191, 38, 188, 141, 208, 175, 86, 52, 254, 146, 134, 204, 43, 40, 79, 127, 106, 1, 18, 4, 10, 2, 8, 127, 24, 39, 18, 19, 10, 13, 10, 4, 117, 116, 105, 97, 18, 5, 55, 50, 52, 51, 49, 16, 185, 215, 17, 26, 64, 98, 225, 18, 145, 187, 225, 213, 198, 229, 6, 6, 240, 177, 0, 28, 112, 160, 126, 193, 177, 221, 161, 96, 79, 5, 192, 224, 168, 253, 161, 12, 33, 9, 118, 215, 22, 219, 239, 73, 133, 79, 37, 218, 83, 238, 115, 44, 232, 16, 163, 242, 174, 100, 175, 162, 213, 142, 194, 58, 69, 84, 81, 3, 70}
	block, _ := testsuite.CreateBlockWithTxs(deliverTx, txData, 1)

	dTx, err := Tx(block, 0)

	assert.NoError(t, err)

	assert.Equal(t, uint64(0), dTx.TimeoutHeight)
	assert.Equal(t, "test ui redelegate tx ", dTx.Memo)
	assert.Equal(t, 1, len(dTx.Messages))
	assert.Equal(t, decimal.NewFromInt(72431), dTx.Fee)
}

func TestDecodeTx_TxV050Signer(t *testing.T) {
	deliverTx := nodeTypes.ResponseDeliverTx{
		Code:      0,
		Data:      []byte{18, 45, 10, 43, 47, 99, 111, 115, 109, 111, 115, 46, 115, 116, 97, 107, 105, 110, 103, 46, 118, 49, 98, 101, 116, 97, 49, 46, 77, 115, 103, 68, 101, 108, 101, 103, 97, 116, 101, 82, 101, 115, 112, 111, 110, 115, 101},
		Log:       `[{\"msg_index\":0,\"events\":[{\"type\":\"coin_received\",\"attributes\":[{\"key\":\"receiver\",\"value\":\"celestia1q0xstyrqame6zl5puekza58jrv8629m5mne0rn\"},{\"key\":\"amount\",\"value\":\"1000000utia\"}]},{\"type\":\"coin_spent\",\"attributes\":[{\"key\":\"spender\",\"value\":\"celestia16etnwjxg6dsjuavjpr9tk822czfeylfm9f7x5g\"},{\"key\":\"amount\",\"value\":\"1000000utia\"}]},{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"/cosmos.bank.v1beta1.MsgSend\"},{\"key\":\"sender\",\"value\":\"celestia16etnwjxg6dsjuavjpr9tk822czfeylfm9f7x5g\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"celestia1q0xstyrqame6zl5puekza58jrv8629m5mne0rn\"},{\"key\":\"sender\",\"value\":\"celestia16etnwjxg6dsjuavjpr9tk822czfeylfm9f7x5g\"},{\"key\":\"amount\",\"value\":\"1000000utia\"}]}]}]`,
		Info:      "",
		GasWanted: 120000,
		GasUsed:   91289,
		Events:    []nodeTypes.Event{},
		Codespace: "",
	}

	txData, err := base64.StdEncoding.DecodeString("CrsBCqIBCiMvY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dEZWxlZ2F0ZRJ7Ci9jZWxlc3RpYTFtZW1jZjZzcWQwMGg2eXpmMGU4ZThlNzA3dGptc3pwZDNsemM2ZRI2Y2VsZXN0aWF2YWxvcGVyMXQ0Z2hmOTg0ejJ5Mnl2bjR4YWp4Y201anoyZzVzZHpxZG05eGVwGhAKBHV0aWESCDI3MDAwMDAwEhRTZW50IHZpYSBDZWxlbml1bS5pbxJkCk4KRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDZihoRCE/LwSjSlI5YpX/zHSqUy/TiAt/V+HupnPyr18SBAoCCAESEgoMCgR1dGlhEgQ0MjgyEOu5ChpAhlo8GUvT+uSwJc880k916MeM+Yl2cBpG6Nv9lmsq2Bp5npoUdSVQ1t8o/evc1EAafhljZQbF/e5geBpapwqZdA==")
	require.NoError(t, err)
	block, _ := testsuite.CreateBlockWithTxs(deliverTx, txData, 1)

	dTx, err := Tx(block, 0)

	require.NoError(t, err)

	require.Len(t, dTx.Messages, 1)
	require.Len(t, dTx.Signers, 1)
	require.Equal(t, "Sent via Celenium.io", dTx.Memo)
	require.EqualValues(t, "4282", dTx.Fee.String())
	for addr, val := range dTx.Signers {
		require.Equal(t, "celestia1memcf6sqd00h6yzf0e8e8e707tjmszpd3lzc6e", addr.String())
		require.Len(t, val, 20)
	}
}

func TestDecodeTx_TxV050Signer2(t *testing.T) {
	deliverTx := nodeTypes.ResponseDeliverTx{
		Code:      0,
		Data:      []byte{18, 45, 10, 43, 47, 99, 111, 115, 109, 111, 115, 46, 115, 116, 97, 107, 105, 110, 103, 46, 118, 49, 98, 101, 116, 97, 49, 46, 77, 115, 103, 68, 101, 108, 101, 103, 97, 116, 101, 82, 101, 115, 112, 111, 110, 115, 101},
		Log:       `[{\"msg_index\":0,\"events\":[{\"type\":\"coin_received\",\"attributes\":[{\"key\":\"receiver\",\"value\":\"celestia1q0xstyrqame6zl5puekza58jrv8629m5mne0rn\"},{\"key\":\"amount\",\"value\":\"1000000utia\"}]},{\"type\":\"coin_spent\",\"attributes\":[{\"key\":\"spender\",\"value\":\"celestia16etnwjxg6dsjuavjpr9tk822czfeylfm9f7x5g\"},{\"key\":\"amount\",\"value\":\"1000000utia\"}]},{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"/cosmos.bank.v1beta1.MsgSend\"},{\"key\":\"sender\",\"value\":\"celestia16etnwjxg6dsjuavjpr9tk822czfeylfm9f7x5g\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"celestia1q0xstyrqame6zl5puekza58jrv8629m5mne0rn\"},{\"key\":\"sender\",\"value\":\"celestia16etnwjxg6dsjuavjpr9tk822czfeylfm9f7x5g\"},{\"key\":\"amount\",\"value\":\"1000000utia\"}]}]}]`,
		Info:      "",
		GasWanted: 120000,
		GasUsed:   91289,
		Events:    []nodeTypes.Event{},
		Codespace: "",
	}

	txData, err := base64.StdEncoding.DecodeString("CpYBCpMBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEnMKL2NlbGVzdGlhMTZldG53anhnNmRzanVhdmpwcjl0azgyMmN6ZmV5bGZtOWY3eDVnEi9jZWxlc3RpYTFxMHhzdHlycWFtZTZ6bDVwdWVremE1OGpydjg2MjltNW1uZTBybhoPCgR1dGlhEgcxMDAwMDAwEh8KCBIECgIIARg1EhMKDQoEdXRpYRIFMTIwMDAQwKkHGkBGkZVTfgXdCxjUIw17iv4kvDMPT60O6r6HayzIs8XucwhgAjfdmuKwYZ20VTZuVF6SDa8UvXnrbMl8bhaImLwR")
	require.NoError(t, err)
	block, _ := testsuite.CreateBlockWithTxs(deliverTx, txData, 1)

	dTx, err := Tx(block, 0)

	require.NoError(t, err)

	require.Len(t, dTx.Messages, 1)
	require.Len(t, dTx.Signers, 1)
	for addr, val := range dTx.Signers {
		require.Equal(t, "celestia16etnwjxg6dsjuavjpr9tk822czfeylfm9f7x5g", addr.String())
		require.Len(t, val, 20)
	}
}

func TestDecodeTx_Tx_PFB(t *testing.T) {
	deliverTx := nodeTypes.ResponseDeliverTx{
		Code:      0,
		Data:      []byte{18, 45, 10, 43, 47, 99, 111, 115, 109, 111, 115, 46, 115, 116, 97, 107, 105, 110, 103, 46, 118, 49, 98, 101, 116, 97, 49, 46, 77, 115, 103, 68, 101, 108, 101, 103, 97, 116, 101, 82, 101, 115, 112, 111, 110, 115, 101},
		Log:       `[{\"msg_index\":0,\"events\":[{\"type\":\"celestia.blob.v1.EventPayForBlobs\",\"attributes\":[{\"key\":\"blob_sizes\",\"value\":\"[684]\"},{\"key\":\"namespaces\",\"value\":\"[\\\"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQ2Vyb0E=\\\"]\"},{\"key\":\"signer\",\"value\":\"\\\"celestia1rky9086t340m7rmkctuj4spxwv2gc62vlwx59v\\\"\"}]},{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"/celestia.blob.v1.MsgPayForBlobs\"}]}]}]`,
		Info:      "",
		GasWanted: 120000,
		GasUsed:   91289,
		Events:    []nodeTypes.Event{},
		Codespace: "",
	}

	txData, err := base64.StdEncoding.DecodeString("CoUCCqABCp0BCiAvY2VsZXN0aWEuYmxvYi52MS5Nc2dQYXlGb3JCbG9icxJ5Ci9jZWxlc3RpYTFya3k5MDg2dDM0MG03cm1rY3R1ajRzcHh3djJnYzYydmx3eDU5dhIdAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQ2Vyb0EaAqwFIiA92mk96XJQMA82kZz4lDP5Fbj4U7ss8LisNXzMW00q0kIBABIeCgkSBAoCCAEYjQUSEQoLCgR1dGlhEgMxODUQ+dAFGkCYFhvYyED7gTt9JbqSSJSFsQfgBcFU/H6n35PgNgZvWUp9EDMknrBVwRNwdHX00Ald9brD/Ir34FDdJAfc8p/tEs0FChwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAENlcm9BEqwFeyJnbG9iYWxTZXF1ZW5jZU51bWJlciI6MTAzOSwiYmxvY2tSYW5nZUNvdmVyZWQiOnsiYmxvY2tTdGFydCI6MzczNjgwMCwiYmxvY2tFbmQiOjM3NDA0MDB9LCJ0aW1lc3RhbXAiOjE3NDcxMzkwNDAsInJvbGx1cFNlcXVlbmNlcyI6W3sicm9sbHVwSWQiOjEsImJhdGNoZXMiOlt7InRpbWVzdGFtcCI6MTc0NzEzOTA0MCwibnVtYmVyIjoxMDM5fV19LHsicm9sbHVwSWQiOjIsImJhdGNoZXMiOlt7InRpbWVzdGFtcCI6MTc0NzEzOTA0MCwidHJhbnNhY3Rpb25zIjpbeyJ0eElkIjoiMDM5YTRmYWIwNmQ0Nzg1OTRmYWUxMmQ5NGYxN2I1ZTlmNDUyMTdiZDUzNDc3YTc5Y2I1MWY2YTQzZDY0ZDQzMSIsInJhd1RyYW5zYWN0aW9uIjoie2Zyb206c29tZXRlc3RhZGRyZXNzLCB0bzogc29tZW9uZWVsc2UsIGFtb3VudDogMC4wMDUsIHRpbWU6MTc0NzEzNTQwNyB9IiwiYmxvY2tIZWlnaHQiOjM3Mzc5MDMsInJvbGx1cElkIjoyfV0sIm51bWJlciI6MTAzOX1dfSx7InJvbGx1cElkIjozLCJiYXRjaGVzIjpbeyJ0aW1lc3RhbXAiOjE3NDcxMzkwNDAsIm51bWJlciI6MTAzOX1dfSx7InJvbGx1cElkIjo0LCJiYXRjaGVzIjpbeyJ0aW1lc3RhbXAiOjE3NDcxMzkwNDAsIm51bWJlciI6MTAzOX1dfSx7InJvbGx1cElkIjo1LCJiYXRjaGVzIjpbeyJ0aW1lc3RhbXAiOjE3NDcxMzkwNDAsIm51bWJlciI6MTAzOX1dfV19GgRCTE9C")
	require.NoError(t, err)
	block, _ := testsuite.CreateBlockWithTxs(deliverTx, txData, 1)

	dTx, err := Tx(block, 0)

	require.NoError(t, err)

	require.Len(t, dTx.Messages, 1)
	require.Len(t, dTx.Signers, 1)
	for addr, val := range dTx.Signers {
		require.Equal(t, "celestia1rky9086t340m7rmkctuj4spxwv2gc62vlwx59v", addr.String())
		require.Len(t, val, 20)
	}
}

func TestDecodeTx_Exec_signal(t *testing.T) {
	deliverTx := nodeTypes.ResponseDeliverTx{
		Code:      0,
		Data:      []byte{18, 45, 10, 43, 47, 99, 111, 115, 109, 111, 115, 46, 115, 116, 97, 107, 105, 110, 103, 46, 118, 49, 98, 101, 116, 97, 49, 46, 77, 115, 103, 68, 101, 108, 101, 103, 97, 116, 101, 82, 101, 115, 112, 111, 110, 115, 101},
		Log:       `[{\"msg_index\":0,\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"/cosmos.authz.v1beta1.MsgExec\"}]}]}]`,
		Info:      "",
		GasWanted: 210000,
		GasUsed:   68808,
		Events:    []nodeTypes.Event{},
		Codespace: "",
	}

	txData, err := base64.StdEncoding.DecodeString("CroBCrcBCh0vY29zbW9zLmF1dGh6LnYxYmV0YTEuTXNnRXhlYxKVAQovY2VsZXN0aWExazJxOGp0ZnlqMmhybm5kenNoeDZ2ZHhxc2F6bDdsbDh4bmN0ZHgSYgokL2NlbGVzdGlhLnNpZ25hbC52MS5Nc2dTaWduYWxWZXJzaW9uEjoKNmNlbGVzdGlhdmFsb3BlcjFxM3Y1Y3VnYzhjZHB1ZDg3dTR6d3kwYTc0dXhrazZ1NHE0Z3g0cBADEmcKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQIknqJ+ODItM0iT8QcFFYcYbIvmEpnFqvSWekVF9uHYlRIECgIIARgBEhMKDQoEdXRpYRIFMjEwMDAQ0OgMGkA4x6+mL9ufZBSHZrcSHtsDXIlZgyby0WmFoP3aYpZfOg+qg0MNMdnaeOf0QM2MN/EfJM8bSg37ytna3yOAZRW4")
	require.NoError(t, err)
	block, _ := testsuite.CreateBlockWithTxs(deliverTx, txData, 1)

	dTx, err := Tx(block, 0)

	require.NoError(t, err)

	require.Len(t, dTx.Messages, 1)
}

func TestDecodeTx_Tx_MsgRegisterEVMAddress(t *testing.T) {
	deliverTx := nodeTypes.ResponseDeliverTx{
		Code:      0,
		Data:      []byte{18, 45, 10, 43, 47, 99, 111, 115, 109, 111, 115, 46, 115, 116, 97, 107, 105, 110, 103, 46, 118, 49, 98, 101, 116, 97, 49, 46, 77, 115, 103, 68, 101, 108, 101, 103, 97, 116, 101, 82, 101, 115, 112, 111, 110, 115, 101},
		Log:       `[{\"msg_index\":0,\"events\":[{\"type\":\"celestia.blob.v1.EventPayForBlobs\",\"attributes\":[{\"key\":\"blob_sizes\",\"value\":\"[684]\"},{\"key\":\"namespaces\",\"value\":\"[\\\"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQ2Vyb0E=\\\"]\"},{\"key\":\"signer\",\"value\":\"\\\"celestia1rky9086t340m7rmkctuj4spxwv2gc62vlwx59v\\\"\"}]},{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"/celestia.blob.v1.MsgPayForBlobs\"}]}]}]`,
		Info:      "",
		GasWanted: 120000,
		GasUsed:   91289,
		Events:    []nodeTypes.Event{},
		Codespace: "",
	}

	txData, err := base64.StdEncoding.DecodeString("CpEBCo4BCiYvY2VsZXN0aWEucWdiLnYxLk1zZ1JlZ2lzdGVyRVZNQWRkcmVzcxJkCjZjZWxlc3RpYXZhbG9wZXIxbmF3OXVxYzl2ems0MGduNnc1OHU5ODhqeTRocWM1bGEzbnZ1bTASKjB4Mjg4NTNENjRiM0QwYzJjY0Q2NzUzM0QzNUE3QzcyMEYyYzk1RTMwZRJnClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECOaLj7A4Ho38QaJLTKtYaMqHDLALGaP36kxuOUOd5aDcSBAoCCAEYAhITCg0KBHV0aWESBTMwMDAwENDoDBpAUWpMhxhq12Kpemsh5zcDqy5Z3V/E4Avt0Dypq6e27R5fVzpUpWsHJg0IJNzy+vUlImHw5mvUU3f3nMAZh7UE7A==")
	require.NoError(t, err)
	block, _ := testsuite.CreateBlockWithTxs(deliverTx, txData, 1)

	dTx, err := Tx(block, 0)

	require.NoError(t, err)

	require.Len(t, dTx.Messages, 1)
	require.Len(t, dTx.Signers, 1)
	for addr, val := range dTx.Signers {
		require.Equal(t, "celestia1naw9uqc9vzk40gn6w58u988jy4hqc5la5vw9df", addr.String())
		require.Len(t, val, 20)
	}
}
func TestDecodeCosmosTx_DelegateMsg(t *testing.T) {
	rawTx := []byte{
		10, 164, 1, 10, 161, 1, 10, 35, 47, 99, 111, 115, 109, 111, 115, 46, 115, 116, 97, 107, 105, 110, 103, 46, 118, 49, 98, 101, 116, 97, 49, 46, 77, 115, 103, 68, 101, 108, 101, 103, 97, 116, 101, 18, 122, 10, 47, 99, 101, 108, 101, 115, 116, 105, 97, 49, 52, 122, 102, 110, 99, 50, 107, 120, 100, 103, 100, 109, 97, 99, 110, 117, 117, 121, 116, 114, 101, 53, 112, 54, 102, 120, 57, 55, 116, 116, 102, 113, 57, 101, 103, 103, 120, 100, 18, 54, 99, 101, 108, 101, 115, 116, 105, 97, 118, 97, 108, 111, 112, 101, 114, 49, 57, 117, 114, 103, 57, 97, 119, 106, 122, 119, 113, 56, 100, 52, 48, 118, 119, 106, 100, 118, 118, 48, 121, 119, 57, 107, 103, 101, 104, 115, 99, 102, 48, 122, 120, 51, 103, 115, 26, 15, 10, 4, 117, 116, 105, 97, 18, 7, 55, 48, 48, 48, 48, 48, 48, 18, 88, 10, 80, 10, 70, 10, 31, 47, 99, 111, 115, 109, 111, 115, 46, 99, 114, 121, 112, 116, 111, 46, 115, 101, 99, 112, 50, 53, 54, 107, 49, 46, 80, 117, 98, 75, 101, 121, 18, 35, 10, 33, 2, 214, 196, 150, 138, 247, 194, 102, 99, 26, 107, 77, 58, 49, 185, 175, 141, 130, 161, 143, 190, 103, 32, 58, 186, 68, 20, 160, 25, 160, 135, 214, 93, 18, 4, 10, 2, 8, 1, 24, 16, 18, 4, 16, 208, 232, 12, 26, 64, 130, 232, 165, 58, 164, 111, 95, 148, 20, 60, 156, 116, 178, 169, 117, 153, 98, 157, 196, 77, 197, 213, 72, 128, 216, 230, 87, 132, 221, 235, 144, 244, 43, 210, 127, 94, 48, 55, 233, 145, 153, 238, 250, 34, 139, 7, 50, 77, 206, 206, 47, 38, 39, 163, 8, 34, 220, 47, 197, 168, 59, 78, 221, 207,
	}

	var d = NewDecodedTx()
	err := decodeCosmosTx(txDecoder, rawTx, &d)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), d.TimeoutHeight)
	assert.Equal(t, "", d.Memo)

	expectedMsgs := []types.Msg{
		&cosmosStakingTypes.MsgDelegate{
			DelegatorAddress: "celestia14zfnc2kxdgdmacnuuytre5p6fx97ttfq9eggxd",
			ValidatorAddress: "celestiavaloper19urg9awjzwq8d40vwjdvv0yw9kgehscf0zx3gs",
			Amount: types.Coin{
				Denom:  "utia",
				Amount: math.NewInt(7000000),
			},
		},
	}
	assert.Equal(t, expectedMsgs, d.Messages)
	assert.Equal(t, "0", d.Fee.String())
}

func TestDecodeCosmosTx_VoteMsg(t *testing.T) {
	rawTx, err := base64.StdEncoding.DecodeString("ClEKTwoWL2Nvc21vcy5nb3YudjEuTXNnVm90ZRI1CAQSL2NlbGVzdGlhMTJ6czdlM244cGpkOHk4ZXgwY3l2NjdldGh2MzBtZWtncXU2NjVyGAESaApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAyJq13zdVvBc4sHiqsxdwmZuhu/+7jp5qybynJAUP4VeEgQKAggBGEQSFAoOCgR1dGlhEgY1MDAwMDAQ85EFGkB0CjjkpeDX/bfNeifAKWUWMSf5l7l8DqsDosnuQK3XMjiTlXN4AthomxLpSDqS/i7fsV7cLnaKV2trwJR5FvTc")
	require.NoError(t, err)

	var d = NewDecodedTx()
	err = decodeCosmosTx(txDecoder, rawTx, &d)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), d.TimeoutHeight)
	assert.Equal(t, "", d.Memo)

	expectedMsgs := []types.Msg{
		&cosmosGovTypesV1.MsgVote{
			ProposalId: 4,
			Voter:      "celestia12zs7e3n8pjd8y8ex0cyv67ethv30mekgqu665r",
			Option:     cosmosGovTypesV1.OptionYes,
		},
	}
	assert.Equal(t, expectedMsgs, d.Messages)
	assert.Equal(t, "500000", d.Fee.String())
}

func TestDecodeCosmosTx_MsgSetToken(t *testing.T) {
	rawTx, err := base64.StdEncoding.DecodeString("Ct8BCtwBCh4vaHlwZXJsYW5lLndhcnAudjEuTXNnU2V0VG9rZW4SuQEKL2NlbGVzdGlhMWxnMGU5bjRwdDI5bHBxMms0cHR1ZTRja3cwOWR4MGF1amxwZTRqEkIweDcyNmY3NTc0NjU3MjVmNjE3MDcwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMjAwMDAwMDAwMDAwMDAwMDEiQjB4NzI2Zjc1NzQ2NTcyNWY2OTczNmQwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMRJmClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDWNM912meXMULRdagfSqXn5vi7myaXeqmuhR3K6DMYlASBAoCCAEYFhISCgwKBHV0aWESBDc3NzUQtt8EGkCEcA6E3QIjVOQI+3ufuuxCj408FXKeRhx/WrvcD2qcqjeebz+DUCKLygNkvCK+otS7gPYDycS5nZw8kBtyuTn1")
	require.NoError(t, err)

	var d = NewDecodedTx()
	err = decodeCosmosTx(txDecoder, rawTx, &d)
	require.NoError(t, err)
	require.Equal(t, uint64(0), d.TimeoutHeight)
	require.Equal(t, "", d.Memo)

	tokenId, err := util.DecodeHexAddress("0x726f757465725f61707000000000000000000000000000020000000000000001")
	require.NoError(t, err)

	ismId, err := util.DecodeHexAddress("0x726f757465725f69736d00000000000000000000000000000000000000000001")
	require.NoError(t, err)

	expectedMsgs := []types.Msg{
		&hyperlaneWarp.MsgSetToken{
			Owner:             "celestia1lg0e9n4pt29lpq2k4ptue4ckw09dx0aujlpe4j",
			RenounceOwnership: false,
			NewOwner:          "",
			TokenId:           tokenId,
			IsmId:             &ismId,
		},
	}
	require.Equal(t, expectedMsgs, d.Messages)
	require.Equal(t, "7775", d.Fee.String())
}

func TestDecodeFee(t *testing.T) {
	testCases := []struct {
		desc        string
		amount      types.Coins
		expectedFee decimal.Decimal
		expectedErr string
	}{
		{
			desc:        "No fee",
			amount:      nil,
			expectedFee: decimal.Zero,
			expectedErr: "",
		},
		{
			desc: "Valid UTIA fee",
			amount: types.Coins{
				types.NewCoin("utia", math.NewInt(1000)),
			},
			expectedFee: decimal.NewFromInt(1000),
			expectedErr: "",
		},
		{
			desc: "Valid TIA fee",
			amount: types.Coins{
				types.NewCoin("tia", math.NewInt(5000000)),
			},
			expectedFee: decimal.NewFromInt(5000000).Shift(6),
			expectedErr: "",
		},
		{
			desc: "Multiple fee currencies",
			amount: types.Coins{
				types.NewCoin("utia", math.NewInt(1000)),
				types.NewCoin("tia", math.NewInt(5000000)),
			},
			expectedFee: decimal.Zero,
			expectedErr: "found fee in 2 currencies",
		},
		{
			desc: "Fee in unknown denom",
			amount: types.Coins{
				types.NewCoin("unknown", math.NewInt(1000)),
			},
			expectedFee: decimal.Zero,
			expectedErr: "couldn't find fee amount in utia or in tia denom",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			fee, err := decodeFee(tc.amount)

			assert.Equal(t, tc.expectedFee, fee)
			if err != nil {
				assert.Equal(t, tc.expectedErr, err.Error())
			} else {
				assert.Equal(t, tc.expectedErr, "")
			}

		})
	}
}

func TestGetFeeInDenom(t *testing.T) {
	testCases := []struct {
		desc        string
		amount      types.Coins
		denom       string
		expectedFee decimal.Decimal
		expectedOk  bool
	}{
		{
			desc: "Valid UTIA fee",
			amount: types.Coins{
				types.NewCoin("utia", math.NewInt(1000)),
			},
			denom:       currency.Utia,
			expectedFee: decimal.NewFromInt(1000),
			expectedOk:  true,
		},
		{
			desc: "Valid TIA fee",
			amount: types.Coins{
				types.NewCoin("tia", math.NewInt(5000000)),
			},
			denom:       currency.Tia,
			expectedFee: decimal.NewFromInt(5000000).Shift(6),
			expectedOk:  true,
		},
		{
			desc: "Fee in unknown denom",
			amount: types.Coins{
				types.NewCoin("unknown", math.NewInt(1000)),
			},
			denom:       currency.Utia,
			expectedFee: decimal.Zero,
			expectedOk:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			fee, ok := getFeeInDenom(tc.amount, tc.denom)

			assert.Equal(t, tc.expectedFee, fee)
			assert.Equal(t, tc.expectedOk, ok)
		})
	}
}
