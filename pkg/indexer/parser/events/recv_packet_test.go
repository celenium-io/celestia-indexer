// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"

	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icaTypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	transferTypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"

	"github.com/stretchr/testify/require"
)

func Test_handleRecvPacket(t *testing.T) {
	tests := []struct {
		name   string
		ctx    *context.Context
		events []storage.Event
		msg    []*storage.Message
		idx    *int
	}{
		{
			name: "recv packet test 1",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 1866988,
					Type:   "message",
					Data: map[string]any{
						"action": "/ibc.core.channel.v1.MsgRecvPacket",
					},
				}, {
					Height: 1866988,
					Type:   "recv_packet",
					Data: map[string]any{
						"packet_channel_ordering":  "ORDER_ORDERED",
						"packet_connection":        "connection-3",
						"packet_data":              "{\"data\":\"CqIBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoEBCkNjZWxlc3RpYTF4enNkbjY1aHlsamNtZW5seHlqbWRtdmdoaGQwdzR1dDI3azNmeDU2anAycDY5ZWg2c3JzOHAzcnNzEi9jZWxlc3RpYTFkc21ueno4cXZuMzQzbHhxMmV3NjczOXNkMDJnc3NweWVzYTRqdRoJCgR1dGlhEgEx\",\"memo\":\"rootulp\",\"type\":\"TYPE_EXECUTE_TX\"}",
						"packet_data_hex":          "7b2264617461223a2243714942436877765932397a6257397a4c6d4a68626d7375646a46695a5852684d53354e633264545a57356b456f4542436b4e6a5a57786c6333527059544634656e4e6b626a593161486c73616d4e745a57357365486c7162575274646d646f61475177647a523164444933617a4e6d65445532616e4179634459355a57673263334a7a4f48417a636e4e7a4569396a5a57786c633352705954466b63323175656e6f3463585a754d7a517a624868784d6d56334e6a637a4f584e6b4d444a6e63334e776557567a5954527164526f4a4367523164476c6845674578222c226d656d6f223a22726f6f74756c70222c2274797065223a22545950455f455845435554455f5458227d",
						"packet_dst_channel":       "channel-2",
						"packet_dst_port":          "icahost",
						"packet_sequence":          "3",
						"packet_src_channel":       "channel-4311",
						"packet_src_port":          "icacontroller-cosmos1epqzuh6myrwrp4zr8zjamcye4nvkkg9xd8ywak",
						"packet_timeout_height":    "0-0",
						"packet_timeout_timestamp": "1725383082324295575",
					},
				}, {
					Height: 1866988,
					Type:   "message",
					Data: map[string]any{
						"module": "ibc_channel",
					},
				}, {
					Height: 1866988,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "1utia",
						"spender": "celestia1xzsdn65hyljcmenlxyjmdmvghhd0w4ut27k3fx56jp2p69eh6srs8p3rss",
					},
				}, {
					Height: 1866988,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "1utia",
						"receiver": "celestia1dsmnzz8qvn343lxq2ew6739sd02gsspyesa4ju",
					},
				}, {
					Height: 1866988,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "1utia",
						"recipient": "celestia1dsmnzz8qvn343lxq2ew6739sd02gsspyesa4ju",
						"sender":    "celestia1xzsdn65hyljcmenlxyjmdmvghhd0w4ut27k3fx56jp2p69eh6srs8p3rss",
					},
				}, {
					Height: 1866988,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1xzsdn65hyljcmenlxyjmdmvghhd0w4ut27k3fx56jp2p69eh6srs8p3rss",
					},
				}, {
					Height: 1866988,
					Type:   "message",
					Data: map[string]any{
						"module": "bank",
					},
				}, {
					Height: 1866988,
					Type:   "ics27_packet",
					Data: map[string]any{
						"host_channel_id": "channel-2",
						"module":          "interchainaccounts",
						"success":         "true",
					},
				}, {
					Height: 1866988,
					Type:   "write_acknowledgement",
					Data: map[string]any{
						"packet_ack":               "{\"result\":\"EiYKJC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmRSZXNwb25zZQ==\"}",
						"packet_ack_hex":           "7b22726573756c74223a224569594b4a43396a62334e7462334d75596d4675617935324d574a6c644745784c6b317a5a314e6c626d52535a584e776232357a5a513d3d227d",
						"packet_connection":        "connection-3",
						"packet_data":              "{\"data\":\"CqIBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoEBCkNjZWxlc3RpYTF4enNkbjY1aHlsamNtZW5seHlqbWRtdmdoaGQwdzR1dDI3azNmeDU2anAycDY5ZWg2c3JzOHAzcnNzEi9jZWxlc3RpYTFkc21ueno4cXZuMzQzbHhxMmV3NjczOXNkMDJnc3NweWVzYTRqdRoJCgR1dGlhEgEx\",\"memo\":\"rootulp\",\"type\":\"TYPE_EXECUTE_TX\"}",
						"packet_data_hex":          "7b2264617461223a2243714942436877765932397a6257397a4c6d4a68626d7375646a46695a5852684d53354e633264545a57356b456f4542436b4e6a5a57786c6333527059544634656e4e6b626a593161486c73616d4e745a57357365486c7162575274646d646f61475177647a523164444933617a4e6d65445532616e4179634459355a57673263334a7a4f48417a636e4e7a4569396a5a57786c633352705954466b63323175656e6f3463585a754d7a517a624868784d6d56334e6a637a4f584e6b4d444a6e63334e776557567a5954527164526f4a4367523164476c6845674578222c226d656d6f223a22726f6f74756c70222c2274797065223a22545950455f455845435554455f5458227d",
						"packet_dst_channel":       "channel-2",
						"packet_dst_port":          "icahost",
						"packet_sequence":          "3",
						"packet_src_channel":       "channel-4311",
						"packet_src_port":          "icacontroller-cosmos1epqzuh6myrwrp4zr8zjamcye4nvkkg9xd8ywak",
						"packet_timeout_height":    "0-0",
						"packet_timeout_timestamp": "1725383082324295575",
					},
				},
			},
			msg: []*storage.Message{
				{
					Type:   types.MsgRecvPacket,
					Height: 1866988,
					Data: map[string]any{
						"Packet": map[string]any{
							"Data": map[string]any{
								"Memo": "rootulp",
								"Type": icaTypes.EXECUTE_TX,
								"Data": []cosmosTypes.Msg{
									&cosmosBankTypes.MsgSend{
										FromAddress: "celestia1xzsdn65hyljcmenlxyjmdmvghhd0w4ut27k3fx56jp2p69eh6srs8p3rss",
										ToAddress:   "celestia1dsmnzz8qvn343lxq2ew6739sd02gsspyesa4ju",
										Amount:      cosmosTypes.NewCoins(cosmosTypes.NewCoin("utia", math.NewInt(1))),
									},
								},
							},
							"DestinationChannel": "channel-2",
							"DestinationPort":    "icahost",
							"Sequence":           3,
							"SourceChannel":      "channel-4311",
							"SourcePort":         "icacontroller-cosmos1epqzuh6myrwrp4zr8zjamcye4nvkkg9xd8ywak",
							"TimeoutHeight": map[string]any{
								"RevisionHeight": 0,
								"RevisionNumber": 0,
							},
							"TimeoutTimestamp": 1725383082324295700,
						},
						"ProofCommitment": "Cr8JCrwJCm9jb21taXRtZW50cy9wb3J0cy9pY2Fjb250cm9sbGVyLWNvc21vczFlcHF6dWg2bXlyd3JwNHpyOHpqYW1jeWU0bnZra2c5eGQ4eXdhay9jaGFubmVscy9jaGFubmVsLTQzMTEvc2VxdWVuY2VzLzMSIDOCpIscr0NGtsrg1HHkMhvGNI7gE5dlR5qpjuymB+A/Gg4IARgBIAEqBgAC6oeoFiIsCAESKAIE6oeoFiC2yEQJEJWHquHWhg/shpu6fOhyTtt2Jrf90zLAwr0UCyAiLAgBEigEBuqHqBYglv6DW7Udd8HWnGac8Tqmn2XL7BK/ab9FC8SERVGMq9AgIiwIARIoBg7qh6gWIK1Vn+IslEiRV+rjuwsUEytK3cQLJyOMaic6y/OeLjP1ICIsCAESKAgW6oeoFiAkf3L0kNPOb3iWG94x1Oo3F7tBbhTIyAFrzQi+pt6rTiAiLAgBEigKKOqHqBYgTaZg3a6jUz0ZxoCGVMv5Ms5Gi6NPmJMb9dAa2fn+Q6UgIi4IARIHDGjqh6gWIBohIJaGaKlZh0VVe2ssuilbDdCi3a0SiB30NGGpltGQmeA4Ii0IARIpDrYB6oeoFiBgbasOp9FmZSOJD++feygAcJYqoaRUFfkzq7ajJQ3LuCAiLQgBEikQpgLqh6gWIAl0SSkvpQjTDxRVrn1CfBfh87LLuW8xmBWLXpOQjt7NICItCAESKRL2BOqHqBYg4MgElmhPULuGOedxNZoAQp1FFnsbG/3yrTPYl4WZa0QgIi8IARIIFPYI6oeoFiAaISAuXh/nYY9vlfQKv/CgyUrPFzhycY1gk3Jw7bqTwF/rMiItCAESKRb4D+qHqBYg2+Rbd6aRYQmx64VbkpBNZ5tTm6ZFoJxSbXhNG1cv8dAgIi0IARIpGNYX6oeoFiCHjG3nSixO/bAilis8FCYwd/EWN9KK7ord/qD8o4JcqCAiLQgBEikatCvqh6gWIP1U5ibnw5lnxJXnEgEF+Sezp3ZOfOd5I46hwrtR2qPWICIvCAESCBy0QeqHqBYgGiEgxVJIDrh0mo4Jo2vnRhShF3Nijplat7z0LF7jnWjbYLMiLggBEioeppYB6oeoFiDyiezkU1qbVkDwyurADIjsoWY/eeML9hW52bHbOWAi+yAiMAgBEgkgrJcC6oeoFiAaISCVF4HnuKW911CS/1Z7RVtH9F3CxNCTe0UL/0JuGijzmSIuCAESKiLq3QPqh6gWIOoqmcYC8BjIhzdpVhEecmVjSEJMhkgBxPHPOYd12zckICIwCAESCSSQoQ/qh6gWIBohIP3AOOA+ESwqm09lXOQEKThWYKYEVChNqPZ5mrO77+HlIi4IARIqJoD3GuqHqBYghG3zHblcVrp+v9Axn2sLv42ZvZ45A7yqeAMLQGEl4F0gIi4IARIqKqzlP+qHqBYg1nEeJj/pEy3BJmJtfQaQT436/rjYi37b+hzYhpjY2r0gIjAIARIJLJDGaeqHqBYgGiEgZ9AeEXfcGGYWs8GFUhtqJjo0hrERT4aQ/kHfqdEEgeoiLwgBEiswvtL0AeqHqBYgGI0c4X+NGJy3R41b7tcFaGe13rFbrOzGuri1P3lp18ogCv4BCvsBCgNpYmMSIJjgI1oJUJC9HROp6RPN2KPTilScZ5Xt5JcAair7ud1YGgkIARgBIAEqAQAiJwgBEgEBGiAsZZMul3Pb18hNf7k3waWY0EPHQC4zEBRXGn11UjB3HCInCAESAQEaIJ6mHb7R8tnlGSTGyBulCjKkS1+sYn2nyOJzwUMa2l9LIiUIARIhAedBsHDxHKeqTnsAdphrBaZD+RqerkYmcU0ZNeAbwDnBIiUIARIhAQb01cMqsbVA5Evp7ZOqln0OG8rkQbPrIC2usl+UYjILIicIARIBARogcy34IRcmnnh/83TNtLTUER4Hm5ZFQJoPMOibHPR0Jlg=",
						"ProofHeight": map[string]any{
							"RevisionHeight": 23396854,
							"RevisionNumber": 0,
						},
						"Signer": "celestia1dsmnzz8qvn343lxq2ew6739sd02gsspyesa4ju",
					},
				},
			},
			idx: testsuite.Ptr(0),
		}, {
			name: "recv packet test 2",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 1866988,
					Type:   "message",
					Data: map[string]any{
						"action": "/ibc.core.channel.v1.MsgRecvPacket",
					},
				}, {
					Height: 1866988,
					Type:   "recv_packet",
					Data: map[string]any{
						"packet_channel_ordering":  "ORDER_UNORDERED",
						"packet_connection":        "connection-7",
						"packet_data":              "{\"amount\":\"4000000\",\"denom\":\"transfer/channel-35/utia\",\"memo\":\"{\\\"forward\\\":{\\\"receiver\\\":\\\"osmo1vkdakqqg5htq5c3wy2kj2geq536q665xdexrtjuwqckpads2c2nsvhhcyv\\\",\\\"port\\\":\\\"transfer\\\",\\\"channel\\\":\\\"channel-2\\\",\\\"timeout\\\":0,\\\"retries\\\":2,\\\"next\\\":{\\\"wasm\\\":{\\\"contract\\\":\\\"osmo1vkdakqqg5htq5c3wy2kj2geq536q665xdexrtjuwqckpads2c2nsvhhcyv\\\",\\\"msg\\\":{\\\"swap_and_action\\\":{\\\"user_swap\\\":{\\\"swap_exact_asset_in\\\":{\\\"swap_venue_name\\\":\\\"osmosis-poolmanager\\\",\\\"operations\\\":[{\\\"pool\\\":\\\"1475\\\",\\\"denom_in\\\":\\\"ibc/D79E7D83AB399BFFF93433E54FAA480C191248FC556924A2A8351AE2638B3877\\\",\\\"denom_out\\\":\\\"factory/osmo1f5vfcph2dvfeqcqkhetwv75fda69z7e5c2dldm3kvgj23crkv6wqcn47a0/umilkTIA\\\"},{\\\"pool\\\":\\\"1694\\\",\\\"denom_in\\\":\\\"factory/osmo1f5vfcph2dvfeqcqkhetwv75fda69z7e5c2dldm3kvgj23crkv6wqcn47a0/umilkTIA\\\",\\\"denom_out\\\":\\\"ibc/69110FF673D70B39904FF056CFDFD58A90BEC3194303F45C32CB91B8B0A738EA\\\"},{\\\"pool\\\":\\\"1698\\\",\\\"denom_in\\\":\\\"ibc/69110FF673D70B39904FF056CFDFD58A90BEC3194303F45C32CB91B8B0A738EA\\\",\\\"denom_out\\\":\\\"ibc/64BA6E31FE887D66C6F8F31C7B1A80C7CA179239677B4088BB55F5EA07DBE273\\\"}]}},\\\"min_asset\\\":{\\\"native\\\":{\\\"denom\\\":\\\"ibc/64BA6E31FE887D66C6F8F31C7B1A80C7CA179239677B4088BB55F5EA07DBE273\\\",\\\"amount\\\":\\\"1095116599248063092\\\"}},\\\"timeout_timestamp\\\":1726673962233433067,\\\"post_swap_action\\\":{\\\"ibc_transfer\\\":{\\\"ibc_info\\\":{\\\"source_channel\\\":\\\"channel-122\\\",\\\"receiver\\\":\\\"inj1syrdh2v2rwf8fhvs5gsmwr4ahcns3m4zvcukwc\\\",\\\"memo\\\":\\\"\\\",\\\"recover_address\\\":\\\"osmo1m84nh75hl474k5d83cunrqfxgmrl523u5q09xw\\\"}}},\\\"affiliates\\\":[{\\\"basis_points_fee\\\":\\\"60\\\",\\\"address\\\":\\\"osmo1my4tk420gjmhggqwvvha6ey9390gqwfree2p4u\\\"},{\\\"basis_points_fee\\\":\\\"15\\\",\\\"address\\\":\\\"osmo1msjnal2glfz6ze8x9kduhg45xppx9sddawfx46\\\"}]}}}}}}\",\"receiver\":\"celestia1m84nh75hl474k5d83cunrqfxgmrl523ud3d923\",\"sender\":\"neutron1m84nh75hl474k5d83cunrqfxgmrl523ucy4h2m\"}",
						"packet_data_hex":          "7b22616d6f756e74223a2234303030303030222c2264656e6f6d223a227472616e736665722f6368616e6e656c2d33352f75746961222c226d656d6f223a227b5c22666f72776172645c223a7b5c2272656365697665725c223a5c226f736d6f31766b64616b717167356874713563337779326b6a32676571353336713636357864657872746a757771636b706164733263326e737668686379765c222c5c22706f72745c223a5c227472616e736665725c222c5c226368616e6e656c5c223a5c226368616e6e656c2d325c222c5c2274696d656f75745c223a302c5c22726574726965735c223a322c5c226e6578745c223a7b5c227761736d5c223a7b5c22636f6e74726163745c223a5c226f736d6f31766b64616b717167356874713563337779326b6a32676571353336713636357864657872746a757771636b706164733263326e737668686379765c222c5c226d73675c223a7b5c22737761705f616e645f616374696f6e5c223a7b5c22757365725f737761705c223a7b5c22737761705f65786163745f61737365745f696e5c223a7b5c22737761705f76656e75655f6e616d655c223a5c226f736d6f7369732d706f6f6c6d616e616765725c222c5c226f7065726174696f6e735c223a5b7b5c22706f6f6c5c223a5c22313437355c222c5c2264656e6f6d5f696e5c223a5c226962632f443739453744383341423339394246464639333433334535344641413438304331393132343846433535363932344132413833353141453236333842333837375c222c5c2264656e6f6d5f6f75745c223a5c22666163746f72792f6f736d6f316635766663706832647666657163716b6865747776373566646136397a3765356332646c646d336b76676a323363726b76367771636e343761302f756d696c6b5449415c227d2c7b5c22706f6f6c5c223a5c22313639345c222c5c2264656e6f6d5f696e5c223a5c22666163746f72792f6f736d6f316635766663706832647666657163716b6865747776373566646136397a3765356332646c646d336b76676a323363726b76367771636e343761302f756d696c6b5449415c222c5c2264656e6f6d5f6f75745c223a5c226962632f363931313046463637334437304233393930344646303536434644464435384139304245433331393433303346343543333243423931423842304137333845415c227d2c7b5c22706f6f6c5c223a5c22313639385c222c5c2264656e6f6d5f696e5c223a5c226962632f363931313046463637334437304233393930344646303536434644464435384139304245433331393433303346343543333243423931423842304137333845415c222c5c2264656e6f6d5f6f75745c223a5c226962632f363442413645333146453838374436364336463846333143374231413830433743413137393233393637374234303838424235354635454130374442453237335c227d5d7d7d2c5c226d696e5f61737365745c223a7b5c226e61746976655c223a7b5c2264656e6f6d5c223a5c226962632f363442413645333146453838374436364336463846333143374231413830433743413137393233393637374234303838424235354635454130374442453237335c222c5c22616d6f756e745c223a5c22313039353131363539393234383036333039325c227d7d2c5c2274696d656f75745f74696d657374616d705c223a313732363637333936323233333433333036372c5c22706f73745f737761705f616374696f6e5c223a7b5c226962635f7472616e736665725c223a7b5c226962635f696e666f5c223a7b5c22736f757263655f6368616e6e656c5c223a5c226368616e6e656c2d3132325c222c5c2272656365697665725c223a5c22696e6a31737972646832763272776638666876733567736d7772346168636e73336d347a7663756b77635c222c5c226d656d6f5c223a5c225c222c5c227265636f7665725f616464726573735c223a5c226f736d6f316d38346e683735686c3437346b3564383363756e72716678676d726c353233753571303978775c227d7d7d2c5c22616666696c69617465735c223a5b7b5c2262617369735f706f696e74735f6665655c223a5c2236305c222c5c22616464726573735c223a5c226f736d6f316d7934746b343230676a6d6867677177767668613665793933393067717766726565327034755c227d2c7b5c2262617369735f706f696e74735f6665655c223a5c2231355c222c5c22616464726573735c223a5c226f736d6f316d736a6e616c32676c667a367a653878396b64756867343578707078397364646177667834365c227d5d7d7d7d7d7d7d222c227265636569766572223a2263656c6573746961316d38346e683735686c3437346b3564383363756e72716678676d726c35323375643364393233222c2273656e646572223a226e657574726f6e316d38346e683735686c3437346b3564383363756e72716678676d726c3532337563793468326d227d",
						"packet_dst_channel":       "channel-8",
						"packet_dst_port":          "transfer",
						"packet_sequence":          "838308",
						"packet_src_channel":       "channel-35",
						"packet_src_port":          "transfer",
						"packet_timeout_height":    "0-2372332",
						"packet_timeout_timestamp": "0",
					},
				}, {
					Height: 1866988,
					Type:   "message",
					Data: map[string]any{
						"module": "ibc_channel",
					},
				}, {
					Height: 1866988,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "4000000utia",
						"spender": "celestia187dz9zlxc3zrltzx5756tu7zew6yu3v0smnfem",
					},
				}, {
					Height: 1866988,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "4000000utia",
						"receiver": "celestia1wzpmvah0aqy7cq948tswayvvstty2djtjm88ee",
					},
				}, {
					Height: 1866988,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "4000000utia",
						"recipient": "celestia1wzpmvah0aqy7cq948tswayvvstty2djtjm88ee",
						"sender":    "celestia187dz9zlxc3zrltzx5756tu7zew6yu3v0smnfem",
					},
				}, {
					Height: 1866988,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia187dz9zlxc3zrltzx5756tu7zew6yu3v0smnfem",
					},
				}, {
					Height: 1866988,
					Type:   "fungible_token_packet",
					Data: map[string]any{
						"amount":   "4000000",
						"denom":    "transfer/channel-35/utia",
						"memo":     "",
						"module":   "transfer",
						"receiver": "celestia1wzpmvah0aqy7cq948tswayvvstty2djtjm88ee",
						"sender":   "neutron1m84nh75hl474k5d83cunrqfxgmrl523ucy4h2m",
						"success":  "true",
					},
				}, {
					Height: 1866988,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "4000000utia",
						"spender": "celestia1wzpmvah0aqy7cq948tswayvvstty2djtjm88ee",
					},
				}, {
					Height: 1866988,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "4000000utia",
						"receiver": "celestia12k2pyuylm9t7ugdvz67h9pg4gmmvhn5vwv5zte",
					},
				}, {
					Height: 1866988,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "4000000utia",
						"recipient": "celestia12k2pyuylm9t7ugdvz67h9pg4gmmvhn5vwv5zte",
						"sender":    "celestia1wzpmvah0aqy7cq948tswayvvstty2djtjm88ee",
					},
				}, {
					Height: 1866988,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1wzpmvah0aqy7cq948tswayvvstty2djtjm88ee",
					},
				}, {
					Height: 1866988,
					Type:   "send_packet",
					Data: map[string]any{
						"packet_channel_ordering":  "ORDER_UNORDERED",
						"packet_connection":        "connection-2",
						"packet_data":              "{\"amount\":\"4000000\",\"denom\":\"utia\",\"memo\":\"{\\\"wasm\\\":{\\\"contract\\\":\\\"osmo1vkdakqqg5htq5c3wy2kj2geq536q665xdexrtjuwqckpads2c2nsvhhcyv\\\",\\\"msg\\\":{\\\"swap_and_action\\\":{\\\"user_swap\\\":{\\\"swap_exact_asset_in\\\":{\\\"swap_venue_name\\\":\\\"osmosis-poolmanager\\\",\\\"operations\\\":[{\\\"pool\\\":\\\"1475\\\",\\\"denom_in\\\":\\\"ibc/D79E7D83AB399BFFF93433E54FAA480C191248FC556924A2A8351AE2638B3877\\\",\\\"denom_out\\\":\\\"factory/osmo1f5vfcph2dvfeqcqkhetwv75fda69z7e5c2dldm3kvgj23crkv6wqcn47a0/umilkTIA\\\"},{\\\"pool\\\":\\\"1694\\\",\\\"denom_in\\\":\\\"factory/osmo1f5vfcph2dvfeqcqkhetwv75fda69z7e5c2dldm3kvgj23crkv6wqcn47a0/umilkTIA\\\",\\\"denom_out\\\":\\\"ibc/69110FF673D70B39904FF056CFDFD58A90BEC3194303F45C32CB91B8B0A738EA\\\"},{\\\"pool\\\":\\\"1698\\\",\\\"denom_in\\\":\\\"ibc/69110FF673D70B39904FF056CFDFD58A90BEC3194303F45C32CB91B8B0A738EA\\\",\\\"denom_out\\\":\\\"ibc/64BA6E31FE887D66C6F8F31C7B1A80C7CA179239677B4088BB55F5EA07DBE273\\\"}]}},\\\"min_asset\\\":{\\\"native\\\":{\\\"denom\\\":\\\"ibc/64BA6E31FE887D66C6F8F31C7B1A80C7CA179239677B4088BB55F5EA07DBE273\\\",\\\"amount\\\":\\\"1095116599248063092\\\"}},\\\"timeout_timestamp\\\":1726673962233433000,\\\"post_swap_action\\\":{\\\"ibc_transfer\\\":{\\\"ibc_info\\\":{\\\"source_channel\\\":\\\"channel-122\\\",\\\"receiver\\\":\\\"inj1syrdh2v2rwf8fhvs5gsmwr4ahcns3m4zvcukwc\\\",\\\"memo\\\":\\\"\\\",\\\"recover_address\\\":\\\"osmo1m84nh75hl474k5d83cunrqfxgmrl523u5q09xw\\\"}}},\\\"affiliates\\\":[{\\\"basis_points_fee\\\":\\\"60\\\",\\\"address\\\":\\\"osmo1my4tk420gjmhggqwvvha6ey9390gqwfree2p4u\\\"},{\\\"basis_points_fee\\\":\\\"15\\\",\\\"address\\\":\\\"osmo1msjnal2glfz6ze8x9kduhg45xppx9sddawfx46\\\"}]}}}}\",\"receiver\":\"osmo1vkdakqqg5htq5c3wy2kj2geq536q665xdexrtjuwqckpads2c2nsvhhcyv\",\"sender\":\"celestia1wzpmvah0aqy7cq948tswayvvstty2djtjm88ee\"}",
						"packet_data_hex":          "7b22616d6f756e74223a2234303030303030222c2264656e6f6d223a2275746961222c226d656d6f223a227b5c227761736d5c223a7b5c22636f6e74726163745c223a5c226f736d6f31766b64616b717167356874713563337779326b6a32676571353336713636357864657872746a757771636b706164733263326e737668686379765c222c5c226d73675c223a7b5c22737761705f616e645f616374696f6e5c223a7b5c22757365725f737761705c223a7b5c22737761705f65786163745f61737365745f696e5c223a7b5c22737761705f76656e75655f6e616d655c223a5c226f736d6f7369732d706f6f6c6d616e616765725c222c5c226f7065726174696f6e735c223a5b7b5c22706f6f6c5c223a5c22313437355c222c5c2264656e6f6d5f696e5c223a5c226962632f443739453744383341423339394246464639333433334535344641413438304331393132343846433535363932344132413833353141453236333842333837375c222c5c2264656e6f6d5f6f75745c223a5c22666163746f72792f6f736d6f316635766663706832647666657163716b6865747776373566646136397a3765356332646c646d336b76676a323363726b76367771636e343761302f756d696c6b5449415c227d2c7b5c22706f6f6c5c223a5c22313639345c222c5c2264656e6f6d5f696e5c223a5c22666163746f72792f6f736d6f316635766663706832647666657163716b6865747776373566646136397a3765356332646c646d336b76676a323363726b76367771636e343761302f756d696c6b5449415c222c5c2264656e6f6d5f6f75745c223a5c226962632f363931313046463637334437304233393930344646303536434644464435384139304245433331393433303346343543333243423931423842304137333845415c227d2c7b5c22706f6f6c5c223a5c22313639385c222c5c2264656e6f6d5f696e5c223a5c226962632f363931313046463637334437304233393930344646303536434644464435384139304245433331393433303346343543333243423931423842304137333845415c222c5c2264656e6f6d5f6f75745c223a5c226962632f363442413645333146453838374436364336463846333143374231413830433743413137393233393637374234303838424235354635454130374442453237335c227d5d7d7d2c5c226d696e5f61737365745c223a7b5c226e61746976655c223a7b5c2264656e6f6d5c223a5c226962632f363442413645333146453838374436364336463846333143374231413830433743413137393233393637374234303838424235354635454130374442453237335c222c5c22616d6f756e745c223a5c22313039353131363539393234383036333039325c227d7d2c5c2274696d656f75745f74696d657374616d705c223a313732363637333936323233333433333030302c5c22706f73745f737761705f616374696f6e5c223a7b5c226962635f7472616e736665725c223a7b5c226962635f696e666f5c223a7b5c22736f757263655f6368616e6e656c5c223a5c226368616e6e656c2d3132325c222c5c2272656365697665725c223a5c22696e6a31737972646832763272776638666876733567736d7772346168636e73336d347a7663756b77635c222c5c226d656d6f5c223a5c225c222c5c227265636f7665725f616464726573735c223a5c226f736d6f316d38346e683735686c3437346b3564383363756e72716678676d726c353233753571303978775c227d7d7d2c5c22616666696c69617465735c223a5b7b5c2262617369735f706f696e74735f6665655c223a5c2236305c222c5c22616464726573735c223a5c226f736d6f316d7934746b343230676a6d6867677177767668613665793933393067717766726565327034755c227d2c7b5c2262617369735f706f696e74735f6665655c223a5c2231355c222c5c22616464726573735c223a5c226f736d6f316d736a6e616c32676c667a367a653878396b64756867343578707078397364646177667834365c227d5d7d7d7d7d222c227265636569766572223a226f736d6f31766b64616b717167356874713563337779326b6a32676571353336713636357864657872746a757771636b706164733263326e73766868637976222c2273656e646572223a2263656c657374696131777a706d76616830617179376371393438747377617976767374747932646a746a6d38386565227d",
						"packet_dst_channel":       "channel-6994",
						"packet_dst_port":          "transfer",
						"packet_sequence":          "1004970",
						"packet_src_channel":       "channel-2",
						"packet_src_port":          "transfer",
						"packet_timeout_height":    "0-0",
						"packet_timeout_timestamp": "1726674267162051887",
					},
				}, {
					Height: 1866988,
					Type:   "message",
					Data: map[string]any{
						"module": "ibc_channel",
					},
				}, {
					Height: 1866988,
					Type:   "ibc_transfer",
					Data: map[string]any{
						"receiver": "osmo1vkdakqqg5htq5c3wy2kj2geq536q665xdexrtjuwqckpads2c2nsvhhcyv",
						"sender":   "celestia1wzpmvah0aqy7cq948tswayvvstty2djtjm88ee",
					},
				}, {
					Height: 1866988,
					Type:   "message",
					Data: map[string]any{
						"module": "transfer",
					},
				}, {
					Height: 1866988,
					Type:   "message",
					Data: map[string]any{
						"action": "/ibc.core.channel.v1.MsgRecvPacket",
					},
				}, {
					Height: 1866988,
					Type:   "recv_packet",
					Data: map[string]any{
						"packet_channel_ordering":  "ORDER_UNORDERED",
						"packet_connection":        "connection-7",
						"packet_data":              "{\"amount\":\"100\",\"denom\":\"transfer/channel-35/utia\",\"receiver\":\"celestia1nsxcgald2c3622hfwflps608tqrj9l3wdcmq9s\",\"sender\":\"neutron1upjaknf6lmnu3p4llldp8jx0whzsxlgetu9zjt\"}",
						"packet_data_hex":          "7b22616d6f756e74223a22313030222c2264656e6f6d223a227472616e736665722f6368616e6e656c2d33352f75746961222c227265636569766572223a2263656c6573746961316e73786367616c64326333363232686677666c70733630387471726a396c337764636d713973222c2273656e646572223a226e657574726f6e3175706a616b6e66366c6d6e753370346c6c6c6470386a783077687a73786c67657475397a6a74227d",
						"packet_dst_channel":       "channel-8",
						"packet_dst_port":          "transfer",
						"packet_sequence":          "838309",
						"packet_src_channel":       "channel-35",
						"packet_src_port":          "transfer",
						"packet_timeout_height":    "0-0",
						"packet_timeout_timestamp": "1726674263466000000",
					},
				}, {
					Height: 1866988,
					Type:   "message",
					Data: map[string]any{
						"module": "ibc_channel",
					},
				}, {
					Height: 1866988,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "100utia",
						"spender": "celestia187dz9zlxc3zrltzx5756tu7zew6yu3v0smnfem",
					},
				}, {
					Height: 1866988,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "100utia",
						"receiver": "celestia1nsxcgald2c3622hfwflps608tqrj9l3wdcmq9s",
					},
				}, {
					Height: 1866988,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "100utia",
						"recipient": "celestia1nsxcgald2c3622hfwflps608tqrj9l3wdcmq9s",
						"sender":    "celestia187dz9zlxc3zrltzx5756tu7zew6yu3v0smnfem",
					},
				}, {
					Height: 1866988,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia187dz9zlxc3zrltzx5756tu7zew6yu3v0smnfem",
					},
				}, {
					Height: 1866988,
					Type:   "fungible_token_packet",
					Data: map[string]any{
						"amount":   "100",
						"denom":    "transfer/channel-35/utia",
						"memo":     "",
						"module":   "transfer",
						"receiver": "celestia1nsxcgald2c3622hfwflps608tqrj9l3wdcmq9s",
						"sender":   "neutron1upjaknf6lmnu3p4llldp8jx0whzsxlgetu9zjt",
						"success":  "true",
					},
				}, {
					Height: 1866988,
					Type:   "write_acknowledgement",
					Data: map[string]any{
						"packet_ack":               "{\"result\":\"AQ==\"}",
						"packet_ack_hex":           "7b22726573756c74223a2241513d3d227d",
						"packet_connection":        "connection-7",
						"packet_data":              "{\"amount\":\"100\",\"denom\":\"transfer/channel-35/utia\",\"receiver\":\"celestia1nsxcgald2c3622hfwflps608tqrj9l3wdcmq9s\",\"sender\":\"neutron1upjaknf6lmnu3p4llldp8jx0whzsxlgetu9zjt\"}",
						"packet_data_hex":          "7b22616d6f756e74223a22313030222c2264656e6f6d223a227472616e736665722f6368616e6e656c2d33352f75746961222c227265636569766572223a2263656c6573746961316e73786367616c64326333363232686677666c70733630387471726a396c337764636d713973222c2273656e646572223a226e657574726f6e3175706a616b6e66366c6d6e753370346c6c6c6470386a783077687a73786c67657475397a6a74227d",
						"packet_dst_channel":       "channel-8",
						"packet_dst_port":          "transfer",
						"packet_sequence":          "838309",
						"packet_src_channel":       "channel-35",
						"packet_src_port":          "transfer",
						"packet_timeout_height":    "0-0",
						"packet_timeout_timestamp": "1726674263466000000",
					},
				}, {
					Height: 1866988,
					Type:   "message",
					Data: map[string]any{
						"module": "ibc_channel",
					},
				},
			},
			msg: []*storage.Message{
				{
					Type:   types.MsgRecvPacket,
					Height: 1866988,
					Data: map[string]any{
						"Packet": map[string]any{
							"Data": transferTypes.FungibleTokenPacketData{
								Amount:   "4000000",
								Denom:    "transfer/channel-35/utia",
								Memo:     "{\\\"forward\\\":{\\\"receiver\\\":\\\"osmo1vkdakqqg5htq5c3wy2kj2geq536q665xdexrtjuwqckpads2c2nsvhhcyv\\\",\\\"port\\\":\\\"transfer\\\",\\\"channel\\\":\\\"channel-2\\\",\\\"timeout\\\":0,\\\"retries\\\":2,\\\"next\\\":{\\\"wasm\\\":{\\\"contract\\\":\\\"osmo1vkdakqqg5htq5c3wy2kj2geq536q665xdexrtjuwqckpads2c2nsvhhcyv\\\",\\\"msg\\\":{\\\"swap_and_action\\\":{\\\"user_swap\\\":{\\\"swap_exact_asset_in\\\":{\\\"swap_venue_name\\\":\\\"osmosis-poolmanager\\\",\\\"operations\\\":[{\\\"pool\\\":\\\"1475\\\",\\\"denom_in\\\":\\\"ibc/D79E7D83AB399BFFF93433E54FAA480C191248FC556924A2A8351AE2638B3877\\\",\\\"denom_out\\\":\\\"factory/osmo1f5vfcph2dvfeqcqkhetwv75fda69z7e5c2dldm3kvgj23crkv6wqcn47a0/umilkTIA\\\"},{\\\"pool\\\":\\\"1694\\\",\\\"denom_in\\\":\\\"factory/osmo1f5vfcph2dvfeqcqkhetwv75fda69z7e5c2dldm3kvgj23crkv6wqcn47a0/umilkTIA\\\",\\\"denom_out\\\":\\\"ibc/69110FF673D70B39904FF056CFDFD58A90BEC3194303F45C32CB91B8B0A738EA\\\"},{\\\"pool\\\":\\\"1698\\\",\\\"denom_in\\\":\\\"ibc/69110FF673D70B39904FF056CFDFD58A90BEC3194303F45C32CB91B8B0A738EA\\\",\\\"denom_out\\\":\\\"ibc/64BA6E31FE887D66C6F8F31C7B1A80C7CA179239677B4088BB55F5EA07DBE273\\\"}]}},\\\"min_asset\\\":{\\\"native\\\":{\\\"denom\\\":\\\"ibc/64BA6E31FE887D66C6F8F31C7B1A80C7CA179239677B4088BB55F5EA07DBE273\\\",\\\"amount\\\":\\\"1095116599248063092\\\"}},\\\"timeout_timestamp\\\":1726673962233433067,\\\"post_swap_action\\\":{\\\"ibc_transfer\\\":{\\\"ibc_info\\\":{\\\"source_channel\\\":\\\"channel-122\\\",\\\"receiver\\\":\\\"inj1syrdh2v2rwf8fhvs5gsmwr4ahcns3m4zvcukwc\\\",\\\"memo\\\":\\\"\\\",\\\"recover_address\\\":\\\"osmo1m84nh75hl474k5d83cunrqfxgmrl523u5q09xw\\\"}}},\\\"affiliates\\\":[{\\\"basis_points_fee\\\":\\\"60\\\",\\\"address\\\":\\\"osmo1my4tk420gjmhggqwvvha6ey9390gqwfree2p4u\\\"},{\\\"basis_points_fee\\\":\\\"15\\\",\\\"address\\\":\\\"osmo1msjnal2glfz6ze8x9kduhg45xppx9sddawfx46\\\"}]}}}}}}",
								Receiver: "celestia1m84nh75hl474k5d83cunrqfxgmrl523ud3d923",
								Sender:   "neutron1m84nh75hl474k5d83cunrqfxgmrl523ucy4h2m",
							},
							"DestinationChannel": "channel-8",
							"DestinationPort":    "transfer",
							"Sequence":           838308,
							"SourceChannel":      "channel-35",
							"SourcePort":         "transfer",
							"TimeoutHeight": map[string]any{
								"RevisionHeight": 2372332,
								"RevisionNumber": 0,
							},
							"TimeoutTimestamp": 0,
						},
						"ProofCommitment": "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",
						"ProofHeight": map[string]any{
							"RevisionHeight": 14715481,
							"RevisionNumber": 1,
						},
						"Signer": "celestia1cdlz8scnf3mmxdnf4njmtp7vz4gps7fsm503qe",
					},
					IbcTransfer: &storage.IbcTransfer{
						Height:        1866988,
						SenderAddress: testsuite.Ptr("neutron1m84nh75hl474k5d83cunrqfxgmrl523ucy4h2m"),
						Receiver: &storage.Address{
							Address: "celestia1m84nh75hl474k5d83cunrqfxgmrl523ud3d923",
						},
						Memo:   "{\\\"forward\\\":{\\\"receiver\\\":\\\"osmo1vkdakqqg5htq5c3wy2kj2geq536q665xdexrtjuwqckpads2c2nsvhhcyv\\\",\\\"port\\\":\\\"transfer\\\",\\\"channel\\\":\\\"channel-2\\\",\\\"timeout\\\":0,\\\"retries\\\":2,\\\"next\\\":{\\\"wasm\\\":{\\\"contract\\\":\\\"osmo1vkdakqqg5htq5c3wy2kj2geq536q665xdexrtjuwqckpads2c2nsvhhcyv\\\",\\\"msg\\\":{\\\"swap_and_action\\\":{\\\"user_swap\\\":{\\\"swap_exact_asset_in\\\":{\\\"swap_venue_name\\\":\\\"osmosis-poolmanager\\\",\\\"operations\\\":[{\\\"pool\\\":\\\"1475\\\",\\\"denom_in\\\":\\\"ibc/D79E7D83AB399BFFF93433E54FAA480C191248FC556924A2A8351AE2638B3877\\\",\\\"denom_out\\\":\\\"factory/osmo1f5vfcph2dvfeqcqkhetwv75fda69z7e5c2dldm3kvgj23crkv6wqcn47a0/umilkTIA\\\"},{\\\"pool\\\":\\\"1694\\\",\\\"denom_in\\\":\\\"factory/osmo1f5vfcph2dvfeqcqkhetwv75fda69z7e5c2dldm3kvgj23crkv6wqcn47a0/umilkTIA\\\",\\\"denom_out\\\":\\\"ibc/69110FF673D70B39904FF056CFDFD58A90BEC3194303F45C32CB91B8B0A738EA\\\"},{\\\"pool\\\":\\\"1698\\\",\\\"denom_in\\\":\\\"ibc/69110FF673D70B39904FF056CFDFD58A90BEC3194303F45C32CB91B8B0A738EA\\\",\\\"denom_out\\\":\\\"ibc/64BA6E31FE887D66C6F8F31C7B1A80C7CA179239677B4088BB55F5EA07DBE273\\\"}]}},\\\"min_asset\\\":{\\\"native\\\":{\\\"denom\\\":\\\"ibc/64BA6E31FE887D66C6F8F31C7B1A80C7CA179239677B4088BB55F5EA07DBE273\\\",\\\"amount\\\":\\\"1095116599248063092\\\"}},\\\"timeout_timestamp\\\":1726673962233433067,\\\"post_swap_action\\\":{\\\"ibc_transfer\\\":{\\\"ibc_info\\\":{\\\"source_channel\\\":\\\"channel-122\\\",\\\"receiver\\\":\\\"inj1syrdh2v2rwf8fhvs5gsmwr4ahcns3m4zvcukwc\\\",\\\"memo\\\":\\\"\\\",\\\"recover_address\\\":\\\"osmo1m84nh75hl474k5d83cunrqfxgmrl523u5q09xw\\\"}}},\\\"affiliates\\\":[{\\\"basis_points_fee\\\":\\\"60\\\",\\\"address\\\":\\\"osmo1my4tk420gjmhggqwvvha6ey9390gqwfree2p4u\\\"},{\\\"basis_points_fee\\\":\\\"15\\\",\\\"address\\\":\\\"osmo1msjnal2glfz6ze8x9kduhg45xppx9sddawfx46\\\"}]}}}}}}",
						Amount: decimal.RequireFromString("4000000"),
						Denom:  "utia",
					},
					IbcChannel: &storage.IbcChannel{
						Id:       "channel-8",
						PortId:   "transfer",
						Received: decimal.RequireFromString("4000000"),
					},
				}, {
					Type:   types.MsgRecvPacket,
					Height: 1866988,
					Data: map[string]any{
						"Packet": map[string]any{
							"Data": transferTypes.FungibleTokenPacketData{
								Amount:   "100",
								Denom:    "transfer/channel-35/utia",
								Receiver: "celestia1nsxcgald2c3622hfwflps608tqrj9l3wdcmq9s",
								Sender:   "neutron1upjaknf6lmnu3p4llldp8jx0whzsxlgetu9zjt",
							},
							"DestinationChannel": "channel-8",
							"DestinationPort":    "transfer",
							"Sequence":           838309,
							"SourceChannel":      "channel-35",
							"SourcePort":         "transfer",
							"TimeoutHeight": map[string]any{
								"RevisionHeight": 0,
								"RevisionNumber": 0,
							},
							"TimeoutTimestamp": 1726674263466000000,
						},
						"ProofCommitment": "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",
						"ProofHeight": map[string]any{
							"RevisionHeight": 14715481,
							"RevisionNumber": 1,
						},
						"Signer": "celestia1cdlz8scnf3mmxdnf4njmtp7vz4gps7fsm503qe",
					},
					IbcTransfer: &storage.IbcTransfer{
						Height:        1866988,
						SenderAddress: testsuite.Ptr("neutron1upjaknf6lmnu3p4llldp8jx0whzsxlgetu9zjt"),
						Receiver: &storage.Address{
							Address: "celestia1nsxcgald2c3622hfwflps608tqrj9l3wdcmq9s",
						},
						Amount: decimal.RequireFromString("100"),
						Denom:  "utia",
					},
					IbcChannel: &storage.IbcChannel{
						Id:       "channel-8",
						PortId:   "transfer",
						Received: decimal.RequireFromString("100"),
					},
				},
			},
			idx: testsuite.Ptr(0),
		}, {
			name: "recv packet test 3",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 1866988,
					Type:   "message",
					Data: map[string]any{
						"action": "/ibc.core.channel.v1.MsgRecvPacket",
					},
				}, {
					Height: 1866988,
					Type:   "recv_packet",
					Data: map[string]any{
						"packet_channel_ordering":  "ORDER_ORDERED",
						"packet_connection":        "connection-99",
						"packet_data":              "{\"data\":\"CrMBCiUvY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dVbmRlbGVnYXRlEokBCkNjZWxlc3RpYTFmM2h4N3JkejVxYXN5Z2ZnNTNsZms0cTZsaDRmcGFqdXZhcXNjZDc2MHV5dWM4bHd4OXhzanFsNnR2EjZjZWxlc3RpYXZhbG9wZXIxcmZsdXRrM2V1dzhkY3dhZWh4d3VnY205cGV3a2RuNTZ4amxoMjYaCgoEdXRpYRICMzMKswEKJS9jb3Ntb3Muc3Rha2luZy52MWJldGExLk1zZ1VuZGVsZWdhdGUSiQEKQ2NlbGVzdGlhMWYzaHg3cmR6NXFhc3lnZmc1M2xmazRxNmxoNGZwYWp1dmFxc2NkNzYwdXl1Yzhsd3g5eHNqcWw2dHYSNmNlbGVzdGlhdmFsb3BlcjF4enJ5OGEzc3MwOHRwd21mcmc0dTgyN3JxdDlqdXczMjBkMDBseBoKCgR1dGlhEgIzNAqzAQolL2Nvc21vcy5zdGFraW5nLnYxYmV0YTEuTXNnVW5kZWxlZ2F0ZRKJAQpDY2VsZXN0aWExZjNoeDdyZHo1cWFzeWdmZzUzbGZrNHE2bGg0ZnBhanV2YXFzY2Q3NjB1eXVjOGx3eDl4c2pxbDZ0dhI2Y2VsZXN0aWF2YWxvcGVyMTl1cmc5YXdqendxOGQ0MHZ3amR2djB5dzlrZ2Voc2NmMHp4M2dzGgoKBHV0aWESAjM0CrMBCiUvY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dVbmRlbGVnYXRlEokBCkNjZWxlc3RpYTFmM2h4N3JkejVxYXN5Z2ZnNTNsZms0cTZsaDRmcGFqdXZhcXNjZDc2MHV5dWM4bHd4OXhzanFsNnR2EjZjZWxlc3RpYXZhbG9wZXIxOWp5ZTc2YW51cTdxdHJoNDV1c3NkemZwODBsMzlhNDk4Y2g1MzIaCgoEdXRpYRICMzQKswEKJS9jb3Ntb3Muc3Rha2luZy52MWJldGExLk1zZ1VuZGVsZWdhdGUSiQEKQ2NlbGVzdGlhMWYzaHg3cmR6NXFhc3lnZmc1M2xmazRxNmxoNGZwYWp1dmFxc2NkNzYwdXl1Yzhsd3g5eHNqcWw2dHYSNmNlbGVzdGlhdmFsb3BlcjF6NGVreHc3Nnd0d2N6ZXhmbmp4OHA5dmNudDJwajdwNmpydndyMBoKCgR1dGlhEgIzNAqzAQolL2Nvc21vcy5zdGFraW5nLnYxYmV0YTEuTXNnVW5kZWxlZ2F0ZRKJAQpDY2VsZXN0aWExZjNoeDdyZHo1cWFzeWdmZzUzbGZrNHE2bGg0ZnBhanV2YXFzY2Q3NjB1eXVjOGx3eDl4c2pxbDZ0dhI2Y2VsZXN0aWF2YWxvcGVyMXpmN3Q4YWh5dDc1dnQwcGN0dWxhMjB2czc2bWw2eXByZjVkbTBjGgoKBHV0aWESAjM0CrMBCiUvY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dVbmRlbGVnYXRlEokBCkNjZWxlc3RpYTFmM2h4N3JkejVxYXN5Z2ZnNTNsZms0cTZsaDRmcGFqdXZhcXNjZDc2MHV5dWM4bHd4OXhzanFsNnR2EjZjZWxlc3RpYXZhbG9wZXIxenp4cnE3dGpsc3F2djBqdHF0bHQ4Z2phYTY4cmFjdWtmNW45eGgaCgoEdXRpYRICMjYKswEKJS9jb3Ntb3Muc3Rha2luZy52MWJldGExLk1zZ1VuZGVsZWdhdGUSiQEKQ2NlbGVzdGlhMWYzaHg3cmR6NXFhc3lnZmc1M2xmazRxNmxoNGZwYWp1dmFxc2NkNzYwdXl1Yzhsd3g5eHNqcWw2dHYSNmNlbGVzdGlhdmFsb3BlcjF5YWNkNGY4aDdtOGZwdXY1Nmh6ZnM0cWRzN2M0NnZqdzc1czd5ehoKCgR1dGlhEgIzNAqzAQolL2Nvc21vcy5zdGFraW5nLnYxYmV0YTEuTXNnVW5kZWxlZ2F0ZRKJAQpDY2VsZXN0aWExZjNoeDdyZHo1cWFzeWdmZzUzbGZrNHE2bGg0ZnBhanV2YXFzY2Q3NjB1eXVjOGx3eDl4c2pxbDZ0dhI2Y2VsZXN0aWF2YWxvcGVyMTl6bDgwNjA2d242aGo1NzVyMjB3eHlla3U3dGV0OXAzNHVmcG5zGgoKBHV0aWESAjMzCrMBCiUvY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dVbmRlbGVnYXRlEokBCkNjZWxlc3RpYTFmM2h4N3JkejVxYXN5Z2ZnNTNsZms0cTZsaDRmcGFqdXZhcXNjZDc2MHV5dWM4bHd4OXhzanFsNnR2EjZjZWxlc3RpYXZhbG9wZXIxcXl1d3FqMGN4ZTZobHpqcnU1ODdueWd3d21naDAzaGE5dmU5YWMaCgoEdXRpYRICMzMKswEKJS9jb3Ntb3Muc3Rha2luZy52MWJldGExLk1zZ1VuZGVsZWdhdGUSiQEKQ2NlbGVzdGlhMWYzaHg3cmR6NXFhc3lnZmc1M2xmazRxNmxoNGZwYWp1dmFxc2NkNzYwdXl1Yzhsd3g5eHNqcWw2dHYSNmNlbGVzdGlhdmFsb3BlcjF5NzRwdGE3Z2Z5YXB2ZmxhZDl6MHh6amQ5OGc5Y3N6MGs0ejdhNxoKCgR1dGlhEgIzNAqzAQolL2Nvc21vcy5zdGFraW5nLnYxYmV0YTEuTXNnVW5kZWxlZ2F0ZRKJAQpDY2VsZXN0aWExZjNoeDdyZHo1cWFzeWdmZzUzbGZrNHE2bGg0ZnBhanV2YXFzY2Q3NjB1eXVjOGx3eDl4c2pxbDZ0dhI2Y2VsZXN0aWF2YWxvcGVyMXljbDM3cXdxMjIzeHRsOTRjNGhncWtqbGZwejM1a2Y3ejM4dWhrGgoKBHV0aWESAjM0CrMBCiUvY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dVbmRlbGVnYXRlEokBCkNjZWxlc3RpYTFmM2h4N3JkejVxYXN5Z2ZnNTNsZms0cTZsaDRmcGFqdXZhcXNjZDc2MHV5dWM4bHd4OXhzanFsNnR2EjZjZWxlc3RpYXZhbG9wZXIxcTN2NWN1Z2M4Y2RwdWQ4N3U0end5MGE3NHV4a2s2dTRxNGd4NHAaCgoEdXRpYRICMzQKswEKJS9jb3Ntb3Muc3Rha2luZy52MWJldGExLk1zZ1VuZGVsZWdhdGUSiQEKQ2NlbGVzdGlhMWYzaHg3cmR6NXFhc3lnZmc1M2xmazRxNmxoNGZwYWp1dmFxc2NkNzYwdXl1Yzhsd3g5eHNqcWw2dHYSNmNlbGVzdGlhdmFsb3BlcjFwajdwdmhjZWQ3c25wdnlqa2FjbGU3NGZkbWhkd2xxZHdleGFsaxoKCgR1dGlhEgIzNAqzAQolL2Nvc21vcy5zdGFraW5nLnYxYmV0YTEuTXNnVW5kZWxlZ2F0ZRKJAQpDY2VsZXN0aWExZjNoeDdyZHo1cWFzeWdmZzUzbGZrNHE2bGg0ZnBhanV2YXFzY2Q3NjB1eXVjOGx3eDl4c2pxbDZ0dhI2Y2VsZXN0aWF2YWxvcGVyMXk5dzg1cWpnaGR0YzJocmwydjd2c3UzdWxhbHljMjdlOTN1bWp4GgoKBHV0aWESAjM0\",\"memo\":\"\",\"type\":\"TYPE_EXECUTE_TX\"}",
						"packet_data_hex":          "7b2264617461223a2243724d42436955765932397a6257397a4c6e4e3059577470626d6375646a46695a5852684d53354e63326456626d526c6247566e5958526c456f6b42436b4e6a5a57786c633352705954466d4d3268344e334a6b656a567859584e355a325a6e4e544e735a6d733063545a736144526d6347467164585a6863584e6a5a4463324d48563564574d34624864344f58687a616e46734e6e5232456a5a6a5a57786c6333527059585a68624739775a584978636d5a73645852724d325631647a686b593364685a5768346433566e5932303563475633613252754e545a34616d786f4d6a596143676f4564585270595249434d7a4d4b7377454b4a53396a62334e7462334d756333526861326c755a7935324d574a6c644745784c6b317a5a3156755a4756735a576468644755536951454b51324e6c6247567a64476c684d57597a61486733636d52364e58466863336c6e5a6d63314d32786d617a52784e6d786f4e475a7759577031646d467863324e6b4e7a597764586c31597a68736433673565484e7163577732644859534e6d4e6c6247567a64476c68646d46736233426c636a4634656e4a354f47457a63334d774f4852776432316d636d6330645467794e334a7864446c716458637a4d6a426b4d44427365426f4b4367523164476c684567497a4e41717a41516f6c4c324e76633231766379357a644746726157356e4c6e5978596d56305954457554584e6e5657356b5a57786c5a3246305a524b4a41517044593256735a584e30615745785a6a4e6f654464795a486f316357467a6557646d5a7a557a62475a724e484532624767305a6e4268616e56325958467a593251334e6a42316558566a4f47783365446c346332707862445a3064684932593256735a584e306157463259577876634756794d546c31636d633559586471656e64784f4751304d485a33616d5232646a4235647a6c725a32566f63324e6d4d4870344d32647a47676f4b4248563061574553416a4d3043724d42436955765932397a6257397a4c6e4e3059577470626d6375646a46695a5852684d53354e63326456626d526c6247566e5958526c456f6b42436b4e6a5a57786c633352705954466d4d3268344e334a6b656a567859584e355a325a6e4e544e735a6d733063545a736144526d6347467164585a6863584e6a5a4463324d48563564574d34624864344f58687a616e46734e6e5232456a5a6a5a57786c6333527059585a68624739775a5849784f5770355a546332595735316354647864484a6f4e44563163334e6b656d5a774f4442734d7a6c684e446b34593267314d7a496143676f4564585270595249434d7a514b7377454b4a53396a62334e7462334d756333526861326c755a7935324d574a6c644745784c6b317a5a3156755a4756735a576468644755536951454b51324e6c6247567a64476c684d57597a61486733636d52364e58466863336c6e5a6d63314d32786d617a52784e6d786f4e475a7759577031646d467863324e6b4e7a597764586c31597a68736433673565484e7163577732644859534e6d4e6c6247567a64476c68646d46736233426c636a46364e475672654863334e6e643064324e365a58686d626d70344f484135646d4e7564444a77616a64774e6d7079646e64794d426f4b4367523164476c684567497a4e41717a41516f6c4c324e76633231766379357a644746726157356e4c6e5978596d56305954457554584e6e5657356b5a57786c5a3246305a524b4a41517044593256735a584e30615745785a6a4e6f654464795a486f316357467a6557646d5a7a557a62475a724e484532624767305a6e4268616e56325958467a593251334e6a42316558566a4f47783365446c346332707862445a3064684932593256735a584e306157463259577876634756794d58706d4e3351345957683564446331646e517763474e30645778684d6a4232637a633262577732655842795a6a566b6254426a47676f4b4248563061574553416a4d3043724d42436955765932397a6257397a4c6e4e3059577470626d6375646a46695a5852684d53354e63326456626d526c6247566e5958526c456f6b42436b4e6a5a57786c633352705954466d4d3268344e334a6b656a567859584e355a325a6e4e544e735a6d733063545a736144526d6347467164585a6863584e6a5a4463324d48563564574d34624864344f58687a616e46734e6e5232456a5a6a5a57786c6333527059585a68624739775a584978656e7034636e45336447707363334632646a427164484630624851345a32706859545934636d466a6457746d4e5734356547676143676f4564585270595249434d6a594b7377454b4a53396a62334e7462334d756333526861326c755a7935324d574a6c644745784c6b317a5a3156755a4756735a576468644755536951454b51324e6c6247567a64476c684d57597a61486733636d52364e58466863336c6e5a6d63314d32786d617a52784e6d786f4e475a7759577031646d467863324e6b4e7a597764586c31597a68736433673565484e7163577732644859534e6d4e6c6247567a64476c68646d46736233426c636a463559574e6b4e475934614464744f475a77645859314e6d68365a6e4d306357527a4e324d304e6e5a71647a6331637a643565686f4b4367523164476c684567497a4e41717a41516f6c4c324e76633231766379357a644746726157356e4c6e5978596d56305954457554584e6e5657356b5a57786c5a3246305a524b4a41517044593256735a584e30615745785a6a4e6f654464795a486f316357467a6557646d5a7a557a62475a724e484532624767305a6e4268616e56325958467a593251334e6a42316558566a4f47783365446c346332707862445a3064684932593256735a584e306157463259577876634756794d546c36624467774e6a41326432343261476f314e7a56794d6a423365486c6c61335533644756304f58417a4e48566d6347357a47676f4b4248563061574553416a4d7a43724d42436955765932397a6257397a4c6e4e3059577470626d6375646a46695a5852684d53354e63326456626d526c6247566e5958526c456f6b42436b4e6a5a57786c633352705954466d4d3268344e334a6b656a567859584e355a325a6e4e544e735a6d733063545a736144526d6347467164585a6863584e6a5a4463324d48563564574d34624864344f58687a616e46734e6e5232456a5a6a5a57786c6333527059585a68624739775a58497863586c31643346714d474e345a545a6f62487071636e55314f446475655764336432316e6144417a61474535646d553559574d6143676f4564585270595249434d7a4d4b7377454b4a53396a62334e7462334d756333526861326c755a7935324d574a6c644745784c6b317a5a3156755a4756735a576468644755536951454b51324e6c6247567a64476c684d57597a61486733636d52364e58466863336c6e5a6d63314d32786d617a52784e6d786f4e475a7759577031646d467863324e6b4e7a597764586c31597a68736433673565484e7163577732644859534e6d4e6c6247567a64476c68646d46736233426c636a46354e7a5277644745335a325a35595842325a6d78685a446c364d486836616d51354f47633559334e364d477330656a64684e786f4b4367523164476c684567497a4e41717a41516f6c4c324e76633231766379357a644746726157356e4c6e5978596d56305954457554584e6e5657356b5a57786c5a3246305a524b4a41517044593256735a584e30615745785a6a4e6f654464795a486f316357467a6557646d5a7a557a62475a724e484532624767305a6e4268616e56325958467a593251334e6a42316558566a4f47783365446c346332707862445a3064684932593256735a584e306157463259577876634756794d586c6a62444d33635864784d6a497a654852734f54526a4e47686e6357747162475a77656a4d3161325933656a4d346457687247676f4b4248563061574553416a4d3043724d42436955765932397a6257397a4c6e4e3059577470626d6375646a46695a5852684d53354e63326456626d526c6247566e5958526c456f6b42436b4e6a5a57786c633352705954466d4d3268344e334a6b656a567859584e355a325a6e4e544e735a6d733063545a736144526d6347467164585a6863584e6a5a4463324d48563564574d34624864344f58687a616e46734e6e5232456a5a6a5a57786c6333527059585a68624739775a58497863544e324e574e315a324d3459325277645751344e335530656e64354d4745334e48563461327332645452784e4764344e48416143676f4564585270595249434d7a514b7377454b4a53396a62334e7462334d756333526861326c755a7935324d574a6c644745784c6b317a5a3156755a4756735a576468644755536951454b51324e6c6247567a64476c684d57597a61486733636d52364e58466863336c6e5a6d63314d32786d617a52784e6d786f4e475a7759577031646d467863324e6b4e7a597764586c31597a68736433673565484e7163577732644859534e6d4e6c6247567a64476c68646d46736233426c636a4677616a6477646d686a5a57513363323577646e6c716132466a624755334e475a6b6257686b643278785a48646c6547467361786f4b4367523164476c684567497a4e41717a41516f6c4c324e76633231766379357a644746726157356e4c6e5978596d56305954457554584e6e5657356b5a57786c5a3246305a524b4a41517044593256735a584e30615745785a6a4e6f654464795a486f316357467a6557646d5a7a557a62475a724e484532624767305a6e4268616e56325958467a593251334e6a42316558566a4f47783365446c346332707862445a3064684932593256735a584e306157463259577876634756794d586b35647a67316357706e61475230597a4a6f636d7779646a64326333557a6457786862486c6a4d6a646c4f544e316257703447676f4b4248563061574553416a4d30222c226d656d6f223a22222c2274797065223a22545950455f455845435554455f5458227d",
						"packet_dst_channel":       "channel-166",
						"packet_dst_port":          "icahost",
						"packet_sequence":          "198",
						"packet_src_channel":       "channel-1332",
						"packet_src_port":          "icacontroller-neutron1z4z4gm7ujgwzdfuq7qtxrzv90ldm8uettvwmnejq8wpm25dc6jlqnh973t.DROP",
						"packet_timeout_height":    "0-0",
						"packet_timeout_timestamp": "1730874973641544615",
					},
				}, {
					Height: 1866988,
					Type:   "message",
					Data: map[string]any{
						"module": "ibc_channel",
					},
				}, {
					Height: 1866988,
					Type:   "withdraw_rewards",
					Data: map[string]any{
						"amount":    "0stake",
						"delegator": "celestia1f3hx7rdz5qasygfg53lfk4q6lh4fpajuvaqscd760uyuc8lwx9xsjql6tv",
						"validator": "celestiavaloper1rflutk3euw8dcwaehxwugcm9pewkdn56xjlh26",
					},
				}, {
					Height: 1866988,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "33utia",
						"spender": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 1866988,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "33utia",
						"receiver": "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
					},
				}, {
					Height: 1866988,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "33utia",
						"recipient": "celestia1tygms3xhhs3yv487phx3dw4a95jn7t7ls3yw4w",
						"sender":    "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 1866988,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 1866988,
					Type:   "unbond",
					Data: map[string]any{
						"amount":          "33utia",
						"completion_time": "2024-11-13T06:36:18Z",
						"validator":       "celestiavaloper1rflutk3euw8dcwaehxwugcm9pewkdn56xjlh26",
					},
				}, {
					Height: 1866988,
					Type:   "message",
					Data: map[string]any{
						"module": "staking",
						"sender": "celestia1f3hx7rdz5qasygfg53lfk4q6lh4fpajuvaqscd760uyuc8lwx9xsjql6tv",
					},
				}, {
					Height: 1866988,
					Type:   "ics27_packet",
					Data: map[string]any{
						"host_channel_id": "channel-166",
						"module":          "interchainaccounts",
						"success":         "true",
					},
				}, {
					Height: 1866988,
					Type:   "write_acknowledgement",
					Data: map[string]any{
						"packet_ack":               "{\"result\":\"Ej4KLS9jb3Ntb3Muc3Rha2luZy52MWJldGExLk1zZ1VuZGVsZWdhdGVSZXNwb25zZRINCgsI4pHRuQYQxL+MVBI+Ci0vY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dVbmRlbGVnYXRlUmVzcG9uc2USDQoLCOKR0bkGEMS/jFQSPgotL2Nvc21vcy5zdGFraW5nLnYxYmV0YTEuTXNnVW5kZWxlZ2F0ZVJlc3BvbnNlEg0KCwjikdG5BhDEv4xUEj4KLS9jb3Ntb3Muc3Rha2luZy52MWJldGExLk1zZ1VuZGVsZWdhdGVSZXNwb25zZRINCgsI4pHRuQYQxL+MVBI+Ci0vY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dVbmRlbGVnYXRlUmVzcG9uc2USDQoLCOKR0bkGEMS/jFQSPgotL2Nvc21vcy5zdGFraW5nLnYxYmV0YTEuTXNnVW5kZWxlZ2F0ZVJlc3BvbnNlEg0KCwjikdG5BhDEv4xUEj4KLS9jb3Ntb3Muc3Rha2luZy52MWJldGExLk1zZ1VuZGVsZWdhdGVSZXNwb25zZRINCgsI4pHRuQYQxL+MVBI+Ci0vY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dVbmRlbGVnYXRlUmVzcG9uc2USDQoLCOKR0bkGEMS/jFQSPgotL2Nvc21vcy5zdGFraW5nLnYxYmV0YTEuTXNnVW5kZWxlZ2F0ZVJlc3BvbnNlEg0KCwjikdG5BhDEv4xUEj4KLS9jb3Ntb3Muc3Rha2luZy52MWJldGExLk1zZ1VuZGVsZWdhdGVSZXNwb25zZRINCgsI4pHRuQYQxL+MVBI+Ci0vY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dVbmRlbGVnYXRlUmVzcG9uc2USDQoLCOKR0bkGEMS/jFQSPgotL2Nvc21vcy5zdGFraW5nLnYxYmV0YTEuTXNnVW5kZWxlZ2F0ZVJlc3BvbnNlEg0KCwjikdG5BhDEv4xUEj4KLS9jb3Ntb3Muc3Rha2luZy52MWJldGExLk1zZ1VuZGVsZWdhdGVSZXNwb25zZRINCgsI4pHRuQYQxL+MVBI+Ci0vY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dVbmRlbGVnYXRlUmVzcG9uc2USDQoLCOKR0bkGEMS/jFQSPgotL2Nvc21vcy5zdGFraW5nLnYxYmV0YTEuTXNnVW5kZWxlZ2F0ZVJlc3BvbnNlEg0KCwjikdG5BhDEv4xU\"}",
						"packet_ack_hex":           "7b22726573756c74223a22456a344b4c53396a62334e7462334d756333526861326c755a7935324d574a6c644745784c6b317a5a3156755a4756735a576468644756535a584e776232357a5a52494e436773493470485275515951784c2b4d5642492b436930765932397a6257397a4c6e4e3059577470626d6375646a46695a5852684d53354e63326456626d526c6247566e5958526c556d567a634739756332555344516f4c434f4b5230626b47454d532f6a46515350676f744c324e76633231766379357a644746726157356e4c6e5978596d56305954457554584e6e5657356b5a57786c5a3246305a564a6c63334276626e4e6c4567304b43776a696b6447354268444576347855456a344b4c53396a62334e7462334d756333526861326c755a7935324d574a6c644745784c6b317a5a3156755a4756735a576468644756535a584e776232357a5a52494e436773493470485275515951784c2b4d5642492b436930765932397a6257397a4c6e4e3059577470626d6375646a46695a5852684d53354e63326456626d526c6247566e5958526c556d567a634739756332555344516f4c434f4b5230626b47454d532f6a46515350676f744c324e76633231766379357a644746726157356e4c6e5978596d56305954457554584e6e5657356b5a57786c5a3246305a564a6c63334276626e4e6c4567304b43776a696b6447354268444576347855456a344b4c53396a62334e7462334d756333526861326c755a7935324d574a6c644745784c6b317a5a3156755a4756735a576468644756535a584e776232357a5a52494e436773493470485275515951784c2b4d5642492b436930765932397a6257397a4c6e4e3059577470626d6375646a46695a5852684d53354e63326456626d526c6247566e5958526c556d567a634739756332555344516f4c434f4b5230626b47454d532f6a46515350676f744c324e76633231766379357a644746726157356e4c6e5978596d56305954457554584e6e5657356b5a57786c5a3246305a564a6c63334276626e4e6c4567304b43776a696b6447354268444576347855456a344b4c53396a62334e7462334d756333526861326c755a7935324d574a6c644745784c6b317a5a3156755a4756735a576468644756535a584e776232357a5a52494e436773493470485275515951784c2b4d5642492b436930765932397a6257397a4c6e4e3059577470626d6375646a46695a5852684d53354e63326456626d526c6247566e5958526c556d567a634739756332555344516f4c434f4b5230626b47454d532f6a46515350676f744c324e76633231766379357a644746726157356e4c6e5978596d56305954457554584e6e5657356b5a57786c5a3246305a564a6c63334276626e4e6c4567304b43776a696b6447354268444576347855456a344b4c53396a62334e7462334d756333526861326c755a7935324d574a6c644745784c6b317a5a3156755a4756735a576468644756535a584e776232357a5a52494e436773493470485275515951784c2b4d5642492b436930765932397a6257397a4c6e4e3059577470626d6375646a46695a5852684d53354e63326456626d526c6247566e5958526c556d567a634739756332555344516f4c434f4b5230626b47454d532f6a46515350676f744c324e76633231766379357a644746726157356e4c6e5978596d56305954457554584e6e5657356b5a57786c5a3246305a564a6c63334276626e4e6c4567304b43776a696b6447354268444576347855227d",
						"packet_connection":        "connection-99",
						"packet_data":              "{\"data\":\"CrMBCiUvY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dVbmRlbGVnYXRlEokBCkNjZWxlc3RpYTFmM2h4N3JkejVxYXN5Z2ZnNTNsZms0cTZsaDRmcGFqdXZhcXNjZDc2MHV5dWM4bHd4OXhzanFsNnR2EjZjZWxlc3RpYXZhbG9wZXIxcmZsdXRrM2V1dzhkY3dhZWh4d3VnY205cGV3a2RuNTZ4amxoMjYaCgoEdXRpYRICMzMKswEKJS9jb3Ntb3Muc3Rha2luZy52MWJldGExLk1zZ1VuZGVsZWdhdGUSiQEKQ2NlbGVzdGlhMWYzaHg3cmR6NXFhc3lnZmc1M2xmazRxNmxoNGZwYWp1dmFxc2NkNzYwdXl1Yzhsd3g5eHNqcWw2dHYSNmNlbGVzdGlhdmFsb3BlcjF4enJ5OGEzc3MwOHRwd21mcmc0dTgyN3JxdDlqdXczMjBkMDBseBoKCgR1dGlhEgIzNAqzAQolL2Nvc21vcy5zdGFraW5nLnYxYmV0YTEuTXNnVW5kZWxlZ2F0ZRKJAQpDY2VsZXN0aWExZjNoeDdyZHo1cWFzeWdmZzUzbGZrNHE2bGg0ZnBhanV2YXFzY2Q3NjB1eXVjOGx3eDl4c2pxbDZ0dhI2Y2VsZXN0aWF2YWxvcGVyMTl1cmc5YXdqendxOGQ0MHZ3amR2djB5dzlrZ2Voc2NmMHp4M2dzGgoKBHV0aWESAjM0CrMBCiUvY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dVbmRlbGVnYXRlEokBCkNjZWxlc3RpYTFmM2h4N3JkejVxYXN5Z2ZnNTNsZms0cTZsaDRmcGFqdXZhcXNjZDc2MHV5dWM4bHd4OXhzanFsNnR2EjZjZWxlc3RpYXZhbG9wZXIxOWp5ZTc2YW51cTdxdHJoNDV1c3NkemZwODBsMzlhNDk4Y2g1MzIaCgoEdXRpYRICMzQKswEKJS9jb3Ntb3Muc3Rha2luZy52MWJldGExLk1zZ1VuZGVsZWdhdGUSiQEKQ2NlbGVzdGlhMWYzaHg3cmR6NXFhc3lnZmc1M2xmazRxNmxoNGZwYWp1dmFxc2NkNzYwdXl1Yzhsd3g5eHNqcWw2dHYSNmNlbGVzdGlhdmFsb3BlcjF6NGVreHc3Nnd0d2N6ZXhmbmp4OHA5dmNudDJwajdwNmpydndyMBoKCgR1dGlhEgIzNAqzAQolL2Nvc21vcy5zdGFraW5nLnYxYmV0YTEuTXNnVW5kZWxlZ2F0ZRKJAQpDY2VsZXN0aWExZjNoeDdyZHo1cWFzeWdmZzUzbGZrNHE2bGg0ZnBhanV2YXFzY2Q3NjB1eXVjOGx3eDl4c2pxbDZ0dhI2Y2VsZXN0aWF2YWxvcGVyMXpmN3Q4YWh5dDc1dnQwcGN0dWxhMjB2czc2bWw2eXByZjVkbTBjGgoKBHV0aWESAjM0CrMBCiUvY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dVbmRlbGVnYXRlEokBCkNjZWxlc3RpYTFmM2h4N3JkejVxYXN5Z2ZnNTNsZms0cTZsaDRmcGFqdXZhcXNjZDc2MHV5dWM4bHd4OXhzanFsNnR2EjZjZWxlc3RpYXZhbG9wZXIxenp4cnE3dGpsc3F2djBqdHF0bHQ4Z2phYTY4cmFjdWtmNW45eGgaCgoEdXRpYRICMjYKswEKJS9jb3Ntb3Muc3Rha2luZy52MWJldGExLk1zZ1VuZGVsZWdhdGUSiQEKQ2NlbGVzdGlhMWYzaHg3cmR6NXFhc3lnZmc1M2xmazRxNmxoNGZwYWp1dmFxc2NkNzYwdXl1Yzhsd3g5eHNqcWw2dHYSNmNlbGVzdGlhdmFsb3BlcjF5YWNkNGY4aDdtOGZwdXY1Nmh6ZnM0cWRzN2M0NnZqdzc1czd5ehoKCgR1dGlhEgIzNAqzAQolL2Nvc21vcy5zdGFraW5nLnYxYmV0YTEuTXNnVW5kZWxlZ2F0ZRKJAQpDY2VsZXN0aWExZjNoeDdyZHo1cWFzeWdmZzUzbGZrNHE2bGg0ZnBhanV2YXFzY2Q3NjB1eXVjOGx3eDl4c2pxbDZ0dhI2Y2VsZXN0aWF2YWxvcGVyMTl6bDgwNjA2d242aGo1NzVyMjB3eHlla3U3dGV0OXAzNHVmcG5zGgoKBHV0aWESAjMzCrMBCiUvY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dVbmRlbGVnYXRlEokBCkNjZWxlc3RpYTFmM2h4N3JkejVxYXN5Z2ZnNTNsZms0cTZsaDRmcGFqdXZhcXNjZDc2MHV5dWM4bHd4OXhzanFsNnR2EjZjZWxlc3RpYXZhbG9wZXIxcXl1d3FqMGN4ZTZobHpqcnU1ODdueWd3d21naDAzaGE5dmU5YWMaCgoEdXRpYRICMzMKswEKJS9jb3Ntb3Muc3Rha2luZy52MWJldGExLk1zZ1VuZGVsZWdhdGUSiQEKQ2NlbGVzdGlhMWYzaHg3cmR6NXFhc3lnZmc1M2xmazRxNmxoNGZwYWp1dmFxc2NkNzYwdXl1Yzhsd3g5eHNqcWw2dHYSNmNlbGVzdGlhdmFsb3BlcjF5NzRwdGE3Z2Z5YXB2ZmxhZDl6MHh6amQ5OGc5Y3N6MGs0ejdhNxoKCgR1dGlhEgIzNAqzAQolL2Nvc21vcy5zdGFraW5nLnYxYmV0YTEuTXNnVW5kZWxlZ2F0ZRKJAQpDY2VsZXN0aWExZjNoeDdyZHo1cWFzeWdmZzUzbGZrNHE2bGg0ZnBhanV2YXFzY2Q3NjB1eXVjOGx3eDl4c2pxbDZ0dhI2Y2VsZXN0aWF2YWxvcGVyMXljbDM3cXdxMjIzeHRsOTRjNGhncWtqbGZwejM1a2Y3ejM4dWhrGgoKBHV0aWESAjM0CrMBCiUvY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dVbmRlbGVnYXRlEokBCkNjZWxlc3RpYTFmM2h4N3JkejVxYXN5Z2ZnNTNsZms0cTZsaDRmcGFqdXZhcXNjZDc2MHV5dWM4bHd4OXhzanFsNnR2EjZjZWxlc3RpYXZhbG9wZXIxcTN2NWN1Z2M4Y2RwdWQ4N3U0end5MGE3NHV4a2s2dTRxNGd4NHAaCgoEdXRpYRICMzQKswEKJS9jb3Ntb3Muc3Rha2luZy52MWJldGExLk1zZ1VuZGVsZWdhdGUSiQEKQ2NlbGVzdGlhMWYzaHg3cmR6NXFhc3lnZmc1M2xmazRxNmxoNGZwYWp1dmFxc2NkNzYwdXl1Yzhsd3g5eHNqcWw2dHYSNmNlbGVzdGlhdmFsb3BlcjFwajdwdmhjZWQ3c25wdnlqa2FjbGU3NGZkbWhkd2xxZHdleGFsaxoKCgR1dGlhEgIzNAqzAQolL2Nvc21vcy5zdGFraW5nLnYxYmV0YTEuTXNnVW5kZWxlZ2F0ZRKJAQpDY2VsZXN0aWExZjNoeDdyZHo1cWFzeWdmZzUzbGZrNHE2bGg0ZnBhanV2YXFzY2Q3NjB1eXVjOGx3eDl4c2pxbDZ0dhI2Y2VsZXN0aWF2YWxvcGVyMXk5dzg1cWpnaGR0YzJocmwydjd2c3UzdWxhbHljMjdlOTN1bWp4GgoKBHV0aWESAjM0\",\"memo\":\"\",\"type\":\"TYPE_EXECUTE_TX\"}",
						"packet_data_hex":          "7b2264617461223a2243724d42436955765932397a6257397a4c6e4e3059577470626d6375646a46695a5852684d53354e63326456626d526c6247566e5958526c456f6b42436b4e6a5a57786c633352705954466d4d3268344e334a6b656a567859584e355a325a6e4e544e735a6d733063545a736144526d6347467164585a6863584e6a5a4463324d48563564574d34624864344f58687a616e46734e6e5232456a5a6a5a57786c6333527059585a68624739775a584978636d5a73645852724d325631647a686b593364685a5768346433566e5932303563475633613252754e545a34616d786f4d6a596143676f4564585270595249434d7a4d4b7377454b4a53396a62334e7462334d756333526861326c755a7935324d574a6c644745784c6b317a5a3156755a4756735a576468644755536951454b51324e6c6247567a64476c684d57597a61486733636d52364e58466863336c6e5a6d63314d32786d617a52784e6d786f4e475a7759577031646d467863324e6b4e7a597764586c31597a68736433673565484e7163577732644859534e6d4e6c6247567a64476c68646d46736233426c636a4634656e4a354f47457a63334d774f4852776432316d636d6330645467794e334a7864446c716458637a4d6a426b4d44427365426f4b4367523164476c684567497a4e41717a41516f6c4c324e76633231766379357a644746726157356e4c6e5978596d56305954457554584e6e5657356b5a57786c5a3246305a524b4a41517044593256735a584e30615745785a6a4e6f654464795a486f316357467a6557646d5a7a557a62475a724e484532624767305a6e4268616e56325958467a593251334e6a42316558566a4f47783365446c346332707862445a3064684932593256735a584e306157463259577876634756794d546c31636d633559586471656e64784f4751304d485a33616d5232646a4235647a6c725a32566f63324e6d4d4870344d32647a47676f4b4248563061574553416a4d3043724d42436955765932397a6257397a4c6e4e3059577470626d6375646a46695a5852684d53354e63326456626d526c6247566e5958526c456f6b42436b4e6a5a57786c633352705954466d4d3268344e334a6b656a567859584e355a325a6e4e544e735a6d733063545a736144526d6347467164585a6863584e6a5a4463324d48563564574d34624864344f58687a616e46734e6e5232456a5a6a5a57786c6333527059585a68624739775a5849784f5770355a546332595735316354647864484a6f4e44563163334e6b656d5a774f4442734d7a6c684e446b34593267314d7a496143676f4564585270595249434d7a514b7377454b4a53396a62334e7462334d756333526861326c755a7935324d574a6c644745784c6b317a5a3156755a4756735a576468644755536951454b51324e6c6247567a64476c684d57597a61486733636d52364e58466863336c6e5a6d63314d32786d617a52784e6d786f4e475a7759577031646d467863324e6b4e7a597764586c31597a68736433673565484e7163577732644859534e6d4e6c6247567a64476c68646d46736233426c636a46364e475672654863334e6e643064324e365a58686d626d70344f484135646d4e7564444a77616a64774e6d7079646e64794d426f4b4367523164476c684567497a4e41717a41516f6c4c324e76633231766379357a644746726157356e4c6e5978596d56305954457554584e6e5657356b5a57786c5a3246305a524b4a41517044593256735a584e30615745785a6a4e6f654464795a486f316357467a6557646d5a7a557a62475a724e484532624767305a6e4268616e56325958467a593251334e6a42316558566a4f47783365446c346332707862445a3064684932593256735a584e306157463259577876634756794d58706d4e3351345957683564446331646e517763474e30645778684d6a4232637a633262577732655842795a6a566b6254426a47676f4b4248563061574553416a4d3043724d42436955765932397a6257397a4c6e4e3059577470626d6375646a46695a5852684d53354e63326456626d526c6247566e5958526c456f6b42436b4e6a5a57786c633352705954466d4d3268344e334a6b656a567859584e355a325a6e4e544e735a6d733063545a736144526d6347467164585a6863584e6a5a4463324d48563564574d34624864344f58687a616e46734e6e5232456a5a6a5a57786c6333527059585a68624739775a584978656e7034636e45336447707363334632646a427164484630624851345a32706859545934636d466a6457746d4e5734356547676143676f4564585270595249434d6a594b7377454b4a53396a62334e7462334d756333526861326c755a7935324d574a6c644745784c6b317a5a3156755a4756735a576468644755536951454b51324e6c6247567a64476c684d57597a61486733636d52364e58466863336c6e5a6d63314d32786d617a52784e6d786f4e475a7759577031646d467863324e6b4e7a597764586c31597a68736433673565484e7163577732644859534e6d4e6c6247567a64476c68646d46736233426c636a463559574e6b4e475934614464744f475a77645859314e6d68365a6e4d306357527a4e324d304e6e5a71647a6331637a643565686f4b4367523164476c684567497a4e41717a41516f6c4c324e76633231766379357a644746726157356e4c6e5978596d56305954457554584e6e5657356b5a57786c5a3246305a524b4a41517044593256735a584e30615745785a6a4e6f654464795a486f316357467a6557646d5a7a557a62475a724e484532624767305a6e4268616e56325958467a593251334e6a42316558566a4f47783365446c346332707862445a3064684932593256735a584e306157463259577876634756794d546c36624467774e6a41326432343261476f314e7a56794d6a423365486c6c61335533644756304f58417a4e48566d6347357a47676f4b4248563061574553416a4d7a43724d42436955765932397a6257397a4c6e4e3059577470626d6375646a46695a5852684d53354e63326456626d526c6247566e5958526c456f6b42436b4e6a5a57786c633352705954466d4d3268344e334a6b656a567859584e355a325a6e4e544e735a6d733063545a736144526d6347467164585a6863584e6a5a4463324d48563564574d34624864344f58687a616e46734e6e5232456a5a6a5a57786c6333527059585a68624739775a58497863586c31643346714d474e345a545a6f62487071636e55314f446475655764336432316e6144417a61474535646d553559574d6143676f4564585270595249434d7a4d4b7377454b4a53396a62334e7462334d756333526861326c755a7935324d574a6c644745784c6b317a5a3156755a4756735a576468644755536951454b51324e6c6247567a64476c684d57597a61486733636d52364e58466863336c6e5a6d63314d32786d617a52784e6d786f4e475a7759577031646d467863324e6b4e7a597764586c31597a68736433673565484e7163577732644859534e6d4e6c6247567a64476c68646d46736233426c636a46354e7a5277644745335a325a35595842325a6d78685a446c364d486836616d51354f47633559334e364d477330656a64684e786f4b4367523164476c684567497a4e41717a41516f6c4c324e76633231766379357a644746726157356e4c6e5978596d56305954457554584e6e5657356b5a57786c5a3246305a524b4a41517044593256735a584e30615745785a6a4e6f654464795a486f316357467a6557646d5a7a557a62475a724e484532624767305a6e4268616e56325958467a593251334e6a42316558566a4f47783365446c346332707862445a3064684932593256735a584e306157463259577876634756794d586c6a62444d33635864784d6a497a654852734f54526a4e47686e6357747162475a77656a4d3161325933656a4d346457687247676f4b4248563061574553416a4d3043724d42436955765932397a6257397a4c6e4e3059577470626d6375646a46695a5852684d53354e63326456626d526c6247566e5958526c456f6b42436b4e6a5a57786c633352705954466d4d3268344e334a6b656a567859584e355a325a6e4e544e735a6d733063545a736144526d6347467164585a6863584e6a5a4463324d48563564574d34624864344f58687a616e46734e6e5232456a5a6a5a57786c6333527059585a68624739775a58497863544e324e574e315a324d3459325277645751344e335530656e64354d4745334e48563461327332645452784e4764344e48416143676f4564585270595249434d7a514b7377454b4a53396a62334e7462334d756333526861326c755a7935324d574a6c644745784c6b317a5a3156755a4756735a576468644755536951454b51324e6c6247567a64476c684d57597a61486733636d52364e58466863336c6e5a6d63314d32786d617a52784e6d786f4e475a7759577031646d467863324e6b4e7a597764586c31597a68736433673565484e7163577732644859534e6d4e6c6247567a64476c68646d46736233426c636a4677616a6477646d686a5a57513363323577646e6c716132466a624755334e475a6b6257686b643278785a48646c6547467361786f4b4367523164476c684567497a4e41717a41516f6c4c324e76633231766379357a644746726157356e4c6e5978596d56305954457554584e6e5657356b5a57786c5a3246305a524b4a41517044593256735a584e30615745785a6a4e6f654464795a486f316357467a6557646d5a7a557a62475a724e484532624767305a6e4268616e56325958467a593251334e6a42316558566a4f47783365446c346332707862445a3064684932593256735a584e306157463259577876634756794d586b35647a67316357706e61475230597a4a6f636d7779646a64326333557a6457786862486c6a4d6a646c4f544e316257703447676f4b4248563061574553416a4d30222c226d656d6f223a22222c2274797065223a22545950455f455845435554455f5458227d",
						"packet_dst_channel":       "channel-166",
						"packet_dst_port":          "icahost",
						"packet_sequence":          "198",
						"packet_src_channel":       "channel-1332",
						"packet_src_port":          "icacontroller-neutron1z4z4gm7ujgwzdfuq7qtxrzv90ldm8uettvwmnejq8wpm25dc6jlqnh973t.DROP",
						"packet_timeout_height":    "0-0",
						"packet_timeout_timestamp": "1730874973641544615",
					},
				}, {
					Height: 1866988,
					Type:   "message",
					Data: map[string]any{
						"module": "ibc_channel",
					},
				},
			},
			msg: []*storage.Message{
				{
					Type:   types.MsgRecvPacket,
					Height: 1866988,
					Data: map[string]any{
						"Packet": map[string]any{
							"Data": map[string]any{
								"Memo": "rootulp",
								"Type": icaTypes.EXECUTE_TX,
								"Data": []cosmosTypes.Msg{
									&cosmosStakingTypes.MsgUndelegate{
										DelegatorAddress: "celestia1f3hx7rdz5qasygfg53lfk4q6lh4fpajuvaqscd760uyuc8lwx9xsjql6tv",
										ValidatorAddress: "celestiavaloper1rflutk3euw8dcwaehxwugcm9pewkdn56xjlh26",
										Amount:           cosmosTypes.NewCoin("utia", math.NewInt(33)),
									},
									&cosmosStakingTypes.MsgUndelegate{
										DelegatorAddress: "celestia1f3hx7rdz5qasygfg53lfk4q6lh4fpajuvaqscd760uyuc8lwx9xsjql6tv",
										ValidatorAddress: "celestiavaloper1xzry8a3ss08tpwmfrg4u827rqt9juw320d00lx",
										Amount:           cosmosTypes.NewCoin("utia", math.NewInt(34)),
									},
									&cosmosStakingTypes.MsgUndelegate{
										DelegatorAddress: "celestia1f3hx7rdz5qasygfg53lfk4q6lh4fpajuvaqscd760uyuc8lwx9xsjql6tv",
										ValidatorAddress: "celestiavaloper19urg9awjzwq8d40vwjdvv0yw9kgehscf0zx3gs",
										Amount:           cosmosTypes.NewCoin("utia", math.NewInt(34)),
									},
									&cosmosStakingTypes.MsgUndelegate{
										DelegatorAddress: "celestia1f3hx7rdz5qasygfg53lfk4q6lh4fpajuvaqscd760uyuc8lwx9xsjql6tv",
										ValidatorAddress: "celestiavaloper19jye76anuq7qtrh45ussdzfp80l39a498ch532",
										Amount:           cosmosTypes.NewCoin("utia", math.NewInt(34)),
									},
									&cosmosStakingTypes.MsgUndelegate{
										DelegatorAddress: "celestia1f3hx7rdz5qasygfg53lfk4q6lh4fpajuvaqscd760uyuc8lwx9xsjql6tv",
										ValidatorAddress: "celestiavaloper1z4ekxw76wtwczexfnjx8p9vcnt2pj7p6jrvwr0",
										Amount:           cosmosTypes.NewCoin("utia", math.NewInt(34)),
									},
								},
							},
							"DestinationChannel": "channel-166",
							"DestinationPort":    "icahost",
							"Sequence":           198,
							"SourceChannel":      "channel-1332",
							"SourcePort":         "icacontroller-neutron1z4z4gm7ujgwzdfuq7qtxrzv90ldm8uettvwmnejq8wpm25dc6jlqnh973t.DROP",
							"TimeoutHeight": map[string]any{
								"RevisionHeight": 0,
								"RevisionNumber": 0,
							},
							"TimeoutTimestamp": 1730874973641544700,
						},
						"ProofCommitment": "CqsICqgICosBY29tbWl0bWVudHMvcG9ydHMvaWNhY29udHJvbGxlci1uZXV0cm9uMXo0ejRnbTd1amd3emRmdXE3cXR4cnp2OTBsZG04dWV0dHZ3bW5lanE4d3BtMjVkYzZqbHFuaDk3M3QuRFJPUC9jaGFubmVscy9jaGFubmVsLTEzMzIvc2VxdWVuY2VzLzE5OBIglNx9jQfMmakdBRompajfMaVrW+DnqwAti+iF/z0eaz4aDggBGAEgASoGAAKIobwTIiwIARIoAgSIobwTIJ2ijpC9p+iuWxc95+/apH9VnCpPIoy8ztJenWHjBtgmICIuCAESBwQGiKG8EyAaISCS4uNNnLAzWZLYBTCiWDgfCZLMzG4xZI+Md9DBCTe7dSIuCAESBwYOiKG8EyAaISAtRp22nW3pTVw6uZztOwZXQ2JW8KKTkoNie19i/gY9ACIuCAESBwgeiKG8EyAaISCaC2kgpW7sPcaXqHJOgedGjHhiFLJCpKCkwORzybM2viIuCAESBww8iKG8EyAaISABCHJqvFDifH2svcfSTatMIwgufE+86LO/Vj/35Lp1ByIvCAESCBCgAYihvBMgGiEgpy53Dp+hbpvptnq+m17vwH9fvVshm+6TRrwXQLxQ/WEiLQgBEikUtgWIobwTIKf/4/pjbUiuPQvALTDFZvMaaJRBPqb+dPJWf6EIdnhXICItCAESKRjkFYihvBMgEHbNgiPAb5T3lcgy9erpNN/DODNkrhKkpyDQ2Fw207cgIi0IARIpGugpiKG8EyAQgZHfpr/Ja7JT6uwvUDKdLscAFstKFGVRcC7LEKIyCSAiLwgBEgge6mqIobwTIBohICumC7lBSFNfLEYuIvxBiL51z+Xm1fzckKoHy6aZVfXTIjAIARIJIIDbAYihvBMgGiEgn+0o46g/cs9qHvyRacusjCx2uk7UKtS5tx3xVylVgm0iMAgBEgkikpMDiKG8EyAaISAti7A48RwWGgts7n6PGIEHgXkzrgBerN53aiMVX3v1BSIuCAESKibE1A2IobwTIGKtRFBN0Qz2D1tFM4SF+HkjnL2kjbzl+3jfrj2niIUvICIwCAESCSjyrBWIobwTIBohIAOAU+cs/ZyTuoQ99kOqty7H/MJcGZvT39MYJfesYB+iIi4IARIqKu6uKIihvBMgeJR3mcbl8UQ/JdyVMNYIKgAv19W5KRboWeNaiQQft3YgIjAIARIJLPDTQYihvBMgGiEg8nfop5yNGlo+b7AtbSyH45cdBgZSqQD7X2ut1qEWixciMAgBEgkurr5giKG8EyAaISBW6XFpn/YYDIozuuDzYD7RxgD+/NiQzjKmLk5dT3xn0yIvCAESKzCEkqMBiKG8EyCmOwYe1iDcHLS+Berl9Dx+0lwKy9I8mJCBGB3Az/WvWSAKpwIKpAIKA2liYxIgMLPK1Fftwld+yiAnuOSp+ry7UIGJqWE6i0YIk7y1WvgaCQgBGAEgASoBACInCAESAQEaIEJm87oCvlYmuoLTv4CaSKuWd5mAMZNjd010gFGYkQ2uIicIARIBARoge2kzcTZqjemVIrTDVXCkKcy+4J9J/wEx04Mw5GvXvfsiJQgBEiEBC3+p/SBuuJ1EEs6m8rgKF5Fh049YqzCdOve1cNdT1voiJwgBEgEBGiAMXTr5zCUtqfpqLdyluIwtUMgQE5uWRY0nv490gobKNCIlCAESIQEE5jcAINBB8Q0Flf3ml5USo/BqiKS0mE5yaLwjx8kHICInCAESAQEaIOKVDxCvKujpFR9/5PeJJlDQmMMNBQjqzXEg2m9yg93u",
						"ProofHeight": map[string]any{
							"RevisionHeight": 20416581,
							"RevisionNumber": 1,
						},
						"Signer": "celestia19hu0gjgp6yk822r83a0g2ytlc7mna3aqlf5f63",
					},
				},
			},
			idx: testsuite.Ptr(0),
		}, {
			name: "recv packet test 4",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 2432340,
					Type:   "message",
					Data: map[string]any{
						"action": "/ibc.core.channel.v1.MsgRecvPacket",
					},
				}, {
					Height: 2432340,
					Type:   "recv_packet",
					Data: map[string]any{
						"packet_channel_ordering":  "ORDER_ORDERED",
						"packet_connection":        "connection-60",
						"packet_data":              "{\"data\":\"CrQBCiMvY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dEZWxlZ2F0ZRKMAQpDY2VsZXN0aWExbHVtYWVtemV0Z3g3ZmE4Z2N3M3dhYWdldnM2cnFoZWpmd3hkdnN6enhsbmVzZWxzaHd4cW01cm14bhI2Y2VsZXN0aWF2YWxvcGVyMTMzdDRncHY0dmhwcWdmbjlncjhsNHU0MjN6cmdsZzhya3FldXByGg0KBHV0aWESBTEwMDAw\",\"memo\":\"perf/celestiavaloper133t4gpv4vhpqgfn9gr8l4u423zrglg8rkqeupr\",\"type\":\"TYPE_EXECUTE_TX\"}",
						"packet_data_hex":          "7b2264617461223a224372514243694d765932397a6257397a4c6e4e3059577470626d6375646a46695a5852684d53354e633264455a57786c5a3246305a524b4d41517044593256735a584e30615745786248567459575674656d56305a3367335a6d45345a324e334d3364685957646c646e4d32636e466f5a57706d6433686b646e4e36656e6873626d567a5a57787a6148643463573031636d313462684932593256735a584e306157463259577876634756794d544d7a6444526e63485930646d68776357646d626a6c6e636a68734e4855304d6a4e36636d64735a7a68796133466c645842794767304b4248563061574553425445774d444177222c226d656d6f223a22706572662f63656c657374696176616c6f7065723133337434677076347668707167666e396772386c34753432337a72676c6738726b7165757072222c2274797065223a22545950455f455845435554455f5458227d",
						"packet_dst_channel":       "channel-44",
						"packet_dst_port":          "icahost",
						"packet_sequence":          "14",
						"packet_src_channel":       "channel-210",
						"packet_src_port":          "icacontroller-celestia.performance",
						"packet_timeout_height":    "0-0",
						"packet_timeout_timestamp": "1727402253142295584",
					},
				}, {
					Height: 2432340,
					Type:   "message",
					Data: map[string]any{
						"module": "ibc_channel",
					},
				}, {
					Height: 2432340,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "10000utia",
						"spender": "celestia1lumaemzetgx7fa8gcw3waagevs6rqhejfwxdvszzxlneselshwxqm5rmxn",
					},
				}, {
					Height: 2432340,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "10000utia",
						"receiver": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 2432340,
					Type:   "delegate",
					Data: map[string]any{
						"amount":     "10000utia",
						"new_shares": "10000.000000000000000000",
						"validator":  "celestiavaloper133t4gpv4vhpqgfn9gr8l4u423zrglg8rkqeupr",
					},
				}, {
					Height: 2432340,
					Type:   "message",
					Data: map[string]any{
						"module": "staking",
						"sender": "celestia1lumaemzetgx7fa8gcw3waagevs6rqhejfwxdvszzxlneselshwxqm5rmxn",
					},
				}, {
					Height: 2432340,
					Type:   "ics27_packet",
					Data: map[string]any{
						"host_channel_id": "channel-44",
						"module":          "interchainaccounts",
						"success":         "true",
					},
				}, {
					Height: 2432340,
					Type:   "write_acknowledgement",
					Data: map[string]any{
						"packet_ack":               "{\"result\":\"Ei0KKy9jb3Ntb3Muc3Rha2luZy52MWJldGExLk1zZ0RlbGVnYXRlUmVzcG9uc2U=\"}",
						"packet_ack_hex":           "7b22726573756c74223a224569304b4b79396a62334e7462334d756333526861326c755a7935324d574a6c644745784c6b317a5a30526c6247566e5958526c556d567a634739756332553d227d",
						"packet_connection":        "connection-60",
						"packet_data":              "{\"data\":\"CrQBCiMvY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dEZWxlZ2F0ZRKMAQpDY2VsZXN0aWExbHVtYWVtemV0Z3g3ZmE4Z2N3M3dhYWdldnM2cnFoZWpmd3hkdnN6enhsbmVzZWxzaHd4cW01cm14bhI2Y2VsZXN0aWF2YWxvcGVyMTMzdDRncHY0dmhwcWdmbjlncjhsNHU0MjN6cmdsZzhya3FldXByGg0KBHV0aWESBTEwMDAw\",\"memo\":\"perf/celestiavaloper133t4gpv4vhpqgfn9gr8l4u423zrglg8rkqeupr\",\"type\":\"TYPE_EXECUTE_TX\"}",
						"packet_data_hex":          "7b2264617461223a224372514243694d765932397a6257397a4c6e4e3059577470626d6375646a46695a5852684d53354e633264455a57786c5a3246305a524b4d41517044593256735a584e30615745786248567459575674656d56305a3367335a6d45345a324e334d3364685957646c646e4d32636e466f5a57706d6433686b646e4e36656e6873626d567a5a57787a6148643463573031636d313462684932593256735a584e306157463259577876634756794d544d7a6444526e63485930646d68776357646d626a6c6e636a68734e4855304d6a4e36636d64735a7a68796133466c645842794767304b4248563061574553425445774d444177222c226d656d6f223a22706572662f63656c657374696176616c6f7065723133337434677076347668707167666e396772386c34753432337a72676c6738726b7165757072222c2274797065223a22545950455f455845435554455f5458227d",
						"packet_dst_channel":       "channel-44",
						"packet_dst_port":          "icahost",
						"packet_sequence":          "14",
						"packet_src_channel":       "channel-210",
						"packet_src_port":          "icacontroller-celestia.performance",
						"packet_timeout_height":    "0-0",
						"packet_timeout_timestamp": "1727402253142295584",
					},
				}, {
					Height: 2432340,
					Type:   "message",
					Data: map[string]any{
						"module": "ibc_channel",
					},
				},

				{
					Height: 2432340,
					Type:   "message",
					Data: map[string]any{
						"action": "/ibc.core.channel.v1.MsgRecvPacket",
					},
				}, {
					Height: 2432340,
					Type:   "recv_packet",
					Data: map[string]any{
						"packet_channel_ordering":  "ORDER_ORDERED",
						"packet_connection":        "connection-60",
						"packet_data":              "{\"data\":\"CrQBCiMvY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dEZWxlZ2F0ZRKMAQpDY2VsZXN0aWExbHVtYWVtemV0Z3g3ZmE4Z2N3M3dhYWdldnM2cnFoZWpmd3hkdnN6enhsbmVzZWxzaHd4cW01cm14bhI2Y2VsZXN0aWF2YWxvcGVyMWNzMzd0dm1haGF2dzh4Y256Y2d5ejM0MnNoMGFsMzdtYTR6cWF0Gg0KBHV0aWESBTEwMDAw\",\"memo\":\"perf/celestiavaloper1cs37tvmahavw8xcnzcgyz342sh0al37ma4zqat\",\"type\":\"TYPE_EXECUTE_TX\"}",
						"packet_data_hex":          "7b2264617461223a224372514243694d765932397a6257397a4c6e4e3059577470626d6375646a46695a5852684d53354e633264455a57786c5a3246305a524b4d41517044593256735a584e30615745786248567459575674656d56305a3367335a6d45345a324e334d3364685957646c646e4d32636e466f5a57706d6433686b646e4e36656e6873626d567a5a57787a6148643463573031636d313462684932593256735a584e306157463259577876634756794d574e7a4d7a6430646d316861474632647a68345932353659326435656a4d304d6e4e6f4d4746734d7a647459545236635746304767304b4248563061574553425445774d444177222c226d656d6f223a22706572662f63656c657374696176616c6f706572316373333774766d61686176773878636e7a6367797a333432736830616c33376d61347a716174222c2274797065223a22545950455f455845435554455f5458227d",
						"packet_dst_channel":       "channel-44",
						"packet_dst_port":          "icahost",
						"packet_sequence":          "15",
						"packet_src_channel":       "channel-210",
						"packet_src_port":          "icacontroller-celestia.performance",
						"packet_timeout_height":    "0-0",
						"packet_timeout_timestamp": "1727402253142295584",
					},
				}, {
					Height: 2432340,
					Type:   "message",
					Data: map[string]any{
						"module": "ibc_channel",
					},
				}, {
					Height: 2432340,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "10000utia",
						"spender": "celestia1lumaemzetgx7fa8gcw3waagevs6rqhejfwxdvszzxlneselshwxqm5rmxn",
					},
				}, {
					Height: 2432340,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "10000utia",
						"receiver": "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
					},
				}, {
					Height: 2432340,
					Type:   "delegate",
					Data: map[string]any{
						"amount":     "10000utia",
						"new_shares": "10000.000000000000000000",
						"validator":  "celestiavaloper1cs37tvmahavw8xcnzcgyz342sh0al37ma4zqat",
					},
				}, {
					Height: 2432340,
					Type:   "message",
					Data: map[string]any{
						"module": "staking",
						"sender": "celestia1lumaemzetgx7fa8gcw3waagevs6rqhejfwxdvszzxlneselshwxqm5rmxn",
					},
				}, {
					Height: 2432340,
					Type:   "ics27_packet",
					Data: map[string]any{
						"host_channel_id": "channel-44",
						"module":          "interchainaccounts",
						"success":         "true",
					},
				}, {
					Height: 2432340,
					Type:   "write_acknowledgement",
					Data: map[string]any{
						"packet_ack":               "{\"result\":\"Ei0KKy9jb3Ntb3Muc3Rha2luZy52MWJldGExLk1zZ0RlbGVnYXRlUmVzcG9uc2U=\"}",
						"packet_ack_hex":           "7b22726573756c74223a224569304b4b79396a62334e7462334d756333526861326c755a7935324d574a6c644745784c6b317a5a30526c6247566e5958526c556d567a634739756332553d227d",
						"packet_connection":        "connection-60",
						"packet_data":              "{\"data\":\"CrQBCiMvY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dEZWxlZ2F0ZRKMAQpDY2VsZXN0aWExbHVtYWVtemV0Z3g3ZmE4Z2N3M3dhYWdldnM2cnFoZWpmd3hkdnN6enhsbmVzZWxzaHd4cW01cm14bhI2Y2VsZXN0aWF2YWxvcGVyMWNzMzd0dm1haGF2dzh4Y256Y2d5ejM0MnNoMGFsMzdtYTR6cWF0Gg0KBHV0aWESBTEwMDAw\",\"memo\":\"perf/celestiavaloper1cs37tvmahavw8xcnzcgyz342sh0al37ma4zqat\",\"type\":\"TYPE_EXECUTE_TX\"}",
						"packet_data_hex":          "7b2264617461223a224372514243694d765932397a6257397a4c6e4e3059577470626d6375646a46695a5852684d53354e633264455a57786c5a3246305a524b4d41517044593256735a584e30615745786248567459575674656d56305a3367335a6d45345a324e334d3364685957646c646e4d32636e466f5a57706d6433686b646e4e36656e6873626d567a5a57787a6148643463573031636d313462684932593256735a584e306157463259577876634756794d574e7a4d7a6430646d316861474632647a68345932353659326435656a4d304d6e4e6f4d4746734d7a647459545236635746304767304b4248563061574553425445774d444177222c226d656d6f223a22706572662f63656c657374696176616c6f706572316373333774766d61686176773878636e7a6367797a333432736830616c33376d61347a716174222c2274797065223a22545950455f455845435554455f5458227d",
						"packet_dst_channel":       "channel-44",
						"packet_dst_port":          "icahost",
						"packet_sequence":          "15",
						"packet_src_channel":       "channel-210",
						"packet_src_port":          "icacontroller-celestia.performance",
						"packet_timeout_height":    "0-0",
						"packet_timeout_timestamp": "1727402253142295584",
					},
				}, {
					Height: 2432340,
					Type:   "message",
					Data: map[string]any{
						"module": "ibc_channel",
					},
				},
			},
			msg: []*storage.Message{
				{
					Type:   types.MsgRecvPacket,
					Height: 2432340,
					Data: map[string]any{
						"Packet": map[string]any{
							"Data": map[string]any{
								"Memo": "perf/celestiavaloper133t4gpv4vhpqgfn9gr8l4u423zrglg8rkqeupr",
								"Type": icaTypes.EXECUTE_TX,
								"Data": []cosmosTypes.Msg{
									&cosmosStakingTypes.MsgDelegate{
										DelegatorAddress: "celestia1lumaemzetgx7fa8gcw3waagevs6rqhejfwxdvszzxlneselshwxqm5rmxn",
										ValidatorAddress: "celestiavaloper133t4gpv4vhpqgfn9gr8l4u423zrglg8rkqeupr",
										Amount:           cosmosTypes.NewCoin("utia", math.NewInt(10000)),
									},
								},
							},
							"DestinationChannel": "channel-44",
							"DestinationPort":    "icahost",
							"Sequence":           14,
							"SourceChannel":      "channel-210",
							"SourcePort":         "icacontroller-celestia.performance",
							"TimeoutHeight": map[string]any{
								"RevisionHeight": 0,
								"RevisionNumber": 0,
							},
							"TimeoutTimestamp": 1727402253142295600,
						},
						"ProofHeight": map[string]any{
							"RevisionHeight": 9291221,
							"RevisionNumber": 2,
						},
						"Signer": "celestia1cts5d9a32lxprwvaw9xt00qnkvndhadp93lwjp",
					},
				}, {
					Type:   types.MsgRecvPacket,
					Height: 2432340,
					Data: map[string]any{
						"Packet": map[string]any{
							"Data": map[string]any{
								"Memo": "perf/celestiavaloper1cs37tvmahavw8xcnzcgyz342sh0al37ma4zqat",
								"Type": icaTypes.EXECUTE_TX,
								"Data": []cosmosTypes.Msg{
									&cosmosStakingTypes.MsgDelegate{
										DelegatorAddress: "celestia1lumaemzetgx7fa8gcw3waagevs6rqhejfwxdvszzxlneselshwxqm5rmxn",
										ValidatorAddress: "celestiavaloper1cs37tvmahavw8xcnzcgyz342sh0al37ma4zqat",
										Amount:           cosmosTypes.NewCoin("utia", math.NewInt(10000)),
									},
								},
							},
							"DestinationChannel": "channel-44",
							"DestinationPort":    "icahost",
							"Sequence":           15,
							"SourceChannel":      "channel-210",
							"SourcePort":         "icacontroller-celestia.performance",
							"TimeoutHeight": map[string]any{
								"RevisionHeight": 0,
								"RevisionNumber": 0,
							},
							"TimeoutTimestamp": 1727402253142295600,
						},
						"ProofHeight": map[string]any{
							"RevisionHeight": 9291221,
							"RevisionNumber": 2,
						},
						"Signer": "celestia1cts5d9a32lxprwvaw9xt00qnkvndhadp93lwjp",
					},
				},
			},
			idx: testsuite.Ptr(0),
		}, {
			name: "recv packet test 5",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 2476842,
					Type:   "message",
					Data: map[string]any{
						"action": "/ibc.core.channel.v1.MsgRecvPacket",
					},
				}, {
					Height: 2476842,
					Type:   "recv_packet",
					Data: map[string]any{
						"packet_channel_ordering":  "ORDER_ORDERED",
						"packet_connection":        "connection-57",
						"packet_data":              "{\"data\":\"CrcBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEpYBCkNjZWxlc3RpYTFuNWp1OG5xNHQ0NHRrbWEyYW1lNnhkZnN4MnRjdG03eDR4dXd4MG1keXRjenJhNHJ0Z2hzbDVlczVkEkNjZWxlc3RpYTE2OXl1cDU2eWR4MmE0MHhrZDcydzJ5cHFqanpnbWw2MnIzMDJ4bWxlN3QwMjY0NGFyazZxN3kweGU1GgoKBHV0aWESAjEzCrgBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEpcBCkNjZWxlc3RpYTFuNWp1OG5xNHQ0NHRrbWEyYW1lNnhkZnN4MnRjdG03eDR4dXd4MG1keXRjenJhNHJ0Z2hzbDVlczVkEkNjZWxlc3RpYTEyN2RocGhuMHM5cXZxODZxM3dxOTVseXM3Zm1hanlod3RkbDNlcjB3MDk0aGpoc3NjNjJzOTQ3cXNrGgsKBHV0aWESAzExNw==\",\"memo\":\"\",\"type\":\"TYPE_EXECUTE_TX\"}",
						"packet_data_hex":          "7b2264617461223a2243726342436877765932397a6257397a4c6d4a68626d7375646a46695a5852684d53354e633264545a57356b45705942436b4e6a5a57786c63335270595446754e5770314f4735784e4851304e485272625745795957316c4e6e686b5a6e4e344d6e526a6447303365445234645864344d47316b6558526a656e4a684e484a305a32687a6244566c637a566b456b4e6a5a57786c63335270595445324f586c3163445532655752344d6d45304d4868725a446379647a4a3563484671616e706e625777324d6e497a4d444a346257786c4e3351774d6a59304e474679617a5a784e336b776547553147676f4b4248563061574553416a457a43726742436877765932397a6257397a4c6d4a68626d7375646a46695a5852684d53354e633264545a57356b45706342436b4e6a5a57786c63335270595446754e5770314f4735784e4851304e485272625745795957316c4e6e686b5a6e4e344d6e526a6447303365445234645864344d47316b6558526a656e4a684e484a305a32687a6244566c637a566b456b4e6a5a57786c63335270595445794e32526f634768754d484d3563585a784f445a784d3364784f54567365584d335a6d3168616e6c6f6433526b62444e6c636a42334d446b306147706f63334e6a4e6a4a7a4f54513363584e724767734b4248563061574553417a45784e773d3d222c226d656d6f223a22222c2274797065223a22545950455f455845435554455f5458227d",
						"packet_dst_channel":       "channel-46",
						"packet_dst_port":          "icahost",
						"packet_sequence":          "1",
						"packet_src_channel":       "channel-27",
						"packet_src_port":          "icacontroller-reward-utia",
						"packet_timeout_height":    "0-0",
						"packet_timeout_timestamp": "1727907245329855685",
					},
				}, {
					Height: 2476842,
					Type:   "message",
					Data: map[string]any{
						"module": "ibc_channel",
					},
				}, {
					Height: 2476842,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "13utia",
						"spender": "celestia1n5ju8nq4t44tkma2ame6xdfsx2tctm7x4xuwx0mdytczra4rtghsl5es5d",
					},
				}, {
					Height: 2476842,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "13utia",
						"receiver": "celestia169yup56ydx2a40xkd72w2ypqjjzgml62r302xmle7t02644ark6q7y0xe5",
					},
				}, {
					Height: 2476842,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "13utia",
						"recipient": "celestia169yup56ydx2a40xkd72w2ypqjjzgml62r302xmle7t02644ark6q7y0xe5",
						"sender":    "celestia1n5ju8nq4t44tkma2ame6xdfsx2tctm7x4xuwx0mdytczra4rtghsl5es5d",
					},
				}, {
					Height: 2476842,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1n5ju8nq4t44tkma2ame6xdfsx2tctm7x4xuwx0mdytczra4rtghsl5es5d",
					},
				}, {
					Height: 2476842,
					Type:   "message",
					Data: map[string]any{
						"module": "bank",
					},
				}, {
					Height: 2476842,
					Type:   "coin_spent",
					Data: map[string]any{
						"amount":  "117utia",
						"spender": "celestia1n5ju8nq4t44tkma2ame6xdfsx2tctm7x4xuwx0mdytczra4rtghsl5es5d",
					},
				}, {
					Height: 2476842,
					Type:   "coin_received",
					Data: map[string]any{
						"amount":   "117utia",
						"receiver": "celestia127dhphn0s9qvq86q3wq95lys7fmajyhwtdl3er0w094hjhssc62s947qsk",
					},
				}, {
					Height: 2476842,
					Type:   "transfer",
					Data: map[string]any{
						"amount":    "117utia",
						"recipient": "celestia127dhphn0s9qvq86q3wq95lys7fmajyhwtdl3er0w094hjhssc62s947qsk",
						"sender":    "celestia1n5ju8nq4t44tkma2ame6xdfsx2tctm7x4xuwx0mdytczra4rtghsl5es5d",
					},
				}, {
					Height: 2476842,
					Type:   "message",
					Data: map[string]any{
						"sender": "celestia1n5ju8nq4t44tkma2ame6xdfsx2tctm7x4xuwx0mdytczra4rtghsl5es5d",
					},
				}, {
					Height: 2476842,
					Type:   "message",
					Data: map[string]any{
						"module": "bank",
					},
				}, {
					Height: 2476842,
					Type:   "ics27_packet",
					Data: map[string]any{
						"host_channel_id": "channel-46",
						"module":          "interchainaccounts",
						"success":         "true",
					},
				}, {
					Height: 2476842,
					Type:   "write_acknowledgement",
					Data: map[string]any{
						"packet_ack":               "{\"result\":\"EiYKJC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmRSZXNwb25zZRImCiQvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kUmVzcG9uc2U=\"}",
						"packet_ack_hex":           "7b22726573756c74223a224569594b4a43396a62334e7462334d75596d4675617935324d574a6c644745784c6b317a5a314e6c626d52535a584e776232357a5a52496d436951765932397a6257397a4c6d4a68626d7375646a46695a5852684d53354e633264545a57356b556d567a634739756332553d227d",
						"packet_connection":        "connection-57",
						"packet_data":              "{\"data\":\"CrcBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEpYBCkNjZWxlc3RpYTFuNWp1OG5xNHQ0NHRrbWEyYW1lNnhkZnN4MnRjdG03eDR4dXd4MG1keXRjenJhNHJ0Z2hzbDVlczVkEkNjZWxlc3RpYTE2OXl1cDU2eWR4MmE0MHhrZDcydzJ5cHFqanpnbWw2MnIzMDJ4bWxlN3QwMjY0NGFyazZxN3kweGU1GgoKBHV0aWESAjEzCrgBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEpcBCkNjZWxlc3RpYTFuNWp1OG5xNHQ0NHRrbWEyYW1lNnhkZnN4MnRjdG03eDR4dXd4MG1keXRjenJhNHJ0Z2hzbDVlczVkEkNjZWxlc3RpYTEyN2RocGhuMHM5cXZxODZxM3dxOTVseXM3Zm1hanlod3RkbDNlcjB3MDk0aGpoc3NjNjJzOTQ3cXNrGgsKBHV0aWESAzExNw==\",\"memo\":\"\",\"type\":\"TYPE_EXECUTE_TX\"}",
						"packet_data_hex":          "7b2264617461223a2243726342436877765932397a6257397a4c6d4a68626d7375646a46695a5852684d53354e633264545a57356b45705942436b4e6a5a57786c63335270595446754e5770314f4735784e4851304e485272625745795957316c4e6e686b5a6e4e344d6e526a6447303365445234645864344d47316b6558526a656e4a684e484a305a32687a6244566c637a566b456b4e6a5a57786c63335270595445324f586c3163445532655752344d6d45304d4868725a446379647a4a3563484671616e706e625777324d6e497a4d444a346257786c4e3351774d6a59304e474679617a5a784e336b776547553147676f4b4248563061574553416a457a43726742436877765932397a6257397a4c6d4a68626d7375646a46695a5852684d53354e633264545a57356b45706342436b4e6a5a57786c63335270595446754e5770314f4735784e4851304e485272625745795957316c4e6e686b5a6e4e344d6e526a6447303365445234645864344d47316b6558526a656e4a684e484a305a32687a6244566c637a566b456b4e6a5a57786c63335270595445794e32526f634768754d484d3563585a784f445a784d3364784f54567365584d335a6d3168616e6c6f6433526b62444e6c636a42334d446b306147706f63334e6a4e6a4a7a4f54513363584e724767734b4248563061574553417a45784e773d3d222c226d656d6f223a22222c2274797065223a22545950455f455845435554455f5458227d",
						"packet_dst_channel":       "channel-46",
						"packet_dst_port":          "icahost",
						"packet_sequence":          "1",
						"packet_src_channel":       "channel-27",
						"packet_src_port":          "icacontroller-reward-utia",
						"packet_timeout_height":    "0-0",
						"packet_timeout_timestamp": "1727907245329855685",
					},
				}, {
					Height: 2476842,
					Type:   "message",
					Data: map[string]any{
						"module": "ibc_channel",
					},
				},
			},
			msg: []*storage.Message{
				{
					Type:   types.MsgRecvPacket,
					Height: 2476842,
					Data: map[string]any{
						"Packet": map[string]any{
							"Data": map[string]any{
								"Memo": "",
								"Type": icaTypes.EXECUTE_TX,
								"Data": []cosmosTypes.Msg{
									&cosmosBankTypes.MsgSend{
										FromAddress: "celestia1n5ju8nq4t44tkma2ame6xdfsx2tctm7x4xuwx0mdytczra4rtghsl5es5d",
										ToAddress:   "celestia169yup56ydx2a40xkd72w2ypqjjzgml62r302xmle7t02644ark6q7y0xe5",
										Amount:      cosmosTypes.NewCoins(cosmosTypes.NewCoin("utia", math.NewInt(113))),
									},
									&cosmosBankTypes.MsgSend{
										FromAddress: "celestia1n5ju8nq4t44tkma2ame6xdfsx2tctm7x4xuwx0mdytczra4rtghsl5es5d",
										ToAddress:   "celestia127dhphn0s9qvq86q3wq95lys7fmajyhwtdl3er0w094hjhssc62s947qsk",
										Amount:      cosmosTypes.NewCoins(cosmosTypes.NewCoin("utia", math.NewInt(117))),
									},
								},
							},
							"DestinationChannel": "channel-44",
							"DestinationPort":    "icahost",
							"Sequence":           14,
							"SourceChannel":      "channel-210",
							"SourcePort":         "icacontroller-celestia.performance",
							"TimeoutHeight": map[string]any{
								"RevisionHeight": 0,
								"RevisionNumber": 0,
							},
							"TimeoutTimestamp": 1727402253142295600,
						},
						"ProofHeight": map[string]any{
							"RevisionHeight": 2243980,
							"RevisionNumber": 1,
						},
						"Signer": "celestia1cla57nqkvv9c76744lkyzrzael536gxf7uxnmy",
					},
				},
			},
			idx: testsuite.Ptr(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ctx.Block = &storage.Block{
				Height: tt.events[0].Height,
				Time:   time.Now().UTC(),
			}
			for i := range tt.msg {
				err := handleRecvPacket(tt.ctx, tt.events, tt.msg[i], tt.idx)
				require.NoError(t, err)
				if tt.msg[i].IbcTransfer != nil {
					require.NotEmpty(t, tt.msg[i].IbcTransfer.ConnectionId)
				}
			}
		})
	}
}
