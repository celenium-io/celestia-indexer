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

// Count mocks base method.
func (m *MockIStats) Count(ctx context.Context, req storage.CountRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, req)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockIStatsMockRecorder) Count(ctx, req any) *IStatsCountCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockIStats)(nil).Count), ctx, req)
	return &IStatsCountCall{Call: call}
}

// IStatsCountCall wrap *gomock.Call
type IStatsCountCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStatsCountCall) Return(arg0 string, arg1 error) *IStatsCountCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStatsCountCall) Do(f func(context.Context, storage.CountRequest) (string, error)) *IStatsCountCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStatsCountCall) DoAndReturn(f func(context.Context, storage.CountRequest) (string, error)) *IStatsCountCall {
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
func (mr *MockIStatsMockRecorder) NamespaceSeries(ctx, timeframe, name, nsId, req any) *IStatsNamespaceSeriesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NamespaceSeries", reflect.TypeOf((*MockIStats)(nil).NamespaceSeries), ctx, timeframe, name, nsId, req)
	return &IStatsNamespaceSeriesCall{Call: call}
}

// IStatsNamespaceSeriesCall wrap *gomock.Call
type IStatsNamespaceSeriesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStatsNamespaceSeriesCall) Return(response []storage.SeriesItem, err error) *IStatsNamespaceSeriesCall {
	c.Call = c.Call.Return(response, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStatsNamespaceSeriesCall) Do(f func(context.Context, storage.Timeframe, string, uint64, storage.SeriesRequest) ([]storage.SeriesItem, error)) *IStatsNamespaceSeriesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStatsNamespaceSeriesCall) DoAndReturn(f func(context.Context, storage.Timeframe, string, uint64, storage.SeriesRequest) ([]storage.SeriesItem, error)) *IStatsNamespaceSeriesCall {
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
func (mr *MockIStatsMockRecorder) Series(ctx, timeframe, name, req any) *IStatsSeriesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Series", reflect.TypeOf((*MockIStats)(nil).Series), ctx, timeframe, name, req)
	return &IStatsSeriesCall{Call: call}
}

// IStatsSeriesCall wrap *gomock.Call
type IStatsSeriesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStatsSeriesCall) Return(arg0 []storage.SeriesItem, arg1 error) *IStatsSeriesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStatsSeriesCall) Do(f func(context.Context, storage.Timeframe, string, storage.SeriesRequest) ([]storage.SeriesItem, error)) *IStatsSeriesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStatsSeriesCall) DoAndReturn(f func(context.Context, storage.Timeframe, string, storage.SeriesRequest) ([]storage.SeriesItem, error)) *IStatsSeriesCall {
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
func (mr *MockIStatsMockRecorder) StakingSeries(ctx, timeframe, name, validatorId, req any) *IStatsStakingSeriesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StakingSeries", reflect.TypeOf((*MockIStats)(nil).StakingSeries), ctx, timeframe, name, validatorId, req)
	return &IStatsStakingSeriesCall{Call: call}
}

// IStatsStakingSeriesCall wrap *gomock.Call
type IStatsStakingSeriesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStatsStakingSeriesCall) Return(response []storage.SeriesItem, err error) *IStatsStakingSeriesCall {
	c.Call = c.Call.Return(response, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStatsStakingSeriesCall) Do(f func(context.Context, storage.Timeframe, string, uint64, storage.SeriesRequest) ([]storage.SeriesItem, error)) *IStatsStakingSeriesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStatsStakingSeriesCall) DoAndReturn(f func(context.Context, storage.Timeframe, string, uint64, storage.SeriesRequest) ([]storage.SeriesItem, error)) *IStatsStakingSeriesCall {
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
func (mr *MockIStatsMockRecorder) Summary(ctx, req any) *IStatsSummaryCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Summary", reflect.TypeOf((*MockIStats)(nil).Summary), ctx, req)
	return &IStatsSummaryCall{Call: call}
}

// IStatsSummaryCall wrap *gomock.Call
type IStatsSummaryCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStatsSummaryCall) Return(arg0 string, arg1 error) *IStatsSummaryCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStatsSummaryCall) Do(f func(context.Context, storage.SummaryRequest) (string, error)) *IStatsSummaryCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStatsSummaryCall) DoAndReturn(f func(context.Context, storage.SummaryRequest) (string, error)) *IStatsSummaryCall {
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
func (mr *MockIStatsMockRecorder) TPS(ctx any) *IStatsTPSCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TPS", reflect.TypeOf((*MockIStats)(nil).TPS), ctx)
	return &IStatsTPSCall{Call: call}
}

// IStatsTPSCall wrap *gomock.Call
type IStatsTPSCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStatsTPSCall) Return(arg0 storage.TPS, arg1 error) *IStatsTPSCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStatsTPSCall) Do(f func(context.Context) (storage.TPS, error)) *IStatsTPSCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStatsTPSCall) DoAndReturn(f func(context.Context) (storage.TPS, error)) *IStatsTPSCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// TxCountForLast24h mocks base method.
func (m *MockIStats) TxCountForLast24h(ctx context.Context) ([]storage.TxCountForLast24hItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TxCountForLast24h", ctx)
	ret0, _ := ret[0].([]storage.TxCountForLast24hItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TxCountForLast24h indicates an expected call of TxCountForLast24h.
func (mr *MockIStatsMockRecorder) TxCountForLast24h(ctx any) *IStatsTxCountForLast24hCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxCountForLast24h", reflect.TypeOf((*MockIStats)(nil).TxCountForLast24h), ctx)
	return &IStatsTxCountForLast24hCall{Call: call}
}

// IStatsTxCountForLast24hCall wrap *gomock.Call
type IStatsTxCountForLast24hCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStatsTxCountForLast24hCall) Return(arg0 []storage.TxCountForLast24hItem, arg1 error) *IStatsTxCountForLast24hCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStatsTxCountForLast24hCall) Do(f func(context.Context) ([]storage.TxCountForLast24hItem, error)) *IStatsTxCountForLast24hCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStatsTxCountForLast24hCall) DoAndReturn(f func(context.Context) ([]storage.TxCountForLast24hItem, error)) *IStatsTxCountForLast24hCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
