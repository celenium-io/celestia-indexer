// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: rollup.go
//
// Generated by this command:
//
//	mockgen -source=rollup.go -destination=mock/rollup.go -package=mock -typed
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

// MockIRollup is a mock of IRollup interface.
type MockIRollup struct {
	ctrl     *gomock.Controller
	recorder *MockIRollupMockRecorder
}

// MockIRollupMockRecorder is the mock recorder for MockIRollup.
type MockIRollupMockRecorder struct {
	mock *MockIRollup
}

// NewMockIRollup creates a new mock instance.
func NewMockIRollup(ctrl *gomock.Controller) *MockIRollup {
	mock := &MockIRollup{ctrl: ctrl}
	mock.recorder = &MockIRollupMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIRollup) EXPECT() *MockIRollupMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockIRollup) Count(ctx context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockIRollupMockRecorder) Count(ctx any) *IRollupCountCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockIRollup)(nil).Count), ctx)
	return &IRollupCountCall{Call: call}
}

// IRollupCountCall wrap *gomock.Call
type IRollupCountCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IRollupCountCall) Return(arg0 int64, arg1 error) *IRollupCountCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IRollupCountCall) Do(f func(context.Context) (int64, error)) *IRollupCountCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IRollupCountCall) DoAndReturn(f func(context.Context) (int64, error)) *IRollupCountCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CursorList mocks base method.
func (m *MockIRollup) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.Rollup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.Rollup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIRollupMockRecorder) CursorList(ctx, id, limit, order, cmp any) *IRollupCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIRollup)(nil).CursorList), ctx, id, limit, order, cmp)
	return &IRollupCursorListCall{Call: call}
}

// IRollupCursorListCall wrap *gomock.Call
type IRollupCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IRollupCursorListCall) Return(arg0 []*storage.Rollup, arg1 error) *IRollupCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IRollupCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Rollup, error)) *IRollupCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IRollupCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Rollup, error)) *IRollupCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIRollup) GetByID(ctx context.Context, id uint64) (*storage.Rollup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.Rollup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIRollupMockRecorder) GetByID(ctx, id any) *IRollupGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIRollup)(nil).GetByID), ctx, id)
	return &IRollupGetByIDCall{Call: call}
}

// IRollupGetByIDCall wrap *gomock.Call
type IRollupGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IRollupGetByIDCall) Return(arg0 *storage.Rollup, arg1 error) *IRollupGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IRollupGetByIDCall) Do(f func(context.Context, uint64) (*storage.Rollup, error)) *IRollupGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IRollupGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.Rollup, error)) *IRollupGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIRollup) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIRollupMockRecorder) IsNoRows(err any) *IRollupIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIRollup)(nil).IsNoRows), err)
	return &IRollupIsNoRowsCall{Call: call}
}

// IRollupIsNoRowsCall wrap *gomock.Call
type IRollupIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IRollupIsNoRowsCall) Return(arg0 bool) *IRollupIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IRollupIsNoRowsCall) Do(f func(error) bool) *IRollupIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IRollupIsNoRowsCall) DoAndReturn(f func(error) bool) *IRollupIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIRollup) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIRollupMockRecorder) LastID(ctx any) *IRollupLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIRollup)(nil).LastID), ctx)
	return &IRollupLastIDCall{Call: call}
}

// IRollupLastIDCall wrap *gomock.Call
type IRollupLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IRollupLastIDCall) Return(arg0 uint64, arg1 error) *IRollupLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IRollupLastIDCall) Do(f func(context.Context) (uint64, error)) *IRollupLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IRollupLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *IRollupLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Leaderboard mocks base method.
func (m *MockIRollup) Leaderboard(ctx context.Context, sortField string, sort storage0.SortOrder, limit, offset int) ([]storage.RollupWithStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Leaderboard", ctx, sortField, sort, limit, offset)
	ret0, _ := ret[0].([]storage.RollupWithStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Leaderboard indicates an expected call of Leaderboard.
func (mr *MockIRollupMockRecorder) Leaderboard(ctx, sortField, sort, limit, offset any) *IRollupLeaderboardCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Leaderboard", reflect.TypeOf((*MockIRollup)(nil).Leaderboard), ctx, sortField, sort, limit, offset)
	return &IRollupLeaderboardCall{Call: call}
}

// IRollupLeaderboardCall wrap *gomock.Call
type IRollupLeaderboardCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IRollupLeaderboardCall) Return(arg0 []storage.RollupWithStats, arg1 error) *IRollupLeaderboardCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IRollupLeaderboardCall) Do(f func(context.Context, string, storage0.SortOrder, int, int) ([]storage.RollupWithStats, error)) *IRollupLeaderboardCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IRollupLeaderboardCall) DoAndReturn(f func(context.Context, string, storage0.SortOrder, int, int) ([]storage.RollupWithStats, error)) *IRollupLeaderboardCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIRollup) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.Rollup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.Rollup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIRollupMockRecorder) List(ctx, limit, offset, order any) *IRollupListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIRollup)(nil).List), ctx, limit, offset, order)
	return &IRollupListCall{Call: call}
}

