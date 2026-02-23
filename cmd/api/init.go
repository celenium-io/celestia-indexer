// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"io"
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
	"github.com/celenium-io/celestia-indexer/cmd/api/hyperlane"
	"github.com/celenium-io/celestia-indexer/cmd/api/ibc_relayer"
	"github.com/celenium-io/celestia-indexer/internal/blob"
	"github.com/celenium-io/celestia-indexer/internal/profiler"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
	"github.com/celenium-io/celestia-indexer/pkg/node"
	nodeApi "github.com/celenium-io/celestia-indexer/pkg/node/dal"
	"github.com/celenium-io/celestia-indexer/pkg/node/rpc"
	"github.com/dipdup-net/go-lib/config"
	"github.com/getsentry/sentry-go"
	"github.com/grafana/pyroscope-go"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
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
	return c.Path() == "/v1/ws"
}

func metricsSkipper(c echo.Context) bool {
	return c.Path() == "/metrics"
}

func postSkipper(c echo.Context) bool {
	if c.Request().Method != http.MethodPost {
		return true
	}
	if strings.HasPrefix(c.Path(), "/v1/blob") {
		return true
	}
	if strings.HasPrefix(c.Path(), "/v1/auth") {
		return true
	}
	return false
}

func gzipSkipper(c echo.Context) bool {
	if c.Path() == "/v1/swagger/doc.json" {
		return true
	}
	if metricsSkipper(c) {
		return true
	}
	return websocketSkipper(c)
}

func observableCacheSkipper(c echo.Context) bool {
	if c.Request().Method != http.MethodGet {
		return true
	}
	if websocketSkipper(c) {
		return true
	}
	if metricsSkipper(c) {
		return true
	}
	if c.Path() == "/v1/head" {
		return true
	}
	if strings.Contains(c.Path(), "/v1/block/:height") {
		return true
	}
	if strings.Contains(c.Path(), "/v1/tx/:hash") {
		return true
	}
	return false
}

func initEcho(cfg ApiConfig, env string) *echo.Echo {
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

	timeout := 30 * time.Second
	if cfg.RequestTimeout > 0 {
		timeout = time.Duration(cfg.RequestTimeout) * time.Second
	}
	e.Use(RequestTimeout(timeout, websocketSkipper))

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: gzipSkipper,
	}))
	e.Use(middleware.DecompressWithConfig(middleware.DecompressConfig{
		Skipper: websocketSkipper,
	}))
	e.Use(middleware.BodyLimit("9M"))
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

	if err := initSentry(e, cfg.SentryDsn, env); err != nil {
		log.Err(err).Msg("sentry")
	}
	e.Server.IdleTimeout = time.Second * 30

	return e
}

var dispatcher *bus.Dispatcher

func initDispatcher(ctx context.Context, db postgres.Storage) {
	d, err := bus.NewDispatcher(db, db.Validator, db.Constants)
	if err != nil {
		panic(err)
	}
	dispatcher = d
	dispatcher.Start(ctx)
}

func initDatabase(cfg config.Database, viewsDir string) postgres.Storage {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := postgres.Create(ctx, cfg, viewsDir, false)
	if err != nil {
		panic(err)
	}
	return db
}

var ttlCache cache.ICache

func initCache(url string, ttl int) {
	if url != "" {
		if ttl == 0 {
			ttl = 5
		}
		c, err := cache.NewValKey(url, time.Duration(ttl)*time.Minute)
		if err != nil {
			panic(err)
		}
		ttlCache = c
	}
}

