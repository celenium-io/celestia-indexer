// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	indexerCfg "github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/dipdup-io/go-lib/config"
	"github.com/dipdup-io/go-lib/testhelpers"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/suite"
)

const testIndexerName = "test_indexer"

// ModuleTestSuite -
type ModuleTestSuite struct {
	suite.Suite
	psqlContainer *testhelpers.PostgreSQLContainer
	storage       postgres.Storage
	fixtures      *testfixtures.Loader
}

// SetupSuite -
func (s *ModuleTestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 180*time.Second)
	defer ctxCancel()

	psqlContainer, err := testhelpers.NewPostgreSQLContainer(ctx, testhelpers.PostgreSQLContainerConfig{
		User:     "user",
		Password: "password",
		Database: "db_test",
		Port:     5432,
		Image:    "timescale/timescaledb-ha:pg15.8-ts2.17.0-all",
	})
	s.Require().NoError(err)
	s.psqlContainer = psqlContainer
	s.T().Cleanup(func() {
		// not t.Context(): it is canceled before cleanup runs
		ctx, ctxCancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer ctxCancel()
		s.Require().NoError(s.psqlContainer.Terminate(ctx))
	})

	strg, err := postgres.Create(ctx, config.Database{
		Kind:     config.DBKindPostgres,
		User:     s.psqlContainer.Config.User,
		Database: s.psqlContainer.Config.Database,
		Password: s.psqlContainer.Config.Password,
		Host:     s.psqlContainer.Config.Host,
		Port:     s.psqlContainer.MappedPort().Int(),
	}, "../../../database", false)
	s.Require().NoError(err)
	s.storage = strg
	s.T().Cleanup(func() {
		s.Require().NoError(s.storage.Close())
	})

	fixtures, err := testfixtures.New(
		testfixtures.Database(strg.Connection().DB().DB),
		testfixtures.Dialect("timescaledb"),
		testfixtures.Directory("../../../test/data"),
		testfixtures.UseAlterConstraint(),
	)
	s.Require().NoError(err)
	s.fixtures = fixtures
}

// SetupTest reloads fixtures: tests move the indexer state and add upgrades
func (s *ModuleTestSuite) SetupTest() {
	s.Require().NoError(s.fixtures.Load())
}

func (s *ModuleTestSuite) TestBlockLast() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 10*time.Second)
	defer ctxCancel()

	module := NewModule(s.storage.Transactable, s.storage.Notificator, indexerCfg.Indexer{Name: testIndexerName})
	module.Start(ctx)

	hash, err := hex.DecodeString("F44BC94BF7D064ADF82618F2691D2353161DE232ECB3091B7E5C89B453C79456")
	s.Require().NoError(err)

	dCtx := decodeContext.NewContext()
	dCtx.Block = &storage.Block{
		Height:          1001,
		Hash:            hash,
		VersionBlock:    11,
		VersionApp:      1,
		ProposerAddress: "81A24EE534DEFE1557A4C7C437E8E8FBC2F834E8",
		Time:            time.Date(2023, 7, 4, 3, 11, 26, 0, time.UTC),
		MessageTypes:    types.NewMsgTypeBitMask(),
	}

	module.MustInput(InputName).Push(dCtx)
	time.Sleep(time.Second)

	block, err := s.storage.Blocks.Last(ctx)
	s.Require().NoError(err)
	s.Require().EqualValues(1001, block.Height)
	s.Require().EqualValues(1, block.VersionApp)
	s.Require().EqualValues(11, block.VersionBlock)
	s.Require().Equal(hash, block.Hash.Bytes())

	state, err := s.storage.State.ByName(ctx, testIndexerName)
	s.Require().NoError(err)
	s.Require().Equal(testIndexerName, state.Name)
	s.Require().EqualValues(1001, state.LastHeight)

	s.Require().NoError(module.Close())
}

