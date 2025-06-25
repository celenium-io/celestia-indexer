// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package decode

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
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
				Amount:  testsuite.Ptr(types.NewCoin("utia", math.OneInt())),
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
				Amount:   testsuite.Ptr(types.NewCoin("utia", math.NewInt(42))),
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

func TestNewCompleteRedelegation(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]any
		wantBody CompleteRedelegation
		wantErr  bool
	}{
		{
			name: "test 1",
			m: map[string]any{
				"amount":                "51000000utia",
				"delegator":             "celestia1jxmxxa04k2tpkwk5u00pj37lmg668ksvh0ydax",
				"destination_validator": "celestiavaloper1clf3nqp89h97umhl4fmcqr642jz6rszcxegjc6",
				"source_validator":      "celestiavaloper1wu24jxpn9j0580ehjz344d58cf3t7lzrrgqmnr",
			},
			wantBody: CompleteRedelegation{
				Amount:        testsuite.Ptr(types.NewCoin("utia", math.NewInt(51000000))),
				Delegator:     "celestia1jxmxxa04k2tpkwk5u00pj37lmg668ksvh0ydax",
				SrcValidator:  "celestiavaloper1wu24jxpn9j0580ehjz344d58cf3t7lzrrgqmnr",
				DestValidator: "celestiavaloper1clf3nqp89h97umhl4fmcqr642jz6rszcxegjc6",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, err := NewCompleteRedelegation(tt.m)
			require.True(t, (err != nil) == tt.wantErr)
			require.Equal(t, tt.wantBody, gotBody)
		})
	}
}

func TestNewCompleteUnbonding(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]any
		wantBody CompleteUnbonding
		wantErr  bool
	}{
		{
			name: "test 1",
			m: map[string]any{
				"amount":    "256000000utia",
				"delegator": "celestia1g60teezwmfdj8xxpnd5kehvp25zfzt25pxxphv",
				"validator": "celestiavaloper1r5xt7twqmh39ky72f4txxjrhlt2z0qwwmdal8c",
			},
			wantBody: CompleteUnbonding{
				Amount:    testsuite.Ptr(types.NewCoin("utia", math.NewInt(256000000))),
				Delegator: "celestia1g60teezwmfdj8xxpnd5kehvp25zfzt25pxxphv",
				Validator: "celestiavaloper1r5xt7twqmh39ky72f4txxjrhlt2z0qwwmdal8c",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, err := NewCompleteUnbonding(tt.m)
			require.True(t, (err != nil) == tt.wantErr)
			require.Equal(t, tt.wantBody, gotBody)
		})
	}
}

func TestNewCommission(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]any
		wantBody Commission
		wantErr  bool
	}{
		{
			name: "test 1",
			m: map[string]any{
				"amount":    "256000000utia",
				"validator": "celestiavaloper1r5xt7twqmh39ky72f4txxjrhlt2z0qwwmdal8c",
			},
			wantBody: Commission{
				Amount:    decimal.RequireFromString("256000000"),
				Validator: "celestiavaloper1r5xt7twqmh39ky72f4txxjrhlt2z0qwwmdal8c",
			},
			wantErr: false,
		}, {
			name: "test 2",
			m: map[string]any{
				"amount":    "469.815871531603829656utia",
				"validator": "celestiavaloper189ecvq5avj0wehrcfnagpd5sd8pup9aqmdglmr",
			},
			wantBody: Commission{
				Amount:    decimal.RequireFromString("469.815871531603829656"),
				Validator: "celestiavaloper189ecvq5avj0wehrcfnagpd5sd8pup9aqmdglmr",
			},
			wantErr: false,
		}, {
			name: "test 3",
			m: map[string]any{
				"amount":    "",
				"validator": "celestiavaloper189ecvq5avj0wehrcfnagpd5sd8pup9aqmdglmr",
			},
			wantBody: Commission{
				Amount:    decimal.Zero,
				Validator: "celestiavaloper189ecvq5avj0wehrcfnagpd5sd8pup9aqmdglmr",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, err := NewCommission(tt.m)
			require.True(t, (err != nil) == tt.wantErr)
			require.Equal(t, tt.wantBody, gotBody)
		})
	}
}

func TestNewRewards(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]any
		wantBody Rewards
		wantErr  bool
	}{
		{
			name: "test 1",
			m: map[string]any{
				"amount":    "256000000utia",
				"validator": "celestiavaloper1r5xt7twqmh39ky72f4txxjrhlt2z0qwwmdal8c",
			},
			wantBody: Rewards{
				Amount:    decimal.RequireFromString("256000000"),
				Validator: "celestiavaloper1r5xt7twqmh39ky72f4txxjrhlt2z0qwwmdal8c",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, err := NewRewards(tt.m)
			require.True(t, (err != nil) == tt.wantErr)
			require.Equal(t, tt.wantBody, gotBody)
		})
	}
}

