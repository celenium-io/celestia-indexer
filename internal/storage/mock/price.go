// Code generated by MockGen. DO NOT EDIT.
// Source: price.go
//
// Generated by this command:
//
//	mockgen -source=price.go -destination=mock/price.go -package=mock -typed
//
// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	storage "github.com/celenium-io/celestia-indexer/internal/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockIPrice is a mock of IPrice interface.
type MockIPrice struct {
	ctrl     *gomock.Controller
	recorder *MockIPriceMockRecorder
}

// MockIPriceMockRecorder is the mock recorder for MockIPrice.
type MockIPriceMockRecorder struct {
	mock *MockIPrice
}

// NewMockIPrice creates a new mock instance.
func NewMockIPrice(ctrl *gomock.Controller) *MockIPrice {
	mock := &MockIPrice{ctrl: ctrl}
	mock.recorder = &MockIPriceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPrice) EXPECT() *MockIPriceMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockIPrice) Get(ctx context.Context, timeframe string, start, end int64, limit int) ([]storage.Price, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, timeframe, start, end, limit)
	ret0, _ := ret[0].([]storage.Price)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockIPriceMockRecorder) Get(ctx, timeframe, start, end, limit any) *IPriceGetCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockIPrice)(nil).Get), ctx, timeframe, start, end, limit)
	return &IPriceGetCall{Call: call}
}

// IPriceGetCall wrap *gomock.Call
type IPriceGetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IPriceGetCall) Return(arg0 []storage.Price, arg1 error) *IPriceGetCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IPriceGetCall) Do(f func(context.Context, string, int64, int64, int) ([]storage.Price, error)) *IPriceGetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IPriceGetCall) DoAndReturn(f func(context.Context, string, int64, int64, int) ([]storage.Price, error)) *IPriceGetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Last mocks base method.
func (m *MockIPrice) Last(ctx context.Context) (storage.Price, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Last", ctx)
	ret0, _ := ret[0].(storage.Price)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Last indicates an expected call of Last.
func (mr *MockIPriceMockRecorder) Last(ctx any) *IPriceLastCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Last", reflect.TypeOf((*MockIPrice)(nil).Last), ctx)
	return &IPriceLastCall{Call: call}
}

// IPriceLastCall wrap *gomock.Call
type IPriceLastCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IPriceLastCall) Return(arg0 storage.Price, arg1 error) *IPriceLastCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IPriceLastCall) Do(f func(context.Context) (storage.Price, error)) *IPriceLastCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IPriceLastCall) DoAndReturn(f func(context.Context) (storage.Price, error)) *IPriceLastCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m *MockIPrice) Save(ctx context.Context, price *storage.Price) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, price)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIPriceMockRecorder) Save(ctx, price any) *IPriceSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIPrice)(nil).Save), ctx, price)
	return &IPriceSaveCall{Call: call}
}

// IPriceSaveCall wrap *gomock.Call
type IPriceSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IPriceSaveCall) Return(arg0 error) *IPriceSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IPriceSaveCall) Do(f func(context.Context, *storage.Price) error) *IPriceSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IPriceSaveCall) DoAndReturn(f func(context.Context, *storage.Price) error) *IPriceSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
