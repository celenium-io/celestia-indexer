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
				filters: &filters{},
			},
			msg: responses.Tx{
				Status: string(types.StatusSuccess),
				MessageTypes: []string{
					string(types.MsgTypeBeginRedelegate),
				},
			},
			want: false,
		}, {
			name: "test 2",
			c: &Client{
				filters: &filters{
					tx: &txFilters{
						msgs: types.NewMsgTypeBitMask(types.MsgTypeBeginRedelegate),
					},
				},
			},
			msg: responses.Tx{
				Status:      string(types.StatusSuccess),
				MsgTypeMask: types.NewMsgTypeBitMask(types.MsgTypeBeginRedelegate),
			},
			want: true,
		}, {
			name: "test 3",
			c: &Client{
				filters: &filters{
					tx: &txFilters{
						msgs: types.NewMsgTypeBitMask(types.MsgTypeSend),
					},
				},
			},
			msg: responses.Tx{
				Status:      string(types.StatusSuccess),
				MsgTypeMask: types.NewMsgTypeBitMask(types.MsgTypeBeginRedelegate),
			},
			want: false,
		}, {
			name: "test 4",
			c: &Client{
				filters: &filters{
					tx: &txFilters{
						status: map[string]struct{}{
							string(types.StatusSuccess): {},
						},
					},
				},
			},
			msg: responses.Tx{
				Status:      string(types.StatusSuccess),
				MsgTypeMask: types.NewMsgTypeBitMask(types.MsgTypeBeginRedelegate),
			},
			want: true,
		}, {
			name: "test 5",
			c: &Client{
				filters: &filters{
					tx: &txFilters{
						status: map[string]struct{}{
							string(types.StatusFailed): {},
						},
					},
				},
			},
			msg: responses.Tx{
				Status:      string(types.StatusSuccess),
				MsgTypeMask: types.NewMsgTypeBitMask(types.MsgTypeBeginRedelegate),
			},
			want: false,
		}, {
			name: "test 6",
			c: &Client{
				filters: &filters{
					tx: &txFilters{
						status: map[string]struct{}{
							string(types.StatusSuccess): {},
						},
						msgs: types.NewMsgTypeBitMask(types.MsgTypeSend),
					},
				},
			},
			msg: responses.Tx{
				Status:      string(types.StatusSuccess),
				MsgTypeMask: types.NewMsgTypeBitMask(types.MsgTypeSend),
			},
			want: true,
		}, {
			name: "test 7",
			c: &Client{
				filters: &filters{
					tx: &txFilters{
						status: map[string]struct{}{
							string(types.StatusFailed): {},
						},
						msgs: types.NewMsgTypeBitMask(types.MsgTypeSend),
					},
				},
			},
			msg: responses.Tx{
				Status:      string(types.StatusSuccess),
				MsgTypeMask: types.NewMsgTypeBitMask(types.MsgTypeSend),
			},
			want: false,
		}, {
			name: "test 8",
			c: &Client{
				filters: &filters{
					tx: &txFilters{
						status: map[string]struct{}{
							string(types.StatusFailed): {},
						},
						msgs: types.NewMsgTypeBitMask(types.MsgTypeBeginRedelegate),
					},
				},
			},
			msg: responses.Tx{
				Status:      string(types.StatusSuccess),
				MsgTypeMask: types.NewMsgTypeBitMask(types.MsgTypeSend),
			},
			want: false,
		}, {
			name: "test 9",
			c: &Client{
				filters: newFilters(),
			},
			msg: responses.Tx{
				Status:      string(types.StatusSuccess),
				MsgTypeMask: types.NewMsgTypeBitMask(types.MsgTypeSend),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(TxFilter).Filter(tt.c, tt.msg)
			require.Equal(t, tt.want, got)
		})
	}
}
