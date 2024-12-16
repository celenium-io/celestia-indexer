// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package tvl

import (
	"context"
	"github.com/celenium-io/celestia-indexer/internal/tvl/l2beat"
	"github.com/celenium-io/celestia-indexer/internal/tvl/lama"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/parser"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/storage"
	"github.com/celenium-io/celestia-indexer/pkg/tvl/receiver"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/stopper"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Indexer struct {
	cfg       config.Config
	l2beatApi l2beat.API
	lamaApi   lama.API
	receiver  *receiver.Module
	//parser    *parser.Module
	//storage   *storage.Module
	stopper modules.Module
	log     zerolog.Logger
}

func New(ctx context.Context, cfg config.Config, stopperModule modules.Module) (Indexer, error) {
	l2beatApi, lamaApi, r, err := createReceiver(cfg)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "while creating receiver module")
	}

	//p, err := createParser(r)
	//if err != nil {
	//	return Indexer{}, errors.Wrap(err, "while creating parser module")
	//}
	//
	//s, err := createStorage(pg, cfg, p)
	//if err != nil {
	//	return Indexer{}, errors.Wrap(err, "while creating storage module")
	//}

	//err = attachStopper(stopperModule, r, p, s)
	//if err != nil {
	//	return Indexer{}, errors.Wrap(err, "while creating stopper module")
	//}

	return Indexer{
		cfg:       cfg,
		l2beatApi: l2beatApi,
		lamaApi:   lamaApi,
		receiver:  r,
		//parser:   p,
		//storage:  s,
		//stopper:  stopperModule,
		log: log.With().Str("module", "TVL indexer").Logger(),
	}, nil
}

func (i *Indexer) Start(ctx context.Context) {
	i.log.Info().Msg("starting...")

	//i.storage.Start(ctx)
	//i.parser.Start(ctx)
	i.receiver.Start(ctx)
}

func (i *Indexer) Close() error {
	i.log.Info().Msg("closing...")

	if err := i.receiver.Close(); err != nil {
		log.Err(err).Msg("closing receiver")
	}

	//if err := i.parser.Close(); err != nil {
	//	log.Err(err).Msg("closing parser")
	//}
	//if err := i.storage.Close(); err != nil {
	//	log.Err(err).Msg("closing storage")
	//}

	if err := i.stopper.Close(); err != nil {
		log.Err(err).Msg("closing stopper")
	}

	return nil
}

func createReceiver(cfg config.Config) (l2beat.API, lama.API, *receiver.Module, error) {
	l2beatDs, ok := cfg.DataSources["l2beat"]

	if !ok {
		return l2beat.API{}, lama.API{}, nil, errors.New("can't find L2Beat api datasource")
	}

	lamaDs, ok := cfg.DataSources["lama"]

	if !ok {
		return l2beat.API{}, lama.API{}, nil, errors.New("can't find DeFi Lama api datasource")
	}

	l2beatApi := l2beat.NewAPI(l2beatDs)
	lamaApi := lama.NewAPI(lamaDs)
	receiverModule := receiver.NewModule(l2beatApi, lamaApi)

	return l2beatApi, lamaApi, &receiverModule, nil
}

//func createParser(receiverModule modules.Module) (*parser.Module, error) {
//	parserModule := parser.NewModule()
//
//	if err := parserModule.AttachTo(receiverModule, receiver.DataOutput, parser.InputName); err != nil {
//		return nil, errors.Wrap(err, "while attaching parser to receiver")
//	}
//
//	return &parserModule, nil
//}
//
//func createStorage(cfg config.Config, parserModule modules.Module) (*storage.Module, error) {
//	storageModule := storage.NewModule(pg.Transactable, cfg.Indexer)
//
//	if err := storageModule.AttachTo(parserModule, parser.OutputName, storage.DataInput); err != nil {
//		return nil, errors.Wrap(err, "while attaching storage to parser")
//	}
//
//	return &storageModule, nil
//}

func attachStopper(
	stopperModule modules.Module,
	receiverModule modules.Module,
	parserModule modules.Module,
	storageModule modules.Module,
) error {
	if err := stopperModule.AttachTo(receiverModule, receiver.StopOutput, stopper.InputName); err != nil {
		return errors.Wrap(err, "while attaching stopper to receiver")
	}

	if err := stopperModule.AttachTo(parserModule, parser.StopOutput, stopper.InputName); err != nil {
		return errors.Wrap(err, "while attaching stopper to parser")
	}

	if err := stopperModule.AttachTo(storageModule, storage.StopOutput, stopper.InputName); err != nil {
		return errors.Wrap(err, "while attaching stopper to storage")
	}

	return nil
}
