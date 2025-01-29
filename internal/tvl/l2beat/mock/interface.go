// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go
//
// Generated by this command:
//
//	mockgen -source=interface.go -destination=mock/interface.go -package=mock -typed
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	l2beat "github.com/celenium-io/celestia-indexer/internal/tvl/l2beat"
	gomock "go.uber.org/mock/gomock"
)

// MockIApi is a mock of IApi interface.
type MockIApi struct {
	ctrl     *gomock.Controller
	recorder *MockIApiMockRecorder
}

// MockIApiMockRecorder is the mock recorder for MockIApi.
type MockIApiMockRecorder struct {
	mock *MockIApi
}

// NewMockIApi creates a new mock instance.
func NewMockIApi(ctrl *gomock.Controller) *MockIApi {
	mock := &MockIApi{ctrl: ctrl}
	mock.recorder = &MockIApiMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIApi) EXPECT() *MockIApiMockRecorder {
	return m.recorder
}

// TVL mocks base method.
func (m *MockIApi) TVL(ctx context.Context, rollupName string, timeframe l2beat.TvlTimeframe) (l2beat.TVLResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TVL", ctx, rollupName, timeframe)
	ret0, _ := ret[0].(l2beat.TVLResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TVL indicates an expected call of TVL.
func (mr *MockIApiMockRecorder) TVL(ctx, rollupName, timeframe any) *MockIApiTVLCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TVL", reflect.TypeOf((*MockIApi)(nil).TVL), ctx, rollupName, timeframe)
	return &MockIApiTVLCall{Call: call}
}

// MockIApiTVLCall wrap *gomock.Call
type MockIApiTVLCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIApiTVLCall) Return(result l2beat.TVLResponse, err error) *MockIApiTVLCall {
	c.Call = c.Call.Return(result, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIApiTVLCall) Do(f func(context.Context, string, l2beat.TvlTimeframe) (l2beat.TVLResponse, error)) *MockIApiTVLCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIApiTVLCall) DoAndReturn(f func(context.Context, string, l2beat.TvlTimeframe) (l2beat.TVLResponse, error)) *MockIApiTVLCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
