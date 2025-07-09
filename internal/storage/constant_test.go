// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstant_TableName(t *testing.T) {
	constant := Constant{}
	assert.Equal(t, "constant", constant.TableName())
}
