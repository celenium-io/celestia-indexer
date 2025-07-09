// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"encoding/hex"
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IHLMailbox interface {
	List(ctx context.Context, limit, offset int) ([]HLMailbox, error)
	ByHash(ctx context.Context, hash []byte) (HLMailbox, error)
}

type HLMailbox struct {
	bun.BaseModel `bun:"hl_mailbox" comment:"Table with hyperlane mailboxes"`

	Id               uint64         `bun:"id,pk,autoincrement" comment:"Internal identity"`
	Height           pkgTypes.Level `bun:"height,notnull"      comment:"The number (height) of this block"`
	Time             time.Time      `bun:"time,pk,notnull"     comment:"The time of block"`
	TxId             uint64         `bun:"tx_id"               comment:"Internal creation transaction id"`
	Mailbox          []byte         `bun:"mailbox,unique"      comment:"Mailbox address"`
	InternalId       uint64         `bun:"internal_id,unique"  comment:"Internal mailbox id"`
	OwnerId          uint64         `bun:"owner_id"            comment:"Owner identity"`
	DefaultIsm       []byte         `bun:"default_ism"         comment:"Default ISM"`
	DefaultHook      []byte         `bun:"default_hook"        comment:"Default hook"`
	RequiredHook     []byte         `bun:"required_hook"       comment:"Required hook"`
	Domain           uint64         `bun:"domain,nullzero"     comment:"Domain"`
	SentMessages     uint64         `bun:"sent_messages"       comment:"Count of sent messages"`
	ReceivedMessages uint64         `bun:"received_messages"   comment:"Count of received messages"`

	Owner *Address `bun:"rel:belongs-to,join:owner_id=id"`
	Tx    *Tx      `bun:"rel:belongs-to,join:tx_id=id"`
}

func (m *HLMailbox) TableName() string {
	return "hl_mailbox"
}

func (m *HLMailbox) String() string {
	return hex.EncodeToString(m.Mailbox)
}
