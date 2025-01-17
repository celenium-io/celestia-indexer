// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IApiKey interface {
	Get(ctx context.Context, key string) (ApiKey, error)
}

type ApiKey struct {
	bun.BaseModel `bun:"apikey" comment:"Table with private api keys"`

	Key         string `bun:"key,pk,notnull" comment:"Key"`
	Description string `bun:"description"    comment:"Additional info about issuer and user"`
	Admin       bool   `bun:"admin"          comment:"Verified user"`
}

func (ApiKey) TableName() string {
	return "apikey"
}
