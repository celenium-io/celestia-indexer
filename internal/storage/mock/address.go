// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: address.go
//
// Generated by this command:
//
//	mockgen -source=address.go -destination=mock/address.go -package=mock -typed
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	storage "github.com/celenium-io/celestia-indexer/internal/storage"
	storage0 "github.com/dipdup-net/indexer-sdk/pkg/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockIAddress is a mock of IAddress interface.
type MockIAddress struct {
	ctrl     *gomock.Controller
	recorder *MockIAddressMockRecorder
}

// MockIAddressMockRecorder is the mock recorder for MockIAddress.
type MockIAddressMockRecorder struct {
	mock *MockIAddress
}

// NewMockIAddress creates a new mock instance.
func NewMockIAddress(ctrl *gomock.Controller) *MockIAddress {
	mock := &MockIAddress{ctrl: ctrl}
	mock.recorder = &MockIAddressMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAddress) EXPECT() *MockIAddressMockRecorder {
	return m.recorder
}

// ByHash mocks base method.
func (m *MockIAddress) ByHash(ctx context.Context, hash []byte) (storage.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByHash", ctx, hash)
	ret0, _ := ret[0].(storage.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByHash indicates an expected call of ByHash.
func (mr *MockIAddressMockRecorder) ByHash(ctx, hash any) *MockIAddressByHashCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByHash", reflect.TypeOf((*MockIAddress)(nil).ByHash), ctx, hash)
	return &MockIAddressByHashCall{Call: call}
}

// MockIAddressByHashCall wrap *gomock.Call
type MockIAddressByHashCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAddressByHashCall) Return(arg0 storage.Address, arg1 error) *MockIAddressByHashCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAddressByHashCall) Do(f func(context.Context, []byte) (storage.Address, error)) *MockIAddressByHashCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAddressByHashCall) DoAndReturn(f func(context.Context, []byte) (storage.Address, error)) *MockIAddressByHashCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CursorList mocks base method.
func (m *MockIAddress) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIAddressMockRecorder) CursorList(ctx, id, limit, order, cmp any) *MockIAddressCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIAddress)(nil).CursorList), ctx, id, limit, order, cmp)
	return &MockIAddressCursorListCall{Call: call}
}

// MockIAddressCursorListCall wrap *gomock.Call
type MockIAddressCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAddressCursorListCall) Return(arg0 []*storage.Address, arg1 error) *MockIAddressCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAddressCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Address, error)) *MockIAddressCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAddressCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Address, error)) *MockIAddressCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIAddress) GetByID(ctx context.Context, id uint64) (*storage.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIAddressMockRecorder) GetByID(ctx, id any) *MockIAddressGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIAddress)(nil).GetByID), ctx, id)
	return &MockIAddressGetByIDCall{Call: call}
}

// MockIAddressGetByIDCall wrap *gomock.Call
type MockIAddressGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAddressGetByIDCall) Return(arg0 *storage.Address, arg1 error) *MockIAddressGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAddressGetByIDCall) Do(f func(context.Context, uint64) (*storage.Address, error)) *MockIAddressGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAddressGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.Address, error)) *MockIAddressGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IdByAddress mocks base method.
func (m *MockIAddress) IdByAddress(ctx context.Context, address string, ids ...uint64) (uint64, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, address}
	for _, a := range ids {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IdByAddress", varargs...)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IdByAddress indicates an expected call of IdByAddress.
func (mr *MockIAddressMockRecorder) IdByAddress(ctx, address any, ids ...any) *MockIAddressIdByAddressCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, address}, ids...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IdByAddress", reflect.TypeOf((*MockIAddress)(nil).IdByAddress), varargs...)
	return &MockIAddressIdByAddressCall{Call: call}
}

// MockIAddressIdByAddressCall wrap *gomock.Call
type MockIAddressIdByAddressCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAddressIdByAddressCall) Return(arg0 uint64, arg1 error) *MockIAddressIdByAddressCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAddressIdByAddressCall) Do(f func(context.Context, string, ...uint64) (uint64, error)) *MockIAddressIdByAddressCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAddressIdByAddressCall) DoAndReturn(f func(context.Context, string, ...uint64) (uint64, error)) *MockIAddressIdByAddressCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IdByHash mocks base method.
func (m *MockIAddress) IdByHash(ctx context.Context, hash ...[]byte) ([]uint64, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx}
	for _, a := range hash {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IdByHash", varargs...)
	ret0, _ := ret[0].([]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IdByHash indicates an expected call of IdByHash.
func (mr *MockIAddressMockRecorder) IdByHash(ctx any, hash ...any) *MockIAddressIdByHashCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, hash...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IdByHash", reflect.TypeOf((*MockIAddress)(nil).IdByHash), varargs...)
	return &MockIAddressIdByHashCall{Call: call}
}

// MockIAddressIdByHashCall wrap *gomock.Call
type MockIAddressIdByHashCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAddressIdByHashCall) Return(arg0 []uint64, arg1 error) *MockIAddressIdByHashCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAddressIdByHashCall) Do(f func(context.Context, ...[]byte) ([]uint64, error)) *MockIAddressIdByHashCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAddressIdByHashCall) DoAndReturn(f func(context.Context, ...[]byte) ([]uint64, error)) *MockIAddressIdByHashCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIAddress) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIAddressMockRecorder) IsNoRows(err any) *MockIAddressIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIAddress)(nil).IsNoRows), err)
	return &MockIAddressIsNoRowsCall{Call: call}
}

