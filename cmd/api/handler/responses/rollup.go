// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/shopspring/decimal"
)

type RollupWithStats struct {
	Id             uint64 `example:"321"                                       format:"integer" json:"id"                    swaggertype:"integer"`
	Name           string `example:"Rollup name"                               format:"string"  json:"name"                  swaggertype:"string"`
	Description    string `example:"Long rollup description"                   format:"string"  json:"description,omitempty" swaggertype:"string"`
	Website        string `example:"https://website.com"                       format:"string"  json:"website,omitempty"     swaggertype:"string"`
	Twitter        string `example:"https://x.com/account"                     format:"string"  json:"twitter,omitempty"     swaggertype:"string"`
	Github         string `example:"https://github.com/account"                format:"string"  json:"github,omitempty"      swaggertype:"string"`
	Logo           string `example:"https://some_link.com/image.png"           format:"string"  json:"logo,omitempty"        swaggertype:"string"`
	Slug           string `example:"rollup_slug"                               format:"string"  json:"slug"                  swaggertype:"string"`
	L2Beat         string `example:"https://l2beat.com/scaling/projects/karak" format:"string"  json:"l2_beat,omitempty"     swaggertype:"string"`
	DeFiLama       string `example:"Manta"                                     format:"string"  json:"defi_lama,omitempty"   swaggertype:"string"`
	Explorer       string `example:"https://explorer.karak.network/"           format:"string"  json:"explorer,omitempty"    swaggertype:"string"`
	BridgeContract string `example:"https://github.com/account"                format:"string"  json:"bridge,omitempty"      swaggertype:"string"`
	Stack          string `example:"op_stack"                                  format:"string"  json:"stack,omitempty"       swaggertype:"string"`
	Type           string `example:"settled"                                   format:"string"  json:"type,omitempty"        swaggertype:"string"`
	Category       string `example:"nft"                                       format:"string"  json:"category,omitempty"    swaggertype:"string"`
	VM             string `example:"evm"                                       format:"string"  json:"vm,omitempty"          swaggertype:"string"`
	Provider       string `example:"name"                                      format:"string"  json:"provider,omitempty"    swaggertype:"string"`
	SettledOn      string `example:"Ethereum"                                  format:"string"  json:"settled_on,omitempty"  swaggertype:"string"`
	Color          string `example:"#123456"                                   format:"string"  json:"color,omitempty"       swaggertype:"string"`
	Compression    string `example:"zip"                                       format:"string"  json:"compression,omitempty" swaggertype:"string"`

	BlobsCount    int64     `example:"2"                         format:"integer"   json:"blobs_count"        swaggertype:"integer"`
	Size          int64     `example:"1000"                      format:"integer"   json:"size"               swaggertype:"integer"`
	LastAction    time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"last_message_time"  swaggertype:"string"`
	FirstAction   time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"first_message_time" swaggertype:"string"`
	Fee           string    `example:"123.456789"                format:"string"    json:"fee"                swaggertype:"string"`
	SizePct       float64   `example:"0.9876"                    format:"float"     json:"size_pct"           swaggertype:"number"`
	FeePct        float64   `example:"0.9876"                    format:"float"     json:"fee_pct"            swaggertype:"number"`
	BlobsCountPct float64   `example:"0.9876"                    format:"float"     json:"blobs_count_pct"    swaggertype:"number"`
	DAPct         float64   `example:"0.9876"                    format:"float"     json:"da_pct"             swaggertype:"number"`
	IsActive      bool      `example:"true"                      format:"boolean"   json:"is_active"          swaggertype:"boolean"`

	Tags  []string `json:"tags,omitempty"`
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
		DeFiLama:       r.DeFiLama,
		Explorer:       r.Explorer,
		BridgeContract: r.BridgeContract,
		Links:          r.Links,
		Stack:          r.Stack,
		Slug:           r.Slug,
		BlobsCount:     r.BlobsCount,
		Size:           r.Size,
		SizePct:        r.SizePct,
		BlobsCountPct:  r.BlobsCountPct,
		DAPct:          r.DAPct,
		FeePct:         r.FeePct,
		LastAction:     r.LastActionTime,
		FirstAction:    r.FirstActionTime,
		Compression:    r.Compression,
		SettledOn:      r.SettledOn,
		Category:       r.Category.String(),
		Type:           r.Type.String(),
		Provider:       r.Provider,
		VM:             r.VM,
		Fee:            r.Fee.StringFixed(0),
		Tags:           r.Tags,
		Color:          r.Color,
		IsActive:       r.IsActive,
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
	DeFiLama       string `example:"Manta"                           format:"string"  json:"defi_lama,omitempty"   swaggertype:"string"`
	Explorer       string `example:"https://explorer.karak.network/" format:"string"  json:"explorer,omitempty"    swaggertype:"string"`
	BridgeContract string `example:"https://github.com/account"      format:"string"  json:"bridge,omitempty"      swaggertype:"string"`
	Stack          string `example:"op_stack"                        format:"string"  json:"stack,omitempty"       swaggertype:"string"`
	Type           string `example:"settled"                         format:"string"  json:"type,omitempty"        swaggertype:"string"`
	Category       string `example:"nft"                             format:"string"  json:"category,omitempty"    swaggertype:"string"`
	Provider       string `example:"name"                            format:"string"  json:"provider,omitempty"    swaggertype:"string"`
	Compression    string `example:"zip"                             format:"string"  json:"compression,omitempty" swaggertype:"string"`
	VM             string `example:"evm"                             format:"string"  json:"vm,omitempty"          swaggertype:"string"`
	SettledOn      string `example:"Ethereum"                        format:"string"  json:"settled_on,omitempty"  swaggertype:"string"`
	Color          string `example:"#123456"                         format:"string"  json:"color,omitempty"       swaggertype:"string"`

	Tags  []string `json:"tags,omitempty"`
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
		DeFiLama:       r.DeFiLama,
		BridgeContract: r.BridgeContract,
		Stack:          r.Stack,
		Explorer:       r.Explorer,
		Links:          r.Links,
		Compression:    r.Compression,
		Category:       r.Category.String(),
		Type:           r.Type.String(),
		Provider:       r.Provider,
		VM:             r.VM,
		SettledOn:      r.SettledOn,
		Tags:           r.Tags,
		Color:          r.Color,
	}
}

