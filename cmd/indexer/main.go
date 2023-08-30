package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dipdup-io/celestia-indexer/pkg/indexer"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "indexer",
	Short: "DipDup Verticals | Celestia Indexer",
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

	indexerModule := indexer.New(*cfg)
	if err := indexerModule.Start(ctx); err != nil {
		log.Panic().Err(err).Msg("indexer module start")
		return
	}

	<-notifyCtx.Done()
	cancel()

	if err := indexerModule.Close(); err != nil {
		log.Panic().Err(err).Msg("stopping indexer")
	}

	log.Info().Msg("stopped")
}
