package repositories

//go:generate mockgen -source=employee_repository.go -destination=employee_repository_gomock.go -package=repositories mountain_service/employee/internal/repositories -imports=gomock=go.uber.org/mock/gomock

import (
	"fmt"
	"gorm.io/gorm"
	"time"

	"api/employee/internal/model"
	"api/shared/utils"
)

const dateFormat = "2006-01-02"

type EmployeeRepository interface {
	Create(employee *model.Employee) error
	GetAll() ([]model.Employee, error)
	GetEmployeeByID(id string, employee *model.Employee) error
	UpdateEmployee(employee *model.Employee) error
	Delete(employeeID uint) error
	ListEmployees(filters map[string]interface{}) ([]model.Employee, error)
}

type employeeRepository struct {
	log utils.Logger
	db  *gorm.DB
}

func NewEmployeeRepository(log utils.Logger, db *gorm.DB) EmployeeRepository {
	return &employeeRepository{log: log, db: db}
}

// Create creates and employee with the hashed version of its password.
func (r *employeeRepository) Create(employee *model.Employee) error {
	hashedPassword, err := utils.HashPassword(employee.Password)
	if err != nil {
		return err
	}
	employee.Password = hashedPassword
	return r.db.Create(employee).Error
}

// GetAll returns all the employees which have deleted_at flag set as NULL in db.
func (r *employeeRepository) GetAll() ([]model.Employee, error) {
	var employees []model.Employee
	err := r.db.Where("deleted_at IS NULL").Find(&employees).Error
	return employees, err
}

// GetEmployeeByID returns employee by its id or error if it cannot be found.
func (r *employeeRepository) GetEmployeeByID(id string, employee *model.Employee) error {
	return r.db.First(employee, "id = ?", id).Error
}

func (r *employeeRepository) ListEmployees(filters map[string]interface{}) ([]model.Employee, error) {
	var employees []model.Employee
	query := r.db.Model(&model.Employee{})
	for key, value := range filters {
		query = query.Where(key+" LIKE ?", "%"+value.(string)+"%")
	}
	r.log.Infof("query: %v", query)
	err := query.Find(&employees).Error
	return employees, err
}

func (r *employeeRepository) UpdateEmployee(employee *model.Employee) error {
	return r.db.Save(employee).Error
}

// Delete marks the employee record as deleted by setting the deleted_at timestamp.
func (r *employeeRepository) Delete(id uint) error {
	// Start by fetching the employee to ensure it exists
	var employee model.Employee
	if err := r.db.First(&employee, id).Error; err != nil {
		return err
	}

	// Permanently delete the employee
	if err := r.db.Unscoped().Delete(&employee).Error; err != nil {
		return err
	}

	return nil
}

func (r *employeeRepository) AssignShift(employeeID uint, shiftDate time.Time, shiftType int, employeeRole string) error {
	// Validate shift type
	if shiftType < 1 || shiftType > 3 {
		return fmt.Errorf("invalid shift type: %d", shiftType)
	}

	// Validate role
	if employeeRole != "Medic" && employeeRole != "Technical" {
		return fmt.Errorf("invalid employee role: %s", employeeRole)
	}

	// Check if the employee is already assigned to this shift on the given date
	var existingShift model.Shift
	if err := r.db.Where("employee_id = ? AND shift_date = ? AND shift_type = ?", employeeID, shiftDate.Format(dateFormat), shiftType).First(&existingShift).Error; err == nil {
		return fmt.Errorf("employee already assigned to this shift on %s", shiftDate.Format(dateFormat))
	}

	// Count current employees in the shift
	var count int64
	err := r.db.Model(&model.Shift{}).
		Where("shift_date = ? AND shift_type = ? AND employee_role = ?", shiftDate.Format(dateFormat), shiftType, employeeRole).
		Count(&count).Error
	if err != nil {
		return err
	}

	// Enforce maximum limits
	if (employeeRole == "Medic" && count >= 2) || (employeeRole == "Technical" && count >= 4) {
		return fmt.Errorf("maximum limit reached for role %s in shift %d on %s", employeeRole, shiftType, shiftDate.Format(dateFormat))
	}

	// Create the shift
	shift := &model.Shift{
		EmployeeID:   employeeID,
		ShiftDate:    shiftDate,
		ShiftType:    shiftType,
		EmployeeRole: employeeRole,
	}
	return r.db.Create(shift).Error
}

func (r *employeeRepository) GetShiftsByEmployeeID(employeeID int) ([]model.Shift, error) {
	var shifts []model.Shift
	err := r.db.Where("employee_id = ?", employeeID).Order("shift_start ASC").Find(&shifts).Error
	return shifts, err
}

func (r *employeeRepository) GetShiftAvailability(role string) (map[string]int, error) {
	var counts struct {
		Role  string
		Count int
	}
	availability := make(map[string]int)

	err := r.db.Model(&model.Employee{}).
		Select("role, COUNT(*) AS count").
		Joins("JOIN shifts ON employees.id = shifts.employee_id").
		Where("? BETWEEN shifts.shift_start AND shifts.shift_end", time.Now()).
		Group("role").
		Scan(&counts).Error
	if err != nil {
		return nil, err
	}

	availability[counts.Role] = counts.Count
	return availability, nil
}

func (r *employeeRepository) GetShiftsForTimeRange(start, end time.Time) ([]model.Shift, error) {
	var shifts []model.Shift
	err := r.db.Where("shift_start >= ? AND shift_end <= ?", start, end).Find(&shifts).Error
	return shifts, err
}

// Delete marks the employee record as deleted by setting the deleted_at timestamp.
//func (r *employeeRepository) Delete(employeeID uint) error {
//	// First, check if the employee is already soft-deleted
//	var employee model.Employee
//	err := r.db.Select("deleted_at").First(&employee, employeeID).Error
//	if err != nil {
//		// Return error if the employee is not found or any other issue occurs
//		return err
//	}
//
//	if employee.DeletedAt.Valid {
//		// Return an error if the employee is already soft-deleted
//		return gorm.ErrRecordNotFound
//	}
//
//	// If not already soft-deleted, mark the employee as deleted
//	return r.db.Model(&model.Employee{}).Where("id = ?", employeeID).Update("deleted_at", time.Now()).Error
//}
