// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package gas

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	q := newQueue(10)

	for i := 0; i < 10000; i++ {
		q.Push(info{
			Height:  uint64(i),
			TxCount: 2,
		})
	}

	var totalTx int64
	for item := range q.All() {
		totalTx += item.TxCount
	}
	require.EqualValues(t, 20, totalTx)
}
