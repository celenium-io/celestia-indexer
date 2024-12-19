// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
	"github.com/celenium-io/celestia-indexer/pkg/tvl"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

var rootCmd = &cobra.Command{
	Use:   "indexer",
	Short: "DipDup Verticals | Celenium TVL Scanner",
}

func main() {
	cfg, err := initConfig()
	if err != nil {
		return
	}

	if err = initLogger(cfg.LogLevel); err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	notifyCtx, notifyCancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer notifyCancel()

	l2beatDs, ok := cfg.DataSources["l2beat"]
	if !ok {
		log.Panic().Err(err).Msg("can't find l2beat data source")
		return
	}

	lamaDs, ok := cfg.DataSources["lama"]
	if !ok {
		log.Panic().Err(err).Msg("can't find lama data source")
		return
	}

	pg, err := postgres.Create(ctx, cfg.Database, cfg.Indexer.ScriptsDir, false)
	if err != nil {
		log.Panic().Err(err).Msg("can't create database connection")
		return
	}

	module := tvl.New(l2beatDs, lamaDs, pg.Rollup, pg.Tvl)
	module.Start(ctx)

	<-notifyCtx.Done()
	cancel()

	if err := module.Close(); err != nil {
		log.Panic().Err(err).Msg("stopping TVL scanner")
	}

	log.Info().Msg("stopped")
}
