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

func NewNumeric(d decimal.Decimal) Numeric {
	return Numeric{d}
}

func NumericFromInt64(v int64) Numeric {
	return Numeric{decimal.NewFromInt(v)}
}

// Value implements driver.Valuer, returning pgtype.Numeric for pgx.
func (n Numeric) Value() (driver.Value, error) {
	pn := pgtype.Numeric{Int: n.Coefficient(), Exp: n.Exponent(), Valid: true}
	return pn.Value()
}

// Scan implements sql.Scanner.
func (n *Numeric) Scan(src any) error {
	// pgtype.Numeric.Scan does not handle []byte; convert to string first.
	if b, ok := src.([]byte); ok {
		src = string(b)
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

// Arithmetic methods that return Numeric instead of decimal.Decimal.

func (n Numeric) Add(d decimal.Decimal) Numeric {
	return Numeric{n.Decimal.Add(d)}
}

func (n Numeric) Sub(d decimal.Decimal) Numeric {
	return Numeric{n.Decimal.Sub(d)}
}

func (n Numeric) Mul(d decimal.Decimal) Numeric {
	return Numeric{n.Decimal.Mul(d)}
}

func (n Numeric) Div(d decimal.Decimal) Numeric {
	return Numeric{n.Decimal.Div(d)}
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

func (n Numeric) Pow(d decimal.Decimal) Numeric {
	return Numeric{n.Decimal.Pow(d)}
}
