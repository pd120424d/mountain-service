package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/employee/internal/model"
	"github.com/pd120424d/mountain-service/api/employee/internal/repositories"
	sharedAuth "github.com/pd120424d/mountain-service/api/shared/auth"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

func TestEmployeeService_RegisterEmployee(t *testing.T) {
	t.Parallel()

	t.Run("it fails with invalid profile type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		req := employeeV1.EmployeeCreateRequest{
			Username:    "testuser",
			FirstName:   "Test",
			LastName:    "User",
			Password:    "ValidPass123!",
			Email:       "test@example.com",
			ProfileType: "InvalidType",
		}

		response, err := service.RegisterEmployee(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "invalid profile type", err.Error())
	})

	t.Run("it fails when username already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		req := employeeV1.EmployeeCreateRequest{
			Username:    "testuser",
			FirstName:   "Test",
			LastName:    "User",
			Password:    "ValidPass123!",
			Email:       "test@example.com",
			ProfileType: "Medic",
		}

		// Mock existing username check
		existingEmployees := []model.Employee{
			{ID: 1, Username: "testuser"},
		}
		emplRepoMock.EXPECT().ListEmployees(gomock.Any()).Return(existingEmployees, nil)

		response, err := service.RegisterEmployee(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "username already exists", err.Error())
	})

	t.Run("it fails when email already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		req := employeeV1.EmployeeCreateRequest{
			Username:    "testuser",
			FirstName:   "Test",
			LastName:    "User",
			Password:    "ValidPass123!",
			Email:       "test@example.com",
			ProfileType: "Medic",
		}

		// Mock username check (no existing username)
		emplRepoMock.EXPECT().ListEmployees(gomock.Any()).Return([]model.Employee{}, nil)

		// Mock existing email check
		existingEmployees := []model.Employee{
			{ID: 1, Email: "test@example.com"},
		}
		emplRepoMock.EXPECT().ListEmployees(gomock.Any()).Return(existingEmployees, nil)

		response, err := service.RegisterEmployee(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "email already exists", err.Error())
	})

	t.Run("it fails when password validation fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		req := employeeV1.EmployeeCreateRequest{
			Username:    "testuser",
			FirstName:   "Test",
			LastName:    "User",
			Password:    "short", // Invalid password
			Email:       "test@example.com",
			ProfileType: "Medic",
		}

		// Mock username check (no existing users)
		emplRepoMock.EXPECT().ListEmployees(gomock.Any()).Return([]model.Employee{}, nil)
		// Mock email check (no existing users)
		emplRepoMock.EXPECT().ListEmployees(gomock.Any()).Return([]model.Employee{}, nil)

		response, err := service.RegisterEmployee(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "password must be")
	})

	t.Run("it fails when database creation fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		req := employeeV1.EmployeeCreateRequest{
			Username:    "testuser",
			FirstName:   "Test",
			LastName:    "User",
			Password:    "Pass123!",
			Email:       "test@example.com",
			ProfileType: "Medic",
		}

		// Mock username check (no existing users)
		emplRepoMock.EXPECT().ListEmployees(gomock.Any()).Return([]model.Employee{}, nil)
		// Mock email check (no existing users)
		emplRepoMock.EXPECT().ListEmployees(gomock.Any()).Return([]model.Employee{}, nil)
		// Mock database creation failure
		emplRepoMock.EXPECT().Create(gomock.Any()).Return(assert.AnError)

		response, err := service.RegisterEmployee(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("it successfully registers employee", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		req := employeeV1.EmployeeCreateRequest{
			Username:    "testuser",
			FirstName:   "Test",
			LastName:    "User",
			Password:    "Pass123!",
			Email:       "test@example.com",
			ProfileType: "Medic",
			Gender:      "M",
			Phone:       "123456789",
		}

		// Mock username check (no existing username)
		emplRepoMock.EXPECT().ListEmployees(gomock.Any()).Return([]model.Employee{}, nil)

		// Mock email check (no existing email)
		emplRepoMock.EXPECT().ListEmployees(gomock.Any()).Return([]model.Employee{}, nil)

		// Mock successful creation
		emplRepoMock.EXPECT().Create(gomock.Any()).DoAndReturn(func(emp *model.Employee) error {
			emp.ID = 1 // Simulate database assigning ID
			return nil
		})

		response, err := service.RegisterEmployee(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "testuser", response.Username)
		assert.Equal(t, "Test", response.FirstName)
		assert.Equal(t, "User", response.LastName)
		assert.Equal(t, "test@example.com", response.Email)
		assert.Equal(t, "Medic", response.ProfileType)
	})
}

