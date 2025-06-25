// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type ListHyperlaneTokens struct {
	Limit     int
	Offset    int
	Sort      sdk.SortOrder
	OwnerId   uint64
	MailboxId uint64
	Type      []types.HLTokenType
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IHLToken interface {
	ByHash(ctx context.Context, id []byte) (HLToken, error)
	List(ctx context.Context, fltrs ListHyperlaneTokens) ([]HLToken, error)
}

type HLToken struct {
	bun.BaseModel `bun:"hl_token" comment:"Table with hyperlane tokens"`

	Id               uint64            `bun:"id,pk,autoincrement"        comment:"Internal identity"`
	Height           pkgTypes.Level    `bun:"height,notnull"             comment:"The number (height) of this block"`
	Time             time.Time         `bun:"time,pk,notnull"            comment:"The time of block"`
	OwnerId          uint64            `bun:"owner_id"                   comment:"Owner internal identity"`
	MailboxId        uint64            `bun:"mailbox_id"                 comment:"Mailbox internal identity"`
	TxId             uint64            `bun:"tx_id"                      comment:"Transaction identity"`
	Type             types.HLTokenType `bun:",type:hyperlane_token_type" comment:"Token type: synthetic or collateral"`
	Denom            string            `bun:"denom"                      comment:"Denom"`
	TokenId          []byte            `bun:"token_id,type:bytea,unique" comment:"Token id"`
	SentTransfers    uint64            `bun:"sent_transfers"             comment:"Sent transfers"`
	ReceiveTransfers uint64            `bun:"received_transfers"         comment:"Receive transfers"`
	Sent             decimal.Decimal   `bun:"sent,type:numeric"          comment:"Sent tokens"`
	Received         decimal.Decimal   `bun:"received,type:numeric"      comment:"Receive tokens"`

	Owner   *Address   `bun:"rel:belongs-to,join:owner_id=id"`
	Mailbox *HLMailbox `bun:"rel:belongs-to,join:mailbox_id=id"`
	Tx      *Tx        `bun:"rel:belongs-to,join:tx_id=id"`
}

func (HLToken) TableName() string {
	return "hl_token"
}

func (t *HLToken) String() string {
	return hex.EncodeToString(t.TokenId)
}
