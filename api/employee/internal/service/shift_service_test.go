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

func TestShiftService_GetShifts(t *testing.T) {
	t.Parallel()

	t.Run("it fails when repository returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		shiftRepoMock.EXPECT().GetEmployeeShiftRowsByEmployeeID(uint(1)).Return(nil, assert.AnError)

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

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		rows := []repositories.EmployeeShiftRow{
			{
				AssignmentID:   111,
				AssignedAt:     time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
				ShiftID:        1,
				ShiftDate:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				ShiftType:      1,
				ShiftCreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				AssignmentID:   222,
				AssignedAt:     time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC),
				ShiftID:        2,
				ShiftDate:      time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
				ShiftType:      2,
				ShiftCreatedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		}

		shiftRepoMock.EXPECT().GetEmployeeShiftRowsByEmployeeID(uint(1)).Return(rows, nil)

		response, err := service.GetShifts(1)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response))
		assert.Equal(t, uint(1), response[0].ID)
		assert.Equal(t, 1, response[0].ShiftType)
		assert.Equal(t, rows[0].AssignedAt, response[0].AssignedAt)
		assert.Equal(t, uint(2), response[1].ID)
		assert.Equal(t, 2, response[1].ShiftType)
		assert.Equal(t, rows[1].AssignedAt, response[1].AssignedAt)
	})
}

