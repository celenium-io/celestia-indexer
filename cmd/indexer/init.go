package main

import (
	"os"
	"strconv"

	"github.com/dipdup-io/celestia-indexer/internal/profiler"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	goLibConfig "github.com/dipdup-net/go-lib/config"
	"github.com/pyroscope-io/client/pyroscope"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	})
}

func initConfig() (*config.Config, error) {
	configPath := rootCmd.PersistentFlags().StringP("config", "c", "dipdup.yml", "path to YAML config file")
	if err := rootCmd.Execute(); err != nil {
		log.Panic().Err(err).Msg("command line execute")
		return nil, err
	}

	if err := rootCmd.MarkFlagRequired("config"); err != nil {
		log.Panic().Err(err).Msg("config command line arg is required")
		return nil, err
	}

	var cfg config.Config
	if err := goLibConfig.Parse(*configPath, &cfg); err != nil {
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
	prscp, err = profiler.New(cfg, "indexer")
	return
}
