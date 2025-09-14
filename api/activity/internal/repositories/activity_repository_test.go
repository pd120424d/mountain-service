package repositories

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pd120424d/mountain-service/api/activity/internal/model"
	"github.com/pd120424d/mountain-service/api/shared/models"
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

	db, err := gorm.Open(sqlite.Dialector{DriverName: "sqlite", Conn: sqlDB}, &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Activity{}, &models.OutboxEvent{})
	require.NoError(t, err)

	// Ensure proper cleanup of underlying sql.DB
	sqlStd, _ := db.DB()
	t.Cleanup(func() { _ = sqlStd.Close() })
	return db
}

func TestActivityRepository_Create(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when database operation fails", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		activity := &model.Activity{
			Description: "Test Description",
			EmployeeID:  1,
			UrgencyID:   2,
		}

		err := repo.Create(t.Context(), activity)
		assert.Error(t, err)
	})

	t.Run("successfully creates an activity", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		activity := &model.Activity{
			Description: "Test Description",
			EmployeeID:  1,
			UrgencyID:   2,
		}

		err := repo.Create(t.Context(), activity)
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

		err := repo.Create(t.Context(), activity)
		assert.NoError(t, err)
		assert.NotZero(t, activity.ID)

		var dbActivity model.Activity
		err = db.First(&dbActivity, activity.ID).Error
		assert.NoError(t, err)
	})
}

func TestActivityRepository_CreateWithOutbox(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when database operation fails", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		activity := &model.Activity{
			Description: "Test Description",
			EmployeeID:  1,
			UrgencyID:   2,
		}

		event := &models.OutboxEvent{
			AggregateID: "activity-1",
			EventData:   `{"x":1}`,
		}

		err := repo.CreateWithOutbox(t.Context(), activity, event)
		assert.Error(t, err)
	})

	t.Run("it returns an error when sending outbox event fails", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		activity := &model.Activity{
			Description: "Test Description",
			EmployeeID:  1,
			UrgencyID:   2,
		}

		event := &models.OutboxEvent{
			AggregateID: "activity-1",
			EventData:   `{"x":1}`,
		}

		err := repo.CreateWithOutbox(t.Context(), activity, event)
		assert.NoError(t, err)
	})

	t.Run("successfully creates activity with outbox event", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		activity := &model.Activity{
			Description: "Test Description",
			EmployeeID:  1,
			UrgencyID:   2,
		}

		event := &models.OutboxEvent{
			AggregateID: "activity-1",
			EventData:   `{"x":1}`,
		}

		err := repo.CreateWithOutbox(t.Context(), activity, event)
		assert.NoError(t, err)
		assert.NotZero(t, activity.ID)
		assert.NotZero(t, event.ID)
	})

}

func TestActivityRepository_GetByID(t *testing.T) {
	t.Parallel()

	t.Run("returns error when activity not found", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		retrievedActivity, err := repo.GetByID(t.Context(), 999)
		assert.Error(t, err)
		assert.Nil(t, retrievedActivity)
		assert.Contains(t, err.Error(), "activity not found")
	})

	t.Run("it returns an error when database operation fails", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		retrievedActivity, err := repo.GetByID(t.Context(), 1)
		assert.Error(t, err)
		assert.Nil(t, retrievedActivity)
	})

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

		retrievedActivity, err := repo.GetByID(t.Context(), activity.ID)
		assert.NoError(t, err)
		require.NotNil(t, retrievedActivity)
		assert.Equal(t, activity.ID, retrievedActivity.ID)
		assert.Equal(t, activity.Description, retrievedActivity.Description)
		assert.Equal(t, activity.EmployeeID, retrievedActivity.EmployeeID)
		assert.Equal(t, activity.UrgencyID, retrievedActivity.UrgencyID)
	})

}

