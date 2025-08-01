// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: ibc_channel.go
//
// Generated by this command:
//
//	mockgen -source=ibc_channel.go -destination=mock/ibc_channel.go -package=mock -typed
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	storage "github.com/celenium-io/celestia-indexer/internal/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockIIbcChannel is a mock of IIbcChannel interface.
type MockIIbcChannel struct {
	ctrl     *gomock.Controller
	recorder *MockIIbcChannelMockRecorder
	isgomock struct{}
}

// MockIIbcChannelMockRecorder is the mock recorder for MockIIbcChannel.
type MockIIbcChannelMockRecorder struct {
	mock *MockIIbcChannel
}

// NewMockIIbcChannel creates a new mock instance.
func NewMockIIbcChannel(ctrl *gomock.Controller) *MockIIbcChannel {
	mock := &MockIIbcChannel{ctrl: ctrl}
	mock.recorder = &MockIIbcChannelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIIbcChannel) EXPECT() *MockIIbcChannelMockRecorder {
	return m.recorder
}

// BusiestChannel1m mocks base method.
func (m *MockIIbcChannel) BusiestChannel1m(ctx context.Context) (storage.BusiestChannel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BusiestChannel1m", ctx)
	ret0, _ := ret[0].(storage.BusiestChannel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BusiestChannel1m indicates an expected call of BusiestChannel1m.
func (mr *MockIIbcChannelMockRecorder) BusiestChannel1m(ctx any) *MockIIbcChannelBusiestChannel1mCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BusiestChannel1m", reflect.TypeOf((*MockIIbcChannel)(nil).BusiestChannel1m), ctx)
	return &MockIIbcChannelBusiestChannel1mCall{Call: call}
}

// MockIIbcChannelBusiestChannel1mCall wrap *gomock.Call
type MockIIbcChannelBusiestChannel1mCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIIbcChannelBusiestChannel1mCall) Return(arg0 storage.BusiestChannel, arg1 error) *MockIIbcChannelBusiestChannel1mCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIIbcChannelBusiestChannel1mCall) Do(f func(context.Context) (storage.BusiestChannel, error)) *MockIIbcChannelBusiestChannel1mCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIIbcChannelBusiestChannel1mCall) DoAndReturn(f func(context.Context) (storage.BusiestChannel, error)) *MockIIbcChannelBusiestChannel1mCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ById mocks base method.
func (m *MockIIbcChannel) ById(ctx context.Context, id string) (storage.IbcChannel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ById", ctx, id)
	ret0, _ := ret[0].(storage.IbcChannel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ById indicates an expected call of ById.
func (mr *MockIIbcChannelMockRecorder) ById(ctx, id any) *MockIIbcChannelByIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ById", reflect.TypeOf((*MockIIbcChannel)(nil).ById), ctx, id)
	return &MockIIbcChannelByIdCall{Call: call}
}

// MockIIbcChannelByIdCall wrap *gomock.Call
type MockIIbcChannelByIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIIbcChannelByIdCall) Return(arg0 storage.IbcChannel, arg1 error) *MockIIbcChannelByIdCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIIbcChannelByIdCall) Do(f func(context.Context, string) (storage.IbcChannel, error)) *MockIIbcChannelByIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIIbcChannelByIdCall) DoAndReturn(f func(context.Context, string) (storage.IbcChannel, error)) *MockIIbcChannelByIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIIbcChannel) List(ctx context.Context, fltrs storage.ListChannelFilters) ([]storage.IbcChannel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, fltrs)
	ret0, _ := ret[0].([]storage.IbcChannel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIIbcChannelMockRecorder) List(ctx, fltrs any) *MockIIbcChannelListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIIbcChannel)(nil).List), ctx, fltrs)
	return &MockIIbcChannelListCall{Call: call}
}

// MockIIbcChannelListCall wrap *gomock.Call
type MockIIbcChannelListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIIbcChannelListCall) Return(arg0 []storage.IbcChannel, arg1 error) *MockIIbcChannelListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIIbcChannelListCall) Do(f func(context.Context, storage.ListChannelFilters) ([]storage.IbcChannel, error)) *MockIIbcChannelListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIIbcChannelListCall) DoAndReturn(f func(context.Context, storage.ListChannelFilters) ([]storage.IbcChannel, error)) *MockIIbcChannelListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// StatsByChain mocks base method.
func (m *MockIIbcChannel) StatsByChain(ctx context.Context, limit, offset int) ([]storage.ChainStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StatsByChain", ctx, limit, offset)
	ret0, _ := ret[0].([]storage.ChainStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StatsByChain indicates an expected call of StatsByChain.
func (mr *MockIIbcChannelMockRecorder) StatsByChain(ctx, limit, offset any) *MockIIbcChannelStatsByChainCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StatsByChain", reflect.TypeOf((*MockIIbcChannel)(nil).StatsByChain), ctx, limit, offset)
	return &MockIIbcChannelStatsByChainCall{Call: call}
}

// MockIIbcChannelStatsByChainCall wrap *gomock.Call
type MockIIbcChannelStatsByChainCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIIbcChannelStatsByChainCall) Return(arg0 []storage.ChainStats, arg1 error) *MockIIbcChannelStatsByChainCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIIbcChannelStatsByChainCall) Do(f func(context.Context, int, int) ([]storage.ChainStats, error)) *MockIIbcChannelStatsByChainCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIIbcChannelStatsByChainCall) DoAndReturn(f func(context.Context, int, int) ([]storage.ChainStats, error)) *MockIIbcChannelStatsByChainCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
