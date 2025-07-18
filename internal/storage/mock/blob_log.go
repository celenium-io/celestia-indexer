// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: blob_log.go
//
// Generated by this command:
//
//	mockgen -source=blob_log.go -destination=mock/blob_log.go -package=mock -typed
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	io "io"
	reflect "reflect"
	time "time"

	storage "github.com/celenium-io/celestia-indexer/internal/storage"
	types "github.com/celenium-io/celestia-indexer/pkg/types"
	storage0 "github.com/dipdup-net/indexer-sdk/pkg/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockIBlobLog is a mock of IBlobLog interface.
type MockIBlobLog struct {
	ctrl     *gomock.Controller
	recorder *MockIBlobLogMockRecorder
}

// MockIBlobLogMockRecorder is the mock recorder for MockIBlobLog.
type MockIBlobLogMockRecorder struct {
	mock *MockIBlobLog
}

// NewMockIBlobLog creates a new mock instance.
func NewMockIBlobLog(ctrl *gomock.Controller) *MockIBlobLog {
	mock := &MockIBlobLog{ctrl: ctrl}
	mock.recorder = &MockIBlobLogMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIBlobLog) EXPECT() *MockIBlobLogMockRecorder {
	return m.recorder
}

// Blob mocks base method.
func (m *MockIBlobLog) Blob(ctx context.Context, height types.Level, nsId uint64, commitment string) (storage.BlobLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Blob", ctx, height, nsId, commitment)
	ret0, _ := ret[0].(storage.BlobLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Blob indicates an expected call of Blob.
func (mr *MockIBlobLogMockRecorder) Blob(ctx, height, nsId, commitment any) *MockIBlobLogBlobCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Blob", reflect.TypeOf((*MockIBlobLog)(nil).Blob), ctx, height, nsId, commitment)
	return &MockIBlobLogBlobCall{Call: call}
}

// MockIBlobLogBlobCall wrap *gomock.Call
type MockIBlobLogBlobCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlobLogBlobCall) Return(arg0 storage.BlobLog, arg1 error) *MockIBlobLogBlobCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlobLogBlobCall) Do(f func(context.Context, types.Level, uint64, string) (storage.BlobLog, error)) *MockIBlobLogBlobCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlobLogBlobCall) DoAndReturn(f func(context.Context, types.Level, uint64, string) (storage.BlobLog, error)) *MockIBlobLogBlobCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ByHeight mocks base method.
func (m *MockIBlobLog) ByHeight(ctx context.Context, height types.Level, fltrs storage.BlobLogFilters) ([]storage.BlobLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByHeight", ctx, height, fltrs)
	ret0, _ := ret[0].([]storage.BlobLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByHeight indicates an expected call of ByHeight.
func (mr *MockIBlobLogMockRecorder) ByHeight(ctx, height, fltrs any) *MockIBlobLogByHeightCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByHeight", reflect.TypeOf((*MockIBlobLog)(nil).ByHeight), ctx, height, fltrs)
	return &MockIBlobLogByHeightCall{Call: call}
}

// MockIBlobLogByHeightCall wrap *gomock.Call
type MockIBlobLogByHeightCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlobLogByHeightCall) Return(arg0 []storage.BlobLog, arg1 error) *MockIBlobLogByHeightCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlobLogByHeightCall) Do(f func(context.Context, types.Level, storage.BlobLogFilters) ([]storage.BlobLog, error)) *MockIBlobLogByHeightCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlobLogByHeightCall) DoAndReturn(f func(context.Context, types.Level, storage.BlobLogFilters) ([]storage.BlobLog, error)) *MockIBlobLogByHeightCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ByNamespace mocks base method.
func (m *MockIBlobLog) ByNamespace(ctx context.Context, nsId uint64, fltrs storage.BlobLogFilters) ([]storage.BlobLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByNamespace", ctx, nsId, fltrs)
	ret0, _ := ret[0].([]storage.BlobLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByNamespace indicates an expected call of ByNamespace.
func (mr *MockIBlobLogMockRecorder) ByNamespace(ctx, nsId, fltrs any) *MockIBlobLogByNamespaceCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByNamespace", reflect.TypeOf((*MockIBlobLog)(nil).ByNamespace), ctx, nsId, fltrs)
	return &MockIBlobLogByNamespaceCall{Call: call}
}

// MockIBlobLogByNamespaceCall wrap *gomock.Call
type MockIBlobLogByNamespaceCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlobLogByNamespaceCall) Return(arg0 []storage.BlobLog, arg1 error) *MockIBlobLogByNamespaceCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlobLogByNamespaceCall) Do(f func(context.Context, uint64, storage.BlobLogFilters) ([]storage.BlobLog, error)) *MockIBlobLogByNamespaceCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlobLogByNamespaceCall) DoAndReturn(f func(context.Context, uint64, storage.BlobLogFilters) ([]storage.BlobLog, error)) *MockIBlobLogByNamespaceCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ByProviders mocks base method.
func (m *MockIBlobLog) ByProviders(ctx context.Context, providers []storage.RollupProvider, fltrs storage.BlobLogFilters) ([]storage.BlobLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByProviders", ctx, providers, fltrs)
	ret0, _ := ret[0].([]storage.BlobLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByProviders indicates an expected call of ByProviders.
func (mr *MockIBlobLogMockRecorder) ByProviders(ctx, providers, fltrs any) *MockIBlobLogByProvidersCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByProviders", reflect.TypeOf((*MockIBlobLog)(nil).ByProviders), ctx, providers, fltrs)
	return &MockIBlobLogByProvidersCall{Call: call}
}

// MockIBlobLogByProvidersCall wrap *gomock.Call
type MockIBlobLogByProvidersCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlobLogByProvidersCall) Return(arg0 []storage.BlobLog, arg1 error) *MockIBlobLogByProvidersCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlobLogByProvidersCall) Do(f func(context.Context, []storage.RollupProvider, storage.BlobLogFilters) ([]storage.BlobLog, error)) *MockIBlobLogByProvidersCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlobLogByProvidersCall) DoAndReturn(f func(context.Context, []storage.RollupProvider, storage.BlobLogFilters) ([]storage.BlobLog, error)) *MockIBlobLogByProvidersCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// BySigner mocks base method.
func (m *MockIBlobLog) BySigner(ctx context.Context, signerId uint64, fltrs storage.BlobLogFilters) ([]storage.BlobLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BySigner", ctx, signerId, fltrs)
	ret0, _ := ret[0].([]storage.BlobLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BySigner indicates an expected call of BySigner.
func (mr *MockIBlobLogMockRecorder) BySigner(ctx, signerId, fltrs any) *MockIBlobLogBySignerCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BySigner", reflect.TypeOf((*MockIBlobLog)(nil).BySigner), ctx, signerId, fltrs)
	return &MockIBlobLogBySignerCall{Call: call}
}

// MockIBlobLogBySignerCall wrap *gomock.Call
type MockIBlobLogBySignerCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlobLogBySignerCall) Return(arg0 []storage.BlobLog, arg1 error) *MockIBlobLogBySignerCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlobLogBySignerCall) Do(f func(context.Context, uint64, storage.BlobLogFilters) ([]storage.BlobLog, error)) *MockIBlobLogBySignerCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlobLogBySignerCall) DoAndReturn(f func(context.Context, uint64, storage.BlobLogFilters) ([]storage.BlobLog, error)) *MockIBlobLogBySignerCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ByTxId mocks base method.
func (m *MockIBlobLog) ByTxId(ctx context.Context, txId uint64, fltrs storage.BlobLogFilters) ([]storage.BlobLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByTxId", ctx, txId, fltrs)
	ret0, _ := ret[0].([]storage.BlobLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByTxId indicates an expected call of ByTxId.
func (mr *MockIBlobLogMockRecorder) ByTxId(ctx, txId, fltrs any) *MockIBlobLogByTxIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByTxId", reflect.TypeOf((*MockIBlobLog)(nil).ByTxId), ctx, txId, fltrs)
	return &MockIBlobLogByTxIdCall{Call: call}
}

// MockIBlobLogByTxIdCall wrap *gomock.Call
type MockIBlobLogByTxIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlobLogByTxIdCall) Return(arg0 []storage.BlobLog, arg1 error) *MockIBlobLogByTxIdCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlobLogByTxIdCall) Do(f func(context.Context, uint64, storage.BlobLogFilters) ([]storage.BlobLog, error)) *MockIBlobLogByTxIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlobLogByTxIdCall) DoAndReturn(f func(context.Context, uint64, storage.BlobLogFilters) ([]storage.BlobLog, error)) *MockIBlobLogByTxIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CountByTxId mocks base method.
func (m *MockIBlobLog) CountByTxId(ctx context.Context, txId uint64) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountByTxId", ctx, txId)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountByTxId indicates an expected call of CountByTxId.
func (mr *MockIBlobLogMockRecorder) CountByTxId(ctx, txId any) *MockIBlobLogCountByTxIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByTxId", reflect.TypeOf((*MockIBlobLog)(nil).CountByTxId), ctx, txId)
	return &MockIBlobLogCountByTxIdCall{Call: call}
}

