// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: api_key.go
//
// Generated by this command:
//
//	mockgen -source=api_key.go -destination=mock/api_key.go -package=mock -typed
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	storage "github.com/celenium-io/celestia-indexer/internal/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockIApiKey is a mock of IApiKey interface.
type MockIApiKey struct {
	ctrl     *gomock.Controller
	recorder *MockIApiKeyMockRecorder
}

// MockIApiKeyMockRecorder is the mock recorder for MockIApiKey.
type MockIApiKeyMockRecorder struct {
	mock *MockIApiKey
}

// NewMockIApiKey creates a new mock instance.
func NewMockIApiKey(ctrl *gomock.Controller) *MockIApiKey {
	mock := &MockIApiKey{ctrl: ctrl}
	mock.recorder = &MockIApiKeyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIApiKey) EXPECT() *MockIApiKeyMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockIApiKey) Get(ctx context.Context, key string) (storage.ApiKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(storage.ApiKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockIApiKeyMockRecorder) Get(ctx, key any) *MockIApiKeyGetCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockIApiKey)(nil).Get), ctx, key)
	return &MockIApiKeyGetCall{Call: call}
}

// MockIApiKeyGetCall wrap *gomock.Call
type MockIApiKeyGetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIApiKeyGetCall) Return(arg0 storage.ApiKey, arg1 error) *MockIApiKeyGetCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIApiKeyGetCall) Do(f func(context.Context, string) (storage.ApiKey, error)) *MockIApiKeyGetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIApiKeyGetCall) DoAndReturn(f func(context.Context, string) (storage.ApiKey, error)) *MockIApiKeyGetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