func initHandlers(ctx context.Context, e *echo.Echo, cfg Config, db postgres.Storage) {
	if cfg.ApiConfig.Prometheus {
		e.GET("/metrics", echoprometheus.NewHandler())
	}

	v1 := e.Group("v1")

	stateHandlers := handler.NewStateHandler(db.State, db.Validator, db.Constants, cfg.Indexer.Name)
	v1.GET("/head", stateHandlers.Head)

	defaultMiddlewareCache := cache.Middleware(ttlCache, nil, nil)
	statsMiddlewareCache := cache.Middleware(ttlCache, nil, func() time.Duration {
		now := time.Now()
		diff := now.Truncate(time.Hour).Add(time.Hour).Sub(now)
		if diff > time.Minute*10 {
			return time.Minute * 10
		}
		return diff
	})

	constantsHandler := handler.NewConstantHandler(db.Constants, db.DenomMetadata, db.Rollup)
	v1.GET("/constants", constantsHandler.Get, defaultMiddlewareCache)
	v1.GET("/enums", constantsHandler.Enums, defaultMiddlewareCache)

	searchHandler := handler.NewSearchHandler(db.Search, db.Address, db.Blocks, db.Tx, db.Namespace, db.Validator, db.Rollup, db.Celestials)
	v1.GET("/search", searchHandler.Search)

	addressHandlers := handler.NewAddressHandler(db.Address, db.Tx, db.BlobLogs, db.Message, db.Delegation, db.Undelegation, db.Redelegation, db.VestingAccounts, db.Grants, db.Celestials, db.Votes, db.State, cfg.Indexer.Name)
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
			addressGroup.GET("/celestials", addressHandlers.Celestials)
			addressGroup.GET("/votes", addressHandlers.Votes)
			addressGroup.GET("/stats/:name/:timeframe", addressHandlers.Stats, statsMiddlewareCache)
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
			heightGroup.GET("", blockHandlers.Get, defaultMiddlewareCache)
			heightGroup.GET("/events", blockHandlers.GetEvents, defaultMiddlewareCache)
			heightGroup.GET("/messages", blockHandlers.GetMessages, defaultMiddlewareCache)
			heightGroup.GET("/stats", blockHandlers.GetStats, defaultMiddlewareCache)
			heightGroup.GET("/blobs", blockHandlers.Blobs, defaultMiddlewareCache)
			heightGroup.GET("/blobs/count", blockHandlers.BlobsCount, defaultMiddlewareCache)
			heightGroup.GET("/ods", blockHandlers.BlockODS, defaultMiddlewareCache)
		}
	}

	txHandlers := handler.NewTxHandler(db.Tx, db.Event, db.Message, db.Namespace, db.BlobLogs, db.State, cfg.Indexer.Name)
	txGroup := v1.Group("/tx")
	{
		txGroup.GET("", txHandlers.List)
		txGroup.GET("/count", txHandlers.Count)
		txGroup.GET("/genesis", txHandlers.Genesis, defaultMiddlewareCache)
		hashGroup := txGroup.Group("/:hash")
		{
			hashGroup.GET("", txHandlers.Get, defaultMiddlewareCache)
			hashGroup.GET("/events", txHandlers.GetEvents, defaultMiddlewareCache)
			hashGroup.GET("/messages", txHandlers.GetMessages, defaultMiddlewareCache)
			hashGroup.GET("/blobs", txHandlers.Blobs, defaultMiddlewareCache)
			hashGroup.GET("/blobs/count", txHandlers.BlobsCount, defaultMiddlewareCache)
		}
	}

	blobReceiver, err := initBlobReceiver(ctx, cfg)
	if err != nil {
		panic(err)
	}

	namespaceHandlers := handler.NewNamespaceHandler(
		db.Namespace,
		db.BlobLogs,
		db.Rollup,
		db.Address,
		db.State,
		cfg.Indexer.Name,
		blobReceiver,
		&node,
	)

	blobGroup := v1.Group("/blob")
	{
		blobGroup.GET("", namespaceHandlers.Blobs)
		blobGroup.POST("", namespaceHandlers.Blob)
		blobGroup.POST("/metadata", namespaceHandlers.BlobMetadata)
		blobGroup.POST("/proofs", namespaceHandlers.BlobProofs)
	}

	namespaceGroup := v1.Group("/namespace")
	{
		namespaceGroup.GET("", namespaceHandlers.List)
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

	validatorsHandler := handler.NewValidatorHandler(db.Validator, db.Blocks, db.BlockSignatures, db.Delegation, db.Constants, db.Jails, db.Votes, db.State, cfg.Indexer.Name)
	validators := v1.Group("/validators")
	{
		validators.GET("", validatorsHandler.List)
		validators.GET("/count", validatorsHandler.Count)
		validators.GET("/metrics", validatorsHandler.TopNMetrics)
		validator := validators.Group("/:id")
		{
			validator.GET("", validatorsHandler.Get)
			validator.GET("/blocks", validatorsHandler.Blocks)
			validator.GET("/uptime", validatorsHandler.Uptime)
			validator.GET("/delegators", validatorsHandler.Delegators)
			validator.GET("/jails", validatorsHandler.Jails)
			validator.GET("/votes", validatorsHandler.Votes)
			validator.GET("/messages", validatorsHandler.Messages)
			validator.GET("/metrics", validatorsHandler.Metrics)
		}
	}

	statsHandler := handler.NewStatsHandler(db.Stats, db.Namespace, db.IbcTransfers, db.IbcChannels, db.HLTransfer, chainStore, db.State)
	stats := v1.Group("/stats")
	{
		stats.GET("/summary/:table/:function", statsHandler.Summary)
		stats.GET("/tps", statsHandler.TPS)
		stats.GET("/changes_24h", statsHandler.Change24hBlockStats)
		stats.GET("/rollup_stats_24h", statsHandler.RollupStats24h)
		stats.GET("/square_size", statsHandler.SquareSize)
		stats.GET("/messages_count_24h", statsHandler.MessagesCount24h)
		stats.GET("/size_groups", statsHandler.SizeGroups, statsMiddlewareCache)

		namespace := stats.Group("/namespace")
		{
			namespace.GET("/usage", statsHandler.NamespaceUsage)
			namespace.GET("/series/:id/:name/:timeframe", statsHandler.NamespaceSeries, statsMiddlewareCache)
		}
		staking := stats.Group("/staking")
		{
			staking.GET("/series/:id/:name/:timeframe", statsHandler.StakingSeries, statsMiddlewareCache)
			staking.GET("/distribution", statsHandler.StakingDistribution, statsMiddlewareCache)
		}
		ibc := stats.Group("/ibc")
		{
			ibc.GET("/series/:id/:name/:timeframe", statsHandler.IbcSeries, statsMiddlewareCache)
			ibc.GET("/chains", statsHandler.IbcByChains, statsMiddlewareCache)
			ibc.GET("/summary", statsHandler.IbcSummary, statsMiddlewareCache)
		}
		hl := stats.Group("/hyperlane")
		{
			hl.GET("/series/:id/:name/:timeframe", statsHandler.HlSeries, statsMiddlewareCache)
			hl.GET("/chains/:name/:timeframe", statsHandler.HlTotalSeries, statsMiddlewareCache)
			hl.GET("/chains", statsHandler.HlByDomain, statsMiddlewareCache)
		}
		series := stats.Group("/series")
		{
			series.GET("/:name/:timeframe", statsHandler.Series, statsMiddlewareCache)
			series.GET("/:name/:timeframe/cumulative", statsHandler.SeriesCumulative, statsMiddlewareCache)
		}
	}

	gasHandler := handler.NewGasHandler(db.State, db.Tx, db.Constants, db.BlockStats, gasTracker)
	gas := v1.Group("/gas")
	{
		gas.GET("/estimate_for_pfb", gasHandler.EstimateForPfb)
		gas.GET("/price", gasHandler.EstimatePrice)
		gas.GET("/price/:priority", gasHandler.EstimatePricePriority)
	}

	vestingHandler := handler.NewVestingHandler(db.VestingPeriods)
	vesting := v1.Group("/vesting")
	{
		vesting.GET("/:id/periods", vestingHandler.Periods)
	}

	proposalHandler := handler.NewProposalsHandler(db.Proposals, db.Votes, db.Address, db.Validator)
	proposal := v1.Group("/proposal")
	{
		proposal.GET("", proposalHandler.List)
		proposal.GET("/:id", proposalHandler.Get)
		proposal.GET("/:id/votes", proposalHandler.Votes)
	}

	relayers, err := ibc_relayer.NewRelayerStore(ctx, "./assets/relayers_celestia.json", db.Address)
	if err != nil {
		log.Warn().Err(err).Msg("init IBC relayers")
	}
	ibcHandler := handler.NewIbcHandler(db.IbcClients, db.IbcConnections, db.IbcChannels, db.IbcTransfers, db.Address, db.Tx, relayers)
	ibc := v1.Group("/ibc")
	{
		ibcClient := ibc.Group("/client")
		{
			ibcClient.GET("", ibcHandler.List)
			ibcClient.GET("/:id", ibcHandler.Get)
		}
		ibcConnection := ibc.Group("/connection")
		{
			ibcConnection.GET("", ibcHandler.ListConnections)
			ibcConnection.GET("/:id", ibcHandler.GetConnection)
		}
		ibcChannel := ibc.Group("/channel")
		{
			ibcChannel.GET("", ibcHandler.ListChannels)
			ibcChannel.GET("/:id", ibcHandler.GetChannel)
		}

		ibcTransfer := ibc.Group("/transfer")
		{
			ibcTransfer.GET("", ibcHandler.ListTransfers)
			ibcTransfer.GET("/:id", ibcHandler.GetIbcTransfer)
		}
		ibc.GET("/relayers", ibcHandler.IbcRelayers)
	}

	hyperlaneHandler := handler.NewHyperlaneHandler(db.HLMailbox, db.HLToken, db.HLTransfer, db.Tx, db.Address, db.HLIGP, chainStore)
	hyperlane := v1.Group("/hyperlane")
	{
		hlMailbox := hyperlane.Group("/mailbox")
		{
			hlMailbox.GET("", hyperlaneHandler.ListMailboxes)
			hlMailbox.GET("/:id", hyperlaneHandler.GetMailbox)
		}
		hlToken := hyperlane.Group("/token")
		{
			hlToken.GET("", hyperlaneHandler.ListTokens)
			hlToken.GET("/:id", hyperlaneHandler.GetToken)
		}
		hyperlaneTransfer := hyperlane.Group("/transfer")
		{
			hyperlaneTransfer.GET("", hyperlaneHandler.ListTransfers)
			hyperlaneTransfer.GET("/:id", hyperlaneHandler.GetTransfer)
		}
		hlIgp := hyperlane.Group("/igp")
		{
			hlIgp.GET("", hyperlaneHandler.ListIgps)
			hlIgp.GET("/:id", hyperlaneHandler.GetIgp)
		}
		hyperlane.GET("/domains", hyperlaneHandler.ListDomains, defaultMiddlewareCache)

		zkismHandler := handler.NewZkISMHandler(db.ZkISM, db.Address, db.Tx)
		hlZkism := hyperlane.Group("/zkism")
		{
			hlZkism.GET("", zkismHandler.List)
			hlZkism.GET("/:id", zkismHandler.Get)
			hlZkism.GET("/:id/updates", zkismHandler.GetUpdates)
			hlZkism.GET("/:id/messages", zkismHandler.GetMessages)
		}
	}

	signalHandler := handler.NewSignalHandler(db.SignalVersion, db.Upgrade, db.Validator, db.Tx, db.Address)
	signalGroup := v1.Group("/signal")
	{
		signalGroup.GET("", signalHandler.List)
		signalGroup.GET("/upgrade", signalHandler.Upgrades)
		signalGroup.GET("/upgrade/:version", signalHandler.Upgrade)
	}

	fwdHandler := handler.NewForwardingsHandler(db.Forwardings, db.Address, db.Tx, chainStore)
	forwarding := v1.Group("/forwarding")
	{
		forwarding.GET("", fwdHandler.List)
		forwarding.GET("/:id", fwdHandler.Get)
	}

	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL: "./docs/swagger.json",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "Celenium API",
		},
		DarkMode:    true,
		ShowSidebar: true,
	})
	if err != nil {
		panic(err)
	}

	v1.GET("/docs", func(c echo.Context) error {
		return c.HTML(http.StatusOK, htmlContent)
	})

	f, err := os.Open("./docs/swagger.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	docsJson, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	v1.GET("/swagger/doc.json", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/json", docsJson)
	})

	if cfg.ApiConfig.Websocket {
		initWebsocket(ctx, v1)
	}

	rollupHandler := handler.NewRollupHandler(db.Rollup, db.RollupProvider, db.Namespace, db.BlobLogs)
	rollups := v1.Group("/rollup")
	{
		rollups.GET("", rollupHandler.Leaderboard)
		rollups.GET("/count", rollupHandler.Count)
		rollups.GET("/day", rollupHandler.LeaderboardDay)
		rollups.GET("/group", rollupHandler.RollupGroupedStats, statsMiddlewareCache)
		rollups.GET("/stats/series/:timeframe", rollupHandler.AllSeries, statsMiddlewareCache)
		rollups.GET("/slug/:slug", rollupHandler.BySlug)
		rollup := rollups.Group("/:id")
		{
			rollup.GET("", rollupHandler.Get)
			rollup.GET("/namespaces", rollupHandler.GetNamespaces)
			rollup.GET("/providers", rollupHandler.GetProviders)
			rollup.GET("/blobs", rollupHandler.GetBlobs)
			rollup.GET("/stats/:name/:timeframe", rollupHandler.Stats, statsMiddlewareCache)
			rollup.GET("/distribution/:name/:timeframe", rollupHandler.Distribution, statsMiddlewareCache)
			rollup.GET("/export", rollupHandler.ExportBlobs)
		}
	}

	log.Info().Msg("API routes:")
	for _, route := range e.Routes() {
		log.Info().Msgf("[%s] %s -> %s", route.Method, route.Path, route.Name)
	}
}

