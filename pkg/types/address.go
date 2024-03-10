// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

import (
	"crypto/ed25519"
	"encoding/hex"
	"math/big"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// celestia prefixes
const (
	AddressPrefixCelestia = "celestia"
	AddressPrefixValoper  = "celestiavaloper"
	AddressPrefixValCons  = "celestiavalcons"
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

func NewAddressFromBytes(data []byte) (Address, error) {
	s, err := bech32.ConvertAndEncode(AddressPrefixCelestia, data)
	if err != nil {
		return "", nil
	}
	return Address(s), nil
}

func NewConsAddressFromBytes(data []byte) (Address, error) {
	s, err := bech32.ConvertAndEncode(AddressPrefixValCons, data)
	if err != nil {
		return "", nil
	}
	return Address(s), nil
}

func NewValoperAddressFromBytes(data []byte) (Address, error) {
	s, err := bech32.ConvertAndEncode(AddressPrefixValoper, data)
	if err != nil {
		return "", nil
	}

	return Address(s), nil
}

func GetConsAddressBytesFromPubKey(data []byte) []byte {
	pk := cryptotypes.PubKey{
		Key: ed25519.PublicKey(data),
	}
	return pk.Address().Bytes()
}
