// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"
	"strconv"
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

type Change24hBlockStats struct {
	TxCount      float64 `example:"0.1234" format:"float" json:"tx_count_24h"       swaggertype:"number"`
	Fee          float64 `example:"0.1234" format:"float" json:"fee_24h"            swaggertype:"number"`
	BytesInBlock float64 `example:"0.1234" format:"float" json:"bytes_in_block_24h" swaggertype:"number"`
	BlobsSize    float64 `example:"0.1234" format:"float" json:"blobs_size_24h"     swaggertype:"number"`
}

func NewChange24hBlockStats(response storage.Change24hBlockStats) Change24hBlockStats {
	return Change24hBlockStats{
		TxCount:      response.TxCount,
		Fee:          response.Fee,
		BytesInBlock: response.BytesInBlock,
		BlobsSize:    response.BlobsSize,
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

type DistributionItem struct {
	Name  string `example:"12"      format:"string" json:"name"  swaggertype:"string"`
	Value string `example:"0.17632" format:"string" json:"value" swaggertype:"string"`
}

func NewDistributionItem(item storage.DistributionItem, tf string) (result DistributionItem) {
	result.Value = item.Value

	switch tf {
	case "day":
		result.Name = time.Weekday(item.Name).String()
	case "hour":
		result.Name = strconv.FormatInt(int64(item.Name), 10)
	default:
		result.Name = strconv.FormatInt(int64(item.Name), 10)
	}

	return
}

type TimeValueItem struct {
	Time  time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"  swaggertype:"string"`
	Value string    `example:"0.17632"                   format:"string"    json:"value" swaggertype:"string"`
}

type SquareSizeResponse map[int][]TimeValueItem

func NewSquareSizeResponse(m map[int][]storage.SeriesItem) SquareSizeResponse {
	response := make(SquareSizeResponse)
	for key, value := range m {
		response[key] = make([]TimeValueItem, len(value))
		for i := range value {
			response[key][i].Time = value[i].Time
			response[key][i].Value = value[i].Value
		}
	}
	return response
}

type RollupStats24h struct {
	Id         int64   `example:"321"                             format:"integer" json:"id,omitempty"   swaggertype:"integer"`
	Name       string  `example:"Rollup name"                     format:"string"  json:"name,omitempty" swaggertype:"string"`
	Logo       string  `example:"https://some_link.com/image.png" format:"string"  json:"logo,omitempty" swaggertype:"string"`
	Size       int64   `example:"123"                             format:"integer" json:"size"           swaggertype:"integer"`
	Fee        float64 `example:"123"                             format:"number"  json:"fee"            swaggertype:"integer"`
	BlobsCount int64   `example:"123"                             format:"integer" json:"blobs_count"    swaggertype:"integer"`
}

func NewRollupStats24h(stats storage.RollupStats24h) RollupStats24h {
	return RollupStats24h{
		Id:         stats.RollupId,
		Name:       stats.Name,
		Logo:       stats.Logo,
		Size:       stats.Size,
		Fee:        stats.Fee,
		BlobsCount: stats.BlobsCount,
	}
}

type CountItem struct {
	Name  string `example:"test"  format:"string" json:"name"  swaggertype:"string"`
	Value int64  `example:"17632" format:"string" json:"value" swaggertype:"string"`
}

func NewCountItem(item storage.CountItem) CountItem {
	return CountItem{
		Name:  item.Name,
		Value: item.Value,
	}
}

type RollupAllSeriesItem struct {
	Time       time.Time `example:"2023-07-04T03:10:57+00:00"       format:"date-time" json:"time"           swaggertype:"string"`
	Name       string    `example:"Rollup name"                     format:"string"    json:"name,omitempty" swaggertype:"string"`
	Logo       string    `example:"https://some_link.com/image.png" format:"string"    json:"logo,omitempty" swaggertype:"string"`
	Size       int64     `example:"123"                             format:"integer"   json:"size"           swaggertype:"integer"`
	Fee        string    `example:"123"                             format:"string"    json:"fee"            swaggertype:"string"`
	BlobsCount int64     `example:"123"                             format:"integer"   json:"blobs_count"    swaggertype:"integer"`
}

func NewRollupAllSeriesItem(stats storage.RollupHistogramItem) RollupAllSeriesItem {
	return RollupAllSeriesItem{
		Time:       stats.Time,
		Name:       stats.Name,
		Logo:       stats.Logo,
		Size:       stats.Size,
		Fee:        stats.Fee,
		BlobsCount: stats.BlobsCount,
	}
}
