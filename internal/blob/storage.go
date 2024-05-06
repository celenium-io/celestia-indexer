// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package blob

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/gabriel-vasile/mimetype"
	blobTypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

type Blob struct {
	*blobTypes.Blob
	Commitment []byte
	Height     uint64
}

func (blob Blob) String() string {
	ns := base64.URLEncoding.EncodeToString(blob.NamespaceId)
	cm := base64.URLEncoding.EncodeToString(blob.Commitment)
	return fmt.Sprintf("%s/%d/%s", ns, blob.Height, cm)
}

func (blob Blob) ContentType() string {
	contentType := mimetype.Detect(blob.Data)
	return contentType.String()
}

//go:generate mockgen -source=$GOFILE -destination=mock.go -package=blob -typed
type Storage interface {
	Save(ctx context.Context, blob Blob) error
	SaveBulk(ctx context.Context, blobs []Blob) error
	Head(ctx context.Context) (uint64, error)
	UpdateHead(ctx context.Context, head uint64) error
}
