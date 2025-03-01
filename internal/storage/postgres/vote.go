// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Vote -
type Vote struct {
	*postgres.Table[*storage.Vote]
}

// NewVote -
func NewVote(db *database.Bun) *Vote {
	return &Vote{
		Table: postgres.NewTable[*storage.Vote](db),
	}
}
