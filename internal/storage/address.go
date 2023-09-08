package storage

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/pkg/types"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IAddress interface {
	storage.Table[*Address]

	ByHash(ctx context.Context, hash string) (Address, error)
}

// Address -
type Address struct {
	bun.BaseModel `bun:"address" comment:"Table with celestia addresses."`

	Id      uint64          `bun:"id,pk,notnull,autoincrement" comment:"Unique internal identity"`
	Height  types.Level     `bun:"height"                      comment:"Block number of the first address occurrence."`
	Hash    string          `bun:",unique:address_hash"        comment:"Address hash."`
	Balance decimal.Decimal `bun:",type:numeric"               comment:"Address balance"`
}

// TableName -
func (Address) TableName() string {
	return "address"
}

func (address Address) String() string {
	return address.Hash
}