func TestShiftService_AssignShift(t *testing.T) {
	t.Parallel()

	t.Run("it fails with invalid employee ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		futureDate := time.Now().AddDate(0, 0, 7)
		futureDateStr := futureDate.Format("2006-01-02")

		req := employeeV1.AssignShiftRequest{
			ShiftDate: futureDateStr,
			ShiftType: 1,
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).Return(assert.AnError)

		response, err := service.AssignShift(1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "employee not found", err.Error())
	})

	t.Run("it fails with invalid shift date format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

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

	t.Run("it fails when shift date is in the past", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		employee := &model.Employee{
			ID:          1,
			ProfileType: model.Medic,
		}

		pastDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC) // Clearly in the past
		pastDateStr := pastDate.Format("2006-01-02")

		req := employeeV1.AssignShiftRequest{
			ShiftDate: pastDateStr,
			ShiftType: 1,
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})

		response, err := service.AssignShift(1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "shift date must be in the future", err.Error())
	})

	t.Run("it fails when shift date is more than 3 months in the future", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		employee := &model.Employee{
			ID:          1,
			ProfileType: model.Medic,
		}

		futureDate := time.Now().AddDate(0, 4, 0) // 4 months in the future
		futureDateStr := futureDate.Format("2006-01-02")

		req := employeeV1.AssignShiftRequest{
			ShiftDate: futureDateStr,
			ShiftType: 1,
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})

		response, err := service.AssignShift(1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "shift date cannot be more than 3 months in the future", err.Error())
	})

	t.Run("it fails when it would result in 3 consecutive shift slots (or insufficient rest)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		employee := &model.Employee{
			ID:          1,
			ProfileType: model.Medic,
		}

		futureDate := time.Now().AddDate(0, 0, 7).Truncate(24 * time.Hour)
		futureDateStr := futureDate.Format("2006-01-02")

		// Request the 3rd shift on the same day, when 1st and 2nd are already assigned => 3 consecutive slots
		req := employeeV1.AssignShiftRequest{
			ShiftDate: futureDateStr,
			ShiftType: 3,
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})

		existingShifts := []model.Shift{
			{ShiftDate: futureDate, ShiftType: 1},
			{ShiftDate: futureDate, ShiftType: 2},
		}

		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(employeeID uint, startDate, endDate time.Time, result *[]model.Shift) error {
			*result = existingShifts
			return nil
		})

		response, err := service.AssignShift(1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "SHIFT_ERRORS.CONSECUTIVE_SHIFTS_LIMIT|3", err.Error())
	})

	t.Run("calendar rest day blocks all shifts on next day after 1+2", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		employee := &model.Employee{ID: 1, ProfileType: model.Medic}
		D := time.Date(2025, 9, 15, 0, 0, 0, 0, time.UTC)
		D1 := D.AddDate(0, 0, 1)

		existingShifts := []model.Shift{{ShiftDate: D, ShiftType: 1}, {ShiftDate: D, ShiftType: 2}}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		}).Times(3)
		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(employeeID uint, startDate, endDate time.Time, result *[]model.Shift) error {
			*result = existingShifts
			return nil
		}).Times(3)

		for _, st := range []int{1, 2, 3} {
			req := employeeV1.AssignShiftRequest{ShiftDate: D1.Format("2006-01-02"), ShiftType: st}
			resp, err := service.AssignShift(1, req)
			assert.Error(t, err)
			assert.Nil(t, resp)
			assert.Equal(t, "SHIFT_ERRORS.CONSECUTIVE_SHIFTS_LIMIT|3", err.Error())
		}
	})

	t.Run("after 1+2, first allowed is 1st shift on D+2", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		employee := &model.Employee{ID: 1, ProfileType: model.Medic}
		D := time.Date(2025, 9, 15, 0, 0, 0, 0, time.UTC)
		D2 := D.AddDate(0, 0, 2)

		existingShifts := []model.Shift{{ShiftDate: D, ShiftType: 1}, {ShiftDate: D, ShiftType: 2}}
		shift := &model.Shift{ID: 10, ShiftDate: D2, ShiftType: 1}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})
		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(employeeID uint, startDate, endDate time.Time, result *[]model.Shift) error {
			*result = existingShifts
			return nil
		})

		req := employeeV1.AssignShiftRequest{ShiftDate: D2.Format("2006-01-02"), ShiftType: 1}
		shiftRepoMock.EXPECT().GetOrCreateShift(gomock.Any(), 1).Return(shift, nil)
		shiftRepoMock.EXPECT().AssignedToShift(uint(1), uint(10)).Return(false, nil)
		shiftRepoMock.EXPECT().CountAssignmentsByProfile(uint(10), model.Medic).Return(int64(0), nil)
		shiftRepoMock.EXPECT().CreateAssignment(uint(1), uint(10)).Return(uint(77), nil)

		resp, err := service.AssignShift(1, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, uint(77), resp.ID)
		assert.Equal(t, 1, resp.ShiftType)
		assert.Equal(t, D2.Format("2006-01-02"), resp.ShiftDate)
	})

	t.Run("2+3 cross-midnight: rest day is D+1; D+2 is allowed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		employee := &model.Employee{ID: 1, ProfileType: model.Medic}
		D := time.Date(2025, 9, 15, 0, 0, 0, 0, time.UTC)
		D2 := D.AddDate(0, 0, 2)

		existingShifts := []model.Shift{{ShiftDate: D, ShiftType: 2}, {ShiftDate: D, ShiftType: 3}}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})
		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(employeeID uint, startDate, endDate time.Time, result *[]model.Shift) error {
			*result = existingShifts
			return nil
		})

		req := employeeV1.AssignShiftRequest{ShiftDate: D2.Format("2006-01-02"), ShiftType: 1}
		shift := &model.Shift{ID: 42, ShiftDate: D2, ShiftType: 1}
		shiftRepoMock.EXPECT().GetOrCreateShift(gomock.Any(), 1).Return(shift, nil)
		shiftRepoMock.EXPECT().AssignedToShift(uint(1), uint(42)).Return(false, nil)
		shiftRepoMock.EXPECT().CountAssignmentsByProfile(uint(42), model.Medic).Return(int64(0), nil)
		shiftRepoMock.EXPECT().CreateAssignment(uint(1), uint(42)).Return(uint(99), nil)

		resp, err := service.AssignShift(1, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, uint(99), resp.ID)
		assert.Equal(t, 1, resp.ShiftType)
		assert.Equal(t, D2.Format("2006-01-02"), resp.ShiftDate)
	})

	t.Run("2+3 cross-midnight: first allowed is D+3 shift 1", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		employee := &model.Employee{ID: 1, ProfileType: model.Medic}
		D := time.Date(2025, 9, 15, 0, 0, 0, 0, time.UTC)
		D3 := D.AddDate(0, 0, 3)

		existingShifts := []model.Shift{{ShiftDate: D, ShiftType: 2}, {ShiftDate: D, ShiftType: 3}}
		shift := &model.Shift{ID: 21, ShiftDate: D3, ShiftType: 1}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})
		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(employeeID uint, startDate, endDate time.Time, result *[]model.Shift) error {
			*result = existingShifts
			return nil
		})

		req := employeeV1.AssignShiftRequest{ShiftDate: D3.Format("2006-01-02"), ShiftType: 1}
		shiftRepoMock.EXPECT().GetOrCreateShift(gomock.Any(), 1).Return(shift, nil)
		shiftRepoMock.EXPECT().AssignedToShift(uint(1), uint(21)).Return(false, nil)
		shiftRepoMock.EXPECT().CountAssignmentsByProfile(uint(21), model.Medic).Return(int64(0), nil)
		shiftRepoMock.EXPECT().CreateAssignment(uint(1), uint(21)).Return(uint(88), nil)

		resp, err := service.AssignShift(1, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, uint(88), resp.ID)
		assert.Equal(t, 1, resp.ShiftType)
		assert.Equal(t, D3.Format("2006-01-02"), resp.ShiftDate)
	})

	t.Run("it fails when create repository call fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		employee := &model.Employee{
			ID:          1,
			ProfileType: model.Medic,
		}

		futureDate := time.Now().UTC().Truncate(24*time.Hour).AddDate(0, 0, 7)
		futureDateStr := futureDate.Format("2006-01-02")

		req := employeeV1.AssignShiftRequest{
			ShiftDate: futureDateStr,
			ShiftType: 1,
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})

		existingShifts := []model.Shift{
			// Keep within rule: not exceeding 2 consecutive days including requested date
			{ShiftDate: futureDate.AddDate(0, 0, -2)},
			{ShiftDate: futureDate.AddDate(0, 0, 2)},
		}

		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(employeeID uint, startDate, endDate time.Time, result *[]model.Shift) error {
			*result = existingShifts
			return nil
		})

		shiftRepoMock.EXPECT().GetOrCreateShift(gomock.Any(), 1).Return(nil, fmt.Errorf("database error")).Times(1)

		response, err := service.AssignShift(1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "failed to create shift", err.Error())
	})

	t.Run("it fails when check assignment repository call fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

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

		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(employeeID uint, startDate, endDate time.Time, result *[]model.Shift) error {
			*result = []model.Shift{}
			return nil
		})

		shiftRepoMock.EXPECT().GetOrCreateShift(gomock.Any(), 1).Return(shift, nil)
		shiftRepoMock.EXPECT().AssignedToShift(uint(1), uint(1)).Return(false, fmt.Errorf("database error")).Times(1)

		response, err := service.AssignShift(1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "failed to check assignment", err.Error())
	})

	t.Run("it fails when employee is already assigned to shift", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

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

		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(employeeID uint, startDate, endDate time.Time, result *[]model.Shift) error {
			*result = []model.Shift{}
			return nil
		})

		shiftRepoMock.EXPECT().GetOrCreateShift(gomock.Any(), 1).Return(shift, nil)
		shiftRepoMock.EXPECT().AssignedToShift(uint(1), uint(1)).Return(true, nil)

		response, err := service.AssignShift(1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "employee is already assigned to this shift", err.Error())
	})

	t.Run("it fails when count assignments repository call fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

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

		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(employeeID uint, startDate, endDate time.Time, result *[]model.Shift) error {
			*result = []model.Shift{}
			return nil
		})

		shiftRepoMock.EXPECT().GetOrCreateShift(gomock.Any(), 1).Return(shift, nil)
		shiftRepoMock.EXPECT().AssignedToShift(uint(1), uint(1)).Return(false, nil)
		shiftRepoMock.EXPECT().CountAssignmentsByProfile(uint(1), model.Medic).Return(int64(1), fmt.Errorf("database error")).Times(1)

		response, err := service.AssignShift(1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("it fails when shift capacity is full for medics", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

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

		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(employeeID uint, startDate, endDate time.Time, result *[]model.Shift) error {
			*result = []model.Shift{}
			return nil
		})

		shiftRepoMock.EXPECT().GetOrCreateShift(gomock.Any(), 1).Return(shift, nil)
		shiftRepoMock.EXPECT().AssignedToShift(uint(1), uint(1)).Return(false, nil)
		shiftRepoMock.EXPECT().CountAssignmentsByProfile(uint(1), model.Medic).Return(int64(2), nil) // Full capacity

		response, err := service.AssignShift(1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "shift capacity is full for Medic staff", err.Error())
	})

	t.Run("it successfully assigns shift", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

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

		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(employeeID uint, startDate, endDate time.Time, result *[]model.Shift) error {
			*result = []model.Shift{}
			return nil
		})

		shiftRepoMock.EXPECT().GetOrCreateShift(gomock.Any(), 1).Return(shift, nil)
		shiftRepoMock.EXPECT().AssignedToShift(uint(1), uint(1)).Return(false, nil)
		shiftRepoMock.EXPECT().CountAssignmentsByProfile(uint(1), model.Medic).Return(int64(1), nil) // Available capacity
		shiftRepoMock.EXPECT().CreateAssignment(uint(1), uint(1)).Return(uint(10), nil)

		response, err := service.AssignShift(1, req)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, uint(10), response.ID)
		assert.Equal(t, futureDateStr, response.ShiftDate)
		assert.Equal(t, 1, response.ShiftType)
	})
}