// MockIBlobLogCountByTxIdCall wrap *gomock.Call
type MockIBlobLogCountByTxIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlobLogCountByTxIdCall) Return(arg0 int, arg1 error) *MockIBlobLogCountByTxIdCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlobLogCountByTxIdCall) Do(f func(context.Context, uint64) (int, error)) *MockIBlobLogCountByTxIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlobLogCountByTxIdCall) DoAndReturn(f func(context.Context, uint64) (int, error)) *MockIBlobLogCountByTxIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CursorList mocks base method.
func (m *MockIBlobLog) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.BlobLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.BlobLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIBlobLogMockRecorder) CursorList(ctx, id, limit, order, cmp any) *MockIBlobLogCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIBlobLog)(nil).CursorList), ctx, id, limit, order, cmp)
	return &MockIBlobLogCursorListCall{Call: call}
}

// MockIBlobLogCursorListCall wrap *gomock.Call
type MockIBlobLogCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlobLogCursorListCall) Return(arg0 []*storage.BlobLog, arg1 error) *MockIBlobLogCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlobLogCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.BlobLog, error)) *MockIBlobLogCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlobLogCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.BlobLog, error)) *MockIBlobLogCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ExportByProviders mocks base method.
func (m *MockIBlobLog) ExportByProviders(ctx context.Context, providers []storage.RollupProvider, from, to time.Time, stream io.Writer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExportByProviders", ctx, providers, from, to, stream)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExportByProviders indicates an expected call of ExportByProviders.
func (mr *MockIBlobLogMockRecorder) ExportByProviders(ctx, providers, from, to, stream any) *MockIBlobLogExportByProvidersCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExportByProviders", reflect.TypeOf((*MockIBlobLog)(nil).ExportByProviders), ctx, providers, from, to, stream)
	return &MockIBlobLogExportByProvidersCall{Call: call}
}

