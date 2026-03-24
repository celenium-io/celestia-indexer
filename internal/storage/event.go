// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"bytes"
	"context"
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/celenium-io/celestia-indexer/internal/pool"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

type EventFilter struct {
	Limit  int
	Offset int
	Time   time.Time
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IEvent interface {
	sdk.Table[*Event]

	ByTxId(ctx context.Context, txId uint64, fltrs EventFilter) ([]Event, error)
	ByBlock(ctx context.Context, height pkgTypes.Level, fltrs EventFilter) ([]Event, error)
}

var _ sdk.Copiable = (*Event)(nil)

// Event -
type Event struct {
	bun.BaseModel `bun:"event" comment:"Table with celestia events."`

	Id       uint64            `bun:"id,pk,notnull,autoincrement"      comment:"Unique internal id"`
	Height   pkgTypes.Level    `bun:"height,notnull"                   comment:"The number (height) of this block" stats:"func:min max,filterable"`
	Time     time.Time         `bun:"time,pk,notnull"                  comment:"The time of block"                 stats:"func:min max,filterable"`
	Position int64             `bun:"position"                         comment:"Position in transaction"`
	Type     types.EventType   `bun:",type:event_type"                 comment:"Event type"                        stats:"filterable"`
	TxId     *uint64           `bun:"tx_id"                            comment:"Transaction id"`
	Data     map[string]string `bun:"data,msgpack,type:bytea,nullzero" comment:"Event data"`
}

// TableName -
func (Event) TableName() string {
	return "event"
}

var msgpackBufPool = pool.New(func() *bytes.Buffer { return new(bytes.Buffer) })

var msgpackEncoderPool = pool.New(func() *msgpack.Encoder {
	return msgpack.NewEncoder(new(bytes.Buffer))
})

func msgPackEncode(data any) ([]byte, error) {
	buf := msgpackBufPool.Get()
	buf.Reset()
	defer msgpackBufPool.Put(buf)

	enc := msgpackEncoderPool.Get()
	enc.Reset(buf)
	defer msgpackEncoderPool.Put(enc)

	if err := enc.Encode(data); err != nil {
		return nil, errors.Wrap(err, "msgpack marshal event's data")
	}
	return bytes.Clone(buf.Bytes()), nil
}

func (e Event) Flat() ([]any, error) {
	var (
		data []byte
		err  error
	)
	if len(e.Data) > 0 {
		data, err = msgPackEncode(e.Data)
		if err != nil {
			return nil, err
		}
	}

	var txID any
	if e.TxId != nil {
		txID = int64(*e.TxId)
	}
	return []any{
		int64(e.Height),
		e.Time,
		e.Position,
		string(e.Type),
		txID,
		data,
	}, nil
}

func (e Event) Columns() []string {
	return []string{
		"height", "time", "position", "type", "tx_id", "data",
	}
}
