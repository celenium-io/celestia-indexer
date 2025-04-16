// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: stats.go
//
// Generated by this command:
//
//	mockgen -source=stats.go -destination=mock/stats.go -package=mock -typed
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	time "time"

	storage "github.com/celenium-io/celestia-indexer/internal/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockIStats is a mock of IStats interface.
type MockIStats struct {
	ctrl     *gomock.Controller
	recorder *MockIStatsMockRecorder
}

// MockIStatsMockRecorder is the mock recorder for MockIStats.
type MockIStatsMockRecorder struct {
	mock *MockIStats
}

// NewMockIStats creates a new mock instance.
func NewMockIStats(ctrl *gomock.Controller) *MockIStats {
	mock := &MockIStats{ctrl: ctrl}
	mock.recorder = &MockIStatsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIStats) EXPECT() *MockIStatsMockRecorder {
	return m.recorder
}

// Change24hBlockStats mocks base method.
func (m *MockIStats) Change24hBlockStats(ctx context.Context) (storage.Change24hBlockStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Change24hBlockStats", ctx)
	ret0, _ := ret[0].(storage.Change24hBlockStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Change24hBlockStats indicates an expected call of Change24hBlockStats.
func (mr *MockIStatsMockRecorder) Change24hBlockStats(ctx any) *MockIStatsChange24hBlockStatsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Change24hBlockStats", reflect.TypeOf((*MockIStats)(nil).Change24hBlockStats), ctx)
	return &MockIStatsChange24hBlockStatsCall{Call: call}
}

// MockIStatsChange24hBlockStatsCall wrap *gomock.Call
type MockIStatsChange24hBlockStatsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIStatsChange24hBlockStatsCall) Return(response storage.Change24hBlockStats, err error) *MockIStatsChange24hBlockStatsCall {
	c.Call = c.Call.Return(response, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIStatsChange24hBlockStatsCall) Do(f func(context.Context) (storage.Change24hBlockStats, error)) *MockIStatsChange24hBlockStatsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIStatsChange24hBlockStatsCall) DoAndReturn(f func(context.Context) (storage.Change24hBlockStats, error)) *MockIStatsChange24hBlockStatsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Count mocks base method.
func (m *MockIStats) Count(ctx context.Context, req storage.CountRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, req)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockIStatsMockRecorder) Count(ctx, req any) *MockIStatsCountCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockIStats)(nil).Count), ctx, req)
	return &MockIStatsCountCall{Call: call}
}

// MockIStatsCountCall wrap *gomock.Call
type MockIStatsCountCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIStatsCountCall) Return(arg0 string, arg1 error) *MockIStatsCountCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIStatsCountCall) Do(f func(context.Context, storage.CountRequest) (string, error)) *MockIStatsCountCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIStatsCountCall) DoAndReturn(f func(context.Context, storage.CountRequest) (string, error)) *MockIStatsCountCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CumulativeSeries mocks base method.
func (m *MockIStats) CumulativeSeries(ctx context.Context, timeframe storage.Timeframe, name string, req storage.SeriesRequest) ([]storage.SeriesItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CumulativeSeries", ctx, timeframe, name, req)
	ret0, _ := ret[0].([]storage.SeriesItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CumulativeSeries indicates an expected call of CumulativeSeries.
func (mr *MockIStatsMockRecorder) CumulativeSeries(ctx, timeframe, name, req any) *MockIStatsCumulativeSeriesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CumulativeSeries", reflect.TypeOf((*MockIStats)(nil).CumulativeSeries), ctx, timeframe, name, req)
	return &MockIStatsCumulativeSeriesCall{Call: call}
}

// MockIStatsCumulativeSeriesCall wrap *gomock.Call
type MockIStatsCumulativeSeriesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIStatsCumulativeSeriesCall) Return(arg0 []storage.SeriesItem, arg1 error) *MockIStatsCumulativeSeriesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIStatsCumulativeSeriesCall) Do(f func(context.Context, storage.Timeframe, string, storage.SeriesRequest) ([]storage.SeriesItem, error)) *MockIStatsCumulativeSeriesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIStatsCumulativeSeriesCall) DoAndReturn(f func(context.Context, storage.Timeframe, string, storage.SeriesRequest) ([]storage.SeriesItem, error)) *MockIStatsCumulativeSeriesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MessagesCount24h mocks base method.
func (m *MockIStats) MessagesCount24h(ctx context.Context) ([]storage.CountItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MessagesCount24h", ctx)
	ret0, _ := ret[0].([]storage.CountItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MessagesCount24h indicates an expected call of MessagesCount24h.
func (mr *MockIStatsMockRecorder) MessagesCount24h(ctx any) *MockIStatsMessagesCount24hCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MessagesCount24h", reflect.TypeOf((*MockIStats)(nil).MessagesCount24h), ctx)
	return &MockIStatsMessagesCount24hCall{Call: call}
}

// MockIStatsMessagesCount24hCall wrap *gomock.Call
type MockIStatsMessagesCount24hCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIStatsMessagesCount24hCall) Return(arg0 []storage.CountItem, arg1 error) *MockIStatsMessagesCount24hCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIStatsMessagesCount24hCall) Do(f func(context.Context) ([]storage.CountItem, error)) *MockIStatsMessagesCount24hCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIStatsMessagesCount24hCall) DoAndReturn(f func(context.Context) ([]storage.CountItem, error)) *MockIStatsMessagesCount24hCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// NamespaceSeries mocks base method.
func (m *MockIStats) NamespaceSeries(ctx context.Context, timeframe storage.Timeframe, name string, nsId uint64, req storage.SeriesRequest) ([]storage.SeriesItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NamespaceSeries", ctx, timeframe, name, nsId, req)
	ret0, _ := ret[0].([]storage.SeriesItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NamespaceSeries indicates an expected call of NamespaceSeries.
func (mr *MockIStatsMockRecorder) NamespaceSeries(ctx, timeframe, name, nsId, req any) *MockIStatsNamespaceSeriesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NamespaceSeries", reflect.TypeOf((*MockIStats)(nil).NamespaceSeries), ctx, timeframe, name, nsId, req)
	return &MockIStatsNamespaceSeriesCall{Call: call}
}

// MockIStatsNamespaceSeriesCall wrap *gomock.Call
type MockIStatsNamespaceSeriesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIStatsNamespaceSeriesCall) Return(response []storage.SeriesItem, err error) *MockIStatsNamespaceSeriesCall {
	c.Call = c.Call.Return(response, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIStatsNamespaceSeriesCall) Do(f func(context.Context, storage.Timeframe, string, uint64, storage.SeriesRequest) ([]storage.SeriesItem, error)) *MockIStatsNamespaceSeriesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIStatsNamespaceSeriesCall) DoAndReturn(f func(context.Context, storage.Timeframe, string, uint64, storage.SeriesRequest) ([]storage.SeriesItem, error)) *MockIStatsNamespaceSeriesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// RollupStats24h mocks base method.
