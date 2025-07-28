package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/employee/internal/model"
	"github.com/pd120424d/mountain-service/api/employee/internal/repositories"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

func TestEmployeeService_GetShifts(t *testing.T) {
	t.Parallel()

	t.Run("it fails when repository returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		shiftRepoMock.EXPECT().GetShiftsByEmployeeID(uint(1), gomock.Any()).Return(assert.AnError)

		response, err := service.GetShifts(1)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "failed to retrieve shifts", err.Error())
	})

	t.Run("it successfully returns shifts", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		shifts := []model.Shift{
			{
				ID:        1,
				ShiftDate: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				ShiftType: 1,
				CreatedAt: time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:        2,
				ShiftDate: time.Date(2025, 1, 16, 0, 0, 0, 0, time.UTC),
				ShiftType: 2,
				CreatedAt: time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
			},
		}

		shiftRepoMock.EXPECT().GetShiftsByEmployeeID(uint(1), gomock.Any()).DoAndReturn(
			func(employeeID uint, result *[]model.Shift) error {
				*result = shifts
				return nil
			})

		response, err := service.GetShifts(1)

		assert.NoError(t, err)
		assert.Len(t, response, 2)
		assert.Equal(t, uint(1), response[0].ID)
		assert.Equal(t, 1, response[0].ShiftType)
		assert.Equal(t, uint(2), response[1].ID)
		assert.Equal(t, 2, response[1].ShiftType)
	})
}

