package main

import (
	"github.com/dipdup-io/celestia-indexer/internal/indexer"
	"github.com/dipdup-net/go-lib/config"
)

type Config struct {
	config.Config `yaml:",inline"`
	LogLevel      string         `yaml:"log_level" validate:"omitempty,oneof=debug trace info warn error fatal panic"`
	Indexer       indexer.Config `yaml:"indexer"`
}

// Substitute -
func (c *Config) Substitute() error {
	if err := c.Config.Substitute(); err != nil {
		return err
	}
	return nil
}
