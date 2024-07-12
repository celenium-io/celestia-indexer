// Code generated by MockGen. DO NOT EDIT.
// Source: tx.go
//
// Generated by this command:
//
//	mockgen -source=tx.go -destination=mock/tx.go -package=mock -typed
//
// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	time "time"

	storage "github.com/celenium-io/celestia-indexer/internal/storage"
	types "github.com/celenium-io/celestia-indexer/pkg/types"
	storage0 "github.com/dipdup-net/indexer-sdk/pkg/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockITx is a mock of ITx interface.
type MockITx struct {
	ctrl     *gomock.Controller
	recorder *MockITxMockRecorder
}

// MockITxMockRecorder is the mock recorder for MockITx.
type MockITxMockRecorder struct {
	mock *MockITx
}

// NewMockITx creates a new mock instance.
func NewMockITx(ctrl *gomock.Controller) *MockITx {
	mock := &MockITx{ctrl: ctrl}
	mock.recorder = &MockITxMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockITx) EXPECT() *MockITxMockRecorder {
	return m.recorder
}

// ByAddress mocks base method.
func (m *MockITx) ByAddress(ctx context.Context, addressId uint64, fltrs storage.TxFilter) ([]storage.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByAddress", ctx, addressId, fltrs)
	ret0, _ := ret[0].([]storage.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByAddress indicates an expected call of ByAddress.
func (mr *MockITxMockRecorder) ByAddress(ctx, addressId, fltrs any) *ITxByAddressCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByAddress", reflect.TypeOf((*MockITx)(nil).ByAddress), ctx, addressId, fltrs)
	return &ITxByAddressCall{Call: call}
}

// ITxByAddressCall wrap *gomock.Call
type ITxByAddressCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ITxByAddressCall) Return(arg0 []storage.Tx, arg1 error) *ITxByAddressCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ITxByAddressCall) Do(f func(context.Context, uint64, storage.TxFilter) ([]storage.Tx, error)) *ITxByAddressCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ITxByAddressCall) DoAndReturn(f func(context.Context, uint64, storage.TxFilter) ([]storage.Tx, error)) *ITxByAddressCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ByHash mocks base method.
func (m *MockITx) ByHash(ctx context.Context, hash []byte) (storage.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByHash", ctx, hash)
	ret0, _ := ret[0].(storage.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByHash indicates an expected call of ByHash.
func (mr *MockITxMockRecorder) ByHash(ctx, hash any) *ITxByHashCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByHash", reflect.TypeOf((*MockITx)(nil).ByHash), ctx, hash)
	return &ITxByHashCall{Call: call}
}

// ITxByHashCall wrap *gomock.Call
type ITxByHashCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ITxByHashCall) Return(arg0 storage.Tx, arg1 error) *ITxByHashCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ITxByHashCall) Do(f func(context.Context, []byte) (storage.Tx, error)) *ITxByHashCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ITxByHashCall) DoAndReturn(f func(context.Context, []byte) (storage.Tx, error)) *ITxByHashCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ByIdWithRelations mocks base method.
func (m *MockITx) ByIdWithRelations(ctx context.Context, id uint64) (storage.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByIdWithRelations", ctx, id)
	ret0, _ := ret[0].(storage.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByIdWithRelations indicates an expected call of ByIdWithRelations.
func (mr *MockITxMockRecorder) ByIdWithRelations(ctx, id any) *ITxByIdWithRelationsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByIdWithRelations", reflect.TypeOf((*MockITx)(nil).ByIdWithRelations), ctx, id)
	return &ITxByIdWithRelationsCall{Call: call}
}

// ITxByIdWithRelationsCall wrap *gomock.Call
type ITxByIdWithRelationsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ITxByIdWithRelationsCall) Return(arg0 storage.Tx, arg1 error) *ITxByIdWithRelationsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ITxByIdWithRelationsCall) Do(f func(context.Context, uint64) (storage.Tx, error)) *ITxByIdWithRelationsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ITxByIdWithRelationsCall) DoAndReturn(f func(context.Context, uint64) (storage.Tx, error)) *ITxByIdWithRelationsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CursorList mocks base method.
func (m *MockITx) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockITxMockRecorder) CursorList(ctx, id, limit, order, cmp any) *ITxCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockITx)(nil).CursorList), ctx, id, limit, order, cmp)
	return &ITxCursorListCall{Call: call}
}

// ITxCursorListCall wrap *gomock.Call
type ITxCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ITxCursorListCall) Return(arg0 []*storage.Tx, arg1 error) *ITxCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ITxCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Tx, error)) *ITxCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ITxCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Tx, error)) *ITxCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Filter mocks base method.
func (m *MockITx) Filter(ctx context.Context, fltrs storage.TxFilter) ([]storage.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Filter", ctx, fltrs)
	ret0, _ := ret[0].([]storage.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Filter indicates an expected call of Filter.
func (mr *MockITxMockRecorder) Filter(ctx, fltrs any) *ITxFilterCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Filter", reflect.TypeOf((*MockITx)(nil).Filter), ctx, fltrs)
	return &ITxFilterCall{Call: call}
}

// ITxFilterCall wrap *gomock.Call
type ITxFilterCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ITxFilterCall) Return(arg0 []storage.Tx, arg1 error) *ITxFilterCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ITxFilterCall) Do(f func(context.Context, storage.TxFilter) ([]storage.Tx, error)) *ITxFilterCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ITxFilterCall) DoAndReturn(f func(context.Context, storage.TxFilter) ([]storage.Tx, error)) *ITxFilterCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Gas mocks base method.
func (m *MockITx) Gas(ctx context.Context, height types.Level, ts time.Time) ([]storage.Gas, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Gas", ctx, height, ts)
	ret0, _ := ret[0].([]storage.Gas)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Gas indicates an expected call of Gas.
func (mr *MockITxMockRecorder) Gas(ctx, height, ts any) *ITxGasCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Gas", reflect.TypeOf((*MockITx)(nil).Gas), ctx, height, ts)
	return &ITxGasCall{Call: call}
}

