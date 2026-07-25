// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/dipdup-io/go-lib/config"
	"github.com/dipdup-io/go-lib/database"
	"github.com/dipdup-io/go-lib/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestCheckDatabaseExists(t *testing.T) {
	ctx, ctxCancel := context.WithTimeout(t.Context(), 180*time.Second)
	defer ctxCancel()

	containerCfg := testhelpers.PostgreSQLContainerConfig{
		User:     "user",
		Password: "password",
		Database: "db_test",
		Port:     5432,
		Image:    "timescale/timescaledb-ha:pg15.8-ts2.17.0-all",
	}

	psqlContainer, err := testhelpers.NewPostgreSQLContainer(ctx, containerCfg)
	require.NoError(t, err)
	t.Cleanup(func() {
		ctx, ctxCancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer ctxCancel()
		require.NoError(t, psqlContainer.Terminate(ctx))
	})

	cfg := config.Database{
		Kind:     config.DBKindPostgres,
		User:     psqlContainer.Config.User,
		Database: psqlContainer.Config.Database,
		Password: psqlContainer.Config.Password,
		Host:     psqlContainer.Config.Host,
		Port:     psqlContainer.MappedPort().Int(),
	}

	db := database.NewBun()
	err = db.Connect(ctx, cfg)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	exists, err := checkTablesExists(ctx, db)
	require.NoError(t, err)
	require.False(t, exists)

	strg, err := Create(ctx, cfg, "../../../database", false)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, strg.Close())
	})

	exists, err = checkTablesExists(ctx, strg.Connection())
	require.NoError(t, err)
	require.True(t, exists)
}