func (s *ModuleTestSuite) TestSignalInUpgradeBlockKeepsAppliedUpgrade() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 10*time.Second)
	defer ctxCancel()

	db := s.storage.Connection().DB()
	// above applied_at_level of the fixture upgrades, so the new signal is counted by the tally
	_, err := db.NewUpdate().Table("state").
		Set("version = 2").
		Set("last_height = 2000").
		Where("name = ?", testIndexerName).
		Exec(ctx)
	s.Require().NoError(err)

	// upgrade to v3 reached quorum at the old version
	tx, err := postgres.BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)
	s.Require().NoError(tx.SaveUpgrades(ctx, &storage.Upgrade{
		Version:      3,
		Height:       900,
		Time:         time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
		VotingPower:  types.NumericFromInt64(2),
		VotedPower:   types.NumericFromInt64(2),
		SignalsCount: 2,
		Status:       types.UpgradeStatusWaitingUpgrade,
	}))
	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	state, err := s.storage.State.ByName(ctx, testIndexerName)
	s.Require().NoError(err)

	module := NewModule(s.storage.Transactable, s.storage.Notificator, indexerCfg.Indexer{Name: testIndexerName})
	module.Start(ctx)
	defer func() {
		ctxCancel()
		s.Require().NoError(module.Close())
	}()

	height := state.LastHeight + 1
	blockTime := state.LastTime.Add(time.Minute)

	// first block at v3: x/signal accepts a signal for the current version, but it must not touch the applied upgrade
	dCtx := decodeContext.NewContext()
	dCtx.Block = &storage.Block{
		Height:          height,
		Hash:            []byte{0x01, 0x02, 0x03},
		VersionBlock:    11,
		VersionApp:      3,
		ProposerAddress: "81A24EE534DEFE1557A4C7C437E8E8FBC2F834E8",
		Time:            blockTime,
		MessageTypes:    types.NewMsgTypeBitMask(),
	}
	validator := storage.EmptyValidator()
	validator.Address = "celestiavaloper17vmk8m246t648hpmde2q7kp4ft9uwrayy09dmw"
	validator.Version = 3
	dCtx.AddValidator(validator)
	dCtx.AddSignal(&storage.SignalVersion{
		Height:    height,
		Time:      blockTime,
		Version:   3,
		Validator: &validator,
		TxId:      1,
		MsgId:     1,
	})
	dCtx.AddUpgrade(storage.Upgrade{Version: 3, SignalsCount: 1, Height: height, Time: blockTime})

	_, err = module.saveBlock(ctx, dCtx)
	s.Require().NoError(err)

	var upgrade storage.Upgrade
	err = db.NewSelect().Model(&upgrade).Where("version = 3").Scan(ctx)
	s.Require().NoError(err)
	s.Require().Equal(types.UpgradeStatusApplied, upgrade.Status)
	s.Require().EqualValues(height, upgrade.AppliedAtLevel)
	s.Require().Equal("2", upgrade.VotedPower.String())
	s.Require().Equal("2", upgrade.VotingPower.String())
	s.Require().EqualValues(2, upgrade.SignalsCount)

	signals, err := s.storage.SignalVersion.List(ctx, storage.ListSignalsFilter{Limit: 10, Version: 3})
	s.Require().NoError(err)
	s.Require().Len(signals, 1)
}

