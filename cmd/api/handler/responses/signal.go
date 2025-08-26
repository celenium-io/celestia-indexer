// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package responses

import (
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
)

type SignalVersion struct {
	Id          uint64         `example:"321"                       format:"int64"     json:"id"           swaggertype:"integer"`
	Height      pkgTypes.Level `example:"100"                       format:"int64"     json:"height"       swaggertype:"integer"`
	Time        time.Time      `example:"2025-07-04T03:10:57+00:00" format:"date-time" json:"time"         swaggertype:"string"`
	VotingPower string         `example:"9348"                      format:"int64"     json:"voting_power" swaggertype:"string"`
	Version     uint64         `example:"1"                         format:"int64"     json:"version"      swaggertype:"integer"`
	MsgId       uint64         `example:"2"                         format:"int64"     json:"msg_id"       swaggertype:"integer"`
	TxId        uint64         `example:"3"                         format:"int64"     json:"tx_id"        swaggertype:"integer"`
}

func NewSignalVersion(signal storage.SignalVersion) SignalVersion {
	result := SignalVersion{
		Id:          signal.Id,
		Height:      signal.Height,
		Time:        signal.Time,
		VotingPower: signal.VotingPower.String(),
		Version:     signal.Version,
		MsgId:       signal.MsgId,
		TxId:        signal.TxId,
	}

	return result
}

type Upgrade struct {
	Id      uint64         `example:"321"                       format:"int64"     json:"id"      swaggertype:"integer"`
	Height  pkgTypes.Level `example:"100"                       format:"int64"     json:"height"  swaggertype:"integer"`
	Time    time.Time      `example:"2025-07-04T03:10:57+00:00" format:"date-time" json:"time"    swaggertype:"string"`
	Version uint64         `example:"1"                         format:"int64"     json:"version" swaggertype:"integer"`
	MsgId   uint64         `example:"2"                         format:"int64"     json:"msg_id"  swaggertype:"integer"`
	TxId    uint64         `example:"3"                         format:"int64"     json:"tx_id"   swaggertype:"integer"`

	Signer *ShortAddress `json:"signer,omitempty"`
}

func NewUpgrade(upgrade storage.Upgrade) Upgrade {
	result := Upgrade{
		Id:      upgrade.Id,
		Height:  upgrade.Height,
		Time:    upgrade.Time,
		Version: upgrade.Version,
		MsgId:   upgrade.MsgId,
		TxId:    upgrade.TxId,
	}

	if upgrade.Signer != nil {
		result.Signer = NewShortAddress(upgrade.Signer)
	}

	return result
}
