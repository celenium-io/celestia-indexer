// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package celestials

import "time"

type ModuleOption func(*Module)

func WithAddressPrefix(prefix string) ModuleOption {
	return func(m *Module) {
		m.prefix = prefix
	}
}

func WithIndexPeriod(period time.Duration) ModuleOption {
	return func(m *Module) {
		m.indexPeriod = period
	}
}

func WithDatabaseTimeout(timeout time.Duration) ModuleOption {
	return func(m *Module) {
		m.databaseTimeout = timeout
	}
}

func WithLimit(limit int64) ModuleOption {
	return func(m *Module) {
		m.limit = limit
	}
}