func TestEmployeeService_AssignShift(t *testing.T) {
	t.Parallel()

	t.Run("it fails with invalid employee ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()

		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		emplRepoMock.EXPECT().GetEmployeeByID(uint(0), gomock.Any()).Return(fmt.Errorf("service failure"))

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		response, err := service.AssignShift(0, employeeV1.AssignShiftRequest{})

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "employee not found", err.Error())
	})

	t.Run("it fails with invalid date format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		employee := &model.Employee{
			ID:          1,
			ProfileType: model.Medic,
		}

		req := employeeV1.AssignShiftRequest{
			ShiftDate: "invalid-date",
			ShiftType: 1,
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})

		response, err := service.AssignShift(1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "invalid shift date format", err.Error())
	})

	t.Run("it fails when shift is in the past", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		employee := &model.Employee{
			ID:          1,
			ProfileType: model.Medic,
		}

		req := employeeV1.AssignShiftRequest{
			ShiftDate: "2020-01-15", // Past date
			ShiftType: 1,
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})

		response, err := service.AssignShift(1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "cannot assign shift in the past", err.Error())
	})

	t.Run("it fails when shift is more than 3 months in advance", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		employee := &model.Employee{
			ID:          1,
			ProfileType: model.Medic,
		}

		futureDate := time.Now().AddDate(0, 4, 0).Format("2006-01-02") // 4 months in future

		req := employeeV1.AssignShiftRequest{
			ShiftDate: futureDate,
			ShiftType: 1,
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})

		response, err := service.AssignShift(1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "cannot assign shift more than 3 months in advance", err.Error())
	})

	t.Run("it fails when concesutive shift validation fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)
		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)
		employee := &model.Employee{
			ID:          1,
			ProfileType: model.Medic,
		}

		futureDate := time.Now().AddDate(0, 0, 7)
		futureDateStr := futureDate.Format("2006-01-02")

		req := employeeV1.AssignShiftRequest{
			ShiftDate: futureDateStr,
			ShiftType: 1,
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})

		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("failed to get shifts"))

		response, err := service.AssignShift(1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "failed to validate consecutive shifts", err.Error())
	})

	t.Run("it fails when would result in more than 2 consecutive shifts", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		employee := &model.Employee{
			ID:          1,
			ProfileType: model.Medic,
		}

		futureDate := time.Now().AddDate(0, 0, 7)
		futureDateStr := futureDate.Format("2006-01-02")

		req := employeeV1.AssignShiftRequest{
			ShiftDate: futureDateStr,
			ShiftType: 1,
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})

		// Mock existing shifts that would create 3 consecutive shifts (2 existing + 1 new)
		existingShifts := []model.Shift{
			{ShiftDate: futureDate.AddDate(0, 0, -2)}, // 2 days before
			{ShiftDate: futureDate.AddDate(0, 0, -1)}, // 1 day before
		}

		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(employeeID uint, startDate, endDate time.Time, result *[]model.Shift) error {
				*result = existingShifts
				return nil
			})

		response, err := service.AssignShift(1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "cannot assign shift: would result in more than 2 consecutive shifts", err.Error())
	})

	t.Run("it allows exactly 2 consecutive shifts", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		employee := &model.Employee{
			ID:          1,
			ProfileType: model.Medic,
		}

		futureDate := time.Now().AddDate(0, 0, 7)
		futureDateStr := futureDate.Format("2006-01-02")

		shift := &model.Shift{
			ID:        1,
			ShiftDate: futureDate,
			ShiftType: 1,
		}

		req := employeeV1.AssignShiftRequest{
			ShiftDate: futureDateStr,
			ShiftType: 1,
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})

		// Mock existing shift that would create exactly 2 consecutive shifts (1 existing + 1 new)
		existingShifts := []model.Shift{
			{ShiftDate: futureDate.AddDate(0, 0, -1)}, // 1 day before
		}

		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(employeeID uint, startDate, endDate time.Time, result *[]model.Shift) error {
				*result = existingShifts
				return nil
			})

		shiftRepoMock.EXPECT().GetOrCreateShift(gomock.Any(), 1).Return(shift, nil)
		shiftRepoMock.EXPECT().AssignedToShift(uint(1), uint(1)).Return(false, nil)
		shiftRepoMock.EXPECT().CountAssignmentsByProfile(uint(1), model.Medic).Return(int64(1), nil)
		shiftRepoMock.EXPECT().CreateAssignment(uint(1), uint(1)).Return(uint(10), nil)

		response, err := service.AssignShift(1, req)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, uint(10), response.ID)
		assert.Equal(t, futureDateStr, response.ShiftDate)
		assert.Equal(t, 1, response.ShiftType)
	})

	t.Run("it fails when GetOrCreateShift fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)
		employee := &model.Employee{
			ID:          1,
			ProfileType: model.Medic,
		}

		futureDate := time.Now().AddDate(0, 0, 7)
		futureDateStr := futureDate.Format("2006-01-02")

		req := employeeV1.AssignShiftRequest{
			ShiftDate: futureDateStr,
			ShiftType: 1,
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})

		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		shiftRepoMock.EXPECT().GetOrCreateShift(gomock.Any(), 1).Return(nil, fmt.Errorf("failed to get or create shift"))

		response, err := service.AssignShift(1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "failed to process shift", err.Error())
	})

	t.Run("it fails when AssignedToShift call fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)
		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)
		employee := &model.Employee{
			ID:          1,
			ProfileType: model.Medic,
		}
		futureDate := time.Now().AddDate(0, 0, 7)
		futureDateStr := futureDate.Format("2006-01-02")
		shift := &model.Shift{
			ID:        1,
			ShiftDate: futureDate,
			ShiftType: 1,
		}

		req := employeeV1.AssignShiftRequest{
			ShiftDate: futureDateStr,
			ShiftType: 1,
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})

		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		shiftRepoMock.EXPECT().GetOrCreateShift(gomock.Any(), 1).Return(shift, nil)
		shiftRepoMock.EXPECT().AssignedToShift(uint(1), uint(1)).Return(false, fmt.Errorf("failed to check assignment"))

		response, err := service.AssignShift(1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "failed to check assignment", err.Error())
	})

	t.Run("it fails when CountAssignmentsByProfile call fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)
		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)
		employee := &model.Employee{
			ID:          1,
			ProfileType: model.Medic,
		}

		futureDate := time.Now().AddDate(0, 0, 7)
		futureDateStr := futureDate.Format("2006-01-02")
		shift := &model.Shift{
			ID:        1,
			ShiftDate: futureDate,
			ShiftType: 1,
		}

		req := employeeV1.AssignShiftRequest{
			ShiftDate: futureDateStr,
			ShiftType: 1,
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})
		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		shiftRepoMock.EXPECT().GetOrCreateShift(gomock.Any(), 1).Return(shift, nil)
		shiftRepoMock.EXPECT().AssignedToShift(uint(1), uint(1)).Return(false, nil)
		shiftRepoMock.EXPECT().CountAssignmentsByProfile(uint(1), model.Medic).Return(int64(0), fmt.Errorf("failed to count assignments"))

		response, err := service.AssignShift(1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "failed to check shift capacity", err.Error())
	})

	t.Run("it fails when employee is already assigned to shift", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		employee := &model.Employee{
			ID:          1,
			ProfileType: model.Medic,
		}
		futureDate := time.Now().AddDate(0, 0, 7)
		futureDateStr := futureDate.Format("2006-01-02")

		shift := &model.Shift{
			ID:        1,
			ShiftDate: futureDate,
			ShiftType: 1,
		}

		req := employeeV1.AssignShiftRequest{
			ShiftDate: futureDateStr,
			ShiftType: 1,
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})

		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		shiftRepoMock.EXPECT().GetOrCreateShift(gomock.Any(), 1).Return(shift, nil)
		shiftRepoMock.EXPECT().AssignedToShift(uint(1), uint(1)).Return(true, nil) // Already assigned

		response, err := service.AssignShift(1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "employee is already assigned to this shift", err.Error())
	})

	t.Run("it fails when shift capacity is full for medics", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		employee := &model.Employee{
			ID:          1,
			ProfileType: model.Medic,
		}
		futureDate := time.Now().AddDate(0, 0, 7)
		futureDateStr := futureDate.Format("2006-01-02")

		shift := &model.Shift{
			ID:        1,
			ShiftDate: futureDate,
			ShiftType: 1,
		}

		req := employeeV1.AssignShiftRequest{
			ShiftDate: futureDateStr,
			ShiftType: 1,
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})

		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		shiftRepoMock.EXPECT().GetOrCreateShift(gomock.Any(), 1).Return(shift, nil)
		shiftRepoMock.EXPECT().AssignedToShift(uint(1), uint(1)).Return(false, nil)
		shiftRepoMock.EXPECT().CountAssignmentsByProfile(uint(1), model.Medic).Return(int64(2), nil) // Full capacity

		response, err := service.AssignShift(1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "maximum capacity for this role reached in the selected shift", err.Error())
	})

	t.Run("it fails when CreateAssignment call fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)
		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		employee := &model.Employee{
			ID:          1,
			ProfileType: model.Medic,
		}
		futureDate := time.Now().AddDate(0, 0, 7)
		futureDateStr := futureDate.Format("2006-01-02")
		shift := &model.Shift{
			ID:        1,
			ShiftDate: futureDate,
			ShiftType: 1,
		}

		req := employeeV1.AssignShiftRequest{
			ShiftDate: futureDateStr,
			ShiftType: 1,
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})

		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		shiftRepoMock.EXPECT().GetOrCreateShift(gomock.Any(), 1).Return(shift, nil)
		shiftRepoMock.EXPECT().AssignedToShift(uint(1), uint(1)).Return(false, nil)
		shiftRepoMock.EXPECT().CountAssignmentsByProfile(uint(1), model.Medic).Return(int64(1), nil)
		shiftRepoMock.EXPECT().CreateAssignment(uint(1), uint(1)).Return(uint(0), fmt.Errorf("failed to create assignment"))

		response, err := service.AssignShift(1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "failed to assign employee", err.Error())
	})

	t.Run("it successfully assigns shift", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		employee := &model.Employee{
			ID:          1,
			ProfileType: model.Medic,
		}
		futureDate := time.Now().AddDate(0, 0, 7) // 7 days from now
		futureDateStr := futureDate.Format("2006-01-02")

		shift := &model.Shift{
			ID:        1,
			ShiftDate: futureDate,
			ShiftType: 1,
		}

		req := employeeV1.AssignShiftRequest{
			ShiftDate: futureDateStr,
			ShiftType: 1,
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})

		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		shiftRepoMock.EXPECT().GetOrCreateShift(gomock.Any(), 1).Return(shift, nil)
		shiftRepoMock.EXPECT().AssignedToShift(uint(1), uint(1)).Return(false, nil)
		shiftRepoMock.EXPECT().CountAssignmentsByProfile(uint(1), model.Medic).Return(int64(1), nil)
		shiftRepoMock.EXPECT().CreateAssignment(uint(1), uint(1)).Return(uint(10), nil)

		response, err := service.AssignShift(1, req)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, uint(10), response.ID)
		assert.Equal(t, futureDateStr, response.ShiftDate)
		assert.Equal(t, 1, response.ShiftType)
	})
}

