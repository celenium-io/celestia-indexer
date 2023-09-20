package rollback

import (
	"context"
	"testing"

	"github.com/dipdup-io/celestia-indexer/internal/consts"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/mock"
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	testAddress     = "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60"
	testHashAddress = []byte{0x96, 0xa, 0xa0, 0x36, 0x6b, 0x25, 0x4e, 0x1e, 0xa7, 0x9b, 0xda, 0x46, 0x7e, 0xb3, 0xaa, 0x5c, 0x97, 0xcb, 0xa5, 0xae}
)

func Test_coinReceived(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]any
		want    *storage.Address
		wantErr bool
	}{
		{
			name: "test 1",
			data: map[string]any{
				"receiver": testAddress,
				"amount":   "123utia",
			},
			want: &storage.Address{
				Hash:    testHashAddress,
				Address: testAddress,
				Balance: storage.Balance{
					Currency: "utia",
					Total:    decimal.RequireFromString("-123"),
				},
			},
		}, {
			name: "test 2",
			data: map[string]any{
				"receiver": testAddress,
				"amount":   nil,
			},
			want: &storage.Address{
				Hash:    testHashAddress,
				Address: testAddress,
				Balance: storage.Balance{
					Currency: "utia",
					Total:    decimal.Zero,
				},
			},
		}, {
			name: "test 3",
			data: map[string]any{
				"receiver": "invalid",
				"amount":   nil,
			},
			wantErr: true,
		}, {
			name:    "test 4",
			data:    nil,
			wantErr: true,
		}, {
			name: "test 5",
			data: map[string]any{
				"receiver": "",
				"amount":   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := coinReceived(tt.data)
			require.Equal(t, tt.wantErr, err != nil)
			if err == nil {
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_coinSpent(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]any
		want    *storage.Address
		wantErr bool
	}{
		{
			name: "test 1",
			data: map[string]any{
				"spender": testAddress,
				"amount":  "123utia",
			},
			want: &storage.Address{
				Hash:    testHashAddress,
				Address: testAddress,
				Balance: storage.Balance{
					Currency: "utia",
					Total:    decimal.RequireFromString("123"),
				},
			},
		}, {
			name: "test 2",
			data: map[string]any{
				"spender": testAddress,
				"amount":  nil,
			},
			want: &storage.Address{
				Hash:    testHashAddress,
				Address: testAddress,
				Balance: storage.Balance{
					Currency: "utia",
					Total:    decimal.Zero,
				},
			},
		}, {
			name: "test 3",
			data: map[string]any{
				"spender": "invalid",
				"amount":  nil,
			},
			wantErr: true,
		}, {
			name:    "test 4",
			data:    nil,
			wantErr: true,
		}, {
			name: "test 5",
			data: map[string]any{
				"spender": "",
				"amount":  nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := coinSpent(tt.data)
			require.Equal(t, tt.wantErr, err != nil)
			if err == nil {
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_getBalanceUpdates(t *testing.T) {
	type args struct {
		deletedAddress map[string]struct{}
		deletedEvents  []storage.Event
	}
	tests := []struct {
		name    string
		args    args
		want    []*storage.Address
		wantErr bool
	}{
		{
			name: "test 1",
			args: args{
				deletedAddress: map[string]struct{}{},
				deletedEvents: []storage.Event{
					{
						Type: types.EventTypeCoinSpent,
						Data: map[string]any{
							"spender": testAddress,
							"amount":  "123utia",
						},
					}, {
						Type: types.EventTypeCoinReceived,
						Data: map[string]any{
							"receiver": testAddress,
							"amount":   "23utia",
						},
					},
				},
			},
			want: []*storage.Address{
				{
					Address: testAddress,
					Hash:    testHashAddress,
					Balance: storage.Balance{
						Currency: consts.DefaultCurrency,
						Total:    decimal.RequireFromString("100"),
					},
					LastHeight: 100,
				},
			},
		}, {
			name: "test 2",
			args: args{
				deletedAddress: map[string]struct{}{
					testAddress: {},
				},
				deletedEvents: []storage.Event{
					{
						Type: types.EventTypeCoinSpent,
						Data: map[string]any{
							"spender": testAddress,
							"amount":  "123utia",
						},
					}, {
						Type: types.EventTypeCoinReceived,
						Data: map[string]any{
							"receiver": testAddress,
							"amount":   "23utia",
						},
					},
				},
			},
			want: []*storage.Address{},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	tx.EXPECT().LastAddressAction(gomock.Any(), gomock.Any()).
		Return(100, nil).
		AnyTimes()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBalanceUpdates(context.Background(), tx, tt.args.deletedAddress, tt.args.deletedEvents)
			require.Equal(t, tt.wantErr, err != nil)
			if err == nil {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
