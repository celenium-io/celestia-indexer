package receiver

import "context"

func (r *Receiver) receiveGenesis(ctx context.Context) error {
	genesis, err := r.api.Genesis(ctx)
	if err != nil {
		return err
	}

	r.log.Info().Msgf("got initial height of genesis block: %d", genesis.InitialHeight)
	return nil
}
