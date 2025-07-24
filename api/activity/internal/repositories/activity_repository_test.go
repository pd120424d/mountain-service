package repositories

import (
	"database/sql"
	"testing"

	"github.com/pd120424d/mountain-service/api/activity/internal/model"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
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

		actorID := uint(1)
		targetID := uint(2)
		activity := &model.Activity{
			Type:        activityV1.ActivityEmployeeCreated,
			Level:       activityV1.ActivityLevelInfo,
			Title:       "Test Title",
			Description: "Test Description",
			ActorID:     &actorID,
			ActorName:   "Test Actor",
			TargetID:    &targetID,
			TargetType:  "Test Target",
			Metadata:    "Test Metadata",
		}

		err := repo.Create(activity)
		assert.NoError(t, err)
		assert.NotZero(t, activity.ID)

		var dbActivity model.Activity
		err = db.First(&dbActivity, activity.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, activityV1.ActivityEmployeeCreated, dbActivity.Type)
		assert.Equal(t, activityV1.ActivityLevelInfo, dbActivity.Level)
		assert.Equal(t, "Test Title", dbActivity.Title)
		assert.Equal(t, "Test Description", dbActivity.Description)
		assert.Equal(t, actorID, *dbActivity.ActorID)
		assert.Equal(t, "Test Actor", dbActivity.ActorName)
		assert.Equal(t, targetID, *dbActivity.TargetID)
		assert.Equal(t, "Test Target", dbActivity.TargetType)
		assert.Equal(t, "Test Metadata", dbActivity.Metadata)
	})

	t.Run("successfully creates activity with nil ActorID and TargetID", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		activity := &model.Activity{
			Type:        activityV1.ActivitySystemReset,
			Level:       activityV1.ActivityLevelInfo,
			Title:       "System Reset",
			Description: "System has been reset",
			ActorID:     nil,
			ActorName:   "system",
			TargetID:    nil,
			TargetType:  "system",
			Metadata:    "{}",
		}

		err := repo.Create(activity)
		assert.NoError(t, err)
		assert.NotZero(t, activity.ID)

		var dbActivity model.Activity
		err = db.First(&dbActivity, activity.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, activityV1.ActivitySystemReset, dbActivity.Type)
		assert.Nil(t, dbActivity.ActorID)
		assert.Nil(t, dbActivity.TargetID)
	})
}

