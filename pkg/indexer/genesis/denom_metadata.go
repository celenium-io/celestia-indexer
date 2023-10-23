// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package genesis

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/node/types"
)

func (module *Module) parseDenomMetadata(raw []types.DenomMetadata, data *parsedData) {
	for i := range raw {
		dm := storage.DenomMetadata{
			Description: raw[i].Description,
			Base:        raw[i].Base,
			Display:     raw[i].Display,
			Name:        raw[i].Name,
			Symbol:      raw[i].Symbol,
			Uri:         raw[i].URI,
			Units:       raw[i].DenomUnits,
		}
		data.denomMetadata = append(data.denomMetadata, dm)
	}
}
