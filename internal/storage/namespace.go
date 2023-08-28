package storage

import (
	"context"
	"fmt"

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

	ID          uint64 `bun:"id,pk,autoincrement"                          comment:"Unique internal identity"`
	Version     byte   `bun:"version,unique:namespace_id_version_idx"      comment:"Namespace version"`
	NamespaceID []byte `bun:"namespace_id,unique:namespace_id_version_idx" comment:"Namespace identity"`
	Size        uint64 `bun:"size"                                         comment:"Namespace size"`
	PfdCount    uint64 `bun:"pfd_count"                                    comment:"Count of pay for blobs messages for the namespace"`
	Reserved    bool   `bun:"reserved"                                     comment:"If namespace is reserved flag is true"`
}

// TableName -
func (Namespace) TableName() string {
	return "namespace"
}

func (ns Namespace) String() string {
	return fmt.Sprintf("%x%x", ns.Version, ns.NamespaceID)
}
