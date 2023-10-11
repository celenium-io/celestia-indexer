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
)

var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "DipDup Verticals | Celestia API",
}

// @title						Swagger Celestia Indexer API
// @version					1.0
// @description				This is docs of Celestia indexer API.
// @host						https://api.celestia.dipdup.net
// @BasePath					/v1
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

	db := initDatabase(cfg.Database)
	e := initEcho(cfg.ApiConfig)
	initHandlers(ctx, e, *cfg, db)

	go func() {
		if err := e.Start(cfg.ApiConfig.Bind); err != nil && errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	cancel()

	if err := wsManager.Close(); err != nil {
		e.Logger.Fatal(err)
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
