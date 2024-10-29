// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
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

// AllSeries mocks base method.
func (m *MockIRollup) AllSeries(ctx context.Context) ([]storage.RollupHistogramItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllSeries", ctx)
	ret0, _ := ret[0].([]storage.RollupHistogramItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllSeries indicates an expected call of AllSeries.
func (mr *MockIRollupMockRecorder) AllSeries(ctx any) *MockIRollupAllSeriesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllSeries", reflect.TypeOf((*MockIRollup)(nil).AllSeries), ctx)
	return &MockIRollupAllSeriesCall{Call: call}
}

// MockIRollupAllSeriesCall wrap *gomock.Call
type MockIRollupAllSeriesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupAllSeriesCall) Return(arg0 []storage.RollupHistogramItem, arg1 error) *MockIRollupAllSeriesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupAllSeriesCall) Do(f func(context.Context) ([]storage.RollupHistogramItem, error)) *MockIRollupAllSeriesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupAllSeriesCall) DoAndReturn(f func(context.Context) ([]storage.RollupHistogramItem, error)) *MockIRollupAllSeriesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ById mocks base method.
func (m *MockIRollup) ById(ctx context.Context, rollupId uint64) (storage.RollupWithStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ById", ctx, rollupId)
	ret0, _ := ret[0].(storage.RollupWithStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ById indicates an expected call of ById.
func (mr *MockIRollupMockRecorder) ById(ctx, rollupId any) *MockIRollupByIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ById", reflect.TypeOf((*MockIRollup)(nil).ById), ctx, rollupId)
	return &MockIRollupByIdCall{Call: call}
}

// MockIRollupByIdCall wrap *gomock.Call
type MockIRollupByIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupByIdCall) Return(arg0 storage.RollupWithStats, arg1 error) *MockIRollupByIdCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupByIdCall) Do(f func(context.Context, uint64) (storage.RollupWithStats, error)) *MockIRollupByIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupByIdCall) DoAndReturn(f func(context.Context, uint64) (storage.RollupWithStats, error)) *MockIRollupByIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// BySlug mocks base method.
func (m *MockIRollup) BySlug(ctx context.Context, slug string) (storage.RollupWithStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BySlug", ctx, slug)
	ret0, _ := ret[0].(storage.RollupWithStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BySlug indicates an expected call of BySlug.
func (mr *MockIRollupMockRecorder) BySlug(ctx, slug any) *MockIRollupBySlugCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BySlug", reflect.TypeOf((*MockIRollup)(nil).BySlug), ctx, slug)
	return &MockIRollupBySlugCall{Call: call}
}

// MockIRollupBySlugCall wrap *gomock.Call
type MockIRollupBySlugCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupBySlugCall) Return(arg0 storage.RollupWithStats, arg1 error) *MockIRollupBySlugCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupBySlugCall) Do(f func(context.Context, string) (storage.RollupWithStats, error)) *MockIRollupBySlugCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupBySlugCall) DoAndReturn(f func(context.Context, string) (storage.RollupWithStats, error)) *MockIRollupBySlugCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
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
func (mr *MockIRollupMockRecorder) Count(ctx any) *MockIRollupCountCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockIRollup)(nil).Count), ctx)
	return &MockIRollupCountCall{Call: call}
}

// MockIRollupCountCall wrap *gomock.Call
type MockIRollupCountCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupCountCall) Return(arg0 int64, arg1 error) *MockIRollupCountCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupCountCall) Do(f func(context.Context) (int64, error)) *MockIRollupCountCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupCountCall) DoAndReturn(f func(context.Context) (int64, error)) *MockIRollupCountCall {
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
func (mr *MockIRollupMockRecorder) CursorList(ctx, id, limit, order, cmp any) *MockIRollupCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIRollup)(nil).CursorList), ctx, id, limit, order, cmp)
	return &MockIRollupCursorListCall{Call: call}
}

// MockIRollupCursorListCall wrap *gomock.Call
type MockIRollupCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupCursorListCall) Return(arg0 []*storage.Rollup, arg1 error) *MockIRollupCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Rollup, error)) *MockIRollupCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Rollup, error)) *MockIRollupCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Distribution mocks base method.
func (m *MockIRollup) Distribution(ctx context.Context, rollupId uint64, series, groupBy string) ([]storage.DistributionItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Distribution", ctx, rollupId, series, groupBy)
	ret0, _ := ret[0].([]storage.DistributionItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Distribution indicates an expected call of Distribution.
func (mr *MockIRollupMockRecorder) Distribution(ctx, rollupId, series, groupBy any) *MockIRollupDistributionCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Distribution", reflect.TypeOf((*MockIRollup)(nil).Distribution), ctx, rollupId, series, groupBy)
	return &MockIRollupDistributionCall{Call: call}
}

// MockIRollupDistributionCall wrap *gomock.Call
type MockIRollupDistributionCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupDistributionCall) Return(items []storage.DistributionItem, err error) *MockIRollupDistributionCall {
	c.Call = c.Call.Return(items, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupDistributionCall) Do(f func(context.Context, uint64, string, string) ([]storage.DistributionItem, error)) *MockIRollupDistributionCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupDistributionCall) DoAndReturn(f func(context.Context, uint64, string, string) ([]storage.DistributionItem, error)) *MockIRollupDistributionCall {
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
func (mr *MockIRollupMockRecorder) GetByID(ctx, id any) *MockIRollupGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIRollup)(nil).GetByID), ctx, id)
	return &MockIRollupGetByIDCall{Call: call}
}

