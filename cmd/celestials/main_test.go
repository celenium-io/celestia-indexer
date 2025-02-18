package main

import (
	"context"
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAddressHandler(t *testing.T) {
	tests := []struct {
		address string
		wantErr bool
	}{
		{
			address: "celestia1a9qy9fuxyteksjhtv2mxvs7z29nt5s3cf6d6ln",
			wantErr: false,
		}, {
			address: "celestia1kw6mw70wafdxgp2n8s4lscx04du8ka6dvul6jy",
			wantErr: false,
		}, {
			address: "osmo13ge29x4e2s63a8ytz2px8gurtyznmue4a69n5275692v3qn3ks8q7cwck7",
			wantErr: true,
		}, {
			address: "0x79FF9170499b0691c3878D6f95519dB05c53C9a1",
			wantErr: true,
		}, {
			address: "celestia1kw6mw70wafdxgp2n8s4lscx04du8ka6dvul6jy",
			wantErr: true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	address := mock.NewMockIAddress(ctrl)

	for _, tt := range tests {
		t.Run(tt.address, func(t *testing.T) {
			if tt.wantErr {
				address.EXPECT().
					IdByHash(gomock.Any(), gomock.Any()).
					Return([]uint64{}, nil).
					AnyTimes()
			} else {
				address.EXPECT().
					IdByHash(gomock.Any(), gomock.Any()).
					Return([]uint64{1}, nil).
					Times(1)
			}

			id, err := addressHandler(context.Background(), address, tt.address)
			require.Equal(t, err != nil, tt.wantErr)
			if err == nil {
				require.Equal(t, uint64(1), id)
			}
		})
	}
}
