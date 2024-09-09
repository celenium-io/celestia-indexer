// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: namespace.go
//
// Generated by this command:
//
//	mockgen -source=namespace.go -destination=mock/namespace.go -package=mock -typed
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

// MockINamespace is a mock of INamespace interface.
type MockINamespace struct {
	ctrl     *gomock.Controller
	recorder *MockINamespaceMockRecorder
}

// MockINamespaceMockRecorder is the mock recorder for MockINamespace.
type MockINamespaceMockRecorder struct {
	mock *MockINamespace
}

// NewMockINamespace creates a new mock instance.
func NewMockINamespace(ctrl *gomock.Controller) *MockINamespace {
	mock := &MockINamespace{ctrl: ctrl}
	mock.recorder = &MockINamespaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockINamespace) EXPECT() *MockINamespaceMockRecorder {
	return m.recorder
}

// ByNamespaceId mocks base method.
func (m *MockINamespace) ByNamespaceId(ctx context.Context, namespaceId []byte) ([]storage.Namespace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByNamespaceId", ctx, namespaceId)
	ret0, _ := ret[0].([]storage.Namespace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByNamespaceId indicates an expected call of ByNamespaceId.
func (mr *MockINamespaceMockRecorder) ByNamespaceId(ctx, namespaceId any) *MockINamespaceByNamespaceIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByNamespaceId", reflect.TypeOf((*MockINamespace)(nil).ByNamespaceId), ctx, namespaceId)
	return &MockINamespaceByNamespaceIdCall{Call: call}
}

// MockINamespaceByNamespaceIdCall wrap *gomock.Call
type MockINamespaceByNamespaceIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockINamespaceByNamespaceIdCall) Return(arg0 []storage.Namespace, arg1 error) *MockINamespaceByNamespaceIdCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockINamespaceByNamespaceIdCall) Do(f func(context.Context, []byte) ([]storage.Namespace, error)) *MockINamespaceByNamespaceIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockINamespaceByNamespaceIdCall) DoAndReturn(f func(context.Context, []byte) ([]storage.Namespace, error)) *MockINamespaceByNamespaceIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ByNamespaceIdAndVersion mocks base method.
func (m *MockINamespace) ByNamespaceIdAndVersion(ctx context.Context, namespaceId []byte, version byte) (storage.Namespace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByNamespaceIdAndVersion", ctx, namespaceId, version)
	ret0, _ := ret[0].(storage.Namespace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByNamespaceIdAndVersion indicates an expected call of ByNamespaceIdAndVersion.
func (mr *MockINamespaceMockRecorder) ByNamespaceIdAndVersion(ctx, namespaceId, version any) *MockINamespaceByNamespaceIdAndVersionCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByNamespaceIdAndVersion", reflect.TypeOf((*MockINamespace)(nil).ByNamespaceIdAndVersion), ctx, namespaceId, version)
	return &MockINamespaceByNamespaceIdAndVersionCall{Call: call}
}

// MockINamespaceByNamespaceIdAndVersionCall wrap *gomock.Call
type MockINamespaceByNamespaceIdAndVersionCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockINamespaceByNamespaceIdAndVersionCall) Return(arg0 storage.Namespace, arg1 error) *MockINamespaceByNamespaceIdAndVersionCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockINamespaceByNamespaceIdAndVersionCall) Do(f func(context.Context, []byte, byte) (storage.Namespace, error)) *MockINamespaceByNamespaceIdAndVersionCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockINamespaceByNamespaceIdAndVersionCall) DoAndReturn(f func(context.Context, []byte, byte) (storage.Namespace, error)) *MockINamespaceByNamespaceIdAndVersionCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CursorList mocks base method.
func (m *MockINamespace) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.Namespace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.Namespace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockINamespaceMockRecorder) CursorList(ctx, id, limit, order, cmp any) *MockINamespaceCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockINamespace)(nil).CursorList), ctx, id, limit, order, cmp)
	return &MockINamespaceCursorListCall{Call: call}
}

