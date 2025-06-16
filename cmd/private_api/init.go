// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/celenium-io/celestia-indexer/cmd/private_api/handler"
	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
	"github.com/dipdup-net/go-lib/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

func initEcho(cfg ApiConfig) *echo.Echo {
	e := echo.New()
	e.Validator = handler.NewCelestiaApiValidator()
	timeout := 30 * time.Second

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper: middleware.DefaultSkipper,
		Timeout: timeout,
	}))

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
	e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Pre(middleware.RemoveTrailingSlash())

	if cfg.RateLimit > 0 {
		e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
			Skipper: middleware.DefaultSkipper,
			Store:   middleware.NewRateLimiterMemoryStore(rate.Limit(cfg.RateLimit)),
		}))

	}

	e.Server.IdleTimeout = time.Second * 30

	return e
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

func initHandlers(e *echo.Echo, db postgres.Storage) {
	v1 := e.Group("v1")
	auth := v1.Group("/auth")
	{
		keyValidator := handler.NewKeyValidator(db.ApiKeys, db.BlobLogs)
		keyMiddleware := middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
			KeyLookup: "header:Authorization",
			Validator: keyValidator.Validate,
		})
		adminMiddleware := AdminMiddleware()

		rollupAuthHandler := handler.NewRollupAuthHandler(db.Rollup, db.Address, db.Namespace, db.Transactable)
		rollup := auth.Group("/rollup")
		{
			rollup.POST("/new", rollupAuthHandler.Create, keyMiddleware)
			rollup.PATCH("/:id", rollupAuthHandler.Update, keyMiddleware)
			rollup.DELETE("/:id", rollupAuthHandler.Delete, keyMiddleware, adminMiddleware)
			rollup.PATCH("/:id/verify", rollupAuthHandler.Verify, keyMiddleware, adminMiddleware)
			rollup.GET("/unverified", rollupAuthHandler.Unverified, keyMiddleware, adminMiddleware)
		}
	}
}
