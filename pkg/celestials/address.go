// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package celestials

import "context"

type IdByHash interface {
	IdByHash(ctx context.Context, hash ...[]byte) ([]uint64, error)
}
