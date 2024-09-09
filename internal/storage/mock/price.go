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
	time "time"

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
func (m *MockIPrice) Get(ctx context.Context, timeframe string, start, end time.Time, limit int) ([]storage.Price, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, timeframe, start, end, limit)
	ret0, _ := ret[0].([]storage.Price)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockIPriceMockRecorder) Get(ctx, timeframe, start, end, limit any) *MockIPriceGetCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockIPrice)(nil).Get), ctx, timeframe, start, end, limit)
	return &MockIPriceGetCall{Call: call}
}

// MockIPriceGetCall wrap *gomock.Call
type MockIPriceGetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIPriceGetCall) Return(arg0 []storage.Price, arg1 error) *MockIPriceGetCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIPriceGetCall) Do(f func(context.Context, string, time.Time, time.Time, int) ([]storage.Price, error)) *MockIPriceGetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIPriceGetCall) DoAndReturn(f func(context.Context, string, time.Time, time.Time, int) ([]storage.Price, error)) *MockIPriceGetCall {
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
func (mr *MockIPriceMockRecorder) Last(ctx any) *MockIPriceLastCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Last", reflect.TypeOf((*MockIPrice)(nil).Last), ctx)
	return &MockIPriceLastCall{Call: call}
}

// MockIPriceLastCall wrap *gomock.Call
type MockIPriceLastCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIPriceLastCall) Return(arg0 storage.Price, arg1 error) *MockIPriceLastCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIPriceLastCall) Do(f func(context.Context) (storage.Price, error)) *MockIPriceLastCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIPriceLastCall) DoAndReturn(f func(context.Context) (storage.Price, error)) *MockIPriceLastCall {
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
func (mr *MockIPriceMockRecorder) Save(ctx, price any) *MockIPriceSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIPrice)(nil).Save), ctx, price)
	return &MockIPriceSaveCall{Call: call}
}

// MockIPriceSaveCall wrap *gomock.Call
type MockIPriceSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIPriceSaveCall) Return(arg0 error) *MockIPriceSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIPriceSaveCall) Do(f func(context.Context, *storage.Price) error) *MockIPriceSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIPriceSaveCall) DoAndReturn(f func(context.Context, *storage.Price) error) *MockIPriceSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
