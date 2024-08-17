package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mountain-service/employee/internal/models"
	"mountain-service/employee/internal/repositories"
	"mountain-service/shared/utils"
)

type EmployeeHandler interface {
	CreateEmployee(ctx *gin.Context)
	GetAllEmployees(c *gin.Context)
}

type employeeHandler struct {
	log  utils.Logger
	repo repositories.EmployeeRepository
}

func NewEmployeeHandler(log utils.Logger, repo repositories.EmployeeRepository) EmployeeHandler {
	return &employeeHandler{log: log, repo: repo}
}

// CreateEmployee handles the HTTP POST request to create a new employee.
func (h *employeeHandler) CreateEmployee(ctx *gin.Context) {
	var employee models.Employee
	if err := ctx.ShouldBindJSON(&employee); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Create(&employee); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, employee)
}

// GetAllEmployees handles the HTTP GET request to retrieve all employees.
func (h *employeeHandler) GetAllEmployees(c *gin.Context) {
	employees, err := h.repo.GetAll()
	if err != nil {
		h.log.Errorf("failed to retrieve employees: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Infof("Successfully retrieved %d employees", len(employees))
	c.JSON(http.StatusOK, employees)
}
