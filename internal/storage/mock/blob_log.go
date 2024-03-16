// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
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

// ByHeight mocks base method.
func (m *MockIBlobLog) ByHeight(ctx context.Context, height types.Level, fltrs storage.BlobLogFilters) ([]storage.BlobLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByHeight", ctx, height, fltrs)
	ret0, _ := ret[0].([]storage.BlobLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByHeight indicates an expected call of ByHeight.
func (mr *MockIBlobLogMockRecorder) ByHeight(ctx, height, fltrs any) *IBlobLogByHeightCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByHeight", reflect.TypeOf((*MockIBlobLog)(nil).ByHeight), ctx, height, fltrs)
	return &IBlobLogByHeightCall{Call: call}
}

// IBlobLogByHeightCall wrap *gomock.Call
type IBlobLogByHeightCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBlobLogByHeightCall) Return(arg0 []storage.BlobLog, arg1 error) *IBlobLogByHeightCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBlobLogByHeightCall) Do(f func(context.Context, types.Level, storage.BlobLogFilters) ([]storage.BlobLog, error)) *IBlobLogByHeightCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBlobLogByHeightCall) DoAndReturn(f func(context.Context, types.Level, storage.BlobLogFilters) ([]storage.BlobLog, error)) *IBlobLogByHeightCall {
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
func (mr *MockIBlobLogMockRecorder) ByNamespace(ctx, nsId, fltrs any) *IBlobLogByNamespaceCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByNamespace", reflect.TypeOf((*MockIBlobLog)(nil).ByNamespace), ctx, nsId, fltrs)
	return &IBlobLogByNamespaceCall{Call: call}
}

// IBlobLogByNamespaceCall wrap *gomock.Call
type IBlobLogByNamespaceCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBlobLogByNamespaceCall) Return(arg0 []storage.BlobLog, arg1 error) *IBlobLogByNamespaceCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBlobLogByNamespaceCall) Do(f func(context.Context, uint64, storage.BlobLogFilters) ([]storage.BlobLog, error)) *IBlobLogByNamespaceCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBlobLogByNamespaceCall) DoAndReturn(f func(context.Context, uint64, storage.BlobLogFilters) ([]storage.BlobLog, error)) *IBlobLogByNamespaceCall {
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
func (mr *MockIBlobLogMockRecorder) ByProviders(ctx, providers, fltrs any) *IBlobLogByProvidersCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByProviders", reflect.TypeOf((*MockIBlobLog)(nil).ByProviders), ctx, providers, fltrs)
	return &IBlobLogByProvidersCall{Call: call}
}

// IBlobLogByProvidersCall wrap *gomock.Call
type IBlobLogByProvidersCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBlobLogByProvidersCall) Return(arg0 []storage.BlobLog, arg1 error) *IBlobLogByProvidersCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBlobLogByProvidersCall) Do(f func(context.Context, []storage.RollupProvider, storage.BlobLogFilters) ([]storage.BlobLog, error)) *IBlobLogByProvidersCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBlobLogByProvidersCall) DoAndReturn(f func(context.Context, []storage.RollupProvider, storage.BlobLogFilters) ([]storage.BlobLog, error)) *IBlobLogByProvidersCall {
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
func (mr *MockIBlobLogMockRecorder) BySigner(ctx, signerId, fltrs any) *IBlobLogBySignerCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BySigner", reflect.TypeOf((*MockIBlobLog)(nil).BySigner), ctx, signerId, fltrs)
	return &IBlobLogBySignerCall{Call: call}
}

// IBlobLogBySignerCall wrap *gomock.Call
type IBlobLogBySignerCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBlobLogBySignerCall) Return(arg0 []storage.BlobLog, arg1 error) *IBlobLogBySignerCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBlobLogBySignerCall) Do(f func(context.Context, uint64, storage.BlobLogFilters) ([]storage.BlobLog, error)) *IBlobLogBySignerCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBlobLogBySignerCall) DoAndReturn(f func(context.Context, uint64, storage.BlobLogFilters) ([]storage.BlobLog, error)) *IBlobLogBySignerCall {
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
func (mr *MockIBlobLogMockRecorder) ByTxId(ctx, txId, fltrs any) *IBlobLogByTxIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByTxId", reflect.TypeOf((*MockIBlobLog)(nil).ByTxId), ctx, txId, fltrs)
	return &IBlobLogByTxIdCall{Call: call}
}

