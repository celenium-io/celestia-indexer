// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"sync/atomic"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celestiaorg/celestia-app/v11/app/frozen"
	"github.com/pkg/errors"
)

// FrozenState tells whether the app/frozen rules are active on the indexed chain. Shared by handlers.
type FrozenState struct {
	state       storage.IState
	indexerName string

	chainId atomic.Pointer[string]
}

func NewFrozenState(state storage.IState, indexerName string) *FrozenState {
	return &FrozenState{state: state, indexerName: indexerName}
}

// activeChainId returns the chain id once the rules apply, "" before that.
// No lock: the value only ever changes from unset to the chain id.
func (f *FrozenState) activeChainId(ctx context.Context) (string, error) {
	if chainId := f.chainId.Load(); chainId != nil {
		return *chainId, nil
	}

	state, err := f.state.ByName(ctx, f.indexerName)
	if err != nil {
		return "", errors.Wrap(err, "failed to get state")
	}
	if state.Version < frozen.ActivationVersion {
		return "", nil
	}
	chainId := state.ChainId
	f.chainId.Store(&chainId)
	return chainId, nil
}
