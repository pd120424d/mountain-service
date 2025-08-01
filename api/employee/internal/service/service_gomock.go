// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package service is a generated GoMock package.
package service

import (
	reflect "reflect"
	time "time"

	v1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	model "github.com/pd120424d/mountain-service/api/employee/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockShiftService is a mock of ShiftService interface.
type MockShiftService struct {
	ctrl     *gomock.Controller
	recorder *MockShiftServiceMockRecorder
}

// MockShiftServiceMockRecorder is the mock recorder for MockShiftService.
type MockShiftServiceMockRecorder struct {
	mock *MockShiftService
}

// NewMockShiftService creates a new mock instance.
func NewMockShiftService(ctrl *gomock.Controller) *MockShiftService {
	mock := &MockShiftService{ctrl: ctrl}
	mock.recorder = &MockShiftServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShiftService) EXPECT() *MockShiftServiceMockRecorder {
	return m.recorder
}

// AssignShift mocks base method.
func (m *MockShiftService) AssignShift(employeeID uint, req v1.AssignShiftRequest) (*v1.AssignShiftResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignShift", employeeID, req)
	ret0, _ := ret[0].(*v1.AssignShiftResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AssignShift indicates an expected call of AssignShift.
func (mr *MockShiftServiceMockRecorder) AssignShift(employeeID, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignShift", reflect.TypeOf((*MockShiftService)(nil).AssignShift), employeeID, req)
}

// GetOnCallEmployees mocks base method.
func (m *MockShiftService) GetOnCallEmployees(currentTime time.Time, shiftBuffer time.Duration) ([]v1.EmployeeResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOnCallEmployees", currentTime, shiftBuffer)
	ret0, _ := ret[0].([]v1.EmployeeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOnCallEmployees indicates an expected call of GetOnCallEmployees.
func (mr *MockShiftServiceMockRecorder) GetOnCallEmployees(currentTime, shiftBuffer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOnCallEmployees", reflect.TypeOf((*MockShiftService)(nil).GetOnCallEmployees), currentTime, shiftBuffer)
}

// GetShiftWarnings mocks base method.
func (m *MockShiftService) GetShiftWarnings(employeeID uint) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShiftWarnings", employeeID)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShiftWarnings indicates an expected call of GetShiftWarnings.
func (mr *MockShiftServiceMockRecorder) GetShiftWarnings(employeeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShiftWarnings", reflect.TypeOf((*MockShiftService)(nil).GetShiftWarnings), employeeID)
}

// GetShifts mocks base method.
func (m *MockShiftService) GetShifts(employeeID uint) ([]v1.ShiftResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShifts", employeeID)
	ret0, _ := ret[0].([]v1.ShiftResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShifts indicates an expected call of GetShifts.
func (mr *MockShiftServiceMockRecorder) GetShifts(employeeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShifts", reflect.TypeOf((*MockShiftService)(nil).GetShifts), employeeID)
}

// GetShiftsAvailability mocks base method.
func (m *MockShiftService) GetShiftsAvailability(employeeID uint, days int) (*v1.ShiftAvailabilityResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShiftsAvailability", employeeID, days)
	ret0, _ := ret[0].(*v1.ShiftAvailabilityResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShiftsAvailability indicates an expected call of GetShiftsAvailability.
func (mr *MockShiftServiceMockRecorder) GetShiftsAvailability(employeeID, days interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShiftsAvailability", reflect.TypeOf((*MockShiftService)(nil).GetShiftsAvailability), employeeID, days)
}

// RemoveShift mocks base method.
func (m *MockShiftService) RemoveShift(employeeID uint, req v1.RemoveShiftRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveShift", employeeID, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveShift indicates an expected call of RemoveShift.
func (mr *MockShiftServiceMockRecorder) RemoveShift(employeeID, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveShift", reflect.TypeOf((*MockShiftService)(nil).RemoveShift), employeeID, req)
}

// MockEmployeeService is a mock of EmployeeService interface.
type MockEmployeeService struct {
	ctrl     *gomock.Controller
	recorder *MockEmployeeServiceMockRecorder
}

// MockEmployeeServiceMockRecorder is the mock recorder for MockEmployeeService.
type MockEmployeeServiceMockRecorder struct {
	mock *MockEmployeeService
}

// NewMockEmployeeService creates a new mock instance.
func NewMockEmployeeService(ctrl *gomock.Controller) *MockEmployeeService {
	mock := &MockEmployeeService{ctrl: ctrl}
	mock.recorder = &MockEmployeeServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmployeeService) EXPECT() *MockEmployeeServiceMockRecorder {
	return m.recorder
}

// DeleteEmployee mocks base method.
func (m *MockEmployeeService) DeleteEmployee(employeeID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEmployee", employeeID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEmployee indicates an expected call of DeleteEmployee.
func (mr *MockEmployeeServiceMockRecorder) DeleteEmployee(employeeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEmployee", reflect.TypeOf((*MockEmployeeService)(nil).DeleteEmployee), employeeID)
}

// GetEmployeeByID mocks base method.
func (m *MockEmployeeService) GetEmployeeByID(employeeID uint) (*model.Employee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmployeeByID", employeeID)
	ret0, _ := ret[0].(*model.Employee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEmployeeByID indicates an expected call of GetEmployeeByID.
func (mr *MockEmployeeServiceMockRecorder) GetEmployeeByID(employeeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmployeeByID", reflect.TypeOf((*MockEmployeeService)(nil).GetEmployeeByID), employeeID)
}

// GetEmployeeByUsername mocks base method.
func (m *MockEmployeeService) GetEmployeeByUsername(username string) (*model.Employee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmployeeByUsername", username)
	ret0, _ := ret[0].(*model.Employee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEmployeeByUsername indicates an expected call of GetEmployeeByUsername.
func (mr *MockEmployeeServiceMockRecorder) GetEmployeeByUsername(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmployeeByUsername", reflect.TypeOf((*MockEmployeeService)(nil).GetEmployeeByUsername), username)
}

// ListEmployees mocks base method.
func (m *MockEmployeeService) ListEmployees() ([]v1.EmployeeResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEmployees")
	ret0, _ := ret[0].([]v1.EmployeeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEmployees indicates an expected call of ListEmployees.
func (mr *MockEmployeeServiceMockRecorder) ListEmployees() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEmployees", reflect.TypeOf((*MockEmployeeService)(nil).ListEmployees))
}

// LoginEmployee mocks base method.
func (m *MockEmployeeService) LoginEmployee(req v1.EmployeeLogin) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginEmployee", req)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginEmployee indicates an expected call of LoginEmployee.
func (mr *MockEmployeeServiceMockRecorder) LoginEmployee(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginEmployee", reflect.TypeOf((*MockEmployeeService)(nil).LoginEmployee), req)
}

// RegisterEmployee mocks base method.
func (m *MockEmployeeService) RegisterEmployee(req v1.EmployeeCreateRequest) (*v1.EmployeeResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterEmployee", req)
	ret0, _ := ret[0].(*v1.EmployeeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterEmployee indicates an expected call of RegisterEmployee.
func (mr *MockEmployeeServiceMockRecorder) RegisterEmployee(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterEmployee", reflect.TypeOf((*MockEmployeeService)(nil).RegisterEmployee), req)
}

// ResetAllData mocks base method.
func (m *MockEmployeeService) ResetAllData() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetAllData")
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetAllData indicates an expected call of ResetAllData.
func (mr *MockEmployeeServiceMockRecorder) ResetAllData() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetAllData", reflect.TypeOf((*MockEmployeeService)(nil).ResetAllData))
}

// UpdateEmployee mocks base method.
func (m *MockEmployeeService) UpdateEmployee(employeeID uint, req v1.EmployeeUpdateRequest) (*v1.EmployeeResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEmployee", employeeID, req)
	ret0, _ := ret[0].(*v1.EmployeeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateEmployee indicates an expected call of UpdateEmployee.
func (mr *MockEmployeeServiceMockRecorder) UpdateEmployee(employeeID, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEmployee", reflect.TypeOf((*MockEmployeeService)(nil).UpdateEmployee), employeeID, req)
}