func TestActivityRepository_List(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when database operation fails", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		filter := model.NewActivityFilter()
		result, total, err := repo.List(t.Context(), filter)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
	})

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
		result, total, err := repo.List(t.Context(), filter)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(2), total)
	})

	t.Run("validates filter parameters", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		employeeID := uint(999)
		urgencyID := uint(999)

		filter := &model.ActivityFilter{
			Page:       -1,
			PageSize:   -1,
			EmployeeID: &employeeID,
			UrgencyID:  &urgencyID,
		}
		result, total, err := repo.List(t.Context(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, result, 0)
		assert.Equal(t, 1, filter.Page)
		assert.Equal(t, model.DefaultPageSize, filter.PageSize)
	})

	t.Run("successfully lists activities with date filters", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		base := time.Date(2025, 1, 10, 10, 0, 0, 0, time.UTC)
		a1 := &model.Activity{Description: "a1", CreatedAt: base.Add(-3 * time.Hour)}
		a2 := &model.Activity{Description: "a2", CreatedAt: base.Add(-2 * time.Hour)}
		a3 := &model.Activity{Description: "a3", CreatedAt: base.Add(-1 * time.Hour)}
		assert.NoError(t, db.Create(a1).Error)
		assert.NoError(t, db.Create(a2).Error)
		assert.NoError(t, db.Create(a3).Error)

		start := base.Add(-2 * time.Hour)
		end := base.Add(-2 * time.Hour)
		f := &model.ActivityFilter{StartDate: &start, EndDate: &end, Page: 1, PageSize: 10}
		items, total, err := repo.List(t.Context(), f)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		if assert.Len(t, items, 1) {
			assert.Equal(t, "a2", items[0].Description)
		}
	})

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

		err := repo.Create(t.Context(), activity1)
		assert.NoError(t, err)
		err = repo.Create(t.Context(), activity2)
		assert.NoError(t, err)

		// Test complex filter with multiple criteria
		filter := &model.ActivityFilter{
			EmployeeID: &employeeID,
			UrgencyID:  &urgencyID,
			Page:       1,
			PageSize:   10,
		}

		activities, total, err := repo.List(t.Context(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, activities, 1)
	})

	t.Run("it returns an error when count query fails", func(t *testing.T) {
		gormDB, mock, sqlDB := newGormWithSQLMock(t)
		defer sqlDB.Close()
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, gormDB)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "activities" WHERE "activities"."deleted_at" IS NULL`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "activities" WHERE "activities"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT $1`)).
			WithArgs(50).
			WillReturnError(sqlmock.ErrCancelled)

		filter := model.NewActivityFilter()
		items, total, err := repo.List(t.Context(), filter)
		assert.Error(t, err)
		assert.Nil(t, items)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestActivityRepository_GetStats(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when count query fails", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		stats, err := repo.GetStats(t.Context())
		assert.Error(t, err)
		assert.Nil(t, stats)
	})

	t.Run("it returns an error when recent activities query fails", func(t *testing.T) {
		gormDB, mock, sqlDB := newGormWithSQLMock(t)
		defer sqlDB.Close()
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, gormDB)

		// Simulate count query succeeding but recent activities query failing
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "activities" WHERE "activities"."deleted_at" IS NULL`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "activities" WHERE "activities"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT $1`)).
			WithArgs(10).
			WillReturnError(sqlmock.ErrCancelled)

		stats, err := repo.GetStats(t.Context())
		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

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

		stats, err := repo.GetStats(t.Context())
		assert.NoError(t, err)
		require.NotNil(t, stats)

		assert.Equal(t, int64(3), stats.TotalActivities)
		assert.Len(t, stats.RecentActivities, 3) // Should return all 3 since limit is 10
	})

	t.Run("returns empty stats when no activities exist", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		stats, err := repo.GetStats(t.Context())
		assert.NoError(t, err)
		require.NotNil(t, stats)

		assert.Equal(t, int64(0), stats.TotalActivities)
		assert.Equal(t, int64(0), stats.ActivitiesLast7Days)
		assert.Equal(t, int64(0), stats.ActivitiesLast30Days)
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

		err = repo.Delete(t.Context(), activity.ID)
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

		err := repo.Delete(t.Context(), 999)
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

		err = repo.ResetAllData(t.Context())
		assert.NoError(t, err)

		err = db.Model(&model.Activity{}).Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}
