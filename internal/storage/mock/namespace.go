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
func (mr *MockINamespaceMockRecorder) ByNamespaceId(ctx, namespaceId any) *INamespaceByNamespaceIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByNamespaceId", reflect.TypeOf((*MockINamespace)(nil).ByNamespaceId), ctx, namespaceId)
	return &INamespaceByNamespaceIdCall{Call: call}
}

// INamespaceByNamespaceIdCall wrap *gomock.Call
type INamespaceByNamespaceIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *INamespaceByNamespaceIdCall) Return(arg0 []storage.Namespace, arg1 error) *INamespaceByNamespaceIdCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *INamespaceByNamespaceIdCall) Do(f func(context.Context, []byte) ([]storage.Namespace, error)) *INamespaceByNamespaceIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *INamespaceByNamespaceIdCall) DoAndReturn(f func(context.Context, []byte) ([]storage.Namespace, error)) *INamespaceByNamespaceIdCall {
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
func (mr *MockINamespaceMockRecorder) ByNamespaceIdAndVersion(ctx, namespaceId, version any) *INamespaceByNamespaceIdAndVersionCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByNamespaceIdAndVersion", reflect.TypeOf((*MockINamespace)(nil).ByNamespaceIdAndVersion), ctx, namespaceId, version)
	return &INamespaceByNamespaceIdAndVersionCall{Call: call}
}

// INamespaceByNamespaceIdAndVersionCall wrap *gomock.Call
type INamespaceByNamespaceIdAndVersionCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *INamespaceByNamespaceIdAndVersionCall) Return(arg0 storage.Namespace, arg1 error) *INamespaceByNamespaceIdAndVersionCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *INamespaceByNamespaceIdAndVersionCall) Do(f func(context.Context, []byte, byte) (storage.Namespace, error)) *INamespaceByNamespaceIdAndVersionCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *INamespaceByNamespaceIdAndVersionCall) DoAndReturn(f func(context.Context, []byte, byte) (storage.Namespace, error)) *INamespaceByNamespaceIdAndVersionCall {
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
func (mr *MockINamespaceMockRecorder) CursorList(ctx, id, limit, order, cmp any) *INamespaceCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockINamespace)(nil).CursorList), ctx, id, limit, order, cmp)
	return &INamespaceCursorListCall{Call: call}
}

// INamespaceCursorListCall wrap *gomock.Call
type INamespaceCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *INamespaceCursorListCall) Return(arg0 []*storage.Namespace, arg1 error) *INamespaceCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *INamespaceCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Namespace, error)) *INamespaceCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *INamespaceCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Namespace, error)) *INamespaceCursorListCall {
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
func (mr *MockINamespaceMockRecorder) GetByID(ctx, id any) *INamespaceGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockINamespace)(nil).GetByID), ctx, id)
	return &INamespaceGetByIDCall{Call: call}
}

// INamespaceGetByIDCall wrap *gomock.Call
type INamespaceGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *INamespaceGetByIDCall) Return(arg0 *storage.Namespace, arg1 error) *INamespaceGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *INamespaceGetByIDCall) Do(f func(context.Context, uint64) (*storage.Namespace, error)) *INamespaceGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *INamespaceGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.Namespace, error)) *INamespaceGetByIDCall {
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
func (mr *MockINamespaceMockRecorder) GetByIds(ctx any, ids ...any) *INamespaceGetByIdsCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, ids...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIds", reflect.TypeOf((*MockINamespace)(nil).GetByIds), varargs...)
	return &INamespaceGetByIdsCall{Call: call}
}

// INamespaceGetByIdsCall wrap *gomock.Call
type INamespaceGetByIdsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *INamespaceGetByIdsCall) Return(ns []storage.Namespace, err error) *INamespaceGetByIdsCall {
	c.Call = c.Call.Return(ns, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *INamespaceGetByIdsCall) Do(f func(context.Context, ...uint64) ([]storage.Namespace, error)) *INamespaceGetByIdsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *INamespaceGetByIdsCall) DoAndReturn(f func(context.Context, ...uint64) ([]storage.Namespace, error)) *INamespaceGetByIdsCall {
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
func (mr *MockINamespaceMockRecorder) IsNoRows(err any) *INamespaceIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockINamespace)(nil).IsNoRows), err)
	return &INamespaceIsNoRowsCall{Call: call}
}

