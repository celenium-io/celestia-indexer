// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDenomMetadata_TableName(t *testing.T) {
	denom_metadata := DenomMetadata{}
	assert.Equal(t, "denom_metadata", denom_metadata.TableName())
}
