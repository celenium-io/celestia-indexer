// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"fmt"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/shopspring/decimal"
)

type Validator struct {
	Id          uint64 `example:"321"                                      json:"id"           swaggertype:"integer"`
	ConsAddress string `example:"E641C7A2C964833E556AEF934FBF166B712874B6" json:"cons_address" swaggertype:"string"`

	Moniker  string `example:"Easy 2 Stake"                   json:"moniker"  swaggertype:"string"`
	Website  string `example:"https://www.easy2stake.com/"    json:"website"  swaggertype:"string"`
	Identity string `example:"2C877AC873132C91"               json:"identity" swaggertype:"string"`
	Contacts string `example:"security@0xfury.com"            json:"contacts" swaggertype:"string"`
	Details  string `example:"Some long text about validator" json:"details"  swaggertype:"string"`

	Rate              string `example:"0.03" json:"rate"                swaggertype:"string"`
	MaxRate           string `example:"0.1"  json:"max_rate"            swaggertype:"string"`
	MaxChangeRate     string `example:"0.01" json:"max_change_rate"     swaggertype:"string"`
	MinSelfDelegation string `example:"1"    json:"min_self_delegation" swaggertype:"string"`
	Stake             string `example:"1"    json:"stake"               swaggertype:"string"`
	Rewards           string `example:"1"    json:"rewards"             swaggertype:"string"`
	Commissions       string `example:"1"    json:"commissions"         swaggertype:"string"`
	VotingPower       string `example:"1"    json:"voting_power"        swaggertype:"string"`

	Jailed bool `example:"false" json:"jailed" swaggertype:"boolean"`

	Address   *ShortAddress `json:"address"`
	Delegator *ShortAddress `json:"delegator"`
}

func NewValidator(val storage.Validator) *Validator {
	if val.Id == 0 { // for genesis block
		return nil
	}
	return &Validator{
		Id: val.Id,
		Delegator: &ShortAddress{
			Hash: val.Delegator,
		},
		Address: &ShortAddress{
			Hash: val.Address,
		},
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
		Stake:             val.Stake.String(),
		Rewards:           val.Rewards.Floor().String(),
		Commissions:       val.Commissions.Floor().String(),
		Jailed:            *val.Jailed,
		VotingPower:       val.Stake.Div(decimal.NewFromInt(1_000_000)).Floor().String(),
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
		threshold  = count
	)

	if threshold > currentLevel {
		threshold = currentLevel
	}

	uptime.Blocks = make([]SignedBlocks, 0)
	for i := currentLevel; i > currentLevel-threshold; i-- {
		if levelIndex < len(levels) && levels[levelIndex] == i {
			levelIndex++
			uptime.Blocks = append(uptime.Blocks, SignedBlocks{
				Signed: true,
				Height: i,
			})
		} else {
			uptime.Blocks = append(uptime.Blocks, SignedBlocks{
				Signed: false,
				Height: i,
			})
		}
	}

	uptime.Uptime = fmt.Sprintf("%.4f", float64(levelIndex)/float64(threshold))
	return uptime
}

type Jail struct {
	Height types.Level `example:"100"                       json:"height" swaggertype:"integer"`
	Time   time.Time   `example:"2023-07-04T03:10:57+00:00" json:"time"   swaggertype:"string"`
	Reason string      `example:"double_sign"               json:"reason" swaggertype:"string"`
	Burned string      `example:"10000000000"               json:"burned" swaggertype:"string"`

	Validator *ShortValidator `json:"validator,omitempty"`
}

func NewJail(jail storage.Jail) Jail {
	j := Jail{
		Height: jail.Height,
		Time:   jail.Time,
		Reason: jail.Reason,
		Burned: jail.Burned.String(),
	}

	if jail.Validator != nil {
		j.Validator = NewShortValidator(*jail.Validator)
	}

	return j
}

type ValidatorCount struct {
	Total    int `example:"100" json:"total"    swaggertype:"integer"`
	Jailed   int `example:"100" json:"jailed"   swaggertype:"integer"`
	Active   int `example:"100" json:"active"   swaggertype:"integer"`
	Inactive int `example:"100" json:"inactive" swaggertype:"integer"`
}
