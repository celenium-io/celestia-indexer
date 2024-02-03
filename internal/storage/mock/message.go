// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: message.go
//
// Generated by this command:
//
//	mockgen -source=message.go -destination=mock/message.go -package=mock -typed
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

// MockIMessage is a mock of IMessage interface.
type MockIMessage struct {
	ctrl     *gomock.Controller
	recorder *MockIMessageMockRecorder
}

// MockIMessageMockRecorder is the mock recorder for MockIMessage.
type MockIMessageMockRecorder struct {
	mock *MockIMessage
}

// NewMockIMessage creates a new mock instance.
func NewMockIMessage(ctrl *gomock.Controller) *MockIMessage {
	mock := &MockIMessage{ctrl: ctrl}
	mock.recorder = &MockIMessageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIMessage) EXPECT() *MockIMessageMockRecorder {
	return m.recorder
}

// ByAddress mocks base method.
func (m *MockIMessage) ByAddress(ctx context.Context, id uint64, filters storage.AddressMsgsFilter) ([]storage.AddressMessageWithTx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByAddress", ctx, id, filters)
	ret0, _ := ret[0].([]storage.AddressMessageWithTx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByAddress indicates an expected call of ByAddress.
func (mr *MockIMessageMockRecorder) ByAddress(ctx, id, filters any) *IMessageByAddressCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByAddress", reflect.TypeOf((*MockIMessage)(nil).ByAddress), ctx, id, filters)
	return &IMessageByAddressCall{Call: call}
}

// IMessageByAddressCall wrap *gomock.Call
type IMessageByAddressCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IMessageByAddressCall) Return(arg0 []storage.AddressMessageWithTx, arg1 error) *IMessageByAddressCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IMessageByAddressCall) Do(f func(context.Context, uint64, storage.AddressMsgsFilter) ([]storage.AddressMessageWithTx, error)) *IMessageByAddressCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IMessageByAddressCall) DoAndReturn(f func(context.Context, uint64, storage.AddressMsgsFilter) ([]storage.AddressMessageWithTx, error)) *IMessageByAddressCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ByTxId mocks base method.
func (m *MockIMessage) ByTxId(ctx context.Context, txId uint64, limit, offset int) ([]storage.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByTxId", ctx, txId, limit, offset)
	ret0, _ := ret[0].([]storage.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByTxId indicates an expected call of ByTxId.
func (mr *MockIMessageMockRecorder) ByTxId(ctx, txId, limit, offset any) *IMessageByTxIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByTxId", reflect.TypeOf((*MockIMessage)(nil).ByTxId), ctx, txId, limit, offset)
	return &IMessageByTxIdCall{Call: call}
}

// IMessageByTxIdCall wrap *gomock.Call
type IMessageByTxIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IMessageByTxIdCall) Return(arg0 []storage.Message, arg1 error) *IMessageByTxIdCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IMessageByTxIdCall) Do(f func(context.Context, uint64, int, int) ([]storage.Message, error)) *IMessageByTxIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IMessageByTxIdCall) DoAndReturn(f func(context.Context, uint64, int, int) ([]storage.Message, error)) *IMessageByTxIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CursorList mocks base method.
func (m *MockIMessage) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIMessageMockRecorder) CursorList(ctx, id, limit, order, cmp any) *IMessageCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIMessage)(nil).CursorList), ctx, id, limit, order, cmp)
	return &IMessageCursorListCall{Call: call}
}

// IMessageCursorListCall wrap *gomock.Call
type IMessageCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IMessageCursorListCall) Return(arg0 []*storage.Message, arg1 error) *IMessageCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IMessageCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Message, error)) *IMessageCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IMessageCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Message, error)) *IMessageCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIMessage) GetByID(ctx context.Context, id uint64) (*storage.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIMessageMockRecorder) GetByID(ctx, id any) *IMessageGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIMessage)(nil).GetByID), ctx, id)
	return &IMessageGetByIDCall{Call: call}
}

