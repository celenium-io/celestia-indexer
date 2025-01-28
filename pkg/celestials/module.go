// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package celestials

import (
	"context"
	"database/sql"
	"maps"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/celestials"
	"github.com/celenium-io/celestia-indexer/internal/celestials/api"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/pkg/errors"
)

type Module struct {
	modules.BaseModule

	celestials celestials.API
	address    IdByHash
	states     storage.ICelestialState
	tx         sdk.Transactable
	state      storage.CelestialState

	celestialsDatasource config.DataSource
	indexerName          string
	network              string
	prefix               string
	indexPeriod          time.Duration
	databaseTimeout      time.Duration
	limit                int64
}

func New(
	celestialsDatasource config.DataSource,
	address IdByHash,
	state storage.ICelestialState,
	tx sdk.Transactable,
	indexerName string,
	network string,
	opts ...ModuleOption,
) *Module {
	module := Module{
		BaseModule:      modules.New("celestials"),
		address:         address,
		states:          state,
		tx:              tx,
		celestials:      api.New(celestialsDatasource.URL),
		indexerName:     indexerName,
		network:         network,
		indexPeriod:     time.Minute,
		databaseTimeout: time.Minute,
		limit:           100,
	}

	for i := range opts {
		opts[i](&module)
	}

	return &module
}

func (m *Module) Close() error {
	m.Log.Info().Msg("closing scanner...")
	m.G.Wait()

	return nil
}

func (m *Module) Start(ctx context.Context) {
	if err := m.getState(ctx); err != nil {
		m.Log.Err(err).Msg("state receiving")
		return
	}
	m.Log.Info().Msg("starting scanner...")
	m.G.GoCtx(ctx, m.receive)
}

func (m *Module) getState(ctx context.Context) error {
	requestCtx, cancel := context.WithTimeout(ctx, m.databaseTimeout)
	defer cancel()

	state, err := m.states.ByName(requestCtx, m.indexerName)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return errors.Wrap(err, "state by name")
		}
		m.state = storage.CelestialState{
			Name:     m.indexerName,
			ChangeId: 0,
		}
		return m.states.Save(ctx, &m.state)
	}
	m.state = state
	return nil
}

func (m *Module) receive(ctx context.Context) {
	if err := m.sync(ctx); err != nil {
		m.Log.Err(err).Msg("sync")
	}

	ticker := time.NewTicker(m.indexPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := m.sync(ctx); err != nil {
				m.Log.Err(err).Msg("sync")
			}
		}
	}
}

func (m *Module) getChanges(ctx context.Context) (celestials.Changes, error) {
	requestCtx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(m.celestialsDatasource.Timeout))
	defer cancel()

	return m.celestials.Changes(
		requestCtx,
		m.network,
		celestials.WithFromChangeId(m.state.ChangeId),
		celestials.WithImages(),
		celestials.WithLimit(m.limit),
	)
}

func (m *Module) sync(ctx context.Context) error {
	var end bool

	for !end {
		changes, err := m.getChanges(ctx)
		if err != nil {
			return errors.Wrap(err, "get changes")
		}

		cids := make(map[string]storage.Celestial)

		for i := range changes.Changes {
			if m.state.ChangeId >= changes.Changes[i].ChangeID {
				continue
			}
			m.state.ChangeId = changes.Changes[i].ChangeID

			prefix, hash, err := types.Address(changes.Changes[i].Address).Decode()
			if err != nil {
				return errors.Wrapf(err, "decoding address %s", changes.Changes[i].Address)
			}
			if m.prefix != "" && prefix != m.prefix {
				return errors.Errorf("invalid address prefix %s", changes.Changes[i].Address)
			}

			addressId, err := m.address.IdByHash(ctx, hash)
			if err != nil {
				return errors.Wrap(err, "address by hash")
			}

			if len(addressId) == 0 {
				return errors.Errorf("can't find address %s", changes.Changes[i].Address)
			}

			cids[changes.Changes[i].CelestialID] = storage.Celestial{
				Id:        changes.Changes[i].CelestialID,
				ImageUrl:  changes.Changes[i].ImageURL,
				AddressId: addressId[0],
				ChangeId:  changes.Changes[i].ChangeID,
			}
		}

		if err := m.save(ctx, cids); err != nil {
			return errors.Wrap(err, "save")
		}
		end = len(changes.Changes) < int(m.limit)
	}
	return nil
}

func (m *Module) save(ctx context.Context, cids map[string]storage.Celestial) error {
	requestCtx, cancel := context.WithTimeout(ctx, m.databaseTimeout)
	defer cancel()

	tx, err := postgres.BeginCelestialTransaction(requestCtx, m.tx)
	if err != nil {
		return errors.Wrap(err, "begin transactions")
	}
	defer tx.Close(requestCtx)

	if err := tx.SaveCelestials(requestCtx, maps.Values(cids)); err != nil {
		return tx.HandleError(requestCtx, errors.Wrap(err, "save celestials"))
	}

	if err := tx.UpdateState(requestCtx, &m.state); err != nil {
		return tx.HandleError(requestCtx, errors.Wrap(err, "update state"))
	}

	return tx.Flush(requestCtx)
}
