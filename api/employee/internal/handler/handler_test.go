package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gorm.io/gorm"

	"github.com/pd120424d/mountain-service/api/employee/internal/model"
	"github.com/pd120424d/mountain-service/api/employee/internal/repositories"
	"github.com/pd120424d/mountain-service/api/shared/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestEmployeeHandler_RegisterEmployee(t *testing.T) {
	log := utils.NewTestLogger()

	t.Run("it returns an error when request payload is invalid json", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		invalidPayload := `{
			"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(invalidPayload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := NewEmployeeHandler(log, nil, nil)
		handler.RegisterEmployee(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "{\"error\":\"Invalid request payload: invalid character '\\\\n' in string literal\"}")
	})

	t.Run("it returns an error when profile type is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		payload := `{
			"username": "test-user",
			"password": "Pass123!",
			"firstName": "Bruce", 
			"lastName": "Lee",
			"gender": "M", 
			"phone": "123456789",
			"email": "test-user@example.com", 
			"profileType": "blabla"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := NewEmployeeHandler(log, nil, nil)
		handler.RegisterEmployee(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid profile type")
	})

	t.Run("it returns an error when it fails to list employees by username", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		payload := `{
			"username": "test-user",
			"password": "Pass123!",
			"firstName": "Bruce", 
			"lastName": "Lee",
			"gender": "M", 
			"phone": "123456789",
			"email": "test-user@example.com", 
			"profileType": "Medic"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		emplRepoMock := repositories.NewMockEmployeeRepository(gomock.NewController(t))
		emplRepoMock.EXPECT().ListEmployees(gomock.Any()).Return(nil, gorm.ErrRecordNotFound).Times(1)

		handler := NewEmployeeHandler(log, emplRepoMock, nil)
		handler.RegisterEmployee(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "{\"error\":\"Failed to check for existing username\"}")
	})

	t.Run("it returns an error when employee with same username already exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		payload := `{
			"username": "test-user",
			"password": "Pass123!",
			"firstName": "Bruce", 
			"lastName": "Lee",
			"gender": "M", 
			"phone": "123456789",
			"email": "test-user@example.com", 
			"profileType": "Medic"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		existingEmployee := model.Employee{Username: "test-user"}

		emplRepoMock := repositories.NewMockEmployeeRepository(gomock.NewController(t))
		emplRepoMock.EXPECT().ListEmployees(gomock.Any()).Return([]model.Employee{existingEmployee}, nil).Times(1)

		handler := NewEmployeeHandler(log, emplRepoMock, nil)
		handler.RegisterEmployee(ctx)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "Username already exists")
	})

	t.Run("it returns an error when it fails to list employees by email", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		payload := `{
			"username": "test-user",
			"password": "Pass123!",
			"firstName": "Bruce", 
			"lastName": "Lee",
			"gender": "M", 
			"phone": "123456789",
			"email": "test-user@example.com", 
			"profileType": "Medic"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		emplRepoMock := repositories.NewMockEmployeeRepository(gomock.NewController(t))
		usernameFilter := map[string]interface{}{
			"username": "test-user",
		}
		emplRepoMock.EXPECT().ListEmployees(usernameFilter).Return([]model.Employee{}, nil).Times(1)
		emplRepoMock.EXPECT().ListEmployees(gomock.Any()).Return(nil, gorm.ErrRecordNotFound).Times(1)

		handler := NewEmployeeHandler(log, emplRepoMock, nil)
		handler.RegisterEmployee(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "{\"error\":\"Failed to check for existing email\"}")
	})

	t.Run("it returns an error when employee with same email already exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		payload := `{
			"username": "test-user",
			"password": "Pass123!",
			"firstName": "Bruce", 
			"lastName": "Lee",
			"gender": "M", 
			"phone": "123456789",
			"email": "test-user@example.com", 
			"profileType": "Medic"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		existingEmployee := model.Employee{Email: "test-user@example.com"}

		emplRepoMock := repositories.NewMockEmployeeRepository(gomock.NewController(t))
		usernameFilter := map[string]interface{}{
			"username": "test-user",
		}
		emplRepoMock.EXPECT().ListEmployees(usernameFilter).Return([]model.Employee{}, nil).Times(1)
		emplRepoMock.EXPECT().ListEmployees(gomock.Any()).Return([]model.Employee{existingEmployee}, nil).Times(1)

		handler := NewEmployeeHandler(log, emplRepoMock, nil)
		handler.RegisterEmployee(ctx)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "Email already exists")
	})

	tests := []struct {
		name     string
		password string
		error    string
	}{
		{
			name:     "it returns an error when password is too short",
			password: "short",
			error:    utils.ErrPasswordLength,
		},
		{
			name:     "it returns an error when password is too long",
			password: "verylongpasswordthatexceedstheallowedlength",
			error:    utils.ErrPasswordLength,
		},
		{
			name:     "it returns an error when password Lees not contain an uppercase letter",
			password: "pass123!",
			error:    utils.ErrPasswordUppercase,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			payload := `{
				"username": "test-user",
				"password": "` + test.password + `",
				"firstName": "Bruce", 
				"lastName": "Lee",
				"gender": "M", 
				"phone": "123456789",
				"email": "test-user@example.com", 
				"profileType": "Medic"
			}`
			ctx.Request = httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(payload))
			ctx.Request.Header.Set("Content-Type", "application/json")

			usernameFilter := map[string]interface{}{
				"username": "test-user",
			}
			emailFilter := map[string]interface{}{
				"email": "test-user@example.com",
			}
			emplRepoMock := repositories.NewMockEmployeeRepository(gomock.NewController(t))
			emplRepoMock.EXPECT().ListEmployees(usernameFilter).Return([]model.Employee{}, nil).Times(1)
			emplRepoMock.EXPECT().ListEmployees(emailFilter).Return([]model.Employee{}, nil).Times(1)

			handler := NewEmployeeHandler(log, emplRepoMock, nil)
			handler.RegisterEmployee(ctx)

			if test.error != "" {
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Contains(t, w.Body.String(), test.error)
			} else {
				assert.Equal(t, http.StatusCreated, w.Code)
			}
		})
	}

	t.Run("it returns an error when it fails to create employee", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		payload := `{
			"username": "test-user",
			"password": "Pass123!",
			"firstName": "Bruce", 
			"lastName": "Lee",
			"gender": "M", 
			"phone": "123456789",
			"email": "test-user@example.com", 
			"profileType": "Medic"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		usernameFilter := map[string]interface{}{
			"username": "test-user",
		}
		emailFilter := map[string]interface{}{
			"email": "test-user@example.com",
		}
		emplRepoMock := repositories.NewMockEmployeeRepository(gomock.NewController(t))
		emplRepoMock.EXPECT().ListEmployees(usernameFilter).Return([]model.Employee{}, nil).Times(1)
		emplRepoMock.EXPECT().ListEmployees(emailFilter).Return([]model.Employee{}, nil).Times(1)
		emplRepoMock.EXPECT().Create(gomock.Any()).Return(gorm.ErrInvalidDB).Times(1)

		handler := NewEmployeeHandler(log, emplRepoMock, nil)
		handler.RegisterEmployee(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "{\"error\":\"invalid db\"}")
	})

	t.Run("it creates an employee when data is valid", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		payload := `{
			"username": "test-user",
			"password": "Pass123!",
			"firstName": "Bruce", 
			"lastName": "Lee",
			"gender": "M", 
			"phone": "123456789",
			"email": "test-user@example.com", 
			"profileType": "Medic"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		usernameFilter := map[string]interface{}{
			"username": "test-user",
		}
		emailFilter := map[string]interface{}{
			"email": "test-user@example.com",
		}
		emplRepoMock := repositories.NewMockEmployeeRepository(gomock.NewController(t))
		emplRepoMock.EXPECT().ListEmployees(usernameFilter).Return([]model.Employee{}, nil).Times(1)
		emplRepoMock.EXPECT().ListEmployees(emailFilter).Return([]model.Employee{}, nil).Times(1)
		emplRepoMock.EXPECT().Create(gomock.Any()).Return(nil).Times(1)

		handler := NewEmployeeHandler(log, emplRepoMock, nil)
		handler.RegisterEmployee(ctx)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "{\"id\":0,\"username\":\"test-user\",\"firstName\":\"Bruce\",\"lastName\":\"Lee\",\"gender\":\"M\",\"phone\":\"123456789\",\"email\":\"test-user@example.com\",\"profilePicture\":\"\",\"profileType\":\"Medic\"}")
	})
}

