// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rollback

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
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
		data    map[string]string
		want    *storage.Address
		wantErr bool
	}{
		{
			name: "test 1",
			data: map[string]string{
				"receiver": testAddress,
				"amount":   "123utia",
			},
			want: &storage.Address{
				Hash:    testHashAddress,
				Address: testAddress,
				Balance: storage.Balance{
					Currency:  "utia",
					Spendable: types.NewNumeric(decimal.RequireFromString("-123")),
				},
			},
		}, {
			name: "test 2",
			data: map[string]string{
				"receiver": testAddress,
				"amount":   "",
			},
			want: &storage.Address{
				Hash:    testHashAddress,
				Address: testAddress,
				Balance: storage.Balance{
					Currency:  "utia",
					Spendable: types.NewNumeric(decimal.Zero),
				},
			},
		}, {
			name: "test 3",
			data: map[string]string{
				"receiver": "invalid",
				"amount":   "",
			},
			wantErr: true,
		}, {
			name:    "test 4",
			data:    nil,
			wantErr: true,
		}, {
			name: "test 5",
			data: map[string]string{
				"receiver": "",
				"amount":   "",
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
		data    map[string]string
		want    *storage.Address
		wantErr bool
	}{
		{
			name: "test 1",
			data: map[string]string{
				"spender": testAddress,
				"amount":  "123utia",
			},
			want: &storage.Address{
				Hash:    testHashAddress,
				Address: testAddress,
				Balance: storage.Balance{
					Currency:  "utia",
					Spendable: types.NewNumeric(decimal.RequireFromString("123")),
				},
			},
		}, {
			name: "test 2",
			data: map[string]string{
				"spender": testAddress,
				"amount":  "",
			},
			want: &storage.Address{
				Hash:    testHashAddress,
				Address: testAddress,
				Balance: storage.Balance{
					Currency:  "utia",
					Spendable: types.NewNumeric(decimal.Zero),
				},
			},
		}, {
			name: "test 3",
			data: map[string]string{
				"spender": "invalid",
				"amount":  "",
			},
			wantErr: true,
		}, {
			name:    "test 4",
			data:    nil,
			wantErr: true,
		}, {
			name: "test 5",
			data: map[string]string{
				"spender": "",
				"amount":  "",
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
						Data: map[string]string{
							"spender": testAddress,
							"amount":  "123utia",
						},
					}, {
						Type: types.EventTypeCoinReceived,
						Data: map[string]string{
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
						Currency:  currency.DefaultCurrency,
						Spendable: types.NewNumeric(decimal.RequireFromString("100")),
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
						Data: map[string]string{
							"spender": testAddress,
							"amount":  "123utia",
						},
					}, {
						Type: types.EventTypeCoinReceived,
						Data: map[string]string{
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
			got, err := getBalanceUpdates(t.Context(), tx, tt.args.deletedAddress, tt.args.deletedEvents)
			require.Equal(t, tt.wantErr, err != nil)
			if err == nil {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
