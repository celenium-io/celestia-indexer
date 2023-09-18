package postgres

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Namespace -
type Namespace struct {
	*postgres.Table[*storage.Namespace]
}

// NewNamespace -
func NewNamespace(db *database.Bun) *Namespace {
	return &Namespace{
		Table: postgres.NewTable[*storage.Namespace](db),
	}
}

// ByNamespaceId -
func (n *Namespace) ByNamespaceId(ctx context.Context, namespaceId []byte) (namespace []storage.Namespace, err error) {
	err = n.DB().NewSelect().Model(&namespace).
		Where("namespace_id = ?", namespaceId).
		Scan(ctx)
	return
}

// ByNamespaceIdAndVersion -
func (n *Namespace) ByNamespaceIdAndVersion(ctx context.Context, namespaceId []byte, version byte) (namespace storage.Namespace, err error) {
	err = n.DB().NewSelect().Model(&namespace).
		Where("namespace_id = ?", namespaceId).
		Where("version = ?", version).
		Scan(ctx)
	return
}

// Messages -
func (n *Namespace) Messages(ctx context.Context, id uint64, limit, offset int) (msgs []storage.NamespaceMessage, err error) {
	query := n.DB().NewSelect().Model(&msgs).
		Where("namespace_message.namespace_id = ?", id).
		Order("namespace_message.time desc").
		Relation("Namespace").
		Relation("Message").
		Relation("Tx")
	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}
	err = query.Scan(ctx)
	return
}

// MessagesByHeight -
func (n *Namespace) MessagesByHeight(ctx context.Context, height uint64, limit, offset int) (msgs []storage.NamespaceMessage, err error) {
	query := n.DB().NewSelect().Model(&msgs).
		Where("namespace_message.height = ?", height).
		Order("namespace_message.time desc").
		Relation("Namespace").
		Relation("Message").
		Relation("Tx")
	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}
	err = query.Scan(ctx)
	return
}

func (n *Namespace) Active(ctx context.Context, top int) (ns []storage.ActiveNamespace, err error) {
	subQuery := n.DB().NewSelect().
		ColumnExpr("namespace_id, max(msg_id) as msg_id, max(height) as height, max(time) as time").
		Model((*storage.NamespaceMessage)(nil)).
		Group("namespace_id").
		Order("msg_id desc")
	subQuery = limitScope(subQuery, top)

	err = n.DB().NewSelect().
		ColumnExpr("action.time as time, action.height as height, namespace.*").
		TableExpr("(?) as action", subQuery).
		Join("LEFT JOIN namespace").
		JoinOn("namespace.id = action.namespace_id").
		Order("msg_id desc").
		Scan(ctx, &ns)
	return
}