// MockIBlobLogExportByProvidersCall wrap *gomock.Call
type MockIBlobLogExportByProvidersCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlobLogExportByProvidersCall) Return(err error) *MockIBlobLogExportByProvidersCall {
	c.Call = c.Call.Return(err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlobLogExportByProvidersCall) Do(f func(context.Context, []storage.RollupProvider, time.Time, time.Time, io.Writer) error) *MockIBlobLogExportByProvidersCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlobLogExportByProvidersCall) DoAndReturn(f func(context.Context, []storage.RollupProvider, time.Time, time.Time, io.Writer) error) *MockIBlobLogExportByProvidersCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIBlobLog) GetByID(ctx context.Context, id uint64) (*storage.BlobLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.BlobLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIBlobLogMockRecorder) GetByID(ctx, id any) *MockIBlobLogGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIBlobLog)(nil).GetByID), ctx, id)
	return &MockIBlobLogGetByIDCall{Call: call}
}

// MockIBlobLogGetByIDCall wrap *gomock.Call
type MockIBlobLogGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlobLogGetByIDCall) Return(arg0 *storage.BlobLog, arg1 error) *MockIBlobLogGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlobLogGetByIDCall) Do(f func(context.Context, uint64) (*storage.BlobLog, error)) *MockIBlobLogGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlobLogGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.BlobLog, error)) *MockIBlobLogGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIBlobLog) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIBlobLogMockRecorder) IsNoRows(err any) *MockIBlobLogIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIBlobLog)(nil).IsNoRows), err)
	return &MockIBlobLogIsNoRowsCall{Call: call}
}

