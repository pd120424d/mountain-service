package repositories

import (
	"database/sql"
	"testing"
	"time"

	urgencyV1 "github.com/pd120424d/mountain-service/api/contracts/urgency/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/pd120424d/mountain-service/api/urgency/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func TestNotificationRepository_Create(t *testing.T) {
	t.Parallel()

	t.Run("successfully creates notification", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		urgency := createTestUrgencyForNotification(t, db)

		notification := &model.Notification{
			UrgencyID:        urgency.ID,
			EmployeeID:       1,
			NotificationType: model.NotificationSMS,
			Recipient:        "+1234567890",
			Message:          "Test notification message",
			Status:           model.NotificationPending,
		}

		err := repo.Create(notification)
		assert.NoError(t, err)
		assert.NotZero(t, notification.ID)

		var dbNotification model.Notification
		err = db.First(&dbNotification, notification.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, urgency.ID, dbNotification.UrgencyID)
		assert.Equal(t, uint(1), dbNotification.EmployeeID)
		assert.Equal(t, model.NotificationSMS, dbNotification.NotificationType)
		assert.Equal(t, "+1234567890", dbNotification.Recipient)
		assert.Equal(t, "Test notification message", dbNotification.Message)
		assert.Equal(t, model.NotificationPending, dbNotification.Status)
	})

	t.Run("creates notification even when urgency does not exist in database", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		notification := &model.Notification{
			UrgencyID:        999, // Non-existent urgency
			EmployeeID:       1,
			NotificationType: model.NotificationSMS,
			Recipient:        "+1234567890",
			Message:          "Test notification message",
			Status:           model.NotificationPending,
		}

		// GORM allows creating notifications with non-existent foreign keys in SQLite
		err := repo.Create(notification)
		assert.NoError(t, err)
		assert.NotZero(t, notification.ID)
	})

	t.Run("returns error when database operation fails", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		notification := &model.Notification{
			UrgencyID:        1,
			EmployeeID:       1,
			NotificationType: model.NotificationSMS,
			Recipient:        "+1234567890",
			Message:          "Test notification message",
			Status:           model.NotificationPending,
		}

		err := repo.Create(notification)
		assert.Error(t, err)
	})
}

func TestNotificationRepository_GetByID(t *testing.T) {
	t.Parallel()

	t.Run("successfully retrieves notification by ID", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		urgency := createTestUrgencyForNotification(t, db)

		notification := &model.Notification{
			UrgencyID:        urgency.ID,
			EmployeeID:       1,
			NotificationType: model.NotificationSMS,
			Recipient:        "+1234567890",
			Message:          "Test notification message",
			Status:           model.NotificationPending,
		}
		err := db.Create(notification).Error
		require.NoError(t, err)

		var retrievedNotification model.Notification
		err = repo.GetByID(notification.ID, &retrievedNotification)
		assert.NoError(t, err)
		assert.Equal(t, notification.ID, retrievedNotification.ID)
		assert.Equal(t, notification.UrgencyID, retrievedNotification.UrgencyID)
		assert.Equal(t, notification.EmployeeID, retrievedNotification.EmployeeID)
		assert.Equal(t, notification.Message, retrievedNotification.Message)
	})

	t.Run("returns error when notification not found", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		var nonExistentNotification model.Notification
		err := repo.GetByID(999, &nonExistentNotification)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("returns error when database operation fails", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		var notification model.Notification
		err := repo.GetByID(1, &notification)
		assert.Error(t, err)
		assert.NotEqual(t, gorm.ErrRecordNotFound, err)
	})
}

