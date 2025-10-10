// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"
)

func (s *StorageTestSuite) TestRollupProviderByRollupId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	providers, err := s.storage.RollupProvider.ByRollupId(ctx, 1)
	s.Require().NoError(err)
	s.Require().Len(providers, 2)

	provider1 := providers[0]
	s.Require().EqualValues(1, provider1.NamespaceId)
	s.Require().EqualValues(1, provider1.AddressId)
	s.Require().NotNil(provider1.Address)
	s.Require().NotNil(provider1.Namespace)
	s.Require().Equal("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", provider1.Address.Address)
	s.Require().Len(provider1.Namespace.NamespaceID, 18)
	s.Require().EqualValues(0, provider1.Namespace.Version)

	provider2 := providers[1]
	s.Require().EqualValues(2, provider2.NamespaceId)
	s.Require().EqualValues(1, provider2.AddressId)
	s.Require().NotNil(provider2.Address)
	s.Require().NotNil(provider2.Namespace)
	s.Require().Equal("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", provider2.Address.Address)
	s.Require().Len(provider2.Namespace.NamespaceID, 18)
	s.Require().EqualValues(1, provider2.Namespace.Version)
}
