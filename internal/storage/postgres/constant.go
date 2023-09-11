package postgres

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-net/go-lib/database"
)

// Constant -
type Constant struct {
	db *database.Bun
}

// NewConstant -
func NewConstant(db *database.Bun) *Constant {
	return &Constant{
		db: db,
	}
}

func (constant *Constant) Get(ctx context.Context, module types.ModuleName, name string) (c storage.Constant, err error) {
	err = constant.db.DB().NewSelect().Model(&c).
		Where("module = ?", module).
		Where("name = ?", name).
		Scan(ctx)
	return
}

func (constant *Constant) ByModule(ctx context.Context, module types.ModuleName) (c []storage.Constant, err error) {
	err = constant.db.DB().NewSelect().Model(&c).
		Where("module = ?", module).
		Scan(ctx)
	return
}

func (constant *Constant) All(ctx context.Context) (c []storage.Constant, err error) {
	err = constant.db.DB().NewSelect().Model(&c).Scan(ctx)
	return
}
