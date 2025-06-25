// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
	pg "github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/stretchr/testify/require"
)

func TestRoutes(t *testing.T) {
	var expectedRoutes = map[string]struct{}{
		"/v1/auth/rollup/new POST":         {},
		"/v1/auth/rollup/:id PATCH":        {},
		"/v1/auth/rollup/:id/verify PATCH": {},
		"/v1/auth/rollup/unverified GET":   {},
		"/v1/auth/rollup/:id DELETE":       {},
		"/v1/auth/bulk POST":               {},
	}

	db := postgres.Storage{
		Storage: &pg.Storage{},
	}
	apiCfg := ApiConfig{
		Bind: "127.0.0.1:9877",
	}

	e := initEcho(apiCfg)
	defer func() {
		err := e.Close()
		require.NoError(t, err)
	}()

	initHandlers(e, db)

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
