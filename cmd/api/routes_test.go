// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
	"github.com/dipdup-net/go-lib/config"
	pg "github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/stretchr/testify/require"
)

func TestRoutes(t *testing.T) {
	var expectedRoutes = map[string]struct{}{
		"/v1/address/:hash/redelegations GET":                 {},
		"/v1/address/:hash/celestials GET":                    {},
		"/v1/validators/count GET":                            {},
		"/v1/validators/:id/blocks GET":                       {},
		"/v1/validators/:id/delegators GET":                   {},
		"/v1/validators/:id/messages GET":                     {},
		"/v1/rollup/:id/stats/:name/:timeframe GET":           {},
		"/v1/gas/estimate_for_pfb GET":                        {},
		"/v1/tx/count GET":                                    {},
		"/v1/namespace/:id/:version/rollups GET":              {},
		"/v1/namespace_by_hash/:hash/:height GET":             {},
		"/v1/validators GET":                                  {},
		"/v1/validators/:id GET":                              {},
		"/v1/stats/tps GET":                                   {},
		"/v1/stats/namespace/series/:id/:name/:timeframe GET": {},
		"/v1/stats/series/:name/:timeframe GET":               {},
		"/v1/stats/series/:name/:timeframe/cumulative GET":    {},
		"/v1/search GET":                                      {},
		"/v1/stats/staking/series/:id/:name/:timeframe GET":   {},
		"/v1/stats/ibc/series/:id/:name/:timeframe GET":       {},
		"/v1/stats/hyperlane/series/:id/:name/:timeframe GET": {},
		"/v1/stats/hyperlane/chains/:name/:timeframe GET":     {},
		"/v1/stats/hyperlane/chains GET":                      {},
		"/v1/rollup/:id GET":                                  {},
		"/v1/address/:hash GET":                               {},
		"/v1/address/:hash/txs GET":                           {},
		"/v1/address/:hash/delegations GET":                   {},
		"/v1/namespace/:id/:version GET":                      {},
		"/v1/tx GET":                                          {},
		"/v1/namespace/:id/:version/blobs GET":                {},
		"/v1/rollup/:id/namespaces GET":                       {},
		"/v1/block GET":                                       {},
		"/v1/tx/:hash/blobs/count GET":                        {},
		"/v1/rollup GET":                                      {},
		"/v1/rollup/day GET":                                  {},
		"/v1/address/:hash/vestings GET":                      {},
		"/v1/tx/:hash/events GET":                             {},
		"/v1/stats/changes_24h GET":                           {},
		"/v1/stats/ibc/chains GET":                            {},
		"/v1/rollup/count GET":                                {},
		"/v1/address/:hash/undelegations GET":                 {},
		"/v1/block/:height/messages GET":                      {},
		"/v1/namespace_by_hash/:hash GET":                     {},
		"/v1/vesting/:id/periods GET":                         {},
		"/v1/constants GET":                                   {},
		"/v1/address GET":                                     {},
		"/v1/block/:height/blobs GET":                         {},
		"/v1/namespace GET":                                   {},
		"/v1/address/count GET":                               {},
		"/v1/address/:hash/blobs GET":                         {},
		"/v1/block/:height/blobs/count GET":                   {},
		"/v1/namespace/:id GET":                               {},
		"/v1/stats/summary/:table/:function GET":              {},
		"/v1/rollup/:id/blobs GET":                            {},
		"/v1/rollup/:id/distribution/:name/:timeframe GET":    {},
		"/v1/address/:hash/messages GET":                      {},
		"/v1/address/:hash/grants GET":                        {},
		"/v1/address/:hash/granters GET":                      {},
		"/v1/blob POST":                                       {},
		"/v1/blob GET":                                        {},
		"/v1/swagger/doc.json GET":                            {},
		"/v1/enums GET":                                       {},
		"/v1/block/count GET":                                 {},
		"/v1/tx/genesis GET":                                  {},
		"/v1/blob/metadata POST":                              {},
		"/v1/validators/:id/jails GET":                        {},
		"/v1/head GET":                                        {},
		"/v1/address/:hash/stats/:name/:timeframe GET":        {},
		"/v1/block/:height GET":                               {},
		"/v1/tx/:hash GET":                                    {},
		"/v1/tx/:hash/messages GET":                           {},
		"/v1/gas/price GET":                                   {},
		"/v1/gas/price/:priority GET":                         {},
		"/v1/block/:height/ods GET":                           {},
		"/v1/tx/:hash/blobs GET":                              {},
		"/v1/namespace/:id/:version/messages GET":             {},
		"/v1/validators/:id/uptime GET":                       {},
		"/v1/stats/namespace/usage GET":                       {},
		"/v1/rollup/slug/:slug GET":                           {},
		"/v1/block/:height/events GET":                        {},
		"/v1/block/:height/stats GET":                         {},
		"/v1/rollup/:id/export GET":                           {},
		"/v1/docs GET":                                        {},
		"/v1/stats/square_size GET":                           {},
		"/v1/stats/size_groups GET":                           {},
		"/v1/stats/rollup_stats_24h GET":                      {},
		"/v1/stats/messages_count_24h GET":                    {},
		"/v1/stats/ibc/summary GET":                           {},
		"/v1/rollup/stats/series/:timeframe GET":              {},
		"/v1/rollup/group GET":                                {},
		"/v1/proposal GET":                                    {},
		"/v1/proposal/:id GET":                                {},
		"/v1/ibc/client GET":                                  {},
		"/v1/ibc/client/:id GET":                              {},
		"/v1/ibc/connection GET":                              {},
		"/v1/ibc/connection/:id GET":                          {},
		"/v1/ibc/channel GET":                                 {},
		"/v1/ibc/channel/:id GET":                             {},
		"/v1/ibc/transfer GET":                                {},
		"/v1/ibc/transfer/:id GET":                            {},
		"/v1/ibc/relayers GET":                                {},
		"/v1/proposal/:id/votes GET":                          {},
		"/v1/address/:hash/votes GET":                         {},
		"/v1/validators/:id/votes GET":                        {},
		"/v1/blob/proofs POST":                                {},
		"/v1/hyperlane/mailbox GET":                           {},
		"/v1/hyperlane/mailbox/:id GET":                       {},
		"/v1/hyperlane/token GET":                             {},
		"/v1/hyperlane/token/:id GET":                         {},
		"/v1/hyperlane/transfer GET":                          {},
		"/v1/hyperlane/transfer/:id GET":                      {},
		"/v1/hyperlane/domains GET":                           {},
		"/v1/hyperlane/igp GET":                               {},
		"/v1/hyperlane/igp/:id GET":                           {},
		"/v1/signal GET":                                      {},
		"/v1/signal/upgrade GET":                              {},
	}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	db := postgres.Storage{
		Storage: &pg.Storage{},
	}
	apiCfg := ApiConfig{
		Bind:         "127.0.0.1:9876",
		BlobReceiver: "dal_node",
	}

	e := initEcho(apiCfg, "development")
	defer func() {
		err := e.Close()
		require.NoError(t, err)
	}()

	initHandlers(ctx, e, Config{
		Config: &config.Config{
			DataSources: map[string]config.DataSource{
				"node_rpc": {},
				"dal_node": {},
			},
		},
		ApiConfig: apiCfg,
	}, db)

	for _, route := range e.Routes() {
		key := fmt.Sprintf("%s %s", route.Path, route.Method)
		_, ok := expectedRoutes[key]
		require.True(t, ok, "routes in expected", key)
	}

	for key := range expectedRoutes {
		parts := strings.Split(key, " ")
		method := parts[1]
		path := parts[0]

		var found bool
		for _, route := range e.Routes() {
			if route.Path == path && route.Method == method {
				found = true
				break
			}
		}
		require.True(t, found, "expected in routes", key)
	}
}
