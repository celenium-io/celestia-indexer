// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	celestials "github.com/celenium-io/celestial-module/pkg/storage"
)

// Address model info
//
//	@Description	Celestia address information
type Address struct {
	Id         uint64         `example:"321"                                             json:"id"           swaggertype:"integer"`
	Height     pkgTypes.Level `example:"100"                                             json:"first_height" swaggertype:"integer"`
	LastHeight pkgTypes.Level `example:"100"                                             json:"last_height"  swaggertype:"integer"`
	Hash       string         `example:"celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60" json:"hash"         swaggertype:"string"`
	Balance    Balance        `json:"balance"`

	Celestials *Celestial `json:"celestials,omitempty"`
}

func NewAddress(addr storage.Address) Address {
	address := Address{
		Id:         addr.Id,
		Height:     addr.Height,
		LastHeight: addr.LastHeight,
		Hash:       addr.Address,
		Balance: Balance{
			Currency:  addr.Balance.Currency,
			Spendable: addr.Balance.Spendable.String(),
			Delegated: addr.Balance.Delegated.String(),
			Unbonding: addr.Balance.Unbonding.String(),
		},
	}
	address.AddCelestails(addr.Celestials)
	return address
}

func (address *Address) AddCelestails(celestials *celestials.Celestial) {
	if celestials != nil {
		address.Celestials = NewCelestial(celestials)
	}
}

// Balance info
//
//	@Description	Balance of address information
type Balance struct {
	Currency  string `example:"utia"        json:"currency"  swaggertype:"string"`
	Spendable string `example:"10000000000" json:"spendable" swaggertype:"string"`
	Delegated string `example:"10000000000" json:"delegated" swaggertype:"string"`
	Unbonding string `example:"10000000000" json:"unbonding" swaggertype:"string"`
}

// Celestial ID
//
//	@Description	Linked celestial id
type Celestial struct {
	Name      string `example:"name"                                                                json:"name"              swaggertype:"string"`
	ImageUrl  string `example:"https://ipfs.io/ipfs/QmUi269vE25fagqhyMCCTNSoiW6x4LHCwwQb3keSrEXAmC" json:"image_url"         swaggertype:"string"`
	IsPrimary bool   `example:"true"                                                                json:"primary,omitempty" swaggertype:"boolean"`
}

func NewCelestial(c *celestials.Celestial) *Celestial {
	return &Celestial{
		ImageUrl:  c.ImageUrl,
		Name:      c.Id,
		IsPrimary: c.Status == celestials.StatusPRIMARY,
	}
}

type ShortAddress struct {
	Hash       string     `example:"celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60" json:"hash" swaggertype:"string"`
	Celestials *Celestial `json:"celestials,omitempty"`
}

func NewShortAddress(address *storage.Address) *ShortAddress {
	if address == nil {
		return nil
	}

	result := new(ShortAddress)
	result.Hash = address.Address
	if address.Celestials != nil {
		result.Celestials = NewCelestial(address.Celestials)
	}
	return result
}
