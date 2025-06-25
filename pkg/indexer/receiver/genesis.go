// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package receiver

import (
	"context"

	"github.com/celenium-io/celestia-indexer/pkg/node/types"
	"github.com/pkg/errors"
)

func (r *Module) receiveGenesis(ctx context.Context) error {
	r.Log.Info().Msg("receiving genesis block")
	genesis, err := r.api.Genesis(ctx)
	if err != nil {
		return err
	}

	moduleAccounts, err := r.cosmosApi.ModuleAccounts(ctx)
	if err != nil {
		return errors.Wrap(err, "module account")
	}

	r.Log.Info().Msgf("got initial height of genesis block: %d", genesis.InitialHeight)
	r.MustOutput(GenesisOutput).Push(types.GenesisOutput{
		Genesis:    genesis,
		ModuleAccs: moduleAccounts,
	})
	genesisDoneInput := r.MustInput(GenesisDoneInput)

	// Wait until the genesis block will be saved
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-genesisDoneInput.Listen():
			return nil
		}
	}
}
