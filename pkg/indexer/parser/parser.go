// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package parser

import (
	"context"

	"github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
)

type Module struct {
	modules.BaseModule

	cfg config.Indexer
}

var _ modules.Module = (*Module)(nil)

const (
	InputName       = "blocks"
	OutputName      = "data"
	OutputBlobsName = "blobs"
	StopOutput      = "stop"
)

func NewModule(cfg config.Indexer) Module {
	m := Module{
		BaseModule: modules.New("parser"),
		cfg:        cfg,
	}
	m.CreateInputWithCapacity(InputName, 128)
	m.CreateOutput(OutputName)
	m.CreateOutput(OutputBlobsName)
	m.CreateOutput(StopOutput)

	return m
}

func (p *Module) Start(ctx context.Context) {
	p.Log.Info().Msg("starting parser module...")
	p.G.GoCtx(ctx, p.listen)
}

func (p *Module) Close() error {
	p.Log.Info().Msg("closing...")
	p.G.Wait()
	return nil
}
