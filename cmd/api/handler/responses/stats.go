// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
)

type TPS struct {
	High              float64 `example:"1.023" format:"float" json:"high"                 swaggertype:"number"`
	Low               float64 `example:"0.123" format:"float" json:"low"                  swaggertype:"number"`
	Current           float64 `example:"0.567" format:"float" json:"current"              swaggertype:"number"`
	ChangeLastHourPct float64 `example:"0.275" format:"float" json:"change_last_hour_pct" swaggertype:"number"`
}

func NewTPS(tps storage.TPS) TPS {
	return TPS{
		High:              tps.High,
		Low:               tps.Low,
		Current:           tps.Current,
		ChangeLastHourPct: tps.ChangeLastHourPct,
	}
}

type TxCountHistogramItem struct {
	Time  time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"  swaggertype:"string"`
	Count int64     `example:"2223424"                   format:"integer"   json:"count" swaggertype:"integer"`
	TPS   float64   `example:"0.13521"                   format:"float"     json:"tps"   swaggertype:"number"`
}

func NewTxCountHistogramItem(item storage.TxCountForLast24hItem) TxCountHistogramItem {
	return TxCountHistogramItem{
		Time:  item.Time,
		Count: item.TxCount,
		TPS:   item.TPS,
	}
}

type NamespaceUsage struct {
	Name        string `example:"00112233"                                                 format:"string"  json:"name"                   swaggertype:"string"`
	Version     *byte  `examle:"1"                                                         format:"byte"    json:"version,omitempty"      swaggertype:"integer"`
	NamespaceID string `example:"4723ce10b187716adfc55ff7e6d9179c226e6b5440b02577cca49d02" format:"binary"  json:"namespace_id,omitempty" swaggertype:"string"`
	Size        int64  `example:"1283518"                                                  format:"integer" json:"size"                   swaggertype:"number"`
}

func NewNamespaceUsage(ns storage.Namespace) NamespaceUsage {
	return NamespaceUsage{
		Name:        decodeName(ns.NamespaceID),
		Size:        ns.Size,
		Version:     &ns.Version,
		NamespaceID: hex.EncodeToString(ns.NamespaceID),
	}
}

type SeriesItem struct {
	Time  time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"          swaggertype:"string"`
	Value string    `example:"0.17632"                   format:"string"    json:"value"         swaggertype:"string"`
	Max   string    `example:"0.17632"                   format:"string"    json:"max,omitempty" swaggertype:"string"`
	Min   string    `example:"0.17632"                   format:"string"    json:"min,omitempty" swaggertype:"string"`
}

func NewSeriesItem(item storage.SeriesItem) SeriesItem {
	return SeriesItem{
		Time:  item.Time,
		Value: item.Value,
		Max:   item.Max,
		Min:   item.Min,
	}
}

type Price struct {
	Time  time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"  swaggertype:"string"`
	Open  string    `example:"0.17632"                   format:"string"    json:"open"  swaggertype:"string"`
	High  string    `example:"0.17632"                   format:"string"    json:"high"  swaggertype:"string"`
	Low   string    `example:"0.17632"                   format:"string"    json:"low"   swaggertype:"string"`
	Close string    `example:"0.17632"                   format:"string"    json:"close" swaggertype:"string"`
}

func NewPrice(price storage.Price) Price {
	return Price{
		Time:  price.Time,
		Open:  price.Open.String(),
		High:  price.High.String(),
		Low:   price.Low.String(),
		Close: price.Close.String(),
	}
}
