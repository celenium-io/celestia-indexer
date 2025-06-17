// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	indexerConfig "github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-net/go-lib/config"
)

type Config struct {
	*config.Config `yaml:",inline"`
	LogLevel       string                `validate:"omitempty,oneof=debug trace info warn error fatal panic" yaml:"log_level"`
	Indexer        indexerConfig.Indexer `validate:"required"                                                yaml:"indexer"`
	ApiConfig      ApiConfig             `validate:"required"                                                yaml:"private_api"`
}

type ApiConfig struct {
	Bind           string  `validate:"required,hostname_port" yaml:"bind"`
	RateLimit      float64 `validate:"omitempty,min=0"        yaml:"rate_limit"`
	RequestTimeout int     `validate:"omitempty,min=1"        yaml:"request_timeout"`
}
