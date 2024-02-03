// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/types"
)

func (s *StorageTestSuite) TestBlockSignatureLevels() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	levels, err := s.storage.BlockSignatures.LevelsByValidator(ctx, 1, 998)
	s.Require().NoError(err)
	s.Require().Len(levels, 2)

	s.Require().Equal([]types.Level{1000, 999}, levels)
}
