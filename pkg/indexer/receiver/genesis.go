package receiver

import "context"

func (r *Module) receiveGenesis(ctx context.Context) error {
	r.log.Info().Msg("receiving genesis block")
	genesis, err := r.api.Genesis(ctx)
	if err != nil {
		return err
	}

	r.log.Info().Msgf("got initial height of genesis block: %d", genesis.InitialHeight)
	r.outputs[GenesisOutput].Push(genesis)

	// Wait until genesis block will be saved
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-r.inputs[GenesisDoneInput].Listen():
			return nil
		}
	}
}
