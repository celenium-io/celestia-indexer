// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
)

type RollupWithStats struct {
	Id          uint64 `example:"321"                             format:"integer" json:"id"                    swaggertype:"integer"`
	Name        string `example:"Rollup name"                     format:"string"  json:"name"                  swaggertype:"string"`
	Description string `example:"Long rollup description"         format:"string"  json:"description,omitempty" swaggertype:"string"`
	Website     string `example:"https://website.com"             format:"string"  json:"website,omitempty"     swaggertype:"string"`
	Twitter     string `example:"https://x.com/account"           format:"string"  json:"twitter,omitempty"     swaggertype:"string"`
	Github      string `example:"https://github.com/account"      format:"string"  json:"github,omitempty"      swaggertype:"string"`
	Logo        string `example:"https://some_link.com/image.png" format:"string"  json:"logo,omitempty"        swaggertype:"string"`

	BlobsCount int64     `example:"2"                         format:"integer"   json:"blobs_count"       swaggertype:"integer"`
	Size       int64     `example:"1000"                      format:"integer"   json:"size"              swaggertype:"integer"`
	LastAction time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"last_message_time" swaggertype:"string"`
}

func NewRollupWithStats(r storage.RollupWithStats) RollupWithStats {
	return RollupWithStats{
		Id:          r.Id,
		Name:        r.Name,
		Description: r.Description,
		Github:      r.GitHub,
		Twitter:     r.Twitter,
		Website:     r.Website,
		Logo:        r.Logo,
		BlobsCount:  r.BlobsCount,
		Size:        r.Size,
		LastAction:  r.LastActionTime,
	}
}

type Rollup struct {
	Id          uint64 `example:"321"                             format:"integer" json:"id"                    swaggertype:"integer"`
	Name        string `example:"Rollup name"                     format:"string"  json:"name"                  swaggertype:"string"`
	Description string `example:"Long rollup description"         format:"string"  json:"description,omitempty" swaggertype:"string"`
	Website     string `example:"https://website.com"             format:"string"  json:"website,omitempty"     swaggertype:"string"`
	Twitter     string `example:"https://x.com/account"           format:"string"  json:"twitter,omitempty"     swaggertype:"string"`
	Github      string `example:"https://github.com/account"      format:"string"  json:"github,omitempty"      swaggertype:"string"`
	Logo        string `example:"https://some_link.com/image.png" format:"string"  json:"logo,omitempty"        swaggertype:"string"`
}

func NewRollup(r *storage.Rollup) Rollup {
	return Rollup{
		Id:          r.Id,
		Name:        r.Name,
		Description: r.Description,
		Github:      r.GitHub,
		Twitter:     r.Twitter,
		Website:     r.Website,
		Logo:        r.Logo,
	}
}
