// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package genesis

import (
	"slices"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
)

var divider = types.MustNumericFromString("1000000")

func parseValidatorsPower(
	ctx *decodeContext.Context,
	maxValidators int,
) {
	if ctx.Validators.Len() == 0 {
		return
	}
	sorted := slices.SortedFunc(ctx.Validators.AllValues(), func(a, b *storage.Validator) int {
		switch {
		case a.Stake.GreaterThan(b.Stake):
			return -1
		case a.Stake.LessThan(b.Stake):
			return 1
		default:
			return 0
		}
	})

	length := min(len(sorted), maxValidators)
	for _, validator := range sorted[:length] {
		power := validator.Stake.Copy().Div(divider).Floor()
		validator.Power = &power
	}
}
