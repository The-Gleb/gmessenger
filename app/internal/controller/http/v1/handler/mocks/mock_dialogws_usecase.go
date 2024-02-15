// Code generated by MockGen. DO NOT EDIT.
// Source: app/internal/controller/http/v1/handler/dialogws.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	dialogws_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/dialogws.go"
	gomock "github.com/golang/mock/gomock"
)

// MockDialogWSUsecase is a mock of DialogWSUsecase interface.
type MockDialogWSUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockDialogWSUsecaseMockRecorder
}

// MockDialogWSUsecaseMockRecorder is the mock recorder for MockDialogWSUsecase.
type MockDialogWSUsecaseMockRecorder struct {
	mock *MockDialogWSUsecase
}

// NewMockDialogWSUsecase creates a new mock instance.
func NewMockDialogWSUsecase(ctrl *gomock.Controller) *MockDialogWSUsecase {
	mock := &MockDialogWSUsecase{ctrl: ctrl}
	mock.recorder = &MockDialogWSUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDialogWSUsecase) EXPECT() *MockDialogWSUsecaseMockRecorder {
	return m.recorder
}

// OpenDialog mocks base method.
func (m *MockDialogWSUsecase) OpenDialog(ctx context.Context, dto dialogws_usecase.OpenDialogDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenDialog", ctx, dto)
	ret0, _ := ret[0].(error)
	return ret0
}

// OpenDialog indicates an expected call of OpenDialog.
func (mr *MockDialogWSUsecaseMockRecorder) OpenDialog(ctx, dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenDialog", reflect.TypeOf((*MockDialogWSUsecase)(nil).OpenDialog), ctx, dto)
}
