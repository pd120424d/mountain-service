package repositories

import (
	"testing"

	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/pd120424d/mountain-service/api/urgency/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Urgency{})
	require.NoError(t, err)

	return db
}

func TestUrgencyRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	log := utils.NewTestLogger()
	repo := NewUrgencyRepository(log, db)

	urgency := &model.Urgency{
		Name:         "Test Urgency",
		Email:        "test@example.com",
		ContactPhone: "123456789",
		Description:  "Test description",
		Level:        model.High,
		Status:       model.Open,
	}

	err := repo.Create(urgency)
	assert.NoError(t, err)
	assert.NotZero(t, urgency.ID)
}

func TestUrgencyRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	log := utils.NewTestLogger()
	repo := NewUrgencyRepository(log, db)

	// Create test urgencies
	urgency1 := &model.Urgency{
		Name:         "Test Urgency 1",
		Email:        "test1@example.com",
		ContactPhone: "123456789",
		Description:  "Test description 1",
		Level:        model.High,
		Status:       "Open",
	}
	urgency2 := &model.Urgency{
		Name:         "Test Urgency 2",
		Email:        "test2@example.com",
		ContactPhone: "987654321",
		Description:  "Test description 2",
		Level:        model.Medium,
		Status:       "In Progress",
	}

	err := repo.Create(urgency1)
	require.NoError(t, err)
	err = repo.Create(urgency2)
	require.NoError(t, err)

	urgencies, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, urgencies, 2)
}

func TestUrgencyRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	log := utils.NewTestLogger()
	repo := NewUrgencyRepository(log, db)

	urgency := &model.Urgency{
		Name:         "Test Urgency",
		Email:        "test@example.com",
		ContactPhone: "123456789",
		Description:  "Test description",
		Level:        model.High,
		Status:       "Open",
	}

	err := repo.Create(urgency)
	require.NoError(t, err)

	var retrieved model.Urgency
	err = repo.GetByID(urgency.ID, &retrieved)
	assert.NoError(t, err)
	assert.Equal(t, urgency.Name, retrieved.Name)
	assert.Equal(t, urgency.Email, retrieved.Email)
}

func TestUrgencyRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	log := utils.NewTestLogger()
	repo := NewUrgencyRepository(log, db)

	urgency := &model.Urgency{
		Name:         "Test Urgency",
		Email:        "test@example.com",
		ContactPhone: "123456789",
		Description:  "Test description",
		Level:        model.High,
		Status:       model.Open,
	}

	err := repo.Create(urgency)
	require.NoError(t, err)

	urgency.Name = "Updated Urgency"
	urgency.Status = model.InProgress

	err = repo.Update(urgency)
	assert.NoError(t, err)

	var updated model.Urgency
	err = repo.GetByID(urgency.ID, &updated)
	require.NoError(t, err)
	assert.Equal(t, "Updated Urgency", updated.Name)
	assert.Equal(t, model.InProgress, updated.Status)
}

func TestUrgencyRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	log := utils.NewTestLogger()
	repo := NewUrgencyRepository(log, db)

	urgency := &model.Urgency{
		Name:         "Test Urgency",
		Email:        "test@example.com",
		ContactPhone: "123456789",
		Description:  "Test description",
		Level:        model.High,
		Status:       "Open",
	}

	err := repo.Create(urgency)
	require.NoError(t, err)

	err = repo.Delete(urgency.ID)
	assert.NoError(t, err)

	var deleted model.Urgency
	err = repo.GetByID(urgency.ID, &deleted)
	assert.Error(t, err) // Should not find deleted record
}

