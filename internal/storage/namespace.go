package storage

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/dipdup-io/celestia-indexer/pkg/types"
	pkgTypes "github.com/dipdup-io/celestia-indexer/pkg/types"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type INamespace interface {
	storage.Table[*Namespace]

	ByNamespaceId(ctx context.Context, namespaceId []byte) ([]Namespace, error)
	ByNamespaceIdAndVersion(ctx context.Context, namespaceId []byte, version byte) (Namespace, error)
	Messages(ctx context.Context, id uint64, limit, offset int) ([]NamespaceMessage, error)
	MessagesByHeight(ctx context.Context, height uint64, limit, offset int) ([]NamespaceMessage, error)
	Active(ctx context.Context, top int) ([]ActiveNamespace, error)
}

// Namespace -
type Namespace struct {
	bun.BaseModel `bun:"namespace" comment:"Table with celestia namespaces."`

	Id          uint64      `bun:"id,pk,autoincrement"                          comment:"Unique internal identity"`
	FirstHeight types.Level `bun:"first_height,notnull"                         comment:"Block height of the first message changing the namespace"`
	Version     byte        `bun:"version,unique:namespace_id_version_idx"      comment:"Namespace version"`
	NamespaceID []byte      `bun:"namespace_id,unique:namespace_id_version_idx" comment:"Namespace identity"`
	Size        uint64      `bun:"size"                                         comment:"Blobs size"`
	PfbCount    uint64      `bun:"pfb_count"                                    comment:"Count of pay for blobs messages for the namespace"`
	Reserved    bool        `bun:"reserved,default:false"                       comment:"If namespace is reserved flag is true"`
}

// TableName -
func (Namespace) TableName() string {
	return "namespace"
}

func (ns Namespace) String() string {
	return fmt.Sprintf("%x%x", ns.Version, ns.NamespaceID)
}

func (ns Namespace) Hash() string {
	return base64.StdEncoding.EncodeToString(append([]byte{ns.Version}, ns.NamespaceID...))
}

type ActiveNamespace struct {
	Namespace
	Height pkgTypes.Level `bun:"height"`
	Time   time.Time      `bun:"time"`
}
