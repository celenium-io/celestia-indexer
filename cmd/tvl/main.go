// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	indexer "github.com/celenium-io/celestia-indexer/pkg/tvl"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/stopper"

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

	stopperModule := stopper.NewModule(cancel)
	indexerModule, err := indexer.New(ctx, *cfg, &stopperModule)
	if err != nil {
		log.Panic().Err(err).Msg("error during indexer module creation")
		return
	}

	stopperModule.Start(ctx)
	indexerModule.Start(ctx)

	<-notifyCtx.Done()
	cancel()

	log.Info().Msg("stopped")
}
