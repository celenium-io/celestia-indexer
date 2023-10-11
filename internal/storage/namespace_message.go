// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"time"

	"github.com/uptrace/bun"
)

type NamespaceMessage struct {
	bun.BaseModel `bun:"namespace_message" comment:"Table with relation messages to namespace."`

	NamespaceId uint64 `bun:"namespace_id,pk" comment:"Namespace internal id"`
	MsgId       uint64 `bun:"msg_id,pk"       comment:"Message id"`
	TxId        uint64 `bun:"tx_id"           comment:"Transaction id"`

	Time   time.Time   `bun:"time,notnull" comment:"Message time"`
	Height types.Level `bun:"height"       comment:"Message block height"`

	Message   *Message   `bun:"rel:belongs-to,join:msg_id=id"`
	Namespace *Namespace `bun:"rel:belongs-to,join:namespace_id=id"`
	Tx        *Tx        `bun:"rel:belongs-to,join:tx_id=id"`
}

func (NamespaceMessage) TableName() string {
	return "namespace_message"
}
