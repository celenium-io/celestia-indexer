// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// StakingLog -
type StakingLog struct {
	*postgres.Table[*storage.StakingLog]
}

// NewStakingLog -
func NewStakingLog(db *database.Bun) *StakingLog {
	return &StakingLog{
		Table: postgres.NewTable[*storage.StakingLog](db),
	}
}