func TestNotificationRepository_GetPendingNotifications(t *testing.T) {
	t.Parallel()

	t.Run("successfully retrieves pending notifications without limit", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		urgency := createTestUrgencyForNotification(t, db)

		pendingNotification1 := &model.Notification{
			UrgencyID:        urgency.ID,
			EmployeeID:       1,
			NotificationType: model.NotificationSMS,
			Recipient:        "+1234567890",
			Message:          "Pending notification 1",
			Status:           model.NotificationPending,
		}
		pendingNotification2 := &model.Notification{
			UrgencyID:        urgency.ID,
			EmployeeID:       2,
			NotificationType: model.NotificationEmail,
			Recipient:        "test@example.com",
			Message:          "Pending notification 2",
			Status:           model.NotificationPending,
		}
		sentNotification := &model.Notification{
			UrgencyID:        urgency.ID,
			EmployeeID:       3,
			NotificationType: model.NotificationSMS,
			Recipient:        "+1234567891",
			Message:          "Sent notification",
			Status:           model.NotificationSent,
		}

		err := db.Create(pendingNotification1).Error
		require.NoError(t, err)
		err = db.Create(pendingNotification2).Error
		require.NoError(t, err)
		err = db.Create(sentNotification).Error
		require.NoError(t, err)

		notifications, err := repo.GetPendingNotifications(0)
		assert.NoError(t, err)
		assert.Len(t, notifications, 2)
		for _, notification := range notifications {
			assert.Equal(t, model.NotificationPending, notification.Status)
			assert.NotNil(t, notification.Urgency)
		}
	})

	t.Run("successfully retrieves pending notifications with limit", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		urgency := createTestUrgencyForNotification(t, db)

		pendingNotification1 := &model.Notification{
			UrgencyID:        urgency.ID,
			EmployeeID:       1,
			NotificationType: model.NotificationSMS,
			Recipient:        "+1234567890",
			Message:          "Pending notification 1",
			Status:           model.NotificationPending,
		}
		pendingNotification2 := &model.Notification{
			UrgencyID:        urgency.ID,
			EmployeeID:       2,
			NotificationType: model.NotificationEmail,
			Recipient:        "test@example.com",
			Message:          "Pending notification 2",
			Status:           model.NotificationPending,
		}

		err := db.Create(pendingNotification1).Error
		require.NoError(t, err)
		err = db.Create(pendingNotification2).Error
		require.NoError(t, err)

		notifications, err := repo.GetPendingNotifications(1)
		assert.NoError(t, err)
		assert.Len(t, notifications, 1)
		assert.Equal(t, model.NotificationPending, notifications[0].Status)
	})

	t.Run("returns error when database operation fails", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		notifications, err := repo.GetPendingNotifications(0)
		assert.Error(t, err)
		assert.Nil(t, notifications)
	})
}

func TestNotificationRepository_GetByUrgencyID(t *testing.T) {
	t.Parallel()

	t.Run("successfully retrieves notifications by urgency ID", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		urgency1 := createTestUrgencyForNotification(t, db)
		urgency2 := &model.Urgency{
			FirstName:    "Marko",
			LastName:     "Markovic",
			Email:        "test2@example.com",
			ContactPhone: "987654321",
			Description:  "Test description 2",
			Level:        urgencyV1.Medium,
			Status:       urgencyV1.InProgress,
		}
		err := db.Create(urgency2).Error
		require.NoError(t, err)

		notification1 := &model.Notification{
			UrgencyID:        urgency1.ID,
			EmployeeID:       1,
			NotificationType: model.NotificationSMS,
			Recipient:        "+1234567890",
			Message:          "Notification for urgency 1",
			Status:           model.NotificationPending,
		}
		notification2 := &model.Notification{
			UrgencyID:        urgency1.ID,
			EmployeeID:       2,
			NotificationType: model.NotificationEmail,
			Recipient:        "test@example.com",
			Message:          "Another notification for urgency 1",
			Status:           model.NotificationSent,
		}
		notification3 := &model.Notification{
			UrgencyID:        urgency2.ID,
			EmployeeID:       1,
			NotificationType: model.NotificationSMS,
			Recipient:        "+1234567890",
			Message:          "Notification for urgency 2",
			Status:           model.NotificationPending,
		}

		err = db.Create(notification1).Error
		require.NoError(t, err)
		err = db.Create(notification2).Error
		require.NoError(t, err)
		err = db.Create(notification3).Error
		require.NoError(t, err)

		notifications, err := repo.GetByUrgencyID(urgency1.ID)
		assert.NoError(t, err)
		assert.Len(t, notifications, 2)

		notifications, err = repo.GetByUrgencyID(urgency2.ID)
		assert.NoError(t, err)
		assert.Len(t, notifications, 1)
	})

	t.Run("returns empty slice when urgency has no notifications", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		notifications, err := repo.GetByUrgencyID(999)
		assert.NoError(t, err)
		assert.Len(t, notifications, 0)
	})

	t.Run("returns error when database operation fails", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		notifications, err := repo.GetByUrgencyID(1)
		assert.Error(t, err)
		assert.Nil(t, notifications)
	})
}

