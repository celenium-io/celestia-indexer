// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package ibc_relayer

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
)

type RelayerStore struct {
	metadata map[uint64]responses.Relayer
	relayers []responses.Relayer
}

func NewRelayerStore(ctx context.Context, pathToFile string, address storage.IAddress) (*RelayerStore, error) {
	file, err := readFile(pathToFile)
	if err != nil {
		return nil, err
	}

	var relayers []responses.Relayer
	if err = json.Unmarshal(file, &relayers); err != nil {
		return nil, err
	}

	if address == nil {
		return nil, errors.New("received 'nil' instead of storage.IAddress")
	}

	data := make(map[uint64]responses.Relayer)
	for _, relayer := range relayers {
		for _, addr := range relayer.Addresses {
			id, err := address.IdByAddress(ctx, addr)
			if err != nil {
				continue
			}
			data[id] = relayer
		}
	}

	rs := &RelayerStore{
		metadata: data,
		relayers: relayers,
	}

	return rs, nil
}

func (s *RelayerStore) List() map[uint64]responses.Relayer {
	if s != nil {
		return s.metadata
	}
	return map[uint64]responses.Relayer{}
}

func (s *RelayerStore) All() []responses.Relayer {
	if s != nil {
		return s.relayers
	}
	return []responses.Relayer{}
}

func readFile(path string) ([]byte, error) {
	wd, _ := os.Getwd()
	p := filepath.Join(wd, path)
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	return data, nil
}
