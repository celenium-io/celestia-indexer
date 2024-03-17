// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"
)

func (s *StorageTestSuite) TestGrantByGrantee() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	grants, err := s.storage.Grants.ByGrantee(ctx, 1, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(grants, 1)

	grant := grants[0]
	s.Require().EqualValues(1, grant.Id)
	s.Require().EqualValues(1000, grant.Height)
	s.Require().EqualValues("/cosmos.staking.v1beta1.MsgDelegate", grant.Authorization)
	s.Require().NotNil(grant.Params)

	s.Require().NotNil(grant.Granter)
	s.Require().EqualValues("celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", grant.Granter.Address)
}

func (s *StorageTestSuite) TestGrantByGranter() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	grants, err := s.storage.Grants.ByGranter(ctx, 2, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(grants, 1)

	grant := grants[0]
	s.Require().EqualValues(1, grant.Id)
	s.Require().EqualValues(1000, grant.Height)
	s.Require().EqualValues("/cosmos.staking.v1beta1.MsgDelegate", grant.Authorization)
	s.Require().NotNil(grant.Params)

	s.Require().NotNil(grant.Grantee)
	s.Require().EqualValues("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", grant.Grantee.Address)
}
