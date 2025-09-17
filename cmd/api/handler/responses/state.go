// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"

	"github.com/celenium-io/celestia-indexer/internal/storage"
)

type State struct {
	Id               uint64         `example:"321"                                                              format:"int64"     json:"id"                 swaggertype:"integer"`
	Version          uint64         `example:"5"                                                                format:"int64"     json:"version"            swaggertype:"integer"`
	Name             string         `example:"indexer"                                                          format:"string"    json:"name"               swaggertype:"string"`
	ChainId          string         `example:"mocha-4"                                                          format:"string"    json:"chain_id"           swaggertype:"string"`
	LastHeight       pkgTypes.Level `example:"100"                                                              format:"int64"     json:"last_height"        swaggertype:"integer"`
	LastHash         string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"string"    json:"hash"               swaggertype:"string"`
	LastTime         time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"last_time"          swaggertype:"string"`
	TotalTx          int64          `example:"23456"                                                            format:"int64"     json:"total_tx"           swaggertype:"integer"`
	TotalAccounts    int64          `example:"43"                                                               format:"int64"     json:"total_accounts"     swaggertype:"integer"`
	TotalFee         string         `example:"312"                                                              format:"string"    json:"total_fee"          swaggertype:"string"`
	TotalBlobsSize   int64          `example:"56789"                                                            format:"int64"     json:"total_blobs_size"   swaggertype:"integer"`
	TotalProposals   int64          `example:"56789"                                                            format:"int64"     json:"total_proposals"    swaggertype:"integer"`
	TotalValidators  int            `example:"100"                                                              format:"int64"     json:"total_validators"   swaggertype:"integer"`
	TotalSupply      string         `example:"312"                                                              format:"string"    json:"total_supply"       swaggertype:"string"`
	TotalVotingPower string         `example:"312"                                                              format:"string"    json:"total_voting_power" swaggertype:"string"`
	TotalNamespaces  int64          `example:"312"                                                              format:"string"    json:"total_namespaces"   swaggertype:"integer"`
	TotalIbcClients  int64          `example:"312"                                                              format:"string"    json:"total_ibc_clients"  swaggertype:"integer"`
	Synced           bool           `example:"true"                                                             format:"boolean"   json:"synced"             swaggertype:"boolean"`
}

func NewState(state storage.State) State {
	return State{
		Id:               state.Id,
		Version:          state.Version,
		Name:             state.Name,
		ChainId:          state.ChainId,
		LastHeight:       state.LastHeight,
		LastHash:         hex.EncodeToString(state.LastHash),
		LastTime:         state.LastTime,
		TotalTx:          state.TotalTx,
		TotalAccounts:    state.TotalAccounts,
		TotalFee:         state.TotalFee.String(),
		TotalBlobsSize:   state.TotalBlobsSize,
		TotalValidators:  state.TotalValidators,
		TotalNamespaces:  state.TotalNamespaces,
		TotalProposals:   state.TotalProposals,
		TotalSupply:      state.TotalSupply.String(),
		TotalVotingPower: state.TotalVotingPower.String(),
		TotalIbcClients:  state.TotalIbcClients,
		Synced:           !state.LastTime.UTC().Add(2 * time.Minute).Before(time.Now().UTC()),
	}
}
