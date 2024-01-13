// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	nodeTypes "github.com/celenium-io/celestia-indexer/pkg/types"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
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

func TestDecodeAuthInfo_WithNilAmount(t *testing.T) {
	rawTx := []byte{
		10, 164, 1, 10, 161, 1, 10, 35, 47, 99, 111, 115, 109, 111, 115, 46, 115, 116, 97, 107, 105, 110, 103, 46, 118, 49, 98, 101, 116, 97, 49, 46, 77, 115, 103, 68, 101, 108, 101, 103, 97, 116, 101, 18, 122, 10, 47, 99, 101, 108, 101, 115, 116, 105, 97, 49, 52, 122, 102, 110, 99, 50, 107, 120, 100, 103, 100, 109, 97, 99, 110, 117, 117, 121, 116, 114, 101, 53, 112, 54, 102, 120, 57, 55, 116, 116, 102, 113, 57, 101, 103, 103, 120, 100, 18, 54, 99, 101, 108, 101, 115, 116, 105, 97, 118, 97, 108, 111, 112, 101, 114, 49, 57, 117, 114, 103, 57, 97, 119, 106, 122, 119, 113, 56, 100, 52, 48, 118, 119, 106, 100, 118, 118, 48, 121, 119, 57, 107, 103, 101, 104, 115, 99, 102, 48, 122, 120, 51, 103, 115, 26, 15, 10, 4, 117, 116, 105, 97, 18, 7, 55, 48, 48, 48, 48, 48, 48, 18, 88, 10, 80, 10, 70, 10, 31, 47, 99, 111, 115, 109, 111, 115, 46, 99, 114, 121, 112, 116, 111, 46, 115, 101, 99, 112, 50, 53, 54, 107, 49, 46, 80, 117, 98, 75, 101, 121, 18, 35, 10, 33, 2, 214, 196, 150, 138, 247, 194, 102, 99, 26, 107, 77, 58, 49, 185, 175, 141, 130, 161, 143, 190, 103, 32, 58, 186, 68, 20, 160, 25, 160, 135, 214, 93, 18, 4, 10, 2, 8, 1, 24, 16, 18, 4, 16, 208, 232, 12, 26, 64, 130, 232, 165, 58, 164, 111, 95, 148, 20, 60, 156, 116, 178, 169, 117, 153, 98, 157, 196, 77, 197, 213, 72, 128, 216, 230, 87, 132, 221, 235, 144, 244, 43, 210, 127, 94, 48, 55, 233, 145, 153, 238, 250, 34, 139, 7, 50, 77, 206, 206, 47, 38, 39, 163, 8, 34, 220, 47, 197, 168, 59, 78, 221, 207,
	}

	authInfo, fee, err := decodeAuthInfo(cfg, rawTx)
	assert.NoError(t, err)
	assert.Equal(t, decimal.Zero, fee)
	assert.Equal(t, uint64(210000), authInfo.Fee.GasLimit)
}

