// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

type BlobLogFilters struct {
	Limit  int
	Offset int
	Sort   sdk.SortOrder
	SortBy string
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IBlobLog interface {
	storage.Table[*BlobLog]

	ByNamespace(ctx context.Context, nsId uint64, fltrs BlobLogFilters) ([]BlobLog, error)
	ByProviders(ctx context.Context, providers []RollupProvider, fltrs BlobLogFilters) ([]BlobLog, error)
}

type BlobLog struct {
	bun.BaseModel `bun:"blob_log" comment:"Table with flatted blob entities."`

	Id         uint64      `bun:"id,pk,autoincrement" comment:"Unique internal identity"`
	Time       time.Time   `bun:"time,notnull,pk"     comment:"Message time"`
	Height     types.Level `bun:"height"              comment:"Message block height"`
	Size       int64       `bun:"size"                comment:"Blob size"`
	Commitment string      `bun:"commitment"          comment:"Blob commitment"`

	SignerId    uint64 `bun:"signer_id"    comment:"Blob signer identity"`
	NamespaceId uint64 `bun:"namespace_id" comment:"Namespace internal id"`
	MsgId       uint64 `bun:"msg_id"       comment:"Message id"`
	TxId        uint64 `bun:"tx_id"        comment:"Transaction id"`

	Message   *Message   `bun:"rel:belongs-to,join:msg_id=id"`
	Namespace *Namespace `bun:"rel:belongs-to,join:namespace_id=id"`
	Tx        *Tx        `bun:"rel:belongs-to,join:tx_id=id"`
	Signer    *Address   `bun:"rel:belongs-to,join:signer_id=id"`
}

func (BlobLog) TableName() string {
	return "blob_log"
}