// IMessageGetByIDCall wrap *gomock.Call
type IMessageGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IMessageGetByIDCall) Return(arg0 *storage.Message, arg1 error) *IMessageGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IMessageGetByIDCall) Do(f func(context.Context, uint64) (*storage.Message, error)) *IMessageGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IMessageGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.Message, error)) *IMessageGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIMessage) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIMessageMockRecorder) IsNoRows(err any) *IMessageIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIMessage)(nil).IsNoRows), err)
	return &IMessageIsNoRowsCall{Call: call}
}

// IMessageIsNoRowsCall wrap *gomock.Call
type IMessageIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IMessageIsNoRowsCall) Return(arg0 bool) *IMessageIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IMessageIsNoRowsCall) Do(f func(error) bool) *IMessageIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IMessageIsNoRowsCall) DoAndReturn(f func(error) bool) *IMessageIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIMessage) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIMessageMockRecorder) LastID(ctx any) *IMessageLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIMessage)(nil).LastID), ctx)
	return &IMessageLastIDCall{Call: call}
}

// IMessageLastIDCall wrap *gomock.Call
type IMessageLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IMessageLastIDCall) Return(arg0 uint64, arg1 error) *IMessageLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IMessageLastIDCall) Do(f func(context.Context) (uint64, error)) *IMessageLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IMessageLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *IMessageLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIMessage) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIMessageMockRecorder) List(ctx, limit, offset, order any) *IMessageListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIMessage)(nil).List), ctx, limit, offset, order)
	return &IMessageListCall{Call: call}
}

// IMessageListCall wrap *gomock.Call
type IMessageListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IMessageListCall) Return(arg0 []*storage.Message, arg1 error) *IMessageListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IMessageListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Message, error)) *IMessageListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IMessageListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Message, error)) *IMessageListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ListWithTx mocks base method.
func (m *MockIMessage) ListWithTx(ctx context.Context, filters storage.MessageListWithTxFilters) ([]storage.MessageWithTx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWithTx", ctx, filters)
	ret0, _ := ret[0].([]storage.MessageWithTx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWithTx indicates an expected call of ListWithTx.
func (mr *MockIMessageMockRecorder) ListWithTx(ctx, filters any) *IMessageListWithTxCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWithTx", reflect.TypeOf((*MockIMessage)(nil).ListWithTx), ctx, filters)
	return &IMessageListWithTxCall{Call: call}
}

// IMessageListWithTxCall wrap *gomock.Call
type IMessageListWithTxCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IMessageListWithTxCall) Return(arg0 []storage.MessageWithTx, arg1 error) *IMessageListWithTxCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IMessageListWithTxCall) Do(f func(context.Context, storage.MessageListWithTxFilters) ([]storage.MessageWithTx, error)) *IMessageListWithTxCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IMessageListWithTxCall) DoAndReturn(f func(context.Context, storage.MessageListWithTxFilters) ([]storage.MessageWithTx, error)) *IMessageListWithTxCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIMessage) Save(ctx context.Context, m *storage.Message) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIMessageMockRecorder) Save(ctx, m any) *IMessageSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIMessage)(nil).Save), ctx, m)
	return &IMessageSaveCall{Call: call}
}

// IMessageSaveCall wrap *gomock.Call
type IMessageSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IMessageSaveCall) Return(arg0 error) *IMessageSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IMessageSaveCall) Do(f func(context.Context, *storage.Message) error) *IMessageSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IMessageSaveCall) DoAndReturn(f func(context.Context, *storage.Message) error) *IMessageSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIMessage) Update(ctx context.Context, m *storage.Message) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIMessageMockRecorder) Update(ctx, m any) *IMessageUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIMessage)(nil).Update), ctx, m)
	return &IMessageUpdateCall{Call: call}
}

// IMessageUpdateCall wrap *gomock.Call
type IMessageUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IMessageUpdateCall) Return(arg0 error) *IMessageUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IMessageUpdateCall) Do(f func(context.Context, *storage.Message) error) *IMessageUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IMessageUpdateCall) DoAndReturn(f func(context.Context, *storage.Message) error) *IMessageUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
