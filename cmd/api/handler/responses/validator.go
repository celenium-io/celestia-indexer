// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"fmt"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/types"
)

type Validator struct {
	Id          uint64 `example:"321"                                                    json:"id"           swaggertype:"integer"`
	Delegator   string `example:"celestia1un77nfm6axkhkupe8fk4xl6fd4adz3y5qk7ph6"        json:"delegator"    swaggertype:"string"`
	Address     string `example:"celestiavaloper1un77nfm6axkhkupe8fk4xl6fd4adz3y59fucpu" json:"address"      swaggertype:"string"`
	ConsAddress string `example:"E641C7A2C964833E556AEF934FBF166B712874B6"               json:"cons_address" swaggertype:"string"`

	Moniker  string `example:"Easy 2 Stake"                   json:"moniker"  swaggertype:"string"`
	Website  string `example:"https://www.easy2stake.com/"    json:"website"  swaggertype:"string"`
	Identity string `example:"2C877AC873132C91"               json:"identity" swaggertype:"string"`
	Contacts string `example:"security@0xfury.com"            json:"contacts" swaggertype:"string"`
	Details  string `example:"Some long text about validator" json:"details"  swaggertype:"string"`

	Rate              string `example:"0.03" json:"rate"                swaggertype:"string"`
	MaxRate           string `example:"0.1"  json:"max_rate"            swaggertype:"string"`
	MaxChangeRate     string `example:"0.01" json:"max_change_rate"     swaggertype:"string"`
	MinSelfDelegation string `example:"1"    json:"min_self_delegation" swaggertype:"string"`
}

func NewValidator(val storage.Validator) *Validator {
	if val.Id == 0 { // for genesis block
		return nil
	}
	return &Validator{
		Id:                val.Id,
		Delegator:         val.Delegator,
		Address:           val.Address,
		ConsAddress:       val.ConsAddress,
		Moniker:           val.Moniker,
		Website:           val.Website,
		Identity:          val.Identity,
		Contacts:          val.Contacts,
		Details:           val.Details,
		Rate:              val.Rate.String(),
		MaxRate:           val.MaxRate.String(),
		MaxChangeRate:     val.MaxChangeRate.String(),
		MinSelfDelegation: val.MinSelfDelegation.String(),
	}
}

type ShortValidator struct {
	Id          uint64 `example:"321"                                      json:"id"           swaggertype:"integer"`
	ConsAddress string `example:"E641C7A2C964833E556AEF934FBF166B712874B6" json:"cons_address" swaggertype:"string"`
	Moniker     string `example:"Easy 2 Stake"                             json:"moniker"      swaggertype:"string"`
}

func NewShortValidator(val storage.Validator) *ShortValidator {
	if val.Id == 0 { // for genesis block
		return nil
	}
	return &ShortValidator{
		Id:          val.Id,
		ConsAddress: val.ConsAddress,
		Moniker:     val.Moniker,
	}
}

type ValidatorUptime struct {
	Uptime string         `example:"0.97" json:"uptime" swaggertype:"string"`
	Blocks []SignedBlocks `json:"blocks"`
}

type SignedBlocks struct {
	Height types.Level `example:"100"  json:"height" swaggertype:"integer"`
	Signed bool        `example:"true" json:"signed" swaggertype:"boolean"`
}

func NewValidatorUptime(levels []types.Level, currentLevel types.Level, count types.Level) (uptime ValidatorUptime) {
	var (
		levelIndex = 0
		blockIndex = 0
		threshold  = count
	)

	if threshold > currentLevel {
		threshold = currentLevel
	}

	uptime.Blocks = make([]SignedBlocks, threshold)
	for i := currentLevel; i > currentLevel-threshold; i-- {
		if levelIndex < len(levels) && levels[levelIndex] == i {
			levelIndex++
			uptime.Blocks[blockIndex] = SignedBlocks{
				Signed: true,
				Height: i,
			}
		} else {
			uptime.Blocks[blockIndex] = SignedBlocks{
				Signed: false,
				Height: i,
			}
		}
		blockIndex++
	}

	uptime.Uptime = fmt.Sprintf("%.4f", float64(levelIndex)/float64(threshold))
	return uptime
}
