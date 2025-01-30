// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	celestials "github.com/celenium-io/celestial-module/pkg/module"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "celestials",
	Short: "DipDup Verticals | Celestials indexer for Celenium",
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

	celestialsDatasource, ok := cfg.DataSources["celestials"]
	if !ok {
		log.Panic().Err(err).Msg("can't find celestials data source")
		return
	}

	pg, err := postgres.Create(ctx, cfg.Database, cfg.Indexer.ScriptsDir, false)
	if err != nil {
		log.Panic().Err(err).Msg("can't create database connection")
		return
	}

	log.Info().Str("chain_id", cfg.Celestials.ChainId).Msg("running module")

	module := celestials.New(
		celestialsDatasource,
		pg.Address,
		pg.Celestials,
		pg.CelestialState,
		pg.Transactable,
		cfg.Indexer.Name,
		cfg.Celestials.ChainId,
		celestials.WithAddressPrefix(types.AddressPrefixCelestia),
	)
	module.Start(ctx)

	<-notifyCtx.Done()
	cancel()

	if err := module.Close(); err != nil {
		log.Panic().Err(err).Msg("stopping celestials indexer")
	}

	log.Info().Msg("stopped")
}
