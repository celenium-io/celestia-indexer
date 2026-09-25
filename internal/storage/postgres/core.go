// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/stats"
	models "github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/postgres/migrations"
	celestials "github.com/celenium-io/celestial-module/pkg/storage"
	celestialsPg "github.com/celenium-io/celestial-module/pkg/storage/postgres"
	"github.com/dipdup-io/go-lib/config"
	"github.com/dipdup-io/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

// Storage -
type Storage struct {
	*postgres.Storage

	cfg        config.Database
	scriptsDir string

	Blocks               models.IBlock
	BlockStats           models.IBlockStats
	BlockSignatures      models.IBlockSignature
	BlobLogs             models.IBlobLog
	Constants            models.IConstant
	DenomMetadata        models.IDenomMetadata
	Tx                   models.ITx
	Message              models.IMessage
	Event                models.IEvent
	Address              models.IAddress
	VestingAccounts      models.IVestingAccount
	VestingPeriods       models.IVestingPeriod
	Namespace            models.INamespace
	State                models.IState
	Stats                models.IStats
	Search               models.ISearch
	Validator            models.IValidator
	ValidatorBondUpdates models.IValidatorBondUpdate
	StakingLogs          models.IStakingLog
	Delegation           models.IDelegation
	Redelegation         models.IRedelegation
	Undelegation         models.IUndelegation
	Jails                models.IJail
	Rollup               models.IRollup
	RollupProvider       models.IRollupProvider
	Grants               models.IGrant
	ApiKeys              models.IApiKey
	Proposals            models.IProposal
	Votes                models.IVote
	IbcClients           models.IIbcClient
	IbcConnections       models.IIbcConnection
	IbcChannels          models.IIbcChannel
	IbcTransfers         models.IIbcTransfer
	HLMailbox            models.IHLMailbox
	HLTransfer           models.IHLTransfer
	HLToken              models.IHLToken
	HLIGP                models.IHLIGP
	HLIGPConfig          models.IHLIGPConfig
	HLGasPayment         models.IHLGasPayment
	SignalVersion        models.ISignalVersion
	Upgrade              models.IUpgrade
	Forwardings          models.IForwarding
	ZkISM                models.IZkISM
	Celestials           celestials.ICelestial
	CelestialState       celestials.ICelestialState
	Notificator          *Notificator

	export models.Export
}

// Create -
func Create(ctx context.Context, cfg config.Database, scriptsDir string, withMigrations bool) (Storage, error) {
	init := initDatabase
	if withMigrations {
		init = initDatabaseWithMigrations
	}
	strg, err := postgres.Create(ctx, cfg, init)
	if err != nil {
		return Storage{}, err
	}

	export := NewExport(strg.Connection())

	s := Storage{
		cfg:                  cfg,
		scriptsDir:           scriptsDir,
		Storage:              strg,
		Blocks:               NewBlocks(strg.Connection().DB()),
		BlockStats:           NewBlockStats(strg.Connection().DB()),
		BlockSignatures:      NewBlockSignature(strg.Connection().DB()),
		BlobLogs:             NewBlobLog(strg.Connection().DB(), export),
		Constants:            NewConstant(strg.Connection().DB()),
		DenomMetadata:        NewDenomMetadata(strg.Connection().DB()),
		Message:              NewMessage(strg.Connection().DB()),
		Event:                NewEvent(strg.Connection().DB()),
		Address:              NewAddress(strg.Connection().DB()),
		VestingAccounts:      NewVestingAccount(strg.Connection().DB()),
		VestingPeriods:       NewVestingPeriod(strg.Connection().DB()),
		Tx:                   NewTx(strg.Connection().DB()),
		State:                NewState(strg.Connection().DB()),
		Namespace:            NewNamespace(strg.Connection().DB()),
		Stats:                NewStats(strg.Connection().DB()),
		Search:               NewSearch(strg.Connection().DB()),
		Validator:            NewValidator(strg.Connection().DB()),
		ValidatorBondUpdates: NewValidatorBondUpdate(strg.Connection().DB()),
		StakingLogs:          NewStakingLog(strg.Connection().DB()),
		Delegation:           NewDelegation(strg.Connection().DB()),
		Redelegation:         NewRedelegation(strg.Connection().DB()),
		Undelegation:         NewUndelegation(strg.Connection().DB()),
		Jails:                NewJail(strg.Connection().DB()),
		Rollup:               NewRollup(strg.Connection().DB()),
		RollupProvider:       NewRollupProvider(strg.Connection().DB()),
		Grants:               NewGrant(strg.Connection().DB()),
		ApiKeys:              NewApiKey(strg.Connection().DB()),
		Proposals:            NewProposal(strg.Connection().DB()),
		Votes:                NewVote(strg.Connection().DB()),
		IbcClients:           NewIbcClient(strg.Connection().DB()),
		IbcConnections:       NewIbcConnection(strg.Connection().DB()),
		IbcChannels:          NewIbcChannel(strg.Connection().DB()),
		IbcTransfers:         NewIbcTransfer(strg.Connection().DB()),
		HLMailbox:            NewHLMailbox(strg.Connection().DB()),
		HLTransfer:           NewHLTransfer(strg.Connection().DB()),
		HLToken:              NewHLToken(strg.Connection().DB()),
		HLIGP:                NewHLIGP(strg.Connection().DB()),
		HLIGPConfig:          NewHLIGPConfig(strg.Connection().DB()),
		HLGasPayment:         NewHLGasPayment(strg.Connection().DB()),
		SignalVersion:        NewSignalVersion(strg.Connection().DB()),
		Upgrade:              NewUpgrade(strg.Connection().DB()),
		Forwardings:          NewForwarding(strg.Connection().DB()),
		ZkISM:                NewZkISM(strg.Connection().DB()),
		Celestials:           celestialsPg.NewCelestials(strg.Connection()),
		CelestialState:       celestialsPg.NewCelestialState(strg.Connection()),
		Notificator:          NewNotificator(strg.Connection().Pool()),

		export: export,
	}

	if err := s.createScripts(ctx, "functions", false); err != nil {
		return s, errors.Wrap(err, "creating views")
	}
	if err := s.createScripts(ctx, "views", true); err != nil {
		return s, errors.Wrap(err, "creating views")
	}
	return s, nil
}