func TestNotificationRepository_GetByEmployeeID(t *testing.T) {
	t.Parallel()

	t.Run("successfully retrieves notifications by employee ID", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		urgency := createTestUrgencyForNotification(t, db)

		notification1 := &model.Notification{
			UrgencyID:        urgency.ID,
			EmployeeID:       1,
			NotificationType: model.NotificationSMS,
			Recipient:        "+1234567890",
			Message:          "Notification for employee 1",
			Status:           model.NotificationPending,
		}
		notification2 := &model.Notification{
			UrgencyID:        urgency.ID,
			EmployeeID:       1,
			NotificationType: model.NotificationEmail,
			Recipient:        "test@example.com",
			Message:          "Another notification for employee 1",
			Status:           model.NotificationSent,
		}
		notification3 := &model.Notification{
			UrgencyID:        urgency.ID,
			EmployeeID:       2,
			NotificationType: model.NotificationSMS,
			Recipient:        "+1234567891",
			Message:          "Notification for employee 2",
			Status:           model.NotificationPending,
		}

		err := db.Create(notification1).Error
		require.NoError(t, err)
		err = db.Create(notification2).Error
		require.NoError(t, err)
		err = db.Create(notification3).Error
		require.NoError(t, err)

		notifications, err := repo.GetByEmployeeID(1)
		assert.NoError(t, err)
		assert.Len(t, notifications, 2)

		notifications, err = repo.GetByEmployeeID(2)
		assert.NoError(t, err)
		assert.Len(t, notifications, 1)
	})

	t.Run("returns empty slice when employee has no notifications", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		notifications, err := repo.GetByEmployeeID(999)
		assert.NoError(t, err)
		assert.Len(t, notifications, 0)
	})

	t.Run("returns error when database operation fails", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		notifications, err := repo.GetByEmployeeID(1)
		assert.Error(t, err)
		assert.Nil(t, notifications)
	})
}

func TestNotificationRepository_Update(t *testing.T) {
	t.Parallel()

	t.Run("successfully updates notification", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		urgency := createTestUrgencyForNotification(t, db)

		notification := &model.Notification{
			UrgencyID:        urgency.ID,
			EmployeeID:       1,
			NotificationType: model.NotificationSMS,
			Recipient:        "+1234567890",
			Message:          "Original message",
			Status:           model.NotificationPending,
		}
		err := db.Create(notification).Error
		require.NoError(t, err)

		notification.Message = "Updated message"
		notification.Status = model.NotificationSent

		err = repo.Update(notification)
		assert.NoError(t, err)

		var updatedNotification model.Notification
		err = db.First(&updatedNotification, notification.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, "Updated message", updatedNotification.Message)
		assert.Equal(t, model.NotificationSent, updatedNotification.Status)
	})

	t.Run("creates new notification when ID does not exist", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		notification := &model.Notification{
			ID:               999, // Non-existent notification
			UrgencyID:        1,
			EmployeeID:       1,
			NotificationType: model.NotificationSMS,
			Recipient:        "+1234567890",
			Message:          "Updated message",
			Status:           model.NotificationSent,
		}

		// GORM's Save() will create a new record if ID doesn't exist
		err := repo.Update(notification)
		assert.NoError(t, err)

		// Verify the notification was created with the specified ID
		var createdNotification model.Notification
		err = db.First(&createdNotification, 999).Error
		assert.NoError(t, err)
		assert.Equal(t, uint(999), createdNotification.ID)
	})

	t.Run("returns error when database operation fails", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		notification := &model.Notification{
			ID:               1,
			UrgencyID:        1,
			EmployeeID:       1,
			NotificationType: model.NotificationSMS,
			Recipient:        "+1234567890",
			Message:          "Updated message",
			Status:           model.NotificationSent,
		}

		err := repo.Update(notification)
		assert.Error(t, err)
	})
}

func TestNotificationRepository_Delete(t *testing.T) {
	t.Parallel()

	t.Run("successfully deletes notification", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		urgency := createTestUrgencyForNotification(t, db)

		notification := &model.Notification{
			UrgencyID:        urgency.ID,
			EmployeeID:       1,
			NotificationType: model.NotificationSMS,
			Recipient:        "+1234567890",
			Message:          "Test notification",
			Status:           model.NotificationPending,
		}
		err := db.Create(notification).Error
		require.NoError(t, err)

		err = repo.Delete(notification.ID)
		assert.NoError(t, err)

		var deletedNotification model.Notification
		err = db.First(&deletedNotification, notification.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("succeeds even when notification does not exist", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		// GORM's Delete() succeeds even if no records are deleted
		err := repo.Delete(999) // Non-existent notification
		assert.NoError(t, err)
	})

	t.Run("returns error when database operation fails", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		err := repo.Delete(1)
		assert.Error(t, err)
	})
}

