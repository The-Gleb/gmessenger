// Code generated by MockGen. DO NOT EDIT.
// Source: app/internal/controller/http/v1/handler/chats.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entity "github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockChatsUsecase is a mock of ChatsUsecase interface.
type MockChatsUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockChatsUsecaseMockRecorder
}

// MockChatsUsecaseMockRecorder is the mock recorder for MockChatsUsecase.
type MockChatsUsecaseMockRecorder struct {
	mock *MockChatsUsecase
}

// NewMockChatsUsecase creates a new mock instance.
func NewMockChatsUsecase(ctrl *gomock.Controller) *MockChatsUsecase {
	mock := &MockChatsUsecase{ctrl: ctrl}
	mock.recorder = &MockChatsUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChatsUsecase) EXPECT() *MockChatsUsecaseMockRecorder {
	return m.recorder
}

// ShowChats mocks base method.
func (m *MockChatsUsecase) ShowChats(ctx context.Context, login string) ([]entity.Chat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShowChats", ctx, login)
	ret0, _ := ret[0].([]entity.Chat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShowChats indicates an expected call of ShowChats.
func (mr *MockChatsUsecaseMockRecorder) ShowChats(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowChats", reflect.TypeOf((*MockChatsUsecase)(nil).ShowChats), ctx, login)
}
