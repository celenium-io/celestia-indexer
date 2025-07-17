// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
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
	Short: "DipDup Verticals | Celenium API",
}

// @title					Celenium API
// @version					1.0
// @description				Celenium API is a powerful tool to access all blockchain data that is processed and indexed by our proprietary indexer. With Celenium API you can retrieve all historical data, off-chain data, blobs and statistics through our REST API. Celenium API indexer are open source, which allows you to not depend on third-party services. You can clone, build and run them independently, giving you full control over all components. If you have any questions or feature requests, please feel free to contact us. We appreciate your feedback!
// @host					api-mainnet.celenium.io
// @schemes					https
// @BasePath				/v1
//
// @contact.name			Support
// @contact.url				https://discord.gg/3k83Przqk8
// @contact.email			celenium@pklabs.me
//
// @externalDocs.description	Full documentation
// @externalDocs.url			https://api-docs.celenium.io/
//
// @x-servers					[{"url": "api-mainnet.celenium.io", "description": "Celenium Mainnet API"},{"url": "api-mocha.celenium.io", "description": "Celenium Mocha API"},{"url": "api-arabica.celenium.io", "description": "Celenium Arabica API"}]
// @query.collection.format	multi
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						apikey
// @description					To authorize your requests you have to select the required tariff on our site. Then you receive api key to authorize. Api key should be passed via request header `apikey`.
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

	initCache(cfg.ApiConfig.Cache)
	db := initDatabase(cfg.Database, cfg.Indexer.ScriptsDir)
	e := initEcho(cfg.ApiConfig, cfg.Environment)
	initDispatcher(ctx, db)
	initGasTracker(ctx, db)
	initHandlers(ctx, e, *cfg, db)
	initChainStore(ctx, *cfg)

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
	if dispatcher != nil {
		if err := dispatcher.Close(); err != nil {
			e.Logger.Fatal(err)
		}
	}
	if prscp != nil {
		if err := prscp.Stop(); err != nil {
			e.Logger.Fatal(err)
		}
	}
	if err := db.Close(); err != nil {
		e.Logger.Fatal(err)
	}
	if ttlCache != nil {
		if err := ttlCache.Close(); err != nil {
			e.Logger.Fatal(err)
		}
	}
	if chainStore != nil {
		if err := chainStore.Close(); err != nil {
			e.Logger.Fatal(err)
		}
	}
}
