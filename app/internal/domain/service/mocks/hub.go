// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/service/client/client.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	entity "github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	client "github.com/The-Gleb/gmessenger/app/internal/domain/service/client"
	gomock "github.com/golang/mock/gomock"
)

// MockHub is a mock of Hub interface.
type MockHub struct {
	ctrl     *gomock.Controller
	recorder *MockHubMockRecorder
}

// MockHubMockRecorder is the mock recorder for MockHub.
type MockHubMockRecorder struct {
	mock *MockHub
}

// NewMockHub creates a new mock instance.
func NewMockHub(ctrl *gomock.Controller) *MockHub {
	mock := &MockHub{ctrl: ctrl}
	mock.recorder = &MockHubMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHub) EXPECT() *MockHubMockRecorder {
	return m.recorder
}

// AddClient mocks base method.
func (m *MockHub) AddClient(c *client.Client) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddClient", c)
}

// AddClient indicates an expected call of AddClient.
func (mr *MockHubMockRecorder) AddClient(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddClient", reflect.TypeOf((*MockHub)(nil).AddClient), c)
}

// RemoveClient mocks base method.
func (m *MockHub) RemoveClient(c *client.Client) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveClient", c)
}

// RemoveClient indicates an expected call of RemoveClient.
func (mr *MockHubMockRecorder) RemoveClient(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveClient", reflect.TypeOf((*MockHub)(nil).RemoveClient), c)
}

// RouteEvent mocks base method.
func (m *MockHub) RouteEvent(event entity.Event, senderClient *client.Client) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RouteEvent", event, senderClient)
}

// RouteEvent indicates an expected call of RouteEvent.
func (mr *MockHubMockRecorder) RouteEvent(event, senderClient interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RouteEvent", reflect.TypeOf((*MockHub)(nil).RouteEvent), event, senderClient)
}
