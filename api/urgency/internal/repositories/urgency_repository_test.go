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
		Status:       "Open",
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
		Status:       "Open",
	}

	err := repo.Create(urgency)
	require.NoError(t, err)

	urgency.Name = "Updated Urgency"
	urgency.Status = "In Progress"

	err = repo.Update(urgency)
	assert.NoError(t, err)

	var updated model.Urgency
	err = repo.GetByID(urgency.ID, &updated)
	require.NoError(t, err)
	assert.Equal(t, "Updated Urgency", updated.Name)
	assert.Equal(t, "In Progress", updated.Status)
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
