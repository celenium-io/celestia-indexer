// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/celenium-io/celestia-indexer/cmd/api/bus"
	"github.com/celenium-io/celestia-indexer/cmd/api/cache"
	"github.com/celenium-io/celestia-indexer/cmd/api/gas"
	"github.com/celenium-io/celestia-indexer/cmd/api/handler"
	"github.com/celenium-io/celestia-indexer/cmd/api/handler/websocket"
	"github.com/celenium-io/celestia-indexer/internal/profiler"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
	nodeApi "github.com/celenium-io/celestia-indexer/pkg/node/dal"
	"github.com/celenium-io/celestia-indexer/pkg/node/rpc"
	"github.com/dipdup-net/go-lib/config"
	"github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"github.com/grafana/pyroscope-go"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/time/rate"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	})
}

func initConfig() (*Config, error) {
	configPath := rootCmd.PersistentFlags().StringP("config", "c", "dipdup.yml", "path to YAML config file")
	if err := rootCmd.Execute(); err != nil {
		log.Panic().Err(err).Msg("command line execute")
		return nil, err
	}

	if err := rootCmd.MarkFlagRequired("config"); err != nil {
		log.Panic().Err(err).Msg("config command line arg is required")
		return nil, err
	}

	var cfg Config
	if err := config.Parse(*configPath, &cfg); err != nil {
		log.Panic().Err(err).Msg("parsing config file")
		return nil, err
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = zerolog.LevelInfoValue
	}

	return &cfg, nil
}

func initLogger(level string) error {
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Panic().Err(err).Msg("parsing log level")
		return err
	}
	zerolog.SetGlobalLevel(logLevel)
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	log.Logger = log.Logger.With().Caller().Logger()

	return nil
}

var prscp *pyroscope.Profiler

func initProflier(cfg *profiler.Config) (err error) {
	prscp, err = profiler.New(cfg, "api")
	return
}

func websocketSkipper(c echo.Context) bool {
	return strings.Contains(c.Request().URL.Path, "ws")
}

func postSkipper(c echo.Context) bool {
	if strings.Contains(c.Request().URL.Path, "blob") {
		return true
	}
	if strings.Contains(c.Request().URL.Path, "auth/rollup") {
		return true
	}
	return false
}

func gzipSkipper(c echo.Context) bool {
	if strings.Contains(c.Request().URL.Path, "swagger") {
		return true
	}
	if strings.Contains(c.Request().URL.Path, "metrics") {
		return true
	}
	return websocketSkipper(c)
}

func cacheSkipper(c echo.Context) bool {
	if c.Request().Method != http.MethodGet {
		return true
	}
	if websocketSkipper(c) {
		return true
	}
	if strings.Contains(c.Request().URL.Path, "metrics") {
		return true
	}
	if strings.Contains(c.Request().URL.Path, "head") {
		return true
	}
	return false
}

func initEcho(cfg ApiConfig, db postgres.Storage, env string) *echo.Echo {
	e := echo.New()
	e.Validator = handler.NewCelestiaApiValidator()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogLatency:   true,
		LogMethod:    true,
		LogUserAgent: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			switch {
			case v.Status == http.StatusOK || v.Status == http.StatusNoContent:
				log.Info().
					Str("uri", v.URI).
					Int("status", v.Status).
					Dur("latency", v.Latency).
					Str("method", v.Method).
					Str("user-agent", v.UserAgent).
					Str("ip", c.RealIP()).
					Msg("request")
			case v.Status >= 500:
				log.Error().
					Str("uri", v.URI).
					Int("status", v.Status).
					Dur("latency", v.Latency).
					Str("method", v.Method).
					Str("user-agent", v.UserAgent).
					Str("ip", c.RealIP()).
					Msg("request")
			default:
				log.Warn().
					Str("uri", v.URI).
					Int("status", v.Status).
					Dur("latency", v.Latency).
					Str("method", v.Method).
					Str("user-agent", v.UserAgent).
					Str("ip", c.RealIP()).
					Msg("request")
			}

			return nil
		},
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: gzipSkipper,
	}))
	e.Use(middleware.DecompressWithConfig(middleware.DecompressConfig{
		Skipper: websocketSkipper,
	}))
	e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.CSRFWithConfig(
		middleware.CSRFConfig{
			Skipper: func(c echo.Context) bool {
				return websocketSkipper(c) || postSkipper(c)
			},
		},
	))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Pre(middleware.RemoveTrailingSlash())

	timeout := 30 * time.Second
	if cfg.RequestTimeout > 0 {
		timeout = time.Duration(cfg.RequestTimeout) * time.Second
	}
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper: websocketSkipper,
		Timeout: timeout,
	}))

	if cfg.Prometheus {
		e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
			Namespace: "celestia_api",
			Skipper:   websocketSkipper,
		}))
	}
	if cfg.RateLimit > 0 {
		e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
			Skipper: websocketSkipper,
			Store:   middleware.NewRateLimiterMemoryStore(rate.Limit(cfg.RateLimit)),
		}))

	}

	if err := initSentry(e, db, cfg.SentryDsn, env); err != nil {
		log.Err(err).Msg("sentry")
	}
	e.Server.IdleTimeout = time.Second * 30

	return e
}

