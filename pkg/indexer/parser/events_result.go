package parser

import (
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/decode"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type eventsResult struct {
	SupplyChange  decimal.Decimal
	InflationRate decimal.Decimal

	Addresses []storage.Address
}

func (er *eventsResult) Fill(events []storage.Event) error {
	er.Addresses = make([]storage.Address, 0)

	for i := range events {
		switch events[i].Type {
		case types.EventTypeBurn:
			amount := decode.Amount(events[i].Data)
			er.SupplyChange = er.SupplyChange.Sub(amount)
		case types.EventTypeMint:
			er.InflationRate = decode.DecimalFromMap(events[i].Data, "inflation_rate")
			amount := decode.Amount(events[i].Data)
			er.SupplyChange = er.SupplyChange.Add(amount)
		case types.EventTypeCoinReceived:
			address, err := parseCoinReceived(events[i].Data, events[i].Height)
			if err != nil {
				return errors.Wrap(err, "parse coin received")
			}
			if address != nil {
				er.Addresses = append(er.Addresses, *address)
			}
		case types.EventTypeCoinSpent:
			address, err := parseCoinSpent(events[i].Data, events[i].Height)
			if err != nil {
				return errors.Wrap(err, "parse coin spent")
			}
			if address != nil {
				er.Addresses = append(er.Addresses, *address)
			}
		}
	}

	return nil
}
