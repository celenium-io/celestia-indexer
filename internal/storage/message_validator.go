// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"fmt"
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/uptrace/bun"
)

type MsgValidator struct {
	bun.BaseModel `bun:"msg_validator" comment:"Table with relation message to validator"`

	ValidatorId uint64         `bun:"validator_id,pk" comment:"Validator internal id"`
	MsgId       uint64         `bun:"msg_id,pk"       comment:"Message internal id"`
	Time        time.Time      `bun:"time,pk,notnull" comment:"The time of block"`
	Height      pkgTypes.Level `bun:",notnull"        comment:"The number (height) of this block"`

	Validator *Validator `bun:"rel:belongs-to,join:validator_id=id"`
	Msg       *Message   `bun:"rel:belongs-to,join:msg_id=id"`
}

func (MsgValidator) TableName() string {
	return "msg_validator"
}

func (m MsgValidator) String() string {
	return fmt.Sprintf("%d_%d", m.ValidatorId, m.MsgId)
}