func TestNotificationRepository_MarkAsSent(t *testing.T) {
	t.Parallel()

	t.Run("successfully marks notification as sent", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		urgency := createTestUrgencyForNotification(t, db)

		notification := &model.Notification{
			UrgencyID:        urgency.ID,
			EmployeeID:       1,
			NotificationType: model.NotificationSMS,
			Recipient:        "+1234567890",
			Message:          "Test notification",
			Status:           model.NotificationPending,
		}
		err := db.Create(notification).Error
		require.NoError(t, err)

		sentAt := time.Now()
		err = repo.MarkAsSent(notification.ID, sentAt)
		assert.NoError(t, err)

		var updatedNotification model.Notification
		err = db.First(&updatedNotification, notification.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, model.NotificationSent, updatedNotification.Status)
		assert.NotNil(t, updatedNotification.SentAt)
		assert.WithinDuration(t, sentAt, *updatedNotification.SentAt, time.Second)
	})

	t.Run("succeeds even when notification does not exist", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		// GORM's Update() succeeds even if no records are updated
		err := repo.MarkAsSent(999, time.Now()) // Non-existent notification
		assert.NoError(t, err)
	})

	t.Run("returns error when database operation fails", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		err := repo.MarkAsSent(1, time.Now())
		assert.Error(t, err)
	})
}

func TestNotificationRepository_MarkAsFailed(t *testing.T) {
	t.Parallel()

	t.Run("successfully marks notification as failed", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		urgency := createTestUrgencyForNotification(t, db)

		notification := &model.Notification{
			UrgencyID:        urgency.ID,
			EmployeeID:       1,
			NotificationType: model.NotificationSMS,
			Recipient:        "+1234567890",
			Message:          "Test notification",
			Status:           model.NotificationPending,
		}
		err := db.Create(notification).Error
		require.NoError(t, err)

		errorMessage := "SMS delivery failed"
		err = repo.MarkAsFailed(notification.ID, errorMessage)
		assert.NoError(t, err)

		var updatedNotification model.Notification
		err = db.First(&updatedNotification, notification.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, model.NotificationFailed, updatedNotification.Status)
		assert.Equal(t, errorMessage, updatedNotification.ErrorMessage)
	})

	t.Run("succeeds even when notification does not exist", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		// GORM's Update() succeeds even if no records are updated
		err := repo.MarkAsFailed(999, "Error message") // Non-existent notification
		assert.NoError(t, err)
	})

	t.Run("returns error when database operation fails", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		err := repo.MarkAsFailed(1, "Error message")
		assert.Error(t, err)
	})
}

func TestNotificationRepository_IncrementAttempts(t *testing.T) {
	t.Parallel()

	t.Run("successfully increments notification attempts", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		urgency := createTestUrgencyForNotification(t, db)

		notification := &model.Notification{
			UrgencyID:        urgency.ID,
			EmployeeID:       1,
			NotificationType: model.NotificationSMS,
			Recipient:        "+1234567890",
			Message:          "Test notification",
			Status:           model.NotificationPending,
			Attempts:         0,
		}
		err := db.Create(notification).Error
		require.NoError(t, err)

		err = repo.IncrementAttempts(notification.ID)
		assert.NoError(t, err)

		var updatedNotification model.Notification
		err = db.First(&updatedNotification, notification.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, 1, updatedNotification.Attempts)
		assert.NotNil(t, updatedNotification.LastAttemptAt)

		err = repo.IncrementAttempts(notification.ID)
		assert.NoError(t, err)

		err = db.First(&updatedNotification, notification.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, 2, updatedNotification.Attempts)
	})

	t.Run("succeeds even when notification does not exist", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		// GORM's Update() succeeds even if no records are updated
		err := repo.IncrementAttempts(999) // Non-existent notification
		assert.NoError(t, err)
	})

	t.Run("returns error when database operation fails", func(t *testing.T) {
		db := setupNotificationTestDB(t)
		log := utils.NewTestLogger()
		repo := NewNotificationRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		err := repo.IncrementAttempts(1)
		assert.Error(t, err)
	})
}

func setupNotificationTestDB(t *testing.T) *gorm.DB {
	sqlDB, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	db, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Urgency{}, &model.Notification{})
	require.NoError(t, err)

	return db
}

func createTestUrgencyForNotification(t *testing.T, db *gorm.DB) *model.Urgency {
	urgency := &model.Urgency{
		FirstName:    "Marko",
		LastName:     "Markovic",
		Email:        "test@example.com",
		ContactPhone: "123456789",
		Description:  "Test description",
		Level:        urgencyV1.High,
		Status:       urgencyV1.Open,
	}
	err := db.Create(urgency).Error
	require.NoError(t, err)
	return urgency
}
