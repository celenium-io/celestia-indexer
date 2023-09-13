package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dipdup-io/celestia-indexer/pkg/stopper"

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
	if err = initProflier(cfg.Profiler); err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	notifyCtx, notifyCancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer notifyCancel()

	stopperModule := stopper.NewModule(cancel)
	stopperInput, err := stopperModule.Input(stopper.InputName)
	if err != nil {
		log.Err(err).Msg("while getting stopper input")
		return
	}

	indexerModule, err := indexer.New(ctx, *cfg, stopperInput)
	if err != nil {
		log.Panic().Err(err).Msg("error during indexer module creation")
		return
	}

	stopperModule.Start(ctx)
	indexerModule.Start(ctx)

	<-notifyCtx.Done()
	cancel()

	if err := indexerModule.Close(); err != nil {
		log.Panic().Err(err).Msg("stopping indexer")
	}

	if prscp != nil {
		if err := prscp.Stop(); err != nil {
			log.Panic().Err(err).Msg("stopping pyroscope")
		}
	}

	log.Info().Msg("stopped")
}
