package repositories

//go:generate mockgen -source=shift_repository.go -destination=shift_repository_gomock.go -package=repositories mountain_service/employee/internal/repositories -imports=gomock=go.uber.org/mock/gomock

import (
	"api/employee/internal/model"
	"api/shared/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ShiftRepository interface {
	AssignEmployee(shiftDate time.Time, shiftType int, employeeID uint, employeeRole string) (uint, error)
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

func (r *shiftRepository) AssignEmployee(shiftDate time.Time, shiftType int, employeeID uint, profileType string) (uint, error) {
	// Step 1: Ensure the shift exists or create it
	var shift model.Shift
	err := r.db.Where("shift_date = ? AND shift_type = ?", shiftDate.Format("2006-01-02"), shiftType).
		FirstOrCreate(&shift, model.Shift{
			ShiftDate: shiftDate,
			ShiftType: shiftType,
		}).Error
	if err != nil {
		return 0, fmt.Errorf("failed to find or create shift: %w", err)
	}

	// Step 2: Check if the employee is already assigned to this shift
	var existingAssignment model.EmployeeShift
	if err := r.db.Where("employee_id = ? AND shift_id = ?", employeeID, shift.ID).First(&existingAssignment).Error; err == nil {
		return 0, model.ErrAlreadyAssigned
	}

	// Step 3: Count current assignments for the given role in this shift
	var count int64
	err = r.db.Model(&model.EmployeeShift{}).
		Where("shift_id = ? AND profile_type = ?", shift.ID, profileType).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count assignments: %w", err)
	}

	// Step 4: Enforce role-based limits
	if (profileType == "Medic" && count >= 2) || (profileType == "Technical" && count >= 4) {
		return 0, model.ErrCapacityReached
	}

	// Step 5: Assign the employee to the shift
	assignment := model.EmployeeShift{
		EmployeeID:  employeeID,
		ShiftID:     shift.ID,
		ProfileType: profileType,
	}
	if err := r.db.Create(&assignment).Error; err != nil {
		return 0, fmt.Errorf("failed to assign employee to shift: %w", err)
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
