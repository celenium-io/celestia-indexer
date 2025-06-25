// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package blob

import (
	"context"
	"encoding/base64"
	"fmt"

	blobTypes "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/gabriel-vasile/mimetype"
)

type Blob struct {
	*blobTypes.Blob
	Commitment []byte
	Height     uint64
}

func (blob Blob) String() string {
	hash := []byte{byte(blob.ShareVersion)}
	ns := base64.URLEncoding.EncodeToString(append(hash, blob.NamespaceId...))
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

func Base64ToUrl(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
