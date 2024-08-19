package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mountain-service/employee/internal/models"
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

		validEmployee := models.Employee{
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

		var createdEmployee models.Employee
		err := json.Unmarshal(w.Body.Bytes(), &createdEmployee)
		assert.Nil(t, err)
		assert.Equal(t, "jdoe", createdEmployee.Username)
		assert.Equal(t, "Pass123!", createdEmployee.Password)
		assert.Equal(t, "John", createdEmployee.FirstName)
		assert.Equal(t, "Doe", createdEmployee.LastName)
	})
}