func TestNewWithdrawReward(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]any
		wantBody WithdrawRewards
		wantErr  bool
	}{
		{
			name: "test 1",
			m: map[string]any{
				"amount":    "256000000utia",
				"delegator": "celestia1g60teezwmfdj8xxpnd5kehvp25zfzt25pxxphv",
				"validator": "celestiavaloper1r5xt7twqmh39ky72f4txxjrhlt2z0qwwmdal8c",
			},
			wantBody: WithdrawRewards{
				Amount:    testsuite.Ptr(types.NewCoin("utia", math.NewInt(256000000))),
				Delegator: "celestia1g60teezwmfdj8xxpnd5kehvp25zfzt25pxxphv",
				Validator: "celestiavaloper1r5xt7twqmh39ky72f4txxjrhlt2z0qwwmdal8c",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, err := NewWithdrawRewards(tt.m)
			require.True(t, (err != nil) == tt.wantErr)
			require.Equal(t, tt.wantBody, gotBody)
		})
	}
}

func TestNewRedelegate(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]any
		wantBody Redelegate
		wantErr  bool
	}{
		{
			name: "test 1",
			m: map[string]any{
				"amount":                "69989816utia",
				"completion_time":       "2024-03-10T22:58:16Z",
				"destination_validator": "celestiavaloper1u825srldhev7t4wnd3hplhrphahjfk7ff3wfdr",
				"source_validator":      "celestiavaloper1rcm7tth05klgkqpucdhm5hexnk49dfda3l3hak",
			},
			wantBody: Redelegate{
				Amount:         testsuite.Ptr(types.NewCoin("utia", math.NewInt(69989816))),
				SrcValidator:   "celestiavaloper1rcm7tth05klgkqpucdhm5hexnk49dfda3l3hak",
				DestValidator:  "celestiavaloper1u825srldhev7t4wnd3hplhrphahjfk7ff3wfdr",
				CompletionTime: time.Date(2024, 3, 10, 22, 58, 16, 0, time.UTC),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, err := NewRedelegate(tt.m)
			require.True(t, (err != nil) == tt.wantErr)
			require.Equal(t, tt.wantBody, gotBody)
		})
	}
}

func TestNewUnbond(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]any
		wantBody Unbond
		wantErr  bool
	}{
		{
			name: "test 1",
			m: map[string]any{
				"amount":          "1000000utia",
				"completion_time": "2024-03-10T23:08:39Z",
				"validator":       "celestiavaloper1ej2es5fjztqjcd4pwa0zyvaevtjd2y5wh8xeg4",
			},
			wantBody: Unbond{
				Amount:         testsuite.Ptr(types.NewCoin("utia", math.NewInt(1000000))),
				Validator:      "celestiavaloper1ej2es5fjztqjcd4pwa0zyvaevtjd2y5wh8xeg4",
				CompletionTime: time.Date(2024, 3, 10, 23, 8, 39, 0, time.UTC),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, err := NewUnbond(tt.m)
			require.True(t, (err != nil) == tt.wantErr)
			require.Equal(t, tt.wantBody, gotBody)
		})
	}
}

func TestNewUpdateClient(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]any
		wantBody UpdateClient
		wantErr  bool
	}{
		{
			name: "test 1",
			m: map[string]any{
				"client_id":        "07-tendermint-145",
				"client_type":      "07-tendermint",
				"consensus_height": "3-884",
			},
			wantBody: UpdateClient{
				Id:              "07-tendermint-145",
				Type:            "07-tendermint",
				Revision:        3,
				ConsensusHeight: 884,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, err := NewUpdateClient(tt.m)
			require.True(t, (err != nil) == tt.wantErr)
			require.Equal(t, tt.wantBody, gotBody)
		})
	}
}

func TestNewConnectionChange(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]any
		wantBody ConnectionChange
	}{
		{
			name: "test 1",
			m: map[string]any{
				"client_id":                  "07-tendermint-145",
				"connection_id":              "connection-97",
				"counterparty_client_id":     "07-tendermint-1",
				"counterparty_connection_id": "connection-1",
			},
			wantBody: ConnectionChange{
				ClientId:                 "07-tendermint-145",
				ConnectionId:             "connection-97",
				CounterpartyClientId:     "07-tendermint-1",
				CounterpartyConnectionId: "connection-1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody := NewConnectionOpen(tt.m)
			require.Equal(t, tt.wantBody, gotBody)
		})
	}
}

func TestNewChannelChange(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]any
		wantBody ChannelChange
	}{
		{
			name: "test 1",
			m: map[string]any{
				"channel_id":              "channel-112",
				"connection_id":           "connection-97",
				"counterparty_channel_id": "channel-1",
				"counterparty_port_id":    "transfer",
				"port_id":                 "transfer",
			},
			wantBody: ChannelChange{
				ChannelId:             "channel-112",
				ConnectionId:          "connection-97",
				CounterpartyChannelId: "channel-1",
				CounterpartyPortId:    "transfer",
				PortId:                "transfer",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody := NewChannelChange(tt.m)
			require.Equal(t, tt.wantBody, gotBody)
		})
	}
}

