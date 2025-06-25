// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"io"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type BlobLogFilters struct {
	Limit      int
	Offset     int
	Sort       sdk.SortOrder
	SortBy     string
	From       time.Time
	To         time.Time
	Commitment string
	Joins      bool
	Signers    []uint64
	Cursor     uint64
}

type ListBlobLogFilters struct {
	Limit      int
	Offset     int
	Sort       sdk.SortOrder
	SortBy     string
	From       time.Time
	To         time.Time
	Commitment string
	Signers    []uint64
	Namespaces []uint64
	Cursor     uint64
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IBlobLog interface {
	sdk.Table[*BlobLog]

	ByNamespace(ctx context.Context, nsId uint64, fltrs BlobLogFilters) ([]BlobLog, error)
	ByProviders(ctx context.Context, providers []RollupProvider, fltrs BlobLogFilters) ([]BlobLog, error)
	BySigner(ctx context.Context, signerId uint64, fltrs BlobLogFilters) ([]BlobLog, error)
	ByTxId(ctx context.Context, txId uint64, fltrs BlobLogFilters) ([]BlobLog, error)
	ByHeight(ctx context.Context, height types.Level, fltrs BlobLogFilters) ([]BlobLog, error)
	CountByTxId(ctx context.Context, txId uint64) (int, error)
	ExportByProviders(ctx context.Context, providers []RollupProvider, from, to time.Time, stream io.Writer) (err error)
	Blob(ctx context.Context, height types.Level, nsId uint64, commitment string) (BlobLog, error)
	ListBlobs(ctx context.Context, fltrs ListBlobLogFilters) ([]BlobLog, error)
}

type BlobLog struct {
	bun.BaseModel `bun:"blob_log" comment:"Table with flatted blob entities."`

	Id          uint64          `bun:"id,pk,autoincrement" comment:"Unique internal identity"`
	Time        time.Time       `bun:"time,notnull,pk"     comment:"Message time"`
	Height      types.Level     `bun:"height"              comment:"Message block height"`
	Size        int64           `bun:"size"                comment:"Blob size"`
	Commitment  string          `bun:"commitment"          comment:"Blob commitment"`
	ContentType string          `bun:"content_type"        comment:"Blob content type"`
	Fee         decimal.Decimal `bun:"fee,type:numeric"    comment:"Fee per blob"`

	SignerId    uint64 `bun:"signer_id"    comment:"Blob signer identity"`
	NamespaceId uint64 `bun:"namespace_id" comment:"Namespace internal id"`
	MsgId       uint64 `bun:"msg_id"       comment:"Message id"`
	TxId        uint64 `bun:"tx_id"        comment:"Transaction id"`

	Message   *Message   `bun:"rel:belongs-to,join:msg_id=id"`
	Namespace *Namespace `bun:"rel:belongs-to,join:namespace_id=id"`
	Tx        *Tx        `bun:"rel:belongs-to,join:tx_id=id"`
	Signer    *Address   `bun:"rel:belongs-to,join:signer_id=id"`
	Rollup    *Rollup    `bun:"rel:belongs-to"`
}

func (BlobLog) TableName() string {
	return "blob_log"
}
