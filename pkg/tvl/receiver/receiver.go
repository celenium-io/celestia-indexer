// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package receiver

import (
	"context"
	"github.com/celenium-io/celestia-indexer/internal/tvl/l2beat"
	"github.com/celenium-io/celestia-indexer/internal/tvl/lama"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
)

const (
	DataOutput = "data"
	StopOutput = "stop"
)

type Module struct {
	modules.BaseModule
	apiL2Beat l2beat.IApi
	apiLama   lama.IApi
}

func NewModule(l2beatApi l2beat.API, lamaApi lama.API) Module {
	receiver := Module{
		BaseModule: modules.New("receiver"),
		apiL2Beat:  l2beatApi,
		apiLama:    lamaApi,
	}

	receiver.CreateOutput(DataOutput)
	receiver.CreateOutput(StopOutput)

	return receiver
}

func (r *Module) Start(ctx context.Context) {
	r.Log.Info().Msg("starting TVL receiver...")
	//r.G.GoCtx(ctx, r.receive)
}

func (r *Module) stopAll() {
	r.MustOutput(StopOutput).Push(struct{}{})
}

var _ modules.Module = (*Module)(nil)
