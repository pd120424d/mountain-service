package repositories

import (
	"context"
	"database/sql"
	"testing"

	urgencyV1 "github.com/pd120424d/mountain-service/api/contracts/urgency/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/pd120424d/mountain-service/api/urgency/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Use pure Go SQLite driver (modernc.org/sqlite) instead of CGO-based mattn/go-sqlite3
	// Create a custom SQLite connection using the pure Go driver
	sqlDB, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	db, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Urgency{}, &model.Notification{})
	require.NoError(t, err)

	return db
}

func TestUrgencyRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	log := utils.NewTestLogger()
	repo := NewUrgencyRepository(log, db)

	urgency := &model.Urgency{
		FirstName:    "Marko",
		LastName:     "Markovic",
		Email:        "test@example.com",
		ContactPhone: "123456789",
		Description:  "Test description",
		Level:        urgencyV1.High,
		Status:       urgencyV1.Open,
	}

	err := repo.Create(context.Background(), urgency)
	assert.NoError(t, err)
	assert.NotZero(t, urgency.ID)
}

func TestUrgencyRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	log := utils.NewTestLogger()
	repo := NewUrgencyRepository(log, db)

	urgency1 := &model.Urgency{
		FirstName:    "Marko",
		LastName:     "Markovic",
		Email:        "test1@example.com",
		ContactPhone: "123456789",
		Description:  "Test description 1",
		Level:        urgencyV1.High,
		Status:       urgencyV1.Open,
	}
	urgency2 := &model.Urgency{
		FirstName:    "Marko",
		LastName:     "Markovic",
		Email:        "test2@example.com",
		ContactPhone: "987654321",
		Description:  "Test description 2",
		Level:        urgencyV1.Medium,
		Status:       urgencyV1.InProgress,
	}

	err := repo.Create(context.Background(), urgency1)
	require.NoError(t, err)
	err = repo.Create(context.Background(), urgency2)
	require.NoError(t, err)

	urgencies, err := repo.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, urgencies, 2)
}

func TestUrgencyRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	log := utils.NewTestLogger()
	repo := NewUrgencyRepository(log, db)

	urgency := &model.Urgency{
		FirstName:    "Marko",
		LastName:     "Markovic",
		Email:        "test@example.com",
		ContactPhone: "123456789",
		Description:  "Test description",
		Level:        urgencyV1.High,
		Status:       urgencyV1.Open,
	}

	err := repo.Create(context.Background(), urgency)
	require.NoError(t, err)

	var retrieved model.Urgency
	err = repo.GetByID(context.Background(), urgency.ID, &retrieved)
	assert.NoError(t, err)
	assert.Equal(t, urgency.FirstName, retrieved.FirstName)
	assert.Equal(t, urgency.LastName, retrieved.LastName)
	assert.Equal(t, urgency.Email, retrieved.Email)
}

func TestUrgencyRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	log := utils.NewTestLogger()
	repo := NewUrgencyRepository(log, db)

	urgency := &model.Urgency{
		FirstName:    "Marko",
		LastName:     "Markovic",
		Email:        "test@example.com",
		ContactPhone: "123456789",
		Description:  "Test description",
		Level:        urgencyV1.High,
		Status:       urgencyV1.Open,
	}

	err := repo.Create(context.Background(), urgency)
	require.NoError(t, err)

	urgency.FirstName = "Marko"
	urgency.LastName = "Markovic"
	urgency.Status = urgencyV1.InProgress

	err = repo.Update(context.Background(), urgency)
	assert.NoError(t, err)

	var updated model.Urgency
	err = repo.GetByID(context.Background(), urgency.ID, &updated)
	require.NoError(t, err)
	assert.Equal(t, "Marko", updated.FirstName)
	assert.Equal(t, "Markovic", updated.LastName)
	assert.Equal(t, urgencyV1.InProgress, updated.Status)
}

func TestUrgencyRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	log := utils.NewTestLogger()
	repo := NewUrgencyRepository(log, db)

	urgency := &model.Urgency{
		FirstName:    "Marko",
		LastName:     "Markovic",
		Email:        "test@example.com",
		ContactPhone: "123456789",
		Description:  "Test description",
		Level:        urgencyV1.High,
		Status:       urgencyV1.Open,
	}

	err := repo.Create(context.Background(), urgency)
	require.NoError(t, err)

	err = repo.Delete(context.Background(), urgency.ID)
	assert.NoError(t, err)

	var deleted model.Urgency
	err = repo.GetByID(context.Background(), urgency.ID, &deleted)
	assert.Error(t, err) // Should not find deleted record
}

