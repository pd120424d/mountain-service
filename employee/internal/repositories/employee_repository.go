package repositories

import (
	"gorm.io/gorm"

	"mountain-service/employee/internal/models"
)

type EmployeeRepository interface {
	Create(employee *models.Employee) error
	GetAll() ([]models.Employee, error)
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (repo *employeeRepository) Create(employee *models.Employee) error {
	return repo.db.Create(employee).Error
}

func (repo *employeeRepository) GetAll() ([]models.Employee, error) {
	var employees []models.Employee
	err := repo.db.Find(&employees).Error
	return employees, err
}
