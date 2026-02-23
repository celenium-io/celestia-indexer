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
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

// Storage -
type Storage struct {
	*postgres.Storage

	cfg        config.Database
	scriptsDir string

	Blocks          models.IBlock
	BlockStats      models.IBlockStats
	BlockSignatures models.IBlockSignature
	BlobLogs        models.IBlobLog
	Constants       models.IConstant
	DenomMetadata   models.IDenomMetadata
	Tx              models.ITx
	Message         models.IMessage
	Event           models.IEvent
	Address         models.IAddress
	VestingAccounts models.IVestingAccount
	VestingPeriods  models.IVestingPeriod
	Namespace       models.INamespace
	State           models.IState
	Stats           models.IStats
	Search          models.ISearch
	Validator       models.IValidator
	StakingLogs     models.IStakingLog
	Delegation      models.IDelegation
	Redelegation    models.IRedelegation
	Undelegation    models.IUndelegation
	Jails           models.IJail
	Rollup          models.IRollup
	RollupProvider  models.IRollupProvider
	Grants          models.IGrant
	ApiKeys         models.IApiKey
	Proposals       models.IProposal
	Votes           models.IVote
	IbcClients      models.IIbcClient
	IbcConnections  models.IIbcConnection
	IbcChannels     models.IIbcChannel
	IbcTransfers    models.IIbcTransfer
	HLMailbox       models.IHLMailbox
	HLTransfer      models.IHLTransfer
	HLToken         models.IHLToken
	HLIGP           models.IHLIGP
	HLIGPConfig     models.IHLIGPConfig
	HLGasPayment    models.IHLGasPayment
	SignalVersion   models.ISignalVersion
	Upgrade         models.IUpgrade
	Forwardings     models.IForwarding
	ZkISM           models.IZkISM
	Celestials      celestials.ICelestial
	CelestialState  celestials.ICelestialState
	Notificator     *Notificator

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

	export := NewExport(cfg)

	s := Storage{
		cfg:             cfg,
		scriptsDir:      scriptsDir,
		Storage:         strg,
		Blocks:          NewBlocks(strg.Connection()),
		BlockStats:      NewBlockStats(strg.Connection()),
		BlockSignatures: NewBlockSignature(strg.Connection()),
		BlobLogs:        NewBlobLog(strg.Connection(), export),
		Constants:       NewConstant(strg.Connection()),
		DenomMetadata:   NewDenomMetadata(strg.Connection()),
		Message:         NewMessage(strg.Connection()),
		Event:           NewEvent(strg.Connection()),
		Address:         NewAddress(strg.Connection()),
		VestingAccounts: NewVestingAccount(strg.Connection()),
		VestingPeriods:  NewVestingPeriod(strg.Connection()),
		Tx:              NewTx(strg.Connection()),
		State:           NewState(strg.Connection()),
		Namespace:       NewNamespace(strg.Connection()),
		Stats:           NewStats(strg.Connection()),
		Search:          NewSearch(strg.Connection()),
		Validator:       NewValidator(strg.Connection()),
		StakingLogs:     NewStakingLog(strg.Connection()),
		Delegation:      NewDelegation(strg.Connection()),
		Redelegation:    NewRedelegation(strg.Connection()),
		Undelegation:    NewUndelegation(strg.Connection()),
		Jails:           NewJail(strg.Connection()),
		Rollup:          NewRollup(strg.Connection()),
		RollupProvider:  NewRollupProvider(strg.Connection()),
		Grants:          NewGrant(strg.Connection()),
		ApiKeys:         NewApiKey(strg.Connection()),
		Proposals:       NewProposal(strg.Connection()),
		Votes:           NewVote(strg.Connection()),
		IbcClients:      NewIbcClient(strg.Connection()),
		IbcConnections:  NewIbcConnection(strg.Connection()),
		IbcChannels:     NewIbcChannel(strg.Connection()),
		IbcTransfers:    NewIbcTransfer(strg.Connection()),
		HLMailbox:       NewHLMailbox(strg.Connection()),
		HLTransfer:      NewHLTransfer(strg.Connection()),
		HLToken:         NewHLToken(strg.Connection()),
		HLIGP:           NewHLIGP(strg.Connection()),
		HLIGPConfig:     NewHLIGPConfig(strg.Connection()),
		HLGasPayment:    NewHLGasPayment(strg.Connection()),
		SignalVersion:   NewSignalVersion(strg.Connection()),
		Upgrade:         NewUpgrade(strg.Connection()),
		Forwardings:     NewForwarding(strg.Connection()),
		ZkISM:           NewZkISM(strg.Connection()),
		Celestials:      celestialsPg.NewCelestials(strg.Connection()),
		CelestialState:  celestialsPg.NewCelestialState(strg.Connection()),
		Notificator:     NewNotificator(cfg, strg.Connection().DB()),

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
			&models.Jail{},
			&models.StakingLog{},
			&models.Vote{},
			&models.IbcTransfer{},
			&models.HLTransfer{},
			&models.SignalVersion{},
			&models.Forwarding{},
			&models.ZkISMUpdate{},
			&models.ZkISMMessage{},
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
	defer migrator.Unlock(ctx) //nolint:errcheck

	_, err := migrator.Migrate(ctx)
	return err
}

func (s Storage) Close() error {
	if err := s.export.Close(); err != nil {
		return err
	}
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
