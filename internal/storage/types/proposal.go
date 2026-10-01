// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

// swagger:enum ProposalStatus
/*
	ENUM(
		inactive,
		active,
		removed,
		applied,
		rejected,
		failed,
		cancelled
	)
*/
//go:generate go-enum --marshal --sql --values --names
type ProposalStatus string

// rank orders statuses by lifecycle: an empty status (vote or deposit) never wins,
// and every final status outranks inactive and active.
func (p ProposalStatus) rank() int {
	switch p {
	case ProposalStatusInactive:
		return 1
	case ProposalStatusActive:
		return 2
	case ProposalStatusRemoved, ProposalStatusApplied, ProposalStatusRejected,
		ProposalStatusFailed, ProposalStatusCancelled:
		return 3
	default:
		return 0
	}
}

// GreaterThan reports whether p is a later lifecycle stage than status.
func (p ProposalStatus) GreaterThan(status ProposalStatus) bool {
	return p.rank() > status.rank()
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

// swagger:enum VoterType
/*
	ENUM(
		address,
		validator
	)
*/
//go:generate go-enum --marshal --sql --values --names
type VoterType string
