package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock_repositories "github.com/pd120424d/mountain-service/api/activity/internal/repositories/mocks"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

func TestActivityService_CreateActivity(t *testing.T) {
	t.Parallel()

	log, err := utils.NewLogger("activity-test")
	assert.NoError(t, err)

	t.Run("successfully creates activity", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mock_repositories.NewMockActivityRepository(mockCtrl)

		svc := NewActivityService(log, mockRepo)

		assert.NotNil(t, svc)
	})
}

func TestActivityHandler_NewActivityHandler(t *testing.T) {
	t.Parallel()

	log, err := utils.NewLogger("activity-test")
	assert.NoError(t, err)

	t.Run("successfully creates handler", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := mock_repositories.NewMockActivityRepository(mockCtrl)
		svc := NewActivityService(log, mockRepo)
		handler := NewActivityHandler(log, svc)

		assert.NotNil(t, handler)
	})
}
