// Code generated by MockGen. DO NOT EDIT.
// Source: shift_repository.go
//
// Generated by this command:
//
//	mockgen -source=shift_repository.go -destination=shift_repository_gomock.go -package=repositories mountain_service/employee/internal/repositories -imports=gomock=go.uber.org/mock/gomock
//
// Package repositories is a generated GoMock package.
package repositories

import (
	model "github.com/pd120424d/mountain-service/api/employee/internal/model"
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockShiftRepository is a mock of ShiftRepository interface.
type MockShiftRepository struct {
	ctrl     *gomock.Controller
	recorder *MockShiftRepositoryMockRecorder
}

// MockShiftRepositoryMockRecorder is the mock recorder for MockShiftRepository.
type MockShiftRepositoryMockRecorder struct {
	mock *MockShiftRepository
}

// NewMockShiftRepository creates a new mock instance.
func NewMockShiftRepository(ctrl *gomock.Controller) *MockShiftRepository {
	mock := &MockShiftRepository{ctrl: ctrl}
	mock.recorder = &MockShiftRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShiftRepository) EXPECT() *MockShiftRepositoryMockRecorder {
	return m.recorder
}

// AssignedToShift mocks base method.
func (m *MockShiftRepository) AssignedToShift(employeeID, shiftID uint) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignedToShift", employeeID, shiftID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AssignedToShift indicates an expected call of AssignedToShift.
func (mr *MockShiftRepositoryMockRecorder) AssignedToShift(employeeID, shiftID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignedToShift", reflect.TypeOf((*MockShiftRepository)(nil).AssignedToShift), employeeID, shiftID)
}

// CountAssignmentsByProfile mocks base method.
func (m *MockShiftRepository) CountAssignmentsByProfile(shiftID uint, profileType model.ProfileType) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountAssignmentsByProfile", shiftID, profileType)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountAssignmentsByProfile indicates an expected call of CountAssignmentsByProfile.
func (mr *MockShiftRepositoryMockRecorder) CountAssignmentsByProfile(shiftID, profileType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountAssignmentsByProfile", reflect.TypeOf((*MockShiftRepository)(nil).CountAssignmentsByProfile), shiftID, profileType)
}

// CreateAssignment mocks base method.
func (m *MockShiftRepository) CreateAssignment(employeeID, shiftID uint) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAssignment", employeeID, shiftID)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAssignment indicates an expected call of CreateAssignment.
func (mr *MockShiftRepositoryMockRecorder) CreateAssignment(employeeID, shiftID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAssignment", reflect.TypeOf((*MockShiftRepository)(nil).CreateAssignment), employeeID, shiftID)
}

// GetOrCreateShift mocks base method.
func (m *MockShiftRepository) GetOrCreateShift(shiftDate time.Time, shiftType int) (*model.Shift, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrCreateShift", shiftDate, shiftType)
	ret0, _ := ret[0].(*model.Shift)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrCreateShift indicates an expected call of GetOrCreateShift.
func (mr *MockShiftRepositoryMockRecorder) GetOrCreateShift(shiftDate, shiftType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrCreateShift", reflect.TypeOf((*MockShiftRepository)(nil).GetOrCreateShift), shiftDate, shiftType)
}

// GetShiftAvailability mocks base method.
func (m *MockShiftRepository) GetShiftAvailability(start, end time.Time) (*model.ShiftsAvailabilityRange, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShiftAvailability", start, end)
	ret0, _ := ret[0].(*model.ShiftsAvailabilityRange)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShiftAvailability indicates an expected call of GetShiftAvailability.
func (mr *MockShiftRepositoryMockRecorder) GetShiftAvailability(start, end interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShiftAvailability", reflect.TypeOf((*MockShiftRepository)(nil).GetShiftAvailability), start, end)
}

// GetShiftsByEmployeeID mocks base method.
func (m *MockShiftRepository) GetShiftsByEmployeeID(employeeID uint, result *[]model.Shift) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShiftsByEmployeeID", employeeID, result)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetShiftsByEmployeeID indicates an expected call of GetShiftsByEmployeeID.
func (mr *MockShiftRepositoryMockRecorder) GetShiftsByEmployeeID(employeeID, result interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShiftsByEmployeeID", reflect.TypeOf((*MockShiftRepository)(nil).GetShiftsByEmployeeID), employeeID, result)
}

// RemoveEmployeeFromShift mocks base method.
func (m *MockShiftRepository) RemoveEmployeeFromShift(assignmentID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveEmployeeFromShift", assignmentID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveEmployeeFromShift indicates an expected call of RemoveEmployeeFromShift.
func (mr *MockShiftRepositoryMockRecorder) RemoveEmployeeFromShift(assignmentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveEmployeeFromShift", reflect.TypeOf((*MockShiftRepository)(nil).RemoveEmployeeFromShift), assignmentID)
}
