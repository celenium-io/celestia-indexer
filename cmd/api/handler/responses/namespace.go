// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
)

type Namespace struct {
	ID          uint64 `example:"321"                                                      format:"integer" json:"id"           swaggertype:"integer"`
	Size        int64  `example:"12345"                                                    format:"integer" json:"size"         swaggertype:"integer"`
	Version     byte   `examle:"1"                                                         format:"byte"    json:"version"      swaggertype:"integer"`
	NamespaceID string `example:"4723ce10b187716adfc55ff7e6d9179c226e6b5440b02577cca49d02" format:"binary"  json:"namespace_id" swaggertype:"string"`
	Hash        string `example:"U3dhZ2dlciByb2Nrcw=="                                     format:"base64"  json:"hash"         swaggertype:"string"`
	Reserved    bool   `example:"true"                                                     json:"reserved"`
	PfbCount    int64  `example:"12"                                                       format:"integer" json:"pfb_count"    swaggertype:"integer"`
}

func NewNamespace(ns storage.Namespace) Namespace {
	return Namespace{
		ID:          ns.Id,
		Size:        ns.Size,
		Version:     ns.Version,
		NamespaceID: hex.EncodeToString(ns.NamespaceID),
		Hash:        ns.Hash(),
		Reserved:    ns.Reserved,
		PfbCount:    ns.PfbCount,
	}
}

func (Namespace) SearchType() string {
	return "namespace"
}

type ActiveNamespace struct {
	ID          uint64         `example:"321"                                                      format:"integer"   json:"id"           swaggertype:"integer"`
	Size        int64          `example:"12345"                                                    format:"integer"   json:"size"         swaggertype:"integer"`
	Version     byte           `examle:"1"                                                         format:"byte"      json:"version"      swaggertype:"integer"`
	NamespaceID string         `example:"4723ce10b187716adfc55ff7e6d9179c226e6b5440b02577cca49d02" format:"binary"    json:"namespace_id" swaggertype:"string"`
	Hash        string         `example:"U3dhZ2dlciByb2Nrcw=="                                     format:"base64"    json:"hash"         swaggertype:"string"`
	Reserved    bool           `example:"true"                                                     json:"reserved"`
	PfbCount    int64          `example:"12"                                                       format:"integer"   json:"pfb_count"    swaggertype:"integer"`
	Height      pkgTypes.Level `example:"100"                                                      format:"int64"     json:"height"       swaggertype:"integer"`
	Time        time.Time      `example:"2023-07-04T03:10:57+00:00"                                format:"date-time" json:"time"         swaggertype:"string"`
}

func NewActiveNamespace(ns storage.ActiveNamespace) ActiveNamespace {
	return ActiveNamespace{
		ID:          ns.Id,
		Size:        ns.Size,
		Version:     ns.Version,
		NamespaceID: hex.EncodeToString(ns.NamespaceID),
		Hash:        ns.Hash(),
		Reserved:    ns.Reserved,
		PfbCount:    ns.PfbCount,
		Height:      ns.Height,
		Time:        ns.Time,
	}
}