// MockIRollupGetByIDCall wrap *gomock.Call
type MockIRollupGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupGetByIDCall) Return(arg0 *storage.Rollup, arg1 error) *MockIRollupGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupGetByIDCall) Do(f func(context.Context, uint64) (*storage.Rollup, error)) *MockIRollupGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.Rollup, error)) *MockIRollupGetByIDCall {
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
func (mr *MockIRollupMockRecorder) IsNoRows(err any) *MockIRollupIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIRollup)(nil).IsNoRows), err)
	return &MockIRollupIsNoRowsCall{Call: call}
}

// MockIRollupIsNoRowsCall wrap *gomock.Call
type MockIRollupIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupIsNoRowsCall) Return(arg0 bool) *MockIRollupIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupIsNoRowsCall) Do(f func(error) bool) *MockIRollupIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupIsNoRowsCall) DoAndReturn(f func(error) bool) *MockIRollupIsNoRowsCall {
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
func (mr *MockIRollupMockRecorder) LastID(ctx any) *MockIRollupLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIRollup)(nil).LastID), ctx)
	return &MockIRollupLastIDCall{Call: call}
}

// MockIRollupLastIDCall wrap *gomock.Call
type MockIRollupLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupLastIDCall) Return(arg0 uint64, arg1 error) *MockIRollupLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupLastIDCall) Do(f func(context.Context) (uint64, error)) *MockIRollupLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *MockIRollupLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Leaderboard mocks base method.
func (m *MockIRollup) Leaderboard(ctx context.Context, fltrs storage.LeaderboardFilters) ([]storage.RollupWithStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Leaderboard", ctx, fltrs)
	ret0, _ := ret[0].([]storage.RollupWithStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Leaderboard indicates an expected call of Leaderboard.
func (mr *MockIRollupMockRecorder) Leaderboard(ctx, fltrs any) *MockIRollupLeaderboardCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Leaderboard", reflect.TypeOf((*MockIRollup)(nil).Leaderboard), ctx, fltrs)
	return &MockIRollupLeaderboardCall{Call: call}
}

// MockIRollupLeaderboardCall wrap *gomock.Call
type MockIRollupLeaderboardCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupLeaderboardCall) Return(arg0 []storage.RollupWithStats, arg1 error) *MockIRollupLeaderboardCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupLeaderboardCall) Do(f func(context.Context, storage.LeaderboardFilters) ([]storage.RollupWithStats, error)) *MockIRollupLeaderboardCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupLeaderboardCall) DoAndReturn(f func(context.Context, storage.LeaderboardFilters) ([]storage.RollupWithStats, error)) *MockIRollupLeaderboardCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LeaderboardDay mocks base method.
func (m *MockIRollup) LeaderboardDay(ctx context.Context, fltrs storage.LeaderboardFilters) ([]storage.RollupWithDayStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LeaderboardDay", ctx, fltrs)
	ret0, _ := ret[0].([]storage.RollupWithDayStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LeaderboardDay indicates an expected call of LeaderboardDay.
func (mr *MockIRollupMockRecorder) LeaderboardDay(ctx, fltrs any) *MockIRollupLeaderboardDayCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeaderboardDay", reflect.TypeOf((*MockIRollup)(nil).LeaderboardDay), ctx, fltrs)
	return &MockIRollupLeaderboardDayCall{Call: call}
}

// MockIRollupLeaderboardDayCall wrap *gomock.Call
type MockIRollupLeaderboardDayCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupLeaderboardDayCall) Return(arg0 []storage.RollupWithDayStats, arg1 error) *MockIRollupLeaderboardDayCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupLeaderboardDayCall) Do(f func(context.Context, storage.LeaderboardFilters) ([]storage.RollupWithDayStats, error)) *MockIRollupLeaderboardDayCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupLeaderboardDayCall) DoAndReturn(f func(context.Context, storage.LeaderboardFilters) ([]storage.RollupWithDayStats, error)) *MockIRollupLeaderboardDayCall {
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
func (mr *MockIRollupMockRecorder) List(ctx, limit, offset, order any) *MockIRollupListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIRollup)(nil).List), ctx, limit, offset, order)
	return &MockIRollupListCall{Call: call}
}

// MockIRollupListCall wrap *gomock.Call
type MockIRollupListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupListCall) Return(arg0 []*storage.Rollup, arg1 error) *MockIRollupListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Rollup, error)) *MockIRollupListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Rollup, error)) *MockIRollupListCall {
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
func (mr *MockIRollupMockRecorder) Namespaces(ctx, rollupId, limit, offset any) *MockIRollupNamespacesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Namespaces", reflect.TypeOf((*MockIRollup)(nil).Namespaces), ctx, rollupId, limit, offset)
	return &MockIRollupNamespacesCall{Call: call}
}

