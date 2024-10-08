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

	binance "github.com/celenium-io/celestia-indexer/internal/binance"
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

// OHLC mocks base method.
func (m *MockIApi) OHLC(ctx context.Context, symbol, interval string, arguments *binance.OHLCArgs) ([]binance.OHLC, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OHLC", ctx, symbol, interval, arguments)
	ret0, _ := ret[0].([]binance.OHLC)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OHLC indicates an expected call of OHLC.
func (mr *MockIApiMockRecorder) OHLC(ctx, symbol, interval, arguments any) *MockIApiOHLCCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OHLC", reflect.TypeOf((*MockIApi)(nil).OHLC), ctx, symbol, interval, arguments)
	return &MockIApiOHLCCall{Call: call}
}

// MockIApiOHLCCall wrap *gomock.Call
type MockIApiOHLCCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIApiOHLCCall) Return(candles []binance.OHLC, err error) *MockIApiOHLCCall {
	c.Call = c.Call.Return(candles, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIApiOHLCCall) Do(f func(context.Context, string, string, *binance.OHLCArgs) ([]binance.OHLC, error)) *MockIApiOHLCCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIApiOHLCCall) DoAndReturn(f func(context.Context, string, string, *binance.OHLCArgs) ([]binance.OHLC, error)) *MockIApiOHLCCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
