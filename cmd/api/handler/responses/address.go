// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package responses

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	celestials "github.com/celenium-io/celestial-module/pkg/storage"
	"github.com/celestiaorg/celestia-app/v11/app/frozen"
	"github.com/celestiaorg/celestia-app/v11/pkg/appconsts"
)

// Address model info
//
//	@Description	Celestia address information
type Address struct {
	Id           uint64         `example:"321"                                             format:"int64"        json:"id"             swaggertype:"integer"`
	Height       pkgTypes.Level `example:"100"                                             format:"int64"        json:"first_height"   swaggertype:"integer"`
	LastHeight   pkgTypes.Level `example:"100"                                             format:"int64"        json:"last_height"    swaggertype:"integer"`
	Hash         string         `example:"celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60" json:"hash"           swaggertype:"string"`
	Name         string         `example:"bonded_tokens_pool"                              json:"name,omitempty" swaggertype:"string"`
	IsForwarding bool           `example:"true"                                            json:"is_forwarding"  swaggertype:"boolean"`
	IsFrozen     bool           `example:"true"                                            json:"is_frozen"      swaggertype:"boolean"`
	IsWhitehat   bool           `example:"false"                                           json:"is_whitehat"    swaggertype:"boolean"`
	Balance      Balance        `json:"balance"`

	Celestials *Celestial `json:"celestials,omitempty"`
}

var frozenRules = frozen.NewRules(nil, nil)

// whitehats mirrors the unexported per-chain map in app/frozen.
var whitehats = map[string]string{
	appconsts.MainnetChainID: frozen.MainnetWhitehat,
	appconsts.MochaChainID:   frozen.MochaWhitehat,
	appconsts.CortoChainID:   frozen.CortoWhitehat,
}

// NewAddress marks frozen and whitehat accounts of chainId; an empty chainId means the rules are not active yet.
func NewAddress(addr storage.Address, chainId string) Address {
	address := Address{
		Id:           addr.Id,
		Height:       addr.Height,
		LastHeight:   addr.LastHeight,
		Hash:         addr.Address,
		Name:         addr.Name,
		IsForwarding: addr.IsForwarding,
		IsFrozen:     frozenRules.Has(chainId, addr.Hash),
		IsWhitehat:   chainId != "" && whitehats[chainId] == addr.Address,
	}
	if addr.DefaultBalance != nil {
		address.Balance = NewBalance(*addr.DefaultBalance)
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

func NewBalance(balance storage.Balance) Balance {
	return Balance{
		Currency:  balance.Currency,
		Spendable: balance.Spendable.String(),
		Delegated: balance.Delegated.String(),
		Unbonding: balance.Unbonding.String(),
	}
}

// Celestial ID
//
//	@Description	Linked celestial id
type Celestial struct {
	Name      string `example:"name"                                                                json:"name"              swaggertype:"string"`
	ImageUrl  string `example:"https://ipfs.io/ipfs/QmUi269vE25fagqhyMCCTNSoiW6x4LHCwwQb3keSrEXAmC" json:"image_url"         swaggertype:"string"`
	Status    string `example:"VERIFIED"                                                            json:"status,omitempty"  swaggertype:"string"`
	IsPrimary bool   `example:"true"                                                                json:"primary,omitempty" swaggertype:"boolean"`
}

func NewCelestial(c *celestials.Celestial) *Celestial {
	return &Celestial{
		ImageUrl:  c.ImageUrl,
		Name:      c.Id,
		IsPrimary: c.Status == celestials.StatusPRIMARY,
		Status:    c.Status.String(),
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
