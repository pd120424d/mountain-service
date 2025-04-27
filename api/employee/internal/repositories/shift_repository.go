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
	GetShiftAvailability(date time.Time) (*model.ShiftsAvailability, error)
	RemoveEmployeeFromShift(assignmentID uint) error
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

func (r *shiftRepository) GetShiftAvailability(date time.Time) (*model.ShiftsAvailability, error) {
	// Initial availability per shift
	availability := model.ShiftsAvailability{
		Availability: map[int]map[model.ProfileType]int{
			1: {model.Medic: 2, model.Technical: 4},
			2: {model.Medic: 2, model.Technical: 4},
			3: {model.Medic: 2, model.Technical: 4},
		},
	}

	// Query assigned employees grouped by shift and role
	var counts []struct {
		ShiftType    int
		EmployeeRole string
		Count        int
	}

	start := date.Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)

	err := r.db.Table("shifts").
		Joins("JOIN employee_shifts ON shifts.id = employee_shifts.shift_id").
		Joins("JOIN employees ON employee_shifts.employee_id = employees.id").
		Select("shifts.shift_type, employees.profile_type AS employee_role, COUNT(*) AS count").
		Where("shift_date >= ? AND shift_date < ?", start, end).
		Group("shifts.shift_type, employees.profile_type").
		Scan(&counts).Error
	if err != nil {
		return nil, err
	}

	// Deduct current counts from max availability
	for _, c := range counts {
		role := model.ProfileTypeFromString(c.EmployeeRole)
		if _, ok := availability.Availability[c.ShiftType][role]; ok {
			availability.Availability[c.ShiftType][role] -= c.Count
		}
	}

	return &availability, nil
}

func (r *shiftRepository) RemoveEmployeeFromShift(assignmentID uint) error {
	err := r.db.First(&model.EmployeeShift{}, assignmentID).Error
	if err != nil {
		return fmt.Errorf("failed to find assignment: %w", err)
	}
	if err := r.db.Delete(&model.EmployeeShift{}, assignmentID).Error; err != nil {
		return fmt.Errorf("failed to remove employee from shift: %w", err)
	}
	return nil
}
