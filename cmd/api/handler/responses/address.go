// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
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
}

func NewAddress(addr storage.Address) Address {
	return Address{
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
