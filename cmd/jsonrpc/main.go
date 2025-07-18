// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

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
	e := initEcho(cfg.JsonRpcConfig, cfg.Environment)
	receiver, err := initBlobReceiver(ctx, *cfg)
	if err != nil {
		return
	}
	initHandlers(e, *cfg, db, receiver)

	go func() {
		if err := e.Start(cfg.JsonRpcConfig.Bind); err != nil && errors.Is(err, http.ErrServerClosed) {
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