func TestEmployeeService_GetShiftsAvailability(t *testing.T) {
	t.Parallel()

	t.Run("it fails when days parameter is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		response, err := service.GetShiftsAvailability(-1)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "invalid days parameter", err.Error())
	})

	t.Run("it fails when GetShiftAvailability call fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)
		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)
		shiftRepoMock.EXPECT().GetShiftAvailability(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("failed to get shift availability"))

		response, err := service.GetShiftsAvailability(7)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "failed to retrieve shift availability", err.Error())
	})

	t.Run("it sucessfully calculates shift availability", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)
		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)
		availability := &model.ShiftsAvailabilityRange{
			Days: map[time.Time][]map[model.ProfileType]int{
				time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC): {
					{model.Medic: 2, model.Technical: 4},
					{model.Medic: 2, model.Technical: 4},
					{model.Medic: 2, model.Technical: 4},
				},
			},
		}
		shiftRepoMock.EXPECT().GetShiftAvailability(gomock.Any(), gomock.Any()).Return(availability, nil)

		expectedResponse := &employeeV1.ShiftAvailabilityResponse{
			Days: map[time.Time]employeeV1.ShiftAvailabilityPerDay{
				time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC): {
					Shift1: employeeV1.ShiftAvailability{
						MedicSlotsAvailable:     2, // Full availability: 2 - 0 assigned
						TechnicalSlotsAvailable: 4, // Full availability: 4 - 0 assigned
					},
					Shift2: employeeV1.ShiftAvailability{
						MedicSlotsAvailable:     2,
						TechnicalSlotsAvailable: 4,
					},
					Shift3: employeeV1.ShiftAvailability{
						MedicSlotsAvailable:     2,
						TechnicalSlotsAvailable: 4,
					},
				},
			},
		}

		response, err := service.GetShiftsAvailability(7)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, len(response.Days))
		assert.Equal(t, expectedResponse, response)
	})

}

