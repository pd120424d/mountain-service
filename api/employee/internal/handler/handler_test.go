package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/employee/internal/model"
	"github.com/pd120424d/mountain-service/api/employee/internal/service"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

func TestEmployeeHandler_RegisterEmployee(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when request payload is invalid json", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		invalidJSON := `{"username": "test", "invalid": json}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(invalidJSON))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.RegisterEmployee(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request payload")
	})

	tests := []struct {
		name     string
		errCode  int
		errMsg   string
		err      error
		response *employeeV1.EmployeeResponse
	}{
		{
			name:    "it returns StatusConflict when username already exists",
			errCode: http.StatusConflict,
			errMsg:  "Username already exists",
			err:     fmt.Errorf("username already exists"),
		},
		{
			name:    "it returns StatusConflict when email already exists",
			errCode: http.StatusConflict,
			errMsg:  "Email already exists",
			err:     fmt.Errorf("email already exists"),
		},
		{
			name:    "it returns StatusBadRequest when profile type is invalid",
			errCode: http.StatusBadRequest,
			errMsg:  "Invalid profile type",
			err:     fmt.Errorf("invalid profile type"),
		},
		{
			name:    "it returns StatusInternalServerError when it fails to check for existing username",
			errCode: http.StatusInternalServerError,
			errMsg:  "Failed to check for existing username",
			err:     fmt.Errorf("failed to check for existing username"),
		},
		{
			name:    "it returns StatusInternalServerError when it fails to check for existing email",
			errCode: http.StatusInternalServerError,
			errMsg:  "Failed to check for existing email",
			err:     fmt.Errorf("failed to check for existing email"),
		},
		{
			name:    "it returns StatusBadRequest when password validation fails",
			errCode: http.StatusBadRequest,
			errMsg:  "password must be between 6 and 10 characters long",
			err:     fmt.Errorf("password must be between 6 and 10 characters long"),
		},
		{
			name:    "it returns StatusInternalServerError when database creation fails",
			errCode: http.StatusInternalServerError,
			errMsg:  "invalid db state",
			err:     fmt.Errorf("invalid db state"),
		},
		{
			name:    "it returns StatusInternalServerError when it fails to register employee with any other reason",
			errCode: http.StatusInternalServerError,
			errMsg:  "Failed to register employee",
			err:     fmt.Errorf("any other error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEmplSvc := service.NewMockEmployeeService(ctrl)
			mockShiftSvc := service.NewMockShiftService(ctrl)
			log := utils.NewTestLogger()
			handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			req := employeeV1.EmployeeCreateRequest{
				FirstName:   "Test",
				LastName:    "User",
				Username:    "existinguser",
				Password:    "Pass123!",
				Email:       "test@example.com",
				Gender:      "Male",
				Phone:       "123456789",
				ProfileType: "Medic",
			}
			payload, _ := json.Marshal(req)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(payload))
			ctx.Request.Header.Set("Content-Type", "application/json")

			mockEmplSvc.EXPECT().RegisterEmployee(req).Return(nil, test.err)

			handler.RegisterEmployee(ctx)

			assert.Equal(t, test.errCode, w.Code)
			assert.Contains(t, w.Body.String(), test.errMsg)
		})
	}

	t.Run("it successfully registers employee", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		req := employeeV1.EmployeeCreateRequest{
			FirstName:   "Test",
			LastName:    "User",
			Username:    "testuser",
			Password:    "Pass123!",
			Email:       "test@example.com",
			Gender:      "Male",
			Phone:       "123456789",
			ProfileType: "Medic",
		}
		payload, _ := json.Marshal(req)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		response := &employeeV1.EmployeeResponse{
			ID:        1,
			Username:  "testuser",
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@example.com",
		}

		mockEmplSvc.EXPECT().RegisterEmployee(req).Return(response, nil)

		handler.RegisterEmployee(ctx)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "testuser")
	})
}

