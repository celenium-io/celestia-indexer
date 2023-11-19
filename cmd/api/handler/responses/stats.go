package responses

import (
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

type GasPriceCandle struct {
	High          string    `example:"0.17632"                   format:"string"    json:"high"               swaggertype:"string"`
	Low           string    `example:"0.11882"                   format:"string"    json:"low"                swaggertype:"string"`
	TotalGasLimit string    `example:"1213134"                   format:"string"    json:"total_gas_limit"    swaggertype:"string"`
	TotalGasUsed  string    `example:"0.45282"                   format:"string"    json:"total_gas_used"     swaggertype:"string"`
	Fee           int64     `example:"1283518"                   format:"integer"   json:"fee"                swaggertype:"number"`
	GasEfficiency string    `example:"0.45282"                   format:"string"    json:"avg_gas_efficiency" swaggertype:"string"`
	AvgGasPrice   string    `example:"0.45282"                   format:"string"    json:"avg_gas_price"      swaggertype:"string"`
	Time          time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"               swaggertype:"string"`
}

func NewGasPriceCandle(item storage.GasCandle) GasPriceCandle {
	return GasPriceCandle{
		Time:          item.Time,
		High:          formatFoat64(item.High),
		Low:           formatFoat64(item.Low),
		Fee:           item.Fee,
		TotalGasLimit: formatFoat64(item.Volume),
		TotalGasUsed:  formatFoat64(float64(item.GasUsed)),
		GasEfficiency: formatFoat64(float64(item.GasUsed) / item.Volume),
		AvgGasPrice:   formatFoat64(float64(item.Fee) / item.Volume),
	}
}

func formatFoat64(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

type NamespaceUsage struct {
	Name string `example:"00112233" format:"string"  json:"name" swaggertype:"string"`
	Size int64  `example:"1283518"  format:"integer" json:"size" swaggertype:"number"`
}

func NewNamespaceUsage(ns storage.Namespace) NamespaceUsage {
	return NamespaceUsage{
		Name: ns.String(),
		Size: ns.Size,
	}
}
