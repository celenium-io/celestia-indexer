// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: block.go
//
// Generated by this command:
//
//	mockgen -source=block.go -destination=mock/block.go -package=mock -typed
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

// MockIBlock is a mock of IBlock interface.
type MockIBlock struct {
	ctrl     *gomock.Controller
	recorder *MockIBlockMockRecorder
}

// MockIBlockMockRecorder is the mock recorder for MockIBlock.
type MockIBlockMockRecorder struct {
	mock *MockIBlock
}

// NewMockIBlock creates a new mock instance.
func NewMockIBlock(ctrl *gomock.Controller) *MockIBlock {
	mock := &MockIBlock{ctrl: ctrl}
	mock.recorder = &MockIBlockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIBlock) EXPECT() *MockIBlockMockRecorder {
	return m.recorder
}

// ByHash mocks base method.
func (m *MockIBlock) ByHash(ctx context.Context, hash []byte) (storage.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByHash", ctx, hash)
	ret0, _ := ret[0].(storage.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByHash indicates an expected call of ByHash.
func (mr *MockIBlockMockRecorder) ByHash(ctx, hash any) *MockIBlockByHashCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByHash", reflect.TypeOf((*MockIBlock)(nil).ByHash), ctx, hash)
	return &MockIBlockByHashCall{Call: call}
}

// MockIBlockByHashCall wrap *gomock.Call
type MockIBlockByHashCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockByHashCall) Return(arg0 storage.Block, arg1 error) *MockIBlockByHashCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockByHashCall) Do(f func(context.Context, []byte) (storage.Block, error)) *MockIBlockByHashCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockByHashCall) DoAndReturn(f func(context.Context, []byte) (storage.Block, error)) *MockIBlockByHashCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ByHeight mocks base method.
func (m *MockIBlock) ByHeight(ctx context.Context, height types.Level) (storage.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByHeight", ctx, height)
	ret0, _ := ret[0].(storage.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByHeight indicates an expected call of ByHeight.
func (mr *MockIBlockMockRecorder) ByHeight(ctx, height any) *MockIBlockByHeightCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByHeight", reflect.TypeOf((*MockIBlock)(nil).ByHeight), ctx, height)
	return &MockIBlockByHeightCall{Call: call}
}

// MockIBlockByHeightCall wrap *gomock.Call
type MockIBlockByHeightCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockByHeightCall) Return(arg0 storage.Block, arg1 error) *MockIBlockByHeightCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockByHeightCall) Do(f func(context.Context, types.Level) (storage.Block, error)) *MockIBlockByHeightCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockByHeightCall) DoAndReturn(f func(context.Context, types.Level) (storage.Block, error)) *MockIBlockByHeightCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ByHeightWithStats mocks base method.
func (m *MockIBlock) ByHeightWithStats(ctx context.Context, height types.Level) (storage.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByHeightWithStats", ctx, height)
	ret0, _ := ret[0].(storage.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByHeightWithStats indicates an expected call of ByHeightWithStats.
func (mr *MockIBlockMockRecorder) ByHeightWithStats(ctx, height any) *MockIBlockByHeightWithStatsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByHeightWithStats", reflect.TypeOf((*MockIBlock)(nil).ByHeightWithStats), ctx, height)
	return &MockIBlockByHeightWithStatsCall{Call: call}
}

// MockIBlockByHeightWithStatsCall wrap *gomock.Call
type MockIBlockByHeightWithStatsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockByHeightWithStatsCall) Return(arg0 storage.Block, arg1 error) *MockIBlockByHeightWithStatsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockByHeightWithStatsCall) Do(f func(context.Context, types.Level) (storage.Block, error)) *MockIBlockByHeightWithStatsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockByHeightWithStatsCall) DoAndReturn(f func(context.Context, types.Level) (storage.Block, error)) *MockIBlockByHeightWithStatsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ByIdWithRelations mocks base method.
func (m *MockIBlock) ByIdWithRelations(ctx context.Context, id uint64) (storage.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByIdWithRelations", ctx, id)
	ret0, _ := ret[0].(storage.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByIdWithRelations indicates an expected call of ByIdWithRelations.
func (mr *MockIBlockMockRecorder) ByIdWithRelations(ctx, id any) *MockIBlockByIdWithRelationsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByIdWithRelations", reflect.TypeOf((*MockIBlock)(nil).ByIdWithRelations), ctx, id)
	return &MockIBlockByIdWithRelationsCall{Call: call}
}

// MockIBlockByIdWithRelationsCall wrap *gomock.Call
type MockIBlockByIdWithRelationsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockByIdWithRelationsCall) Return(arg0 storage.Block, arg1 error) *MockIBlockByIdWithRelationsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockByIdWithRelationsCall) Do(f func(context.Context, uint64) (storage.Block, error)) *MockIBlockByIdWithRelationsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockByIdWithRelationsCall) DoAndReturn(f func(context.Context, uint64) (storage.Block, error)) *MockIBlockByIdWithRelationsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ByProposer mocks base method.
func (m *MockIBlock) ByProposer(ctx context.Context, proposerId uint64, limit, offset int) ([]storage.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByProposer", ctx, proposerId, limit, offset)
	ret0, _ := ret[0].([]storage.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByProposer indicates an expected call of ByProposer.
func (mr *MockIBlockMockRecorder) ByProposer(ctx, proposerId, limit, offset any) *MockIBlockByProposerCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByProposer", reflect.TypeOf((*MockIBlock)(nil).ByProposer), ctx, proposerId, limit, offset)
	return &MockIBlockByProposerCall{Call: call}
}

// MockIBlockByProposerCall wrap *gomock.Call
type MockIBlockByProposerCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockByProposerCall) Return(arg0 []storage.Block, arg1 error) *MockIBlockByProposerCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockByProposerCall) Do(f func(context.Context, uint64, int, int) ([]storage.Block, error)) *MockIBlockByProposerCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockByProposerCall) DoAndReturn(f func(context.Context, uint64, int, int) ([]storage.Block, error)) *MockIBlockByProposerCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CursorList mocks base method.
func (m *MockIBlock) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIBlockMockRecorder) CursorList(ctx, id, limit, order, cmp any) *MockIBlockCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIBlock)(nil).CursorList), ctx, id, limit, order, cmp)
	return &MockIBlockCursorListCall{Call: call}
}

