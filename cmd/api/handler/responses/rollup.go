// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
)

type RollupWithStats struct {
	Id             uint64 `example:"321"                             format:"integer" json:"id"                    swaggertype:"integer"`
	Name           string `example:"Rollup name"                     format:"string"  json:"name"                  swaggertype:"string"`
	Description    string `example:"Long rollup description"         format:"string"  json:"description,omitempty" swaggertype:"string"`
	Website        string `example:"https://website.com"             format:"string"  json:"website,omitempty"     swaggertype:"string"`
	Twitter        string `example:"https://x.com/account"           format:"string"  json:"twitter,omitempty"     swaggertype:"string"`
	Github         string `example:"https://github.com/account"      format:"string"  json:"github,omitempty"      swaggertype:"string"`
	Logo           string `example:"https://some_link.com/image.png" format:"string"  json:"logo,omitempty"        swaggertype:"string"`
	Slug           string `example:"rollup_slug"                     format:"string"  json:"slug"                  swaggertype:"string"`
	L2Beat         string `example:"https://github.com/account"      format:"string"  json:"l2_beat,omitempty"     swaggertype:"string"`
	BridgeContract string `example:"https://github.com/account"      format:"string"  json:"bridge,omitempty"      swaggertype:"string"`

	BlobsCount  int64     `example:"2"                         format:"integer"   json:"blobs_count"        swaggertype:"integer"`
	Size        int64     `example:"1000"                      format:"integer"   json:"size"               swaggertype:"integer"`
	LastAction  time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"last_message_time"  swaggertype:"string"`
	FirstAction time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"first_message_time" swaggertype:"string"`
	Fee         string    `example:"123.456789"                format:"string"    json:"fee"                swaggertype:"string"`

	Links []string `json:"links,omitempty"`
}

func NewRollupWithStats(r storage.RollupWithStats) RollupWithStats {
	return RollupWithStats{
		Id:             r.Id,
		Name:           r.Name,
		Description:    r.Description,
		Github:         r.GitHub,
		Twitter:        r.Twitter,
		Website:        r.Website,
		Logo:           r.Logo,
		L2Beat:         r.L2Beat,
		BridgeContract: r.BridgeContract,
		Links:          r.Links,
		Slug:           r.Slug,
		BlobsCount:     r.BlobsCount,
		Size:           r.Size,
		LastAction:     r.LastActionTime,
		FirstAction:    r.FirstActionTime,
		Fee:            r.Fee.StringFixed(0),
	}
}

type Rollup struct {
	Id             uint64 `example:"321"                             format:"integer" json:"id"                    swaggertype:"integer"`
	Name           string `example:"Rollup name"                     format:"string"  json:"name"                  swaggertype:"string"`
	Description    string `example:"Long rollup description"         format:"string"  json:"description,omitempty" swaggertype:"string"`
	Website        string `example:"https://website.com"             format:"string"  json:"website,omitempty"     swaggertype:"string"`
	Twitter        string `example:"https://x.com/account"           format:"string"  json:"twitter,omitempty"     swaggertype:"string"`
	Github         string `example:"https://github.com/account"      format:"string"  json:"github,omitempty"      swaggertype:"string"`
	Logo           string `example:"https://some_link.com/image.png" format:"string"  json:"logo,omitempty"        swaggertype:"string"`
	Slug           string `example:"rollup_slug"                     format:"string"  json:"slug"                  swaggertype:"string"`
	L2Beat         string `example:"https://github.com/account"      format:"string"  json:"l2_beat,omitempty"     swaggertype:"string"`
	BridgeContract string `example:"https://github.com/account"      format:"string"  json:"bridge,omitempty"      swaggertype:"string"`

	Links []string `json:"links,omitempty"`
}

func NewRollup(r *storage.Rollup) Rollup {
	return Rollup{
		Id:             r.Id,
		Name:           r.Name,
		Description:    r.Description,
		Github:         r.GitHub,
		Twitter:        r.Twitter,
		Website:        r.Website,
		Logo:           r.Logo,
		Slug:           r.Slug,
		L2Beat:         r.L2Beat,
		BridgeContract: r.BridgeContract,
		Links:          r.Links,
	}
}