func TestActivityRepository_GetByID(t *testing.T) {
	t.Parallel()

	t.Run("successfully retrieves activity by ID", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		actorID := uint(1)
		targetID := uint(2)
		activity := &model.Activity{
			Type:        activityV1.ActivityEmployeeCreated,
			Level:       activityV1.ActivityLevelInfo,
			Title:       "Test Title",
			Description: "Test Description",
			ActorID:     &actorID,
			ActorName:   "Test Actor",
			TargetID:    &targetID,
			TargetType:  "employee",
			Metadata:    "{}",
		}
		err := db.Create(activity).Error
		require.NoError(t, err)

		retrievedActivity, err := repo.GetByID(activity.ID)
		assert.NoError(t, err)
		require.NotNil(t, retrievedActivity)
		assert.Equal(t, activity.ID, retrievedActivity.ID)
		assert.Equal(t, activity.Type, retrievedActivity.Type)
		assert.Equal(t, activity.Level, retrievedActivity.Level)
		assert.Equal(t, activity.Title, retrievedActivity.Title)
		assert.Equal(t, activity.Description, retrievedActivity.Description)
		assert.Equal(t, *activity.ActorID, *retrievedActivity.ActorID)
		assert.Equal(t, activity.ActorName, retrievedActivity.ActorName)
		assert.Equal(t, *activity.TargetID, *retrievedActivity.TargetID)
		assert.Equal(t, activity.TargetType, retrievedActivity.TargetType)
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

		actorID := uint(1)
		targetID := uint(2)
		activities := []*model.Activity{
			{
				Type:        activityV1.ActivityEmployeeCreated,
				Level:       activityV1.ActivityLevelInfo,
				Title:       "Employee Created",
				Description: "Employee was created",
				ActorID:     &actorID,
				ActorName:   "Admin",
				TargetID:    &targetID,
				TargetType:  "employee",
				Metadata:    "{}",
			},
			{
				Type:        activityV1.ActivityUrgencyCreated,
				Level:       activityV1.ActivityLevelWarning,
				Title:       "Urgency Created",
				Description: "Urgency was created",
				ActorID:     &actorID,
				ActorName:   "Admin",
				TargetID:    &targetID,
				TargetType:  "urgency",
				Metadata:    "{}",
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

	t.Run("successfully lists activities with type filter", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		actorID := uint(1)
		targetID := uint(2)
		activities := []*model.Activity{
			{
				Type:        activityV1.ActivityEmployeeCreated,
				Level:       activityV1.ActivityLevelInfo,
				Title:       "Employee Created",
				Description: "Employee was created",
				ActorID:     &actorID,
				ActorName:   "Admin",
				TargetID:    &targetID,
				TargetType:  "employee",
				Metadata:    "{}",
			},
			{
				Type:        activityV1.ActivityUrgencyCreated,
				Level:       activityV1.ActivityLevelWarning,
				Title:       "Urgency Created",
				Description: "Urgency was created",
				ActorID:     &actorID,
				ActorName:   "Admin",
				TargetID:    &targetID,
				TargetType:  "urgency",
				Metadata:    "{}",
			},
		}

		for _, activity := range activities {
			err := db.Create(activity).Error
			require.NoError(t, err)
		}

		filter := model.NewActivityFilter()
		employeeType := activityV1.ActivityEmployeeCreated
		filter.Type = &employeeType
		result, total, err := repo.List(filter)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, activityV1.ActivityEmployeeCreated, result[0].Type)
	})

	t.Run("successfully lists activities with level filter", func(t *testing.T) {
		db := setupActivityTestDB(t)
		log := utils.NewTestLogger()
		repo := NewActivityRepository(log, db)

		actorID := uint(1)
		targetID := uint(2)
		activities := []*model.Activity{
			{
				Type:        activityV1.ActivityEmployeeCreated,
				Level:       activityV1.ActivityLevelInfo,
				Title:       "Employee Created",
				Description: "Employee was created",
				ActorID:     &actorID,
				ActorName:   "Admin",
				TargetID:    &targetID,
				TargetType:  "employee",
				Metadata:    "{}",
			},
			{
				Type:        activityV1.ActivityUrgencyCreated,
				Level:       activityV1.ActivityLevelWarning,
				Title:       "Urgency Created",
				Description: "Urgency was created",
				ActorID:     &actorID,
				ActorName:   "Admin",
				TargetID:    &targetID,
				TargetType:  "urgency",
				Metadata:    "{}",
			},
		}

		for _, activity := range activities {
			err := db.Create(activity).Error
			require.NoError(t, err)
		}

		filter := model.NewActivityFilter()
		warningLevel := activityV1.ActivityLevelWarning
		filter.Level = &warningLevel
		result, total, err := repo.List(filter)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, activityV1.ActivityLevelWarning, result[0].Level)
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

		actorID := uint(1)
		targetID := uint(2)
		activity := &model.Activity{
			Type:        activityV1.ActivityEmployeeCreated,
			Level:       activityV1.ActivityLevelInfo,
			Title:       "Test Title",
			Description: "Test Description",
			ActorID:     &actorID,
			ActorName:   "Test Actor",
			TargetID:    &targetID,
			TargetType:  "employee",
			Metadata:    "{}",
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

		actorID := uint(1)
		targetID := uint(2)
		activities := []*model.Activity{
			{
				Type:        activityV1.ActivityEmployeeCreated,
				Level:       activityV1.ActivityLevelInfo,
				Title:       "Employee Created",
				Description: "Employee was created",
				ActorID:     &actorID,
				ActorName:   "Admin",
				TargetID:    &targetID,
				TargetType:  "employee",
				Metadata:    "{}",
			},
			{
				Type:        activityV1.ActivityUrgencyCreated,
				Level:       activityV1.ActivityLevelWarning,
				Title:       "Urgency Created",
				Description: "Urgency was created",
				ActorID:     &actorID,
				ActorName:   "Admin",
				TargetID:    &targetID,
				TargetType:  "urgency",
				Metadata:    "{}",
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

		actorID := uint(1)
		targetID := uint(2)
		activities := []*model.Activity{
			{
				Type:        activityV1.ActivityEmployeeCreated,
				Level:       activityV1.ActivityLevelInfo,
				Title:       "Employee Created 1",
				Description: "Employee was created",
				ActorID:     &actorID,
				ActorName:   "Admin",
				TargetID:    &targetID,
				TargetType:  "employee",
				Metadata:    "{}",
			},
			{
				Type:        activityV1.ActivityEmployeeCreated,
				Level:       activityV1.ActivityLevelInfo,
				Title:       "Employee Created 2",
				Description: "Another employee was created",
				ActorID:     &actorID,
				ActorName:   "Admin",
				TargetID:    &targetID,
				TargetType:  "employee",
				Metadata:    "{}",
			},
			{
				Type:        activityV1.ActivityUrgencyCreated,
				Level:       activityV1.ActivityLevelWarning,
				Title:       "Urgency Created",
				Description: "Urgency was created",
				ActorID:     &actorID,
				ActorName:   "Admin",
				TargetID:    &targetID,
				TargetType:  "urgency",
				Metadata:    "{}",
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
		assert.Equal(t, int64(2), stats.ActivitiesByType[activityV1.ActivityEmployeeCreated])
		assert.Equal(t, int64(1), stats.ActivitiesByType[activityV1.ActivityUrgencyCreated])
		assert.Equal(t, int64(2), stats.ActivitiesByLevel[activityV1.ActivityLevelInfo])
		assert.Equal(t, int64(1), stats.ActivitiesByLevel[activityV1.ActivityLevelWarning])
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
		assert.Empty(t, stats.ActivitiesByType)
		assert.Empty(t, stats.ActivitiesByLevel)
		assert.Empty(t, stats.RecentActivities)
		assert.Equal(t, int64(0), stats.ActivitiesLast24h)
		assert.Equal(t, int64(0), stats.ActivitiesLast7Days)
		assert.Equal(t, int64(0), stats.ActivitiesLast30Days)
	})
}
