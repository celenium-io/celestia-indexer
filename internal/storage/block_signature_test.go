// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTableName(t *testing.T) {
	blockSignature := BlockSignature{}
	assert.Equal(t, "block_signature", blockSignature.TableName())
}
