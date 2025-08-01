// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: hl_mailbox.go
//
// Generated by this command:
//
//	mockgen -source=hl_mailbox.go -destination=mock/hl_mailbox.go -package=mock -typed
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	storage "github.com/celenium-io/celestia-indexer/internal/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockIHLMailbox is a mock of IHLMailbox interface.
type MockIHLMailbox struct {
	ctrl     *gomock.Controller
	recorder *MockIHLMailboxMockRecorder
}

// MockIHLMailboxMockRecorder is the mock recorder for MockIHLMailbox.
type MockIHLMailboxMockRecorder struct {
	mock *MockIHLMailbox
}

// NewMockIHLMailbox creates a new mock instance.
func NewMockIHLMailbox(ctrl *gomock.Controller) *MockIHLMailbox {
	mock := &MockIHLMailbox{ctrl: ctrl}
	mock.recorder = &MockIHLMailboxMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIHLMailbox) EXPECT() *MockIHLMailboxMockRecorder {
	return m.recorder
}

// ByHash mocks base method.
func (m *MockIHLMailbox) ByHash(ctx context.Context, hash []byte) (storage.HLMailbox, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByHash", ctx, hash)
	ret0, _ := ret[0].(storage.HLMailbox)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByHash indicates an expected call of ByHash.
func (mr *MockIHLMailboxMockRecorder) ByHash(ctx, hash any) *MockIHLMailboxByHashCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByHash", reflect.TypeOf((*MockIHLMailbox)(nil).ByHash), ctx, hash)
	return &MockIHLMailboxByHashCall{Call: call}
}

// MockIHLMailboxByHashCall wrap *gomock.Call
type MockIHLMailboxByHashCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIHLMailboxByHashCall) Return(arg0 storage.HLMailbox, arg1 error) *MockIHLMailboxByHashCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIHLMailboxByHashCall) Do(f func(context.Context, []byte) (storage.HLMailbox, error)) *MockIHLMailboxByHashCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIHLMailboxByHashCall) DoAndReturn(f func(context.Context, []byte) (storage.HLMailbox, error)) *MockIHLMailboxByHashCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIHLMailbox) List(ctx context.Context, limit, offset int) ([]storage.HLMailbox, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset)
	ret0, _ := ret[0].([]storage.HLMailbox)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIHLMailboxMockRecorder) List(ctx, limit, offset any) *MockIHLMailboxListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIHLMailbox)(nil).List), ctx, limit, offset)
	return &MockIHLMailboxListCall{Call: call}
}

// MockIHLMailboxListCall wrap *gomock.Call
type MockIHLMailboxListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIHLMailboxListCall) Return(arg0 []storage.HLMailbox, arg1 error) *MockIHLMailboxListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIHLMailboxListCall) Do(f func(context.Context, int, int) ([]storage.HLMailbox, error)) *MockIHLMailboxListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIHLMailboxListCall) DoAndReturn(f func(context.Context, int, int) ([]storage.HLMailbox, error)) *MockIHLMailboxListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