// MockIBlobLogIsNoRowsCall wrap *gomock.Call
type MockIBlobLogIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlobLogIsNoRowsCall) Return(arg0 bool) *MockIBlobLogIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlobLogIsNoRowsCall) Do(f func(error) bool) *MockIBlobLogIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlobLogIsNoRowsCall) DoAndReturn(f func(error) bool) *MockIBlobLogIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIBlobLog) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIBlobLogMockRecorder) LastID(ctx any) *MockIBlobLogLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIBlobLog)(nil).LastID), ctx)
	return &MockIBlobLogLastIDCall{Call: call}
}

// MockIBlobLogLastIDCall wrap *gomock.Call
type MockIBlobLogLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlobLogLastIDCall) Return(arg0 uint64, arg1 error) *MockIBlobLogLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlobLogLastIDCall) Do(f func(context.Context) (uint64, error)) *MockIBlobLogLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlobLogLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *MockIBlobLogLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIBlobLog) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.BlobLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.BlobLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIBlobLogMockRecorder) List(ctx, limit, offset, order any) *MockIBlobLogListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIBlobLog)(nil).List), ctx, limit, offset, order)
	return &MockIBlobLogListCall{Call: call}
}

// MockIBlobLogListCall wrap *gomock.Call
type MockIBlobLogListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlobLogListCall) Return(arg0 []*storage.BlobLog, arg1 error) *MockIBlobLogListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlobLogListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.BlobLog, error)) *MockIBlobLogListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlobLogListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.BlobLog, error)) *MockIBlobLogListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ListBlobs mocks base method.
func (m *MockIBlobLog) ListBlobs(ctx context.Context, fltrs storage.ListBlobLogFilters) ([]storage.BlobLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListBlobs", ctx, fltrs)
	ret0, _ := ret[0].([]storage.BlobLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListBlobs indicates an expected call of ListBlobs.
func (mr *MockIBlobLogMockRecorder) ListBlobs(ctx, fltrs any) *MockIBlobLogListBlobsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListBlobs", reflect.TypeOf((*MockIBlobLog)(nil).ListBlobs), ctx, fltrs)
	return &MockIBlobLogListBlobsCall{Call: call}
}

// MockIBlobLogListBlobsCall wrap *gomock.Call
type MockIBlobLogListBlobsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlobLogListBlobsCall) Return(arg0 []storage.BlobLog, arg1 error) *MockIBlobLogListBlobsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlobLogListBlobsCall) Do(f func(context.Context, storage.ListBlobLogFilters) ([]storage.BlobLog, error)) *MockIBlobLogListBlobsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlobLogListBlobsCall) DoAndReturn(f func(context.Context, storage.ListBlobLogFilters) ([]storage.BlobLog, error)) *MockIBlobLogListBlobsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIBlobLog) Save(ctx context.Context, m *storage.BlobLog) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIBlobLogMockRecorder) Save(ctx, m any) *MockIBlobLogSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIBlobLog)(nil).Save), ctx, m)
	return &MockIBlobLogSaveCall{Call: call}
}

// MockIBlobLogSaveCall wrap *gomock.Call
type MockIBlobLogSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlobLogSaveCall) Return(arg0 error) *MockIBlobLogSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlobLogSaveCall) Do(f func(context.Context, *storage.BlobLog) error) *MockIBlobLogSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlobLogSaveCall) DoAndReturn(f func(context.Context, *storage.BlobLog) error) *MockIBlobLogSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIBlobLog) Update(ctx context.Context, m *storage.BlobLog) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIBlobLogMockRecorder) Update(ctx, m any) *MockIBlobLogUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIBlobLog)(nil).Update), ctx, m)
	return &MockIBlobLogUpdateCall{Call: call}
}

// MockIBlobLogUpdateCall wrap *gomock.Call
type MockIBlobLogUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlobLogUpdateCall) Return(arg0 error) *MockIBlobLogUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlobLogUpdateCall) Do(f func(context.Context, *storage.BlobLog) error) *MockIBlobLogUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlobLogUpdateCall) DoAndReturn(f func(context.Context, *storage.BlobLog) error) *MockIBlobLogUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