// MockIBlockCursorListCall wrap *gomock.Call
type MockIBlockCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockCursorListCall) Return(arg0 []*storage.Block, arg1 error) *MockIBlockCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Block, error)) *MockIBlockCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Block, error)) *MockIBlockCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIBlock) GetByID(ctx context.Context, id uint64) (*storage.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIBlockMockRecorder) GetByID(ctx, id any) *MockIBlockGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIBlock)(nil).GetByID), ctx, id)
	return &MockIBlockGetByIDCall{Call: call}
}

// MockIBlockGetByIDCall wrap *gomock.Call
type MockIBlockGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockGetByIDCall) Return(arg0 *storage.Block, arg1 error) *MockIBlockGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockGetByIDCall) Do(f func(context.Context, uint64) (*storage.Block, error)) *MockIBlockGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.Block, error)) *MockIBlockGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIBlock) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIBlockMockRecorder) IsNoRows(err any) *MockIBlockIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIBlock)(nil).IsNoRows), err)
	return &MockIBlockIsNoRowsCall{Call: call}
}

// MockIBlockIsNoRowsCall wrap *gomock.Call
type MockIBlockIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockIsNoRowsCall) Return(arg0 bool) *MockIBlockIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockIsNoRowsCall) Do(f func(error) bool) *MockIBlockIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockIsNoRowsCall) DoAndReturn(f func(error) bool) *MockIBlockIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Last mocks base method.
func (m *MockIBlock) Last(ctx context.Context) (storage.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Last", ctx)
	ret0, _ := ret[0].(storage.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Last indicates an expected call of Last.
func (mr *MockIBlockMockRecorder) Last(ctx any) *MockIBlockLastCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Last", reflect.TypeOf((*MockIBlock)(nil).Last), ctx)
	return &MockIBlockLastCall{Call: call}
}

// MockIBlockLastCall wrap *gomock.Call
type MockIBlockLastCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockLastCall) Return(arg0 storage.Block, arg1 error) *MockIBlockLastCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockLastCall) Do(f func(context.Context) (storage.Block, error)) *MockIBlockLastCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockLastCall) DoAndReturn(f func(context.Context) (storage.Block, error)) *MockIBlockLastCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIBlock) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIBlockMockRecorder) LastID(ctx any) *MockIBlockLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIBlock)(nil).LastID), ctx)
	return &MockIBlockLastIDCall{Call: call}
}

