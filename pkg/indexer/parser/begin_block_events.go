package parser

import (
	"strings"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/shopspring/decimal"
)

type beginBlockEventsResult struct {
	SupplyChange  decimal.Decimal
	InflationRate decimal.Decimal
}

func (bbe *beginBlockEventsResult) Fill(events []storage.Event) error {
	for i := range events {
		switch events[i].Type {
		case types.EventTypeBurn:
			bbe.SupplyChange = bbe.SupplyChange.Sub(getDecimalFromMap(events[i].Data, "amount"))
		case types.EventTypeMint:
			bbe.InflationRate = getDecimalFromMap(events[i].Data, "inflation_rate")
			bbe.SupplyChange = bbe.SupplyChange.Add(getDecimalFromMap(events[i].Data, "amount"))
		}
	}

	return nil
}

func getDecimalFromMap(m map[string]any, key string) decimal.Decimal {
	val, ok := m[key]
	if !ok {
		return decimal.Zero
	}
	str, ok := val.(string)
	if !ok {
		return decimal.Zero
	}
	str = strings.TrimSuffix(str, "utia")
	dec, err := decimal.NewFromString(str)
	if err != nil {
		return decimal.Zero
	}
	return dec
}
