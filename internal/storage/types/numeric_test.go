// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

import (
	"math/big"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestNewNumeric(t *testing.T) {
	d := decimal.RequireFromString("123.456")
	n := NewNumeric(d)
	require.True(t, n.Equal(d))
}

func TestNumericFromInt64(t *testing.T) {
	n := NumericFromInt64(42)
	require.True(t, n.Equal(decimal.NewFromInt(42)))
}

func TestNumeric_Decimal(t *testing.T) {
	d := decimal.RequireFromString("-99.99")
	n := NewNumeric(d)
	require.True(t, n.Equal(d))
}

func TestNumeric_Zero(t *testing.T) {
	n := NewNumeric(decimal.Zero)
	require.True(t, n.IsZero())
}

func TestNumeric_Value(t *testing.T) {
	tests := []struct {
		name string
		val  string
	}{
		{"zero", "0"},
		{"positive integer", "12345"},
		{"negative", "-100"},
		{"decimal", "123.456"},
		{"large", "99999999999999999999"},
		{"small fraction", "0.000001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := decimal.RequireFromString(tt.val)
			n := NewNumeric(d)

			v, err := n.Value()
			require.NoError(t, err)
			require.NotNil(t, v)

			// Scan back through pgtype.Numeric to verify round-trip
			var restored Numeric
			err = restored.Scan(v)
			require.NoError(t, err)
			require.True(t, restored.Equal(d), "expected %s, got %s", d, restored.String())
		})
	}
}

func TestNumeric_ScanString(t *testing.T) {
	var n Numeric
	err := n.Scan("123.456")
	require.NoError(t, err)
	require.True(t, n.Equal(decimal.RequireFromString("123.456")))
}

func TestNumeric_ScanBytes(t *testing.T) {
	var n Numeric
	err := n.Scan([]byte("17263"))
	require.NoError(t, err)
	require.True(t, n.Equal(decimal.NewFromInt(17263)))
}

func TestNumeric_ScanFloat64(t *testing.T) {
	var n Numeric
	err := n.Scan(float64(3.14))
	require.NoError(t, err)
	require.True(t, n.Equal(decimal.NewFromFloat(3.14)))
}

func TestNumeric_ScanInt64(t *testing.T) {
	var n Numeric
	err := n.Scan(int64(999))
	require.NoError(t, err)
	require.True(t, n.Equal(decimal.NewFromInt(999)))
}

func TestNumeric_ScanNil(t *testing.T) {
	n := NewNumeric(decimal.NewFromInt(42))
	err := n.Scan(nil)
	require.NoError(t, err)
	require.True(t, n.IsZero())
}

func TestNumeric_ScanInvalid(t *testing.T) {
	var n Numeric
	err := n.Scan("not_a_number")
	require.Error(t, err)
}

func TestNumeric_NumericValue(t *testing.T) {
	d := decimal.RequireFromString("456.789")
	n := NewNumeric(d)

	pn, err := n.NumericValue()
	require.NoError(t, err)
	require.True(t, pn.Valid)
	require.NotNil(t, pn.Int)

	// Reconstruct from pgtype.Numeric
	restored := decimal.NewFromBigInt(new(big.Int).Set(pn.Int), pn.Exp)
	require.True(t, restored.Equal(d))
}

func TestNumeric_ScanNumeric(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		pn := pgtype.Numeric{
			Int:   big.NewInt(123456),
			Exp:   -3,
			Valid: true,
		}
		var n Numeric
		err := n.ScanNumeric(pn)
		require.NoError(t, err)
		require.True(t, n.Equal(decimal.RequireFromString("123.456")))
	})

	t.Run("invalid", func(t *testing.T) {
		pn := pgtype.Numeric{Valid: false}
		var n Numeric
		err := n.ScanNumeric(pn)
		require.NoError(t, err)
		require.True(t, n.IsZero())
	})

	t.Run("nil int", func(t *testing.T) {
		pn := pgtype.Numeric{Valid: true, Int: nil}
		var n Numeric
		err := n.ScanNumeric(pn)
		require.NoError(t, err)
		require.True(t, n.IsZero())
	})
}

func TestNumeric_NumericValueScanRoundTrip(t *testing.T) {
	values := []string{
		"0", "1", "-1", "123.456", "-999.999",
		"100000000000", "0.000000001",
	}

	for _, v := range values {
		t.Run(v, func(t *testing.T) {
			original := NumericFromString(v)

			pn, err := original.NumericValue()
			require.NoError(t, err)

			var restored Numeric
			err = restored.ScanNumeric(pn)
			require.NoError(t, err)
			require.True(t, restored.Equal(original.Decimal),
				"expected %s, got %s", original.String(), restored.String())
		})
	}
}

func TestNumeric_MarshalJSON(t *testing.T) {
	n := NumericFromString("123.456")
	data, err := json.Marshal(n)
	require.NoError(t, err)
	require.Equal(t, `"123.456"`, string(data))
}

func TestNumeric_UnmarshalJSON(t *testing.T) {
	var n Numeric
	err := json.Unmarshal([]byte(`"789.012"`), &n)
	require.NoError(t, err)
	require.True(t, n.Equal(decimal.RequireFromString("789.012")))
}

func TestNumeric_JSONRoundTrip(t *testing.T) {
	type wrapper struct {
		Amount Numeric `json:"amount"`
	}

	original := wrapper{Amount: NumericFromString("-42.5")}
	data, err := json.Marshal(original)
	require.NoError(t, err)

	var restored wrapper
	err = json.Unmarshal(data, &restored)
	require.NoError(t, err)
	require.True(t, restored.Amount.Equal(original.Amount.Decimal))
}

func TestNumeric_UnmarshalJSON_Invalid(t *testing.T) {
	var n Numeric
	err := json.Unmarshal([]byte(`"not_a_number"`), &n)
	require.Error(t, err)
}
