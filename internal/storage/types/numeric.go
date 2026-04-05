// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

import (
	"database/sql/driver"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

// Numeric is a wrapper around decimal.Decimal that implements pgx-native
// encoding/decoding via pgtype.Numeric, avoiding the string round-trip.
type Numeric struct {
	decimal.Decimal
}

var (
	_ pgtype.NumericValuer  = Numeric{}
	_ pgtype.NumericScanner = (*Numeric)(nil)
)

func NewNumeric(d decimal.Decimal) Numeric {
	return Numeric{d}
}

func NumericFromInt64(v int64) Numeric {
	return Numeric{decimal.NewFromInt(v)}
}

func NumericFromString(s string) (Numeric, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return Numeric{}, err
	}
	return Numeric{d}, nil
}

func MustNumericFromString(s string) Numeric {
	return Numeric{decimal.RequireFromString(s)}
}

func NumericFromFloat64(v float64) Numeric {
	return Numeric{decimal.NewFromFloat(v)}
}

// Value implements driver.Valuer, returning pgtype.Numeric for pgx.
func (n Numeric) Value() (driver.Value, error) {
	pn := pgtype.Numeric{Int: n.Coefficient(), Exp: n.Exponent(), Valid: true}
	return pn.Value()
}

// Scan implements sql.Scanner.
func (n *Numeric) Scan(src any) error {
	if src == nil {
		n.Decimal = decimal.Decimal{}
		return nil
	}
	// pgtype.Numeric.Scan handles only string and nil.
	// database/sql may deliver []byte, float64 or int64 depending on the column type.
	switch v := src.(type) {
	case []byte:
		src = string(v)
	case float64:
		n.Decimal = decimal.NewFromFloat(v)
		return nil
	case int64:
		n.Decimal = decimal.NewFromInt(v)
		return nil
	}
	var pn pgtype.Numeric
	if err := pn.Scan(src); err != nil {
		return err
	}
	if !pn.Valid || pn.Int == nil {
		n.Decimal = decimal.Decimal{}
		return nil
	}
	n.Decimal = decimal.NewFromBigInt(new(big.Int).Set(pn.Int), pn.Exp)
	return nil
}

// NumericValue implements pgx NumericValuer for native pgx COPY encoding.
func (n Numeric) NumericValue() (pgtype.Numeric, error) {
	return pgtype.Numeric{Int: n.Coefficient(), Exp: n.Exponent(), Valid: true}, nil
}

// ScanNumeric implements pgx NumericScanner for native pgx decoding.
func (n *Numeric) ScanNumeric(v pgtype.Numeric) error {
	if !v.Valid || v.Int == nil {
		n.Decimal = decimal.Decimal{}
		return nil
	}
	n.Decimal = decimal.NewFromBigInt(new(big.Int).Set(v.Int), v.Exp)
	return nil
}

// MarshalJSON implements json.Marshaler.
func (n Numeric) MarshalJSON() ([]byte, error) {
	return n.Decimal.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
func (n *Numeric) UnmarshalJSON(data []byte) error {
	return n.Decimal.UnmarshalJSON(data)
}

// Arithmetic methods that accept Numeric and return Numeric.

func (n Numeric) Add(d Numeric) Numeric {
	return Numeric{n.Decimal.Add(d.Decimal)}
}

func (n Numeric) Sub(d Numeric) Numeric {
	return Numeric{n.Decimal.Sub(d.Decimal)}
}

func (n Numeric) Mul(d Numeric) Numeric {
	return Numeric{n.Decimal.Mul(d.Decimal)}
}

func (n Numeric) Div(d Numeric) Numeric {
	return Numeric{n.Decimal.Div(d.Decimal)}
}

func (n Numeric) Neg() Numeric {
	return Numeric{n.Decimal.Neg()}
}

func (n Numeric) Copy() Numeric {
	return Numeric{n.Decimal.Copy()}
}

func (n Numeric) Floor() Numeric {
	return Numeric{n.Decimal.Floor()}
}

func (n Numeric) Pow(d Numeric) Numeric {
	return Numeric{n.Decimal.Pow(d.Decimal)}
}

// Comparison methods that accept Numeric.

func (n Numeric) GreaterThan(d Numeric) bool {
	return n.Decimal.GreaterThan(d.Decimal)
}

func (n Numeric) GreaterThanOrEqual(d Numeric) bool {
	return n.Decimal.GreaterThanOrEqual(d.Decimal)
}

func (n Numeric) LessThan(d Numeric) bool {
	return n.Decimal.LessThan(d.Decimal)
}

func (n Numeric) LessThanOrEqual(d Numeric) bool {
	return n.Decimal.LessThanOrEqual(d.Decimal)
}

func (n Numeric) Equal(d Numeric) bool {
	return n.Decimal.Equal(d.Decimal)
}