// MockINamespaceCursorListCall wrap *gomock.Call
type MockINamespaceCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockINamespaceCursorListCall) Return(arg0 []*storage.Namespace, arg1 error) *MockINamespaceCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockINamespaceCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Namespace, error)) *MockINamespaceCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockINamespaceCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Namespace, error)) *MockINamespaceCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockINamespace) GetByID(ctx context.Context, id uint64) (*storage.Namespace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.Namespace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockINamespaceMockRecorder) GetByID(ctx, id any) *MockINamespaceGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockINamespace)(nil).GetByID), ctx, id)
	return &MockINamespaceGetByIDCall{Call: call}
}

// MockINamespaceGetByIDCall wrap *gomock.Call
type MockINamespaceGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockINamespaceGetByIDCall) Return(arg0 *storage.Namespace, arg1 error) *MockINamespaceGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockINamespaceGetByIDCall) Do(f func(context.Context, uint64) (*storage.Namespace, error)) *MockINamespaceGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockINamespaceGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.Namespace, error)) *MockINamespaceGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByIds mocks base method.
func (m *MockINamespace) GetByIds(ctx context.Context, ids ...uint64) ([]storage.Namespace, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx}
	for _, a := range ids {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetByIds", varargs...)
	ret0, _ := ret[0].([]storage.Namespace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByIds indicates an expected call of GetByIds.
func (mr *MockINamespaceMockRecorder) GetByIds(ctx any, ids ...any) *MockINamespaceGetByIdsCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, ids...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIds", reflect.TypeOf((*MockINamespace)(nil).GetByIds), varargs...)
	return &MockINamespaceGetByIdsCall{Call: call}
}

// MockINamespaceGetByIdsCall wrap *gomock.Call
type MockINamespaceGetByIdsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockINamespaceGetByIdsCall) Return(ns []storage.Namespace, err error) *MockINamespaceGetByIdsCall {
	c.Call = c.Call.Return(ns, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockINamespaceGetByIdsCall) Do(f func(context.Context, ...uint64) ([]storage.Namespace, error)) *MockINamespaceGetByIdsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockINamespaceGetByIdsCall) DoAndReturn(f func(context.Context, ...uint64) ([]storage.Namespace, error)) *MockINamespaceGetByIdsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockINamespace) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockINamespaceMockRecorder) IsNoRows(err any) *MockINamespaceIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockINamespace)(nil).IsNoRows), err)
	return &MockINamespaceIsNoRowsCall{Call: call}
}

// MockINamespaceIsNoRowsCall wrap *gomock.Call
type MockINamespaceIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockINamespaceIsNoRowsCall) Return(arg0 bool) *MockINamespaceIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockINamespaceIsNoRowsCall) Do(f func(error) bool) *MockINamespaceIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockINamespaceIsNoRowsCall) DoAndReturn(f func(error) bool) *MockINamespaceIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockINamespace) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockINamespaceMockRecorder) LastID(ctx any) *MockINamespaceLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockINamespace)(nil).LastID), ctx)
	return &MockINamespaceLastIDCall{Call: call}
}

// MockINamespaceLastIDCall wrap *gomock.Call
type MockINamespaceLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockINamespaceLastIDCall) Return(arg0 uint64, arg1 error) *MockINamespaceLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockINamespaceLastIDCall) Do(f func(context.Context) (uint64, error)) *MockINamespaceLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockINamespaceLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *MockINamespaceLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockINamespace) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.Namespace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.Namespace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockINamespaceMockRecorder) List(ctx, limit, offset, order any) *MockINamespaceListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockINamespace)(nil).List), ctx, limit, offset, order)
	return &MockINamespaceListCall{Call: call}
}

