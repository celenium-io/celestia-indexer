package responses

import (
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
