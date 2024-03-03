// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package events

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/stretchr/testify/require"
)

func Test_handleUnjail(t *testing.T) {
	tests := []struct {
		name   string
		ctx    *context.Context
		events []storage.Event
		msg    *storage.Message
		idx    *int
	}{
		{
			name: "test 1",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 844287,
					Type:   "message",
					Data: map[string]any{
						"action": "/cosmos.slashing.v1beta1.MsgUnjail",
					},
				}, {
					Height: 844287,
					Type:   "message",
					Data: map[string]any{
						"module": "slashing",
						"sender": "celestiavaloper1qe8uuf5x69c526h4nzxwv4ltftr73v7q5qhs58",
					},
				},
			},
			msg: &storage.Message{
				Type:   types.MsgUnjail,
				Height: 844287,
			},
			idx: testsuite.Ptr(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleUnjail(tt.ctx, tt.events, tt.msg, tt.idx)
			require.NoError(t, err)
		})
	}
}