func (s *ModuleTestSuite) TestWithdrawnSignalLowersQuorum() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 10*time.Second)
	defer ctxCancel()

	db := s.storage.Connection().DB()
	_, err := db.NewUpdate().Table("state").
		Set("version = 3").
		Set("last_height = 3000").
		Where("name = ?", testIndexerName).
		Exec(ctx)
	s.Require().NoError(err)

	// both fixture validators (power 1 each) voted for v4: quorum 2 of 2
	tx, err := postgres.BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)
	signalTime := time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC)
	s.Require().NoError(storage.Insert(ctx, tx,
		&storage.SignalVersion{Height: 2990, ValidatorId: 1, Version: 4, Time: signalTime, TxId: 1, MsgId: 1},
		&storage.SignalVersion{Height: 2991, ValidatorId: 2, Version: 4, Time: signalTime, TxId: 2, MsgId: 2},
	))
	s.Require().NoError(tx.SaveUpgrades(ctx, &storage.Upgrade{
		Version:      4,
		Height:       2990,
		Time:         signalTime,
		VotingPower:  types.NumericFromInt64(2),
		VotedPower:   types.NumericFromInt64(2),
		SignalsCount: 2,
		Status:       types.UpgradeStatusWaitingUpgrade,
	}))
	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	state, err := s.storage.State.ByName(ctx, testIndexerName)
	s.Require().NoError(err)

	module := NewModule(s.storage.Transactable, s.storage.Notificator, indexerCfg.Indexer{Name: testIndexerName})
	module.Start(ctx)
	defer func() {
		ctxCancel()
		s.Require().NoError(module.Close())
	}()

	height := state.LastHeight + 1
	blockTime := state.LastTime.Add(time.Minute)

	// validator 2 withdraws its vote by signaling the current version
	dCtx := decodeContext.NewContext()
	dCtx.Block = &storage.Block{
		Height:          height,
		Hash:            []byte{0x04, 0x05, 0x06},
		VersionBlock:    11,
		VersionApp:      3,
		ProposerAddress: "81A24EE534DEFE1557A4C7C437E8E8FBC2F834E8",
		Time:            blockTime,
		MessageTypes:    types.NewMsgTypeBitMask(),
	}
	validator := storage.EmptyValidator()
	validator.Address = "celestiavaloper189ecvq5avj0wehrcfnagpd5sd8pup9aqmdglmr"
	validator.Version = 3
	dCtx.AddValidator(validator)
	dCtx.AddSignal(&storage.SignalVersion{
		Height:    height,
		Time:      blockTime,
		Version:   3,
		Validator: &validator,
		TxId:      3,
		MsgId:     3,
	})
	dCtx.AddUpgrade(storage.Upgrade{Version: 3, SignalsCount: 1, Height: height, Time: blockTime})

	_, err = module.saveBlock(ctx, dCtx)
	s.Require().NoError(err)

	var upgrade storage.Upgrade
	err = db.NewSelect().Model(&upgrade).Where("version = 4").Scan(ctx)
	s.Require().NoError(err)
	s.Require().Equal("1", upgrade.VotedPower.String())
	s.Require().Equal("2", upgrade.VotingPower.String())
	s.Require().Equal(types.UpgradeStatusProcessing, upgrade.Status)
	s.Require().EqualValues(2, upgrade.SignalsCount)
}

func (s *ModuleTestSuite) TestAppliedUpgradeClosesRound() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 10*time.Second)
	defer ctxCancel()

	db := s.storage.Connection().DB()
	_, err := db.NewUpdate().Table("state").
		Set("version = 4").
		Set("last_height = 4000").
		Where("name = ?", testIndexerName).
		Exec(ctx)
	s.Require().NoError(err)

	// the round for v5 was never closed by an indexed MsgTryUpgrade (e.g. it was sent via authz)
	tx, err := postgres.BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)
	signalTime := time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC)
	s.Require().NoError(storage.Insert(ctx, tx,
		&storage.SignalVersion{Height: 3990, ValidatorId: 1, Version: 5, Time: signalTime, TxId: 1, MsgId: 1},
		&storage.SignalVersion{Height: 3991, ValidatorId: 2, Version: 6, Time: signalTime, TxId: 2, MsgId: 2},
	))
	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	state, err := s.storage.State.ByName(ctx, testIndexerName)
	s.Require().NoError(err)

	module := NewModule(s.storage.Transactable, s.storage.Notificator, indexerCfg.Indexer{Name: testIndexerName})
	module.Start(ctx)
	defer func() {
		ctxCancel()
		s.Require().NoError(module.Close())
	}()

	dCtx := decodeContext.NewContext()
	dCtx.Block = &storage.Block{
		Height:          state.LastHeight + 1,
		Hash:            []byte{0x07, 0x08, 0x09},
		VersionBlock:    11,
		VersionApp:      5,
		ProposerAddress: "81A24EE534DEFE1557A4C7C437E8E8FBC2F834E8",
		Time:            state.LastTime.Add(time.Minute),
		MessageTypes:    types.NewMsgTypeBitMask(),
	}
	_, err = module.saveBlock(ctx, dCtx)
	s.Require().NoError(err)

	signals, err := s.storage.SignalVersion.List(ctx, storage.ListSignalsFilter{Limit: 10, Sort: sdk.SortOrderAsc, From: signalTime})
	s.Require().NoError(err)
	var counted, notCounted *storage.SignalVersion
	for i := range signals {
		switch signals[i].Height {
		case 3990:
			counted = &signals[i]
		case 3991:
			notCounted = &signals[i]
		}
	}
	s.Require().NotNil(counted)
	s.Require().NotNil(notCounted)
	s.Require().Equal("1", counted.VotingPower.String())
	s.Require().Equal("0", notCounted.VotingPower.String())
}

