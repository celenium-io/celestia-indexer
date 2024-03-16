// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
)

type Vesting struct {
	Id        uint64            `example:"12"                                                               format:"integer"   json:"id"                   swaggertype:"integer"`
	Height    pkgTypes.Level    `example:"100"                                                              format:"integer"   json:"height"               swaggertype:"integer"`
	Time      time.Time         `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"                 swaggertype:"string"`
	StartTime time.Time         `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"start_time,omitempty" swaggertype:"string"`
	EndTime   time.Time         `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"end_time,omitempty"   swaggertype:"string"`
	Hash      string            `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"hash,omitempty"       swaggertype:"string"`
	Type      types.VestingType `example:"delayed"                                                          format:"string"    json:"type"                 swaggertype:"string"`
	Amount    string            `example:"123.13333"                                                        format:"string"    json:"amount"               swaggertype:"string"`
}

func NewVesting(v storage.VestingAccount) Vesting {
	vesting := Vesting{
		Id:     v.Id,
		Height: v.Height,
		Time:   v.Time,
		Type:   v.Type,
		Amount: v.Amount.String(),
	}

	if v.StartTime != nil {
		vesting.StartTime = *v.StartTime
	}

	if v.EndTime != nil {
		vesting.EndTime = *v.EndTime
	}

	if v.Tx != nil {
		vesting.Hash = hex.EncodeToString(v.Tx.Hash)
	}

	return vesting
}

type VestingPeriod struct {
	Time   time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"   swaggertype:"string"`
	Amount string    `example:"123.13333"                 format:"string"    json:"amount" swaggertype:"string"`
}

func NewVestingPeriod(v storage.VestingPeriod) VestingPeriod {
	return VestingPeriod{
		Time:   v.Time,
		Amount: v.Amount.String(),
	}
}
