package repositories

//go:generate mockgen -source=employee_repository.go -destination=employee_repository_gomock.go -package=repositories mountain_service/employee/internal/repositories -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"fmt"
	"maps"
	"slices"

	"github.com/pd120424d/mountain-service/api/employee/internal/model"
	"github.com/pd120424d/mountain-service/api/shared/utils"

	"gorm.io/gorm"
)

type EmployeeRepository interface {
	Create(employee *model.Employee) error
	GetAll() ([]model.Employee, error)
	GetEmployeeByID(id string, employee *model.Employee) error
	GetEmployeeByUsername(username string) (*model.Employee, error)
	UpdateEmployee(employee *model.Employee) error
	Delete(employeeID uint) error
	ListEmployees(filters map[string]interface{}) ([]model.Employee, error)
}

type employeeRepository struct {
	log utils.Logger
	db  *gorm.DB
}

func NewEmployeeRepository(log utils.Logger, db *gorm.DB) EmployeeRepository {
	return &employeeRepository{log: log.WithName("employeeRepository"), db: db}
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

// GetEmployeeByUsername returns employee by its username or error if it cannot be found.
func (r *employeeRepository) GetEmployeeByUsername(username string) (*model.Employee, error) {
	var employee model.Employee
	err := r.db.Where("username = ?", username).First(&employee).Error
	return &employee, err
}

func (r *employeeRepository) ListEmployees(filters map[string]any) ([]model.Employee, error) {
	allowedColumns := r.allowedColumns()
	var employees []model.Employee
	query := r.db.Model(&model.Employee{})

	// Extract and sort filter keys
	filterKeys := slices.Collect(maps.Keys(filters))
	slices.Sort(filterKeys)

	// Apply filters safely
	for _, key := range filterKeys {
		// Validate key
		if _, ok := allowedColumns[key]; !ok {
			return nil, fmt.Errorf("invalid filter key: %s", key)
		}

		// Extract value
		value := filters[key]

		switch v := value.(type) {
		case string:
			// Use LIKE for string fields
			query = query.Where(fmt.Sprintf("%s LIKE ?", key), fmt.Sprintf("%%%s%%", v))
		case int, int32, int64, float32, float64, bool:
			// Use exact match for non-string types
			query = query.Where(fmt.Sprintf("%s = ?", key), v)
		default:
			return nil, fmt.Errorf("unsupported type for filter key: %s", key)
		}
	}

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

func (r *employeeRepository) allowedColumns() map[string]bool {
	return map[string]bool{
		"id":           true,
		"username":     true,
		"first_name":   true,
		"last_name":    true,
		"gender":       true,
		"phone":        true,
		"email":        true,
		"profile_type": true,
	}
}
