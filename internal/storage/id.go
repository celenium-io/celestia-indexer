// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
)

const maxHeight = 0xffffffffff
const maxPosition = 0xffffff

func idFromHeightAndPosition(height types.Level, position int64) (uint64, error) {
	if height > maxHeight {
		return 0, errors.Errorf("can't get id: overflow max height value %d", height)
	}

	// 5 bytes for height and 3 bytes for position
	if height > 0 {
		if position > maxPosition {
			return 0, errors.Errorf("can't get id: overflow max position value %d", position)
		}
		return uint64(height)<<24 | uint64(position), nil
	}

	if position-1 > maxPosition {
		return 0, errors.Errorf("can't get id: overflow max position value %d", position)
	}
	// for genesis block for avoiding zero id
	return uint64(height)<<24 | uint64(position+1), nil
}
