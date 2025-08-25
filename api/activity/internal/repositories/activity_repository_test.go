package repositories

import (
	"database/sql"
	"testing"

	"github.com/pd120424d/mountain-service/api/activity/internal/model"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func setupActivityTestDB(t *testing.T) *gorm.DB {
	sqlDB, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	db, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Activity{})
	require.NoError(t, err)

	return db
}

func TestActivityRepository_Create(t *testing.T) {
	t.Parallel()

	t.Run("successfully creates an activity", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		activity := &model.Activity{
			Description: "Test Description",
			EmployeeID:  1,
			UrgencyID:   2,
		}

		err := repo.Create(activity)
		assert.NoError(t, err)
		assert.NotZero(t, activity.ID)

		var dbActivity model.Activity
		err = db.First(&dbActivity, activity.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, "Test Description", dbActivity.Description)
		assert.Equal(t, uint(1), dbActivity.EmployeeID)
		assert.Equal(t, uint(2), dbActivity.UrgencyID)
	})

	t.Run("successfully creates activity with different level", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		activity := &model.Activity{
			Description: "System has been reset",
			EmployeeID:  3,
			UrgencyID:   4,
		}

		err := repo.Create(activity)
		assert.NoError(t, err)
		assert.NotZero(t, activity.ID)

		var dbActivity model.Activity
		err = db.First(&dbActivity, activity.ID).Error
		assert.NoError(t, err)
	})
}

func TestActivityRepository_GetByID(t *testing.T) {
	t.Parallel()

	t.Run("successfully retrieves activity by ID", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		activity := &model.Activity{
			Description: "Test Description",
			EmployeeID:  1,
			UrgencyID:   2,
		}
		err := db.Create(activity).Error
		require.NoError(t, err)

		retrievedActivity, err := repo.GetByID(activity.ID)
		assert.NoError(t, err)
		require.NotNil(t, retrievedActivity)
		assert.Equal(t, activity.ID, retrievedActivity.ID)
		assert.Equal(t, activity.Description, retrievedActivity.Description)
		assert.Equal(t, activity.EmployeeID, retrievedActivity.EmployeeID)
		assert.Equal(t, activity.UrgencyID, retrievedActivity.UrgencyID)
	})

	t.Run("returns error when activity not found", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		retrievedActivity, err := repo.GetByID(999)
		assert.Error(t, err)
		assert.Nil(t, retrievedActivity)
		assert.Contains(t, err.Error(), "activity not found")
	})
}

func TestActivityRepository_List(t *testing.T) {
	t.Parallel()

	t.Run("successfully lists activities with no filters", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		activities := []*model.Activity{
			{
				Description: "Employee was assigned to urgency",
				EmployeeID:  1,
				UrgencyID:   2,
			},
			{
				Description: "Employee resolved urgency",
				EmployeeID:  1,
				UrgencyID:   2,
			},
		}

		for _, activity := range activities {
			err := db.Create(activity).Error
			require.NoError(t, err)
		}

		filter := model.NewActivityFilter()
		result, total, err := repo.List(filter)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(2), total)
	})

	t.Run("validates filter parameters", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		filter := &model.ActivityFilter{
			Page:     -1,
			PageSize: -1,
		}
		result, total, err := repo.List(filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, result, 0)
		assert.Equal(t, 1, filter.Page)
		assert.Equal(t, model.DefaultPageSize, filter.PageSize)
	})
}

func TestActivityRepository_Delete(t *testing.T) {
	t.Parallel()

	t.Run("successfully deletes activity", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		activity := &model.Activity{
			Description: "Test Description",
		}
		err := db.Create(activity).Error
		require.NoError(t, err)

		err = repo.Delete(activity.ID)
		assert.NoError(t, err)

		var deletedActivity model.Activity
		err = db.First(&deletedActivity, activity.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("returns error when activity not found for deletion", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		err := repo.Delete(999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "activity not found")
	})
}

func TestActivityRepository_ResetAllData(t *testing.T) {
	t.Parallel()

	t.Run("successfully resets all activity data", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		activities := []*model.Activity{
			{
				Description: "Employee was created",
				EmployeeID:  1,
				UrgencyID:   2,
			},
			{
				Description: "Urgency was created",
				EmployeeID:  1,
				UrgencyID:   2,
			},
		}

		for _, activity := range activities {
			err := db.Create(activity).Error
			require.NoError(t, err)
		}

		var count int64
		err := db.Model(&model.Activity{}).Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(2), count)

		err = repo.ResetAllData()
		assert.NoError(t, err)

		err = db.Model(&model.Activity{}).Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestActivityRepository_GetStats(t *testing.T) {
	t.Parallel()

	t.Run("successfully retrieves activity statistics", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		activities := []*model.Activity{
			{
				Description: "Employee was created",
			},
			{
				Description: "Another employee was created",
			},
			{
				Description: "Urgency was created",
			},
		}

		for _, activity := range activities {
			err := db.Create(activity).Error
			require.NoError(t, err)
		}

		stats, err := repo.GetStats()
		assert.NoError(t, err)
		require.NotNil(t, stats)

		assert.Equal(t, int64(3), stats.TotalActivities)
		assert.Len(t, stats.RecentActivities, 3) // Should return all 3 since limit is 10
	})

	t.Run("returns empty stats when no activities exist", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		stats, err := repo.GetStats()
		assert.NoError(t, err)
		require.NotNil(t, stats)

		assert.Equal(t, int64(0), stats.TotalActivities)
		assert.Equal(t, int64(0), stats.ActivitiesLast7Days)
		assert.Equal(t, int64(0), stats.ActivitiesLast30Days)
	})
}

func TestActivityRepository_DatabaseConnection(t *testing.T) {
	t.Parallel()

	t.Run("it handles database connection properly", func(t *testing.T) {
		// Test that the repository works with a proper database connection
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		activity := &model.Activity{
			Description: "Test",
		}

		err := repo.Create(activity)
		assert.NoError(t, err)
		assert.NotZero(t, activity.ID)
	})
}

func TestActivityRepository_QueryBuilding(t *testing.T) {
	t.Parallel()

	t.Run("it builds complex queries correctly", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		// Create test data
		employeeID := uint(1)
		urgencyID := uint(2)

		activity1 := &model.Activity{
			Description: "Test employee creation",
			EmployeeID:  employeeID,
		}
		activity2 := &model.Activity{
			Description: "Test urgency creation",
			EmployeeID:  employeeID,
			UrgencyID:   urgencyID,
		}

		err := repo.Create(activity1)
		assert.NoError(t, err)
		err = repo.Create(activity2)
		assert.NoError(t, err)

		// Test complex filter with multiple criteria
		filter := &model.ActivityFilter{
			EmployeeID: &employeeID,
			UrgencyID:  &urgencyID,
			Page:       1,
			PageSize:   10,
		}

		activities, total, err := repo.List(filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, activities, 1)
	})
}
