// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/service/session.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entity "github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockSessionStorage is a mock of SessionStorage interface.
type MockSessionStorage struct {
	ctrl     *gomock.Controller
	recorder *MockSessionStorageMockRecorder
}

// MockSessionStorageMockRecorder is the mock recorder for MockSessionStorage.
type MockSessionStorageMockRecorder struct {
	mock *MockSessionStorage
}

// NewMockSessionStorage creates a new mock instance.
func NewMockSessionStorage(ctrl *gomock.Controller) *MockSessionStorage {
	mock := &MockSessionStorage{ctrl: ctrl}
	mock.recorder = &MockSessionStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionStorage) EXPECT() *MockSessionStorageMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockSessionStorage) Create(ctx context.Context, session entity.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, session)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockSessionStorageMockRecorder) Create(ctx, session interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSessionStorage)(nil).Create), ctx, session)
}

// Delete mocks base method.
func (m *MockSessionStorage) Delete(ctx context.Context, token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSessionStorageMockRecorder) Delete(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSessionStorage)(nil).Delete), ctx, token)
}

// GetByToken mocks base method.
func (m *MockSessionStorage) GetByToken(ctx context.Context, token string) (entity.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByToken", ctx, token)
	ret0, _ := ret[0].(entity.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByToken indicates an expected call of GetByToken.
func (mr *MockSessionStorageMockRecorder) GetByToken(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByToken", reflect.TypeOf((*MockSessionStorage)(nil).GetByToken), ctx, token)
}