func TestShiftService_GetShiftsAvailability(t *testing.T) {
	t.Parallel()

	t.Run("it fails when days parameter is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		emplId := uint(1)

		response, err := service.GetShiftsAvailability(emplId, 91)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "days must be between 1 and 90", err.Error())
	})

	t.Run("it fails when repository returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		emplId := uint(1)

		shiftRepoMock.EXPECT().GetShiftAvailabilityWithEmployeeStatus(emplId, gomock.Any(), gomock.Any()).Return(nil, assert.AnError)

		response, err := service.GetShiftsAvailability(emplId, 7)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "failed to retrieve shift availability", err.Error())
	})

	t.Run("it successfully returns shift availability", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		today := time.Now().Truncate(24 * time.Hour)
		tomorrow := today.AddDate(0, 0, 1)

		shiftRepoMock.EXPECT().GetShiftAvailabilityWithEmployeeStatus(uint(1), gomock.Any(), gomock.Any()).Return(&model.ShiftsAvailabilityWithEmployeeStatus{
			Days: map[time.Time][]model.ShiftAvailabilityWithStatus{
				today: {
					{MedicSlotsAvailable: 1, TechnicalSlotsAvailable: 2, IsAssignedToEmployee: false, IsFullyBooked: false},
					{MedicSlotsAvailable: 2, TechnicalSlotsAvailable: 4, IsAssignedToEmployee: true, IsFullyBooked: false},
					{MedicSlotsAvailable: 0, TechnicalSlotsAvailable: 1, IsAssignedToEmployee: false, IsFullyBooked: false},
				},
				tomorrow: {
					{MedicSlotsAvailable: 2, TechnicalSlotsAvailable: 4, IsAssignedToEmployee: false, IsFullyBooked: false},
					{MedicSlotsAvailable: 1, TechnicalSlotsAvailable: 3, IsAssignedToEmployee: false, IsFullyBooked: false},
					{MedicSlotsAvailable: 2, TechnicalSlotsAvailable: 4, IsAssignedToEmployee: false, IsFullyBooked: false},
				},
			},
		}, nil)

		response, err := service.GetShiftsAvailability(1, 2)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Days))

		assert.Contains(t, response.Days, today)
		assert.Contains(t, response.Days, tomorrow)

		todayShifts := response.Days[today]
		assert.Equal(t, 1, todayShifts.FirstShift.MedicSlotsAvailable)
		assert.Equal(t, 2, todayShifts.FirstShift.TechnicalSlotsAvailable)
		assert.Equal(t, false, todayShifts.FirstShift.IsAssignedToEmployee)
		assert.Equal(t, false, todayShifts.FirstShift.IsFullyBooked)

		assert.Equal(t, 2, todayShifts.SecondShift.MedicSlotsAvailable)
		assert.Equal(t, 4, todayShifts.SecondShift.TechnicalSlotsAvailable)
		assert.Equal(t, true, todayShifts.SecondShift.IsAssignedToEmployee)
		assert.Equal(t, false, todayShifts.SecondShift.IsFullyBooked)

		assert.Equal(t, 0, todayShifts.ThirdShift.MedicSlotsAvailable)
		assert.Equal(t, 1, todayShifts.ThirdShift.TechnicalSlotsAvailable)
		assert.Equal(t, false, todayShifts.ThirdShift.IsAssignedToEmployee)
		assert.Equal(t, false, todayShifts.ThirdShift.IsFullyBooked)
	})
}