// MockIAddressIsNoRowsCall wrap *gomock.Call
type MockIAddressIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAddressIsNoRowsCall) Return(arg0 bool) *MockIAddressIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAddressIsNoRowsCall) Do(f func(error) bool) *MockIAddressIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAddressIsNoRowsCall) DoAndReturn(f func(error) bool) *MockIAddressIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIAddress) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIAddressMockRecorder) LastID(ctx any) *MockIAddressLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIAddress)(nil).LastID), ctx)
	return &MockIAddressLastIDCall{Call: call}
}

// MockIAddressLastIDCall wrap *gomock.Call
type MockIAddressLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAddressLastIDCall) Return(arg0 uint64, arg1 error) *MockIAddressLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAddressLastIDCall) Do(f func(context.Context) (uint64, error)) *MockIAddressLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAddressLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *MockIAddressLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIAddress) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIAddressMockRecorder) List(ctx, limit, offset, order any) *MockIAddressListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIAddress)(nil).List), ctx, limit, offset, order)
	return &MockIAddressListCall{Call: call}
}

// MockIAddressListCall wrap *gomock.Call
type MockIAddressListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAddressListCall) Return(arg0 []*storage.Address, arg1 error) *MockIAddressListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAddressListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Address, error)) *MockIAddressListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAddressListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Address, error)) *MockIAddressListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ListWithBalance mocks base method.
func (m *MockIAddress) ListWithBalance(ctx context.Context, filters storage.AddressListFilter) ([]storage.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWithBalance", ctx, filters)
	ret0, _ := ret[0].([]storage.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWithBalance indicates an expected call of ListWithBalance.
func (mr *MockIAddressMockRecorder) ListWithBalance(ctx, filters any) *MockIAddressListWithBalanceCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWithBalance", reflect.TypeOf((*MockIAddress)(nil).ListWithBalance), ctx, filters)
	return &MockIAddressListWithBalanceCall{Call: call}
}

// MockIAddressListWithBalanceCall wrap *gomock.Call
type MockIAddressListWithBalanceCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAddressListWithBalanceCall) Return(arg0 []storage.Address, arg1 error) *MockIAddressListWithBalanceCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAddressListWithBalanceCall) Do(f func(context.Context, storage.AddressListFilter) ([]storage.Address, error)) *MockIAddressListWithBalanceCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAddressListWithBalanceCall) DoAndReturn(f func(context.Context, storage.AddressListFilter) ([]storage.Address, error)) *MockIAddressListWithBalanceCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIAddress) Save(ctx context.Context, m *storage.Address) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIAddressMockRecorder) Save(ctx, m any) *MockIAddressSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIAddress)(nil).Save), ctx, m)
	return &MockIAddressSaveCall{Call: call}
}

// MockIAddressSaveCall wrap *gomock.Call
type MockIAddressSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAddressSaveCall) Return(arg0 error) *MockIAddressSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAddressSaveCall) Do(f func(context.Context, *storage.Address) error) *MockIAddressSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAddressSaveCall) DoAndReturn(f func(context.Context, *storage.Address) error) *MockIAddressSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Series mocks base method.
func (m *MockIAddress) Series(ctx context.Context, addressId uint64, timeframe storage.Timeframe, column string, req storage.SeriesRequest) ([]storage.HistogramItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Series", ctx, addressId, timeframe, column, req)
	ret0, _ := ret[0].([]storage.HistogramItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Series indicates an expected call of Series.
func (mr *MockIAddressMockRecorder) Series(ctx, addressId, timeframe, column, req any) *MockIAddressSeriesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Series", reflect.TypeOf((*MockIAddress)(nil).Series), ctx, addressId, timeframe, column, req)
	return &MockIAddressSeriesCall{Call: call}
}

// MockIAddressSeriesCall wrap *gomock.Call
type MockIAddressSeriesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAddressSeriesCall) Return(items []storage.HistogramItem, err error) *MockIAddressSeriesCall {
	c.Call = c.Call.Return(items, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAddressSeriesCall) Do(f func(context.Context, uint64, storage.Timeframe, string, storage.SeriesRequest) ([]storage.HistogramItem, error)) *MockIAddressSeriesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAddressSeriesCall) DoAndReturn(f func(context.Context, uint64, storage.Timeframe, string, storage.SeriesRequest) ([]storage.HistogramItem, error)) *MockIAddressSeriesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIAddress) Update(ctx context.Context, m *storage.Address) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIAddressMockRecorder) Update(ctx, m any) *MockIAddressUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIAddress)(nil).Update), ctx, m)
	return &MockIAddressUpdateCall{Call: call}
}

// MockIAddressUpdateCall wrap *gomock.Call
type MockIAddressUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAddressUpdateCall) Return(arg0 error) *MockIAddressUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAddressUpdateCall) Do(f func(context.Context, *storage.Address) error) *MockIAddressUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAddressUpdateCall) DoAndReturn(f func(context.Context, *storage.Address) error) *MockIAddressUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
