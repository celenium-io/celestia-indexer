package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type INamespace interface {
	storage.Table[*Namespace]

	ByNamespaceId(ctx context.Context, namespaceId []byte) ([]Namespace, error)
	ByNamespaceIdAndVersion(ctx context.Context, namespaceId []byte, version byte) (Namespace, error)
}

// Namespace -
type Namespace struct {
	bun.BaseModel `bun:"namespace" comment:"Table with celestia namespaces."`

	ID          uint64 `bun:"id,type:bigint,pk,notnull" comment:"Unique internal identity"`
	Version     byte   `bun:"version"                   comment:"Namespace version"`
	NamespaceID []byte `bun:"namespace_id"              comment:"Namespace identity"`
	Size        uint64 `bun:"size"                      comment:"Namespace size"`
	Reserved    bool   `bun:"reserved"                  comment:"If namespace is reserved flag is true"`
}

// TableName -
func (Namespace) TableName() string {
	return "namespace"
}