func TestShiftService_RemoveShift(t *testing.T) {
	t.Parallel()

	t.Run("it fails when shift date format is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		req := employeeV1.RemoveShiftRequest{
			ShiftDate: "invalid-date",
			ShiftType: 1,
		}

		err := service.RemoveShift(1, req)

		assert.Error(t, err)
		assert.Equal(t, "invalid shift date format", err.Error())
	})

	t.Run("it fails when repository returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		futureDate := time.Now().AddDate(0, 0, 7) // 7 days from now
		futureDateStr := futureDate.Format("2006-01-02")

		req := employeeV1.RemoveShiftRequest{
			ShiftDate: futureDateStr,
			ShiftType: 1,
		}

		shiftRepoMock.EXPECT().RemoveEmployeeFromShiftByDetails(uint(1), gomock.Any(), 1).Return(assert.AnError)

		err := service.RemoveShift(1, req)

		assert.Error(t, err)
		assert.Equal(t, "failed to remove shift", err.Error())
	})

	t.Run("it successfully removes shift", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		futureDate := time.Now().AddDate(0, 0, 7) // 7 days from now
		futureDateStr := futureDate.Format("2006-01-02")

		req := employeeV1.RemoveShiftRequest{
			ShiftDate: futureDateStr,
			ShiftType: 1,
		}

		shiftRepoMock.EXPECT().RemoveEmployeeFromShiftByDetails(uint(1), gomock.Any(), 1).Return(nil)

		err := service.RemoveShift(1, req)

		assert.NoError(t, err)
	})
}

