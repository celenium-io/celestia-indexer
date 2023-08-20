package storage

import (
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

// IAddress -
type INamespace interface {
	storage.Table[*Namespace]
}

// Namespace -
type Namespace struct {
	bun.BaseModel `bun:"namespace" comment:"Table with celestia namespaces."`

	ID          uint64 `bun:"id,type:bigint,pk,notnull" comment:"Unique internal identity"`
	Version     byte   `bun:"version"                   comment:"Namespace version"`
	NamespaceID []byte `bun:"namespace_id"              comment:"Namespace identity"`
	Size        uint64 `bun:"size"                      comment:"Namespace size"`
}

// TableName -
func (Namespace) TableName() string {
	return "namespace"
}
