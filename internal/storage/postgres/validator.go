package postgres

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Validator -
type Validator struct {
	*postgres.Table[*storage.Validator]
}

// NewValidator -
func NewValidator(db *database.Bun) *Validator {
	return &Validator{
		Table: postgres.NewTable[*storage.Validator](db),
	}
}

func (v *Validator) ByAddress(ctx context.Context, address string) (validator storage.Validator, err error) {
	err = v.DB().NewSelect().Model(&validator).
		Where("address = ?", address).
		Scan(ctx)
	return
}
