// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"bytes"
	"encoding/hex"
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"

	"github.com/celenium-io/celestia-indexer/internal/storage"
)

type Namespace struct {
	ID              uint64         `example:"321"                                                      format:"integer"   json:"id"                swaggertype:"integer"`
	Size            int64          `example:"12345"                                                    format:"integer"   json:"size"              swaggertype:"integer"`
	Version         byte           `examle:"1"                                                         format:"byte"      json:"version"           swaggertype:"integer"`
	NamespaceID     string         `example:"4723ce10b187716adfc55ff7e6d9179c226e6b5440b02577cca49d02" format:"binary"    json:"namespace_id"      swaggertype:"string"`
	Hash            string         `example:"U3dhZ2dlciByb2Nrcw=="                                     format:"base64"    json:"hash"              swaggertype:"string"`
	PfbCount        int64          `example:"12"                                                       format:"integer"   json:"pfb_count"         swaggertype:"integer"`
	LastHeight      pkgTypes.Level `example:"100"                                                      format:"int64"     json:"last_height"       swaggertype:"integer"`
	LastMessageTime time.Time      `example:"2023-07-04T03:10:57+00:00"                                format:"date-time" json:"last_message_time" swaggertype:"string"`
	Name            string         `example:"name"                                                     format:"string"    json:"name"              swaggertype:"string"`
	Reserved        bool           `example:"true"                                                     json:"reserved"`
}

func NewNamespace(ns storage.Namespace) Namespace {
	return Namespace{
		ID:              ns.Id,
		Size:            ns.Size,
		Version:         ns.Version,
		NamespaceID:     hex.EncodeToString(ns.NamespaceID),
		Name:            decodeName(ns.NamespaceID),
		Hash:            ns.Hash(),
		Reserved:        ns.Reserved,
		PfbCount:        ns.PfbCount,
		LastHeight:      ns.LastHeight,
		LastMessageTime: ns.LastMessageTime,
	}
}

func (Namespace) SearchType() string {
	return "namespace"
}

func decodeName(nsId []byte) string {
	var (
		trimmed     = bytes.TrimLeft(nsId, "\x00")
		data        = make([]byte, 0)
		isDecodable = true
	)
	for i := range trimmed {
		if trimmed[i] < 0x20 || trimmed[i] > 0x7f {
			isDecodable = false
		}
		data = append(data, trimmed[i])
	}

	if isDecodable {
		return string(data)
	}
	return hex.EncodeToString(data)
}