type ShortRollup struct {
	Id   uint64 `example:"321"                             format:"integer" json:"id"             swaggertype:"integer"`
	Name string `example:"Rollup name"                     format:"string"  json:"name"           swaggertype:"string"`
	Logo string `example:"https://some_link.com/image.png" format:"string"  json:"logo,omitempty" swaggertype:"string"`
	Slug string `example:"rollup_slug"                     format:"string"  json:"slug"           swaggertype:"string"`
}

func NewShortRollup(r *storage.Rollup) *ShortRollup {
	if r == nil {
		return nil
	}
	return &ShortRollup{
		Id:   r.Id,
		Name: r.Name,
		Logo: r.Logo,
		Slug: r.Slug,
	}
}

type RollupWithDayStats struct {
	Id             uint64 `example:"321"                                       format:"integer" json:"id"                    swaggertype:"integer"`
	Name           string `example:"Rollup name"                               format:"string"  json:"name"                  swaggertype:"string"`
	Description    string `example:"Long rollup description"                   format:"string"  json:"description,omitempty" swaggertype:"string"`
	Website        string `example:"https://website.com"                       format:"string"  json:"website,omitempty"     swaggertype:"string"`
	Twitter        string `example:"https://x.com/account"                     format:"string"  json:"twitter,omitempty"     swaggertype:"string"`
	Github         string `example:"https://github.com/account"                format:"string"  json:"github,omitempty"      swaggertype:"string"`
	Logo           string `example:"https://some_link.com/image.png"           format:"string"  json:"logo,omitempty"        swaggertype:"string"`
	Slug           string `example:"rollup_slug"                               format:"string"  json:"slug"                  swaggertype:"string"`
	L2Beat         string `example:"https://l2beat.com/scaling/projects/karak" format:"string"  json:"l2_beat,omitempty"     swaggertype:"string"`
	DeFiLama       string `example:"Manta"                                     format:"string"  json:"defi_lama,omitempty"   swaggertype:"string"`
	Explorer       string `example:"https://explorer.karak.network/"           format:"string"  json:"explorer,omitempty"    swaggertype:"string"`
	BridgeContract string `example:"https://github.com/account"                format:"string"  json:"bridge,omitempty"      swaggertype:"string"`
	Stack          string `example:"op_stack"                                  format:"string"  json:"stack,omitempty"       swaggertype:"string"`
	Type           string `example:"settled"                                   format:"string"  json:"type,omitempty"        swaggertype:"string"`
	Category       string `example:"nft"                                       format:"string"  json:"category,omitempty"    swaggertype:"string"`
	Provider       string `example:"name"                                      format:"string"  json:"provider,omitempty"    swaggertype:"string"`
	Compression    string `example:"zip"                                       format:"string"  json:"compression,omitempty" swaggertype:"string"`
	VM             string `example:"evm"                                       format:"string"  json:"vm,omitempty"          swaggertype:"string"`
	SettledOn      string `example:"Ethereum"                                  format:"string"  json:"settled_on,omitempty"  swaggertype:"string"`
	Color          string `example:"#123456"                                   format:"string"  json:"color,omitempty"       swaggertype:"string"`

	AvgSize        int64   `example:"100" format:"integer" json:"avg_size"        swaggertype:"integer"`
	BlobsCount     int64   `example:"100" format:"integer" json:"blobs_count"     swaggertype:"integer"`
	TotalSize      int64   `example:"100" format:"integer" json:"total_size"      swaggertype:"integer"`
	Throghput      int64   `example:"100" format:"integer" json:"throughput"      swaggertype:"integer"`
	NamespaceCount int64   `example:"100" format:"integer" json:"namespace_count" swaggertype:"integer"`
	PfbCount       int64   `example:"100" format:"integer" json:"pfb_count"       swaggertype:"integer"`
	TotalFee       string  `example:"100" format:"string"  json:"total_fee"       swaggertype:"string"`
	MBPrice        string  `example:"100" format:"string"  json:"mb_price"        swaggertype:"string"`
	FeePerPfb      string  `example:"100" format:"string"  json:"fee_per_pfb"     swaggertype:"string"`
	BlobsPerPfb    float64 `example:"100" format:"float"   json:"blobs_per_pfb"   swaggertype:"number"`
}