// MockINamespaceListCall wrap *gomock.Call
type MockINamespaceListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockINamespaceListCall) Return(arg0 []*storage.Namespace, arg1 error) *MockINamespaceListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockINamespaceListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Namespace, error)) *MockINamespaceListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockINamespaceListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Namespace, error)) *MockINamespaceListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ListWithSort mocks base method.
func (m *MockINamespace) ListWithSort(ctx context.Context, sortField string, sort storage0.SortOrder, limit, offset int) ([]storage.Namespace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWithSort", ctx, sortField, sort, limit, offset)
	ret0, _ := ret[0].([]storage.Namespace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWithSort indicates an expected call of ListWithSort.
func (mr *MockINamespaceMockRecorder) ListWithSort(ctx, sortField, sort, limit, offset any) *MockINamespaceListWithSortCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWithSort", reflect.TypeOf((*MockINamespace)(nil).ListWithSort), ctx, sortField, sort, limit, offset)
	return &MockINamespaceListWithSortCall{Call: call}
}

// MockINamespaceListWithSortCall wrap *gomock.Call
type MockINamespaceListWithSortCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockINamespaceListWithSortCall) Return(ns []storage.Namespace, err error) *MockINamespaceListWithSortCall {
	c.Call = c.Call.Return(ns, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockINamespaceListWithSortCall) Do(f func(context.Context, string, storage0.SortOrder, int, int) ([]storage.Namespace, error)) *MockINamespaceListWithSortCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockINamespaceListWithSortCall) DoAndReturn(f func(context.Context, string, storage0.SortOrder, int, int) ([]storage.Namespace, error)) *MockINamespaceListWithSortCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Messages mocks base method.
func (m *MockINamespace) Messages(ctx context.Context, id uint64, limit, offset int) ([]storage.NamespaceMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Messages", ctx, id, limit, offset)
	ret0, _ := ret[0].([]storage.NamespaceMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Messages indicates an expected call of Messages.
func (mr *MockINamespaceMockRecorder) Messages(ctx, id, limit, offset any) *MockINamespaceMessagesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Messages", reflect.TypeOf((*MockINamespace)(nil).Messages), ctx, id, limit, offset)
	return &MockINamespaceMessagesCall{Call: call}
}

// MockINamespaceMessagesCall wrap *gomock.Call
type MockINamespaceMessagesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockINamespaceMessagesCall) Return(arg0 []storage.NamespaceMessage, arg1 error) *MockINamespaceMessagesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockINamespaceMessagesCall) Do(f func(context.Context, uint64, int, int) ([]storage.NamespaceMessage, error)) *MockINamespaceMessagesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockINamespaceMessagesCall) DoAndReturn(f func(context.Context, uint64, int, int) ([]storage.NamespaceMessage, error)) *MockINamespaceMessagesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockINamespace) Save(ctx context.Context, m *storage.Namespace) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockINamespaceMockRecorder) Save(ctx, m any) *MockINamespaceSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockINamespace)(nil).Save), ctx, m)
	return &MockINamespaceSaveCall{Call: call}
}

// MockINamespaceSaveCall wrap *gomock.Call
type MockINamespaceSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockINamespaceSaveCall) Return(arg0 error) *MockINamespaceSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockINamespaceSaveCall) Do(f func(context.Context, *storage.Namespace) error) *MockINamespaceSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockINamespaceSaveCall) DoAndReturn(f func(context.Context, *storage.Namespace) error) *MockINamespaceSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockINamespace) Update(ctx context.Context, m *storage.Namespace) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockINamespaceMockRecorder) Update(ctx, m any) *MockINamespaceUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockINamespace)(nil).Update), ctx, m)
	return &MockINamespaceUpdateCall{Call: call}
}

// MockINamespaceUpdateCall wrap *gomock.Call
type MockINamespaceUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockINamespaceUpdateCall) Return(arg0 error) *MockINamespaceUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockINamespaceUpdateCall) Do(f func(context.Context, *storage.Namespace) error) *MockINamespaceUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockINamespaceUpdateCall) DoAndReturn(f func(context.Context, *storage.Namespace) error) *MockINamespaceUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
