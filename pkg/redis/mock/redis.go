// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/redis/redis.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockRedisMethod is a mock of RedisMethod interface.
type MockRedisMethod struct {
	ctrl     *gomock.Controller
	recorder *MockRedisMethodMockRecorder
}

// MockRedisMethodMockRecorder is the mock recorder for MockRedisMethod.
type MockRedisMethodMockRecorder struct {
	mock *MockRedisMethod
}

// NewMockRedisMethod creates a new mock instance.
func NewMockRedisMethod(ctrl *gomock.Controller) *MockRedisMethod {
	mock := &MockRedisMethod{ctrl: ctrl}
	mock.recorder = &MockRedisMethodMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRedisMethod) EXPECT() *MockRedisMethodMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockRedisMethod) Get(key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRedisMethodMockRecorder) Get(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRedisMethod)(nil).Get), key)
}

// Set mocks base method.
func (m *MockRedisMethod) Set(key string, value interface{}, expiration time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", key, value, expiration)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockRedisMethodMockRecorder) Set(key, value, expiration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockRedisMethod)(nil).Set), key, value, expiration)
}
