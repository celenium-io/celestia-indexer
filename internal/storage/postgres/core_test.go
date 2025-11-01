// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/stretchr/testify/require"
)

func TestCheckDatabaseExists(t *testing.T) {
	ctx, ctxCancel := context.WithTimeout(t.Context(), 180*time.Second)
	defer ctxCancel()

	containerCfg := database.PostgreSQLContainerConfig{
		User:     "user",
		Password: "password",
		Database: "db_test",
		Port:     5432,
		Image:    "timescale/timescaledb-ha:pg15.8-ts2.17.0-all",
	}

	psqlContainer, err := database.NewPostgreSQLContainer(ctx, containerCfg)
	require.NoError(t, err)

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

	exists, err := checkTablesExists(ctx, db)
	require.NoError(t, err)
	require.False(t, exists)
	require.NoError(t, db.Close())

	strg, err := Create(ctx, cfg, "../../../database", false)
	require.NoError(t, err)

	exists, err = checkTablesExists(ctx, strg.Connection())
	require.NoError(t, err)
	require.True(t, exists)

	require.NoError(t, strg.Close())
	require.NoError(t, psqlContainer.Terminate(ctx))
}
