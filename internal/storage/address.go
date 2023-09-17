package storage

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/pkg/types"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IAddress interface {
	storage.Table[*Address]

	ByHash(ctx context.Context, hash []byte) (Address, error)
}

// Address -
type Address struct {
	bun.BaseModel `bun:"address" comment:"Table with celestia addresses."`

	Id         uint64      `bun:"id,pk,notnull,autoincrement" comment:"Unique internal identity"`
	Height     types.Level `bun:"height"                      comment:"Block number of the first address occurrence."`
	LastHeight types.Level `bun:"last_height"                 comment:"Block number of the last address occurence "`
	Hash       []byte      `bun:"hash"                        comment:"Address hash."`
	Address    string      `bun:"address,unique:address_idx"  comment:"Human-readable ddress."`

	Balance Balance `bun:"rel:has-one,join:id=id"`
}

// TableName -
func (Address) TableName() string {
	return "address"
}

func (address Address) String() string {
	return address.Address
}
