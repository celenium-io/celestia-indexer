// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_saveMessages(t *testing.T) {
	type args struct {
		messages []*storage.Message
		addrToId map[string]uint64
	}

	now := time.Now()

	tests := []struct {
		name                      string
		args                      args
		wantNamespaceMessageCount int
		wantMsgAddress            int
		wantErr                   bool
	}{
		{
			name: "test without namespaces and validators",
			args: args{
				messages: []*storage.Message{
					{
						Height:   100,
						Time:     now,
						Position: 0,
						Type:     types.MsgSend,
						TxId:     1,
						Addresses: []storage.AddressWithType{
							{
								Type: types.MsgAddressTypeFromAddress,
								Address: storage.Address{
									Address:    "address1",
									Height:     100,
									LastHeight: 100,
								},
							}, {
								Type: types.MsgAddressTypeToAddress,
								Address: storage.Address{
									Address:    "address2",
									Height:     100,
									LastHeight: 100,
								},
							},
						},
					},
				},
				addrToId: map[string]uint64{
					"address1": 1,
					"address2": 2,
				},
			},
			wantNamespaceMessageCount: 0,
			wantMsgAddress:            2,
			wantErr:                   false,
		}, {
			name: "test with namespaces and without validators",
			args: args{
				messages: []*storage.Message{
					{
						Height:   100,
						Time:     now,
						Position: 0,
						Type:     types.MsgSend,
						TxId:     1,
						Addresses: []storage.AddressWithType{
							{
								Type: types.MsgAddressTypeFromAddress,
								Address: storage.Address{
									Address:    "address1",
									Height:     100,
									LastHeight: 100,
								},
							}, {
								Type: types.MsgAddressTypeToAddress,
								Address: storage.Address{
									Address:    "address2",
									Height:     100,
									LastHeight: 100,
								},
							},
						},
					}, {
						Height:   100,
						Time:     now,
						Position: 0,
						Type:     types.MsgPayForBlobs,
						TxId:     2,
						Addresses: []storage.AddressWithType{
							{
								Type: types.MsgAddressTypeSigner,
								Address: storage.Address{
									Address:    "address3",
									Height:     100,
									LastHeight: 100,
								},
							},
						},
						Namespace: []storage.Namespace{
							{
								Id:          1,
								FirstHeight: 100,
								Version:     0,
								NamespaceID: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7},
								Size:        1000,
								PfbCount:    1,
								Reserved:    false,
							}, {
								Id:          2,
								FirstHeight: 100,
								Version:     1,
								NamespaceID: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7},
								Size:        1000,
								PfbCount:    1,
								Reserved:    false,
							},
						},
					},
				},
				addrToId: map[string]uint64{
					"address1": 1,
					"address2": 2,
					"address3": 3,
				},
			},
			wantNamespaceMessageCount: 2,
			wantMsgAddress:            3,
			wantErr:                   false,
		}, {
			name: "test without namespaces and with validators",
			args: args{
				messages: []*storage.Message{
					{
						Height:   100,
						Time:     now,
						Position: 0,
						Type:     types.MsgSend,
						TxId:     1,
						Addresses: []storage.AddressWithType{
							{
								Type: types.MsgAddressTypeFromAddress,
								Address: storage.Address{
									Address:    "address1",
									Height:     100,
									LastHeight: 100,
								},
							}, {
								Type: types.MsgAddressTypeToAddress,
								Address: storage.Address{
									Address:    "address2",
									Height:     100,
									LastHeight: 100,
								},
							},
						},
					}, {
						Height:   100,
						Time:     now,
						Position: 0,
						Type:     types.MsgCreateValidator,
						TxId:     2,
						Addresses: []storage.AddressWithType{
							{
								Type: types.MsgAddressTypeDelegator,
								Address: storage.Address{
									Address:    "address3",
									Height:     100,
									LastHeight: 100,
								},
							}, {
								Type: types.MsgAddressTypeValidator,
								Address: storage.Address{
									Address:    "address1",
									Height:     100,
									LastHeight: 100,
								},
							},
						},
					},
				},
				addrToId: map[string]uint64{
					"address1": 1,
					"address2": 2,
					"address3": 3,
				},
			},
			wantNamespaceMessageCount: 0,
			wantMsgAddress:            4,
			wantErr:                   false,
		}, {
			name: "test with duplicate namespaces",
			args: args{
				messages: []*storage.Message{
					{
						Height:   100,
						Time:     now,
						Position: 0,
						Type:     types.MsgSend,
						TxId:     1,
						Addresses: []storage.AddressWithType{
							{
								Type: types.MsgAddressTypeFromAddress,
								Address: storage.Address{
									Address:    "address1",
									Height:     100,
									LastHeight: 100,
								},
							}, {
								Type: types.MsgAddressTypeToAddress,
								Address: storage.Address{
									Address:    "address2",
									Height:     100,
									LastHeight: 100,
								},
							},
						},
					}, {
						Height:   100,
						Time:     now,
						Position: 0,
						Type:     types.MsgPayForBlobs,
						TxId:     2,
						Addresses: []storage.AddressWithType{
							{
								Type: types.MsgAddressTypeSigner,
								Address: storage.Address{
									Address:    "address3",
									Height:     100,
									LastHeight: 100,
								},
							},
						},
						Namespace: []storage.Namespace{
							{
								Id:          1,
								FirstHeight: 100,
								Version:     0,
								NamespaceID: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7},
								Size:        1000,
								PfbCount:    1,
								Reserved:    false,
							}, {
								FirstHeight: 100,
								Version:     0,
								NamespaceID: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7},
								Size:        1000,
								PfbCount:    1,
								Reserved:    false,
							},
						},
					},
				},
				addrToId: map[string]uint64{
					"address1": 1,
					"address2": 2,
					"address3": 3,
				},
			},
			wantNamespaceMessageCount: 1,
			wantMsgAddress:            3,
			wantErr:                   false,
		}, {
			name: "test with duplicate message addresses",
			args: args{
				messages: []*storage.Message{
					{
						Height:   100,
						Time:     now,
						Position: 0,
						Type:     types.MsgMultiSend,
						TxId:     1,
						Addresses: []storage.AddressWithType{
							{
								Type: types.MsgAddressTypeOutput,
								Address: storage.Address{
									Address:    "address1",
									Height:     100,
									LastHeight: 100,
								},
							}, {
								Type: types.MsgAddressTypeOutput,
								Address: storage.Address{
									Address:    "address1",
									Height:     100,
									LastHeight: 100,
								},
							},
						},
					},
				},
				addrToId: map[string]uint64{
					"address1": 1,
				},
			},
			wantNamespaceMessageCount: 0,
			wantMsgAddress:            1,
			wantErr:                   false,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	module := Module{
		validatorsByConsAddress: make(map[string]uint64),
	}

	for _, tt := range tests {
		tx := mock.NewMockTransaction(ctrl)

		tx.EXPECT().
			SaveNamespaceMessage(gomock.Any(), gomock.Any()).
			Times(1).
			DoAndReturn(func(_ context.Context, nsMsg ...storage.NamespaceMessage) error {
				require.Equal(t, tt.wantNamespaceMessageCount, len(nsMsg))
				return nil
			})

		tx.EXPECT().
			SaveMessages(gomock.Any(), gomock.Any()).
			Times(1).
			DoAndReturn(func(_ context.Context, msgs ...*storage.Message) error {
				require.Equal(t, len(tt.args.messages), len(msgs))
				return nil
			})

		tx.EXPECT().
			SaveMsgAddresses(gomock.Any(), gomock.Any()).
			Times(1).
			DoAndReturn(func(_ context.Context, msgAddr ...storage.MsgAddress) error {
				require.Equal(t, tt.wantMsgAddress, len(msgAddr))
				return nil
			})

		tx.EXPECT().
			SaveBlobLogs(gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		tx.EXPECT().
			SaveVestingAccounts(gomock.Any(), gomock.Any()).
			MaxTimes(1).
			Return(nil)

		tx.EXPECT().
			SaveVestingPeriods(gomock.Any(), gomock.Any()).
			MaxTimes(1).
			Return(nil)

		tx.EXPECT().
			SaveGrants(gomock.Any(), gomock.Any()).
			MaxTimes(1).
			MinTimes(1).
			Return(nil)

		t.Run(tt.name, func(t *testing.T) {
			err := module.saveMessages(t.Context(), tx, tt.args.messages, tt.args.addrToId)
			require.Equal(t, tt.wantErr, err != nil)
		})
	}
}
