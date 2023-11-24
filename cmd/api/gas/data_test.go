// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
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
	err := q.Range(func(item info) (bool, error) {
		totalTx += item.TxCount
		return false, nil
	})
	require.NoError(t, err)
	require.EqualValues(t, 20, totalTx)
}