func TestEmployeeService_RemoveShift(t *testing.T) {
	t.Parallel()

	t.Run("it fails when shift date format is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		req := employeeV1.RemoveShiftRequest{
			ShiftDate: "invalid-date",
			ShiftType: 1,
		}

		err := service.RemoveShift(1, req)

		assert.Error(t, err)
		assert.Equal(t, "invalid shift date format", err.Error())
	})

	t.Run("it fails when RemoveEmployeeFromShiftByDetails call fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)
		shiftRepoMock.EXPECT().RemoveEmployeeFromShiftByDetails(uint(1), gomock.Any(), gomock.Any()).Return(fmt.Errorf("failed to remove shift"))

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)
		req := employeeV1.RemoveShiftRequest{
			ShiftDate: "2025-02-03",
			ShiftType: 1,
		}

		err := service.RemoveShift(1, req)

		assert.Error(t, err)
		assert.Equal(t, "failed to remove shift", err.Error())
	})

	t.Run("it successfully removes shift assignment", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)
		shiftRepoMock.EXPECT().RemoveEmployeeFromShiftByDetails(uint(1), gomock.Any(), gomock.Any()).Return(nil)
		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)
		req := employeeV1.RemoveShiftRequest{
			ShiftDate: "2025-02-03",
			ShiftType: 1,
		}

		err := service.RemoveShift(1, req)

		assert.NoError(t, err)
	})

}

