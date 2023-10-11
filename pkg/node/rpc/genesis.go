// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package rpc

import (
	"context"
	"strconv"

	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
)

type GenesisChunk struct {
	Chunk int64  `json:"chunk,string"`
	Total int64  `json:"total,string"`
	Data  []byte `json:"data"`
}

func (api *API) Genesis(ctx context.Context) (types.Genesis, error) {
	path := "genesis_chunked"

	genesisData := make([]byte, 0)
	var chunk int64
	var total int64

	for chunk == 0 || chunk < total {
		args := map[string]string{
			"chunk": strconv.FormatInt(chunk, 10),
		}

		var gr types.Response[GenesisChunk]
		if err := api.get(ctx, path, args, &gr); err != nil {
			return types.Genesis{}, errors.Wrap(err, "genesis block request")
		}

		if gr.Error != nil {
			return types.Genesis{}, errors.Wrapf(types.ErrRequest, "request %d error: %s", gr.Id, gr.Error.Error())
		}

		chunk += 1
		total = gr.Result.Total
		genesisData = append(genesisData, gr.Result.Data...)
	}

	var genesis types.Genesis
	err := json.Unmarshal(genesisData, &genesis)
	return genesis, err
}
