// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rpc

import (
	"context"
	"strconv"

	"github.com/bytedance/sonic"

	"github.com/celenium-io/celestia-indexer/pkg/node/types"
	jxpkg "github.com/go-faster/jx"
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

		var gc GenesisChunk
		err := api.getStream(ctx, path, args, func(d *jxpkg.Decoder) error {
			return jxResponse(d, func(d *jxpkg.Decoder) error {
				return jxGenesisChunk(d, &gc)
			})
		})
		if err != nil {
			return types.Genesis{}, errors.Wrap(err, "genesis block request")
		}

		chunk++
		total = gc.Total
		genesisData = append(genesisData, gc.Data...)
	}

	var genesis types.Genesis
	err := sonic.ConfigFastest.Unmarshal(genesisData, &genesis)
	return genesis, err
}