func TestEmployeeService_GetOnCallEmployees(t *testing.T) {
	t.Parallel()

	t.Run("it fails when GetOnCallEmployees call fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		shiftRepoMock.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("failed to get on-call employees"))

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		response, err := service.GetOnCallEmployees(time.Now(), 0)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "failed to retrieve on-call employees", err.Error())
	})

	t.Run("it successfully returns on-call employees", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		employees := []model.Employee{
			{
				ID:          1,
				Username:    "johndoe",
				FirstName:   "John",
				LastName:    "Doe",
				Email:       "john@example.com",
				ProfileType: model.Medic,
			},
			{
				ID:          2,
				Username:    "janesmith",
				FirstName:   "Jane",
				LastName:    "Smith",
				Email:       "jane@example.com",
				ProfileType: model.Technical,
			},
		}
		shiftRepoMock.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return(employees, nil)

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		response, err := service.GetOnCallEmployees(time.Now(), 0)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response))
	})

}

func TestEmployeeService_GetShiftWarnings(t *testing.T) {
	t.Parallel()

	t.Run("it fails when GetEmployeeByID call fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).Return(fmt.Errorf("failed to get employee"))

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		response, err := service.GetShiftWarnings(1)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "employee not found", err.Error())
	})

	t.Run("it fails when GetShiftsByEmployeeIDInDateRange call fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = model.Employee{ID: 1, ProfileType: model.Medic}
			return nil
		}).Times(1)

		// Mock GetShiftAvailability to return uncovered shifts so the flow continues
		shiftRepoMock.EXPECT().GetShiftAvailability(gomock.Any(), gomock.Any()).Return(&model.ShiftsAvailabilityRange{
			Days: map[time.Time][]map[model.ProfileType]int{
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC): {
					{model.Medic: 0, model.Technical: 4}, // Uncovered shift - needs medic
				},
			},
		}, nil).Times(1)

		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("failed to get shifts in date range")).Times(1)

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		response, err := service.GetShiftWarnings(1)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "failed to get shifts in date range", err.Error())
	})

	t.Run("it returns empty list when all shifts are adequately covered", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = model.Employee{ID: 1, ProfileType: model.Medic}
			return nil
		}).Times(1)

		// Mock GetShiftAvailability to return fully covered shifts (no uncovered shifts)
		// When all shifts have adequate coverage, findUncoveredShifts returns empty slice
		shiftRepoMock.EXPECT().GetShiftAvailability(gomock.Any(), gomock.Any()).Return(&model.ShiftsAvailabilityRange{
			Days: map[time.Time][]map[model.ProfileType]int{
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC): {
					{model.Medic: 2, model.Technical: 4}, // Shift 1 - fully covered
					{model.Medic: 2, model.Technical: 4}, // Shift 2 - fully covered
					{model.Medic: 2, model.Technical: 4}, // Shift 3 - fully covered
				},
			},
		}, nil).Times(1)

		// Since shifts are fully covered, GetShiftsByEmployeeIDInDateRange should NOT be called

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		response, err := service.GetShiftWarnings(1)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 0, len(response))
	})

	t.Run("it returns empty list when there are uncovered shifts but employee has met quota", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = model.Employee{ID: 1, ProfileType: model.Medic}
			return nil
		}).Times(1)

		// Mock GetShiftAvailability to return some uncovered shifts
		shiftRepoMock.EXPECT().GetShiftAvailability(gomock.Any(), gomock.Any()).Return(&model.ShiftsAvailabilityRange{
			Days: map[time.Time][]map[model.ProfileType]int{
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC): {
					{model.Medic: 0, model.Technical: 4}, // Shift 1 - needs medic coverage
					{model.Medic: 2, model.Technical: 4}, // Shift 2 - fully covered
					{model.Medic: 2, model.Technical: 4}, // Shift 3 - fully covered
				},
			},
		}, nil).Times(1)

		// Mock employee shifts showing they've met their quota (10 shifts = 5 per week for 2 weeks)
		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(employeeID uint, startDate, endDate time.Time, result *[]model.Shift) error {
				now := time.Now()
				shifts := make([]model.Shift, 10)
				// First week: 5 shifts (days 0, 1, 2, 3, 4)
				for i := 0; i < 5; i++ {
					shifts[i] = model.Shift{
						ID:        uint(i + 1),
						ShiftDate: now.AddDate(0, 0, i),
						ShiftType: 1,
					}
				}
				// Second week: 5 shifts (days 7, 8, 9, 10, 11)
				for i := 5; i < 10; i++ {
					shifts[i] = model.Shift{
						ID:        uint(i + 1),
						ShiftDate: now.AddDate(0, 0, i+2), // +2 to skip weekend and start second week
						ShiftType: 1,
					}
				}
				*result = shifts
				return nil
			}).Times(1)

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		response, err := service.GetShiftWarnings(1)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 0, len(response))
	})

	t.Run("it returns warnings when there are uncovered shifts and employee hasn't met quota", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = model.Employee{ID: 1, ProfileType: model.Medic}
			return nil
		}).Times(1)

		// Mock GetShiftAvailability to return some uncovered shifts
		shiftRepoMock.EXPECT().GetShiftAvailability(gomock.Any(), gomock.Any()).Return(&model.ShiftsAvailabilityRange{
			Days: map[time.Time][]map[model.ProfileType]int{
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC): {
					{model.Medic: 0, model.Technical: 4}, // Shift 1 - needs medic coverage
					{model.Medic: 0, model.Technical: 2}, // Shift 2 - needs medic coverage
					{model.Medic: 2, model.Technical: 4}, // Shift 3 - fully covered
				},
			},
		}, nil).Times(1)

		// Mock employee shifts showing they haven't met their quota (only 3 shifts)
		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(employeeID uint, startDate, endDate time.Time, result *[]model.Shift) error {
				now := time.Now()
				shifts := make([]model.Shift, 3)
				for i := 0; i < 3; i++ {
					shifts[i] = model.Shift{
						ID:        uint(i + 1),
						ShiftDate: now.AddDate(0, 0, i),
						ShiftType: 1,
					}
				}
				*result = shifts
				return nil
			}).Times(1)

		service := NewEmployeeService(log, emplRepoMock, shiftRepoMock)

		response, err := service.GetShiftWarnings(1)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, len(response))
		assert.Contains(t, response[0], "2 shifts in the next 2 weeks that need Medic coverage")
	})

}

