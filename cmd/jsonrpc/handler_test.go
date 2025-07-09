// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	nodeMock "github.com/celenium-io/celestia-indexer/pkg/node/mock"
	"github.com/celenium-io/celestia-indexer/pkg/node/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// HandlerTestSuite -
type HandlerTestSuite struct {
	suite.Suite
	ns           *mock.MockINamespace
	logs         *mock.MockIBlobLog
	blobReceiver *nodeMock.MockDalApi
	echo         *echo.Echo
	handler      *BlobHandler
	ctrl         *gomock.Controller
}

// SetupSuite -
func (s *HandlerTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.ctrl = gomock.NewController(s.T())
	s.ns = mock.NewMockINamespace(s.ctrl)
	s.logs = mock.NewMockIBlobLog(s.ctrl)
	s.blobReceiver = nodeMock.NewMockDalApi(s.ctrl)
	s.handler = NewBlobHandler(s.blobReceiver, s.logs, s.ns)
}

// TearDownSuite -
func (s *HandlerTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteHandler_Run(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (s *HandlerTestSuite) TestGet() {
	data := testsuite.RandomText(100)
	s.blobReceiver.EXPECT().
		Blob(gomock.Any(), pkgTypes.Level(6357101), "AAAAAAAAAAAAAAAAAAAAAAAAAMod4SqXFNgavI8=", "TekzHEK3JoBFapZdQFUoPwtuzUgpEd6OWyhGGzFTv3s=").
		Return(types.Blob{
			Namespace:    "AAAAAAAAAAAAAAAAAAAAAAAAAMod4SqXFNgavI8=",
			Commitment:   "TekzHEK3JoBFapZdQFUoPwtuzUgpEd6OWyhGGzFTv3s=",
			ShareVersion: 0,
			Data:         data,
		}, nil).
		Times(1)

	response, err := s.handler.Get(s.T().Context(), 6357101, "AAAAAAAAAAAAAAAAAAAAAAAAAMod4SqXFNgavI8=", "TekzHEK3JoBFapZdQFUoPwtuzUgpEd6OWyhGGzFTv3s=")
	s.Require().NoError(err)
	s.Require().EqualValues("AAAAAAAAAAAAAAAAAAAAAAAAAMod4SqXFNgavI8=", response.Namespace)
	s.Require().EqualValues("TekzHEK3JoBFapZdQFUoPwtuzUgpEd6OWyhGGzFTv3s=", response.Commitment)
	s.Require().EqualValues(data, response.Data)
	s.Require().EqualValues(0, response.ShareVersion)
}

func (s *HandlerTestSuite) TestGetAll() {
	data := testsuite.RandomText(100)

	s.ns.EXPECT().
		ByNamespaceIdAndVersion(gomock.Any(), gomock.Any(), byte(0)).
		Return(storage.Namespace{
			Id: 1,
		}, nil).
		Times(1)

	s.logs.EXPECT().
		ByNamespace(gomock.Any(), uint64(1), storage.BlobLogFilters{
			Limit:  100,
			Height: 6357101,
		}).
		Return([]storage.BlobLog{
			{
				Commitment: "TekzHEK3JoBFapZdQFUoPwtuzUgpEd6OWyhGGzFTv3s=",
			},
		}, nil).
		Times(1)

	s.blobReceiver.EXPECT().
		Blob(gomock.Any(), pkgTypes.Level(6357101), "AAAAAAAAAAAAAAAAAAAAAAAAAMod4SqXFNgavI8=", "TekzHEK3JoBFapZdQFUoPwtuzUgpEd6OWyhGGzFTv3s=").
		Return(types.Blob{
			Namespace:    "AAAAAAAAAAAAAAAAAAAAAAAAAMod4SqXFNgavI8=",
			Commitment:   "TekzHEK3JoBFapZdQFUoPwtuzUgpEd6OWyhGGzFTv3s=",
			ShareVersion: 0,
			Data:         data,
		}, nil).
		Times(1)

	response, err := s.handler.GetAll(s.T().Context(), 6357101, []string{"AAAAAAAAAAAAAAAAAAAAAAAAAMod4SqXFNgavI8="})
	s.Require().NoError(err)
	s.Require().Len(response, 1)

	s.Require().EqualValues("AAAAAAAAAAAAAAAAAAAAAAAAAMod4SqXFNgavI8=", response[0].Namespace)
	s.Require().EqualValues("TekzHEK3JoBFapZdQFUoPwtuzUgpEd6OWyhGGzFTv3s=", response[0].Commitment)
	s.Require().EqualValues(data, response[0].Data)
	s.Require().EqualValues(0, response[0].ShareVersion)
}
