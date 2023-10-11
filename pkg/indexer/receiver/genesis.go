// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package receiver

import "context"

func (r *Module) receiveGenesis(ctx context.Context) error {
	r.Log.Info().Msg("receiving genesis block")
	genesis, err := r.api.Genesis(ctx)
	if err != nil {
		return err
	}

	r.Log.Info().Msgf("got initial height of genesis block: %d", genesis.InitialHeight)
	r.MustOutput(GenesisOutput).Push(genesis)
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
