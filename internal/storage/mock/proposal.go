// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: proposal.go
//
// Generated by this command:
//
//	mockgen -source=proposal.go -destination=mock/proposal.go -package=mock -typed
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

// MockIProposal is a mock of IProposal interface.
type MockIProposal struct {
	ctrl     *gomock.Controller
	recorder *MockIProposalMockRecorder
}

// MockIProposalMockRecorder is the mock recorder for MockIProposal.
type MockIProposalMockRecorder struct {
	mock *MockIProposal
}

// NewMockIProposal creates a new mock instance.
func NewMockIProposal(ctrl *gomock.Controller) *MockIProposal {
	mock := &MockIProposal{ctrl: ctrl}
	mock.recorder = &MockIProposalMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIProposal) EXPECT() *MockIProposalMockRecorder {
	return m.recorder
}

// CursorList mocks base method.
func (m *MockIProposal) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.Proposal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.Proposal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIProposalMockRecorder) CursorList(ctx, id, limit, order, cmp any) *MockIProposalCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIProposal)(nil).CursorList), ctx, id, limit, order, cmp)
	return &MockIProposalCursorListCall{Call: call}
}

// MockIProposalCursorListCall wrap *gomock.Call
type MockIProposalCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIProposalCursorListCall) Return(arg0 []*storage.Proposal, arg1 error) *MockIProposalCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIProposalCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Proposal, error)) *MockIProposalCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIProposalCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Proposal, error)) *MockIProposalCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIProposal) GetByID(ctx context.Context, id uint64) (*storage.Proposal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.Proposal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIProposalMockRecorder) GetByID(ctx, id any) *MockIProposalGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIProposal)(nil).GetByID), ctx, id)
	return &MockIProposalGetByIDCall{Call: call}
}

// MockIProposalGetByIDCall wrap *gomock.Call
type MockIProposalGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIProposalGetByIDCall) Return(arg0 *storage.Proposal, arg1 error) *MockIProposalGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIProposalGetByIDCall) Do(f func(context.Context, uint64) (*storage.Proposal, error)) *MockIProposalGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIProposalGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.Proposal, error)) *MockIProposalGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIProposal) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIProposalMockRecorder) IsNoRows(err any) *MockIProposalIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIProposal)(nil).IsNoRows), err)
	return &MockIProposalIsNoRowsCall{Call: call}
}

// MockIProposalIsNoRowsCall wrap *gomock.Call
type MockIProposalIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIProposalIsNoRowsCall) Return(arg0 bool) *MockIProposalIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIProposalIsNoRowsCall) Do(f func(error) bool) *MockIProposalIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIProposalIsNoRowsCall) DoAndReturn(f func(error) bool) *MockIProposalIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIProposal) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIProposalMockRecorder) LastID(ctx any) *MockIProposalLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIProposal)(nil).LastID), ctx)
	return &MockIProposalLastIDCall{Call: call}
}

// MockIProposalLastIDCall wrap *gomock.Call
type MockIProposalLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIProposalLastIDCall) Return(arg0 uint64, arg1 error) *MockIProposalLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIProposalLastIDCall) Do(f func(context.Context) (uint64, error)) *MockIProposalLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIProposalLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *MockIProposalLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIProposal) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.Proposal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.Proposal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIProposalMockRecorder) List(ctx, limit, offset, order any) *MockIProposalListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIProposal)(nil).List), ctx, limit, offset, order)
	return &MockIProposalListCall{Call: call}
}

// MockIProposalListCall wrap *gomock.Call
type MockIProposalListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIProposalListCall) Return(arg0 []*storage.Proposal, arg1 error) *MockIProposalListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIProposalListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Proposal, error)) *MockIProposalListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIProposalListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Proposal, error)) *MockIProposalListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ListWithFilters mocks base method.
func (m *MockIProposal) ListWithFilters(ctx context.Context, filters storage.ListProposalFilters) ([]storage.Proposal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWithFilters", ctx, filters)
	ret0, _ := ret[0].([]storage.Proposal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWithFilters indicates an expected call of ListWithFilters.
func (mr *MockIProposalMockRecorder) ListWithFilters(ctx, filters any) *MockIProposalListWithFiltersCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWithFilters", reflect.TypeOf((*MockIProposal)(nil).ListWithFilters), ctx, filters)
	return &MockIProposalListWithFiltersCall{Call: call}
}

// MockIProposalListWithFiltersCall wrap *gomock.Call
type MockIProposalListWithFiltersCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIProposalListWithFiltersCall) Return(proposals []storage.Proposal, err error) *MockIProposalListWithFiltersCall {
	c.Call = c.Call.Return(proposals, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIProposalListWithFiltersCall) Do(f func(context.Context, storage.ListProposalFilters) ([]storage.Proposal, error)) *MockIProposalListWithFiltersCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIProposalListWithFiltersCall) DoAndReturn(f func(context.Context, storage.ListProposalFilters) ([]storage.Proposal, error)) *MockIProposalListWithFiltersCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIProposal) Save(ctx context.Context, m *storage.Proposal) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIProposalMockRecorder) Save(ctx, m any) *MockIProposalSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIProposal)(nil).Save), ctx, m)
	return &MockIProposalSaveCall{Call: call}
}

// MockIProposalSaveCall wrap *gomock.Call
type MockIProposalSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIProposalSaveCall) Return(arg0 error) *MockIProposalSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIProposalSaveCall) Do(f func(context.Context, *storage.Proposal) error) *MockIProposalSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIProposalSaveCall) DoAndReturn(f func(context.Context, *storage.Proposal) error) *MockIProposalSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIProposal) Update(ctx context.Context, m *storage.Proposal) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIProposalMockRecorder) Update(ctx, m any) *MockIProposalUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIProposal)(nil).Update), ctx, m)
	return &MockIProposalUpdateCall{Call: call}
}

// MockIProposalUpdateCall wrap *gomock.Call
type MockIProposalUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIProposalUpdateCall) Return(arg0 error) *MockIProposalUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIProposalUpdateCall) Do(f func(context.Context, *storage.Proposal) error) *MockIProposalUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIProposalUpdateCall) DoAndReturn(f func(context.Context, *storage.Proposal) error) *MockIProposalUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