func initSentry(e *echo.Echo, dsn, environment string) error {
	if dsn == "" {
		return nil
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		AttachStacktrace: true,
		Environment:      environment,
		EnableTracing:    true,
		TracesSampleRate: 0.1,
		Release:          os.Getenv("TAG"),
		IgnoreTransactions: []string{
			"GET /v1/ws",
		},
	}); err != nil {
		return errors.Wrap(err, "initialization")
	}

	e.Use(SentryMiddleware())

	return nil
}

var (
	wsManager *websocket.Manager
)

func initWebsocket(ctx context.Context, group *echo.Group) {
	observer := dispatcher.Observe(storage.ChannelHead, storage.ChannelBlock)
	wsManager = websocket.NewManager(observer)
	if gasTracker != nil {
		gasTracker.SubscribeOnCompute(wsManager.GasTrackerHandler)
	}
	wsManager.Start(ctx)
	group.GET("/ws", wsManager.Handle)
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

func initBlobReceiver(_ context.Context, cfg Config) (node.DalApi, error) {
	switch cfg.ApiConfig.BlobReceiver {
	case "celenium_blobs":
		datasource, ok := cfg.DataSources[cfg.ApiConfig.BlobReceiver]
		if !ok {
			return nil, errors.Errorf("unknown data source pointed in blob_receiver: %s", cfg.ApiConfig.BlobReceiver)
		}
		celeniumBlobReceiver := blob.NewCelenium(datasource)
		return celeniumBlobReceiver, nil
	default:
		datasource, ok := cfg.DataSources[cfg.ApiConfig.BlobReceiver]
		if !ok {
			return nil, errors.Errorf("unknown data source pointed in blob_receiver: %s", cfg.ApiConfig.BlobReceiver)
		}

		return nodeApi.New(datasource.URL).
			WithAuthToken(os.Getenv("CELESTIA_NODE_AUTH_TOKEN")).
			WithRateLimit(datasource.RequestsPerSecond), nil
	}
}

var chainStore *hyperlane.ChainStore

func initChainStore(ctx context.Context, url string) {
	if url != "" {
		chainStore = hyperlane.NewChainStore(url)
		chainStore.Start(ctx)
	}
}
