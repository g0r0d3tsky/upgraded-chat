// Code generated by MockGen. DO NOT EDIT.
// Source: message.go
//
// Generated by this command:
//
//	mockgen -source=message.go -destination=mocks/repositoryMock.go
//

// Package mock_cache is a generated GoMock package.
package mock_cache

import (
	domain "2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/domain"
	context "context"
	reflect "reflect"

	redis "github.com/redis/go-redis/v9"
	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetAmountMessage mocks base method.
func (m *MockRepository) GetAmountMessage(ctx context.Context, amount int) ([]*domain.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAmountMessage", ctx, amount)
	ret0, _ := ret[0].([]*domain.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAmountMessage indicates an expected call of GetAmountMessage.
func (mr *MockRepositoryMockRecorder) GetAmountMessage(ctx, amount any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAmountMessage", reflect.TypeOf((*MockRepository)(nil).GetAmountMessage), ctx, amount)
}

// MockRedis is a mock of Redis interface.
type MockRedis struct {
	ctrl     *gomock.Controller
	recorder *MockRedisMockRecorder
}

// MockRedisMockRecorder is the mock recorder for MockRedis.
type MockRedisMockRecorder struct {
	mock *MockRedis
}

// NewMockRedis creates a new mock instance.
func NewMockRedis(ctrl *gomock.Controller) *MockRedis {
	mock := &MockRedis{ctrl: ctrl}
	mock.recorder = &MockRedisMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRedis) EXPECT() *MockRedisMockRecorder {
	return m.recorder
}

// LLen mocks base method.
func (m *MockRedis) LLen(ctx context.Context, key string) *redis.IntCmd {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LLen", ctx, key)
	ret0, _ := ret[0].(*redis.IntCmd)
	return ret0
}

// LLen indicates an expected call of LLen.
func (mr *MockRedisMockRecorder) LLen(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LLen", reflect.TypeOf((*MockRedis)(nil).LLen), ctx, key)
}

// LRange mocks base method.
func (m *MockRedis) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LRange", ctx, key, start, stop)
	ret0, _ := ret[0].(*redis.StringSliceCmd)
	return ret0
}

// LRange indicates an expected call of LRange.
func (mr *MockRedisMockRecorder) LRange(ctx, key, start, stop any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LRange", reflect.TypeOf((*MockRedis)(nil).LRange), ctx, key, start, stop)
}

// RPush mocks base method.
func (m *MockRedis) RPush(ctx context.Context, key string, values ...any) *redis.IntCmd {
	m.ctrl.T.Helper()
	varargs := []any{ctx, key}
	for _, a := range values {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RPush", varargs...)
	ret0, _ := ret[0].(*redis.IntCmd)
	return ret0
}

// RPush indicates an expected call of RPush.
func (mr *MockRedisMockRecorder) RPush(ctx, key any, values ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, key}, values...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RPush", reflect.TypeOf((*MockRedis)(nil).RPush), varargs...)
}