// MockIBlockLastIDCall wrap *gomock.Call
type MockIBlockLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockLastIDCall) Return(arg0 uint64, arg1 error) *MockIBlockLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockLastIDCall) Do(f func(context.Context) (uint64, error)) *MockIBlockLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *MockIBlockLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIBlock) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIBlockMockRecorder) List(ctx, limit, offset, order any) *MockIBlockListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIBlock)(nil).List), ctx, limit, offset, order)
	return &MockIBlockListCall{Call: call}
}

// MockIBlockListCall wrap *gomock.Call
type MockIBlockListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockListCall) Return(arg0 []*storage.Block, arg1 error) *MockIBlockListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Block, error)) *MockIBlockListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Block, error)) *MockIBlockListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ListWithStats mocks base method.
func (m *MockIBlock) ListWithStats(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWithStats", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWithStats indicates an expected call of ListWithStats.
func (mr *MockIBlockMockRecorder) ListWithStats(ctx, limit, offset, order any) *MockIBlockListWithStatsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWithStats", reflect.TypeOf((*MockIBlock)(nil).ListWithStats), ctx, limit, offset, order)
	return &MockIBlockListWithStatsCall{Call: call}
}

// MockIBlockListWithStatsCall wrap *gomock.Call
type MockIBlockListWithStatsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockListWithStatsCall) Return(arg0 []*storage.Block, arg1 error) *MockIBlockListWithStatsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockListWithStatsCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Block, error)) *MockIBlockListWithStatsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockListWithStatsCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Block, error)) *MockIBlockListWithStatsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIBlock) Save(ctx context.Context, m *storage.Block) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIBlockMockRecorder) Save(ctx, m any) *MockIBlockSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIBlock)(nil).Save), ctx, m)
	return &MockIBlockSaveCall{Call: call}
}

// MockIBlockSaveCall wrap *gomock.Call
type MockIBlockSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockSaveCall) Return(arg0 error) *MockIBlockSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockSaveCall) Do(f func(context.Context, *storage.Block) error) *MockIBlockSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockSaveCall) DoAndReturn(f func(context.Context, *storage.Block) error) *MockIBlockSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Time mocks base method.
func (m *MockIBlock) Time(ctx context.Context, height types.Level) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Time", ctx, height)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Time indicates an expected call of Time.
func (mr *MockIBlockMockRecorder) Time(ctx, height any) *MockIBlockTimeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Time", reflect.TypeOf((*MockIBlock)(nil).Time), ctx, height)
	return &MockIBlockTimeCall{Call: call}
}

// MockIBlockTimeCall wrap *gomock.Call
type MockIBlockTimeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockTimeCall) Return(arg0 time.Time, arg1 error) *MockIBlockTimeCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockTimeCall) Do(f func(context.Context, types.Level) (time.Time, error)) *MockIBlockTimeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockTimeCall) DoAndReturn(f func(context.Context, types.Level) (time.Time, error)) *MockIBlockTimeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIBlock) Update(ctx context.Context, m *storage.Block) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIBlockMockRecorder) Update(ctx, m any) *MockIBlockUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIBlock)(nil).Update), ctx, m)
	return &MockIBlockUpdateCall{Call: call}
}

// MockIBlockUpdateCall wrap *gomock.Call
type MockIBlockUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockUpdateCall) Return(arg0 error) *MockIBlockUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockUpdateCall) Do(f func(context.Context, *storage.Block) error) *MockIBlockUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockUpdateCall) DoAndReturn(f func(context.Context, *storage.Block) error) *MockIBlockUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
