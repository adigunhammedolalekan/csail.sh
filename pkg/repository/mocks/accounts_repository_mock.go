// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/saas/hostgolang/pkg/repository (interfaces: AccountRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	types "github.com/saas/hostgolang/pkg/types"
	reflect "reflect"
)

// MockAccountRepository is a mock of AccountRepository interface
type MockAccountRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAccountRepositoryMockRecorder
}

// MockAccountRepositoryMockRecorder is the mock recorder for MockAccountRepository
type MockAccountRepositoryMockRecorder struct {
	mock *MockAccountRepository
}

// NewMockAccountRepository creates a new mock instance
func NewMockAccountRepository(ctrl *gomock.Controller) *MockAccountRepository {
	mock := &MockAccountRepository{ctrl: ctrl}
	mock.recorder = &MockAccountRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAccountRepository) EXPECT() *MockAccountRepositoryMockRecorder {
	return m.recorder
}

// AuthenticateAccount mocks base method
func (m *MockAccountRepository) AuthenticateAccount(arg0 *types.AuthenticateAccountOpts) (*types.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthenticateAccount", arg0)
	ret0, _ := ret[0].(*types.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthenticateAccount indicates an expected call of AuthenticateAccount
func (mr *MockAccountRepositoryMockRecorder) AuthenticateAccount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthenticateAccount", reflect.TypeOf((*MockAccountRepository)(nil).AuthenticateAccount), arg0)
}

// CreateAccount mocks base method
func (m *MockAccountRepository) CreateAccount(arg0 *types.NewAccountOpts) (*types.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", arg0)
	ret0, _ := ret[0].(*types.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount
func (mr *MockAccountRepositoryMockRecorder) CreateAccount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockAccountRepository)(nil).CreateAccount), arg0)
}

// GetAccountByEmail mocks base method
func (m *MockAccountRepository) GetAccountByEmail(arg0 string) (*types.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountByEmail", arg0)
	ret0, _ := ret[0].(*types.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountByEmail indicates an expected call of GetAccountByEmail
func (mr *MockAccountRepositoryMockRecorder) GetAccountByEmail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountByEmail", reflect.TypeOf((*MockAccountRepository)(nil).GetAccountByEmail), arg0)
}
