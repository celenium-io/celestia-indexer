// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"testing"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/mock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestModule_saveSigners(t *testing.T) {
	type args struct {
		addrToId map[string]uint64
		txs      []storage.Tx
	}
	tests := []struct {
		name    string
		args    args
		want    []storage.Signer
		wantErr bool
	}{
		{
			name: "test 1",
			args: args{
				addrToId: map[string]uint64{
					"address1": 1,
					"address2": 2,
				},
				txs: []storage.Tx{
					{
						Id: 1,
						Signers: []storage.Address{
							{
								Address: "address1",
							}, {
								Address: "address2",
							},
						},
					},
				},
			},
			want: []storage.Signer{
				{
					TxId:      1,
					AddressId: 1,
				}, {
					TxId:      1,
					AddressId: 2,
				},
			},
			wantErr: false,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		tx := mock.NewMockTransaction(ctrl)
		tx.EXPECT().
			SaveSigners(gomock.Any(), gomock.Any()).
			MaxTimes(1).
			MinTimes(1).
			DoAndReturn(func(_ context.Context, addresses ...storage.Signer) error {
				require.Equal(t, tt.want, addresses)
				return nil
			})

		t.Run(tt.name, func(t *testing.T) {
			err := saveSigners(context.Background(), tx, tt.args.addrToId, tt.args.txs)
			require.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func Test_saveAddresses(t *testing.T) {
	tests := []struct {
		name      string
		addresses map[string]*storage.Address
		addr      map[string]uint64
		total     int64
		wantErr   bool
	}{
		{
			name:      "test 1",
			addresses: map[string]*storage.Address{},
			addr:      nil,
			total:     0,
			wantErr:   false,
		}, {
			name: "test 2",
			addresses: map[string]*storage.Address{
				"address1": {
					Address:    "address1",
					Height:     100,
					LastHeight: 100,
					Balance: storage.Balance{
						Currency: "utia",
						Total:    decimal.RequireFromString("1"),
					},
				},
			},
			addr: map[string]uint64{
				"address1": 1,
			},
			total:   1,
			wantErr: false,
		}, {
			name: "test 3",
			addresses: map[string]*storage.Address{
				"address1": {
					Address:    "address1",
					Height:     100,
					LastHeight: 101,
					Balance: storage.Balance{
						Currency: "utia",
						Total:    decimal.RequireFromString("1"),
					},
				},
			},
			addr: map[string]uint64{
				"address1": 1,
			},
			total:   0,
			wantErr: false,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		tx := mock.NewMockTransaction(ctrl)

		tx.EXPECT().
			SaveAddresses(gomock.Any(), gomock.Any()).
			MaxTimes(1).
			MinTimes(0).
			DoAndReturn(func(_ context.Context, addresses ...*storage.Address) (int64, error) {
				require.Equal(t, len(tt.addresses), len(addresses))
				var count int64
				for i := range addresses {
					addresses[i].Id = uint64(i + 1)

					if addresses[i].Height == addresses[i].LastHeight {
						count++
					}
				}
				return count, nil
			})

		tx.EXPECT().
			SaveBalances(gomock.Any(), gomock.Any()).
			MaxTimes(1).
			MinTimes(0).
			DoAndReturn(func(_ context.Context, balances ...storage.Balance) error {
				require.Equal(t, len(tt.addresses), len(balances))
				return nil
			})

		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := saveAddresses(context.Background(), tx, tt.addresses)
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.addr, got)
			require.Equal(t, tt.total, got1)
		})
	}
}
