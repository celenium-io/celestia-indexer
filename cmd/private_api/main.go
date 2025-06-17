// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
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
	Use:   "private-api",
	Short: "DipDup Verticals | Celenium API",
}

func main() {
	cfg, err := initConfig()
	if err != nil {
		return
	}

	if err = initLogger(cfg.LogLevel); err != nil {
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	db := initDatabase(cfg.Database, cfg.Indexer.ScriptsDir)
	e := initEcho(cfg.ApiConfig)
	initHandlers(e, db)

	go func() {
		if err := e.Start(cfg.ApiConfig.Bind); err != nil && errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	if err := db.Close(); err != nil {
		e.Logger.Fatal(err)
	}
}
