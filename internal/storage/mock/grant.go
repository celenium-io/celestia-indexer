// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: grant.go
//
// Generated by this command:
//
//	mockgen -source=grant.go -destination=mock/grant.go -package=mock -typed
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

// MockIGrant is a mock of IGrant interface.
type MockIGrant struct {
	ctrl     *gomock.Controller
	recorder *MockIGrantMockRecorder
}

// MockIGrantMockRecorder is the mock recorder for MockIGrant.
type MockIGrantMockRecorder struct {
	mock *MockIGrant
}

// NewMockIGrant creates a new mock instance.
func NewMockIGrant(ctrl *gomock.Controller) *MockIGrant {
	mock := &MockIGrant{ctrl: ctrl}
	mock.recorder = &MockIGrantMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIGrant) EXPECT() *MockIGrantMockRecorder {
	return m.recorder
}

// ByGrantee mocks base method.
func (m *MockIGrant) ByGrantee(ctx context.Context, id uint64, limit, offset int) ([]storage.Grant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByGrantee", ctx, id, limit, offset)
	ret0, _ := ret[0].([]storage.Grant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByGrantee indicates an expected call of ByGrantee.
func (mr *MockIGrantMockRecorder) ByGrantee(ctx, id, limit, offset any) *IGrantByGranteeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByGrantee", reflect.TypeOf((*MockIGrant)(nil).ByGrantee), ctx, id, limit, offset)
	return &IGrantByGranteeCall{Call: call}
}

// IGrantByGranteeCall wrap *gomock.Call
type IGrantByGranteeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IGrantByGranteeCall) Return(arg0 []storage.Grant, arg1 error) *IGrantByGranteeCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IGrantByGranteeCall) Do(f func(context.Context, uint64, int, int) ([]storage.Grant, error)) *IGrantByGranteeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IGrantByGranteeCall) DoAndReturn(f func(context.Context, uint64, int, int) ([]storage.Grant, error)) *IGrantByGranteeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ByGranter mocks base method.
func (m *MockIGrant) ByGranter(ctx context.Context, id uint64, limit, offset int) ([]storage.Grant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByGranter", ctx, id, limit, offset)
	ret0, _ := ret[0].([]storage.Grant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByGranter indicates an expected call of ByGranter.
func (mr *MockIGrantMockRecorder) ByGranter(ctx, id, limit, offset any) *IGrantByGranterCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByGranter", reflect.TypeOf((*MockIGrant)(nil).ByGranter), ctx, id, limit, offset)
	return &IGrantByGranterCall{Call: call}
}

// IGrantByGranterCall wrap *gomock.Call
type IGrantByGranterCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IGrantByGranterCall) Return(arg0 []storage.Grant, arg1 error) *IGrantByGranterCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IGrantByGranterCall) Do(f func(context.Context, uint64, int, int) ([]storage.Grant, error)) *IGrantByGranterCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IGrantByGranterCall) DoAndReturn(f func(context.Context, uint64, int, int) ([]storage.Grant, error)) *IGrantByGranterCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CursorList mocks base method.
func (m *MockIGrant) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.Grant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.Grant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIGrantMockRecorder) CursorList(ctx, id, limit, order, cmp any) *IGrantCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIGrant)(nil).CursorList), ctx, id, limit, order, cmp)
	return &IGrantCursorListCall{Call: call}
}

// IGrantCursorListCall wrap *gomock.Call
type IGrantCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IGrantCursorListCall) Return(arg0 []*storage.Grant, arg1 error) *IGrantCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IGrantCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Grant, error)) *IGrantCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IGrantCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Grant, error)) *IGrantCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIGrant) GetByID(ctx context.Context, id uint64) (*storage.Grant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.Grant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIGrantMockRecorder) GetByID(ctx, id any) *IGrantGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIGrant)(nil).GetByID), ctx, id)
	return &IGrantGetByIDCall{Call: call}
}

// IGrantGetByIDCall wrap *gomock.Call
type IGrantGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IGrantGetByIDCall) Return(arg0 *storage.Grant, arg1 error) *IGrantGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IGrantGetByIDCall) Do(f func(context.Context, uint64) (*storage.Grant, error)) *IGrantGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IGrantGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.Grant, error)) *IGrantGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIGrant) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIGrantMockRecorder) IsNoRows(err any) *IGrantIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIGrant)(nil).IsNoRows), err)
	return &IGrantIsNoRowsCall{Call: call}
}

// IGrantIsNoRowsCall wrap *gomock.Call
type IGrantIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IGrantIsNoRowsCall) Return(arg0 bool) *IGrantIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IGrantIsNoRowsCall) Do(f func(error) bool) *IGrantIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IGrantIsNoRowsCall) DoAndReturn(f func(error) bool) *IGrantIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIGrant) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIGrantMockRecorder) LastID(ctx any) *IGrantLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIGrant)(nil).LastID), ctx)
	return &IGrantLastIDCall{Call: call}
}

// IGrantLastIDCall wrap *gomock.Call
type IGrantLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IGrantLastIDCall) Return(arg0 uint64, arg1 error) *IGrantLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IGrantLastIDCall) Do(f func(context.Context) (uint64, error)) *IGrantLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IGrantLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *IGrantLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIGrant) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.Grant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.Grant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIGrantMockRecorder) List(ctx, limit, offset, order any) *IGrantListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIGrant)(nil).List), ctx, limit, offset, order)
	return &IGrantListCall{Call: call}
}

// IGrantListCall wrap *gomock.Call
type IGrantListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IGrantListCall) Return(arg0 []*storage.Grant, arg1 error) *IGrantListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IGrantListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Grant, error)) *IGrantListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IGrantListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Grant, error)) *IGrantListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIGrant) Save(ctx context.Context, m *storage.Grant) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIGrantMockRecorder) Save(ctx, m any) *IGrantSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIGrant)(nil).Save), ctx, m)
	return &IGrantSaveCall{Call: call}
}

// IGrantSaveCall wrap *gomock.Call
type IGrantSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IGrantSaveCall) Return(arg0 error) *IGrantSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IGrantSaveCall) Do(f func(context.Context, *storage.Grant) error) *IGrantSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IGrantSaveCall) DoAndReturn(f func(context.Context, *storage.Grant) error) *IGrantSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIGrant) Update(ctx context.Context, m *storage.Grant) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIGrantMockRecorder) Update(ctx, m any) *IGrantUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIGrant)(nil).Update), ctx, m)
	return &IGrantUpdateCall{Call: call}
}

// IGrantUpdateCall wrap *gomock.Call
type IGrantUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IGrantUpdateCall) Return(arg0 error) *IGrantUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IGrantUpdateCall) Do(f func(context.Context, *storage.Grant) error) *IGrantUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IGrantUpdateCall) DoAndReturn(f func(context.Context, *storage.Grant) error) *IGrantUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
