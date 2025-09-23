package postgres

import "github.com/dipdup-net/go-lib/database"

type HLIGPConfig struct {
	*database.Bun
}

func NewHLIGPConfig(conn *database.Bun) *HLIGPConfig {
	return &HLIGPConfig{conn}
}