// IBlobLogByTxIdCall wrap *gomock.Call
type IBlobLogByTxIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBlobLogByTxIdCall) Return(arg0 []storage.BlobLog, arg1 error) *IBlobLogByTxIdCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBlobLogByTxIdCall) Do(f func(context.Context, uint64, storage.BlobLogFilters) ([]storage.BlobLog, error)) *IBlobLogByTxIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBlobLogByTxIdCall) DoAndReturn(f func(context.Context, uint64, storage.BlobLogFilters) ([]storage.BlobLog, error)) *IBlobLogByTxIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CountByHeight mocks base method.
func (m *MockIBlobLog) CountByHeight(ctx context.Context, height types.Level) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountByHeight", ctx, height)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountByHeight indicates an expected call of CountByHeight.
func (mr *MockIBlobLogMockRecorder) CountByHeight(ctx, height any) *IBlobLogCountByHeightCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByHeight", reflect.TypeOf((*MockIBlobLog)(nil).CountByHeight), ctx, height)
	return &IBlobLogCountByHeightCall{Call: call}
}

// IBlobLogCountByHeightCall wrap *gomock.Call
type IBlobLogCountByHeightCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBlobLogCountByHeightCall) Return(arg0 int, arg1 error) *IBlobLogCountByHeightCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBlobLogCountByHeightCall) Do(f func(context.Context, types.Level) (int, error)) *IBlobLogCountByHeightCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBlobLogCountByHeightCall) DoAndReturn(f func(context.Context, types.Level) (int, error)) *IBlobLogCountByHeightCall {
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
func (mr *MockIBlobLogMockRecorder) CountByTxId(ctx, txId any) *IBlobLogCountByTxIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByTxId", reflect.TypeOf((*MockIBlobLog)(nil).CountByTxId), ctx, txId)
	return &IBlobLogCountByTxIdCall{Call: call}
}

// IBlobLogCountByTxIdCall wrap *gomock.Call
type IBlobLogCountByTxIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBlobLogCountByTxIdCall) Return(arg0 int, arg1 error) *IBlobLogCountByTxIdCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBlobLogCountByTxIdCall) Do(f func(context.Context, uint64) (int, error)) *IBlobLogCountByTxIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBlobLogCountByTxIdCall) DoAndReturn(f func(context.Context, uint64) (int, error)) *IBlobLogCountByTxIdCall {
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
func (mr *MockIBlobLogMockRecorder) CursorList(ctx, id, limit, order, cmp any) *IBlobLogCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIBlobLog)(nil).CursorList), ctx, id, limit, order, cmp)
	return &IBlobLogCursorListCall{Call: call}
}

// IBlobLogCursorListCall wrap *gomock.Call
type IBlobLogCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBlobLogCursorListCall) Return(arg0 []*storage.BlobLog, arg1 error) *IBlobLogCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBlobLogCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.BlobLog, error)) *IBlobLogCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBlobLogCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.BlobLog, error)) *IBlobLogCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ExportByProviders mocks base method.
func (m *MockIBlobLog) ExportByProviders(ctx context.Context, providers []storage.RollupProvider, stream io.Writer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExportByProviders", ctx, providers, stream)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExportByProviders indicates an expected call of ExportByProviders.
func (mr *MockIBlobLogMockRecorder) ExportByProviders(ctx, providers, stream any) *IBlobLogExportByProvidersCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExportByProviders", reflect.TypeOf((*MockIBlobLog)(nil).ExportByProviders), ctx, providers, stream)
	return &IBlobLogExportByProvidersCall{Call: call}
}