func TestEmployeeService_LoginEmployee(t *testing.T) {
	t.Parallel()

	// Note: Empty username and invalid password format tests are now handled at DTO level

	t.Run("it fails when employee not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		req := employeeV1.EmployeeLogin{
			Username: "nonexistent",
			Password: "Pass123!",
		}

		emplRepoMock.EXPECT().GetEmployeeByUsername("nonexistent").Return(nil, assert.AnError)

		token, err := service.LoginEmployee(context.Background(), req)

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	t.Run("it fails with incorrect password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		req := employeeV1.EmployeeLogin{
			Username: "testuser",
			Password: "Pass123!",
		}

		// Create employee with different password hash
		employee := &model.Employee{
			ID:       1,
			Username: "testuser",
			Password: "$2a$10$differenthashvalue", // Different hash that won't match
		}

		emplRepoMock.EXPECT().GetEmployeeByUsername("testuser").Return(employee, nil)

		token, err := service.LoginEmployee(context.Background(), req)

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	t.Run("it fails with invalid credentials", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		req := employeeV1.EmployeeLogin{
			Username: "testuser",
			Password: "Pass123!",
		}

		emplRepoMock.EXPECT().GetEmployeeByUsername("testuser").Return(nil, assert.AnError)

		token, err := service.LoginEmployee(context.Background(), req)

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	t.Run("it successfully logs in employee", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		req := employeeV1.EmployeeLogin{
			Username: "testuser",
			Password: "Pass123!",
		}

		// Create employee with hashed password for "Pass123!"
		employee := &model.Employee{
			ID:       1,
			Username: "testuser",
			Password: "$2a$10$umEwWgSPqYCkyOEuAMNd7.1mmhRhJVZ3JO1AFq8Z/3bM6uRrwFgDC", // "Pass123!" hashed
		}

		emplRepoMock.EXPECT().GetEmployeeByUsername("testuser").Return(employee, nil)

		token, err := service.LoginEmployee(context.Background(), req)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestEmployeeService_LogoutEmployee(t *testing.T) {
	t.Parallel()

	t.Run("it succeeds when blacklist is available", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		// Create a mock blacklist using gomock
		tokenID := "test-token-123"
		expiresAt := time.Now().Add(time.Hour)
		mockBlacklist := sharedAuth.NewMockTokenBlacklist(ctrl)
		mockBlacklist.EXPECT().BlacklistToken(gomock.Any(), tokenID, expiresAt).Return(nil)
		service := NewEmployeeService(log, emplRepoMock, mockBlacklist)

		err := service.LogoutEmployee(context.Background(), tokenID, expiresAt)

		assert.NoError(t, err)
	})

	t.Run("it succeeds when blacklist is not available", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)
		// Don't set blacklist

		tokenID := "test-token-123"
		expiresAt := time.Now().Add(time.Hour)

		err := service.LogoutEmployee(context.Background(), tokenID, expiresAt)

		assert.NoError(t, err) // Should not error even without blacklist
	})

	t.Run("it fails when blacklist returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		// Create a mock blacklist that returns error using gomock
		mockBlacklist := sharedAuth.NewMockTokenBlacklist(ctrl)
		mockBlacklist.EXPECT().BlacklistToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("mock blacklist error"))
		service = NewEmployeeService(log, emplRepoMock, mockBlacklist)

		tokenID := "test-token-123"
		expiresAt := time.Now().Add(time.Hour)

		err := service.LogoutEmployee(context.Background(), tokenID, expiresAt)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to logout")
	})
}

func TestEmployeeService_ListEmployees(t *testing.T) {
	t.Parallel()

	t.Run("it fails when repository returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		emplRepoMock.EXPECT().GetAll().Return(nil, assert.AnError)

		employees, err := service.ListEmployees(context.Background())

		assert.Error(t, err)
		assert.Nil(t, employees)
		assert.Equal(t, "failed to retrieve employees", err.Error())
	})

	t.Run("it successfully returns employees", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		employees := []model.Employee{
			{
				ID:          1,
				Username:    "user1",
				FirstName:   "Marko",
				LastName:    "Markovic",
				Email:       "john@example.com",
				ProfileType: model.Medic,
			},
			{
				ID:          2,
				Username:    "user2",
				FirstName:   "Marko",
				LastName:    "Markovic",
				Email:       "marko@example.com",
				ProfileType: model.Technical,
			},
		}

		emplRepoMock.EXPECT().GetAll().Return(employees, nil)

		response, err := service.ListEmployees(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response))
		assert.Equal(t, uint(1), response[0].ID)
		assert.Equal(t, "user1", response[0].Username)
		assert.Equal(t, "Medic", response[0].ProfileType)
		assert.Equal(t, uint(2), response[1].ID)
		assert.Equal(t, "user2", response[1].Username)
		assert.Equal(t, "Technical", response[1].ProfileType)
	})
}