func TestUrgencyRepository_List(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	log := utils.NewTestLogger()
	repo := NewUrgencyRepository(log, db)

	urgency1 := &model.Urgency{
		FirstName:    "Emergency",
		LastName:     "Service",
		Email:        "rescue@example.com",
		ContactPhone: "123456789",
		Description:  "Mountain rescue needed",
		Level:        urgencyV1.Critical,
		Status:       urgencyV1.Open,
	}
	urgency2 := &model.Urgency{
		FirstName:    "Medical",
		LastName:     "Help",
		Email:        "medical@example.com",
		ContactPhone: "987654321",
		Description:  "Medical assistance required",
		Level:        urgencyV1.High,
		Status:       urgencyV1.InProgress,
	}
	urgency3 := &model.Urgency{
		FirstName:    "Equipment",
		LastName:     "Issue",
		Email:        "equipment@example.com",
		ContactPhone: "555666777",
		Description:  "Equipment malfunction",
		Level:        urgencyV1.Medium,
		Status:       urgencyV1.Resolved,
	}

	require.NoError(t, repo.Create(context.Background(), urgency1))
	require.NoError(t, repo.Create(context.Background(), urgency2))
	require.NoError(t, repo.Create(context.Background(), urgency3))

	t.Run("it returns all urgencies when no filters are provided", func(t *testing.T) {
		urgencies, err := repo.List(context.Background(), map[string]interface{}{})
		assert.NoError(t, err)
		assert.Len(t, urgencies, 3)
	})

	t.Run("it returns urgencies filtered by level", func(t *testing.T) {
		filters := map[string]interface{}{
			"level": "High",
		}
		urgencies, err := repo.List(context.Background(), filters)
		assert.NoError(t, err)
		assert.Len(t, urgencies, 1)
		assert.Equal(t, "Medical", urgencies[0].FirstName)
		assert.Equal(t, "Help", urgencies[0].LastName)
	})

	t.Run("it returns urgencies filtered by status", func(t *testing.T) {
		filters := map[string]interface{}{
			"status": "Open",
		}
		urgencies, err := repo.List(context.Background(), filters)
		assert.NoError(t, err)
		assert.Len(t, urgencies, 1)
		assert.Equal(t, "Emergency", urgencies[0].FirstName)
		assert.Equal(t, "Service", urgencies[0].LastName)
	})

	t.Run("it returns urgencies filtered by first_name (partial match)", func(t *testing.T) {
		filters := map[string]interface{}{
			"first_name": "Medical",
		}
		urgencies, err := repo.List(context.Background(), filters)
		assert.NoError(t, err)
		assert.Len(t, urgencies, 1)
		assert.Equal(t, "Medical", urgencies[0].FirstName)
		assert.Equal(t, "Help", urgencies[0].LastName)
	})

	t.Run("it returns urgencies filtered by email (partial match)", func(t *testing.T) {
		filters := map[string]interface{}{
			"email": "rescue",
		}
		urgencies, err := repo.List(context.Background(), filters)
		assert.NoError(t, err)
		assert.Len(t, urgencies, 1)
		assert.Equal(t, "Emergency", urgencies[0].FirstName)
		assert.Equal(t, "Service", urgencies[0].LastName)
	})

	t.Run("it returns urgencies filtered by multiple filters", func(t *testing.T) {
		filters := map[string]interface{}{
			"level":  "Medium",
			"status": "Resolved",
		}
		urgencies, err := repo.List(context.Background(), filters)
		assert.NoError(t, err)
		assert.Len(t, urgencies, 1)
		assert.Equal(t, "Equipment", urgencies[0].FirstName)
		assert.Equal(t, "Issue", urgencies[0].LastName)
	})

	t.Run("it returns an error when an invalid filter key is provided", func(t *testing.T) {
		filters := map[string]interface{}{
			"invalid_field": "value",
		}
		urgencies, err := repo.List(context.Background(), filters)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid filter key: invalid_field")
		assert.Nil(t, urgencies)
	})

	t.Run("it returns an error when an unsupported filter value type is provided", func(t *testing.T) {
		filters := map[string]interface{}{
			"first_name": []string{"test"},
		}
		urgencies, err := repo.List(context.Background(), filters)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported type for filter key: first_name")
		assert.Nil(t, urgencies)
	})

	t.Run("it returns an empty slice when no matches are found", func(t *testing.T) {
		filters := map[string]interface{}{
			"first_name": "NonExistent",
		}
		urgencies, err := repo.List(context.Background(), filters)
		assert.NoError(t, err)
		assert.Len(t, urgencies, 0)
	})
}

