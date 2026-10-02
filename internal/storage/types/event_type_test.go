// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Postgres rejects enum labels longer than NAMEDATALEN-1 bytes, which breaks database creation.
func TestEnumLabelsFitPostgres(t *testing.T) {
	const maxLabel = 63
	for _, name := range EventTypeNames() {
		require.LessOrEqual(t, len(name), maxLabel, "event_type label %q is too long: add a short label to eventTypeByChainName", name)
	}
	for _, name := range MsgTypeNames() {
		require.LessOrEqual(t, len(name), maxLabel, "msg_type label %q is too long", name)
	}
}

func TestParseChainEventType(t *testing.T) {
	tests := []struct {
		name string
		want EventType
	}{
		{
			name: "hyperlane.core.interchain_security.v1.EventRemoveRoutingIsmDomain",
			want: EventTypeHyperlanecoreismv1EventRemoveRoutingIsmDomain,
		}, {
			name: "hyperlane.core.interchain_security.v1.EventAnnounceStorageLocation",
			want: EventTypeHyperlanecoreismv1EventAnnounceStorageLocation,
		}, {
			name: "hyperlane.core.interchain_security.v1.EventCreateMessageIdMultisigIsm",
			want: EventTypeHyperlanecoreismv1EventCreateMessageIdMultisigIsm,
		}, {
			name: "hyperlane.core.interchain_security.v1.EventCreateMerkleRootMultisigIsm",
			want: EventTypeHyperlanecoreismv1EventCreateMerkleRootMultisigIsm,
		}, {
			// names that fit are stored as is
			name: "hyperlane.core.interchain_security.v1.EventSetRoutingIsmDomain",
			want: EventTypeHyperlanecoreinterchainSecurityv1EventSetRoutingIsmDomain,
		}, {
			name: "coin_spent",
			want: EventTypeCoinSpent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseChainEventType(tt.name)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}

	_, err := ParseChainEventType("hyperlane.core.interchain_security.v1.EventUnknown")
	require.Error(t, err)

	for name, typ := range eventTypeByChainName {
		require.NotEqual(t, name, typ.String(), "alias must differ from the chain name")
	}
}
