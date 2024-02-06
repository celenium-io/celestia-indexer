// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
	"github.com/celenium-io/celestia-indexer/pkg/quotes"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "indexer",
	Short: "DipDup Verticals | Celenium Quotes Indexer",
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

	binanceDatasource, ok := cfg.DataSources["binance"]
	if !ok {
		log.Panic().Err(err).Msg("can't find binance data source")
		return
	}

	pg, err := postgres.Create(ctx, cfg.Database, cfg.Indexer.ScriptsDir)
	if err != nil {
		log.Panic().Err(err).Msg("can't create database connection")
		return
	}

	module := quotes.New(binanceDatasource, pg.Price)
	module.Start(ctx)

	<-notifyCtx.Done()
	cancel()

	if err := module.Close(); err != nil {
		log.Panic().Err(err).Msg("stopping quotes module")
	}

	log.Info().Msg("stopped")
}
