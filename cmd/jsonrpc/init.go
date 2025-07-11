// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/blob"
	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
	"github.com/celenium-io/celestia-indexer/pkg/node"
	nodeApi "github.com/celenium-io/celestia-indexer/pkg/node/dal"
	"github.com/dipdup-net/go-lib/config"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
)

var rootCmd = &cobra.Command{
	Use:   "jsonrpc",
	Short: "DipDup Verticals | Celenium Json Rpc",
}

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

func initDatabase(cfg config.Database, viewsDir string) postgres.Storage {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := postgres.Create(ctx, cfg, viewsDir, false)
	if err != nil {
		panic(err)
	}
	return db
}

func initEcho(cfg JsonRpcConfig, env string) *echo.Echo {
	e := echo.New()

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
	e.Use(RequestTimeout(timeout, nil))

	e.Use(middleware.Gzip())
	e.Use(middleware.Decompress())
	e.Use(middleware.BodyLimit("9M"))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		Skipper: func(c echo.Context) bool {
			return true
		},
	}))
	e.Pre(middleware.RemoveTrailingSlash())

	if cfg.Prometheus {
		e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
			Namespace: "celestia_json_rpc",
		}))
	}
	if cfg.RateLimit > 0 {
		e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
			Store: middleware.NewRateLimiterMemoryStore(rate.Limit(cfg.RateLimit)),
		}))

	}

	if err := initSentry(e, cfg.SentryDsn, env); err != nil {
		log.Err(err).Msg("sentry")
	}
	e.Server.IdleTimeout = time.Second * 30

	return e
}

func initBlobReceiver(ctx context.Context, cfg Config) (node.DalApi, error) {
	switch cfg.JsonRpcConfig.BlobReceiver {
	case "r2":
		r2 := blob.NewR2(blob.R2Config{
			BucketName:      os.Getenv("R2_BUCKET"),
			AccountId:       os.Getenv("R2_ACCOUNT_ID"),
			AccessKeyId:     os.Getenv("R2_ACCESS_KEY_ID"),
			AccessKeySecret: os.Getenv("R2_ACCESS_KEY_SECRET"),
		})
		err := r2.Init(ctx)
		return r2, err
	case "celenium_blobs":
		datasource, ok := cfg.DataSources[cfg.JsonRpcConfig.BlobReceiver]
		if !ok {
			return nil, errors.Errorf("unknown data source pointed in blob_receiver: %s", cfg.JsonRpcConfig.BlobReceiver)
		}
		celeniumBlobReceiver := blob.NewCelenium(datasource)
		return celeniumBlobReceiver, nil
	default:
		datasource, ok := cfg.DataSources[cfg.JsonRpcConfig.BlobReceiver]
		if !ok {
			return nil, errors.Errorf("unknown data source pointed in blob_receiver: %s", cfg.JsonRpcConfig.BlobReceiver)
		}

		return nodeApi.New(datasource.URL).
			WithAuthToken(os.Getenv("CELESTIA_NODE_AUTH_TOKEN")).
			WithRateLimit(datasource.RequestsPerSecond), nil
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
	}); err != nil {
		return errors.Wrap(err, "initialization")
	}

	e.Use(SentryMiddleware())

	return nil
}

func initHandlers(e *echo.Echo, cfg Config, db postgres.Storage, receiver node.DalApi) {
	if cfg.JsonRpcConfig.Prometheus {
		e.GET("/metrics", echoprometheus.NewHandler())
	}

	rpcServer := jsonrpc.NewServer()
	handler := NewBlobHandler(receiver, db.BlobLogs, db.Namespace)

	rpcServer.Register("blob", handler)
	e.POST("", WrapHandler(rpcServer))
}

// WrapHandler wraps `http.Handler` into `echo.HandlerFunc`.
func WrapHandler(h http.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
