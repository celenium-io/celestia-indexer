// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConstant_TableName(t *testing.T) {
	constant := Constant{}
	assert.Equal(t, "constant", constant.TableName())
}
