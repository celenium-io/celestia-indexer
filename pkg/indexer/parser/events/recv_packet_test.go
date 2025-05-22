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

func Test_handleRecvPacket(t *testing.T) {
	tests := []struct {
		name   string
		ctx    *context.Context
		events []storage.Event
		msg    *storage.Message
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
			msg: &storage.Message{
				Type:   types.MsgRecvPacket,
				Height: 1866988,
				Data: map[string]any{
					"Packet": map[string]any{
						"Data":               "eyJkYXRhIjoiQ3FJQkNod3ZZMjl6Ylc5ekxtSmhibXN1ZGpGaVpYUmhNUzVOYzJkVFpXNWtFb0VCQ2tOalpXeGxjM1JwWVRGNGVuTmtialkxYUhsc2FtTnRaVzVzZUhscWJXUnRkbWRvYUdRd2R6UjFkREkzYXpObWVEVTJhbkF5Y0RZNVpXZzJjM0p6T0hBemNuTnpFaTlqWld4bGMzUnBZVEZrYzIxdWVubzRjWFp1TXpRemJIaHhNbVYzTmpjek9YTmtNREpuYzNOd2VXVnpZVFJxZFJvSkNnUjFkR2xoRWdFeCIsIm1lbW8iOiJyb290dWxwIiwidHlwZSI6IlRZUEVfRVhFQ1VURV9UWCJ9",
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
				},
			},
			msg: &storage.Message{
				Type:   types.MsgRecvPacket,
				Height: 1866988,
				Data: map[string]any{
					"Packet": map[string]any{
						"Data":               "eyJhbW91bnQiOiI0MDAwMDAwIiwiZGVub20iOiJ0cmFuc2Zlci9jaGFubmVsLTM1L3V0aWEiLCJtZW1vIjoie1wiZm9yd2FyZFwiOntcInJlY2VpdmVyXCI6XCJvc21vMXZrZGFrcXFnNWh0cTVjM3d5MmtqMmdlcTUzNnE2NjV4ZGV4cnRqdXdxY2twYWRzMmMybnN2aGhjeXZcIixcInBvcnRcIjpcInRyYW5zZmVyXCIsXCJjaGFubmVsXCI6XCJjaGFubmVsLTJcIixcInRpbWVvdXRcIjowLFwicmV0cmllc1wiOjIsXCJuZXh0XCI6e1wid2FzbVwiOntcImNvbnRyYWN0XCI6XCJvc21vMXZrZGFrcXFnNWh0cTVjM3d5MmtqMmdlcTUzNnE2NjV4ZGV4cnRqdXdxY2twYWRzMmMybnN2aGhjeXZcIixcIm1zZ1wiOntcInN3YXBfYW5kX2FjdGlvblwiOntcInVzZXJfc3dhcFwiOntcInN3YXBfZXhhY3RfYXNzZXRfaW5cIjp7XCJzd2FwX3ZlbnVlX25hbWVcIjpcIm9zbW9zaXMtcG9vbG1hbmFnZXJcIixcIm9wZXJhdGlvbnNcIjpbe1wicG9vbFwiOlwiMTQ3NVwiLFwiZGVub21faW5cIjpcImliYy9ENzlFN0Q4M0FCMzk5QkZGRjkzNDMzRTU0RkFBNDgwQzE5MTI0OEZDNTU2OTI0QTJBODM1MUFFMjYzOEIzODc3XCIsXCJkZW5vbV9vdXRcIjpcImZhY3Rvcnkvb3NtbzFmNXZmY3BoMmR2ZmVxY3FraGV0d3Y3NWZkYTY5ejdlNWMyZGxkbTNrdmdqMjNjcmt2NndxY240N2EwL3VtaWxrVElBXCJ9LHtcInBvb2xcIjpcIjE2OTRcIixcImRlbm9tX2luXCI6XCJmYWN0b3J5L29zbW8xZjV2ZmNwaDJkdmZlcWNxa2hldHd2NzVmZGE2OXo3ZTVjMmRsZG0za3ZnajIzY3JrdjZ3cWNuNDdhMC91bWlsa1RJQVwiLFwiZGVub21fb3V0XCI6XCJpYmMvNjkxMTBGRjY3M0Q3MEIzOTkwNEZGMDU2Q0ZERkQ1OEE5MEJFQzMxOTQzMDNGNDVDMzJDQjkxQjhCMEE3MzhFQVwifSx7XCJwb29sXCI6XCIxNjk4XCIsXCJkZW5vbV9pblwiOlwiaWJjLzY5MTEwRkY2NzNENzBCMzk5MDRGRjA1NkNGREZENThBOTBCRUMzMTk0MzAzRjQ1QzMyQ0I5MUI4QjBBNzM4RUFcIixcImRlbm9tX291dFwiOlwiaWJjLzY0QkE2RTMxRkU4ODdENjZDNkY4RjMxQzdCMUE4MEM3Q0ExNzkyMzk2NzdCNDA4OEJCNTVGNUVBMDdEQkUyNzNcIn1dfX0sXCJtaW5fYXNzZXRcIjp7XCJuYXRpdmVcIjp7XCJkZW5vbVwiOlwiaWJjLzY0QkE2RTMxRkU4ODdENjZDNkY4RjMxQzdCMUE4MEM3Q0ExNzkyMzk2NzdCNDA4OEJCNTVGNUVBMDdEQkUyNzNcIixcImFtb3VudFwiOlwiMTA5NTExNjU5OTI0ODA2MzA5MlwifX0sXCJ0aW1lb3V0X3RpbWVzdGFtcFwiOjE3MjY2NzM5NjIyMzM0MzMwNjcsXCJwb3N0X3N3YXBfYWN0aW9uXCI6e1wiaWJjX3RyYW5zZmVyXCI6e1wiaWJjX2luZm9cIjp7XCJzb3VyY2VfY2hhbm5lbFwiOlwiY2hhbm5lbC0xMjJcIixcInJlY2VpdmVyXCI6XCJpbmoxc3lyZGgydjJyd2Y4Zmh2czVnc213cjRhaGNuczNtNHp2Y3Vrd2NcIixcIm1lbW9cIjpcIlwiLFwicmVjb3Zlcl9hZGRyZXNzXCI6XCJvc21vMW04NG5oNzVobDQ3NGs1ZDgzY3VucnFmeGdtcmw1MjN1NXEwOXh3XCJ9fX0sXCJhZmZpbGlhdGVzXCI6W3tcImJhc2lzX3BvaW50c19mZWVcIjpcIjYwXCIsXCJhZGRyZXNzXCI6XCJvc21vMW15NHRrNDIwZ2ptaGdncXd2dmhhNmV5OTM5MGdxd2ZyZWUycDR1XCJ9LHtcImJhc2lzX3BvaW50c19mZWVcIjpcIjE1XCIsXCJhZGRyZXNzXCI6XCJvc21vMW1zam5hbDJnbGZ6NnplOHg5a2R1aGc0NXhwcHg5c2RkYXdmeDQ2XCJ9XX19fX19fSIsInJlY2VpdmVyIjoiY2VsZXN0aWExbTg0bmg3NWhsNDc0azVkODNjdW5ycWZ4Z21ybDUyM3VkM2Q5MjMiLCJzZW5kZXIiOiJuZXV0cm9uMW04NG5oNzVobDQ3NGs1ZDgzY3VucnFmeGdtcmw1MjN1Y3k0aDJtIn0=",
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
			},
			idx: testsuite.Ptr(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleRecvPacket(tt.ctx, tt.events, tt.msg, tt.idx)
			require.NoError(t, err)
		})
	}
}
