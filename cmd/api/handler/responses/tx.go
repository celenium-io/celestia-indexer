// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
)

type Tx struct {
	Id            uint64         `example:"321"                                                              format:"int64"     json:"id"                  swaggertype:"integer"`
	Height        pkgTypes.Level `example:"100"                                                              format:"int64"     json:"height"              swaggertype:"integer"`
	Position      int64          `example:"11"                                                               format:"int64"     json:"position"            swaggertype:"integer"`
	GasWanted     int64          `example:"9348"                                                             format:"int64"     json:"gas_wanted"          swaggertype:"integer"`
	GasUsed       int64          `example:"4253"                                                             format:"int64"     json:"gas_used"            swaggertype:"integer"`
	TimeoutHeight uint64         `example:"0"                                                                format:"int64"     json:"timeout_height"      swaggertype:"integer"`
	EventsCount   int64          `example:"2"                                                                format:"int64"     json:"events_count"        swaggertype:"integer"`
	MessagesCount int64          `example:"1"                                                                format:"int64"     json:"messages_count"      swaggertype:"integer"`
	Hash          string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"hash"                swaggertype:"string"`
	Fee           string         `example:"9348"                                                             format:"int64"     json:"fee"                 swaggertype:"string"`
	Error         string         `example:""                                                                 format:"string"    json:"error,omitempty"     swaggertype:"string"`
	Codespace     string         `example:"sdk"                                                              format:"string"    json:"codespace,omitempty" swaggertype:"string"`
	Memo          string         `example:"Transfer to private account"                                      format:"string"    json:"memo,omitempty"      swaggertype:"string"`
	Time          time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"                swaggertype:"string"`

	Messages []Message `json:"messages,omitempty"`

	MessageTypes []types.MsgType `example:"MsgSend,MsgUnjail" json:"message_types"`
	Status       types.Status    `example:"success"           json:"status"`

	MsgTypeMask types.MsgTypeBits `json:"-"`
}

func NewTx(tx storage.Tx) Tx {
	result := Tx{
		Id:            tx.Id,
		Height:        tx.Height,
		Time:          tx.Time,
		Position:      tx.Position,
		GasWanted:     tx.GasWanted,
		GasUsed:       tx.GasUsed,
		TimeoutHeight: tx.TimeoutHeight,
		EventsCount:   tx.EventsCount,
		MessagesCount: tx.MessagesCount,
		Fee:           tx.Fee.String(),
		Status:        tx.Status,
		Error:         tx.Error,
		Codespace:     tx.Codespace,
		Hash:          hex.EncodeToString(tx.Hash),
		Memo:          tx.Memo,
		MessageTypes:  tx.MessageTypes.Names(),
		MsgTypeMask:   tx.MessageTypes,
		Messages:      make([]Message, 0),
	}

	for i := range tx.Messages {
		result.Messages = append(result.Messages, NewMessage(tx.Messages[i]))
	}

	return result
}

func (Tx) SearchType() string {
	return "tx"
}

type TxForAddress struct {
	MessagesCount int64  `example:"1"                                                                format:"int64"  json:"messages_count" swaggertype:"integer"`
	Hash          string `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary" json:"hash"           swaggertype:"string"`
	Fee           string `example:"9348"                                                             format:"int64"  json:"fee"            swaggertype:"string"`

	MessageTypes []types.MsgType `example:"MsgSend,MsgUnjail" json:"message_types"`
	Status       types.Status    `example:"success"           json:"status"`
}

func NewTxForAddress(tx *storage.Tx) TxForAddress {
	result := TxForAddress{
		MessagesCount: tx.MessagesCount,
		Fee:           tx.Fee.String(),
		Status:        tx.Status,
		Hash:          hex.EncodeToString(tx.Hash),
		MessageTypes:  tx.MessageTypes.Names(),
	}

	return result
}
