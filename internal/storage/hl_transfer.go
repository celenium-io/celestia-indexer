// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type ListHyperlaneTransfers struct {
	Limit     int
	Offset    int
	Sort      sdk.SortOrder
	MailboxId uint64
	AddressId uint64
	TokenId   uint64
	RelayerId uint64
	Type      []types.HLTransferType
	Domain    uint64
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IHLTransfer interface {
	List(ctx context.Context, filters ListHyperlaneTransfers) ([]HLTransfer, error)
}

type HLTransfer struct {
	bun.BaseModel `bun:"hl_transfer" comment:"Table with hyperlane transfers"`

	Id                  uint64               `bun:"id,pk,autoincrement"           comment:"Internal identity"`
	Height              pkgTypes.Level       `bun:"height,notnull"                comment:"The number (height) of this block"`
	Time                time.Time            `bun:"time,pk,notnull"               comment:"The time of block"`
	TxId                uint64               `bun:"tx_id"                         comment:"Transaction identity"`
	MailboxId           uint64               `bun:"mailbox_id"                    comment:"Mailbox address"`
	RelayerId           uint64               `bun:"relayer_id"                    comment:"Relayer address"`
	TokenId             uint64               `bun:"token_id"                      comment:"Token id"`
	Counterparty        uint64               `bun:"counterparty"                  comment:"Counterparty domain"`
	AddressId           uint64               `bun:"address_id"                    comment:"Internal celestia address identity"`
	CounterpartyAddress string               `bun:"counterparty_address"          comment:"Counterparty address"`
	Version             byte                 `bun:"version"                       comment:"Version"`
	Nonce               uint32               `bun:"nonce"                         comment:"Nonce"`
	Body                []byte               `bun:"body,type:bytea,nullzero"      comment:"Body"`
	Metadata            []byte               `bun:"metadata,type:bytea,nullzero"  comment:"Metadata"`
	Type                types.HLTransferType `bun:",type:hyperlane_transfer_type" comment:"Transfer type"`
	Amount              decimal.Decimal      `bun:"amount,type:numeric"           comment:"Amount"`
	Denom               string               `bun:"denom"                         comment:"Denom"`

	Mailbox *HLMailbox `bun:"rel:belongs-to,join:mailbox_id=id"`
	Relayer *Address   `bun:"rel:belongs-to,join:relayer_id=id"`
	Address *Address   `bun:"rel:belongs-to,join:address_id=id"`
	Token   *HLToken   `bun:"rel:belongs-to,join:token_id=id"`
	Tx      *Tx        `bun:"rel:belongs-to,join:tx_id=id"`
}

func (m *HLTransfer) TableName() string {
	return "hl_transfer"
}