func NewRollupWithDayStats(r storage.RollupWithDayStats) RollupWithDayStats {
	response := RollupWithDayStats{
		Id:             r.Id,
		Name:           r.Name,
		Description:    r.Description,
		Github:         r.GitHub,
		Twitter:        r.Twitter,
		Website:        r.Website,
		Logo:           r.Logo,
		L2Beat:         r.L2Beat,
		DeFiLama:       r.DeFiLama,
		Explorer:       r.Explorer,
		BridgeContract: r.BridgeContract,
		Stack:          r.Stack,
		Compression:    r.Compression,
		Category:       r.Category.String(),
		Type:           r.Type.String(),
		Provider:       r.Provider,
		VM:             r.VM,
		SettledOn:      r.SettledOn,
		Slug:           r.Slug,
		BlobsCount:     r.BlobsCount,
		AvgSize:        int64(r.AvgSize),
		TotalSize:      r.TotalSize,
		Throghput:      r.Throghput,
		NamespaceCount: r.NamespaceCount,
		PfbCount:       r.PfbCount,
		TotalFee:       r.TotalFee.String(),
		MBPrice:        r.MBPrice.String(),
		FeePerPfb:      decimal.Zero.String(),
		Color:          r.Color,
	}

	if r.PfbCount > 0 {
		response.BlobsPerPfb = float64(r.BlobsCount / r.PfbCount)
		response.FeePerPfb = r.TotalFee.Div(decimal.NewFromInt(r.PfbCount)).String()
	}

	return response
}

type RollupGroupedStats struct {
	Fee        float64 `example:"123.456789" format:"string"  json:"fee"         swaggertype:"string"`
	Size       float64 `example:"1000"       format:"integer" json:"size"        swaggertype:"integer"`
	BlobsCount int64   `example:"2"          format:"integer" json:"blobs_count" swaggertype:"integer"`
	Group      string  `example:"group"      format:"string"  json:"group"       swaggertype:"string"`
}

func NewRollupGroupedStats(r storage.RollupGroupedStats) RollupGroupedStats {
	response := RollupGroupedStats{
		Fee:        r.Fee,
		Size:       r.Size,
		BlobsCount: r.BlobsCount,
		Group:      r.Group,
	}

	return response
}
