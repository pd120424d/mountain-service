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
	CountAssignmentsByProfile(shiftID uint, profileType string) (int64, error)
	CreateAssignment(employeeID, shiftID uint, profileType string) (uint, error)
	GetShiftsByEmployeeID(employeeID uint, result *[]model.Shift) error
	GetShiftAvailability(date time.Time) (map[int]map[model.ProfileType]int, error)
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

func (r *shiftRepository) CountAssignmentsByProfile(shiftID uint, profileType string) (int64, error) {
	var count int64
	err := r.db.Model(&model.EmployeeShift{}).
		Where("shift_id = ? AND profile_type = ?", shiftID, profileType).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count assignments: %w", err)
	}
	return count, nil
}

func (r *shiftRepository) CreateAssignment(employeeID, shiftID uint, profileType string) (uint, error) {
	assignment := model.EmployeeShift{
		EmployeeID:  employeeID,
		ShiftID:     shiftID,
		ProfileType: profileType,
	}
	if err := r.db.Create(&assignment).Error; err != nil {
		return 0, fmt.Errorf("failed to create assignment: %w", err)
	}
	return assignment.ID, nil
}

func (r *shiftRepository) GetShiftsByEmployeeID(employeeID uint, result *[]model.Shift) error {
	return r.db.Table("employee_shifts").
		Select("employee_shifts.id, shifts.shift_date, shifts.shift_type, employee_shifts.profile_type").
		Joins("JOIN shifts ON employee_shifts.shift_id = shifts.id").
		Where("employee_shifts.employee_id = ?", employeeID).
		Order("shifts.shift_date ASC, shifts.shift_type ASC").
		Scan(result).Error
}

func (r *shiftRepository) GetShiftAvailability(date time.Time) (map[int]map[model.ProfileType]int, error) {
	// Initialize availability map
	availability := map[int]map[model.ProfileType]int{
		1: {model.Medic: 2, model.Technical: 4},
		2: {model.Medic: 2, model.Technical: 4},
		3: {model.Medic: 2, model.Technical: 4},
	}

	// Query current counts
	var counts []struct {
		ShiftType    int
		EmployeeRole string
		Count        int
	}
	err := r.db.Model(&model.Shift{}).
		Select("shift_type, profile_type, COUNT(*) AS count").
		Where("shift_date = ?", date.Format(time.DateOnly)).
		Group("shift_type, profile_type").
		Scan(&counts).Error
	if err != nil {
		return nil, err
	}

	// Deduct current counts from availability
	for _, count := range counts {
		role := model.ProfileTypeFromString(count.EmployeeRole)
		if _, ok := availability[count.ShiftType][role]; ok {
			availability[count.ShiftType][role] -= count.Count
		}
	}

	return availability, nil
}

func (r *shiftRepository) RemoveEmployeeFromShift(assignmentID uint) error {
	if err := r.db.Where("id = ?", assignmentID).Delete(&model.EmployeeShift{}).Error; err != nil {
		return fmt.Errorf("failed to remove employee from shift: %w", err)
	}
	return nil
}
