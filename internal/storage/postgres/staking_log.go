// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// StakingLog -
type StakingLog struct {
	*postgres.Table[*storage.StakingLog]
}

// NewStakingLog -
func NewStakingLog(db bun.IDB) storage.IStakingLog {
	return &StakingLog{
		Table: postgres.NewTable[*storage.StakingLog](db),
	}
}
