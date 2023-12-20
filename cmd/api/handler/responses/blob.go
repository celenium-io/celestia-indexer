// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/node/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
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

type BlobLog struct {
	Commitment string         `example:"vbGakK59+Non81TE3ULg5Ve5ufT9SFm/bCyY+WLR3gg="    format:"base64"    json:"commitment" swaggertype:"string"`
	Size       int64          `example:"10"                                              format:"integer"   json:"size"       swaggertype:"integer"`
	Height     pkgTypes.Level `example:"100"                                             format:"integer"   json:"height"     swaggertype:"integer"`
	Time       time.Time      `example:"2023-07-04T03:10:57+00:00"                       format:"date-time" json:"time"       swaggertype:"string"`
	Signer     string         `example:"celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60" format:"string"    json:"signer"     swaggertype:"string"`

	Namespace *Namespace `json:"namespace,omitempty"`
}

func NewBlobLog(blob storage.BlobLog) BlobLog {
	b := BlobLog{
		Commitment: blob.Commitment,
		Size:       blob.Size,
		Height:     blob.Height,
		Time:       blob.Time,
	}

	if blob.Namespace != nil {
		ns := NewNamespace(*blob.Namespace)
		b.Namespace = &ns
	}
	if blob.Signer != nil {
		b.Signer = blob.Signer.Address
	}

	return b
}