func TestUrgencyRepository_List(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	log := utils.NewTestLogger()
	repo := NewUrgencyRepository(log, db)

	urgency1 := &model.Urgency{
		Name:         "Emergency Rescue",
		Email:        "rescue@example.com",
		ContactPhone: "123456789",
		Description:  "Mountain rescue needed",
		Level:        model.Critical,
		Status:       model.Open,
	}
	urgency2 := &model.Urgency{
		Name:         "Medical Help",
		Email:        "medical@example.com",
		ContactPhone: "987654321",
		Description:  "Medical assistance required",
		Level:        model.High,
		Status:       model.InProgress,
	}
	urgency3 := &model.Urgency{
		Name:         "Equipment Issue",
		Email:        "equipment@example.com",
		ContactPhone: "555666777",
		Description:  "Equipment malfunction",
		Level:        model.Medium,
		Status:       model.Resolved,
	}

	require.NoError(t, repo.Create(urgency1))
	require.NoError(t, repo.Create(urgency2))
	require.NoError(t, repo.Create(urgency3))

	t.Run("it returns all urgencies when no filters are provided", func(t *testing.T) {
		urgencies, err := repo.List(map[string]interface{}{})
		assert.NoError(t, err)
		assert.Len(t, urgencies, 3)
	})

	t.Run("it returns urgencies filtered by level", func(t *testing.T) {
		filters := map[string]interface{}{
			"level": "High",
		}
		urgencies, err := repo.List(filters)
		assert.NoError(t, err)
		assert.Len(t, urgencies, 1)
		assert.Equal(t, "Medical Help", urgencies[0].Name)
	})

	t.Run("it returns urgencies filtered by status", func(t *testing.T) {
		filters := map[string]interface{}{
			"status": "Open",
		}
		urgencies, err := repo.List(filters)
		assert.NoError(t, err)
		assert.Len(t, urgencies, 1)
		assert.Equal(t, "Emergency Rescue", urgencies[0].Name)
	})

	t.Run("it returns urgencies filtered by name (partial match)", func(t *testing.T) {
		filters := map[string]interface{}{
			"name": "Medical",
		}
		urgencies, err := repo.List(filters)
		assert.NoError(t, err)
		assert.Len(t, urgencies, 1)
		assert.Equal(t, "Medical Help", urgencies[0].Name)
	})

	t.Run("it returns urgencies filtered by email (partial match)", func(t *testing.T) {
		filters := map[string]interface{}{
			"email": "rescue",
		}
		urgencies, err := repo.List(filters)
		assert.NoError(t, err)
		assert.Len(t, urgencies, 1)
		assert.Equal(t, "Emergency Rescue", urgencies[0].Name)
	})

	t.Run("it returns urgencies filtered by multiple filters", func(t *testing.T) {
		filters := map[string]interface{}{
			"level":  "Medium",
			"status": "Resolved",
		}
		urgencies, err := repo.List(filters)
		assert.NoError(t, err)
		assert.Len(t, urgencies, 1)
		assert.Equal(t, "Equipment Issue", urgencies[0].Name)
	})

	t.Run("it returns an error when an invalid filter key is provided", func(t *testing.T) {
		filters := map[string]interface{}{
			"invalid_field": "value",
		}
		urgencies, err := repo.List(filters)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid filter key: invalid_field")
		assert.Nil(t, urgencies)
	})

	t.Run("it returns an error when an unsupported filter value type is provided", func(t *testing.T) {
		filters := map[string]interface{}{
			"name": []string{"test"},
		}
		urgencies, err := repo.List(filters)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported type for filter key: name")
		assert.Nil(t, urgencies)
	})

	t.Run("it returns an empty slice when no matches are found", func(t *testing.T) {
		filters := map[string]interface{}{
			"name": "NonExistent",
		}
		urgencies, err := repo.List(filters)
		assert.NoError(t, err)
		assert.Len(t, urgencies, 0)
	})
}

func TestUrgencyRepository_ResetAllData(t *testing.T) {
	db := setupTestDB(t)
	log := utils.NewTestLogger()
	repo := NewUrgencyRepository(log, db)

	urgency1 := &model.Urgency{
		Name:         "Test Urgency 1",
		Email:        "test1@example.com",
		ContactPhone: "123456789",
		Description:  "Test description 1",
		Level:        model.High,
		Status:       model.Open,
	}
	urgency2 := &model.Urgency{
		Name:         "Test Urgency 2",
		Email:        "test2@example.com",
		ContactPhone: "987654321",
		Description:  "Test description 2",
		Level:        model.Medium,
		Status:       model.InProgress,
	}

	require.NoError(t, repo.Create(urgency1))
	require.NoError(t, repo.Create(urgency2))

	urgencies, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, urgencies, 2)

	err = repo.Delete(urgency1.ID)
	require.NoError(t, err)

	urgencies, err = repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, urgencies, 1)

	err = repo.ResetAllData()
	assert.NoError(t, err)

	urgencies, err = repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, urgencies, 0)

	var count int64
	err = db.Unscoped().Model(&model.Urgency{}).Count(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
