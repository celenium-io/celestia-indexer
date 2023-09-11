package storage

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IConstant interface {
	Get(ctx context.Context, module types.ModuleName, name string) (Constant, error)
	ByModule(ctx context.Context, module types.ModuleName) ([]Constant, error)
	All(ctx context.Context) ([]Constant, error)
}

type Constant struct {
	bun.BaseModel `bun:"table:constant" comment:"Table with celestia constants."`

	Module types.ModuleName `bun:"module,pk,type:module_name" comment:"Module name which declares constant"`
	Name   string           `bun:"name,pk,type:text"          comment:"Constant name"`
	Value  string           `bun:"value,type:text"            comment:"Constant value"`
}

func (Constant) TableName() string {
	return "constant"
}