// INamespaceIsNoRowsCall wrap *gomock.Call
type INamespaceIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *INamespaceIsNoRowsCall) Return(arg0 bool) *INamespaceIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *INamespaceIsNoRowsCall) Do(f func(error) bool) *INamespaceIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *INamespaceIsNoRowsCall) DoAndReturn(f func(error) bool) *INamespaceIsNoRowsCall {
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
func (mr *MockINamespaceMockRecorder) LastID(ctx any) *INamespaceLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockINamespace)(nil).LastID), ctx)
	return &INamespaceLastIDCall{Call: call}
}

// INamespaceLastIDCall wrap *gomock.Call
type INamespaceLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *INamespaceLastIDCall) Return(arg0 uint64, arg1 error) *INamespaceLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *INamespaceLastIDCall) Do(f func(context.Context) (uint64, error)) *INamespaceLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *INamespaceLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *INamespaceLastIDCall {
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
func (mr *MockINamespaceMockRecorder) List(ctx, limit, offset, order any) *INamespaceListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockINamespace)(nil).List), ctx, limit, offset, order)
	return &INamespaceListCall{Call: call}
}

// INamespaceListCall wrap *gomock.Call
type INamespaceListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *INamespaceListCall) Return(arg0 []*storage.Namespace, arg1 error) *INamespaceListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *INamespaceListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Namespace, error)) *INamespaceListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *INamespaceListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Namespace, error)) *INamespaceListCall {
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
func (mr *MockINamespaceMockRecorder) ListWithSort(ctx, sortField, sort, limit, offset any) *INamespaceListWithSortCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWithSort", reflect.TypeOf((*MockINamespace)(nil).ListWithSort), ctx, sortField, sort, limit, offset)
	return &INamespaceListWithSortCall{Call: call}
}

// INamespaceListWithSortCall wrap *gomock.Call
type INamespaceListWithSortCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *INamespaceListWithSortCall) Return(ns []storage.Namespace, err error) *INamespaceListWithSortCall {
	c.Call = c.Call.Return(ns, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *INamespaceListWithSortCall) Do(f func(context.Context, string, storage0.SortOrder, int, int) ([]storage.Namespace, error)) *INamespaceListWithSortCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *INamespaceListWithSortCall) DoAndReturn(f func(context.Context, string, storage0.SortOrder, int, int) ([]storage.Namespace, error)) *INamespaceListWithSortCall {
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
func (mr *MockINamespaceMockRecorder) Messages(ctx, id, limit, offset any) *INamespaceMessagesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Messages", reflect.TypeOf((*MockINamespace)(nil).Messages), ctx, id, limit, offset)
	return &INamespaceMessagesCall{Call: call}
}

// INamespaceMessagesCall wrap *gomock.Call
type INamespaceMessagesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *INamespaceMessagesCall) Return(arg0 []storage.NamespaceMessage, arg1 error) *INamespaceMessagesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *INamespaceMessagesCall) Do(f func(context.Context, uint64, int, int) ([]storage.NamespaceMessage, error)) *INamespaceMessagesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *INamespaceMessagesCall) DoAndReturn(f func(context.Context, uint64, int, int) ([]storage.NamespaceMessage, error)) *INamespaceMessagesCall {
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
func (mr *MockINamespaceMockRecorder) Save(ctx, m any) *INamespaceSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockINamespace)(nil).Save), ctx, m)
	return &INamespaceSaveCall{Call: call}
}

// INamespaceSaveCall wrap *gomock.Call
type INamespaceSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *INamespaceSaveCall) Return(arg0 error) *INamespaceSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *INamespaceSaveCall) Do(f func(context.Context, *storage.Namespace) error) *INamespaceSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *INamespaceSaveCall) DoAndReturn(f func(context.Context, *storage.Namespace) error) *INamespaceSaveCall {
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
func (mr *MockINamespaceMockRecorder) Update(ctx, m any) *INamespaceUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockINamespace)(nil).Update), ctx, m)
	return &INamespaceUpdateCall{Call: call}
}

// INamespaceUpdateCall wrap *gomock.Call
type INamespaceUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *INamespaceUpdateCall) Return(arg0 error) *INamespaceUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *INamespaceUpdateCall) Do(f func(context.Context, *storage.Namespace) error) *INamespaceUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *INamespaceUpdateCall) DoAndReturn(f func(context.Context, *storage.Namespace) error) *INamespaceUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
