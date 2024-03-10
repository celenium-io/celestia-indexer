// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"encoding/hex"
	"strings"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/shopspring/decimal"
)

func parseValidatorUpdates(ctx *context.Context, updates []types.ValidatorUpdate) {
	for i := range updates {
		if updates[i].Power != nil {
			continue
		}
		key := updates[i].PubKey.Sum.Value.Ed25519
		consAddressBytes := types.GetConsAddressBytesFromPubKey(key)
		consAddress := strings.ToUpper(hex.EncodeToString(consAddressBytes))

		jailed := true
		ctx.AddJail(storage.Jail{
			Height: ctx.Block.Height,
			Time:   ctx.Block.Time,
			Burned: decimal.Zero,
			Validator: &storage.Validator{
				ConsAddress: consAddress,
				Stake:       decimal.Zero,
				Jailed:      &jailed,
			},
		})
	}
}
