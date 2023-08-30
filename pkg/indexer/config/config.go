package config

import "github.com/dipdup-net/go-lib/config"

type Config struct {
	config.Config `yaml:",inline"`
	LogLevel      string  `validate:"omitempty,oneof=debug trace info warn error fatal panic" yaml:"log_level"`
	Indexer       Indexer `yaml:"indexer"`
}

type Indexer struct {
	Name         string `validate:"omitempty"       yaml:"name"`
	ThreadsCount uint32 `validate:"omitempty,min=1" yaml:"threads_count"`
}

// Substitute -
func (c *Config) Substitute() error {
	if err := c.Config.Substitute(); err != nil {
		return err
	}
	return nil
}
