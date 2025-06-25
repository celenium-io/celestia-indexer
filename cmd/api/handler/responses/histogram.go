// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package responses

import (
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
)

type HistogramItem struct {
	Time  time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"  swaggertype:"string"`
	Value string    `example:"2223424"                   format:"string"    json:"value" swaggertype:"string"`
}

func NewHistogramItem(item storage.HistogramItem) HistogramItem {
	return HistogramItem{
		Time:  item.Time,
		Value: item.Value,
	}
}
