package service

import (
	"fmt"
	"strings"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/employee/internal/model"
	"github.com/pd120424d/mountain-service/api/employee/internal/repositories"
	sharedAuth "github.com/pd120424d/mountain-service/api/shared/auth"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type employeeService struct {
	log      utils.Logger
	emplRepo repositories.EmployeeRepository
}

func NewEmployeeService(log utils.Logger, emplRepo repositories.EmployeeRepository) EmployeeService {
	return &employeeService{
		log:      log.WithName("employeeService"),
		emplRepo: emplRepo,
	}
}

func (s *employeeService) RegisterEmployee(req employeeV1.EmployeeCreateRequest) (*employeeV1.EmployeeResponse, error) {
	s.log.Infof("Creating new employee with data: %s", req.ToString())

	profileType := model.ProfileTypeFromString(req.ProfileType)
	if !profileType.Valid() {
		s.log.Errorf("invalid profile type: %s", req.ProfileType)
		return nil, fmt.Errorf("invalid profile type")
	}

	employee := model.Employee{
		Username:       req.Username,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Password:       req.Password,
		Gender:         req.Gender,
		Phone:          req.Phone,
		Email:          req.Email,
		ProfilePicture: req.ProfilePicture,
		ProfileType:    profileType,
	}

	// Check for existing username
	usernameFilter := map[string]any{
		"username": employee.Username,
	}
	existingEmployees, err := s.emplRepo.ListEmployees(usernameFilter)
	if err != nil {
		s.log.Errorf("failed to register employee, checking for employee with username failed: %v", err)
		return nil, fmt.Errorf("failed to check for existing username")
	}
	if len(existingEmployees) > 0 {
		s.log.Errorf("failed to register employee: username %s already exists", employee.Username)
		return nil, fmt.Errorf("username already exists")
	}

	// Check for existing email
	emailFilter := map[string]any{
		"email": employee.Email,
	}
	existingEmployees, err = s.emplRepo.ListEmployees(emailFilter)
	if err != nil {
		s.log.Errorf("failed to register employee, checking for employee with email failed: %v", err)
		return nil, fmt.Errorf("failed to check for existing email")
	}
	if len(existingEmployees) > 0 {
		s.log.Errorf("failed to register employee: email %s already exists", employee.Email)
		return nil, fmt.Errorf("email already exists")
	}

	// Validate password
	if err := utils.ValidatePassword(employee.Password); err != nil {
		s.log.Errorf("failed to validate password: %v", err)
		return nil, err
	}

	// Create employee
	if err := s.emplRepo.Create(&employee); err != nil {
		s.log.Errorf("failed to create employee: %v", err)
		// Propagate specific database errors
		if strings.Contains(err.Error(), "invalid db") {
			return nil, fmt.Errorf("invalid db")
		}
		return nil, err
	}

	response := &employeeV1.EmployeeResponse{
		ID:             employee.ID,
		Username:       employee.Username,
		FirstName:      employee.FirstName,
		LastName:       employee.LastName,
		Gender:         employee.Gender,
		Phone:          employee.Phone,
		Email:          employee.Email,
		ProfilePicture: employee.ProfilePicture,
		ProfileType:    employee.ProfileType.String(),
	}

	s.log.Infof("Successfully created employee with ID %d", employee.ID)
	return response, nil
}

func (s *employeeService) LoginEmployee(req employeeV1.EmployeeLogin) (string, error) {
	s.log.Info("Processing employee login")

	// Validate request
	if req.Username == "" {
		s.log.Error("empty username provided")
		return "", fmt.Errorf("invalid credentials")
	}

	if err := utils.ValidatePassword(req.Password); err != nil {
		s.log.Errorf("failed to validate password: %v", err)
		return "", fmt.Errorf("invalid credentials")
	}

	employee, err := s.emplRepo.GetEmployeeByUsername(req.Username)
	if err != nil {
		s.log.Errorf("failed to retrieve employee: %v", err)
		return "", fmt.Errorf("invalid credentials")
	}

	if !sharedAuth.CheckPassword(employee.Password, req.Password) {
		s.log.Error("failed to verify password")
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := sharedAuth.GenerateJWT(employee.ID, employee.Role())
	if err != nil {
		s.log.Errorf("failed to generate token: %v", err)
		return "", fmt.Errorf("failed to generate token")
	}

	s.log.Info("Successfully validated employee and generated JWT token")
	return token, nil
}

func (s *employeeService) ListEmployees() ([]employeeV1.EmployeeResponse, error) {
	s.log.Info("Retrieving list of employees")

	employees, err := s.emplRepo.GetAll()
	if err != nil {
		s.log.Errorf("failed to retrieve employees: %v", err)
		// Propagate specific database errors
		if strings.Contains(err.Error(), "record not found") {
			return nil, fmt.Errorf("record not found")
		}
		return nil, fmt.Errorf("failed to retrieve employees")
	}

	response := make([]employeeV1.EmployeeResponse, 0, len(employees))
	for _, employee := range employees {
		response = append(response, employeeV1.EmployeeResponse{
			ID:             employee.ID,
			Username:       employee.Username,
			FirstName:      employee.FirstName,
			LastName:       employee.LastName,
			Gender:         employee.Gender,
			Phone:          employee.Phone,
			Email:          employee.Email,
			ProfilePicture: employee.ProfilePicture,
			ProfileType:    employee.ProfileType.String(),
		})
	}

	s.log.Infof("Successfully retrieved %d employees", len(response))
	return response, nil
}

func (s *employeeService) UpdateEmployee(employeeID uint, req employeeV1.EmployeeUpdateRequest) (*employeeV1.EmployeeResponse, error) {
	s.log.Infof("Updating employee with ID %d", employeeID)

	var employee model.Employee
	if err := s.emplRepo.GetEmployeeByID(employeeID, &employee); err != nil {
		s.log.Errorf("failed to get employee: %v", err)
		return nil, fmt.Errorf("employee not found")
	}

	// Validate email if provided
	if validationErr := utils.ValidateOptionalEmail(req.Email); validationErr != nil {
		s.log.Errorf("failed to update employee, validation failed: %v", validationErr)
		return nil, validationErr
	}

	model.MapUpdateRequestToEmployee(&req, &employee)

	if err := s.emplRepo.UpdateEmployee(&employee); err != nil {
		s.log.Errorf("failed to update employee: %v", err)
		return nil, fmt.Errorf("failed to update employee")
	}

	response := employee.UpdateResponseFromEmployee()
	s.log.Infof("Successfully updated employee with ID %d", employeeID)
	return &response, nil
}

func (s *employeeService) DeleteEmployee(employeeID uint) error {
	s.log.Infof("Deleting employee with ID %d", employeeID)

	if err := s.emplRepo.Delete(employeeID); err != nil {
		s.log.Errorf("failed to delete employee: %v", err)
		return fmt.Errorf("failed to delete employee")
	}

	s.log.Infof("Employee with ID %d was deleted", employeeID)
	return nil
}

func (s *employeeService) GetEmployeeByID(employeeID uint) (*model.Employee, error) {
	s.log.Infof("Getting employee by ID %d", employeeID)

	var employee model.Employee
	if err := s.emplRepo.GetEmployeeByID(employeeID, &employee); err != nil {
		s.log.Errorf("failed to get employee: %v", err)
		return nil, fmt.Errorf("employee not found")
	}

	return &employee, nil
}

func (s *employeeService) GetEmployeeByUsername(username string) (*model.Employee, error) {
	s.log.Infof("Getting employee by username %s", username)

	employee, err := s.emplRepo.GetEmployeeByUsername(username)
	if err != nil {
		s.log.Errorf("failed to get employee: %v", err)
		return nil, fmt.Errorf("employee not found")
	}

	return employee, nil
}

func (s *employeeService) ResetAllData() error {
	s.log.Warn("Resetting all employee data")

	err := s.emplRepo.ResetAllData()
	if err != nil {
		s.log.Errorf("Failed to reset all data: %v", err)
		return fmt.Errorf("failed to reset data")
	}

	s.log.Info("Successfully reset all employee data")
	return nil
}
