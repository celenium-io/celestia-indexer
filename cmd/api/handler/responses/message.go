// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
)

type Message struct {
	Id       uint64         `example:"321"                       format:"int64"     json:"id"              swaggertype:"integer"`
	Height   pkgTypes.Level `example:"100"                       format:"int64"     json:"height"          swaggertype:"integer"`
	Time     time.Time      `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"            swaggertype:"string"`
	Position int64          `example:"2"                         format:"int64"     json:"position"        swaggertype:"integer"`
	TxId     uint64         `example:"11"                        format:"int64"     json:"tx_id,omitempty" swaggertype:"integer"`

	Type types.MsgType `example:"MsgCreatePeriodicVestingAccount" json:"type"`

	Data map[string]any `json:"data"`

	Tx *Tx `json:"tx,omitempty"`
}

func NewMessage(msg storage.Message) Message {
	return Message{
		Id:       msg.Id,
		Height:   msg.Height,
		Time:     msg.Time,
		Position: msg.Position,
		Type:     msg.Type,
		TxId:     msg.TxId,
		Data:     msg.Data,
	}
}

func NewMessageForAddress(msg storage.MsgAddress) Message {
	return Message{
		Id:       msg.MsgId,
		Height:   msg.Msg.Height,
		Time:     msg.Msg.Time,
		Position: msg.Msg.Position,
		TxId:     msg.Msg.TxId,
		Type:     msg.Msg.Type,
		Data:     msg.Msg.Data,
	}
}

func NewMessageWithTx(msg storage.MessageWithTx) Message {
	message := Message{
		Id:       msg.Id,
		Height:   msg.Height,
		Time:     msg.Time,
		Position: msg.Position,
		Type:     msg.Type,
		TxId:     msg.TxId,
		Data:     msg.Data,
	}

	if msg.Tx != nil {
		tx := NewTx(*msg.Tx)
		message.Tx = &tx
	}

	return message
}