func (s *ModuleTestSuite) TestTryUpgradeClosesRound() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 10*time.Second)
	defer ctxCancel()

	db := s.storage.Connection().DB()
	_, err := db.NewUpdate().Table("state").
		Set("version = 3").
		Set("last_height = 3000").
		Where("name = ?", testIndexerName).
		Exec(ctx)
	s.Require().NoError(err)

	// both fixture validators (power 1 each) voted: 1 for v4 and later switched to v5, 2 for v4
	tx, err := postgres.BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)
	signalTime := time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC)
	s.Require().NoError(storage.Insert(ctx, tx,
		&storage.SignalVersion{Height: 2980, ValidatorId: 1, Version: 4, Time: signalTime, TxId: 1, MsgId: 1},
		&storage.SignalVersion{Height: 2990, ValidatorId: 2, Version: 4, Time: signalTime, TxId: 2, MsgId: 2},
		&storage.SignalVersion{Height: 2995, ValidatorId: 1, Version: 5, Time: signalTime, TxId: 3, MsgId: 3},
	))
	_, err = tx.Tx().NewUpdate().Table("validator").Set("version = 5").Where("id = 1").Exec(ctx)
	s.Require().NoError(err)
	_, err = tx.Tx().NewUpdate().Table("validator").Set("version = 4").Where("id = 2").Exec(ctx)
	s.Require().NoError(err)
	// validator 2 left the active set, so v5 has the whole power of 1
	_, err = tx.Tx().NewUpdate().Table("validator").Set("power = 0").Where("id = 2").Exec(ctx)
	s.Require().NoError(err)
	s.Require().NoError(tx.SaveUpgrades(ctx,
		&storage.Upgrade{Version: 4, Height: 2980, Time: signalTime, SignalsCount: 2},
		&storage.Upgrade{Version: 5, Height: 2995, Time: signalTime, SignalsCount: 1},
	))
	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	state, err := s.storage.State.ByName(ctx, testIndexerName)
	s.Require().NoError(err)

	module := NewModule(s.storage.Transactable, s.storage.Notificator, indexerCfg.Indexer{Name: testIndexerName})
	module.Start(ctx)
	defer func() {
		ctxCancel()
		s.Require().NoError(module.Close())
	}()

	const signer = "celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8"
	height := state.LastHeight + 1
	blockTime := state.LastTime.Add(time.Minute)

	dCtx := decodeContext.NewContext()
	dCtx.Block = &storage.Block{
		Height:          height,
		Hash:            []byte{0x0a, 0x0b, 0x0c},
		VersionBlock:    11,
		VersionApp:      3,
		ChainId:         "celestia",
		ProposerAddress: "81A24EE534DEFE1557A4C7C437E8E8FBC2F834E8",
		Time:            blockTime,
		MessageTypes:    types.NewMsgTypeBitMask(),
	}
	signerAddress := &storage.Address{Address: signer, Height: height, LastHeight: height}
	s.Require().NoError(dCtx.AddAddress(signerAddress))
	dCtx.TryUpgrade = &storage.Upgrade{
		Height:    height,
		Time:      blockTime,
		EndTime:   blockTime,
		EndHeight: height,
		Signer:    signerAddress,
		TxId:      7,
		MsgId:     7,
	}

	_, err = module.saveBlock(ctx, dCtx)
	s.Require().NoError(err)

	upgrade, err := s.storage.Upgrade.ByVersion(ctx, 5)
	s.Require().NoError(err)
	s.Require().Equal(types.UpgradeStatusWaitingUpgrade, upgrade.Status)
	s.Require().EqualValues(height, upgrade.EndHeight)
	s.Require().EqualValues(height+232_616, upgrade.ExpectedHeight)
	s.Require().NotZero(upgrade.SignerId)
	s.Require().NotNil(upgrade.Signer)
	s.Require().Equal(signer, upgrade.Signer.Address)
	s.Require().Equal("1", upgrade.VotedPower.String())

	signals, err := s.storage.SignalVersion.List(ctx, storage.ListSignalsFilter{Limit: 10, Sort: sdk.SortOrderAsc, From: signalTime})
	s.Require().NoError(err)
	power := make(map[int64]string)
	for i := range signals {
		if signals[i].Height >= 2980 {
			power[int64(signals[i].Height)] = signals[i].VotingPower.String()
		}
	}
	// counted: validator 1 for v5; not counted: its old v4 vote and the unbonded validator 2
	s.Require().Equal(map[int64]string{2980: "0", 2990: "0", 2995: "1"}, power)
}

