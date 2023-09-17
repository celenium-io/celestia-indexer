package storage

import (
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IBalance interface {
	storage.Table[*Balance]
}

type Balance struct {
	bun.BaseModel `bun:"balance" comment:"Table with account balances."`

	Id       uint64          `bun:"id,pk,notnull,autoincrement" comment:"Unique internal identity"`
	Currency string          `bun:"currency,pk,notnull"         comment:"Balance currency"`
	Total    decimal.Decimal `bun:"total,type:numeric"          comment:"Total account balance"`
}

func (Balance) TableName() string {
	return "balance"
}

func (b Balance) IsEmpty() bool {
	return b.Currency == "" && b.Total.IsZero()
}
