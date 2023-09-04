package postgres

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Tx -
type Tx struct {
	*postgres.Table[*storage.Tx]
}

// NewTx -
func NewTx(db *database.Bun) *Tx {
	return &Tx{
		Table: postgres.NewTable[*storage.Tx](db),
	}
}

func (tx *Tx) ByHash(ctx context.Context, hash []byte) (transaction storage.Tx, err error) {
	err = tx.DB().NewSelect().Model(&transaction).
		Where("hash = ?", hash).
		Scan(ctx)
	return
}

func (tx *Tx) Filter(ctx context.Context, fltrs storage.TxFilter) (txs []storage.Tx, err error) {
	query := tx.DB().NewSelect().Model(&txs).Offset(fltrs.Offset)
	query = limitScope(query, fltrs.Limit)
	query = sortScope(query, "id", fltrs.Sort)

	if !fltrs.MessageTypes.Empty() {
		query = query.Where("message_types & ? > 0", fltrs.MessageTypes)
	}

	if len(fltrs.Status) > 0 {
		query = query.WhereGroup(" AND ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			for i := range fltrs.Status {
				sq = sq.WhereOr("status = ?", fltrs.Status[i])
			}
			return sq
		})
	}
	if fltrs.Height > 0 {
		query = query.Where("height = ?", fltrs.Height)
	}

	if !fltrs.TimeFrom.IsZero() {
		query = query.Where("time >= ?", fltrs.TimeFrom)
	}
	if !fltrs.TimeTo.IsZero() {
		query = query.Where("time < ?", fltrs.TimeTo)
	}

	err = query.Scan(ctx)
	return
}

func (tx *Tx) ByIdWithRelations(ctx context.Context, id uint64) (transaction storage.Tx, err error) {
	err = tx.DB().NewSelect().Model(&transaction).
		Where("id = ?", id).
		Relation("Messages").
		Scan(ctx)
	return
}