func TestNewFungibleTokenPacket(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]any
		wantBody FungibleTokenPacket
	}{
		{
			name: "test 1",
			m: map[string]any{
				"amount":   "699567",
				"denom":    "transfer/channel-6994/utia",
				"memo":     "",
				"module":   "transfer",
				"receiver": "celestia1j8qsd8f7mdcm5skfs50tv9nat9g00qh8y2zd39",
				"sender":   "osmo1m8wg4vxkefhs374qxmmqpyusgz289wmulex5qdwpfx7jnrxzer5s9cv83q",
				"success":  "true",
			},
			wantBody: FungibleTokenPacket{
				Amount:   decimal.RequireFromString("699567"),
				Denom:    "transfer/channel-6994/utia",
				Memo:     "",
				Module:   "transfer",
				Receiver: "celestia1j8qsd8f7mdcm5skfs50tv9nat9g00qh8y2zd39",
				Sender:   "osmo1m8wg4vxkefhs374qxmmqpyusgz289wmulex5qdwpfx7jnrxzer5s9cv83q",
				Success:  "true",
			},
		}, {
			name: "test 2",
			m: map[string]any{
				"acknowledgement": "result:\"\\001\" ",
				"amount":          "4745268",
				"denom":           "utia",
				"memo":            "",
				"module":          "transfer",
				"receiver":        "neutron1q7pct93xm8qunmnjyu2feezjznhyedsf968u72p6ma2mcg755t0susl2gv",
				"sender":          "celestia13qe9fxcd63ym5gt4fc235ugdv9zzjejuwky7glqq8xtdc66r9g6sn4vfr6",
			},
			wantBody: FungibleTokenPacket{
				Amount:   decimal.RequireFromString("4745268"),
				Denom:    "utia",
				Memo:     "",
				Module:   "transfer",
				Receiver: "neutron1q7pct93xm8qunmnjyu2feezjznhyedsf968u72p6ma2mcg755t0susl2gv",
				Sender:   "celestia13qe9fxcd63ym5gt4fc235ugdv9zzjejuwky7glqq8xtdc66r9g6sn4vfr6",
			},
		}, {
			name: "test 3",
			m: map[string]any{
				"success": "\u0001",
			},
			wantBody: FungibleTokenPacket{
				Success: "\u0001",
				Amount:  decimal.Zero,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody := NewFungibleTokenPacket(tt.m)
			require.Equal(t, tt.wantBody, gotBody)
		})
	}
}

func TestNewRecvPacket(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]any
		wantBody RecvPacket
	}{
		{
			name: "test 1",
			m: map[string]any{
				"packet_channel_ordering":  "ORDER_UNORDERED",
				"packet_connection":        "connection-2",
				"packet_data":              "{\"denom\":\"transfer/channel-6994/utia\",\"amount\":\"699567\",\"sender\":\"osmo1m8wg4vxkefhs374qxmmqpyusgz289wmulex5qdwpfx7jnrxzer5s9cv83q\",\"receiver\":\"celestia1j8qsd8f7mdcm5skfs50tv9nat9g00qh8y2zd39\"}",
				"packet_data_hex":          "7b2264656e6f6d223a227472616e736665722f6368616e6e656c2d363939342f75746961222c22616d6f756e74223a22363939353637222c2273656e646572223a226f736d6f316d3877673476786b6566687333373471786d6d7170797573677a323839776d756c657835716477706678376a6e72787a65723573396376383371222c227265636569766572223a2263656c6573746961316a387173643866376d64636d35736b667335307476396e61743967303071683879327a643339227d",
				"packet_dst_channel":       "channel-2",
				"packet_dst_port":          "transfer",
				"packet_sequence":          "1388301",
				"packet_src_channel":       "channel-6994",
				"packet_src_port":          "transfer",
				"packet_timeout_height":    "0-0",
				"packet_timeout_timestamp": "1747320319278917000",
			},
			wantBody: RecvPacket{
				Ordering:      "ORDER_UNORDERED",
				Connection:    "connection-2",
				Data:          "{\"denom\":\"transfer/channel-6994/utia\",\"amount\":\"699567\",\"sender\":\"osmo1m8wg4vxkefhs374qxmmqpyusgz289wmulex5qdwpfx7jnrxzer5s9cv83q\",\"receiver\":\"celestia1j8qsd8f7mdcm5skfs50tv9nat9g00qh8y2zd39\"}",
				DstChannel:    "channel-2",
				DstPort:       "transfer",
				Sequence:      1388301,
				SrcChannel:    "channel-6994",
				SrcPort:       "transfer",
				Timeout:       time.Date(2025, 05, 15, 14, 45, 19, 278917000, time.UTC),
				TimeoutHeight: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, err := NewRecvPacket(tt.m)
			require.NoError(t, err)
			require.Equal(t, tt.wantBody, gotBody)
		})
	}
}
