// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package websocket

import (
	"testing"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/stretchr/testify/require"
)

func TestTxFilter_Filter(t *testing.T) {
	tests := []struct {
		name string
		c    *Client
		msg  responses.Tx
		want bool
	}{
		{
			name: "test 1",
			c: &Client{
				filters: &Filters{},
			},
			msg: responses.Tx{
				Status: types.StatusSuccess,
				MessageTypes: []types.MsgType{
					types.MsgBeginRedelegate,
				},
			},
			want: false,
		}, {
			name: "test 2",
			c: &Client{
				filters: &Filters{
					tx: &txFilters{
						msgs: types.NewMsgTypeBitMask(types.MsgBeginRedelegate),
					},
				},
			},
			msg: responses.Tx{
				Status:      types.StatusSuccess,
				MsgTypeMask: types.NewMsgTypeBitMask(types.MsgBeginRedelegate),
			},
			want: true,
		}, {
			name: "test 3",
			c: &Client{
				filters: &Filters{
					tx: &txFilters{
						msgs: types.NewMsgTypeBitMask(types.MsgSend),
					},
				},
			},
			msg: responses.Tx{
				Status:      types.StatusSuccess,
				MsgTypeMask: types.NewMsgTypeBitMask(types.MsgBeginRedelegate),
			},
			want: false,
		}, {
			name: "test 4",
			c: &Client{
				filters: &Filters{
					tx: &txFilters{
						status: map[types.Status]struct{}{
							types.StatusSuccess: {},
						},
					},
				},
			},
			msg: responses.Tx{
				Status:      types.StatusSuccess,
				MsgTypeMask: types.NewMsgTypeBitMask(types.MsgBeginRedelegate),
			},
			want: true,
		}, {
			name: "test 5",
			c: &Client{
				filters: &Filters{
					tx: &txFilters{
						status: map[types.Status]struct{}{
							types.StatusFailed: {},
						},
					},
				},
			},
			msg: responses.Tx{
				Status:      types.StatusSuccess,
				MsgTypeMask: types.NewMsgTypeBitMask(types.MsgBeginRedelegate),
			},
			want: false,
		}, {
			name: "test 6",
			c: &Client{
				filters: &Filters{
					tx: &txFilters{
						status: map[types.Status]struct{}{
							types.StatusSuccess: {},
						},
						msgs: types.NewMsgTypeBitMask(types.MsgSend),
					},
				},
			},
			msg: responses.Tx{
				Status:      types.StatusSuccess,
				MsgTypeMask: types.NewMsgTypeBitMask(types.MsgSend),
			},
			want: true,
		}, {
			name: "test 7",
			c: &Client{
				filters: &Filters{
					tx: &txFilters{
						status: map[types.Status]struct{}{
							types.StatusFailed: {},
						},
						msgs: types.NewMsgTypeBitMask(types.MsgSend),
					},
				},
			},
			msg: responses.Tx{
				Status:      types.StatusSuccess,
				MsgTypeMask: types.NewMsgTypeBitMask(types.MsgSend),
			},
			want: false,
		}, {
			name: "test 8",
			c: &Client{
				filters: &Filters{
					tx: &txFilters{
						status: map[types.Status]struct{}{
							types.StatusFailed: {},
						},
						msgs: types.NewMsgTypeBitMask(types.MsgBeginRedelegate),
					},
				},
			},
			msg: responses.Tx{
				Status:      types.StatusSuccess,
				MsgTypeMask: types.NewMsgTypeBitMask(types.MsgSend),
			},
			want: false,
		}, {
			name: "test 9",
			c: &Client{
				filters: newFilters(),
			},
			msg: responses.Tx{
				Status:      types.StatusSuccess,
				MsgTypeMask: types.NewMsgTypeBitMask(types.MsgSend),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(TxFilter).Filter(tt.c, &tt.msg)
			require.Equal(t, tt.want, got)
		})
	}
}
