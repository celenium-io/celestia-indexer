// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package config

import (
	"github.com/dipdup-io/celestia-indexer/internal/profiler"
	"github.com/dipdup-net/go-lib/config"
)

type Config struct {
	config.Config `yaml:",inline"`
	LogLevel      string           `validate:"omitempty,oneof=debug trace info warn error fatal panic" yaml:"log_level"`
	Indexer       Indexer          `yaml:"indexer"`
	Profiler      *profiler.Config `validate:"omitempty"                                               yaml:"profiler"`
}

type Indexer struct {
	Name         string `validate:"omitempty"       yaml:"name"`
	ThreadsCount uint32 `validate:"omitempty,min=1" yaml:"threads_count"`
	StartLevel   int64  `validate:"omitempty"       yaml:"start_level"`
	BlockPeriod  int64  `validate:"omitempty"       yaml:"block_period"`
}

// Substitute -
func (c *Config) Substitute() error {
	if err := c.Config.Substitute(); err != nil {
		return err
	}
	return nil
}
