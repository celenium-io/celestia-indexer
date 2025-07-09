// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package gas

import (
	"context"
	"io"
)

//go:generate mockgen -source=$GOFILE -destination=mock.go -package=gas -typed
type ITracker interface {
	io.Closer

	Start(ctx context.Context)
	Init(ctx context.Context) error
	SubscribeOnCompute(handler ComputeHandler)
	State() GasPrice
}