func TestDecodeAuthInfo_WithFee(t *testing.T) {
	rawTx := []byte{
		10, 171, 1, 10, 168, 1, 10, 35, 47, 99, 111, 115, 109, 111, 115, 46, 115, 116, 97, 107, 105, 110, 103, 46, 118, 49, 98, 101, 116, 97, 49, 46, 77, 115, 103, 68, 101, 108, 101, 103, 97, 116, 101, 18, 128, 1, 10, 47, 99, 101, 108, 101, 115, 116, 105, 97, 49, 55, 97, 100, 115, 106, 107, 117, 101, 99, 103, 106, 104, 101, 117, 103, 114, 100, 114, 119, 100, 113, 118, 57, 117, 104, 51, 113, 107, 114, 102, 109, 106, 57, 120, 122, 97, 119, 120, 18, 54, 99, 101, 108, 101, 115, 116, 105, 97, 118, 97, 108, 111, 112, 101, 114, 49, 55, 97, 100, 115, 106, 107, 117, 101, 99, 103, 106, 104, 101, 117, 103, 114, 100, 114, 119, 100, 113, 118, 57, 117, 104, 51, 113, 107, 114, 102, 109, 106, 113, 101, 113, 121, 99, 113, 26, 21, 10, 4, 117, 116, 105, 97, 18, 13, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 18, 104, 10, 81, 10, 70, 10, 31, 47, 99, 111, 115, 109, 111, 115, 46, 99, 114, 121, 112, 116, 111, 46, 115, 101, 99, 112, 50, 53, 54, 107, 49, 46, 80, 117, 98, 75, 101, 121, 18, 35, 10, 33, 2, 5, 5, 146, 95, 90, 69, 253, 244, 240, 130, 93, 143, 158, 212, 70, 117, 227, 56, 38, 141, 84, 101, 29, 76, 145, 143, 105, 95, 140, 136, 230, 156, 18, 4, 10, 2, 8, 1, 24, 170, 1, 18, 19, 10, 13, 10, 4, 117, 116, 105, 97, 18, 5, 50, 49, 48, 48, 48, 16, 208, 232, 12, 26, 64, 93, 57, 117, 108, 143, 235, 212, 126, 28, 128, 252, 240, 168, 77, 60, 219, 10, 189, 241, 178, 117, 145, 177, 79, 112, 156, 36, 73, 6, 88, 0, 182, 72, 92, 192, 27, 7, 4, 51, 165, 1, 44, 21, 25, 78, 128, 31, 101, 86, 247, 159, 82, 136, 212, 79, 118, 139, 241, 135, 91, 205, 125, 100, 77,
	}

	authInfo, fee, err := decodeAuthInfo(cfg, rawTx)
	assert.NoError(t, err)
	assert.Equal(t, decimal.NewFromInt(21000), fee)
	assert.Equal(t, uint64(210000), authInfo.Fee.GasLimit)
}

func TestDecodeCosmosTx_DelegateMsg(t *testing.T) {
	rawTx := []byte{
		10, 164, 1, 10, 161, 1, 10, 35, 47, 99, 111, 115, 109, 111, 115, 46, 115, 116, 97, 107, 105, 110, 103, 46, 118, 49, 98, 101, 116, 97, 49, 46, 77, 115, 103, 68, 101, 108, 101, 103, 97, 116, 101, 18, 122, 10, 47, 99, 101, 108, 101, 115, 116, 105, 97, 49, 52, 122, 102, 110, 99, 50, 107, 120, 100, 103, 100, 109, 97, 99, 110, 117, 117, 121, 116, 114, 101, 53, 112, 54, 102, 120, 57, 55, 116, 116, 102, 113, 57, 101, 103, 103, 120, 100, 18, 54, 99, 101, 108, 101, 115, 116, 105, 97, 118, 97, 108, 111, 112, 101, 114, 49, 57, 117, 114, 103, 57, 97, 119, 106, 122, 119, 113, 56, 100, 52, 48, 118, 119, 106, 100, 118, 118, 48, 121, 119, 57, 107, 103, 101, 104, 115, 99, 102, 48, 122, 120, 51, 103, 115, 26, 15, 10, 4, 117, 116, 105, 97, 18, 7, 55, 48, 48, 48, 48, 48, 48, 18, 88, 10, 80, 10, 70, 10, 31, 47, 99, 111, 115, 109, 111, 115, 46, 99, 114, 121, 112, 116, 111, 46, 115, 101, 99, 112, 50, 53, 54, 107, 49, 46, 80, 117, 98, 75, 101, 121, 18, 35, 10, 33, 2, 214, 196, 150, 138, 247, 194, 102, 99, 26, 107, 77, 58, 49, 185, 175, 141, 130, 161, 143, 190, 103, 32, 58, 186, 68, 20, 160, 25, 160, 135, 214, 93, 18, 4, 10, 2, 8, 1, 24, 16, 18, 4, 16, 208, 232, 12, 26, 64, 130, 232, 165, 58, 164, 111, 95, 148, 20, 60, 156, 116, 178, 169, 117, 153, 98, 157, 196, 77, 197, 213, 72, 128, 216, 230, 87, 132, 221, 235, 144, 244, 43, 210, 127, 94, 48, 55, 233, 145, 153, 238, 250, 34, 139, 7, 50, 77, 206, 206, 47, 38, 39, 163, 8, 34, 220, 47, 197, 168, 59, 78, 221, 207,
	}

	timeoutHeight, memo, messages, err := decodeCosmosTx(decoder, rawTx)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), timeoutHeight)
	assert.Equal(t, "", memo)

	expectedMsgs := []types.Msg{
		&cosmosStakingTypes.MsgDelegate{
			DelegatorAddress: "celestia14zfnc2kxdgdmacnuuytre5p6fx97ttfq9eggxd",
			ValidatorAddress: "celestiavaloper19urg9awjzwq8d40vwjdvv0yw9kgehscf0zx3gs",
			Amount: types.Coin{
				Denom:  "utia",
				Amount: types.NewInt(7000000),
			},
		},
	}
	assert.Equal(t, expectedMsgs, messages)
}