var dispatcher *bus.Dispatcher

func initDispatcher(ctx context.Context, db postgres.Storage) {
	d, err := bus.NewDispatcher(db, db.Blocks, db.Validator)
	if err != nil {
		panic(err)
	}
	dispatcher = d
	dispatcher.Start(ctx)
}

func initDatabase(cfg config.Database, viewsDir string) postgres.Storage {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := postgres.Create(ctx, cfg, viewsDir)
	if err != nil {
		panic(err)
	}
	return db
}

func initHandlers(ctx context.Context, e *echo.Echo, cfg Config, db postgres.Storage) {
	v1 := e.Group("v1")

	stateHandlers := handler.NewStateHandler(db.State, db.Validator, cfg.Indexer.Name)
	v1.GET("/head", stateHandlers.Head)
	constantsHandler := handler.NewConstantHandler(db.Constants, db.DenomMetadata, db.Address)
	v1.GET("/constants", constantsHandler.Get)
	v1.GET("/enums", constantsHandler.Enums)

	searchHandler := handler.NewSearchHandler(db.Search, db.Address, db.Blocks, db.Tx, db.Namespace, db.Validator, db.Rollup)
	v1.GET("/search", searchHandler.Search)

	addressHandlers := handler.NewAddressHandler(db.Address, db.Tx, db.BlobLogs, db.Message, db.Delegation, db.Undelegation, db.Redelegation, db.VestingAccounts, db.Grants, db.State, cfg.Indexer.Name)
	addressesGroup := v1.Group("/address")
	{
		addressesGroup.GET("", addressHandlers.List)
		addressesGroup.GET("/count", addressHandlers.Count)
		addressGroup := addressesGroup.Group("/:hash")
		{
			addressGroup.GET("", addressHandlers.Get)
			addressGroup.GET("/txs", addressHandlers.Transactions)
			addressGroup.GET("/messages", addressHandlers.Messages)
			addressGroup.GET("/blobs", addressHandlers.Blobs)
			addressGroup.GET("/delegations", addressHandlers.Delegations)
			addressGroup.GET("/undelegations", addressHandlers.Undelegations)
			addressGroup.GET("/redelegations", addressHandlers.Redelegations)
			addressGroup.GET("/vestings", addressHandlers.Vestings)
			addressGroup.GET("/grants", addressHandlers.Grants)
			addressGroup.GET("/granters", addressHandlers.Grantee)
			addressGroup.GET("/stats/:name/:timeframe", addressHandlers.Stats)
		}
	}
	ds, ok := cfg.DataSources["node_rpc"]
	if !ok {
		panic("can't find node data source: node_rpc")
	}
	node := rpc.NewAPI(ds)

	blockHandlers := handler.NewBlockHandler(db.Blocks, db.BlockStats, db.Event, db.Namespace, db.Message, db.BlobLogs, db.State, &node, cfg.Indexer.Name)
	blockGroup := v1.Group("/block")
	{
		blockGroup.GET("", blockHandlers.List)
		blockGroup.GET("/count", blockHandlers.Count)
		heightGroup := blockGroup.Group("/:height")
		{
			heightGroup.GET("", blockHandlers.Get)
			heightGroup.GET("/events", blockHandlers.GetEvents)
			heightGroup.GET("/messages", blockHandlers.GetMessages)
			heightGroup.GET("/stats", blockHandlers.GetStats)
			heightGroup.GET("/namespace", blockHandlers.GetNamespaces)
			heightGroup.GET("/namespace/count", blockHandlers.GetNamespacesCount)
			heightGroup.GET("/blobs", blockHandlers.Blobs)
			heightGroup.GET("/blobs/count", blockHandlers.BlobsCount)
			heightGroup.GET("/ods", blockHandlers.BlockODS)
		}
	}

	txHandlers := handler.NewTxHandler(db.Tx, db.Event, db.Message, db.Namespace, db.BlobLogs, db.State, cfg.Indexer.Name)
	txGroup := v1.Group("/tx")
	{
		txGroup.GET("", txHandlers.List)
		txGroup.GET("/count", txHandlers.Count)
		txGroup.GET("/genesis", txHandlers.Genesis)
		hashGroup := txGroup.Group("/:hash")
		{
			hashGroup.GET("", txHandlers.Get)
			hashGroup.GET("/events", txHandlers.GetEvents)
			hashGroup.GET("/messages", txHandlers.GetMessages)
			hashGroup.GET("/namespace", txHandlers.Namespaces)
			hashGroup.GET("/namespace/count", txHandlers.NamespacesCount)
			hashGroup.GET("/blobs", txHandlers.Blobs)
			hashGroup.GET("/blobs/count", txHandlers.BlobsCount)
		}
	}

	datasource, ok := cfg.DataSources[cfg.ApiConfig.BlobReceiver]
	if !ok {
		panic(fmt.Sprintf("unknown data source pointed in blob_receiver: %s", cfg.ApiConfig.BlobReceiver))
	}

	blobReceiver := nodeApi.New(datasource.URL).
		WithAuthToken(os.Getenv("CELESTIA_NODE_AUTH_TOKEN")).
		WithRateLimit(datasource.RequestsPerSecond)

	namespaceHandlers := handler.NewNamespaceHandler(db.Namespace, db.BlobLogs, db.Rollup, db.State, cfg.Indexer.Name, blobReceiver)
	v1.POST("/blob", namespaceHandlers.Blob)

	namespaceGroup := v1.Group("/namespace")
	{
		namespaceGroup.GET("", namespaceHandlers.List)
		namespaceGroup.GET("/count", namespaceHandlers.Count)
		namespaceGroup.GET("/active", namespaceHandlers.GetActive)
		namespaceGroup.GET("/:id", namespaceHandlers.Get)
		namespaceGroup.GET("/:id/:version", namespaceHandlers.GetWithVersion)
		namespaceGroup.GET("/:id/:version/messages", namespaceHandlers.GetMessages)
		namespaceGroup.GET("/:id/:version/blobs", namespaceHandlers.GetBlobLogs)
		namespaceGroup.GET("/:id/:version/rollups", namespaceHandlers.Rollups)
	}

	namespaceByHash := v1.Group("/namespace_by_hash")
	{
		namespaceByHash.GET("/:hash", namespaceHandlers.GetByHash)
		namespaceByHash.GET("/:hash/:height", namespaceHandlers.GetBlobs)
	}

	validatorsHandler := handler.NewValidatorHandler(db.Validator, db.Blocks, db.BlockSignatures, db.Delegation, db.Constants, db.Jails, db.State, cfg.Indexer.Name)
	validators := v1.Group("/validators")
	{
		validators.GET("", validatorsHandler.List)
		validators.GET("/count", validatorsHandler.Count)
		validator := validators.Group("/:id")
		{
			validator.GET("", validatorsHandler.Get)
			validator.GET("/blocks", validatorsHandler.Blocks)
			validator.GET("/uptime", validatorsHandler.Uptime)
			validator.GET("/delegators", validatorsHandler.Delegators)
			validator.GET("/jails", validatorsHandler.Jails)
		}
	}

	statsHandler := handler.NewStatsHandler(db.Stats, db.Namespace, db.Price, db.State)
	stats := v1.Group("/stats")
	{
		stats.GET("/summary/:table/:function", statsHandler.Summary)
		stats.GET("/histogram/:table/:function/:timeframe", statsHandler.Histogram)
		stats.GET("/tps", statsHandler.TPS)
		stats.GET("/tx_count_24h", statsHandler.TxCountHourly24h)

		price := stats.Group("/price")
		{
			price.GET("/current", statsHandler.PriceCurrent)
			price.GET("/series/:timeframe", statsHandler.PriceSeries)
		}

		namespace := stats.Group("/namespace")
		{
			namespace.GET("/usage", statsHandler.NamespaceUsage)
			namespace.GET("/series/:id/:name/:timeframe", statsHandler.NamespaceSeries)
		}
		staking := stats.Group("/staking")
		{
			staking.GET("/series/:id/:name/:timeframe", statsHandler.StakingSeries)
		}
		series := stats.Group("/series")
		{
			series.GET("/:name/:timeframe", statsHandler.Series)
		}
	}

	gasHandler := handler.NewGasHandler(db.State, db.Tx, db.Constants, db.BlockStats, gasTracker)
	gas := v1.Group("/gas")
	{
		gas.GET("/estimate_for_pfb", gasHandler.EstimateForPfb)
		gas.GET("/price", gasHandler.EstimatePrice)
	}

	vestingHandler := handler.NewVestingHandler(db.VestingPeriods)
	vesting := v1.Group("/vesting")
	{
		vesting.GET("/:id/periods", vestingHandler.Periods)
	}

	if cfg.ApiConfig.Prometheus {
		v1.GET("/metrics", echoprometheus.NewHandler())
	}

	v1.GET("/swagger/*", echoSwagger.WrapHandler)

	if cfg.ApiConfig.Websocket {
		initWebsocket(ctx, v1)
	}

	rollupHandler := handler.NewRollupHandler(db.Rollup, db.Namespace, db.BlobLogs)
	rollups := v1.Group("/rollup")
	{
		rollups.GET("", rollupHandler.Leaderboard)
		rollups.GET("/count", rollupHandler.Count)
		rollups.GET("/slug/:slug", rollupHandler.BySlug)
		rollup := rollups.Group("/:id")
		{
			rollup.GET("", rollupHandler.Get)
			rollup.GET("/namespaces", rollupHandler.GetNamespaces)
			rollup.GET("/blobs", rollupHandler.GetBlobs)
			rollup.GET("/stats/:name/:timeframe", rollupHandler.Stats)
			rollup.GET("/distribution/:name/:timeframe", rollupHandler.Distribution)
			rollup.GET("/export", rollupHandler.ExportBlobs)
		}
	}

	auth := v1.Group("/auth")
	{
		keyMiddleware := middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
			KeyLookup: "header:Authorization",
			Validator: func(key string, c echo.Context) (bool, error) {
				return key == os.Getenv("API_AUTH_KEY"), nil
			},
		})
		auth.Use(keyMiddleware)

		rollupAuthHandler := handler.NewRollupAuthHandler(db.Rollup, db.Address, db.Namespace, db.Transactable)
		rollup := auth.Group("/rollup")
		{
			rollup.POST("/new", rollupAuthHandler.Create)
			rollup.PATCH("/:id", rollupAuthHandler.Update)
			rollup.DELETE("/:id", rollupAuthHandler.Delete)
		}
	}

	log.Info().Msg("API routes:")
	for _, route := range e.Routes() {
		log.Info().Msgf("[%s] %s -> %s", route.Method, route.Path, route.Name)
	}
}