func TestUrgencyRepository_ResetAllData(t *testing.T) {
	db := setupTestDB(t)
	log := utils.NewTestLogger()
	repo := NewUrgencyRepository(log, db)

	urgency1 := &model.Urgency{
		FirstName:    "Marko",
		LastName:     "Markovic",
		Email:        "test1@example.com",
		ContactPhone: "123456789",
		Description:  "Test description 1",
		Level:        urgencyV1.High,
		Status:       urgencyV1.Open,
	}
	urgency2 := &model.Urgency{
		FirstName:    "Marko",
		LastName:     "Markovic",
		Email:        "test2@example.com",
		ContactPhone: "987654321",
		Description:  "Test description 2",
		Level:        urgencyV1.Medium,
		Status:       urgencyV1.InProgress,
	}

	require.NoError(t, repo.Create(context.Background(), urgency1))
	require.NoError(t, repo.Create(context.Background(), urgency2))

	urgencies, err := repo.GetAll(context.Background())
	require.NoError(t, err)
	assert.Len(t, urgencies, 2)

	err = repo.Delete(context.Background(), urgency1.ID)
	require.NoError(t, err)

	urgencies, err = repo.GetAll(context.Background())
	require.NoError(t, err)
	assert.Len(t, urgencies, 1)

	err = repo.ResetAllData(context.Background())
	assert.NoError(t, err)

	urgencies, err = repo.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, urgencies, 0)

	var count int64
	err = db.Unscoped().Model(&model.Urgency{}).Count(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)

}

func TestUrgencyRepository_ListPaginated(t *testing.T) {
	db := setupTestDB(t)
	log := utils.NewTestLogger()
	repo := NewUrgencyRepository(log, db)

	// Seed 35 rows
	for i := 0; i < 35; i++ {
		u := &model.Urgency{
			FirstName:    "F",
			LastName:     "L",
			Email:        "e@x.com",
			ContactPhone: "1",
			Description:  "d",
			Level:        urgencyV1.Medium,
			Status:       urgencyV1.Open,
		}
		require.NoError(t, repo.Create(context.Background(), u))
	}

	page := 2
	pageSize := 10
	items, total, err := repo.ListPaginated(context.Background(), page, pageSize, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(35), total)
	assert.Len(t, items, 10)

	// Page beyond total pages returns empty slice
	items, total, err = repo.ListPaginated(context.Background(), 4, 10, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(35), total)
	assert.Len(t, items, 5)
}

func TestUrgencyRepository_ListUnassignedIDs(t *testing.T) {
	db := setupTestDB(t)
	log := utils.NewTestLogger()
	repo := NewUrgencyRepository(log, db)

	// u1: open & unassigned -> SortPriority 1 (should be returned)
	u1 := &model.Urgency{FirstName: "A", LastName: "U", ContactPhone: "1", Location: "L", Description: "d", Level: urgencyV1.Medium, Status: urgencyV1.Open, SortPriority: 1}
	require.NoError(t, repo.Create(context.Background(), u1))

	// u2: open & assigned -> SortPriority 2 (should be excluded)
	emp := uint(10)
	u2 := &model.Urgency{FirstName: "B", LastName: "A", ContactPhone: "1", Location: "L", Description: "d", Level: urgencyV1.Medium, Status: urgencyV1.Open, AssignedEmployeeID: &emp, SortPriority: 2}
	require.NoError(t, repo.Create(context.Background(), u2))

	// u3: in_progress -> SortPriority 3 (should be excluded)
	u3 := &model.Urgency{FirstName: "C", LastName: "P", ContactPhone: "1", Location: "L", Description: "d", Level: urgencyV1.Medium, Status: urgencyV1.InProgress, SortPriority: 3}
	require.NoError(t, repo.Create(context.Background(), u3))

	ids, err := repo.ListUnassignedIDs(context.Background())
	assert.NoError(t, err)
	if assert.Len(t, ids, 1) {
		assert.Equal(t, u1.ID, ids[0])
	}
}
