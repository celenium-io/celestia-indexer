package storage

import (
	"context"

	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IDenomMetadata interface {
	All(ctx context.Context) ([]DenomMetadata, error)
}

type DenomMetadata struct {
	bun.BaseModel `bun:"table:denom_metadata" comment:"Table with celestia coins metadata."`

	Id          uint64 `bun:"id,pk,notnull,autoincrement" comment:"Internal unique identity"`
	Description string `bun:"description,type:text"       comment:"Denom description"`
	Base        string `bun:"base,type:text"              comment:"Denom base"`
	Display     string `bun:"display,type:text"           comment:"Denom display"`
	Name        string `bun:"name,type:text"              comment:"Denom name"`
	Symbol      string `bun:"symbol,type:text"            comment:"Denom symbol"`
	Uri         string `bun:"uri,type:text"               comment:"Denom uri"`

	Units []byte `bun:"units,type:bytea" comment:"Denom units information"`
}

func (DenomMetadata) TableName() string {
	return "denom_metadata"
}