func TestShiftService_GetOnCallEmployees(t *testing.T) {
	t.Parallel()

	t.Run("it fails when repository returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		currentTime := time.Now()
		shiftBuffer := time.Hour

		shiftRepoMock.EXPECT().GetOnCallEmployees(currentTime, shiftBuffer).Return(nil, assert.AnError)

		response, err := service.GetOnCallEmployees(currentTime, shiftBuffer)

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

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		currentTime := time.Now()
		shiftBuffer := time.Hour

		employees := []model.Employee{
			{
				ID:          1,
				Username:    "medic1",
				FirstName:   "Marko",
				LastName:    "Markovic",
				ProfileType: model.Medic,
			},
			{
				ID:          2,
				Username:    "tech1",
				FirstName:   "Marko",
				LastName:    "Markovic",
				ProfileType: model.Technical,
			},
		}

		shiftRepoMock.EXPECT().GetOnCallEmployees(currentTime, shiftBuffer).Return(employees, nil)

		response, err := service.GetOnCallEmployees(currentTime, shiftBuffer)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response))
		assert.Equal(t, uint(1), response[0].ID)
		assert.Equal(t, "medic1", response[0].Username)
		assert.Equal(t, uint(2), response[1].ID)
		assert.Equal(t, "tech1", response[1].Username)
	})
}

func TestShiftService_GetShiftWarnings(t *testing.T) {
	t.Parallel()

	t.Run("it fails when employee not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).Return(assert.AnError)

		warnings, err := service.GetShiftWarnings(1)

		assert.Error(t, err)
		assert.Nil(t, warnings)
		assert.Equal(t, "employee not found", err.Error())
	})

	t.Run("it returns no warnings when coverage is adequate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		employee := &model.Employee{
			ID:          1,
			ProfileType: model.Medic,
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})

		// Mock GetShiftAvailability to return shifts with at least one Medic assigned (no zero-coverage for Medic)
		shiftRepoMock.EXPECT().GetShiftAvailability(gomock.Any(), gomock.Any()).Return(&model.ShiftsAvailabilityRange{
			Days: map[time.Time][]map[model.ProfileType]int{
				time.Now().Truncate(24 * time.Hour): {
					{model.Medic: 1, model.Technical: 4}, // at least one medic assigned
					{model.Medic: 1, model.Technical: 4},
					{model.Medic: 1, model.Technical: 4},
				},
			},
		}, nil)

		warnings, err := service.GetShiftWarnings(1)

		assert.NoError(t, err)
		assert.Equal(t, 0, len(warnings))
	})

	t.Run("it returns warnings when coverage is inadequate and employee hasn't met quota", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)
		shiftRepoMock := repositories.NewMockShiftRepository(ctrl)

		service := NewShiftService(log, emplRepoMock, shiftRepoMock)

		employee := &model.Employee{
			ID:          1,
			ProfileType: model.Medic,
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = *employee
			return nil
		})

		// Mock GetShiftAvailability to include at least one shift with zero Medics assigned (2 slots available)
		shiftRepoMock.EXPECT().GetShiftAvailability(gomock.Any(), gomock.Any()).Return(&model.ShiftsAvailabilityRange{
			Days: map[time.Time][]map[model.ProfileType]int{
				time.Now().Truncate(24 * time.Hour): {
					{model.Medic: 2, model.Technical: 4}, // zero medics assigned triggers warning flow
					{model.Medic: 1, model.Technical: 4},
					{model.Medic: 1, model.Technical: 4},
				},
			},
		}, nil)

		// Mock employee shifts - less than 5 per week in next 2 weeks, dates within the 14-day window
		shiftRepoMock.EXPECT().GetShiftsByEmployeeIDInDateRange(uint(1), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(employeeID uint, startDate, endDate time.Time, result *[]model.Shift) error {
			t1 := time.Now().Truncate(24*time.Hour).AddDate(0, 0, 1)
			t2 := time.Now().Truncate(24*time.Hour).AddDate(0, 0, 2)
			*result = []model.Shift{
				{ID: 1, ShiftDate: t1},
				{ID: 2, ShiftDate: t2},
			}
			return nil
		})

		warnings, err := service.GetShiftWarnings(1)

		assert.NoError(t, err)
		assert.NotNil(t, warnings)
		assert.Equal(t, 1, len(warnings))
		assert.Equal(t, "SHIFT_WARNINGS.INSUFFICIENT_SHIFTS|2|14|5", warnings[0])
	})
}
