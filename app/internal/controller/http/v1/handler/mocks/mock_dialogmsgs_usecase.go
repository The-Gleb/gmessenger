// Code generated by MockGen. DO NOT EDIT.
// Source: app/internal/controller/http/v1/handler/dialogmsgs.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entity "github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	dialogmsgs_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/dialogmsgs"
	gomock "github.com/golang/mock/gomock"
)

// MockDialogMsgsUsecase is a mock of DialogMsgsUsecase interface.
type MockDialogMsgsUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockDialogMsgsUsecaseMockRecorder
}

// MockDialogMsgsUsecaseMockRecorder is the mock recorder for MockDialogMsgsUsecase.
type MockDialogMsgsUsecaseMockRecorder struct {
	mock *MockDialogMsgsUsecase
}

// NewMockDialogMsgsUsecase creates a new mock instance.
func NewMockDialogMsgsUsecase(ctrl *gomock.Controller) *MockDialogMsgsUsecase {
	mock := &MockDialogMsgsUsecase{ctrl: ctrl}
	mock.recorder = &MockDialogMsgsUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDialogMsgsUsecase) EXPECT() *MockDialogMsgsUsecaseMockRecorder {
	return m.recorder
}

// GetDialogMessages mocks base method.
func (m *MockDialogMsgsUsecase) GetDialogMessages(ctx context.Context, dto dialogmsgs_usecase.GetDialogMessagesDTO) ([]entity.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDialogMessages", ctx, dto)
	ret0, _ := ret[0].([]entity.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDialogMessages indicates an expected call of GetDialogMessages.
func (mr *MockDialogMsgsUsecaseMockRecorder) GetDialogMessages(ctx, dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDialogMessages", reflect.TypeOf((*MockDialogMsgsUsecase)(nil).GetDialogMessages), ctx, dto)
}
