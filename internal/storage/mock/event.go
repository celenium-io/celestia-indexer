// Code generated by MockGen. DO NOT EDIT.
// Source: event.go
//
// Generated by this command:
//
//	mockgen -source=event.go -destination=mock/event.go -package=mock -typed
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	storage "github.com/celenium-io/celestia-indexer/internal/storage"
	types "github.com/celenium-io/celestia-indexer/pkg/types"
	storage0 "github.com/dipdup-net/indexer-sdk/pkg/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockIEvent is a mock of IEvent interface.
type MockIEvent struct {
	ctrl     *gomock.Controller
	recorder *MockIEventMockRecorder
}

// MockIEventMockRecorder is the mock recorder for MockIEvent.
type MockIEventMockRecorder struct {
	mock *MockIEvent
}

// NewMockIEvent creates a new mock instance.
func NewMockIEvent(ctrl *gomock.Controller) *MockIEvent {
	mock := &MockIEvent{ctrl: ctrl}
	mock.recorder = &MockIEventMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIEvent) EXPECT() *MockIEventMockRecorder {
	return m.recorder
}

// ByBlock mocks base method.
func (m *MockIEvent) ByBlock(ctx context.Context, height types.Level, fltrs storage.EventFilter) ([]storage.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByBlock", ctx, height, fltrs)
	ret0, _ := ret[0].([]storage.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByBlock indicates an expected call of ByBlock.
func (mr *MockIEventMockRecorder) ByBlock(ctx, height, fltrs any) *MockIEventByBlockCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByBlock", reflect.TypeOf((*MockIEvent)(nil).ByBlock), ctx, height, fltrs)
	return &MockIEventByBlockCall{Call: call}
}

// MockIEventByBlockCall wrap *gomock.Call
type MockIEventByBlockCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIEventByBlockCall) Return(arg0 []storage.Event, arg1 error) *MockIEventByBlockCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIEventByBlockCall) Do(f func(context.Context, types.Level, storage.EventFilter) ([]storage.Event, error)) *MockIEventByBlockCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIEventByBlockCall) DoAndReturn(f func(context.Context, types.Level, storage.EventFilter) ([]storage.Event, error)) *MockIEventByBlockCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ByTxId mocks base method.
func (m *MockIEvent) ByTxId(ctx context.Context, txId uint64, fltrs storage.EventFilter) ([]storage.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByTxId", ctx, txId, fltrs)
	ret0, _ := ret[0].([]storage.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByTxId indicates an expected call of ByTxId.
func (mr *MockIEventMockRecorder) ByTxId(ctx, txId, fltrs any) *MockIEventByTxIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByTxId", reflect.TypeOf((*MockIEvent)(nil).ByTxId), ctx, txId, fltrs)
	return &MockIEventByTxIdCall{Call: call}
}

// MockIEventByTxIdCall wrap *gomock.Call
type MockIEventByTxIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIEventByTxIdCall) Return(arg0 []storage.Event, arg1 error) *MockIEventByTxIdCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIEventByTxIdCall) Do(f func(context.Context, uint64, storage.EventFilter) ([]storage.Event, error)) *MockIEventByTxIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIEventByTxIdCall) DoAndReturn(f func(context.Context, uint64, storage.EventFilter) ([]storage.Event, error)) *MockIEventByTxIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CursorList mocks base method.
func (m *MockIEvent) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIEventMockRecorder) CursorList(ctx, id, limit, order, cmp any) *MockIEventCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIEvent)(nil).CursorList), ctx, id, limit, order, cmp)
	return &MockIEventCursorListCall{Call: call}
}

// MockIEventCursorListCall wrap *gomock.Call
type MockIEventCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIEventCursorListCall) Return(arg0 []*storage.Event, arg1 error) *MockIEventCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIEventCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Event, error)) *MockIEventCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIEventCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Event, error)) *MockIEventCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIEvent) GetByID(ctx context.Context, id uint64) (*storage.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIEventMockRecorder) GetByID(ctx, id any) *MockIEventGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIEvent)(nil).GetByID), ctx, id)
	return &MockIEventGetByIDCall{Call: call}
}

// MockIEventGetByIDCall wrap *gomock.Call
type MockIEventGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIEventGetByIDCall) Return(arg0 *storage.Event, arg1 error) *MockIEventGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIEventGetByIDCall) Do(f func(context.Context, uint64) (*storage.Event, error)) *MockIEventGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIEventGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.Event, error)) *MockIEventGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIEvent) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIEventMockRecorder) IsNoRows(err any) *MockIEventIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIEvent)(nil).IsNoRows), err)
	return &MockIEventIsNoRowsCall{Call: call}
}

// MockIEventIsNoRowsCall wrap *gomock.Call
type MockIEventIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIEventIsNoRowsCall) Return(arg0 bool) *MockIEventIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIEventIsNoRowsCall) Do(f func(error) bool) *MockIEventIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIEventIsNoRowsCall) DoAndReturn(f func(error) bool) *MockIEventIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIEvent) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIEventMockRecorder) LastID(ctx any) *MockIEventLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIEvent)(nil).LastID), ctx)
	return &MockIEventLastIDCall{Call: call}
}

// MockIEventLastIDCall wrap *gomock.Call
type MockIEventLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIEventLastIDCall) Return(arg0 uint64, arg1 error) *MockIEventLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIEventLastIDCall) Do(f func(context.Context) (uint64, error)) *MockIEventLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIEventLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *MockIEventLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIEvent) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIEventMockRecorder) List(ctx, limit, offset, order any) *MockIEventListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIEvent)(nil).List), ctx, limit, offset, order)
	return &MockIEventListCall{Call: call}
}

// MockIEventListCall wrap *gomock.Call
type MockIEventListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIEventListCall) Return(arg0 []*storage.Event, arg1 error) *MockIEventListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIEventListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Event, error)) *MockIEventListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIEventListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Event, error)) *MockIEventListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIEvent) Save(ctx context.Context, m *storage.Event) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIEventMockRecorder) Save(ctx, m any) *MockIEventSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIEvent)(nil).Save), ctx, m)
	return &MockIEventSaveCall{Call: call}
}

// MockIEventSaveCall wrap *gomock.Call
type MockIEventSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIEventSaveCall) Return(arg0 error) *MockIEventSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIEventSaveCall) Do(f func(context.Context, *storage.Event) error) *MockIEventSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIEventSaveCall) DoAndReturn(f func(context.Context, *storage.Event) error) *MockIEventSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIEvent) Update(ctx context.Context, m *storage.Event) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIEventMockRecorder) Update(ctx, m any) *MockIEventUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIEvent)(nil).Update), ctx, m)
	return &MockIEventUpdateCall{Call: call}
}

// MockIEventUpdateCall wrap *gomock.Call
type MockIEventUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIEventUpdateCall) Return(arg0 error) *MockIEventUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIEventUpdateCall) Do(f func(context.Context, *storage.Event) error) *MockIEventUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIEventUpdateCall) DoAndReturn(f func(context.Context, *storage.Event) error) *MockIEventUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