func TestEmployeeHandler_LoginEmployee(t *testing.T) {
	// Note: Not using t.Parallel() because this test modifies environment variables

	t.Run("it returns an error when request payload is invalid json", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		invalidJSON := `{"username": "test", "invalid": json}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(invalidJSON))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.LoginEmployee(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request payload")
	})

	t.Run("it returns an error when admin login has invalid password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		req := employeeV1.EmployeeLogin{
			Username: "admin",
			Password: "wrongpass",
		}
		payload, _ := json.Marshal(req)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.LoginEmployee(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid credentials")
	})

	t.Run("it successfully logs in admin", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Set admin password and JWT secret before test
		os.Setenv("ADMIN_PASSWORD", "admin123")
		os.Setenv("JWT_SECRET", "test-secret-key")
		defer func() {
			os.Unsetenv("ADMIN_PASSWORD")
			os.Unsetenv("JWT_SECRET")
		}()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		req := employeeV1.EmployeeLogin{
			Username: "admin",
			Password: "admin123",
		}
		payload, _ := json.Marshal(req)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.LoginEmployee(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "token")
	})

	t.Run("it returns an error when credentials are invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		req := employeeV1.EmployeeLogin{
			Username: "testuser",
			Password: "wrongpass",
		}
		payload, _ := json.Marshal(req)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		mockEmplSvc.EXPECT().LoginEmployee(req).Return("", fmt.Errorf("invalid credentials"))

		handler.LoginEmployee(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid credentials")
	})

	t.Run("it returns an error when employee login fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		req := employeeV1.EmployeeLogin{
			Username: "testuser",
			Password: "Pass123!",
		}
		payload, _ := json.Marshal(req)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		mockEmplSvc.EXPECT().LoginEmployee(req).Return("", fmt.Errorf("any other error"))

		handler.LoginEmployee(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to login user")
	})

	t.Run("it successfully logs in employee", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		req := employeeV1.EmployeeLogin{
			Username: "testuser",
			Password: "Pass123!",
		}
		payload, _ := json.Marshal(req)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		mockEmplSvc.EXPECT().LoginEmployee(req).Return("jwt-token", nil)

		handler.LoginEmployee(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "jwt-token")
	})
}

func TestEmployeeHandler_OAuth2Token(t *testing.T) {
	// Note: Not using t.Parallel() because this test modifies environment variables

	t.Run("it returns an error when username is not provided", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		// Send form data with only password
		formData := "password=test"
		ctx.Request = httptest.NewRequest(http.MethodPost, "/oauth2/token", strings.NewReader(formData))
		ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		handler.OAuth2Token(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "username and password are required")
	})

	t.Run("it returns an error when password is not provided", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		// Send form data with only username
		formData := "username=test"
		ctx.Request = httptest.NewRequest(http.MethodPost, "/oauth2/token", strings.NewReader(formData))
		ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		handler.OAuth2Token(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "username and password are required")
	})

	t.Run("it returns an error when admin password is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		os.Setenv("ADMIN_PASSWORD", "admin123")
		defer os.Unsetenv("ADMIN_PASSWORD")

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		formData := "username=admin&password=wrongpass"
		ctx.Request = httptest.NewRequest(http.MethodPost, "/oauth2/token", strings.NewReader(formData))
		ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		handler.OAuth2Token(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid credentials")
	})

	t.Run("it successfully authenticates admin via OAuth2", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		os.Setenv("ADMIN_PASSWORD", "admin123")
		os.Setenv("JWT_SECRET", "test-secret-key")
		defer func() {
			os.Unsetenv("ADMIN_PASSWORD")
			os.Unsetenv("JWT_SECRET")
		}()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		// Create form data for OAuth2 endpoint with correct admin credentials
		formData := "username=admin&password=admin123"
		ctx.Request = httptest.NewRequest(http.MethodPost, "/oauth2/token", strings.NewReader(formData))
		ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		handler.OAuth2Token(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "access_token")
		assert.Contains(t, w.Body.String(), "Bearer")
		assert.Contains(t, w.Body.String(), "expires_in")
	})

	t.Run("it returns an error when LoginEmployee call fails for regular employee", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		// Create form data for OAuth2 endpoint with correct employee credentials
		formData := "username=testuser&password=Pass123!"
		ctx.Request = httptest.NewRequest(http.MethodPost, "/oauth2/token", strings.NewReader(formData))
		ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		mockEmplSvc.EXPECT().LoginEmployee(employeeV1.EmployeeLogin{
			Username: "testuser",
			Password: "Pass123!",
		}).Return("", fmt.Errorf("invalid credentials provided"))

		handler.OAuth2Token(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid credentials")
	})

	t.Run("it successfully authenticates employee via OAuth2", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		// Create form data for OAuth2 endpoint with correct employee credentials
		formData := "username=testuser&password=Pass123!"
		ctx.Request = httptest.NewRequest(http.MethodPost, "/oauth2/token", strings.NewReader(formData))
		ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		mockEmplSvc.EXPECT().LoginEmployee(employeeV1.EmployeeLogin{
			Username: "testuser",
			Password: "Pass123!",
		}).Return("jwt-token", nil)

		handler.OAuth2Token(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "access_token")
		assert.Contains(t, w.Body.String(), "Bearer")
		assert.Contains(t, w.Body.String(), "expires_in")
	})
}

func TestEmployeeHandler_ListEmployees(t *testing.T) {
	t.Parallel()

	t.Run("it returns StatusInternalServerError when service call fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		mockEmplSvc.EXPECT().ListEmployees().Return(nil, fmt.Errorf("any other error"))

		handler.ListEmployees(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to retrieve employees")
	})

	t.Run("it successfully returns employees", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		employees := []employeeV1.EmployeeResponse{
			{
				ID:        1,
				Username:  "testuser",
				FirstName: "Test",
				LastName:  "User",
				Email:     "test@example.com",
			},
		}

		mockEmplSvc.EXPECT().ListEmployees().Return(employees, nil)

		handler.ListEmployees(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "testuser")
	})
}

func TestEmployeeHandler_UpdateEmployee(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when employee ID is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "invalid"}}

		req := employeeV1.EmployeeUpdateRequest{
			FirstName: "Updated",
			LastName:  "User",
		}
		payload, _ := json.Marshal(req)
		ctx.Request = httptest.NewRequest(http.MethodPut, "/employees/invalid", bytes.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateEmployee(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid employee ID")
	})

	t.Run("it returns an error when request payload is invalid json", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}

		invalidJSON := `{"firstName": "test", "invalid": json}`
		ctx.Request = httptest.NewRequest(http.MethodPut, "/employees/1", strings.NewReader(invalidJSON))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateEmployee(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request payload")
	})

	tests := []struct {
		name    string
		errCode int
		errMsg  string
		err     error
	}{
		{
			name:    "it returns StatusNotFound when employee not found",
			errCode: http.StatusNotFound,
			errMsg:  "Employee not found",
			err:     fmt.Errorf("employee not found"),
		},
		{
			name:    "it returns StatusBadRequest when email validation fails",
			errCode: http.StatusBadRequest,
			errMsg:  "invalid email format",
			err:     fmt.Errorf("mail: invalid email format"),
		},
		{
			name:    "it returns StatusInternalServerError when update fails with other error",
			errCode: http.StatusInternalServerError,
			errMsg:  "Failed to update employee",
			err:     fmt.Errorf("database connection failed"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEmplSvc := service.NewMockEmployeeService(ctrl)
			mockShiftSvc := service.NewMockShiftService(ctrl)
			log := utils.NewTestLogger()
			handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = gin.Params{{Key: "id", Value: "1"}}

			req := employeeV1.EmployeeUpdateRequest{
				FirstName: "Updated",
				LastName:  "User",
				Email:     "updated@example.com",
			}
			payload, _ := json.Marshal(req)
			ctx.Request = httptest.NewRequest(http.MethodPut, "/employees/1", bytes.NewReader(payload))
			ctx.Request.Header.Set("Content-Type", "application/json")

			mockEmplSvc.EXPECT().UpdateEmployee(uint(1), req).Return(nil, test.err)

			handler.UpdateEmployee(ctx)

			assert.Equal(t, test.errCode, w.Code)
			assert.Contains(t, w.Body.String(), test.errMsg)
		})
	}

	t.Run("it successfully updates employee", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}

		req := employeeV1.EmployeeUpdateRequest{
			FirstName: "Updated",
			LastName:  "User",
			Email:     "updated@example.com",
		}
		payload, _ := json.Marshal(req)
		ctx.Request = httptest.NewRequest(http.MethodPut, "/employees/1", bytes.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		response := &employeeV1.EmployeeResponse{
			ID:        1,
			Username:  "testuser",
			FirstName: "Updated",
			LastName:  "User",
			Email:     "updated@example.com",
		}

		mockEmplSvc.EXPECT().UpdateEmployee(uint(1), req).Return(response, nil)

		handler.UpdateEmployee(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Updated")
		assert.Contains(t, w.Body.String(), "updated@example.com")
	})
}

func TestEmployeeHandler_DeleteEmployee(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when employee ID is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "invalid"}}

		ctx.Request = httptest.NewRequest(http.MethodDelete, "/employees/invalid", nil)

		handler.DeleteEmployee(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid employee ID")
	})

	t.Run("it returns StatusInternalServerError when delete fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}

		ctx.Request = httptest.NewRequest(http.MethodDelete, "/employees/1", nil)

		mockEmplSvc.EXPECT().DeleteEmployee(uint(1)).Return(fmt.Errorf("database error"))

		handler.DeleteEmployee(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to delete employee")
	})

	t.Run("it successfully deletes employee", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}

		ctx.Request = httptest.NewRequest(http.MethodDelete, "/employees/1", nil)

		mockEmplSvc.EXPECT().DeleteEmployee(uint(1)).Return(nil)

		handler.DeleteEmployee(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Employee deleted successfully")
	})
}

func TestEmployeeHandler_AssignShift(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when employee ID is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "invalid"}}

		req := employeeV1.AssignShiftRequest{
			ShiftDate: "2024-01-15",
			ShiftType: 1,
		}
		payload, _ := json.Marshal(req)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/employees/invalid/shifts", bytes.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.AssignShift(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid employee ID")
	})

	t.Run("it returns an error when employee ID is zero or negative", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "0"}}

		req := employeeV1.AssignShiftRequest{
			ShiftDate: "2024-01-15",
			ShiftType: 1,
		}
		payload, _ := json.Marshal(req)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/employees/0/shifts", bytes.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.AssignShift(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid employee ID")
	})

	t.Run("it returns an error when request payload is invalid json", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}

		invalidJSON := `{"shiftDate": "2024-01-15", "invalid": json}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/employees/1/shifts", strings.NewReader(invalidJSON))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.AssignShift(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	tests := []struct {
		name    string
		errCode int
		errMsg  string
		err     error
	}{
		{
			name:    "it returns StatusNotFound when employee not found",
			errCode: http.StatusNotFound,
			errMsg:  "employee not found",
			err:     fmt.Errorf("employee not found"),
		},
		{
			name:    "it returns StatusBadRequest when shift date format is invalid",
			errCode: http.StatusBadRequest,
			errMsg:  "invalid shift date format",
			err:     fmt.Errorf("invalid shift date format"),
		},
		{
			name:    "it returns StatusBadRequest when shift is in the past",
			errCode: http.StatusBadRequest,
			errMsg:  "cannot assign shift in the past",
			err:     fmt.Errorf("cannot assign shift in the past"),
		},
		{
			name:    "it returns StatusBadRequest when shift is too far in advance",
			errCode: http.StatusBadRequest,
			errMsg:  "cannot assign shift more than 3 months in advance",
			err:     fmt.Errorf("cannot assign shift more than 3 months in advance"),
		},
		{
			name:    "it returns StatusConflict when employee already assigned to shift",
			errCode: http.StatusConflict,
			errMsg:  "employee is already assigned to this shift",
			err:     fmt.Errorf("employee is already assigned to this shift"),
		},
		{
			name:    "it returns StatusConflict when maximum capacity reached",
			errCode: http.StatusConflict,
			errMsg:  "maximum capacity for this role reached in the selected shift",
			err:     fmt.Errorf("maximum capacity for this role reached in the selected shift"),
		},
		{
			name:    "it returns StatusInternalServerError when assign fails with other error",
			errCode: http.StatusInternalServerError,
			errMsg:  "internal server error",
			err:     fmt.Errorf("database connection failed"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEmplSvc := service.NewMockEmployeeService(ctrl)
			mockShiftSvc := service.NewMockShiftService(ctrl)
			log := utils.NewTestLogger()
			handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = gin.Params{{Key: "id", Value: "1"}}

			req := employeeV1.AssignShiftRequest{
				ShiftDate: "2024-01-15",
				ShiftType: 1,
			}
			payload, _ := json.Marshal(req)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/employees/1/shifts", bytes.NewReader(payload))
			ctx.Request.Header.Set("Content-Type", "application/json")

			mockShiftSvc.EXPECT().AssignShift(uint(1), req).Return(nil, test.err)

			handler.AssignShift(ctx)

			assert.Equal(t, test.errCode, w.Code)
			assert.Contains(t, w.Body.String(), test.errMsg)
		})
	}

	t.Run("it successfully assigns shift", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}

		req := employeeV1.AssignShiftRequest{
			ShiftDate: "2024-01-15",
			ShiftType: 1,
		}
		payload, _ := json.Marshal(req)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/employees/1/shifts", bytes.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		response := &employeeV1.AssignShiftResponse{
			ID:        1,
			ShiftDate: "2024-01-15",
			ShiftType: 1,
		}

		mockShiftSvc.EXPECT().AssignShift(uint(1), req).Return(response, nil)

		handler.AssignShift(ctx)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "2024-01-15")
	})
}

