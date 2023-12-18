// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/stats"
	models "github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

// Storage -
type Storage struct {
	*postgres.Storage

	cfg        config.Database
	scriptsDir string

	Blocks        models.IBlock
	BlockStats    models.IBlockStats
	BlobLogs      models.IBlobLog
	Constants     models.IConstant
	DenomMetadata models.IDenomMetadata
	Tx            models.ITx
	Message       models.IMessage
	Event         models.IEvent
	Address       models.IAddress
	Namespace     models.INamespace
	State         models.IState
	Stats         models.IStats
	Validator     models.IValidator
	Notificator   *Notificator
}

// Create -
func Create(ctx context.Context, cfg config.Database, scriptsDir string) (Storage, error) {
	strg, err := postgres.Create(ctx, cfg, initDatabase)
	if err != nil {
		return Storage{}, err
	}

	s := Storage{
		cfg:           cfg,
		scriptsDir:    scriptsDir,
		Storage:       strg,
		Blocks:        NewBlocks(strg.Connection()),
		BlockStats:    NewBlockStats(strg.Connection()),
		BlobLogs:      NewBlobLog(strg.Connection()),
		Constants:     NewConstant(strg.Connection()),
		DenomMetadata: NewDenomMetadata(strg.Connection()),
		Message:       NewMessage(strg.Connection()),
		Event:         NewEvent(strg.Connection()),
		Address:       NewAddress(strg.Connection()),
		Tx:            NewTx(strg.Connection()),
		State:         NewState(strg.Connection()),
		Namespace:     NewNamespace(strg.Connection()),
		Stats:         NewStats(strg.Connection()),
		Validator:     NewValidator(strg.Connection()),
		Notificator:   NewNotificator(cfg, strg.Connection().DB()),
	}

	if err := s.createScripts(ctx, strg.Connection(), "functions", false); err != nil {
		return s, errors.Wrap(err, "creating views")
	}
	if err := s.createScripts(ctx, strg.Connection(), "views", true); err != nil {
		return s, errors.Wrap(err, "creating views")
	}

	return s, nil
}

func initDatabase(ctx context.Context, conn *database.Bun) error {
	if err := createTypes(ctx, conn); err != nil {
		return errors.Wrap(err, "creating custom types")
	}

	// register many-to-many relationships
	conn.DB().RegisterModel(
		(*models.NamespaceMessage)(nil),
		(*models.Signer)(nil),
		(*models.MsgAddress)(nil),
	)

	if err := database.CreateTables(ctx, conn, models.Models...); err != nil {
		if err := conn.Close(); err != nil {
			return err
		}
		return err
	}

	if err := database.MakeComments(ctx, conn, models.Models...); err != nil {
		if err := conn.Close(); err != nil {
			return err
		}
		return errors.Wrap(err, "make comments")
	}

	if err := createHypertables(ctx, conn); err != nil {
		if err := conn.Close(); err != nil {
			return err
		}
		return errors.Wrap(err, "create hypertables")
	}

	return createIndices(ctx, conn)
}

func (s Storage) CreateListener() models.Listener {
	return NewNotificator(s.cfg, s.Notificator.db)
}

func createHypertables(ctx context.Context, conn *database.Bun) error {
	return conn.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		for _, model := range []storage.Model{
			&models.Block{},
			&models.BlockStats{},
			&models.Tx{},
			&models.Message{},
			&models.Event{},
			&models.NamespaceMessage{},
			&models.BlobLog{},
		} {
			if _, err := tx.ExecContext(ctx,
				`SELECT create_hypertable(?, 'time', chunk_time_interval => INTERVAL '1 month', if_not_exists => TRUE);`,
				model.TableName(),
			); err != nil {
				return err
			}

			if err := stats.InitModel(model); err != nil {
				return err
			}
		}

		if err := stats.InitModel(&models.Validator{}); err != nil {
			return err
		}
		return nil
	})
}
