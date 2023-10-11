// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/uptrace/bun"
)

type MsgAddress struct {
	bun.BaseModel `bun:"msg_address" comment:"Table with relation tx to address"`

	AddressId uint64               `bun:"address_id,pk"             comment:"Address internal id"`
	MsgId     uint64               `bun:"msg_id,pk"                 comment:"Message internal id"`
	Type      types.MsgAddressType `bun:",pk,type:msg_address_type" comment:"The reason why address link to transaction"`

	Address *Address `bun:"rel:belongs-to,join:address_id=id"`
	Msg     *Message `bun:"rel:belongs-to,join:msg_id=id"`
}

func (MsgAddress) TableName() string {
	return "msg_address"
}

type AddressWithType struct {
	Address

	Type types.MsgAddressType
}
