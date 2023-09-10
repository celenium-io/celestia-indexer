package types

import (
	"encoding/hex"
	"math/big"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// celestia prefixes
const (
	AddressPrefixCelestia = "celestia"
	AddressPrefixValoper  = "celestiavaloper"
)

type Address string

// Decode decodes address, returning the human-readable part and the data part excluding the checksum.
func (a Address) Decode() (string, []byte, error) {
	return bech32.DecodeAndConvert(a.String())
}

func (a Address) String() string {
	return string(a)
}

func (a Address) Decimal() (decimal.Decimal, error) {
	_, data, err := a.Decode()
	if err != nil {
		return decimal.Zero, err
	}

	if bi, ok := new(big.Int).SetString(hex.EncodeToString(data), 16); ok {
		return decimal.NewFromBigInt(bi, 0), nil
	}
	return decimal.Zero, errors.Errorf("invalid decoded address: %x %s", data, a)
}
