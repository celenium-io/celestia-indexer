// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package lama

import "context"

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IApi interface {
	TVL(ctx context.Context, arguments *TVLArgs) (result []TVLResponse, err error)
}