func TestEmployeeHandler_LoginEmployee(t *testing.T) {
	log := utils.NewTestLogger()

	t.Run("it returns an error when request payload is invalid json", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		invalidPayload := `{
			"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(invalidPayload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := NewEmployeeHandler(log, nil, nil)
		handler.LoginEmployee(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "{\"error\":\"Invalid request payload: invalid character '\\\\n' in string literal\"}")
	})

	t.Run("it returns an error when employee Lees not exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		payload := `{
			"username": "test-user",
			"password": "Pass123!"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		emplRepoMock := repositories.NewMockEmployeeRepository(gomock.NewController(t))
		emplRepoMock.EXPECT().GetEmployeeByUsername("test-user").Return(nil, gorm.ErrRecordNotFound).Times(1)

		handler := NewEmployeeHandler(log, emplRepoMock, nil)
		handler.LoginEmployee(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid credentials")
	})

	t.Run("it returns an error when password is incorrect", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		payload := `{
			"username": "test-user",
			"password": "WrongPassword!"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		emplRepoMock := repositories.NewMockEmployeeRepository(gomock.NewController(t))
		emplRepoMock.EXPECT().GetEmployeeByUsername("test-user").Return(&model.Employee{Password: "Pass123!"}, nil).Times(1)

		handler := NewEmployeeHandler(log, emplRepoMock, nil)
		handler.LoginEmployee(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid credentials")
	})

	t.Run("it returns a JWT token when login is successful", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		payload := `{
			"username": "test-user",
			"password": "Pass123!"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		emplRepoMock := repositories.NewMockEmployeeRepository(gomock.NewController(t))
		emplRepoMock.EXPECT().GetEmployeeByUsername("test-user").Return(&model.Employee{Username: "test-user", Password: "$2a$10$wq8KS0Dy7tGWM5pnCqPhfO.uY1vvVzZb5.CWsqqCyEQv89Uu6QDaK"}, nil).Times(1)

		handler := NewEmployeeHandler(log, emplRepoMock, nil)
		handler.LoginEmployee(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "token")
	})

}

func TestEmployeeHandler_ListEmployees(t *testing.T) {
	log := utils.NewTestLogger()

	t.Run("it returns an error when it fails to retrieve employees", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		emplRepoMock := repositories.NewMockEmployeeRepository(gomock.NewController(t))
		emplRepoMock.EXPECT().GetAll().Return(nil, gorm.ErrRecordNotFound).Times(1)

		handler := NewEmployeeHandler(log, emplRepoMock, nil)
		handler.ListEmployees(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "{\"error\":\"record not found\"}")
	})

	t.Run("it returns an empty list when no employees exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		emplRepoMock := repositories.NewMockEmployeeRepository(gomock.NewController(t))
		emplRepoMock.EXPECT().GetAll().Return([]model.Employee{}, nil).Times(1)

		handler := NewEmployeeHandler(log, emplRepoMock, nil)
		handler.ListEmployees(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "[]")
	})

	t.Run("it returns a list of employees when employees exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		employees := []model.Employee{
			{Username: "test-user", FirstName: "Bruce", LastName: "Lee", Password: "Pass123!"},
			{Username: "asmith", FirstName: "Alice", LastName: "Smith", Password: "Pass123!"},
		}

		emplRepoMock := repositories.NewMockEmployeeRepository(gomock.NewController(t))
		emplRepoMock.EXPECT().GetAll().Return(employees, nil).Times(1)

		handler := NewEmployeeHandler(log, emplRepoMock, nil)
		handler.ListEmployees(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "[{\"id\":0,\"username\":\"test-user\",\"firstName\":\"Bruce\",\"lastName\":\"Lee\",\"gender\":\"\",\"phone\":\"\",\"email\":\"\",\"profilePicture\":\"\",\"profileType\":\"Unknown\"},{\"id\":0,\"username\":\"asmith\",\"firstName\":\"Alice\",\"lastName\":\"Smith\",\"gender\":\"\",\"phone\":\"\",\"email\":\"\",\"profilePicture\":\"\",\"profileType\":\"Unknown\"}]")
	})

}

