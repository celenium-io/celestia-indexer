// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/goccy/go-json"
)

type Constants struct {
	Module        map[string]Params `json:"module"`
	DenomMetadata []DenomMetadata   `json:"denom_metadata"`
}

type Params map[string]string

type DenomMetadata struct {
	Description string `example:"Some description"    json:"description" swaggertype:"string"`
	Base        string `example:"utia"                json:"base"        swaggertype:"string"`
	Display     string `example:"TIA"                 json:"display"     swaggertype:"string"`
	Name        string `example:"TIA"                 json:"name"        swaggertype:"string"`
	Symbol      string `example:"TIA"                 json:"symbol"      swaggertype:"string"`
	Uri         string `example:"https://example.com" json:"uri"         swaggertype:"string"`

	Units json.RawMessage `json:"units"`
}

func NewConstants(consts []storage.Constant, denomMetadata []storage.DenomMetadata) Constants {
	response := Constants{
		Module:        make(map[string]Params),
		DenomMetadata: make([]DenomMetadata, len(denomMetadata)),
	}

	for i := range consts {
		if params, ok := response.Module[string(consts[i].Module)]; ok {
			params[consts[i].Name] = consts[i].Value
		} else {
			response.Module[string(consts[i].Module)] = Params{
				consts[i].Name: consts[i].Value,
			}
		}
	}

	for i := range denomMetadata {
		response.DenomMetadata[i].Base = denomMetadata[i].Base
		response.DenomMetadata[i].Symbol = denomMetadata[i].Symbol
		response.DenomMetadata[i].Name = denomMetadata[i].Name
		response.DenomMetadata[i].Description = denomMetadata[i].Description
		response.DenomMetadata[i].Display = denomMetadata[i].Display
		response.DenomMetadata[i].Uri = denomMetadata[i].Uri
		response.DenomMetadata[i].Units = denomMetadata[i].Units
	}

	return response
}