// IBlobLogExportByProvidersCall wrap *gomock.Call
type IBlobLogExportByProvidersCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBlobLogExportByProvidersCall) Return(err error) *IBlobLogExportByProvidersCall {
	c.Call = c.Call.Return(err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBlobLogExportByProvidersCall) Do(f func(context.Context, []storage.RollupProvider, io.Writer) error) *IBlobLogExportByProvidersCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBlobLogExportByProvidersCall) DoAndReturn(f func(context.Context, []storage.RollupProvider, io.Writer) error) *IBlobLogExportByProvidersCall {
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
func (mr *MockIBlobLogMockRecorder) GetByID(ctx, id any) *IBlobLogGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIBlobLog)(nil).GetByID), ctx, id)
	return &IBlobLogGetByIDCall{Call: call}
}

// IBlobLogGetByIDCall wrap *gomock.Call
type IBlobLogGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBlobLogGetByIDCall) Return(arg0 *storage.BlobLog, arg1 error) *IBlobLogGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBlobLogGetByIDCall) Do(f func(context.Context, uint64) (*storage.BlobLog, error)) *IBlobLogGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBlobLogGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.BlobLog, error)) *IBlobLogGetByIDCall {
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
func (mr *MockIBlobLogMockRecorder) IsNoRows(err any) *IBlobLogIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIBlobLog)(nil).IsNoRows), err)
	return &IBlobLogIsNoRowsCall{Call: call}
}

// IBlobLogIsNoRowsCall wrap *gomock.Call
type IBlobLogIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBlobLogIsNoRowsCall) Return(arg0 bool) *IBlobLogIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBlobLogIsNoRowsCall) Do(f func(error) bool) *IBlobLogIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBlobLogIsNoRowsCall) DoAndReturn(f func(error) bool) *IBlobLogIsNoRowsCall {
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
func (mr *MockIBlobLogMockRecorder) LastID(ctx any) *IBlobLogLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIBlobLog)(nil).LastID), ctx)
	return &IBlobLogLastIDCall{Call: call}
}

// IBlobLogLastIDCall wrap *gomock.Call
type IBlobLogLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBlobLogLastIDCall) Return(arg0 uint64, arg1 error) *IBlobLogLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBlobLogLastIDCall) Do(f func(context.Context) (uint64, error)) *IBlobLogLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBlobLogLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *IBlobLogLastIDCall {
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
func (mr *MockIBlobLogMockRecorder) List(ctx, limit, offset, order any) *IBlobLogListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIBlobLog)(nil).List), ctx, limit, offset, order)
	return &IBlobLogListCall{Call: call}
}

// IBlobLogListCall wrap *gomock.Call
type IBlobLogListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBlobLogListCall) Return(arg0 []*storage.BlobLog, arg1 error) *IBlobLogListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBlobLogListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.BlobLog, error)) *IBlobLogListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBlobLogListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.BlobLog, error)) *IBlobLogListCall {
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
func (mr *MockIBlobLogMockRecorder) Save(ctx, m any) *IBlobLogSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIBlobLog)(nil).Save), ctx, m)
	return &IBlobLogSaveCall{Call: call}
}

// IBlobLogSaveCall wrap *gomock.Call
type IBlobLogSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBlobLogSaveCall) Return(arg0 error) *IBlobLogSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBlobLogSaveCall) Do(f func(context.Context, *storage.BlobLog) error) *IBlobLogSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBlobLogSaveCall) DoAndReturn(f func(context.Context, *storage.BlobLog) error) *IBlobLogSaveCall {
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
func (mr *MockIBlobLogMockRecorder) Update(ctx, m any) *IBlobLogUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIBlobLog)(nil).Update), ctx, m)
	return &IBlobLogUpdateCall{Call: call}
}

// IBlobLogUpdateCall wrap *gomock.Call
type IBlobLogUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBlobLogUpdateCall) Return(arg0 error) *IBlobLogUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBlobLogUpdateCall) Do(f func(context.Context, *storage.BlobLog) error) *IBlobLogUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBlobLogUpdateCall) DoAndReturn(f func(context.Context, *storage.BlobLog) error) *IBlobLogUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
