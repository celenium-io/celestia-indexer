package storage

import (
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

// IAddress -
type IAddress interface {
	storage.Table[*Address]
}

// Address -
type Address struct {
	bun.BaseModel `bun:"address" comment:"Table with celestia addresses."`

	ID      uint64          `bun:"id,type:bigint,pk,notnull" comment:"Unique internal identity"`
	Height  uint64          `bun:"height"                    comment:"Block number of the first address occurrence."`
	Hash    []byte          `bun:",unique:address_hash"      comment:"Address hash."`
	Balance decimal.Decimal `bun:",type:numeric"             comment:"Address balance"`
}

// TableName -
func (Address) TableName() string {
	return "address"
}
