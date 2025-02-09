// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestNewDistributionItem(t *testing.T) {
	tests := []struct {
		name       string
		item       storage.DistributionItem
		tf         string
		wantResult DistributionItem
	}{
		{
			name: "Sunday",
			item: storage.DistributionItem{
				Name:  0,
				Value: "10",
			},
			tf: "day",
			wantResult: DistributionItem{
				Name:  "Sunday",
				Value: "10",
			},
		}, {
			name: "Monday",
			item: storage.DistributionItem{
				Name:  1,
				Value: "10",
			},
			tf: "day",
			wantResult: DistributionItem{
				Name:  "Monday",
				Value: "10",
			},
		}, {
			name: "Tuesday",
			item: storage.DistributionItem{
				Name:  2,
				Value: "10",
			},
			tf: "day",
			wantResult: DistributionItem{
				Name:  "Tuesday",
				Value: "10",
			},
		}, {
			name: "Wednesday",
			item: storage.DistributionItem{
				Name:  3,
				Value: "10",
			},
			tf: "day",
			wantResult: DistributionItem{
				Name:  "Wednesday",
				Value: "10",
			},
		}, {
			name: "Thursday",
			item: storage.DistributionItem{
				Name:  4,
				Value: "10",
			},
			tf: "day",
			wantResult: DistributionItem{
				Name:  "Thursday",
				Value: "10",
			},
		}, {
			name: "Friday",
			item: storage.DistributionItem{
				Name:  5,
				Value: "10",
			},
			tf: "day",
			wantResult: DistributionItem{
				Name:  "Friday",
				Value: "10",
			},
		}, {
			name: "Saturday",
			item: storage.DistributionItem{
				Name:  6,
				Value: "10",
			},
			tf: "day",
			wantResult: DistributionItem{
				Name:  "Saturday",
				Value: "10",
			},
		}, {
			name: "Sunday",
			item: storage.DistributionItem{
				Name:  7,
				Value: "10",
			},
			tf: "day",
			wantResult: DistributionItem{
				Name:  "Sunday",
				Value: "10",
			},
		}, {
			name: "10 hour",
			item: storage.DistributionItem{
				Name:  10,
				Value: "10",
			},
			tf: "hour",
			wantResult: DistributionItem{
				Name:  "10",
				Value: "10",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := NewDistributionItem(tt.item, tt.tf)
			require.Equal(t, tt.wantResult, gotResult)
		})
	}
}
