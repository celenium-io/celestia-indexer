package math

import "github.com/shopspring/decimal"

var powerDivider = decimal.NewFromInt(1_000_000)

func VotingPower(stake decimal.Decimal) decimal.Decimal {
	return stake.Div(powerDivider).Floor()
}
