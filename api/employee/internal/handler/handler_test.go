package handler

import (
	"api/employee/internal/repositories"
	"api/shared/utils"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"api/employee/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestEmployeeHandler_CreateEmployee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmplRepo := repositories.NewMockEmployeeRepository(ctrl)
	mockShiftRepo := repositories.NewMockShiftRepository(ctrl)
	log := utils.NewLogger() // Use a mocked or real logger depending on the situation
	handler := NewEmployeeHandler(log, mockEmplRepo, mockShiftRepo)

	gin.SetMode(gin.TestMode)

	t.Run("it returns an error when password validation fails", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		usernameFilter := map[string]interface{}{
			"username": "jdoe",
		}
		emailFilter := map[string]interface{}{
			"email": "jdoe@example.com",
		}
		mockEmplRepo.EXPECT().ListEmployees(usernameFilter).Return([]model.Employee{}, nil).Times(1)
		mockEmplRepo.EXPECT().ListEmployees(emailFilter).Return([]model.Employee{}, nil).Times(1)
		mockEmplRepo.EXPECT().Create(gomock.Any()).Return(nil).Times(0)
		invalidEmployee := `{
			"username": "jdoe",
			"password": "short", 
			"firstName": "John", 
			"lastName": "Doe",
			"gender": "M", 
			"phone": "123456789",
			"email": "jdoe@example.com", 
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
			"username": "jdoe",
		}
		emailFilter := map[string]interface{}{
			"email": "jdoe@example.com",
		}
		mockEmplRepo.EXPECT().ListEmployees(usernameFilter).Return([]model.Employee{{}}, nil).Times(1)
		mockEmplRepo.EXPECT().ListEmployees(emailFilter).Return([]model.Employee{}, nil).Times(0)
		mockEmplRepo.EXPECT().Create(gomock.Any()).Return(nil).Times(0)
		existingEmployee := `{
			"username": "jdoe",
			"password": "Pass123!",
			"firstName": "John", 
			"lastName": "Doe",
			"gender": "M", 
			"phone": "123456789",
			"email": "jdoe@example.com", 
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
			"username": "jdoe",
		}
		emailFilter := map[string]interface{}{
			"email": "jdoe@example.com",
		}
		mockEmplRepo.EXPECT().ListEmployees(usernameFilter).Return([]model.Employee{}, nil).Times(1)
		mockEmplRepo.EXPECT().ListEmployees(emailFilter).Return([]model.Employee{{}}, nil).Times(1)
		mockEmplRepo.EXPECT().Create(gomock.Any()).Return(nil).Times(0)
		existingEmployee := `{
			"username": "jdoe",
			"password": "Pass123!",
			"firstName": "John", 
			"lastName": "Doe",
			"gender": "M", 
			"phone": "123456789",
			"email": "jdoe@example.com", 
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
			"username": "jdoe",
		}
		emailFilter := map[string]interface{}{
			"email": "jdoe@example.com",
		}
		mockEmplRepo.EXPECT().ListEmployees(usernameFilter).Return([]model.Employee{}, nil).Times(1)
		mockEmplRepo.EXPECT().ListEmployees(emailFilter).Return([]model.Employee{}, nil).Times(1)
		mockEmplRepo.EXPECT().Create(gomock.Any()).Return(nil).Times(1)
		validEmployee := `{
			"username": "jdoe",
			"password": "Pass123!",
			"firstName": "John", 
			"lastName": "Doe",
			"gender": "M", 
			"phone": "123456789",
			"email": "jdoe@example.com", 
			"profileType": "Medic"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(validEmployee))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.RegisterEmployee(ctx)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, `{"id":0,"username":"jdoe","firstName":"John","lastName":"Doe","gender":"M","phoneNumber":"123456789","email":"jdoe@example.com","profilePicture":"","profileType":"Medic"}`, w.Body.String())
	})
}

func TestEmployeeHandler_GetAllEmployees(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositories.NewMockEmployeeRepository(ctrl)
	mockShiftRepo := repositories.NewMockShiftRepository(ctrl)
	log := utils.NewLogger()
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
			{Username: "jdoe", FirstName: "John", LastName: "Doe", Password: "Pass123!"},
			{Username: "asmith", FirstName: "Alice", LastName: "Smith", Password: "Pass123!"},
		}

		mockRepo.EXPECT().GetAll().Return(employees, nil).Times(1)

		handler.ListEmployees(ctx)

		expectedJSON := `[
            {"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"Username":"jdoe","Password":"Pass123!","FirstName":"John","LastName":"Doe","Gender":"","Phone":"","Email":"","ProfilePicture":"","ProfileType":""},
            {"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"Username":"asmith","Password":"Pass123!","FirstName":"Alice","LastName":"Smith","Gender":"","Phone":"","Email":"","ProfilePicture":"","ProfileType":""}
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
	log := utils.NewLogger()
	handler := NewEmployeeHandler(log, mockRepo, mockShiftRepo)

	gin.SetMode(gin.TestMode)

	t.Run("it returns an error when employee does not exist", func(t *testing.T) {
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
