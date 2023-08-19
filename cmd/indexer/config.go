package main

import (
	"github.com/dipdup-net/go-lib/config"
)

type Config struct {
	config.Config `yaml:",inline"`
	LogLevel      string  `yaml:"log_level" validate:"omitempty,oneof=debug trace info warn error fatal panic"`
	Indexer       Indexer `yaml:"indexer"`
}

type Indexer struct {
	Name    string `yaml:"name" validate:"omitempty"`
	Timeout int    `yaml:"timeout" validate:"omitempty"`
	Node    *Node  `yaml:"node" validate:"omitempty"`
}

type Node struct {
	Url string `yaml:"url" validate:"omitempty,url"`
	Rps int    `yaml:"requests_per_second" validate:"omitempty,min=1"`
}

// Substitute -
func (c *Config) Substitute() error {
	if err := c.Config.Substitute(); err != nil {
		return err
	}
	return nil
}
