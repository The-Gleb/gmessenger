// Code generated by MockGen. DO NOT EDIT.
// Source: app/internal/controller/http/v1/handler/register.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entity "github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	register_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/register"
	gomock "github.com/golang/mock/gomock"
)

// MockRegisterUsecase is a mock of RegisterUsecase interface.
type MockRegisterUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockRegisterUsecaseMockRecorder
}

// MockRegisterUsecaseMockRecorder is the mock recorder for MockRegisterUsecase.
type MockRegisterUsecaseMockRecorder struct {
	mock *MockRegisterUsecase
}

// NewMockRegisterUsecase creates a new mock instance.
func NewMockRegisterUsecase(ctrl *gomock.Controller) *MockRegisterUsecase {
	mock := &MockRegisterUsecase{ctrl: ctrl}
	mock.recorder = &MockRegisterUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRegisterUsecase) EXPECT() *MockRegisterUsecaseMockRecorder {
	return m.recorder
}

// Register mocks base method.
func (m *MockRegisterUsecase) Register(ctx context.Context, usecaseDTO register_usecase.RegisterUserDTO) (entity.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, usecaseDTO)
	ret0, _ := ret[0].(entity.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockRegisterUsecaseMockRecorder) Register(ctx, usecaseDTO interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockRegisterUsecase)(nil).Register), ctx, usecaseDTO)
}
