// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestBlock_SearchType(t *testing.T) {
	block := Block{}
	searchType := block.SearchType()
	assert.EqualValues(t, "block", searchType)
}

func TestNewBlock(t *testing.T) {
	type args struct {
		block     storage.Block
		withStats bool
	}
	tests := []struct {
		name string
		args args
		want Block
	}{
		{
			name: "without stats",
			args: args{
				block: storage.Block{
					Id:           1000,
					Height:       1000,
					Time:         time.Time{},
					VersionBlock: 10,
					VersionApp:   11,
					MessageTypes: storageTypes.MsgTypeBits{
						Bits: 1,
					},
					Hash:               []byte{0x01},
					ParentHash:         []byte{0x02},
					LastCommitHash:     []byte{0x03},
					DataHash:           []byte{0x04},
					ValidatorsHash:     []byte{0x05},
					NextValidatorsHash: []byte{0x06},
					ConsensusHash:      []byte{0x07},
					AppHash:            []byte{0x08},
					LastResultsHash:    []byte{0x09},
					EvidenceHash:       []byte{0x10},
					ProposerAddress:    []byte{0x11},
					ChainId:            "dipdup",
				},
				withStats: false,
			},
			want: Block{
				Id:                 1000,
				Height:             1000,
				Time:               time.Time{},
				VersionBlock:       "10",
				VersionApp:         "11",
				Hash:               []byte{0x01},
				ParentHash:         []byte{0x02},
				LastCommitHash:     []byte{0x03},
				DataHash:           []byte{0x04},
				ValidatorsHash:     []byte{0x05},
				NextValidatorsHash: []byte{0x06},
				ConsensusHash:      []byte{0x07},
				AppHash:            []byte{0x08},
				LastResultsHash:    []byte{0x09},
				EvidenceHash:       []byte{0x10},
				ProposerAddress:    []byte{0x11},
				MessageTypes:       []storageTypes.MsgType{storageTypes.MsgUnknown},
				Stats:              nil,
			},
		},

		{
			name: "with stats",
			args: args{
				block: storage.Block{
					Id:           1000,
					Height:       1000,
					Time:         time.Time{},
					VersionBlock: 10,
					VersionApp:   11,
					MessageTypes: storageTypes.MsgTypeBits{
						Bits: 1,
					},
					Hash:               []byte{0x01},
					ParentHash:         []byte{0x02},
					LastCommitHash:     []byte{0x03},
					DataHash:           []byte{0x04},
					ValidatorsHash:     []byte{0x05},
					NextValidatorsHash: []byte{0x06},
					ConsensusHash:      []byte{0x07},
					AppHash:            []byte{0x08},
					LastResultsHash:    []byte{0x09},
					EvidenceHash:       []byte{0x10},
					ProposerAddress:    []byte{0x11},
					ChainId:            "dipdup",
					Stats: storage.BlockStats{
						Id:            1000,
						Height:        1000,
						Time:          time.Time{},
						TxCount:       6,
						EventsCount:   10,
						BlobsSize:     1234,
						BlockTime:     11000,
						SupplyChange:  decimal.NewFromInt(123),
						InflationRate: decimal.NewFromFloat(0.08),
						Fee:           decimal.NewFromInt(125),
						MessagesCounts: map[storageTypes.MsgType]int64{
							storageTypes.MsgSend:        1,
							storageTypes.MsgPayForBlobs: 2,
							storageTypes.MsgDelegate:    3,
						},
					},
				},
				withStats: true,
			},
			want: Block{
				Id:                 1000,
				Height:             1000,
				Time:               time.Time{},
				VersionBlock:       "10",
				VersionApp:         "11",
				Hash:               []byte{0x01},
				ParentHash:         []byte{0x02},
				LastCommitHash:     []byte{0x03},
				DataHash:           []byte{0x04},
				ValidatorsHash:     []byte{0x05},
				NextValidatorsHash: []byte{0x06},
				ConsensusHash:      []byte{0x07},
				AppHash:            []byte{0x08},
				LastResultsHash:    []byte{0x09},
				EvidenceHash:       []byte{0x10},
				ProposerAddress:    []byte{0x11},
				MessageTypes:       []storageTypes.MsgType{storageTypes.MsgUnknown},
				Stats: &BlockStats{
					TxCount:       6,
					EventsCount:   10,
					BlobsSize:     1234,
					Fee:           "125",
					SupplyChange:  "123",
					InflationRate: "0.08",
					BlockTime:     11000,
					MessagesCounts: map[storageTypes.MsgType]int64{
						storageTypes.MsgSend:        1,
						storageTypes.MsgPayForBlobs: 2,
						storageTypes.MsgDelegate:    3,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBlock(tt.args.block, tt.args.withStats); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBlockStats(t *testing.T) {
	type args struct {
		stats storage.BlockStats
	}
	tests := []struct {
		name string
		args args
		want *BlockStats
	}{
		{
			name: "with msgs counts",
			args: args{
				stats: storage.BlockStats{
					Id:            100,
					Height:        100,
					Time:          time.Time{},
					TxCount:       6,
					EventsCount:   10,
					BlobsSize:     1234,
					BlockTime:     11000,
					SupplyChange:  decimal.NewFromInt(123),
					InflationRate: decimal.NewFromFloat(0.08),
					Fee:           decimal.NewFromInt(125),
					MessagesCounts: map[storageTypes.MsgType]int64{
						storageTypes.MsgSend:        1,
						storageTypes.MsgPayForBlobs: 2,
						storageTypes.MsgDelegate:    3,
					},
				},
			},
			want: &BlockStats{
				TxCount:       6,
				EventsCount:   10,
				BlobsSize:     1234,
				Fee:           "125",
				SupplyChange:  "123",
				InflationRate: "0.08",
				BlockTime:     11000,
				MessagesCounts: map[storageTypes.MsgType]int64{
					storageTypes.MsgSend:        1,
					storageTypes.MsgPayForBlobs: 2,
					storageTypes.MsgDelegate:    3,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBlockStats(tt.args.stats); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBlockStats() = %v, want %v", got, tt.want)
			}
		})
	}
}
