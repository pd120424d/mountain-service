package handler

import (
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mountain-service/employee/internal/model"
	"mountain-service/employee/internal/repositories"
	"mountain-service/shared/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestEmployeeHandler_CreateEmployee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositories.NewMockEmployeeRepository(ctrl)
	log := utils.NewLogger() // Use a mocked or real logger depending on the situation
	handler := NewEmployeeHandler(log, mockRepo)

	gin.SetMode(gin.TestMode)

	t.Run("it returns an error when password validation fails", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		invalidEmployee := `{"username": "jdoe", "password": "short", "first_name": "John", "last_name": "Doe"}`
		c.Request = httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(invalidEmployee))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateEmployee(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), utils.ErrPasswordLength)
	})

	t.Run("it creates an employee when data is valid", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		validEmployee := model.Employee{
			Username:  "jdoe",
			Password:  "Pass123!",
			FirstName: "John",
			LastName:  "Doe",
		}

		mockRepo.EXPECT().Create(&validEmployee).Return(nil).Times(1)

		c.Request = httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(`{
            "username": "jdoe",
            "password": "Pass123!",
            "firstName": "John",
            "lastName": "Doe"
        }`))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateEmployee(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var createdEmployee model.Employee
		err := json.Unmarshal(w.Body.Bytes(), &createdEmployee)
		assert.Nil(t, err)
		assert.Equal(t, "jdoe", createdEmployee.Username)
		assert.Equal(t, "Pass123!", createdEmployee.Password)
		assert.Equal(t, "John", createdEmployee.FirstName)
		assert.Equal(t, "Doe", createdEmployee.LastName)
	})
}

func TestEmployeeHandler_GetAllEmployees(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositories.NewMockEmployeeRepository(ctrl)
	log := utils.NewLogger()
	handler := NewEmployeeHandler(log, mockRepo)

	gin.SetMode(gin.TestMode)

	t.Run("it returns an empty list when no employees exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockRepo.EXPECT().GetAll().Return([]model.Employee{}, nil).Times(1)

		handler.ListEmployees(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `[]`, w.Body.String())
	})

	t.Run("it returns a list of employees when employees exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		employees := []model.Employee{
			{Username: "jdoe", FirstName: "John", LastName: "Doe", Password: "Pass123!"},
			{Username: "asmith", FirstName: "Alice", LastName: "Smith", Password: "Pass123!"},
		}

		mockRepo.EXPECT().GetAll().Return(employees, nil).Times(1)

		handler.ListEmployees(c)

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
	log := utils.NewLogger()
	handler := NewEmployeeHandler(log, mockRepo)

	gin.SetMode(gin.TestMode)

	t.Run("it returns an error when employee does not exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockRepo.EXPECT().Delete(uint(1)).Return(gorm.ErrRecordNotFound).Times(1)

		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		handler.DeleteEmployee(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to delete employee")
	})

	t.Run("it deletes an existing employee", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockRepo.EXPECT().Delete(uint(1)).Return(nil).Times(1)

		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		handler.DeleteEmployee(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Employee deleted successfully")
	})
}
