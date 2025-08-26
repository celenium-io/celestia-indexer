// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestUpgradeList() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	for _, fltrs := range []storage.ListUpgradesFilter{
		{
			Offset: 0,
			Sort:   sdk.SortOrderDesc,
			TxId:   testsuite.Ptr(uint64(1)),
		}, {
			Offset: 0,
			Sort:   sdk.SortOrderDesc,
			Height: 1010,
		}, {
			Offset:   0,
			Sort:     sdk.SortOrderDesc,
			SignerId: testsuite.Ptr(uint64(1)),
		}, {
			Limit:  1,
			Offset: 0,
			Sort:   sdk.SortOrderAsc,
		},
	} {

		upgrades, err := s.storage.Upgrade.List(ctx, fltrs)
		s.Require().NoError(err)
		s.Require().Len(upgrades, 1)

		upgrade := upgrades[0]
		s.Require().EqualValues(1, upgrade.Id)
		s.Require().EqualValues(1010, upgrade.Height)
		s.Require().EqualValues(1499, upgrade.Version)
		s.Require().EqualValues(1, upgrade.MsgId)
		s.Require().EqualValues(1, upgrade.TxId)
		s.Require().NotNil(upgrade.Signer)
		s.Require().EqualValues("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", upgrade.Signer.Address)
	}
}