func TestEmployeeService_getMaxCapacityForProfile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		profileType model.ProfileType
		expected    int
	}{
		{model.Medic, 2},
		{model.Technical, 4},
		{model.Administrator, 0},
	}
	for _, test := range tests {

		t.Run(fmt.Sprintf("it returns correct capacity for %s", test.profileType), func(t *testing.T) {
			log := utils.NewTestLogger()
			service := &employeeService{log: log}

			capacity := service.getMaxCapacityForProfile(test.profileType)

			assert.Equal(t, test.expected, capacity)
		})
	}
}

func TestEmployeeService_validateConsecutiveShifts(t *testing.T) {
	t.Parallel()

	t.Run("it fails when GetShiftsByEmployeeIDInDateRange call fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)
		service := &employeeService{log: log, shiftsRepo: shiftRepoMock}

		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("failed to get shifts"))

		err := service.validateConsecutiveShifts(1, time.Now())

		assert.Error(t, err)
		assert.Equal(t, "failed to validate consecutive shifts", err.Error())
	})

	t.Run("it fails when employee would have more than 2 consecutive shifts", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)
		service := &employeeService{log: log, shiftsRepo: shiftRepoMock}
		newShiftDate := time.Now().AddDate(0, 0, 7)
		existingShifts := []model.Shift{
			{ShiftDate: newShiftDate.AddDate(0, 0, -2)}, // 2 days before
			{ShiftDate: newShiftDate.AddDate(0, 0, -1)}, // 1 day before
		}
		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(employeeID uint, startDate, endDate time.Time, result *[]model.Shift) error {
				*result = existingShifts
				return nil
			})

		err := service.validateConsecutiveShifts(1, newShiftDate)

		assert.Error(t, err)
		assert.Equal(t, "cannot assign shift: would result in more than 2 consecutive shifts", err.Error())
	})
}

