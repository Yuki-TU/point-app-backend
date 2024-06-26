// Code generated by MockGen. DO NOT EDIT.
// Source: ./batch/controller/usecase.go
//
// Generated by this command:
//
//	mockgen -source=./batch/controller/usecase.go -destination=./batch/controller/_mock/mock_usecase.go
//

// Package mock_controller is a generated GoMock package.
package mock_controller

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockResetSendablePointer is a mock of ResetSendablePointer interface.
type MockResetSendablePointer struct {
	ctrl     *gomock.Controller
	recorder *MockResetSendablePointerMockRecorder
}

// MockResetSendablePointerMockRecorder is the mock recorder for MockResetSendablePointer.
type MockResetSendablePointerMockRecorder struct {
	mock *MockResetSendablePointer
}

// NewMockResetSendablePointer creates a new mock instance.
func NewMockResetSendablePointer(ctrl *gomock.Controller) *MockResetSendablePointer {
	mock := &MockResetSendablePointer{ctrl: ctrl}
	mock.recorder = &MockResetSendablePointerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResetSendablePointer) EXPECT() *MockResetSendablePointerMockRecorder {
	return m.recorder
}

// ResetPoint mocks base method.
func (m *MockResetSendablePointer) ResetPoint(ctx context.Context, initialSendablePoint int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetPoint", ctx, initialSendablePoint)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetPoint indicates an expected call of ResetPoint.
func (mr *MockResetSendablePointerMockRecorder) ResetPoint(ctx, initialSendablePoint any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetPoint", reflect.TypeOf((*MockResetSendablePointer)(nil).ResetPoint), ctx, initialSendablePoint)
}
