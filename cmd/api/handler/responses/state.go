// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
)

type State struct {
	Id             uint64         `example:"321"                                                              format:"int64"     json:"id"               swaggertype:"integer"`
	Name           string         `example:"indexer"                                                          format:"string"    json:"name"             swaggertype:"string"`
	LastHeight     pkgTypes.Level `example:"100"                                                              format:"int64"     json:"last_height"      swaggertype:"integer"`
	LastHash       string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"string"    json:"hash"             swaggertype:"string"`
	LastTime       time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"last_time"        swaggertype:"string"`
	TotalTx        int64          `example:"23456"                                                            format:"int64"     json:"total_tx"         swaggertype:"integer"`
	TotalAccounts  int64          `example:"43"                                                               format:"int64"     json:"total_accounts"   swaggertype:"integer"`
	TotalFee       string         `example:"312"                                                              format:"string"    json:"total_fee"        swaggertype:"string"`
	TotalBlobsSize int64          `example:"56789"                                                            format:"int64"     json:"total_blobs_size" swaggertype:"integer"`
	TotalSupply    string         `example:"312"                                                              format:"string"    json:"total_supply"     swaggertype:"string"`
}

func NewState(state storage.State) State {
	return State{
		Id:             state.Id,
		Name:           state.Name,
		LastHeight:     state.LastHeight,
		LastHash:       hex.EncodeToString(state.LastHash),
		LastTime:       state.LastTime,
		TotalTx:        state.TotalTx,
		TotalAccounts:  state.TotalAccounts,
		TotalFee:       state.TotalFee.String(),
		TotalBlobsSize: state.TotalBlobsSize,
		TotalSupply:    state.TotalSupply.String(),
	}
}