func (m *MockIStats) RollupStats24h(ctx context.Context) ([]storage.RollupStats24h, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RollupStats24h", ctx)
	ret0, _ := ret[0].([]storage.RollupStats24h)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RollupStats24h indicates an expected call of RollupStats24h.
func (mr *MockIStatsMockRecorder) RollupStats24h(ctx any) *MockIStatsRollupStats24hCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RollupStats24h", reflect.TypeOf((*MockIStats)(nil).RollupStats24h), ctx)
	return &MockIStatsRollupStats24hCall{Call: call}
}

// MockIStatsRollupStats24hCall wrap *gomock.Call
type MockIStatsRollupStats24hCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIStatsRollupStats24hCall) Return(arg0 []storage.RollupStats24h, arg1 error) *MockIStatsRollupStats24hCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIStatsRollupStats24hCall) Do(f func(context.Context) ([]storage.RollupStats24h, error)) *MockIStatsRollupStats24hCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIStatsRollupStats24hCall) DoAndReturn(f func(context.Context) ([]storage.RollupStats24h, error)) *MockIStatsRollupStats24hCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Series mocks base method.
func (m *MockIStats) Series(ctx context.Context, timeframe storage.Timeframe, name string, req storage.SeriesRequest) ([]storage.SeriesItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Series", ctx, timeframe, name, req)
	ret0, _ := ret[0].([]storage.SeriesItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Series indicates an expected call of Series.
func (mr *MockIStatsMockRecorder) Series(ctx, timeframe, name, req any) *MockIStatsSeriesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Series", reflect.TypeOf((*MockIStats)(nil).Series), ctx, timeframe, name, req)
	return &MockIStatsSeriesCall{Call: call}
}

// MockIStatsSeriesCall wrap *gomock.Call
type MockIStatsSeriesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIStatsSeriesCall) Return(arg0 []storage.SeriesItem, arg1 error) *MockIStatsSeriesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIStatsSeriesCall) Do(f func(context.Context, storage.Timeframe, string, storage.SeriesRequest) ([]storage.SeriesItem, error)) *MockIStatsSeriesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIStatsSeriesCall) DoAndReturn(f func(context.Context, storage.Timeframe, string, storage.SeriesRequest) ([]storage.SeriesItem, error)) *MockIStatsSeriesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SizeGroups mocks base method.
func (m *MockIStats) SizeGroups(ctx context.Context, timeFilter *time.Time) ([]storage.SizeGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SizeGroups", ctx, timeFilter)
	ret0, _ := ret[0].([]storage.SizeGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SizeGroups indicates an expected call of SizeGroups.
func (mr *MockIStatsMockRecorder) SizeGroups(ctx, timeFilter any) *MockIStatsSizeGroupsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SizeGroups", reflect.TypeOf((*MockIStats)(nil).SizeGroups), ctx, timeFilter)
	return &MockIStatsSizeGroupsCall{Call: call}
}

// MockIStatsSizeGroupsCall wrap *gomock.Call
type MockIStatsSizeGroupsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIStatsSizeGroupsCall) Return(arg0 []storage.SizeGroup, arg1 error) *MockIStatsSizeGroupsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIStatsSizeGroupsCall) Do(f func(context.Context, *time.Time) ([]storage.SizeGroup, error)) *MockIStatsSizeGroupsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIStatsSizeGroupsCall) DoAndReturn(f func(context.Context, *time.Time) ([]storage.SizeGroup, error)) *MockIStatsSizeGroupsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SquareSize mocks base method.
func (m *MockIStats) SquareSize(ctx context.Context, from, to *time.Time) (map[int][]storage.SeriesItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SquareSize", ctx, from, to)
	ret0, _ := ret[0].(map[int][]storage.SeriesItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SquareSize indicates an expected call of SquareSize.
func (mr *MockIStatsMockRecorder) SquareSize(ctx, from, to any) *MockIStatsSquareSizeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SquareSize", reflect.TypeOf((*MockIStats)(nil).SquareSize), ctx, from, to)
	return &MockIStatsSquareSizeCall{Call: call}
}

// MockIStatsSquareSizeCall wrap *gomock.Call
type MockIStatsSquareSizeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIStatsSquareSizeCall) Return(arg0 map[int][]storage.SeriesItem, arg1 error) *MockIStatsSquareSizeCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIStatsSquareSizeCall) Do(f func(context.Context, *time.Time, *time.Time) (map[int][]storage.SeriesItem, error)) *MockIStatsSquareSizeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIStatsSquareSizeCall) DoAndReturn(f func(context.Context, *time.Time, *time.Time) (map[int][]storage.SeriesItem, error)) *MockIStatsSquareSizeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// StakingSeries mocks base method.
func (m *MockIStats) StakingSeries(ctx context.Context, timeframe storage.Timeframe, name string, validatorId uint64, req storage.SeriesRequest) ([]storage.SeriesItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StakingSeries", ctx, timeframe, name, validatorId, req)
	ret0, _ := ret[0].([]storage.SeriesItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StakingSeries indicates an expected call of StakingSeries.
func (mr *MockIStatsMockRecorder) StakingSeries(ctx, timeframe, name, validatorId, req any) *MockIStatsStakingSeriesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StakingSeries", reflect.TypeOf((*MockIStats)(nil).StakingSeries), ctx, timeframe, name, validatorId, req)
	return &MockIStatsStakingSeriesCall{Call: call}
}

// MockIStatsStakingSeriesCall wrap *gomock.Call
type MockIStatsStakingSeriesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIStatsStakingSeriesCall) Return(response []storage.SeriesItem, err error) *MockIStatsStakingSeriesCall {
	c.Call = c.Call.Return(response, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIStatsStakingSeriesCall) Do(f func(context.Context, storage.Timeframe, string, uint64, storage.SeriesRequest) ([]storage.SeriesItem, error)) *MockIStatsStakingSeriesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIStatsStakingSeriesCall) DoAndReturn(f func(context.Context, storage.Timeframe, string, uint64, storage.SeriesRequest) ([]storage.SeriesItem, error)) *MockIStatsStakingSeriesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Summary mocks base method.
