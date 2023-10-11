// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// State -
type State struct {
	*postgres.Table[*storage.State]
}

// NewState -
func NewState(db *database.Bun) *State {
	return &State{
		Table: postgres.NewTable[*storage.State](db),
	}
}

// ByName -
func (s *State) ByName(ctx context.Context, name string) (state storage.State, err error) {
	err = s.DB().NewSelect().Model(&state).Where("name = ?", name).Scan(ctx)
	return
}
