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

// Histogram mocks base method.
func (m *MockIStats) Histogram(ctx context.Context, req storage.HistogramRequest) ([]storage.HistogramItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Histogram", ctx, req)
	ret0, _ := ret[0].([]storage.HistogramItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Histogram indicates an expected call of Histogram.
func (mr *MockIStatsMockRecorder) Histogram(ctx, req any) *IStatsHistogramCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Histogram", reflect.TypeOf((*MockIStats)(nil).Histogram), ctx, req)
	return &IStatsHistogramCall{Call: call}
}

// IStatsHistogramCall wrap *gomock.Call
type IStatsHistogramCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStatsHistogramCall) Return(arg0 []storage.HistogramItem, arg1 error) *IStatsHistogramCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStatsHistogramCall) Do(f func(context.Context, storage.HistogramRequest) ([]storage.HistogramItem, error)) *IStatsHistogramCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStatsHistogramCall) DoAndReturn(f func(context.Context, storage.HistogramRequest) ([]storage.HistogramItem, error)) *IStatsHistogramCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// HistogramCount mocks base method.
func (m *MockIStats) HistogramCount(ctx context.Context, req storage.HistogramCountRequest) ([]storage.HistogramItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HistogramCount", ctx, req)
	ret0, _ := ret[0].([]storage.HistogramItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HistogramCount indicates an expected call of HistogramCount.
func (mr *MockIStatsMockRecorder) HistogramCount(ctx, req any) *IStatsHistogramCountCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HistogramCount", reflect.TypeOf((*MockIStats)(nil).HistogramCount), ctx, req)
	return &IStatsHistogramCountCall{Call: call}
}

// IStatsHistogramCountCall wrap *gomock.Call
type IStatsHistogramCountCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IStatsHistogramCountCall) Return(arg0 []storage.HistogramItem, arg1 error) *IStatsHistogramCountCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IStatsHistogramCountCall) Do(f func(context.Context, storage.HistogramCountRequest) ([]storage.HistogramItem, error)) *IStatsHistogramCountCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IStatsHistogramCountCall) DoAndReturn(f func(context.Context, storage.HistogramCountRequest) ([]storage.HistogramItem, error)) *IStatsHistogramCountCall {
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
