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
func (mr *MockIJailMockRecorder) ByValidator(ctx, id, limit, offset any) *MockIJailByValidatorCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByValidator", reflect.TypeOf((*MockIJail)(nil).ByValidator), ctx, id, limit, offset)
	return &MockIJailByValidatorCall{Call: call}
}

// MockIJailByValidatorCall wrap *gomock.Call
type MockIJailByValidatorCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIJailByValidatorCall) Return(arg0 []storage.Jail, arg1 error) *MockIJailByValidatorCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIJailByValidatorCall) Do(f func(context.Context, uint64, int, int) ([]storage.Jail, error)) *MockIJailByValidatorCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIJailByValidatorCall) DoAndReturn(f func(context.Context, uint64, int, int) ([]storage.Jail, error)) *MockIJailByValidatorCall {
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
func (mr *MockIJailMockRecorder) CursorList(ctx, id, limit, order, cmp any) *MockIJailCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIJail)(nil).CursorList), ctx, id, limit, order, cmp)
	return &MockIJailCursorListCall{Call: call}
}

// MockIJailCursorListCall wrap *gomock.Call
type MockIJailCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIJailCursorListCall) Return(arg0 []*storage.Jail, arg1 error) *MockIJailCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIJailCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Jail, error)) *MockIJailCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIJailCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Jail, error)) *MockIJailCursorListCall {
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
func (mr *MockIJailMockRecorder) GetByID(ctx, id any) *MockIJailGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIJail)(nil).GetByID), ctx, id)
	return &MockIJailGetByIDCall{Call: call}
}

// MockIJailGetByIDCall wrap *gomock.Call
type MockIJailGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIJailGetByIDCall) Return(arg0 *storage.Jail, arg1 error) *MockIJailGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIJailGetByIDCall) Do(f func(context.Context, uint64) (*storage.Jail, error)) *MockIJailGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIJailGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.Jail, error)) *MockIJailGetByIDCall {
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
func (mr *MockIJailMockRecorder) IsNoRows(err any) *MockIJailIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIJail)(nil).IsNoRows), err)
	return &MockIJailIsNoRowsCall{Call: call}
}

// MockIJailIsNoRowsCall wrap *gomock.Call
type MockIJailIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIJailIsNoRowsCall) Return(arg0 bool) *MockIJailIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIJailIsNoRowsCall) Do(f func(error) bool) *MockIJailIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIJailIsNoRowsCall) DoAndReturn(f func(error) bool) *MockIJailIsNoRowsCall {
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
func (mr *MockIJailMockRecorder) LastID(ctx any) *MockIJailLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIJail)(nil).LastID), ctx)
	return &MockIJailLastIDCall{Call: call}
}

// MockIJailLastIDCall wrap *gomock.Call
type MockIJailLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIJailLastIDCall) Return(arg0 uint64, arg1 error) *MockIJailLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIJailLastIDCall) Do(f func(context.Context) (uint64, error)) *MockIJailLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIJailLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *MockIJailLastIDCall {
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
func (mr *MockIJailMockRecorder) List(ctx, limit, offset, order any) *MockIJailListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIJail)(nil).List), ctx, limit, offset, order)
	return &MockIJailListCall{Call: call}
}

// MockIJailListCall wrap *gomock.Call
type MockIJailListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIJailListCall) Return(arg0 []*storage.Jail, arg1 error) *MockIJailListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIJailListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Jail, error)) *MockIJailListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIJailListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Jail, error)) *MockIJailListCall {
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
func (mr *MockIJailMockRecorder) Save(ctx, m any) *MockIJailSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIJail)(nil).Save), ctx, m)
	return &MockIJailSaveCall{Call: call}
}

// MockIJailSaveCall wrap *gomock.Call
type MockIJailSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIJailSaveCall) Return(arg0 error) *MockIJailSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIJailSaveCall) Do(f func(context.Context, *storage.Jail) error) *MockIJailSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIJailSaveCall) DoAndReturn(f func(context.Context, *storage.Jail) error) *MockIJailSaveCall {
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
func (mr *MockIJailMockRecorder) Update(ctx, m any) *MockIJailUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIJail)(nil).Update), ctx, m)
	return &MockIJailUpdateCall{Call: call}
}

// MockIJailUpdateCall wrap *gomock.Call
type MockIJailUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIJailUpdateCall) Return(arg0 error) *MockIJailUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIJailUpdateCall) Do(f func(context.Context, *storage.Jail) error) *MockIJailUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIJailUpdateCall) DoAndReturn(f func(context.Context, *storage.Jail) error) *MockIJailUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