// MockIRollupNamespacesCall wrap *gomock.Call
type MockIRollupNamespacesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupNamespacesCall) Return(namespaceIds []uint64, err error) *MockIRollupNamespacesCall {
	c.Call = c.Call.Return(namespaceIds, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupNamespacesCall) Do(f func(context.Context, uint64, int, int) ([]uint64, error)) *MockIRollupNamespacesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupNamespacesCall) DoAndReturn(f func(context.Context, uint64, int, int) ([]uint64, error)) *MockIRollupNamespacesCall {
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
func (mr *MockIRollupMockRecorder) Providers(ctx, rollupId any) *MockIRollupProvidersCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Providers", reflect.TypeOf((*MockIRollup)(nil).Providers), ctx, rollupId)
	return &MockIRollupProvidersCall{Call: call}
}

// MockIRollupProvidersCall wrap *gomock.Call
type MockIRollupProvidersCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupProvidersCall) Return(providers []storage.RollupProvider, err error) *MockIRollupProvidersCall {
	c.Call = c.Call.Return(providers, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupProvidersCall) Do(f func(context.Context, uint64) ([]storage.RollupProvider, error)) *MockIRollupProvidersCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupProvidersCall) DoAndReturn(f func(context.Context, uint64) ([]storage.RollupProvider, error)) *MockIRollupProvidersCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// RollupsByNamespace mocks base method.
func (m *MockIRollup) RollupsByNamespace(ctx context.Context, namespaceId uint64, limit, offset int) ([]storage.Rollup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RollupsByNamespace", ctx, namespaceId, limit, offset)
	ret0, _ := ret[0].([]storage.Rollup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RollupsByNamespace indicates an expected call of RollupsByNamespace.
func (mr *MockIRollupMockRecorder) RollupsByNamespace(ctx, namespaceId, limit, offset any) *MockIRollupRollupsByNamespaceCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RollupsByNamespace", reflect.TypeOf((*MockIRollup)(nil).RollupsByNamespace), ctx, namespaceId, limit, offset)
	return &MockIRollupRollupsByNamespaceCall{Call: call}
}

// MockIRollupRollupsByNamespaceCall wrap *gomock.Call
type MockIRollupRollupsByNamespaceCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupRollupsByNamespaceCall) Return(rollups []storage.Rollup, err error) *MockIRollupRollupsByNamespaceCall {
	c.Call = c.Call.Return(rollups, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupRollupsByNamespaceCall) Do(f func(context.Context, uint64, int, int) ([]storage.Rollup, error)) *MockIRollupRollupsByNamespaceCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupRollupsByNamespaceCall) DoAndReturn(f func(context.Context, uint64, int, int) ([]storage.Rollup, error)) *MockIRollupRollupsByNamespaceCall {
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
func (mr *MockIRollupMockRecorder) Save(ctx, m any) *MockIRollupSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIRollup)(nil).Save), ctx, m)
	return &MockIRollupSaveCall{Call: call}
}

// MockIRollupSaveCall wrap *gomock.Call
type MockIRollupSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupSaveCall) Return(arg0 error) *MockIRollupSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupSaveCall) Do(f func(context.Context, *storage.Rollup) error) *MockIRollupSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupSaveCall) DoAndReturn(f func(context.Context, *storage.Rollup) error) *MockIRollupSaveCall {
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
func (mr *MockIRollupMockRecorder) Series(ctx, rollupId, timeframe, column, req any) *MockIRollupSeriesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Series", reflect.TypeOf((*MockIRollup)(nil).Series), ctx, rollupId, timeframe, column, req)
	return &MockIRollupSeriesCall{Call: call}
}

// MockIRollupSeriesCall wrap *gomock.Call
type MockIRollupSeriesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupSeriesCall) Return(items []storage.HistogramItem, err error) *MockIRollupSeriesCall {
	c.Call = c.Call.Return(items, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupSeriesCall) Do(f func(context.Context, uint64, string, string, storage.SeriesRequest) ([]storage.HistogramItem, error)) *MockIRollupSeriesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupSeriesCall) DoAndReturn(f func(context.Context, uint64, string, string, storage.SeriesRequest) ([]storage.HistogramItem, error)) *MockIRollupSeriesCall {
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
func (mr *MockIRollupMockRecorder) Update(ctx, m any) *MockIRollupUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIRollup)(nil).Update), ctx, m)
	return &MockIRollupUpdateCall{Call: call}
}

// MockIRollupUpdateCall wrap *gomock.Call
type MockIRollupUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupUpdateCall) Return(arg0 error) *MockIRollupUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupUpdateCall) Do(f func(context.Context, *storage.Rollup) error) *MockIRollupUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupUpdateCall) DoAndReturn(f func(context.Context, *storage.Rollup) error) *MockIRollupUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
