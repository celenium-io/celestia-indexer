// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	transferTypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/shopspring/decimal"
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
					Data: types.PackedBytes{
						"Acknowledgement": "eyJyZXN1bHQiOiJBUT09In0=",
						"Packet": map[string]any{
							"Data": transferTypes.FungibleTokenPacketData{
								Amount:   "561801",
								Denom:    "utia",
								Receiver: "osmo1345fue0f2zwmfef4d48qfe38k0wfvca657jkm0",
								Sender:   "celestia1ycjgmkjvjkmdwvjeuutxf6vxmfw9yk5cava5xt",
							},
							"DestinationChannel": "channel-6994",
							"DestinationPort":    "transfer",
							"Sequence":           1004901,
							"SourceChannel":      "channel-2",
							"SourcePort":         "transfer",
							"TimeoutHeight": map[string]any{
								"RevisionHeight": 0,
								"RevisionNumber": 0,
							},
							"TimeoutTimestamp": 1726667321908511000,
						},
						"ProofAcked": "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",
						"ProofHeight": map[string]any{
							"RevisionHeight": 21043235,
							"RevisionNumber": 1,
						},
						"Signer": "celestia16p6lrlxf7f03c0ka8cv4sznr29rym27us7u4v6",
					},
					IbcTransfer: &storage.IbcTransfer{
						Height:          2371609,
						Amount:          decimal.RequireFromString("561801"),
						Denom:           "utia",
						ReceiverAddress: testsuite.Ptr("osmo1345fue0f2zwmfef4d48qfe38k0wfvca657jkm0"),
						Sender: &storage.Address{
							Address: "celestia1ycjgmkjvjkmdwvjeuutxf6vxmfw9yk5cava5xt",
						},
						ChannelId: "channel-2",
						Port:      "transfer",
						Sequence:  1004901,
					},
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
					Data: types.PackedBytes{
						"Acknowledgement": "eyJyZXN1bHQiOiJleUpqYjI1MGNtRmpkRjl5WlhOMWJIUWlPbTUxYkd3c0ltbGlZMTloWTJzaU9pSmxlVXA1V2xoT01XSklVV2xQYVVwQ1ZWUXdPVWx1TUQwaWZRPT0ifQ==",
						"Packet": map[string]any{
							"Data": transferTypes.FungibleTokenPacketData{
								Amount:   "12000000",
								Denom:    "utia",
								Receiver: "osmo1vkdakqqg5htq5c3wy2kj2geq536q665xdexrtjuwqckpads2c2nsvhhcyv",
								Sender:   "celestia1nc44cmtgmp6cwch2ccfp6txdelu6qtz963ksrw",
								Memo:     "{\"wasm\":{\"contract\":\"osmo1vkdakqqg5htq5c3wy2kj2geq536q665xdexrtjuwqckpads2c2nsvhhcyv\",\"msg\":{\"swap_and_action\":{\"user_swap\":{\"swap_exact_asset_in\":{\"swap_venue_name\":\"osmosis-poolmanager\",\"operations\":[{\"pool\":\"1247\",\"denom_in\":\"ibc/D79E7D83AB399BFFF93433E54FAA480C191248FC556924A2A8351AE2638B3877\",\"denom_out\":\"ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4\"},{\"pool\":\"1319\",\"denom_in\":\"ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4\",\"denom_out\":\"ibc/64BA6E31FE887D66C6F8F31C7B1A80C7CA179239677B4088BB55F5EA07DBE273\"}]}},\"min_asset\":{\"native\":{\"denom\":\"ibc/64BA6E31FE887D66C6F8F31C7B1A80C7CA179239677B4088BB55F5EA07DBE273\",\"amount\":\"3298454343374009626\"}},\"timeout_timestamp\":1726667021008444464,\"post_swap_action\":{\"ibc_transfer\":{\"ibc_info\":{\"source_channel\":\"channel-122\",\"receiver\":\"inj14amztqem07qvyyty8k4p6s2jp8ylsvlax0g42f\",\"memo\":\"\",\"recover_address\":\"osmo1nc44cmtgmp6cwch2ccfp6txdelu6qtz9rq5s03\"}}},\"affiliates\":[{\"basis_points_fee\":\"60\",\"address\":\"osmo1my4tk420gjmhggqwvvha6ey9390gqwfree2p4u\"},{\"basis_points_fee\":\"15\",\"address\":\"osmo1msjnal2glfz6ze8x9kduhg45xppx9sddawfx46\"}]}}}}",
							},
							"DestinationChannel": "channel-6994",
							"DestinationPort":    "transfer",
							"Sequence":           1004900,
							"SourceChannel":      "channel-2",
							"SourcePort":         "transfer",
							"TimeoutHeight": map[string]any{
								"RevisionHeight": 21043368,
								"RevisionNumber": 1,
							},
							"TimeoutTimestamp": 0,
						},
						"ProofAcked": "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",
						"ProofHeight": map[string]any{
							"RevisionHeight": 21043234,
							"RevisionNumber": 1,
						},
						"Signer": "celestia16p6lrlxf7f03c0ka8cv4sznr29rym27us7u4v6",
					},
					IbcTransfer: &storage.IbcTransfer{
						Height:          2371609,
						Amount:          decimal.RequireFromString("12000000"),
						Denom:           "utia",
						ReceiverAddress: testsuite.Ptr("osmo1vkdakqqg5htq5c3wy2kj2geq536q665xdexrtjuwqckpads2c2nsvhhcyv"),
						Sender: &storage.Address{
							Address: "celestia1nc44cmtgmp6cwch2ccfp6txdelu6qtz963ksrw",
						},
						ChannelId: "channel-6994",
						Port:      "transfer",
						Sequence:  1004900,
						Memo:      "{\"wasm\":{\"contract\":\"osmo1vkdakqqg5htq5c3wy2kj2geq536q665xdexrtjuwqckpads2c2nsvhhcyv\",\"msg\":{\"swap_and_action\":{\"user_swap\":{\"swap_exact_asset_in\":{\"swap_venue_name\":\"osmosis-poolmanager\",\"operations\":[{\"pool\":\"1247\",\"denom_in\":\"ibc/D79E7D83AB399BFFF93433E54FAA480C191248FC556924A2A8351AE2638B3877\",\"denom_out\":\"ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4\"},{\"pool\":\"1319\",\"denom_in\":\"ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4\",\"denom_out\":\"ibc/64BA6E31FE887D66C6F8F31C7B1A80C7CA179239677B4088BB55F5EA07DBE273\"}]}},\"min_asset\":{\"native\":{\"denom\":\"ibc/64BA6E31FE887D66C6F8F31C7B1A80C7CA179239677B4088BB55F5EA07DBE273\",\"amount\":\"3298454343374009626\"}},\"timeout_timestamp\":1726667021008444464,\"post_swap_action\":{\"ibc_transfer\":{\"ibc_info\":{\"source_channel\":\"channel-122\",\"receiver\":\"inj14amztqem07qvyyty8k4p6s2jp8ylsvlax0g42f\",\"memo\":\"\",\"recover_address\":\"osmo1nc44cmtgmp6cwch2ccfp6txdelu6qtz9rq5s03\"}}},\"affiliates\":[{\"basis_points_fee\":\"60\",\"address\":\"osmo1my4tk420gjmhggqwvvha6ey9390gqwfree2p4u\"},{\"basis_points_fee\":\"15\",\"address\":\"osmo1msjnal2glfz6ze8x9kduhg45xppx9sddawfx46\"}]}}}}",
					},
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
