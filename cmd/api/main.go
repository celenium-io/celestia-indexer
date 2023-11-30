// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	_ "github.com/celenium-io/celestia-indexer/cmd/api/docs"
)

var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "DipDup Verticals | Celenium API",
}

// @title						Swagger Celenium API
// @version					    1.0
// @description				    This is docs of Celenium API.
// @host						api.celenium.io
//
// @query.collection.format	multi
func main() {
	cfg, err := initConfig()
	if err != nil {
		return
	}

	if err = initLogger(cfg.LogLevel); err != nil {
		return
	}

	if err := initProflier(cfg.Profiler); err != nil {
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	db := initDatabase(cfg.Database, cfg.Indexer.ScriptsDir)
	e := initEcho(cfg.ApiConfig, cfg.Environment)
	initDispatcher(ctx, db)
	initGasTracker(ctx, db)
	initHandlers(ctx, e, *cfg, db)
	initCache(ctx, e)

	go func() {
		if err := e.Start(cfg.ApiConfig.Bind); err != nil && errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	cancel()

	if gasTracker != nil {
		if err := gasTracker.Close(); err != nil {
			e.Logger.Fatal(err)
		}
	}

	if wsManager != nil {
		if err := wsManager.Close(); err != nil {
			e.Logger.Fatal(err)
		}
	}
	if endpointCache != nil {
		if err := endpointCache.Close(); err != nil {
			e.Logger.Fatal(err)
		}
	}
	if dispatcher != nil {
		if err := dispatcher.Close(); err != nil {
			e.Logger.Fatal(err)
		}
	}

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	if prscp != nil {
		if err := prscp.Stop(); err != nil {
			e.Logger.Fatal(err)
		}
	}
}
