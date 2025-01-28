// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package celestials

import "context"

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type API interface {
	Changes(ctx context.Context, chainId string, opts ...ChangeOption) (Changes, error)
}
