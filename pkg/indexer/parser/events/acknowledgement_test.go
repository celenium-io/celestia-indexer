// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package events

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/stretchr/testify/require"
)

func Test_handleAcknowledgement(t *testing.T) {
	tests := []struct {
		name   string
		ctx    *context.Context
		events []storage.Event
		msg    []*storage.Message
		idx    *int

		isTransferNil bool
	}{
		{
			name:          "test 1",
			ctx:           context.NewContext(),
			isTransferNil: true,
			events: []storage.Event{
				{
					Height: 2371609,
					Type:   "message",
					Data: map[string]any{
						"action": "/ibc.core.channel.v1.MsgAcknowledgement",
					},
				}, {
					Height: 2371609,
					Type:   "acknowledge_packet",
					Data: map[string]any{
						"packet_channel_ordering":  "ORDER_UNORDERED",
						"packet_connection":        "connection-2",
						"packet_dst_channel":       "channel-6994",
						"packet_dst_port":          "transfer",
						"packet_sequence":          "1004901",
						"packet_src_channel":       "channel-2",
						"packet_src_port":          "transfer",
						"packet_timeout_height":    "0-0",
						"packet_timeout_timestamp": "1726667321908511033",
					},
				}, {
					Height: 2371609,
					Type:   "message",
					Data: map[string]any{
						"module": "ibc_channel",
					},
				},
				{
					Height: 2371609,
					Type:   "write_acknowledgement",
					Data: map[string]any{
						"packet_ack":               "{\"result\":\"AQ==\"}",
						"packet_ack_hex":           "7b22726573756c74223a2241513d3d227d",
						"packet_connection":        "connection-18",
						"packet_data":              "{\"amount\":\"561801\",\"denom\":\"transfer/channel-119/utia\",\"memo\":\"{\\\"forward\\\":{\\\"receiver\\\":\\\"osmo1345fue0f2zwmfef4d48qfe38k0wfvca657jkm0\\\",\\\"port\\\":\\\"transfer\\\",\\\"channel\\\":\\\"channel-2\\\"}}\",\"receiver\":\"celestia1345fue0f2zwmfef4d48qfe38k0wfvca6d0skhs\",\"sender\":\"umee1345fue0f2zwmfef4d48qfe38k0wfvca6wnuef0\"}",
						"packet_data_hex":          "7b22616d6f756e74223a22353631383031222c2264656e6f6d223a227472616e736665722f6368616e6e656c2d3131392f75746961222c226d656d6f223a227b5c22666f72776172645c223a7b5c2272656365697665725c223a5c226f736d6f313334356675653066327a776d6665663464343871666533386b3077667663613635376a6b6d305c222c5c22706f72745c223a5c227472616e736665725c222c5c226368616e6e656c5c223a5c226368616e6e656c2d325c227d7d222c227265636569766572223a2263656c6573746961313334356675653066327a776d6665663464343871666533386b307766766361366430736b6873222c2273656e646572223a22756d6565313334356675653066327a776d6665663464343871666533386b30776676636136776e75656630227d",
						"packet_dst_channel":       "channel-19",
						"packet_dst_port":          "transfer",
						"packet_sequence":          "1658",
						"packet_src_channel":       "channel-119",
						"packet_src_port":          "transfer",
						"packet_timeout_height":    "0-2371753",
						"packet_timeout_timestamp": "0",
					},
				}, {
					Height: 2371609,
					Type:   "message",
					Data: map[string]any{
						"module": "ibc_channel",
					},
				},
			},
			msg: []*storage.Message{
				{
					Type:   types.MsgAcknowledgement,
					Height: 2371609,
				},
			},
			idx: testsuite.Ptr(0),
		}, {
			name:          "test 2",
			ctx:           context.NewContext(),
			isTransferNil: false,
			events: []storage.Event{
				{
					Height: 2371609,
					Type:   "message",
					Data: map[string]any{
						"action": "/ibc.core.channel.v1.MsgAcknowledgement",
					},
				}, {
					Height: 2371609,
					Type:   "acknowledge_packet",
					Data: map[string]any{
						"packet_channel_ordering":  "ORDER_UNORDERED",
						"packet_connection":        "connection-2",
						"packet_dst_channel":       "channel-6994",
						"packet_dst_port":          "transfer",
						"packet_sequence":          "1004900",
						"packet_src_channel":       "channel-2",
						"packet_src_port":          "transfer",
						"packet_timeout_height":    "1-21043368",
						"packet_timeout_timestamp": "0",
					},
				}, {
					Height: 2371609,
					Type:   "message",
					Data: map[string]any{
						"module": "ibc_channel",
					},
				},
				{
					Height: 2371609,
					Type:   "fungible_token_packet",
					Data: map[string]any{"acknowledgement": "result:\"{\\\"contract_result\\\":null,\\\"ibc_ack\\\":\\\"eyJyZXN1bHQiOiJBUT09In0=\\\"}\" ",
						"amount":   "12000000",
						"denom":    "utia",
						"memo":     "{\"wasm\":{\"contract\":\"osmo1vkdakqqg5htq5c3wy2kj2geq536q665xdexrtjuwqckpads2c2nsvhhcyv\",\"msg\":{\"swap_and_action\":{\"user_swap\":{\"swap_exact_asset_in\":{\"swap_venue_name\":\"osmosis-poolmanager\",\"operations\":[{\"pool\":\"1247\",\"denom_in\":\"ibc/D79E7D83AB399BFFF93433E54FAA480C191248FC556924A2A8351AE2638B3877\",\"denom_out\":\"ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4\"},{\"pool\":\"1319\",\"denom_in\":\"ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4\",\"denom_out\":\"ibc/64BA6E31FE887D66C6F8F31C7B1A80C7CA179239677B4088BB55F5EA07DBE273\"}]}},\"min_asset\":{\"native\":{\"denom\":\"ibc/64BA6E31FE887D66C6F8F31C7B1A80C7CA179239677B4088BB55F5EA07DBE273\",\"amount\":\"3298454343374009626\"}},\"timeout_timestamp\":1726667021008444464,\"post_swap_action\":{\"ibc_transfer\":{\"ibc_info\":{\"source_channel\":\"channel-122\",\"receiver\":\"inj14amztqem07qvyyty8k4p6s2jp8ylsvlax0g42f\",\"memo\":\"\",\"recover_address\":\"osmo1nc44cmtgmp6cwch2ccfp6txdelu6qtz9rq5s03\"}}},\"affiliates\":[{\"basis_points_fee\":\"60\",\"address\":\"osmo1my4tk420gjmhggqwvvha6ey9390gqwfree2p4u\"},{\"basis_points_fee\":\"15\",\"address\":\"osmo1msjnal2glfz6ze8x9kduhg45xppx9sddawfx46\"}]}}}}",
						"module":   "transfer",
						"receiver": "osmo1vkdakqqg5htq5c3wy2kj2geq536q665xdexrtjuwqckpads2c2nsvhhcyv",
						"sender":   "celestia1nc44cmtgmp6cwch2ccfp6txdelu6qtz963ksrw",
					},
				}, {
					Height: 2371609,
					Type:   "fungible_token_packet",
					Data: map[string]any{"acknowledgement": "result:\"{\\\"contract_result\\\":null,\\\"ibc_ack\\\":\\\"eyJyZXN1bHQiOiJBUT09In0=\\\"}\" ",
						"success": "{\"contract_result\":null,\"ibc_ack\":\"eyJyZXN1bHQiOiJBUT09In0=\"}",
					},
				},
			},
			msg: []*storage.Message{
				{
					Type:   types.MsgAcknowledgement,
					Height: 2371609,
				},
			},
			idx: testsuite.Ptr(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := range tt.msg {
				err := handleAcknowledgement(tt.ctx, tt.events, tt.msg[i], tt.idx)
				require.NoError(t, err)
				if tt.isTransferNil {
					require.Nil(t, tt.msg[i].IbcTransfer)
				} else {
					require.NotNil(t, tt.msg[i].IbcTransfer)
				}
			}
		})
	}
}
