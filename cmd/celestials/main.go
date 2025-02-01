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
	"github.com/pkg/errors"
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

	addressHandler := func(ctx context.Context, address string) (uint64, error) {
		prefix, hash, err := types.Address(address).Decode()
		if err != nil {
			return 0, errors.Wrap(err, "decoding address")
		}
		if prefix != types.AddressPrefixCelestia {
			return 0, errors.Errorf("invalid prefix: %s", prefix)
		}
		addressId, err := pg.Address.IdByHash(ctx, hash)
		if err != nil || len(address) == 0 {
			return 0, errors.Errorf("can't find address %s in database", address)
		}
		return addressId[0], nil
	}

	module := celestials.New(
		celestialsDatasource,
		addressHandler,
		pg.Celestials,
		pg.CelestialState,
		pg.Transactable,
		cfg.Indexer.Name,
		cfg.Celestials.ChainId,
	)
	module.Start(ctx)

	<-notifyCtx.Done()
	cancel()

	if err := module.Close(); err != nil {
		log.Panic().Err(err).Msg("stopping celestials indexer")
	}

	log.Info().Msg("stopped")
}
