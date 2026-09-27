// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package mock

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	gomock "go.uber.org/mock/gomock"
)

// TxRepos holds typed mocks for storage.TxRepos, so tests set expectations without type assertions.
type TxRepos struct {
	Address          *MockIAddress
	Blocks           *MockIBlock
	BondUpdates      *MockIValidatorBondUpdate
	Constants        *MockIConstant
	Delegation       *MockIDelegation
	HyperlaneMailbox *MockIHLMailbox
	HyperlaneToken   *MockIHLToken
	HyperlaneIgp     *MockIHLIGP
	IbcConnections   *MockIIbcConnection
	Namespace        *MockINamespace
	Proposals        *MockIProposal
	State            *MockIState
	Validators       *MockIValidator
	Votes            *MockIVote
	ZkIsm            *MockIZkISM
}

func NewTxRepos(ctrl *gomock.Controller) *TxRepos {
	return &TxRepos{
		Address:          NewMockIAddress(ctrl),
		Blocks:           NewMockIBlock(ctrl),
		BondUpdates:      NewMockIValidatorBondUpdate(ctrl),
		Constants:        NewMockIConstant(ctrl),
		Delegation:       NewMockIDelegation(ctrl),
		HyperlaneMailbox: NewMockIHLMailbox(ctrl),
		HyperlaneToken:   NewMockIHLToken(ctrl),
		HyperlaneIgp:     NewMockIHLIGP(ctrl),
		IbcConnections:   NewMockIIbcConnection(ctrl),
		Namespace:        NewMockINamespace(ctrl),
		Proposals:        NewMockIProposal(ctrl),
		State:            NewMockIState(ctrl),
		Validators:       NewMockIValidator(ctrl),
		Votes:            NewMockIVote(ctrl),
		ZkIsm:            NewMockIZkISM(ctrl),
	}
}

// Repos returns the mocks as storage.TxRepos.
func (r *TxRepos) Repos() storage.TxRepos {
	return storage.TxRepos{
		Address:          r.Address,
		Blocks:           r.Blocks,
		BondUpdates:      r.BondUpdates,
		Constants:        r.Constants,
		Delegation:       r.Delegation,
		HyperlaneMailbox: r.HyperlaneMailbox,
		HyperlaneToken:   r.HyperlaneToken,
		HyperlaneIgp:     r.HyperlaneIgp,
		IbcConnections:   r.IbcConnections,
		Namespace:        r.Namespace,
		Proposals:        r.Proposals,
		State:            r.State,
		Validators:       r.Validators,
		Votes:            r.Votes,
		ZkIsm:            r.ZkIsm,
	}
}