func TestEmployeeService_countConsecutiveShifts(t *testing.T) {
	t.Parallel()

	t.Run("it counts single shift", func(t *testing.T) {
		log := utils.NewTestLogger()
		service := &employeeService{log: log}

		shiftDates := map[string]bool{
			"2025-01-15": true,
		}
		centerDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)

		count := service.countConsecutiveShifts(shiftDates, centerDate)

		assert.Equal(t, 1, count)
	})

	t.Run("it counts consecutive shifts", func(t *testing.T) {
		log := utils.NewTestLogger()
		service := &employeeService{log: log}

		shiftDates := map[string]bool{
			"2025-01-13": true,
			"2025-01-14": true,
			"2025-01-15": true,
			"2025-01-16": true,
			"2025-01-17": true,
		}
		centerDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)

		count := service.countConsecutiveShifts(shiftDates, centerDate)

		assert.Equal(t, 5, count)
	})

	t.Run("it counts consecutive shifts with gaps", func(t *testing.T) {
		log := utils.NewTestLogger()
		service := &employeeService{log: log}

		shiftDates := map[string]bool{
			"2025-01-13": true,
			"2025-01-14": true,
			"2025-01-15": true,
			// Gap on 2025-01-16
			"2025-01-17": true,
		}
		centerDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)

		count := service.countConsecutiveShifts(shiftDates, centerDate)

		assert.Equal(t, 3, count) // Only counts up to the gap
	})
}

func TestEmployeeService_checkWeeklyQuota(t *testing.T) {
	t.Parallel()

	t.Run("it returns false when employee has not met weekly quota", func(t *testing.T) {
		log := utils.NewTestLogger()
		service := &employeeService{log: log}

		shifts := []model.Shift{
			{ShiftDate: time.Now()},
			{ShiftDate: time.Now().AddDate(0, 0, 1)},
		}
		startDate := time.Now()

		metQuota := service.checkWeeklyQuota(shifts, startDate)

		assert.False(t, metQuota)
	})

	t.Run("it returns true when employee has met weekly quota", func(t *testing.T) {
		log := utils.NewTestLogger()
		service := &employeeService{log: log}

		shifts := []model.Shift{
			{ShiftDate: time.Now()},
			{ShiftDate: time.Now().AddDate(0, 0, 1)},
			{ShiftDate: time.Now().AddDate(0, 0, 2)},
			{ShiftDate: time.Now().AddDate(0, 0, 3)},
			{ShiftDate: time.Now().AddDate(0, 0, 4)},
			{ShiftDate: time.Now().AddDate(0, 0, 5)},
			{ShiftDate: time.Now().AddDate(0, 0, 8)},
			{ShiftDate: time.Now().AddDate(0, 0, 9)},
			{ShiftDate: time.Now().AddDate(0, 0, 10)},
			{ShiftDate: time.Now().AddDate(0, 0, 11)},
			{ShiftDate: time.Now().AddDate(0, 0, 12)},
		}
		startDate := time.Now()

		metQuota := service.checkWeeklyQuota(shifts, startDate)

		assert.True(t, metQuota)
	})
}

func TestEmployeeService_findUncoveredShifts(t *testing.T) {
	t.Parallel()

	t.Run("it fails when GetShiftAvailability call fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)
		service := &employeeService{log: log, shiftsRepo: shiftRepoMock}
		shiftRepoMock.EXPECT().GetShiftAvailability(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("failed to get shift availability"))

		_, err := service.findUncoveredShifts(model.Medic, time.Now(), time.Now().AddDate(0, 0, 7))

		assert.Error(t, err)
		assert.Equal(t, "failed to get shift availability: failed to get shift availability", err.Error())
	})

	t.Run("it correctly returns uncovered shifts for medic profile", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)
		service := &employeeService{log: log, shiftsRepo: shiftRepoMock}
		availability := &model.ShiftsAvailabilityRange{
			Days: map[time.Time][]map[model.ProfileType]int{
				time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC): {
					{model.Medic: 0, model.Technical: 4}, // Shift 1 - no medic coverage
					{model.Medic: 0, model.Technical: 4}, // Shift 2 - no medic coverage
					{model.Medic: 0, model.Technical: 4}, // Shift 3 - no medic coverage
				},
			},
		}
		shiftRepoMock.EXPECT().GetShiftAvailability(gomock.Any(), gomock.Any()).Return(availability, nil)

		uncoveredShifts, err := service.findUncoveredShifts(model.Medic, time.Now(), time.Now().AddDate(0, 0, 7))

		assert.NoError(t, err)
		assert.NotNil(t, uncoveredShifts)
		assert.Equal(t, 3, len(uncoveredShifts))
	})
}
