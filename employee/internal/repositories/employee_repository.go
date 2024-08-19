package repositories

//go:generate mockgen -source=employee_repository.go -destination=employee_repository_gomock.go -package=repositories mountain_service/employee/internal/repositories -imports=gomock=go.uber.org/mock/gomock

import (
	"gorm.io/gorm"
	"time"

	"mountain-service/employee/internal/models"
	"mountain-service/shared/utils"
)

type EmployeeRepository interface {
	Create(employee *models.Employee) error
	GetAll() ([]models.Employee, error)
	Delete(employeeID uint) error
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

// Create creates and employee with the hashed version of its password.
func (repo *employeeRepository) Create(employee *models.Employee) error {
	hashedPassword, err := utils.HashPassword(employee.Password)
	if err != nil {
		return err
	}
	employee.Password = hashedPassword
	return repo.db.Create(employee).Error
}

// GetAll returns all the employees which have deleted_at flag set as NULL in db.
func (repo *employeeRepository) GetAll() ([]models.Employee, error) {
	var employees []models.Employee
	err := repo.db.Where("deleted_at IS NULL").Find(&employees).Error
	return employees, err
}

// Delete marks the employee record as deleted by setting the deleted_at timestamp.
func (repo *employeeRepository) Delete(employeeID uint) error {
	// First, check if the employee is already soft-deleted
	var employee models.Employee
	err := repo.db.Select("deleted_at").First(&employee, employeeID).Error
	if err != nil {
		// Return error if the employee is not found or any other issue occurs
		return err
	}

	if employee.DeletedAt.Valid {
		// Return an error if the employee is already soft-deleted
		return gorm.ErrRecordNotFound
	}

	// If not already soft-deleted, mark the employee as deleted
	return repo.db.Model(&models.Employee{}).Where("id = ?", employeeID).Update("deleted_at", time.Now()).Error
}
