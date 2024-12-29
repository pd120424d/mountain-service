package repositories

//go:generate mockgen -source=shift_repository.go -destination=shift_repository_gomock.go -package=repositories mountain_service/employee/internal/repositories -imports=gomock=go.uber.org/mock/gomock

import (
	"api/employee/internal/model"
	"api/shared/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ShiftRepository interface {
	AssignEmployee(shiftDate time.Time, shiftType int, employeeID uint, employeeRole string) error
	GetShiftsByEmployeeID(employeeID uint) ([]model.Shift, error)
	GetShiftAvailability(date time.Time) (map[int]map[model.ProfileType]int, error)
	RemoveEmployeeFromShift(shiftDate time.Time, shiftType int, employeeID uint) error
}

type shiftRepository struct {
	log utils.Logger
	db  *gorm.DB
}

func NewShiftRepository(log utils.Logger, db *gorm.DB) ShiftRepository {
	return &shiftRepository{log: log.WithName("shiftRepository"), db: db}
}

func (r *shiftRepository) AssignEmployee(shiftDate time.Time, shiftType int, employeeID uint, profileType string) error {
	// Check if the employee is already assigned to this shift
	var existingShift model.Shift
	if err := r.db.Where("employee_id = ? AND shift_date = ? AND shift_type = ?", employeeID, shiftDate.Format(time.DateOnly), shiftType).First(&existingShift).Error; err == nil {
		return model.ErrAlreadyAssigned
	}

	// Count current employees in the shift for the given profile type
	var count int64
	err := r.db.Model(&model.Shift{}).
		Where("shift_date = ? AND shift_type = ? AND employee_role = ?", shiftDate.Format(time.DateOnly), shiftType, profileType).
		Count(&count).Error
	if err != nil {
		return err
	}

	// Enforce maximum limits
	if (profileType == "Medic" && count >= 2) || (profileType == "Technical" && count >= 4) {
		return model.ErrCapacityReached
	}

	// Assign the employee
	shift := &model.Shift{
		EmployeeID:   employeeID,
		ShiftDate:    shiftDate,
		ShiftType:    shiftType,
		EmployeeRole: profileType,
	}
	return r.db.Create(shift).Error
}

func (r *shiftRepository) GetShiftsByEmployeeID(employeeID uint) ([]model.Shift, error) {
	var shifts []model.Shift
	err := r.db.Where("employee_id = ?", employeeID).Order("shift_date ASC, shift_type ASC").Find(&shifts).Error
	return shifts, err
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

func (r *shiftRepository) RemoveEmployeeFromShift(shiftDate time.Time, shiftType int, employeeID uint) error {
	err := r.db.Where("employee_id = ? AND shift_date = ? AND shift_type = ?", employeeID, shiftDate.Format("2006-01-02"), shiftType).
		Delete(&model.Shift{}).Error
	if err != nil {
		return fmt.Errorf("failed to remove employee from shift: %w", err)
	}
	return nil
}
