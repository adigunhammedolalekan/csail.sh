// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/saas/hostgolang/pkg/services (interfaces: K8sService)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	types "github.com/saas/hostgolang/pkg/types"
	reflect "reflect"
)

// MockK8sService is a mock of K8sService interface
type MockK8sService struct {
	ctrl     *gomock.Controller
	recorder *MockK8sServiceMockRecorder
}

// MockK8sServiceMockRecorder is the mock recorder for MockK8sService
type MockK8sServiceMockRecorder struct {
	mock *MockK8sService
}

// NewMockK8sService creates a new mock instance
func NewMockK8sService(ctrl *gomock.Controller) *MockK8sService {
	mock := &MockK8sService{ctrl: ctrl}
	mock.recorder = &MockK8sServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockK8sService) EXPECT() *MockK8sServiceMockRecorder {
	return m.recorder
}

// AddDomain mocks base method
func (m *MockK8sService) AddDomain(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDomain", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddDomain indicates an expected call of AddDomain
func (mr *MockK8sServiceMockRecorder) AddDomain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDomain", reflect.TypeOf((*MockK8sService)(nil).AddDomain), arg0, arg1)
}

// DeployService mocks base method
func (m *MockK8sService) DeployService(arg0 *types.CreateDeploymentOpts) (*types.DeploymentResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeployService", arg0)
	ret0, _ := ret[0].(*types.DeploymentResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeployService indicates an expected call of DeployService
func (mr *MockK8sServiceMockRecorder) DeployService(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeployService", reflect.TypeOf((*MockK8sService)(nil).DeployService), arg0)
}

// GetLogs mocks base method
func (m *MockK8sService) GetLogs(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogs", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLogs indicates an expected call of GetLogs
func (mr *MockK8sServiceMockRecorder) GetLogs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogs", reflect.TypeOf((*MockK8sService)(nil).GetLogs), arg0)
}

// ListRunningPods mocks base method
func (m *MockK8sService) ListRunningPods(arg0 string) ([]types.Instance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRunningPods", arg0)
	ret0, _ := ret[0].([]types.Instance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRunningPods indicates an expected call of ListRunningPods
func (mr *MockK8sServiceMockRecorder) ListRunningPods(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRunningPods", reflect.TypeOf((*MockK8sService)(nil).ListRunningPods), arg0)
}

// PodExec mocks base method
func (m *MockK8sService) PodExec(arg0, arg1 string, arg2 []string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PodExec", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PodExec indicates an expected call of PodExec
func (mr *MockK8sServiceMockRecorder) PodExec(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PodExec", reflect.TypeOf((*MockK8sService)(nil).PodExec), arg0, arg1, arg2)
}

// RemoveDomain mocks base method
func (m *MockK8sService) RemoveDomain(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveDomain", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveDomain indicates an expected call of RemoveDomain
func (mr *MockK8sServiceMockRecorder) RemoveDomain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveDomain", reflect.TypeOf((*MockK8sService)(nil).RemoveDomain), arg0, arg1)
}

// ScaleApp mocks base method
func (m *MockK8sService) ScaleApp(arg0 string, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScaleApp", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ScaleApp indicates an expected call of ScaleApp
func (mr *MockK8sServiceMockRecorder) ScaleApp(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScaleApp", reflect.TypeOf((*MockK8sService)(nil).ScaleApp), arg0, arg1)
}

// UpdateEnvs mocks base method
func (m *MockK8sService) UpdateEnvs(arg0 string, arg1 []types.Environment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEnvs", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEnvs indicates an expected call of UpdateEnvs
func (mr *MockK8sServiceMockRecorder) UpdateEnvs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEnvs", reflect.TypeOf((*MockK8sService)(nil).UpdateEnvs), arg0, arg1)
}
