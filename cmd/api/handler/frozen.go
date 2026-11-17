// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"sync"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celestiaorg/celestia-app/v11/app/frozen"
	"github.com/pkg/errors"
)

// frozenState tells whether the app/frozen rules are active on the indexed chain.
type frozenState struct {
	state       storage.IState
	indexerName string

	mx      sync.Mutex
	chainId string
}

func newFrozenState(state storage.IState, indexerName string) *frozenState {
	return &frozenState{state: state, indexerName: indexerName}
}

// activeChainId returns the chain id once the rules apply, "" before that.
func (f *frozenState) activeChainId(ctx context.Context) (string, error) {
	f.mx.Lock()
	defer f.mx.Unlock()

	if f.chainId != "" {
		return f.chainId, nil
	}

	state, err := f.state.ByName(ctx, f.indexerName)
	if err != nil {
		return "", errors.Wrap(err, "failed to get state")
	}
	if state.Version < frozen.ActivationVersion {
		return "", nil
	}
	f.chainId = state.ChainId
	return f.chainId, nil
}
