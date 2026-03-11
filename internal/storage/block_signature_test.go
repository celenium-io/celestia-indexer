// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTableName(t *testing.T) {
	blockSignature := BlockSignature{}
	require.Equal(t, "block_signature", blockSignature.TableName())
}
