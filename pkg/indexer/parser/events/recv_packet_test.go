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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleRecvPacket(tt.ctx, tt.events, tt.msg, tt.idx)
			require.NoError(t, err)
		})
	}
}
