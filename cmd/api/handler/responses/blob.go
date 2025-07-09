// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/base64"
	"encoding/hex"
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
	Id          uint64         `example:"200"                                          format:"integer"   json:"id"           swaggertype:"integer"`
	Commitment  string         `example:"vbGakK59+Non81TE3ULg5Ve5ufT9SFm/bCyY+WLR3gg=" format:"base64"    json:"commitment"   swaggertype:"string"`
	Size        int64          `example:"10"                                           format:"integer"   json:"size"         swaggertype:"integer"`
	Height      pkgTypes.Level `example:"100"                                          format:"integer"   json:"height"       swaggertype:"integer"`
	Time        time.Time      `example:"2023-07-04T03:10:57+00:00"                    format:"date-time" json:"time"         swaggertype:"string"`
	ContentType string         `example:"image/png"                                    format:"string"    json:"content_type" swaggertype:"string"`
	Namespace   *Namespace     `json:"namespace,omitempty"`
	Tx          *Tx            `json:"tx,omitempty"`
	Rollup      *ShortRollup   `json:"rollup,omitempty"`
	Signer      *ShortAddress  `json:"signer,omitempty"`
}

func NewBlobLog(blob storage.BlobLog) BlobLog {
	b := BlobLog{
		Id:          blob.Id,
		Commitment:  blob.Commitment,
		Size:        blob.Size,
		Height:      blob.Height,
		Time:        blob.Time,
		ContentType: blob.ContentType,
		Rollup:      NewShortRollup(blob.Rollup),
		Signer:      NewShortAddress(blob.Signer),
	}

	if blob.Namespace != nil {
		ns := NewNamespace(*blob.Namespace)
		b.Namespace = &ns
	}
	if blob.Tx != nil {
		tx := NewTx(*blob.Tx)
		b.Tx = &tx
	}

	return b
}

type LightBlobLog struct {
	Id          uint64         `example:"200"                                                              format:"integer"   json:"id"           swaggertype:"integer"`
	Commitment  string         `example:"vbGakK59+Non81TE3ULg5Ve5ufT9SFm/bCyY+WLR3gg="                     format:"base64"    json:"commitment"   swaggertype:"string"`
	Size        int64          `example:"10"                                                               format:"integer"   json:"size"         swaggertype:"integer"`
	Height      pkgTypes.Level `example:"100"                                                              format:"integer"   json:"height"       swaggertype:"integer"`
	Time        time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"         swaggertype:"string"`
	ContentType string         `example:"image/png"                                                        format:"string"    json:"content_type" swaggertype:"string"`
	Namespace   string         `example:"AAAAAAAAAAAAAAAAAAAAAAAAAAAAs2bWWU6FOB0="                         format:"base64"    json:"namespace"    swaggertype:"string"`
	TxHash      string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"tx_hash"      swaggertype:"string"`

	Signer *ShortAddress `json:"signer,omitempty"`
}

func NewLightBlobLog(blob storage.BlobLog) LightBlobLog {
	b := LightBlobLog{
		Id:          blob.Id,
		Commitment:  blob.Commitment,
		Size:        blob.Size,
		Height:      blob.Height,
		Time:        blob.Time,
		ContentType: blob.ContentType,
		Signer:      NewShortAddress(blob.Signer),
	}

	if blob.Namespace != nil {
		b.Namespace = blob.Namespace.Hash()
	}
	if blob.Tx != nil {
		b.TxHash = hex.EncodeToString(blob.Tx.Hash)
	}

	return b
}