func (s *ModuleTestSuite) TestPowerChangeRecountsWithoutSignals() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 10*time.Second)
	defer ctxCancel()

	db := s.storage.Connection().DB()
	_, err := db.NewUpdate().Table("state").
		Set("version = 3").
		Set("last_height = 3000").
		Where("name = ?", testIndexerName).
		Exec(ctx)
	s.Require().NoError(err)

	// validator 1 voted for v4: 1 of 2, no quorum
	tx, err := postgres.BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)
	signalTime := time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC)
	s.Require().NoError(storage.Insert(ctx, tx,
		&storage.SignalVersion{Height: 2990, ValidatorId: 1, Version: 4, Time: signalTime, TxId: 1, MsgId: 1},
	))
	s.Require().NoError(tx.SaveUpgrades(ctx, &storage.Upgrade{
		Version:      4,
		Height:       2990,
		Time:         signalTime,
		VotingPower:  types.NumericFromInt64(2),
		VotedPower:   types.NumericFromInt64(1),
		SignalsCount: 1,
	}))
	// validator 2 leaves the active set before the next block
	_, err = tx.Tx().NewUpdate().Table("validator").Set("power = 0").Where("id = 2").Exec(ctx)
	s.Require().NoError(err)
	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	state, err := s.storage.State.ByName(ctx, testIndexerName)
	s.Require().NoError(err)

	module := NewModule(s.storage.Transactable, s.storage.Notificator, indexerCfg.Indexer{Name: testIndexerName})
	module.Start(ctx)
	defer func() {
		ctxCancel()
		s.Require().NoError(module.Close())
	}()

	// a block without signals
	dCtx := decodeContext.NewContext()
	dCtx.Block = &storage.Block{
		Height:          state.LastHeight + 1,
		Hash:            []byte{0x0d, 0x0e, 0x0f},
		VersionBlock:    11,
		VersionApp:      3,
		ProposerAddress: "81A24EE534DEFE1557A4C7C437E8E8FBC2F834E8",
		Time:            state.LastTime.Add(time.Minute),
		MessageTypes:    types.NewMsgTypeBitMask(),
	}
	_, err = module.saveBlock(ctx, dCtx)
	s.Require().NoError(err)

	// total power dropped to 1, so validator 1 alone reaches the quorum
	upgrade, err := s.storage.Upgrade.ByVersion(ctx, 4)
	s.Require().NoError(err)
	s.Require().Equal("1", upgrade.VotingPower.String())
	s.Require().Equal("1", upgrade.VotedPower.String())
	s.Require().Equal(types.UpgradeStatusWaitingUpgrade, upgrade.Status)
	s.Require().EqualValues(1, upgrade.SignalsCount)
}

func TestSuiteModule_Run(t *testing.T) {
	suite.Run(t, new(ModuleTestSuite))
}
