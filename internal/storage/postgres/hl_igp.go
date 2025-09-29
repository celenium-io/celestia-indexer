package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

type HLIGP struct {
	*database.Bun
}

func NewHLIGP(conn *database.Bun) *HLIGP {
	return &HLIGP{conn}
}

func (hl *HLIGP) List(ctx context.Context, limit, offset int) (igp []storage.HLIGP, err error) {
	query := hl.DB().NewSelect().
		Model((*storage.HLIGP)(nil))

	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}

	err = hl.DB().NewSelect().
		TableExpr("(?) as igp", query).
		ColumnExpr("igp.*").
		ColumnExpr("hl_igp_config.gas_overhead as config__gas_overhead, hl_igp_config.gas_price as config__gas_price, hl_igp_config.remote_domain as config__remote_domain, hl_igp_config.token_exchange_rate as config__token_exchange_rate").
		Join("left join hl_igp_config on hl_igp_config.id = igp.id").
		Scan(ctx, &igp)
	return
}

func (hl *HLIGP) ByHash(ctx context.Context, hash []byte) (igp storage.HLIGP, err error) {
	query := hl.DB().NewSelect().
		Model(&igp).
		Where("igp_id = ?", hash).
		Limit(1)

	err = hl.DB().NewSelect().
		TableExpr("(?) as igp", query).
		ColumnExpr("igp.*").
		ColumnExpr("hl_igp_config.gas_overhead as config__gas_overhead, hl_igp_config.gas_price as config__gas_price, hl_igp_config.remote_domain as config__remote_domain, hl_igp_config.token_exchange_rate as config__token_exchange_rate").
		Join("left join hl_igp_config on hl_igp_config.id = igp.id").
		Scan(ctx, &igp)
	return
}
