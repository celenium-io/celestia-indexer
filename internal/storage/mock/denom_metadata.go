// Code generated by MockGen. DO NOT EDIT.
// Source: denom_metadata.go
//
// Generated by this command:
//
//	mockgen -source=denom_metadata.go -destination=mock/denom_metadata.go -package=mock -typed
//
// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	storage "github.com/dipdup-io/celestia-indexer/internal/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockIDenomMetadata is a mock of IDenomMetadata interface.
type MockIDenomMetadata struct {
	ctrl     *gomock.Controller
	recorder *MockIDenomMetadataMockRecorder
}

// MockIDenomMetadataMockRecorder is the mock recorder for MockIDenomMetadata.
type MockIDenomMetadataMockRecorder struct {
	mock *MockIDenomMetadata
}

// NewMockIDenomMetadata creates a new mock instance.
func NewMockIDenomMetadata(ctrl *gomock.Controller) *MockIDenomMetadata {
	mock := &MockIDenomMetadata{ctrl: ctrl}
	mock.recorder = &MockIDenomMetadataMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIDenomMetadata) EXPECT() *MockIDenomMetadataMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *MockIDenomMetadata) All(ctx context.Context) ([]storage.DenomMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", ctx)
	ret0, _ := ret[0].([]storage.DenomMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockIDenomMetadataMockRecorder) All(ctx any) *IDenomMetadataAllCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockIDenomMetadata)(nil).All), ctx)
	return &IDenomMetadataAllCall{Call: call}
}

// IDenomMetadataAllCall wrap *gomock.Call
type IDenomMetadataAllCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IDenomMetadataAllCall) Return(arg0 []storage.DenomMetadata, arg1 error) *IDenomMetadataAllCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IDenomMetadataAllCall) Do(f func(context.Context) ([]storage.DenomMetadata, error)) *IDenomMetadataAllCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IDenomMetadataAllCall) DoAndReturn(f func(context.Context) ([]storage.DenomMetadata, error)) *IDenomMetadataAllCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}