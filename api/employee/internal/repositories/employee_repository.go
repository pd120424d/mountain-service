package repositories

//go:generate mockgen -source=employee_repository.go -destination=employee_repository_gomock.go -package=repositories mountain_service/employee/internal/repositories -imports=gomock=go.uber.org/mock/gomock

import (
	"api/employee/internal/model"
	"api/shared/utils"
	"gorm.io/gorm"
)

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

// TODO:  Soft delete
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
