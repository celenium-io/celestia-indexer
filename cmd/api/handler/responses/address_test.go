// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package responses

import (
	"testing"

	"github.com/celestiaorg/celestia-app/v11/pkg/appconsts"
	"github.com/stretchr/testify/require"
)

// A new public network after a celestia-app bump must get a whitehat here, or is_whitehat silently stays false.
func TestWhitehatsCoverPublicNetworks(t *testing.T) {
	for _, chainId := range appconsts.PublicNetworks {
		require.NotEmpty(t, whitehats[chainId], "no whitehat for public network %s", chainId)
	}
}
