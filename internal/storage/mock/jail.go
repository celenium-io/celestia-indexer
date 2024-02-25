// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: jail.go
//
// Generated by this command:
//
//	mockgen -source=jail.go -destination=mock/jail.go -package=mock -typed
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

// MockIJail is a mock of IJail interface.
type MockIJail struct {
	ctrl     *gomock.Controller
	recorder *MockIJailMockRecorder
}

// MockIJailMockRecorder is the mock recorder for MockIJail.
type MockIJailMockRecorder struct {
	mock *MockIJail
}

// NewMockIJail creates a new mock instance.
func NewMockIJail(ctrl *gomock.Controller) *MockIJail {
	mock := &MockIJail{ctrl: ctrl}
	mock.recorder = &MockIJailMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIJail) EXPECT() *MockIJailMockRecorder {
	return m.recorder
}

// ByValidator mocks base method.
func (m *MockIJail) ByValidator(ctx context.Context, id uint64, limit, offset int) ([]storage.Jail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByValidator", ctx, id, limit, offset)
	ret0, _ := ret[0].([]storage.Jail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByValidator indicates an expected call of ByValidator.
func (mr *MockIJailMockRecorder) ByValidator(ctx, id, limit, offset any) *IJailByValidatorCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByValidator", reflect.TypeOf((*MockIJail)(nil).ByValidator), ctx, id, limit, offset)
	return &IJailByValidatorCall{Call: call}
}

// IJailByValidatorCall wrap *gomock.Call
type IJailByValidatorCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IJailByValidatorCall) Return(arg0 []storage.Jail, arg1 error) *IJailByValidatorCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IJailByValidatorCall) Do(f func(context.Context, uint64, int, int) ([]storage.Jail, error)) *IJailByValidatorCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IJailByValidatorCall) DoAndReturn(f func(context.Context, uint64, int, int) ([]storage.Jail, error)) *IJailByValidatorCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CursorList mocks base method.
func (m *MockIJail) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.Jail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.Jail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIJailMockRecorder) CursorList(ctx, id, limit, order, cmp any) *IJailCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIJail)(nil).CursorList), ctx, id, limit, order, cmp)
	return &IJailCursorListCall{Call: call}
}

// IJailCursorListCall wrap *gomock.Call
type IJailCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IJailCursorListCall) Return(arg0 []*storage.Jail, arg1 error) *IJailCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IJailCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Jail, error)) *IJailCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IJailCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Jail, error)) *IJailCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIJail) GetByID(ctx context.Context, id uint64) (*storage.Jail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.Jail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIJailMockRecorder) GetByID(ctx, id any) *IJailGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIJail)(nil).GetByID), ctx, id)
	return &IJailGetByIDCall{Call: call}
}

// IJailGetByIDCall wrap *gomock.Call
type IJailGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IJailGetByIDCall) Return(arg0 *storage.Jail, arg1 error) *IJailGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IJailGetByIDCall) Do(f func(context.Context, uint64) (*storage.Jail, error)) *IJailGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IJailGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.Jail, error)) *IJailGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIJail) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIJailMockRecorder) IsNoRows(err any) *IJailIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIJail)(nil).IsNoRows), err)
	return &IJailIsNoRowsCall{Call: call}
}

// IJailIsNoRowsCall wrap *gomock.Call
type IJailIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IJailIsNoRowsCall) Return(arg0 bool) *IJailIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IJailIsNoRowsCall) Do(f func(error) bool) *IJailIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IJailIsNoRowsCall) DoAndReturn(f func(error) bool) *IJailIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIJail) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIJailMockRecorder) LastID(ctx any) *IJailLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIJail)(nil).LastID), ctx)
	return &IJailLastIDCall{Call: call}
}

// IJailLastIDCall wrap *gomock.Call
type IJailLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IJailLastIDCall) Return(arg0 uint64, arg1 error) *IJailLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IJailLastIDCall) Do(f func(context.Context) (uint64, error)) *IJailLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IJailLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *IJailLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIJail) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.Jail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.Jail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIJailMockRecorder) List(ctx, limit, offset, order any) *IJailListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIJail)(nil).List), ctx, limit, offset, order)
	return &IJailListCall{Call: call}
}

// IJailListCall wrap *gomock.Call
type IJailListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IJailListCall) Return(arg0 []*storage.Jail, arg1 error) *IJailListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IJailListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Jail, error)) *IJailListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IJailListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Jail, error)) *IJailListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIJail) Save(ctx context.Context, m *storage.Jail) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIJailMockRecorder) Save(ctx, m any) *IJailSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIJail)(nil).Save), ctx, m)
	return &IJailSaveCall{Call: call}
}

// IJailSaveCall wrap *gomock.Call
type IJailSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IJailSaveCall) Return(arg0 error) *IJailSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IJailSaveCall) Do(f func(context.Context, *storage.Jail) error) *IJailSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IJailSaveCall) DoAndReturn(f func(context.Context, *storage.Jail) error) *IJailSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIJail) Update(ctx context.Context, m *storage.Jail) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIJailMockRecorder) Update(ctx, m any) *IJailUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIJail)(nil).Update), ctx, m)
	return &IJailUpdateCall{Call: call}
}

// IJailUpdateCall wrap *gomock.Call
type IJailUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IJailUpdateCall) Return(arg0 error) *IJailUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IJailUpdateCall) Do(f func(context.Context, *storage.Jail) error) *IJailUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IJailUpdateCall) DoAndReturn(f func(context.Context, *storage.Jail) error) *IJailUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
