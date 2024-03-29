// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: vesting_period.go
//
// Generated by this command:
//
//	mockgen -source=vesting_period.go -destination=mock/vesting_period.go -package=mock -typed
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

// MockIVestingPeriod is a mock of IVestingPeriod interface.
type MockIVestingPeriod struct {
	ctrl     *gomock.Controller
	recorder *MockIVestingPeriodMockRecorder
}

// MockIVestingPeriodMockRecorder is the mock recorder for MockIVestingPeriod.
type MockIVestingPeriodMockRecorder struct {
	mock *MockIVestingPeriod
}

// NewMockIVestingPeriod creates a new mock instance.
func NewMockIVestingPeriod(ctrl *gomock.Controller) *MockIVestingPeriod {
	mock := &MockIVestingPeriod{ctrl: ctrl}
	mock.recorder = &MockIVestingPeriodMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIVestingPeriod) EXPECT() *MockIVestingPeriodMockRecorder {
	return m.recorder
}

// ByVesting mocks base method.
func (m *MockIVestingPeriod) ByVesting(ctx context.Context, id uint64, limit, offset int) ([]storage.VestingPeriod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByVesting", ctx, id, limit, offset)
	ret0, _ := ret[0].([]storage.VestingPeriod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByVesting indicates an expected call of ByVesting.
func (mr *MockIVestingPeriodMockRecorder) ByVesting(ctx, id, limit, offset any) *IVestingPeriodByVestingCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByVesting", reflect.TypeOf((*MockIVestingPeriod)(nil).ByVesting), ctx, id, limit, offset)
	return &IVestingPeriodByVestingCall{Call: call}
}

// IVestingPeriodByVestingCall wrap *gomock.Call
type IVestingPeriodByVestingCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IVestingPeriodByVestingCall) Return(arg0 []storage.VestingPeriod, arg1 error) *IVestingPeriodByVestingCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IVestingPeriodByVestingCall) Do(f func(context.Context, uint64, int, int) ([]storage.VestingPeriod, error)) *IVestingPeriodByVestingCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IVestingPeriodByVestingCall) DoAndReturn(f func(context.Context, uint64, int, int) ([]storage.VestingPeriod, error)) *IVestingPeriodByVestingCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CursorList mocks base method.
func (m *MockIVestingPeriod) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.VestingPeriod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.VestingPeriod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIVestingPeriodMockRecorder) CursorList(ctx, id, limit, order, cmp any) *IVestingPeriodCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIVestingPeriod)(nil).CursorList), ctx, id, limit, order, cmp)
	return &IVestingPeriodCursorListCall{Call: call}
}

// IVestingPeriodCursorListCall wrap *gomock.Call
type IVestingPeriodCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IVestingPeriodCursorListCall) Return(arg0 []*storage.VestingPeriod, arg1 error) *IVestingPeriodCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IVestingPeriodCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.VestingPeriod, error)) *IVestingPeriodCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IVestingPeriodCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.VestingPeriod, error)) *IVestingPeriodCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIVestingPeriod) GetByID(ctx context.Context, id uint64) (*storage.VestingPeriod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.VestingPeriod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIVestingPeriodMockRecorder) GetByID(ctx, id any) *IVestingPeriodGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIVestingPeriod)(nil).GetByID), ctx, id)
	return &IVestingPeriodGetByIDCall{Call: call}
}

// IVestingPeriodGetByIDCall wrap *gomock.Call
type IVestingPeriodGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IVestingPeriodGetByIDCall) Return(arg0 *storage.VestingPeriod, arg1 error) *IVestingPeriodGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IVestingPeriodGetByIDCall) Do(f func(context.Context, uint64) (*storage.VestingPeriod, error)) *IVestingPeriodGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IVestingPeriodGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.VestingPeriod, error)) *IVestingPeriodGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIVestingPeriod) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIVestingPeriodMockRecorder) IsNoRows(err any) *IVestingPeriodIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIVestingPeriod)(nil).IsNoRows), err)
	return &IVestingPeriodIsNoRowsCall{Call: call}
}

// IVestingPeriodIsNoRowsCall wrap *gomock.Call
type IVestingPeriodIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IVestingPeriodIsNoRowsCall) Return(arg0 bool) *IVestingPeriodIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IVestingPeriodIsNoRowsCall) Do(f func(error) bool) *IVestingPeriodIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IVestingPeriodIsNoRowsCall) DoAndReturn(f func(error) bool) *IVestingPeriodIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIVestingPeriod) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIVestingPeriodMockRecorder) LastID(ctx any) *IVestingPeriodLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIVestingPeriod)(nil).LastID), ctx)
	return &IVestingPeriodLastIDCall{Call: call}
}

// IVestingPeriodLastIDCall wrap *gomock.Call
type IVestingPeriodLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IVestingPeriodLastIDCall) Return(arg0 uint64, arg1 error) *IVestingPeriodLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IVestingPeriodLastIDCall) Do(f func(context.Context) (uint64, error)) *IVestingPeriodLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IVestingPeriodLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *IVestingPeriodLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIVestingPeriod) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.VestingPeriod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.VestingPeriod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIVestingPeriodMockRecorder) List(ctx, limit, offset, order any) *IVestingPeriodListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIVestingPeriod)(nil).List), ctx, limit, offset, order)
	return &IVestingPeriodListCall{Call: call}
}

// IVestingPeriodListCall wrap *gomock.Call
type IVestingPeriodListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IVestingPeriodListCall) Return(arg0 []*storage.VestingPeriod, arg1 error) *IVestingPeriodListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IVestingPeriodListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.VestingPeriod, error)) *IVestingPeriodListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IVestingPeriodListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.VestingPeriod, error)) *IVestingPeriodListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIVestingPeriod) Save(ctx context.Context, m *storage.VestingPeriod) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIVestingPeriodMockRecorder) Save(ctx, m any) *IVestingPeriodSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIVestingPeriod)(nil).Save), ctx, m)
	return &IVestingPeriodSaveCall{Call: call}
}

// IVestingPeriodSaveCall wrap *gomock.Call
type IVestingPeriodSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IVestingPeriodSaveCall) Return(arg0 error) *IVestingPeriodSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IVestingPeriodSaveCall) Do(f func(context.Context, *storage.VestingPeriod) error) *IVestingPeriodSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IVestingPeriodSaveCall) DoAndReturn(f func(context.Context, *storage.VestingPeriod) error) *IVestingPeriodSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIVestingPeriod) Update(ctx context.Context, m *storage.VestingPeriod) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIVestingPeriodMockRecorder) Update(ctx, m any) *IVestingPeriodUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIVestingPeriod)(nil).Update), ctx, m)
	return &IVestingPeriodUpdateCall{Call: call}
}

// IVestingPeriodUpdateCall wrap *gomock.Call
type IVestingPeriodUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IVestingPeriodUpdateCall) Return(arg0 error) *IVestingPeriodUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IVestingPeriodUpdateCall) Do(f func(context.Context, *storage.VestingPeriod) error) *IVestingPeriodUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IVestingPeriodUpdateCall) DoAndReturn(f func(context.Context, *storage.VestingPeriod) error) *IVestingPeriodUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
