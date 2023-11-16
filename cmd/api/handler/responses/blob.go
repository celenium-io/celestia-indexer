// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/base64"
	"net/http"

	"github.com/celenium-io/celestia-indexer/pkg/node/types"
)

type Blob struct {
	Namespace    string `example:"AAAAAAAAAAAAAAAAAAAAAAAAAAAAs2bWWU6FOB0="     format:"base64"  json:"namespace"     swaggertype:"string"`
	Data         string `example:"b2sgZGVtbyBkYQ=="                             format:"base64"  json:"data"          swaggertype:"string"`
	ShareVersion int    `example:"0"                                            format:"integer" json:"share_version" swaggertype:"integer"`
	Commitment   string `example:"vbGakK59+Non81TE3ULg5Ve5ufT9SFm/bCyY+WLR3gg=" format:"base64"  json:"commitment"    swaggertype:"string"`
	ContentType  string `example:"image/png"                                    format:"string"  json:"content_type"  swaggertype:"string"`
}

func NewBlob(blob types.Blob) (Blob, error) {
	b := Blob{
		Namespace:    blob.Namespace,
		Data:         blob.Data,
		Commitment:   blob.Commitment,
		ShareVersion: blob.ShareVersion,
	}

	data, err := base64.StdEncoding.DecodeString(blob.Data)
	if err != nil {
		return b, err
	}
	b.ContentType = http.DetectContentType(data)
	return b, nil
}
