package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

type HLGasPayment struct {
	*database.Bun
}

func NewHLGasPayment(conn *database.Bun) *HLGasPayment {
	return &HLGasPayment{conn}
}

func (hl *HLGasPayment) List(ctx context.Context, limit, offset int) (payments []storage.HLGasPayment, err error) {
	query := hl.DB().NewSelect().
		Model((*storage.HLGasPayment)(nil))

	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}

	err = query.Scan(ctx, &payments)

	return
}
