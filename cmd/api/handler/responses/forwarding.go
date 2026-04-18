// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"
	"time"

	"github.com/celenium-io/celestia-indexer/cmd/api/hyperlane"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
)

type Forwarding struct {
	Id          uint64         `example:"321"                                                                          format:"int64"     json:"id"           swaggertype:"integer"`
	Height      pkgTypes.Level `example:"100"                                                                          format:"int64"     json:"height"       swaggertype:"integer"`
	Time        time.Time      `example:"2023-07-04T03:10:57+00:00"                                                    format:"date-time" json:"time"         swaggertype:"string"`
	TxHash      string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF"             format:"binary"    json:"tx_hash"      swaggertype:"string"`
	DestDomain  uint64         `example:"11155111"                                                                     format:"int64"     json:"dest_domain"  swaggertype:"integer"`
	DestAddress string         `example:"0x000000000000000000000000d5e85e86fc692cedad6d6992f1f0ccf273e39913"           format:"binary"    json:"dest_address" swaggertype:"string"`
	Amount      string         `example:"1000000"                                                                      format:"string"    json:"amount"       swaggertype:"string"`
	Denom       string         `example:"hyperlane/0x726f757465725f61707000000000000000000000000000020000000000000024" format:"string"    json:"denom"        swaggertype:"string"`
	MessageId   string         `example:"0xac8852bd411c0c88cdadfe9b2386b2bcd702f35479c25a4b2d2cc3fb49d095d4"           format:"string"    json:"message_id"   swaggertype:"string"`
	TokenId     string         `example:"12652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF345"        format:"string"    json:"token_id"     swaggertype:"string"`

	Chain          *ChainMetadata    `json:"chain,omitempty"`
	ForwardAddress *ShortAddress     `json:"forward_address,omitempty"`
	Inputs         []ForwardingInput `json:"inputs,omitempty"`
}

func NewForwarding(forwarding storage.Forwarding, store hyperlane.IChainStore) Forwarding {
	response := Forwarding{
		Id:         forwarding.Id,
		Time:       forwarding.Time,
		Height:     forwarding.Height,
		DestDomain: forwarding.DestDomain,
		Amount:     forwarding.Amount.String(),
		Denom:      forwarding.Denom,
		MessageId:  forwarding.MessageId,
	}

	if forwarding.Tx != nil {
		response.TxHash = hex.EncodeToString(forwarding.Tx.Hash)
	}

	if forwarding.Address != nil {
		response.ForwardAddress = NewShortAddress(forwarding.Address)
	}

	if len(forwarding.DestRecipient) > 0 {
		response.DestAddress = hex.EncodeToString(forwarding.DestRecipient)
	}

	if forwarding.Token != nil {
		response.TokenId = hex.EncodeToString(forwarding.Token.TokenId)
	}

	if store != nil {
		response.Chain = NewChainMetadata(forwarding.DestDomain, store)
	}

	return response
}

type ForwardingInput struct {
	Height pkgTypes.Level `example:"100"                                                              format:"int64"     json:"height"   swaggertype:"integer"`
	Time   time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"     swaggertype:"string"`
	TxHash string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"tx_hash"  swaggertype:"string"`
	From   string         `example:"0x000000000000000000000000123456789abcdef123456789abcdef12345609" format:"string"    json:"from"     swaggertype:"string"`
	Amount string         `example:"123445"                                                           format:"string"    json:"received" swaggertype:"string"`
	Denom  string         `example:"utia"                                                             format:"string"    json:"denom"    swaggertype:"string"`

	Chain *ChainMetadata `json:"chain,omitempty"`
}

func NewForwardingInputFromHyperlaneTransfer(input storage.ForwardingInput, store hyperlane.IChainStore) ForwardingInput {
	response := ForwardingInput{
		Height: input.Height,
		Time:   input.Time,
		TxHash: hex.EncodeToString(input.TxHash),
		From:   input.From,
		Amount: input.Amount,
		Denom:  input.Denom,
	}
	if input.Counterparty > 0 {
		response.Chain = NewChainMetadata(input.Counterparty, store)
	}
	return response
}
