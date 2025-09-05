package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/employee/internal/model"
	"github.com/pd120424d/mountain-service/api/employee/internal/repositories"
	sharedAuth "github.com/pd120424d/mountain-service/api/shared/auth"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type employeeService struct {
	log       utils.Logger
	emplRepo  repositories.EmployeeRepository
	blacklist sharedAuth.TokenBlacklist
}

func NewEmployeeService(log utils.Logger, emplRepo repositories.EmployeeRepository, blacklist sharedAuth.TokenBlacklist) EmployeeService {
	return &employeeService{
		log:       log.WithName("employeeService"),
		emplRepo:  emplRepo,
		blacklist: blacklist,
	}
}

func (s *employeeService) RegisterEmployee(ctx context.Context, req employeeV1.EmployeeCreateRequest) (*employeeV1.EmployeeResponse, error) {
	log := s.log.WithContext(ctx)
	log.Infof("Creating new employee with data: %s", req.ToString())

	profileType := model.ProfileTypeFromString(req.ProfileType)
	if !profileType.Valid() {
		log.Errorf("invalid profile type: %s", req.ProfileType)
		return nil, fmt.Errorf("invalid profile type")
	}

	employee := model.Employee{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
		Gender:    req.Gender,
		Phone:     req.Phone,
		Email:     req.Email,

		ProfileType: profileType,
	}

	// Check for existing username
	usernameFilter := map[string]any{
		"username": employee.Username,
	}
	existingEmployees, err := s.emplRepo.ListEmployees(ctx, usernameFilter)
	if err != nil {
		log.Errorf("failed to register employee, checking for employee with username failed: %v", err)
		return nil, fmt.Errorf("failed to check for existing username")
	}
	if len(existingEmployees) > 0 {
		log.Errorf("failed to register employee: username %s already exists", employee.Username)
		return nil, fmt.Errorf("username already exists")
	}

	// Check for existing email
	emailFilter := map[string]any{
		"email": employee.Email,
	}
	existingEmployees, err = s.emplRepo.ListEmployees(ctx, emailFilter)
	if err != nil {
		log.Errorf("failed to register employee, checking for employee with email failed: %v", err)
		return nil, fmt.Errorf("failed to check for existing email")
	}
	if len(existingEmployees) > 0 {
		log.Errorf("failed to register employee: email %s already exists", employee.Email)
		return nil, fmt.Errorf("email already exists")
	}

	// Validate password
	if err := utils.ValidatePassword(employee.Password); err != nil {
		log.Errorf("failed to validate password: %v", err)
		return nil, err
	}

	// Create employee
	if err := s.emplRepo.Create(ctx, &employee); err != nil {
		log.Errorf("failed to create employee: %v", err)
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

	log.Infof("Successfully created employee with ID %d", employee.ID)
	return response, nil
}

func (s *employeeService) LoginEmployee(ctx context.Context, req employeeV1.EmployeeLogin) (string, error) {
	log := s.log.WithContext(ctx)
	log.Info("Processing employee login")

	employee, err := s.emplRepo.GetEmployeeByUsername(ctx, req.Username)
	if err != nil {
		log.Errorf("failed to retrieve employee: %v", err)
		return "", fmt.Errorf("invalid credentials")
	}

	if !sharedAuth.CheckPassword(employee.Password, req.Password) {
		log.Error("failed to verify password")
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := sharedAuth.GenerateJWT(employee.ID, employee.Role())
	if err != nil {
		log.Errorf("failed to generate token: %v", err)
		return "", fmt.Errorf("failed to generate token")
	}

	log.Info("Successfully validated employee and generated JWT token")
	return token, nil
}

func (s *employeeService) LogoutEmployee(ctx context.Context, tokenID string, expiresAt time.Time) error {
	log := s.log.WithContext(ctx)
	log.Info("Processing employee logout")

	if s.blacklist == nil {
		s.log.Warn("Token blacklist not available, logout will not invalidate token")
		return nil
	}

	if err := s.blacklist.BlacklistToken(ctx, tokenID, expiresAt); err != nil {
		log.Errorf("failed to blacklist token: %v", err)
		return fmt.Errorf("failed to logout: %w", err)
	}

	log.Info("Successfully logged out employee and blacklisted token")
	return nil
}

func (s *employeeService) ListEmployees(ctx context.Context) ([]employeeV1.EmployeeResponse, error) {
	log := s.log.WithContext(ctx)
	log.Info("Retrieving list of employees")

	employees, err := s.emplRepo.GetAll(ctx)
	if err != nil {
		log.Errorf("failed to retrieve employees: %v", err)
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

	log.Infof("Successfully retrieved %d employees", len(response))
	return response, nil
}

func (s *employeeService) UpdateEmployee(ctx context.Context, employeeID uint, req employeeV1.EmployeeUpdateRequest) (*employeeV1.EmployeeResponse, error) {
	log := s.log.WithContext(ctx)
	log.Infof("Updating employee with ID %d", employeeID)

	var employee model.Employee
	if err := s.emplRepo.GetEmployeeByID(ctx, employeeID, &employee); err != nil {
		log.Errorf("failed to get employee: %v", err)
		return nil, fmt.Errorf("employee not found")
	}

	model.MapUpdateRequestToEmployee(&req, &employee)

	if err := s.emplRepo.UpdateEmployee(ctx, &employee); err != nil {
		log.Errorf("failed to update employee: %v", err)
		return nil, fmt.Errorf("failed to update employee")
	}

	response := employee.UpdateResponseFromEmployee()
	log.Infof("Successfully updated employee with ID %d", employeeID)
	return &response, nil
}

func (s *employeeService) DeleteEmployee(ctx context.Context, employeeID uint) error {
	log := s.log.WithContext(ctx)
	log.Infof("Deleting employee with ID %d", employeeID)

	if err := s.emplRepo.Delete(ctx, employeeID); err != nil {
		log.Errorf("failed to delete employee: %v", err)
		return fmt.Errorf("failed to delete employee")
	}

	log.Infof("Employee with ID %d was deleted", employeeID)
	return nil
}

func (s *employeeService) GetEmployeeByID(ctx context.Context, employeeID uint) (*model.Employee, error) {
	log := s.log.WithContext(ctx)
	log.Infof("Getting employee by ID %d", employeeID)

	var employee model.Employee
	if err := s.emplRepo.GetEmployeeByID(ctx, employeeID, &employee); err != nil {
		log.Errorf("failed to get employee: %v", err)
		return nil, fmt.Errorf("employee not found")
	}

	return &employee, nil
}

func (s *employeeService) GetEmployeeByUsername(ctx context.Context, username string) (*model.Employee, error) {
	log := s.log.WithContext(ctx)
	log.Infof("Getting employee by username %s", username)

	employee, err := s.emplRepo.GetEmployeeByUsername(ctx, username)
	if err != nil {
		log.Errorf("failed to get employee: %v", err)
		return nil, fmt.Errorf("employee not found")
	}

	return employee, nil
}

func (s *employeeService) ResetAllData(ctx context.Context) error {
	log := s.log.WithContext(ctx)
	log.Warn("Resetting all employee data")

	err := s.emplRepo.ResetAllData(ctx)
	if err != nil {
		log.Errorf("Failed to reset all data: %v", err)
		return fmt.Errorf("failed to reset data")
	}

	log.Info("Successfully reset all employee data")
	return nil
}