func (m *MockIStats) Summary(ctx context.Context, req storage.SummaryRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Summary", ctx, req)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Summary indicates an expected call of Summary.
func (mr *MockIStatsMockRecorder) Summary(ctx, req any) *MockIStatsSummaryCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Summary", reflect.TypeOf((*MockIStats)(nil).Summary), ctx, req)
	return &MockIStatsSummaryCall{Call: call}
}

// MockIStatsSummaryCall wrap *gomock.Call
type MockIStatsSummaryCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIStatsSummaryCall) Return(arg0 string, arg1 error) *MockIStatsSummaryCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIStatsSummaryCall) Do(f func(context.Context, storage.SummaryRequest) (string, error)) *MockIStatsSummaryCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIStatsSummaryCall) DoAndReturn(f func(context.Context, storage.SummaryRequest) (string, error)) *MockIStatsSummaryCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// TPS mocks base method.
func (m *MockIStats) TPS(ctx context.Context) (storage.TPS, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TPS", ctx)
	ret0, _ := ret[0].(storage.TPS)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TPS indicates an expected call of TPS.
func (mr *MockIStatsMockRecorder) TPS(ctx any) *MockIStatsTPSCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TPS", reflect.TypeOf((*MockIStats)(nil).TPS), ctx)
	return &MockIStatsTPSCall{Call: call}
}

// MockIStatsTPSCall wrap *gomock.Call
type MockIStatsTPSCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIStatsTPSCall) Return(arg0 storage.TPS, arg1 error) *MockIStatsTPSCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIStatsTPSCall) Do(f func(context.Context) (storage.TPS, error)) *MockIStatsTPSCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIStatsTPSCall) DoAndReturn(f func(context.Context) (storage.TPS, error)) *MockIStatsTPSCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