func initSentry(e *echo.Echo, db postgres.Storage, dsn, environment string) error {
	if dsn == "" {
		return nil
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		AttachStacktrace: true,
		Environment:      environment,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		Release:          os.Getenv("TAG"),
	}); err != nil {
		return errors.Wrap(err, "initialization")
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(sentryotel.NewSentryPropagator())

	db.SetTracer(tp)

	e.Use(SentryMiddleware())

	return nil
}

var (
	wsManager     *websocket.Manager
	endpointCache *cache.Cache
)

func initWebsocket(ctx context.Context, group *echo.Group) {
	observer := dispatcher.Observe(storage.ChannelHead, storage.ChannelBlock)
	wsManager = websocket.NewManager(observer)
	wsManager.Start(ctx)
	group.GET("/ws", wsManager.Handle)
}

func initCache(ctx context.Context, e *echo.Echo) {
	observer := dispatcher.Observe(storage.ChannelHead)
	endpointCache = cache.NewCache(cache.Config{
		MaxEntitiesCount: 1000,
	}, observer)
	e.Use(cache.Middleware(endpointCache, cacheSkipper))
	endpointCache.Start(ctx)
}

var gasTracker *gas.Tracker

func initGasTracker(ctx context.Context, db postgres.Storage) {
	observer := dispatcher.Observe(storage.ChannelBlock)
	gasTracker = gas.NewTracker(db.State, db.BlockStats, db.Tx, observer)
	if err := gasTracker.Init(ctx); err != nil {
		panic(err)
	}
	gasTracker.Start(ctx)
}
