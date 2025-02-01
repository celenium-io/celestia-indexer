// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import "github.com/celenium-io/celestia-indexer/pkg/indexer/config"

type CelestialsConfig struct {
	ChainId string `validate:"required" yaml:"chain_id"`
}

type Config struct {
	*config.Config `yaml:",inline"`

	Celestials CelestialsConfig `validate:"required" yaml:"celestials"`
}
