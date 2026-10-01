// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProposalStatusGreaterThan(t *testing.T) {
	type testCase struct {
		status ProposalStatus
		other  ProposalStatus
		want   bool
	}

	finals := []ProposalStatus{
		ProposalStatusRemoved,
		ProposalStatusApplied,
		ProposalStatusRejected,
		ProposalStatusFailed,
		ProposalStatusCancelled,
	}

	tests := make([]testCase, 0, 7+6*len(finals))
	tests = append(tests, []testCase{
		{ProposalStatusInactive, "", true},
		{ProposalStatusActive, "", true},
		{ProposalStatusActive, ProposalStatusInactive, true},
		{"", ProposalStatusInactive, false},
		{ProposalStatusInactive, ProposalStatusActive, false},
		{ProposalStatusActive, ProposalStatusActive, false},
		{"", "", false},
	}...)
	for _, final := range finals {
		tests = append(tests,
			testCase{final, "", true},
			testCase{final, ProposalStatusInactive, true},
			testCase{final, ProposalStatusActive, true},
			testCase{"", final, false},
			testCase{ProposalStatusInactive, final, false},
			testCase{ProposalStatusActive, final, false},
		)
	}

	for _, tt := range tests {
		t.Run(string(tt.status)+" over "+string(tt.other), func(t *testing.T) {
			require.Equal(t, tt.want, tt.status.GreaterThan(tt.other))
		})
	}
}
