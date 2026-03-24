// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"testing"

	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/testutil"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
)

func mustMsgToMap(t *testing.T, msg cosmosTypes.Msg) storageTypes.PackedBytes {
	t.Helper()
	return testutil.MustMsgToMap(t, msg)
}
