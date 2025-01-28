// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type ICelestial interface {
	ById(ctx context.Context, id string) (Celestial, error)
	ByAddressId(ctx context.Context, addressId uint64, limit, offset int) ([]Celestial, error)
}

type Celestial struct {
	bun.BaseModel `bun:"celestial" comment:"Table with celestial ids."`

	Id        string `bun:"id,pk,notnull" comment:"Celestial id"`
	AddressId uint64 `bun:"address_id"    comment:"Internal address identity for connected address"`
	ImageUrl  string `bun:"image_url"     comment:"Image url"`
	ChangeId  int64  `bun:"change_id"     comment:"Id of the last change of celestial id"`
}

func (Celestial) TableName() string {
	return "celestial"
}

func (cid Celestial) String() string {
	return fmt.Sprintf("%s %s", cid.Id, cid.ImageUrl)
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type ICelestialState interface {
	ByName(ctx context.Context, name string) (CelestialState, error)
	Save(ctx context.Context, state *CelestialState) error
}

type CelestialState struct {
	bun.BaseModel `bun:"celestial_state" comment:"Table with celestial ids."`

	Name     string `bun:"name,pk,notnull" comment:"Celestial id indexer name"`
	ChangeId int64  `bun:"change_id"       comment:"Id of the last change of celestial id"`
}

func (CelestialState) TableName() string {
	return "celestial_state"
}

func (cid CelestialState) String() string {
	return fmt.Sprintf("%s %d", cid.Name, cid.ChangeId)
}