func TestDecodeFee(t *testing.T) {
	testCases := []struct {
		desc        string
		authInfo    tx.AuthInfo
		expectedFee decimal.Decimal
		expectedErr string
	}{
		{
			desc:        "No fee",
			authInfo:    tx.AuthInfo{},
			expectedFee: decimal.Zero,
			expectedErr: "",
		},
		{
			desc: "Valid UTIA fee",
			authInfo: tx.AuthInfo{
				Fee: &tx.Fee{
					Amount: types.Coins{
						types.NewCoin("utia", types.NewInt(1000)),
					},
				},
			},
			expectedFee: decimal.NewFromInt(1000),
			expectedErr: "",
		},
		{
			desc: "Valid TIA fee",
			authInfo: tx.AuthInfo{
				Fee: &tx.Fee{
					Amount: types.Coins{
						types.NewCoin("tia", types.NewInt(5000000)),
					},
				},
			},
			expectedFee: decimal.NewFromInt(5000000).Shift(6),
			expectedErr: "",
		},
		{
			desc: "Multiple fee currencies",
			authInfo: tx.AuthInfo{
				Fee: &tx.Fee{
					Amount: types.Coins{
						types.NewCoin("utia", types.NewInt(1000)),
						types.NewCoin("tia", types.NewInt(5000000)),
					},
				},
			},
			expectedFee: decimal.Zero,
			expectedErr: "found fee in 2 currencies",
		},
		{
			desc: "Fee in unknown denom",
			authInfo: tx.AuthInfo{
				Fee: &tx.Fee{
					Amount: types.Coins{
						types.NewCoin("unknown", types.NewInt(1000)),
					},
				},
			},
			expectedFee: decimal.Zero,
			expectedErr: "couldn't find fee amount in utia or in tia denom",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			fee, err := decodeFee(tc.authInfo)

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
		denom       currency.Denom
		expectedFee decimal.Decimal
		expectedOk  bool
	}{
		{
			desc: "Valid UTIA fee",
			amount: types.Coins{
				types.NewCoin("utia", types.NewInt(1000)),
			},
			denom:       currency.Utia,
			expectedFee: decimal.NewFromInt(1000),
			expectedOk:  true,
		},
		{
			desc: "Valid TIA fee",
			amount: types.Coins{
				types.NewCoin("tia", types.NewInt(5000000)),
			},
			denom:       currency.Tia,
			expectedFee: decimal.NewFromInt(5000000).Shift(6),
			expectedOk:  true,
		},
		{
			desc: "Fee in unknown denom",
			amount: types.Coins{
				types.NewCoin("unknown", types.NewInt(1000)),
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
