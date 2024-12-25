package repositories

import (
	"api/employee/internal/model"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ShiftRepository interface {
	AssignEmployee(shiftDate time.Time, shiftType int, employeeID uint, employeeRole string) error
	GetShiftsByEmployeeID(employeeID uint) ([]model.Shift, error)
	GetShiftAvailability(date time.Time) (map[int]map[string]int, error)
}

type shiftRepository struct {
	db *gorm.DB
}

func NewShiftRepository(db *gorm.DB) ShiftRepository {
	return &shiftRepository{db: db}
}

func (r *shiftRepository) AssignEmployee(shiftDate time.Time, shiftType int, employeeID uint, employeeRole string) error {
	// Validate shift type
	if shiftType < 1 || shiftType > 3 {
		return fmt.Errorf("invalid shift type: %d", shiftType)
	}

	// Validate role
	if employeeRole != "Medic" && employeeRole != "Technical" {
		return fmt.Errorf("invalid employee role: %s", employeeRole)
	}

	// Check if the employee is already assigned to this shift
	var existingShift model.Shift
	if err := r.db.Where("employee_id = ? AND shift_date = ? AND shift_type = ?", employeeID, shiftDate.Format("2006-01-02"), shiftType).First(&existingShift).Error; err == nil {
		return fmt.Errorf("employee already assigned to this shift on %s", shiftDate.Format("2006-01-02"))
	}

	// Count current employees in the shift
	var count int64
	err := r.db.Model(&model.Shift{}).
		Where("shift_date = ? AND shift_type = ? AND employee_role = ?", shiftDate.Format("2006-01-02"), shiftType, employeeRole).
		Count(&count).Error
	if err != nil {
		return err
	}

	// Enforce maximum limits
	if (employeeRole == "Medic" && count >= 2) || (employeeRole == "Technical" && count >= 4) {
		return fmt.Errorf("maximum limit reached for role %s in shift %d on %s", employeeRole, shiftType, shiftDate.Format("2006-01-02"))
	}

	// Assign the employee
	shift := &model.Shift{
		EmployeeID:   employeeID,
		ShiftDate:    shiftDate,
		ShiftType:    shiftType,
		EmployeeRole: employeeRole,
	}
	return r.db.Create(shift).Error
}

func (r *shiftRepository) GetShiftsByEmployeeID(employeeID uint) ([]model.Shift, error) {
	var shifts []model.Shift
	err := r.db.Where("employee_id = ?", employeeID).Order("shift_date ASC, shift_type ASC").Find(&shifts).Error
	return shifts, err
}

func (r *shiftRepository) GetShiftAvailability(date time.Time) (map[int]map[string]int, error) {
	// Initialize availability limits
	availability := map[int]map[string]int{
		1: {"Medic": 2, "Technical": 4},
		2: {"Medic": 2, "Technical": 4},
		3: {"Medic": 2, "Technical": 4},
	}

	// Query current counts
	var counts []struct {
		ShiftType    int
		EmployeeRole string
		Count        int
	}
	err := r.db.Model(&model.Shift{}).
		Select("shift_type, employee_role, COUNT(*) AS count").
		Where("shift_date = ?", date.Format("2006-01-02")).
		Group("shift_type, employee_role").
		Scan(&counts).Error
	if err != nil {
		return nil, err
	}

	// Deduct current counts from availability
	for _, count := range counts {
		if _, ok := availability[count.ShiftType][count.EmployeeRole]; ok {
			availability[count.ShiftType][count.EmployeeRole] -= count.Count
		}
	}

	return availability, nil
}
