// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upFillNetworkMinGasPrice, downFillNetworkMinGasPrice)
}

// LegacyDec form of appconsts.DefaultNetworkMinGasPrice, set by CIP-006 in app version 2
const cip6NetworkMinGasPrice = "0.000001000000000000"

// Chains launched before app version 2 have no minfee in genesis, so the constant was saved empty.
func upFillNetworkMinGasPrice(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `UPDATE constant
		SET value = ?
		WHERE module = 'minfee' AND name = 'network_min_gas_price' AND value = ''`,
		cip6NetworkMinGasPrice,
	)
	return errors.Wrap(err, "fill network min gas price")
}

func downFillNetworkMinGasPrice(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `UPDATE constant
		SET value = ''
		WHERE module = 'minfee' AND name = 'network_min_gas_price' AND value = ?`,
		cip6NetworkMinGasPrice,
	)
	return errors.Wrap(err, "clear network min gas price")
}
