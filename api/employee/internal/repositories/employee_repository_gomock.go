// Code generated by MockGen. DO NOT EDIT.
// Source: employee_repository.go

// Package repositories is a generated GoMock package.
package repositories

import (
	model "api/employee/internal/model"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockEmployeeRepository is a mock of EmployeeRepository interface.
type MockEmployeeRepository struct {
	ctrl     *gomock.Controller
	recorder *MockEmployeeRepositoryMockRecorder
}

// MockEmployeeRepositoryMockRecorder is the mock recorder for MockEmployeeRepository.
type MockEmployeeRepositoryMockRecorder struct {
	mock *MockEmployeeRepository
}

// NewMockEmployeeRepository creates a new mock instance.
func NewMockEmployeeRepository(ctrl *gomock.Controller) *MockEmployeeRepository {
	mock := &MockEmployeeRepository{ctrl: ctrl}
	mock.recorder = &MockEmployeeRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmployeeRepository) EXPECT() *MockEmployeeRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockEmployeeRepository) Create(employee *model.Employee) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", employee)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockEmployeeRepositoryMockRecorder) Create(employee interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockEmployeeRepository)(nil).Create), employee)
}

// Delete mocks base method.
func (m *MockEmployeeRepository) Delete(employeeID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", employeeID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockEmployeeRepositoryMockRecorder) Delete(employeeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockEmployeeRepository)(nil).Delete), employeeID)
}

// GetAll mocks base method.
func (m *MockEmployeeRepository) GetAll() ([]model.Employee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]model.Employee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockEmployeeRepositoryMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockEmployeeRepository)(nil).GetAll))
}

// GetEmployeeByID mocks base method.
func (m *MockEmployeeRepository) GetEmployeeByID(id string, employee *model.Employee) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmployeeByID", id, employee)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetEmployeeByID indicates an expected call of GetEmployeeByID.
func (mr *MockEmployeeRepositoryMockRecorder) GetEmployeeByID(id, employee interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmployeeByID", reflect.TypeOf((*MockEmployeeRepository)(nil).GetEmployeeByID), id, employee)
}

// ListEmployees mocks base method.
func (m *MockEmployeeRepository) ListEmployees(filters map[string]interface{}) ([]model.Employee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEmployees", filters)
	ret0, _ := ret[0].([]model.Employee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEmployees indicates an expected call of ListEmployees.
func (mr *MockEmployeeRepositoryMockRecorder) ListEmployees(filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEmployees", reflect.TypeOf((*MockEmployeeRepository)(nil).ListEmployees), filters)
}

// UpdateEmployee mocks base method.
func (m *MockEmployeeRepository) UpdateEmployee(employee *model.Employee) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEmployee", employee)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEmployee indicates an expected call of UpdateEmployee.
func (mr *MockEmployeeRepositoryMockRecorder) UpdateEmployee(employee interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEmployee", reflect.TypeOf((*MockEmployeeRepository)(nil).UpdateEmployee), employee)
}