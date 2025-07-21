package repositories

//go:generate mockgen -source=shift_repository.go -destination=shift_repository_gomock.go -package=repositories mountain_service/employee/internal/repositories -imports=gomock=go.uber.org/mock/gomock

import (
	"errors"
	"fmt"
	"time"

	"github.com/pd120424d/mountain-service/api/employee/internal/model"
	"github.com/pd120424d/mountain-service/api/shared/utils"

	"gorm.io/gorm"
)

type ShiftRepository interface {
	GetOrCreateShift(shiftDate time.Time, shiftType int) (*model.Shift, error)
	AssignedToShift(employeeID, shiftID uint) (bool, error)
	CountAssignmentsByProfile(shiftID uint, profileType model.ProfileType) (int64, error)
	CreateAssignment(employeeID, shiftID uint) (uint, error)
	GetShiftsByEmployeeID(employeeID uint, result *[]model.Shift) error
	GetShiftAvailability(start, end time.Time) (*model.ShiftsAvailabilityRange, error)
	RemoveEmployeeFromShiftByDetails(employeeID uint, shiftDate time.Time, shiftType int) error
}

type shiftRepository struct {
	log utils.Logger
	db  *gorm.DB
}

func NewShiftRepository(log utils.Logger, db *gorm.DB) ShiftRepository {
	return &shiftRepository{log: log.WithName("shiftRepository"), db: db}
}

func (r *shiftRepository) GetOrCreateShift(shiftDate time.Time, shiftType int) (*model.Shift, error) {
	var shift model.Shift
	err := r.db.FirstOrCreate(&shift, model.Shift{
		ShiftDate: shiftDate,
		ShiftType: shiftType,
	}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find or create shift: %w", err)
	}
	return &shift, nil
}

func (r *shiftRepository) AssignedToShift(employeeID, shiftID uint) (bool, error) {
	var existing model.EmployeeShift
	err := r.db.Where("employee_id = ? AND shift_id = ?", employeeID, shiftID).First(&existing).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, fmt.Errorf("failed to check assignment: %w", err)
}

func (r *shiftRepository) CountAssignmentsByProfile(shiftID uint, profileType model.ProfileType) (int64, error) {
	var count int64
	err := r.db.Table("employee_shifts").
		Joins("JOIN employees ON employee_shifts.employee_id = employees.id").
		Where("employee_shifts.shift_id = ? AND employees.profile_type = ?", shiftID, profileType).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count assignments: %w", err)
	}
	return count, nil
}

func (r *shiftRepository) CreateAssignment(employeeID, shiftID uint) (uint, error) {
	assignment := model.EmployeeShift{
		EmployeeID: employeeID,
		ShiftID:    shiftID,
	}
	if err := r.db.Create(&assignment).Error; err != nil {
		return 0, fmt.Errorf("failed to create assignment: %w", err)
	}
	return assignment.ID, nil
}

func (r *shiftRepository) GetShiftsByEmployeeID(employeeID uint, result *[]model.Shift) error {
	return r.db.Table("employee_shifts").
		Select("employee_shifts.id, shifts.shift_date, shifts.shift_type").
		Joins("JOIN shifts ON employee_shifts.shift_id = shifts.id").
		Where("employee_shifts.employee_id = ?", employeeID).
		Order("shifts.shift_date ASC, shifts.shift_type ASC").
		Scan(result).Error
}

func (r *shiftRepository) GetShiftAvailability(start, end time.Time) (*model.ShiftsAvailabilityRange, error) {
	result := model.ShiftsAvailabilityRange{
		Days: map[time.Time][]map[model.ProfileType]int{},
	}

	// Initial availability per shift
	for d := start; d.Before(end); d = d.Add(24 * time.Hour) {
		day := d.Truncate(24 * time.Hour)
		result.Days[day] = []map[model.ProfileType]int{
			{model.Medic: 2, model.Technical: 4},
			{model.Medic: 2, model.Technical: 4},
			{model.Medic: 2, model.Technical: 4},
		}
	}

	// Query assigned employees grouped by shift and role
	var counts []struct {
		ShiftDate    time.Time
		ShiftType    int
		EmployeeRole string
		Count        int
	}

	err := r.db.Table("shifts").
		Joins("JOIN employee_shifts ON shifts.id = employee_shifts.shift_id").
		Joins("JOIN employees ON employee_shifts.employee_id = employees.id").
		Select("shifts.shift_date, shifts.shift_type, employees.profile_type AS employee_role, COUNT(*) AS count").
		Where("shift_date >= ? AND shift_date < ?", start, end).
		Group("shifts.shift_date, shifts.shift_type, employees.profile_type").
		Scan(&counts).Error
	if err != nil {
		return nil, err
	}

	// Deduct from initialized capacities
	for _, c := range counts {
		day := c.ShiftDate.Truncate(24 * time.Hour)
		role := model.ProfileTypeFromString(c.EmployeeRole)
		shiftIndex := c.ShiftType - 1

		if shifts, ok := result.Days[day]; ok && shiftIndex >= 0 && shiftIndex < len(shifts) {
			shifts[shiftIndex][role] -= c.Count
		}
	}

	return &result, nil
}

func (r *shiftRepository) RemoveEmployeeFromShiftByDetails(employeeID uint, shiftDate time.Time, shiftType int) error {
	var shift model.Shift
	err := r.db.Where("shift_date = ? AND shift_type = ?", shiftDate, shiftType).First(&shift).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("shift not found for date %s and type %d", shiftDate.Format(time.DateOnly), shiftType)
		}
		return fmt.Errorf("failed to find shift: %w", err)
	}

	var assignment model.EmployeeShift
	err = r.db.Where("employee_id = ? AND shift_id = ?", employeeID, shift.ID).First(&assignment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("employee is not assigned to this shift")
		}
		return fmt.Errorf("failed to find assignment: %w", err)
	}

	if err := r.db.Delete(&assignment).Error; err != nil {
		return fmt.Errorf("failed to remove employee from shift: %w", err)
	}

	return nil
}