func TestEmployeeHandler_GetShifts(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when employee ID is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "invalid"}}

		ctx.Request = httptest.NewRequest(http.MethodGet, "/employees/invalid/shifts", nil)

		handler.GetShifts(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid employee ID")
	})

	t.Run("it returns an error when employee ID is zero or negative", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "0"}}

		ctx.Request = httptest.NewRequest(http.MethodGet, "/employees/0/shifts", nil)

		handler.GetShifts(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid employee ID")
	})

	t.Run("it returns StatusInternalServerError when service call fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}

		ctx.Request = httptest.NewRequest(http.MethodGet, "/employees/1/shifts", nil)

		mockShiftSvc.EXPECT().GetShifts(uint(1)).Return(nil, fmt.Errorf("database error"))

		handler.GetShifts(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "internal server error")
	})

	t.Run("it successfully returns shifts", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}

		ctx.Request = httptest.NewRequest(http.MethodGet, "/employees/1/shifts", nil)

		shifts := []employeeV1.ShiftResponse{
			{
				ID:        1,
				ShiftType: 1,
			},
		}

		mockShiftSvc.EXPECT().GetShifts(uint(1)).Return(shifts, nil)

		handler.GetShifts(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"id":1`)
	})
}

func TestEmployeeHandler_GetShiftsAvailability(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when days parameter is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest(http.MethodGet, "/shifts/availability?days=invalid", nil)
		ctx.Set("employeeID", uint(1))

		handler.GetShiftsAvailability(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid days parameter")
	})

	t.Run("it returns an error when days parameter is zero or negative", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest(http.MethodGet, "/shifts/availability?days=0", nil)
		ctx.Set("employeeID", uint(1))

		handler.GetShiftsAvailability(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid days parameter")
	})

	t.Run("it returns StatusInternalServerError when service call fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest(http.MethodGet, "/shifts/availability?days=7", nil)
		ctx.Set("employeeID", uint(1))

		mockShiftSvc.EXPECT().GetShiftsAvailability(uint(1), 7).Return(nil, fmt.Errorf("database error"))

		handler.GetShiftsAvailability(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "internal server error")
	})

	t.Run("it successfully returns shifts availability with default days", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest(http.MethodGet, "/shifts/availability", nil)
		ctx.Set("employeeID", uint(1))

		response := &employeeV1.ShiftAvailabilityResponse{}

		mockShiftSvc.EXPECT().GetShiftsAvailability(uint(1), 7).Return(response, nil)

		handler.GetShiftsAvailability(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("it successfully returns shifts availability with custom days", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest(http.MethodGet, "/shifts/availability?days=14", nil)
		ctx.Set("employeeID", uint(1))

		response := &employeeV1.ShiftAvailabilityResponse{}

		mockShiftSvc.EXPECT().GetShiftsAvailability(uint(1), 14).Return(response, nil)

		handler.GetShiftsAvailability(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestEmployeeHandler_RemoveShift(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when employee ID is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "invalid"}}

		req := employeeV1.RemoveShiftRequest{
			ShiftDate: "2024-01-15",
			ShiftType: 1,
		}
		payload, _ := json.Marshal(req)
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/employees/invalid/shifts", bytes.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.RemoveShift(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid employee ID")
	})

	t.Run("it returns an error when employee ID is zero or negative", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "0"}}

		req := employeeV1.RemoveShiftRequest{
			ShiftDate: "2024-01-15",
			ShiftType: 1,
		}
		payload, _ := json.Marshal(req)
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/employees/0/shifts", bytes.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.RemoveShift(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid employee ID")
	})

	t.Run("it returns an error when request payload is invalid json", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}

		invalidJSON := `{"shiftDate": "2024-01-15", "invalid": json}`
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/employees/1/shifts", strings.NewReader(invalidJSON))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.RemoveShift(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("it returns StatusBadRequest when shift date format is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}

		req := employeeV1.RemoveShiftRequest{
			ShiftDate: "2024-01-15",
			ShiftType: 1,
		}
		payload, _ := json.Marshal(req)
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/employees/1/shifts", bytes.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		mockShiftSvc.EXPECT().RemoveShift(uint(1), req).Return(fmt.Errorf("invalid shift date format"))

		handler.RemoveShift(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid shift date format")
	})

	t.Run("it returns StatusInternalServerError when remove fails with other error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}

		req := employeeV1.RemoveShiftRequest{
			ShiftDate: "2024-01-15",
			ShiftType: 1,
		}
		payload, _ := json.Marshal(req)
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/employees/1/shifts", bytes.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		mockShiftSvc.EXPECT().RemoveShift(uint(1), req).Return(fmt.Errorf("database error"))

		handler.RemoveShift(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "internal server error")
	})

	t.Run("it successfully removes shift", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}

		req := employeeV1.RemoveShiftRequest{
			ShiftDate: "2024-01-15",
			ShiftType: 1,
		}
		payload, _ := json.Marshal(req)
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/employees/1/shifts", bytes.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		mockShiftSvc.EXPECT().RemoveShift(uint(1), req).Return(nil)

		handler.RemoveShift(ctx)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

func TestEmployeeHandler_ResetAllData(t *testing.T) {
	t.Parallel()

	t.Run("it returns StatusInternalServerError when reset fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest(http.MethodDelete, "/admin/reset", nil)

		mockEmplSvc.EXPECT().ResetAllData().Return(fmt.Errorf("database error"))

		handler.ResetAllData(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to reset data")
	})

	t.Run("it successfully resets all data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest(http.MethodDelete, "/admin/reset", nil)

		mockEmplSvc.EXPECT().ResetAllData().Return(nil)

		handler.ResetAllData(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "All data has been successfully reset")
	})
}

func TestEmployeeHandler_GetOnCallEmployees(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when shift_buffer parameter is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest(http.MethodGet, "/on-call?shift_buffer=invalid", nil)

		handler.GetOnCallEmployees(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid shift_buffer format")
	})

	t.Run("it returns StatusInternalServerError when service call fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest(http.MethodGet, "/on-call", nil)

		mockShiftSvc.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("database error"))

		handler.GetOnCallEmployees(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to retrieve on-call employees")
	})

	t.Run("it successfully returns on-call employees without shift_buffer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest(http.MethodGet, "/on-call", nil)

		employees := []employeeV1.EmployeeResponse{
			{
				ID:        1,
				Username:  "testuser",
				FirstName: "Test",
				LastName:  "User",
			},
		}

		mockShiftSvc.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return(employees, nil)

		handler.GetOnCallEmployees(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "testuser")
	})

	t.Run("it successfully returns on-call employees with shift_buffer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest(http.MethodGet, "/on-call?shift_buffer=1h", nil)

		employees := []employeeV1.EmployeeResponse{
			{
				ID:        1,
				Username:  "testuser",
				FirstName: "Test",
				LastName:  "User",
			},
		}

		mockShiftSvc.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return(employees, nil)

		handler.GetOnCallEmployees(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "testuser")
	})
}

func TestEmployeeHandler_CheckActiveEmergencies(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when employee ID is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "invalid"}}

		ctx.Request = httptest.NewRequest(http.MethodGet, "/employees/invalid/emergencies", nil)

		handler.CheckActiveEmergencies(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid employee ID")
	})

	t.Run("it returns StatusNotFound when employee not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}

		ctx.Request = httptest.NewRequest(http.MethodGet, "/employees/1/emergencies", nil)

		mockEmplSvc.EXPECT().GetEmployeeByID(uint(1)).Return(nil, fmt.Errorf("employee not found"))

		handler.CheckActiveEmergencies(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Employee not found")
	})

	t.Run("it successfully returns active emergencies status", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}

		ctx.Request = httptest.NewRequest(http.MethodGet, "/employees/1/emergencies", nil)

		// Import the model package to use the correct type
		employee := &model.Employee{
			ID:       1,
			Username: "testuser",
		}

		mockEmplSvc.EXPECT().GetEmployeeByID(uint(1)).Return(employee, nil)

		handler.CheckActiveEmergencies(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "hasActiveEmergencies")
		assert.Contains(t, w.Body.String(), "false") // Currently always returns false as per TODO
	})
}

func TestEmployeeHandler_GetShiftWarnings(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when employee ID is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "invalid"}}

		ctx.Request = httptest.NewRequest(http.MethodGet, "/employees/invalid/shift-warnings", nil)

		handler.GetShiftWarnings(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid employee ID")
	})

	t.Run("it returns an error when employee ID is zero or negative", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "0"}}

		ctx.Request = httptest.NewRequest(http.MethodGet, "/employees/0/shift-warnings", nil)

		handler.GetShiftWarnings(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid employee ID")
	})

	t.Run("it returns StatusNotFound when employee not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}

		ctx.Request = httptest.NewRequest(http.MethodGet, "/employees/1/shift-warnings", nil)

		mockShiftSvc.EXPECT().GetShiftWarnings(uint(1)).Return(nil, fmt.Errorf("employee not found"))

		handler.GetShiftWarnings(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "employee not found")
	})

	t.Run("it returns StatusInternalServerError when service call fails with other error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}

		ctx.Request = httptest.NewRequest(http.MethodGet, "/employees/1/shift-warnings", nil)

		mockShiftSvc.EXPECT().GetShiftWarnings(uint(1)).Return(nil, fmt.Errorf("database connection failed"))

		handler.GetShiftWarnings(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "internal server error")
	})

	t.Run("it successfully returns shift warnings when warnings exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}

		ctx.Request = httptest.NewRequest(http.MethodGet, "/employees/1/shift-warnings", nil)

		warnings := []string{
			"You have only 3 shifts scheduled in the next 2 weeks. Consider scheduling more shifts to meet the 5 days/week quota.",
			"There is insufficient coverage for some shifts in the next 2 weeks.",
		}

		mockShiftSvc.EXPECT().GetShiftWarnings(uint(1)).Return(warnings, nil)

		handler.GetShiftWarnings(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "warnings")
		assert.Contains(t, w.Body.String(), "only 3 shifts scheduled")
		assert.Contains(t, w.Body.String(), "insufficient coverage")
	})

	t.Run("it successfully returns empty warnings when no warnings exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}

		ctx.Request = httptest.NewRequest(http.MethodGet, "/employees/1/shift-warnings", nil)

		warnings := []string{} // Empty warnings array

		mockShiftSvc.EXPECT().GetShiftWarnings(uint(1)).Return(warnings, nil)

		handler.GetShiftWarnings(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "warnings")
		assert.Contains(t, w.Body.String(), "[]") // Empty array in JSON
	})
}

func TestEmployeeHandler_GetAdminShiftsAvailability(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when days parameter is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/admin/shifts/availability?days=invalid", nil)

		handler.GetAdminShiftsAvailability(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Days must be a number between 1 and 90")
	})

	t.Run("it returns an error when days parameter is zero or negative", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/admin/shifts/availability?days=0", nil)

		handler.GetAdminShiftsAvailability(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Days must be a number between 1 and 90")
	})

	t.Run("it returns an error when days parameter exceeds maximum", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEmplSvc := service.NewMockEmployeeService(ctrl)
		mockShiftSvc := service.NewMockShiftService(ctrl)
		log := utils.NewTestLogger()
		handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/admin/shifts/availability?days=100", nil)

		handler.GetAdminShiftsAvailability(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Days must be a number between 1 and 90")
	})

	tests := []struct {
		name           string
		queryParams    string
		setupMocks     func(*service.MockShiftService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "it returns StatusInternalServerError when service call fails",
			queryParams: "?days=7",
			setupMocks: func(mockShiftSvc *service.MockShiftService) {
				mockShiftSvc.EXPECT().GetShiftsAvailability(uint(1), 7).Return(nil, fmt.Errorf("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"error":"Failed to retrieve shifts availability"`,
		},
		{
			name:        "it successfully returns admin shifts availability with default days",
			queryParams: "",
			setupMocks: func(mockShiftSvc *service.MockShiftService) {
				testDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				response := &employeeV1.ShiftAvailabilityResponse{
					Days: map[time.Time]employeeV1.ShiftAvailabilityPerDay{
						testDate: {
							FirstShift: employeeV1.ShiftAvailability{
								MedicSlotsAvailable:     2,
								TechnicalSlotsAvailable: 4,
								IsAssignedToEmployee:    false,
								IsFullyBooked:           false,
							},
							SecondShift: employeeV1.ShiftAvailability{
								MedicSlotsAvailable:     1,
								TechnicalSlotsAvailable: 2,
								IsAssignedToEmployee:    false,
								IsFullyBooked:           false,
							},
							ThirdShift: employeeV1.ShiftAvailability{
								MedicSlotsAvailable:     2,
								TechnicalSlotsAvailable: 4,
								IsAssignedToEmployee:    false,
								IsFullyBooked:           false,
							},
						},
					},
				}
				mockShiftSvc.EXPECT().GetShiftsAvailability(uint(1), 7).Return(response, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"days"`,
		},
		{
			name:        "it successfully returns admin shifts availability with custom days",
			queryParams: "?days=14",
			setupMocks: func(mockShiftSvc *service.MockShiftService) {
				testDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				response := &employeeV1.ShiftAvailabilityResponse{
					Days: map[time.Time]employeeV1.ShiftAvailabilityPerDay{
						testDate: {
							FirstShift: employeeV1.ShiftAvailability{
								MedicSlotsAvailable:     2,
								TechnicalSlotsAvailable: 4,
								IsAssignedToEmployee:    false,
								IsFullyBooked:           false,
							},
							SecondShift: employeeV1.ShiftAvailability{
								MedicSlotsAvailable:     1,
								TechnicalSlotsAvailable: 2,
								IsAssignedToEmployee:    false,
								IsFullyBooked:           false,
							},
							ThirdShift: employeeV1.ShiftAvailability{
								MedicSlotsAvailable:     2,
								TechnicalSlotsAvailable: 4,
								IsAssignedToEmployee:    false,
								IsFullyBooked:           false,
							},
						},
					},
				}
				mockShiftSvc.EXPECT().GetShiftsAvailability(uint(1), 14).Return(response, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"days"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEmplSvc := service.NewMockEmployeeService(ctrl)
			mockShiftSvc := service.NewMockShiftService(ctrl)
			tt.setupMocks(mockShiftSvc)

			log := utils.NewTestLogger()
			handler := NewEmployeeHandler(log, mockEmplSvc, mockShiftSvc)

			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest("GET", "/admin/shifts/availability"+tt.queryParams, nil)

			handler.GetAdminShiftsAvailability(ctx)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}