func TestEmployeeService_UpdateEmployee(t *testing.T) {
	t.Parallel()

	t.Run("it fails when employee not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		req := employeeV1.EmployeeUpdateRequest{
			FirstName: "Updated",
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).Return(assert.AnError)

		response, err := service.UpdateEmployee(context.Background(), 1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "employee not found", err.Error())
	})

	// Note: Email validation test removed - validation is now handled at DTO level

	t.Run("it fails when update repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		req := employeeV1.EmployeeUpdateRequest{
			FirstName: "Updated",
		}

		employee := model.Employee{
			ID:        1,
			Username:  "testuser",
			FirstName: "Original",
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = employee
			return nil
		})

		emplRepoMock.EXPECT().UpdateEmployee(gomock.Any()).Return(assert.AnError)

		response, err := service.UpdateEmployee(context.Background(), 1, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "failed to update employee", err.Error())
	})

	t.Run("it successfully updates employee", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		req := employeeV1.EmployeeUpdateRequest{
			FirstName: "Updated",
			LastName:  "Name",
		}

		employee := model.Employee{
			ID:        1,
			Username:  "testuser",
			FirstName: "Original",
			LastName:  "Name",
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = employee
			return nil
		})

		emplRepoMock.EXPECT().UpdateEmployee(gomock.Any()).Return(nil)

		response, err := service.UpdateEmployee(context.Background(), 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "Updated", response.FirstName)
	})
}

func TestEmployeeService_DeleteEmployee(t *testing.T) {
	t.Parallel()

	t.Run("it fails when repository returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		emplRepoMock.EXPECT().Delete(uint(1)).Return(assert.AnError)

		err := service.DeleteEmployee(context.Background(), 1)

		assert.Error(t, err)
		assert.Equal(t, "failed to delete employee", err.Error())
	})

	t.Run("it successfully deletes employee", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		emplRepoMock.EXPECT().Delete(uint(1)).Return(nil)

		err := service.DeleteEmployee(context.Background(), 1)

		assert.NoError(t, err)
	})
}

func TestEmployeeService_GetEmployeeByID(t *testing.T) {
	t.Parallel()

	t.Run("it fails when employee not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).Return(assert.AnError)

		employee, err := service.GetEmployeeByID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, employee)
		assert.Equal(t, "employee not found", err.Error())
	})

	t.Run("it successfully returns employee", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		expectedEmployee := model.Employee{
			ID:        1,
			Username:  "testuser",
			FirstName: "Test",
			LastName:  "User",
		}

		emplRepoMock.EXPECT().GetEmployeeByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, emp *model.Employee) error {
			*emp = expectedEmployee
			return nil
		})

		employee, err := service.GetEmployeeByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, employee)
		assert.Equal(t, uint(1), employee.ID)
		assert.Equal(t, "testuser", employee.Username)
	})
}

func TestEmployeeService_GetEmployeeByUsername(t *testing.T) {
	t.Parallel()

	t.Run("it fails when employee not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		emplRepoMock.EXPECT().GetEmployeeByUsername("nonexistent").Return(nil, assert.AnError)

		employee, err := service.GetEmployeeByUsername(context.Background(), "nonexistent")

		assert.Error(t, err)
		assert.Nil(t, employee)
		assert.Equal(t, "employee not found", err.Error())
	})

	t.Run("it successfully returns employee", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		expectedEmployee := &model.Employee{
			ID:        1,
			Username:  "testuser",
			FirstName: "Test",
			LastName:  "User",
		}

		emplRepoMock.EXPECT().GetEmployeeByUsername("testuser").Return(expectedEmployee, nil)

		employee, err := service.GetEmployeeByUsername(context.Background(), "testuser")

		assert.NoError(t, err)
		assert.NotNil(t, employee)
		assert.Equal(t, uint(1), employee.ID)
		assert.Equal(t, "testuser", employee.Username)
	})
}

func TestEmployeeService_ResetAllData(t *testing.T) {
	t.Parallel()

	t.Run("it fails when repository returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		emplRepoMock.EXPECT().ResetAllData().Return(assert.AnError)

		err := service.ResetAllData(context.Background())

		assert.Error(t, err)
		assert.Equal(t, "failed to reset data", err.Error())
	})

	t.Run("it successfully resets all data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		emplRepoMock := repositories.NewMockEmployeeRepository(ctrl)

		service := NewEmployeeService(log, emplRepoMock, nil)

		emplRepoMock.EXPECT().ResetAllData().Return(nil)

		err := service.ResetAllData(context.Background())

		assert.NoError(t, err)
	})
}
