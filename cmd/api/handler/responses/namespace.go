package responses

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
)

type Namespace struct {
	ID          uint64 `example:"321"                                                      format:"integer" json:"id"           swaggertype:"integer"`
	Size        uint64 `example:"12345"                                                    format:"integer" json:"size"         swaggertype:"integer"`
	Version     byte   `examle:"1"                                                         format:"byte"    json:"version"      swaggertype:"integer"`
	NamespaceID string `example:"4723ce10b187716adfc55ff7e6d9179c226e6b5440b02577cca49d02" format:"binary"  json:"namespace_id" swaggertype:"string"`
	Hash        string `example:"U3dhZ2dlciByb2Nrcw=="                                     format:"base64"  json:"hash"         swaggertype:"string"`
	Reserved    bool   `example:"true"                                                     json:"reserved"`
}

func NewNamespace(ns storage.Namespace) Namespace {
	return Namespace{
		ID:          ns.Id,
		Size:        ns.Size,
		Version:     ns.Version,
		NamespaceID: hex.EncodeToString(ns.NamespaceID),
		Hash:        base64.URLEncoding.EncodeToString(append([]byte{ns.Version}, ns.NamespaceID...)),
		Reserved:    ns.Reserved,
	}
}

func (Namespace) SearchType() string {
	return "namespace"
}