// IRollupListCall wrap *gomock.Call
type IRollupListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IRollupListCall) Return(arg0 []*storage.Rollup, arg1 error) *IRollupListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IRollupListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Rollup, error)) *IRollupListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IRollupListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Rollup, error)) *IRollupListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Namespaces mocks base method.
func (m *MockIRollup) Namespaces(ctx context.Context, rollupId uint64, limit, offset int) ([]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Namespaces", ctx, rollupId, limit, offset)
	ret0, _ := ret[0].([]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Namespaces indicates an expected call of Namespaces.
func (mr *MockIRollupMockRecorder) Namespaces(ctx, rollupId, limit, offset any) *IRollupNamespacesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Namespaces", reflect.TypeOf((*MockIRollup)(nil).Namespaces), ctx, rollupId, limit, offset)
	return &IRollupNamespacesCall{Call: call}
}

// IRollupNamespacesCall wrap *gomock.Call
type IRollupNamespacesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IRollupNamespacesCall) Return(namespaceIds []uint64, err error) *IRollupNamespacesCall {
	c.Call = c.Call.Return(namespaceIds, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IRollupNamespacesCall) Do(f func(context.Context, uint64, int, int) ([]uint64, error)) *IRollupNamespacesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IRollupNamespacesCall) DoAndReturn(f func(context.Context, uint64, int, int) ([]uint64, error)) *IRollupNamespacesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Providers mocks base method.
func (m *MockIRollup) Providers(ctx context.Context, rollupId uint64) ([]storage.RollupProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Providers", ctx, rollupId)
	ret0, _ := ret[0].([]storage.RollupProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Providers indicates an expected call of Providers.
func (mr *MockIRollupMockRecorder) Providers(ctx, rollupId any) *IRollupProvidersCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Providers", reflect.TypeOf((*MockIRollup)(nil).Providers), ctx, rollupId)
	return &IRollupProvidersCall{Call: call}
}

// IRollupProvidersCall wrap *gomock.Call
type IRollupProvidersCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IRollupProvidersCall) Return(providers []storage.RollupProvider, err error) *IRollupProvidersCall {
	c.Call = c.Call.Return(providers, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IRollupProvidersCall) Do(f func(context.Context, uint64) ([]storage.RollupProvider, error)) *IRollupProvidersCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IRollupProvidersCall) DoAndReturn(f func(context.Context, uint64) ([]storage.RollupProvider, error)) *IRollupProvidersCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIRollup) Save(ctx context.Context, m *storage.Rollup) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIRollupMockRecorder) Save(ctx, m any) *IRollupSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIRollup)(nil).Save), ctx, m)
	return &IRollupSaveCall{Call: call}
}

// IRollupSaveCall wrap *gomock.Call
type IRollupSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IRollupSaveCall) Return(arg0 error) *IRollupSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IRollupSaveCall) Do(f func(context.Context, *storage.Rollup) error) *IRollupSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IRollupSaveCall) DoAndReturn(f func(context.Context, *storage.Rollup) error) *IRollupSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Series mocks base method.
func (m *MockIRollup) Series(ctx context.Context, rollupId uint64, timeframe, column string, req storage.SeriesRequest) ([]storage.HistogramItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Series", ctx, rollupId, timeframe, column, req)
	ret0, _ := ret[0].([]storage.HistogramItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Series indicates an expected call of Series.
func (mr *MockIRollupMockRecorder) Series(ctx, rollupId, timeframe, column, req any) *IRollupSeriesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Series", reflect.TypeOf((*MockIRollup)(nil).Series), ctx, rollupId, timeframe, column, req)
	return &IRollupSeriesCall{Call: call}
}

// IRollupSeriesCall wrap *gomock.Call
type IRollupSeriesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IRollupSeriesCall) Return(items []storage.HistogramItem, err error) *IRollupSeriesCall {
	c.Call = c.Call.Return(items, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IRollupSeriesCall) Do(f func(context.Context, uint64, string, string, storage.SeriesRequest) ([]storage.HistogramItem, error)) *IRollupSeriesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IRollupSeriesCall) DoAndReturn(f func(context.Context, uint64, string, string, storage.SeriesRequest) ([]storage.HistogramItem, error)) *IRollupSeriesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Stats mocks base method.
func (m *MockIRollup) Stats(ctx context.Context, rollupId uint64) (storage.RollupStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stats", ctx, rollupId)
	ret0, _ := ret[0].(storage.RollupStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stats indicates an expected call of Stats.
func (mr *MockIRollupMockRecorder) Stats(ctx, rollupId any) *IRollupStatsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stats", reflect.TypeOf((*MockIRollup)(nil).Stats), ctx, rollupId)
	return &IRollupStatsCall{Call: call}
}

// IRollupStatsCall wrap *gomock.Call
type IRollupStatsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IRollupStatsCall) Return(arg0 storage.RollupStats, arg1 error) *IRollupStatsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IRollupStatsCall) Do(f func(context.Context, uint64) (storage.RollupStats, error)) *IRollupStatsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IRollupStatsCall) DoAndReturn(f func(context.Context, uint64) (storage.RollupStats, error)) *IRollupStatsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIRollup) Update(ctx context.Context, m *storage.Rollup) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIRollupMockRecorder) Update(ctx, m any) *IRollupUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIRollup)(nil).Update), ctx, m)
	return &IRollupUpdateCall{Call: call}
}

// IRollupUpdateCall wrap *gomock.Call
type IRollupUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IRollupUpdateCall) Return(arg0 error) *IRollupUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IRollupUpdateCall) Do(f func(context.Context, *storage.Rollup) error) *IRollupUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IRollupUpdateCall) DoAndReturn(f func(context.Context, *storage.Rollup) error) *IRollupUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
