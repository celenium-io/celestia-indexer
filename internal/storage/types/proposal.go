// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

// swagger:enum ProposalStatus
/*
	ENUM(
		inactive,
		active,
		removed,
		applied,
		rejected
	)
*/
//go:generate go-enum --marshal --sql --values --names
type ProposalStatus string

func (p ProposalStatus) GreaterThan(status ProposalStatus) bool {
	switch status {
	case ProposalStatusInactive:
		return false
	case ProposalStatusRemoved:
		return true
	case ProposalStatusRejected:
		return true
	case ProposalStatusActive:
		return p == ProposalStatusInactive
	case ProposalStatusApplied:
		return true
	}
	return false
}

// swagger:enum ProposalType
/*
	ENUM(
		param_changed,
		text,
		client_update,
		community_pool_spend
	)
*/
//go:generate go-enum --marshal --sql --values --names
type ProposalType string

// swagger:enum VoteOption
/*
	ENUM(
		yes,
		no,
		no_with_veto,
		abstain
	)
*/
//go:generate go-enum --marshal --sql --values --names
type VoteOption string
