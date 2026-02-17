// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
)

type Forwarding struct {
	Id             uint64          `example:"321"                                                              format:"int64"     json:"id"              swaggertype:"integer"`
	Height         pkgTypes.Level  `example:"100"                                                              format:"int64"     json:"height"          swaggertype:"integer"`
	Time           time.Time       `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"            swaggertype:"string"`
	TxHash         string          `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"tx_hash"         swaggertype:"string"`
	DestDomain     uint64          `example:"123456789"                                                        format:"int64"     json:"dest_domain"     swaggertype:"integer"`
	DestAddress    []byte          `example:"0x000000000000000000000000123456789abcdef123456789abcdef12345609" format:"binary"    json:"dest_address"    swaggertype:"string"`
	ForwardAddress string          `example:"celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60"                  format:"string"    json:"forward_address" swaggertype:"string"`
	SuccessCount   uint64          `example:"100"                                                              format:"int64"     json:"success_count"   swaggertype:"integer"`
	FailedCount    uint64          `example:"10"                                                               format:"int64"     json:"failed_count"    swaggertype:"integer"`
	Transfers      json.RawMessage `json:"transfers,omitempty"`
}

func NewForwarding(forwarding storage.Forwarding) Forwarding {
	response := Forwarding{
		Id:           forwarding.Id,
		Time:         forwarding.Time,
		Height:       forwarding.Height,
		DestDomain:   forwarding.DestDomain,
		DestAddress:  forwarding.DestRecipient,
		SuccessCount: forwarding.SuccessCount,
		FailedCount:  forwarding.FailedCount,
		Transfers:    forwarding.Transfers,
	}

	if forwarding.Tx != nil {
		response.TxHash = hex.EncodeToString(forwarding.Tx.Hash)
	}

	if forwarding.Address != nil {
		response.ForwardAddress = forwarding.Address.Address
	}

	return response
}