func TestEmployeeHandler_UpdateEmployee(t *testing.T) {
	log := utils.NewTestLogger()

	t.Run("it returns an error when employee Lees not exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		emplRepoMock := repositories.NewMockEmployeeRepository(gomock.NewController(t))
		emplRepoMock.EXPECT().GetEmployeeByID("1", gomock.Any()).Return(gorm.ErrRecordNotFound).Times(1)

		handler := NewEmployeeHandler(log, emplRepoMock, nil)
		handler.UpdateEmployee(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "{\"error\":\"Employee not found\"}")
	})

	t.Run("it returns an error when request payload is invalid json", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		invalidPayload := `{
			"
		}`

		emplRepoMock := repositories.NewMockEmployeeRepository(gomock.NewController(t))
		emplRepoMock.EXPECT().GetEmployeeByID("1", gomock.Any()).Return(nil).Times(1)

		ctx.Request = httptest.NewRequest(http.MethodPut, "/employees/1", strings.NewReader(invalidPayload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := NewEmployeeHandler(log, emplRepoMock, nil)
		handler.UpdateEmployee(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "{\"error\":\"Invalid request payload\"}")
	})

	t.Run("it returns an error when validation fails", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		emplRepoMock := repositories.NewMockEmployeeRepository(gomock.NewController(t))
		emplRepoMock.EXPECT().GetEmployeeByID("1", gomock.Any()).Return(nil).Times(1)

		payload := `{
			"firstName": "B",
			"lastName": "L",
			"age": 10,
			"email": "invalid-email.com"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPut, "/employees/1", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := NewEmployeeHandler(log, emplRepoMock, nil)
		handler.UpdateEmployee(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "{\"error\":\"mail: missing '@' or angle-addr\"}")
	})

	t.Run("it returns an error when it fails to update an existing employee", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		emplRepoMock := repositories.NewMockEmployeeRepository(gomock.NewController(t))
		emplRepoMock.EXPECT().GetEmployeeByID("1", gomock.Any()).Return(nil).Times(1)

		payload := `{
			"firstName": "Bruce",
			"lastName": "Lee",
			"age": 30,
			"email": "test-user@example.com"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPut, "/employees/1", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		emplRepoMock.EXPECT().UpdateEmployee(gomock.Any()).Return(gorm.ErrRecordNotFound).Times(1)

		handler := NewEmployeeHandler(log, emplRepoMock, nil)
		handler.UpdateEmployee(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "{\"error\":\"Failed to update employee\"}")
	})

	t.Run("it successfully updates an existing employee when it exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		emplRepoMock := repositories.NewMockEmployeeRepository(gomock.NewController(t))
		emplRepoMock.EXPECT().GetEmployeeByID("1", gomock.Any()).Return(nil).Times(1)

		payload := `{
			"firstName": "Bruce",
			"lastName": "Lee",
			"age": 30,
			"email": "test-user@example.com"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPut, "/employees/1", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		emplRepoMock.EXPECT().UpdateEmployee(gomock.Any()).Return(nil).Times(1)

		handler := NewEmployeeHandler(log, emplRepoMock, nil)
		handler.UpdateEmployee(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "{\"id\":0,\"username\":\"\",\"firstName\":\"Bruce\",\"lastName\":\"Lee\",\"gender\":\"\",\"phone\":\"\",\"email\":\"test-user@example.com\",\"profilePicture\":\"\",\"profileType\":\"Unknown\"}")
	})

}

func TestEmployeeHandler_CreateEmployee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmplRepo := repositories.NewMockEmployeeRepository(ctrl)
	mockShiftRepo := repositories.NewMockShiftRepository(ctrl)
	log := utils.NewTestLogger()
	handler := NewEmployeeHandler(log, mockEmplRepo, mockShiftRepo)

	gin.SetMode(gin.TestMode)

	t.Run("it returns an error when password validation fails", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		usernameFilter := map[string]interface{}{
			"username": "test-user",
		}
		emailFilter := map[string]interface{}{
			"email": "test-user@example.com",
		}
		mockEmplRepo.EXPECT().ListEmployees(usernameFilter).Return([]model.Employee{}, nil).Times(1)
		mockEmplRepo.EXPECT().ListEmployees(emailFilter).Return([]model.Employee{}, nil).Times(1)
		mockEmplRepo.EXPECT().Create(gomock.Any()).Return(nil).Times(0)
		invalidEmployee := `{
			"username": "test-user",
			"password": "short", 
			"firstName": "Bruce", 
			"lastName": "Lee",
			"gender": "M", 
			"phone": "123456789",
			"email": "test-user@example.com", 
			"profileType": "Medic"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(invalidEmployee))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.RegisterEmployee(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), utils.ErrPasswordLength)
	})

	t.Run("it returns an error when employee already exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		usernameFilter := map[string]interface{}{
			"username": "test-user",
		}
		emailFilter := map[string]interface{}{
			"email": "test-user@example.com",
		}
		mockEmplRepo.EXPECT().ListEmployees(usernameFilter).Return([]model.Employee{{}}, nil).Times(1)
		mockEmplRepo.EXPECT().ListEmployees(emailFilter).Return([]model.Employee{}, nil).Times(0)
		mockEmplRepo.EXPECT().Create(gomock.Any()).Return(nil).Times(0)
		existingEmployee := `{
			"username": "test-user",
			"password": "Pass123!",
			"firstName": "Bruce", 
			"lastName": "Lee",
			"gender": "M", 
			"phone": "123456789",
			"email": "test-user@example.com", 
			"profileType": "Medic"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(existingEmployee))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.RegisterEmployee(ctx)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "Username already exists")
	})

	t.Run("it returns an error when email already exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		usernameFilter := map[string]interface{}{
			"username": "test-user",
		}
		emailFilter := map[string]interface{}{
			"email": "test-user@example.com",
		}
		mockEmplRepo.EXPECT().ListEmployees(usernameFilter).Return([]model.Employee{}, nil).Times(1)
		mockEmplRepo.EXPECT().ListEmployees(emailFilter).Return([]model.Employee{{}}, nil).Times(1)
		mockEmplRepo.EXPECT().Create(gomock.Any()).Return(nil).Times(0)
		existingEmployee := `{
			"username": "test-user",
			"password": "Pass123!",
			"firstName": "Bruce", 
			"lastName": "Lee",
			"gender": "M", 
			"phone": "123456789",
			"email": "test-user@example.com", 
			"profileType": "Medic"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(existingEmployee))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.RegisterEmployee(ctx)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "Email already exists")
	})

	t.Run("it creates an employee when data is valid", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		usernameFilter := map[string]interface{}{
			"username": "test-user",
		}
		emailFilter := map[string]interface{}{
			"email": "test-user@example.com",
		}
		mockEmplRepo.EXPECT().ListEmployees(usernameFilter).Return([]model.Employee{}, nil).Times(1)
		mockEmplRepo.EXPECT().ListEmployees(emailFilter).Return([]model.Employee{}, nil).Times(1)
		mockEmplRepo.EXPECT().Create(gomock.Any()).Return(nil).Times(1)
		validEmployee := `{
			"username": "test-user",
			"password": "Pass123!",
			"firstName": "Bruce", 
			"lastName": "Lee",
			"gender": "M", 
			"phone": "123456789",
			"email": "test-user@example.com", 
			"profileType": "Medic"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(validEmployee))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.RegisterEmployee(ctx)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, `{"id":0,"username":"test-user","firstName":"Bruce","lastName":"Lee","gender":"M","phone":"123456789","email":"test-user@example.com","profilePicture":"","profileType":"Medic"}`, w.Body.String())
	})
}

func TestEmployeeHandler_GetAllEmployees(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositories.NewMockEmployeeRepository(ctrl)
	mockShiftRepo := repositories.NewMockShiftRepository(ctrl)
	log := utils.NewTestLogger()
	handler := NewEmployeeHandler(log, mockRepo, mockShiftRepo)

	gin.SetMode(gin.TestMode)

	t.Run("it returns an empty list when no employees exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		mockRepo.EXPECT().GetAll().Return([]model.Employee{}, nil).Times(1)

		handler.ListEmployees(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `[]`, w.Body.String())
	})

	t.Run("it returns a list of employees when employees exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		employees := []model.Employee{
			{Username: "test-user", FirstName: "Bruce", LastName: "Lee", Password: "Pass123!"},
			{Username: "asmith", FirstName: "Alice", LastName: "Smith", Password: "Pass123!"},
		}

		mockRepo.EXPECT().GetAll().Return(employees, nil).Times(1)

		handler.ListEmployees(ctx)

		expectedJSON := `[
			{"id":0,"username":"test-user","firstName":"Bruce","lastName":"Lee","gender":"","phone":"","email":"","profilePicture":"","profileType":"Unknown"},
			{"id":0,"username":"asmith","firstName":"Alice","lastName":"Smith","gender":"","phone":"","email":"","profilePicture":"","profileType":"Unknown"}
		]`

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expectedJSON, w.Body.String())
	})
}

func TestEmployeeHandler_DeleteEmployee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositories.NewMockEmployeeRepository(ctrl)
	mockShiftRepo := repositories.NewMockShiftRepository(ctrl)
	log := utils.NewTestLogger()
	handler := NewEmployeeHandler(log, mockRepo, mockShiftRepo)

	gin.SetMode(gin.TestMode)

	t.Run("it returns an error when employee Lees not exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		mockRepo.EXPECT().Delete(uint(1)).Return(gorm.ErrRecordNotFound).Times(1)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}
		handler.DeleteEmployee(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to delete employee")
	})

	t.Run("it deletes an existing employee", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		mockRepo.EXPECT().Delete(uint(1)).Return(nil).Times(1)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}
		handler.DeleteEmployee(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Employee deleted successfully")
	})
}