// ITxGasCall wrap *gomock.Call
type ITxGasCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ITxGasCall) Return(arg0 []storage.Gas, arg1 error) *ITxGasCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ITxGasCall) Do(f func(context.Context, types.Level, time.Time) ([]storage.Gas, error)) *ITxGasCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ITxGasCall) DoAndReturn(f func(context.Context, types.Level, time.Time) ([]storage.Gas, error)) *ITxGasCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Genesis mocks base method.
func (m *MockITx) Genesis(ctx context.Context, limit, offset int, sortOrder storage0.SortOrder) ([]storage.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Genesis", ctx, limit, offset, sortOrder)
	ret0, _ := ret[0].([]storage.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Genesis indicates an expected call of Genesis.
func (mr *MockITxMockRecorder) Genesis(ctx, limit, offset, sortOrder any) *ITxGenesisCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Genesis", reflect.TypeOf((*MockITx)(nil).Genesis), ctx, limit, offset, sortOrder)
	return &ITxGenesisCall{Call: call}
}

// ITxGenesisCall wrap *gomock.Call
type ITxGenesisCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ITxGenesisCall) Return(arg0 []storage.Tx, arg1 error) *ITxGenesisCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ITxGenesisCall) Do(f func(context.Context, int, int, storage0.SortOrder) ([]storage.Tx, error)) *ITxGenesisCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ITxGenesisCall) DoAndReturn(f func(context.Context, int, int, storage0.SortOrder) ([]storage.Tx, error)) *ITxGenesisCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockITx) GetByID(ctx context.Context, id uint64) (*storage.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockITxMockRecorder) GetByID(ctx, id any) *ITxGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockITx)(nil).GetByID), ctx, id)
	return &ITxGetByIDCall{Call: call}
}

// ITxGetByIDCall wrap *gomock.Call
type ITxGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ITxGetByIDCall) Return(arg0 *storage.Tx, arg1 error) *ITxGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ITxGetByIDCall) Do(f func(context.Context, uint64) (*storage.Tx, error)) *ITxGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ITxGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.Tx, error)) *ITxGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IdByHash mocks base method.
func (m *MockITx) IdByHash(ctx context.Context, hash []byte) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IdByHash", ctx, hash)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IdByHash indicates an expected call of IdByHash.
func (mr *MockITxMockRecorder) IdByHash(ctx, hash any) *ITxIdByHashCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IdByHash", reflect.TypeOf((*MockITx)(nil).IdByHash), ctx, hash)
	return &ITxIdByHashCall{Call: call}
}

// ITxIdByHashCall wrap *gomock.Call
type ITxIdByHashCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ITxIdByHashCall) Return(arg0 uint64, arg1 error) *ITxIdByHashCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ITxIdByHashCall) Do(f func(context.Context, []byte) (uint64, error)) *ITxIdByHashCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ITxIdByHashCall) DoAndReturn(f func(context.Context, []byte) (uint64, error)) *ITxIdByHashCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockITx) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockITxMockRecorder) IsNoRows(err any) *ITxIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockITx)(nil).IsNoRows), err)
	return &ITxIsNoRowsCall{Call: call}
}

// ITxIsNoRowsCall wrap *gomock.Call
type ITxIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ITxIsNoRowsCall) Return(arg0 bool) *ITxIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ITxIsNoRowsCall) Do(f func(error) bool) *ITxIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ITxIsNoRowsCall) DoAndReturn(f func(error) bool) *ITxIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockITx) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockITxMockRecorder) LastID(ctx any) *ITxLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockITx)(nil).LastID), ctx)
	return &ITxLastIDCall{Call: call}
}

// ITxLastIDCall wrap *gomock.Call
type ITxLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ITxLastIDCall) Return(arg0 uint64, arg1 error) *ITxLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ITxLastIDCall) Do(f func(context.Context) (uint64, error)) *ITxLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ITxLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *ITxLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockITx) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockITxMockRecorder) List(ctx, limit, offset, order any) *ITxListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockITx)(nil).List), ctx, limit, offset, order)
	return &ITxListCall{Call: call}
}

// ITxListCall wrap *gomock.Call
type ITxListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ITxListCall) Return(arg0 []*storage.Tx, arg1 error) *ITxListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ITxListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Tx, error)) *ITxListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ITxListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Tx, error)) *ITxListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockITx) Save(ctx context.Context, m *storage.Tx) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockITxMockRecorder) Save(ctx, m any) *ITxSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockITx)(nil).Save), ctx, m)
	return &ITxSaveCall{Call: call}
}

// ITxSaveCall wrap *gomock.Call
type ITxSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ITxSaveCall) Return(arg0 error) *ITxSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ITxSaveCall) Do(f func(context.Context, *storage.Tx) error) *ITxSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ITxSaveCall) DoAndReturn(f func(context.Context, *storage.Tx) error) *ITxSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockITx) Update(ctx context.Context, m *storage.Tx) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockITxMockRecorder) Update(ctx, m any) *ITxUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockITx)(nil).Update), ctx, m)
	return &ITxUpdateCall{Call: call}
}

// ITxUpdateCall wrap *gomock.Call
type ITxUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ITxUpdateCall) Return(arg0 error) *ITxUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ITxUpdateCall) Do(f func(context.Context, *storage.Tx) error) *ITxUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ITxUpdateCall) DoAndReturn(f func(context.Context, *storage.Tx) error) *ITxUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
