// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDenomMetadata_TableName(t *testing.T) {
	denom_metadata := DenomMetadata{}
	assert.Equal(t, "denom_metadata", denom_metadata.TableName())
}