func initDatabase(ctx context.Context, conn *database.Bun) error {
	if err := createExtensions(ctx, conn); err != nil {
		return errors.Wrap(err, "create extensions")
	}
	if err := createTypes(ctx, conn); err != nil {
		return errors.Wrap(err, "creating custom types")
	}

	// register many-to-many relationships
	conn.DB().RegisterModel(
		(*models.NamespaceMessage)(nil),
		(*models.MsgAddress)(nil),
		(*models.MsgValidator)(nil),
		(*models.RollupProvider)(nil),
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

func initDatabaseWithMigrations(ctx context.Context, conn *database.Bun) error {
	exists, err := checkTablesExists(ctx, conn)
	if err != nil {
		return errors.Wrap(err, "check table exists")
	}

	if exists {
		if err := migrateDatabase(ctx, conn); err != nil {
			return errors.Wrap(err, "migrate database")
		}
	}

	return initDatabase(ctx, conn)
}

func (s Storage) CreateListener() models.Listener {
	return NewNotificator(s.Notificator.pool)
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
			&models.Jail{},
			&models.StakingLog{},
			&models.Vote{},
			&models.IbcTransfer{},
			&models.HLTransfer{},
			&models.SignalVersion{},
			&models.Forwarding{},
			&models.ZkISMUpdate{},
			&models.ZkISMMessage{},
			&models.ValidatorBondUpdate{},
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

func createExtensions(ctx context.Context, conn *database.Bun) error {
	return conn.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.ExecContext(ctx, "CREATE EXTENSION IF NOT EXISTS pg_trgm;")
		return err
	})
}

func migrateDatabase(ctx context.Context, db *database.Bun) error {
	migrator := migrate.NewMigrator(db.DB(), migrations.Migrations)
	if err := migrator.Init(ctx); err != nil {
		return err
	}
	if err := migrator.Lock(ctx); err != nil {
		return err
	}
	defer func() {
		if err := migrator.Unlock(ctx); err != nil {
			log.Err(err).Msg("migrator.Unlock")
		}
	}()

	ms, err := migrator.MigrationsWithStatus(ctx)
	if err != nil {
		return errors.Wrap(err, "migrator.MigrationsWithStatus")
	}
	if len(ms.Unapplied()) == 0 {
		return nil
	}

	if _, err := migrator.Migrate(ctx); err != nil {
		return errors.Wrap(err, "migrator.Migrate")
	}
	return nil
}

func (s Storage) Close() error {
	if err := s.Storage.Close(); err != nil {
		return err
	}
	return nil
}

func checkTablesExists(ctx context.Context, db *database.Bun) (bool, error) {
	var exists bool
	err := db.DB().NewRaw(`SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE  table_schema = 'public'
		AND    table_name   = 'state'
	)`).Scan(ctx, &exists)
	return exists, err
}
