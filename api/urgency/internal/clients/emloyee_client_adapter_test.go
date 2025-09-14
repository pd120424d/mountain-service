package clients

import (
	"testing"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	s2semployee "github.com/pd120424d/mountain-service/api/shared/s2s/employee"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestNewEmployeeClientFromS2S(t *testing.T) {
	t.Parallel()

	t.Run("it calls inner client GetEmployeeByID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := utils.NewTestLogger()
		inner := s2semployee.NewMockClient(ctrl)
		client := NewEmployeeClientFromS2S(inner, logger)

		inner.EXPECT().GetEmployeeByID(gomock.Any(), uint(1)).Return(&employeeV1.EmployeeResponse{ID: 1}, nil)
		_, err := client.GetEmployeeByID(t.Context(), 1)
		assert.NoError(t, err)
	})

	t.Run("it calls inner client GetAllEmployees", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := utils.NewTestLogger()
		inner := s2semployee.NewMockClient(ctrl)
		client := NewEmployeeClientFromS2S(inner, logger)

		inner.EXPECT().GetAllEmployees(gomock.Any()).Return([]employeeV1.EmployeeResponse{{ID: 1}}, nil)
		_, err := client.GetAllEmployees(t.Context())
		assert.NoError(t, err)
	})

	t.Run("it calls inner client GetOnCallEmployees", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := utils.NewTestLogger()
		inner := s2semployee.NewMockClient(ctrl)
		client := NewEmployeeClientFromS2S(inner, logger)

		inner.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return([]employeeV1.EmployeeResponse{{ID: 1}}, nil)
		_, err := client.GetOnCallEmployees(t.Context(), 0)
		assert.NoError(t, err)
	})

	t.Run("it calls inner client CheckActiveEmergencies", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := utils.NewTestLogger()
		inner := s2semployee.NewMockClient(ctrl)
		client := NewEmployeeClientFromS2S(inner, logger)

		inner.EXPECT().CheckActiveEmergencies(gomock.Any(), uint(1)).Return(true, nil)
		_, err := client.CheckActiveEmergencies(t.Context(), 1)
		assert.NoError(t, err)
	})

}
