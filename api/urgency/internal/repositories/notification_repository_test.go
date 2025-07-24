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
		Name:         "Test Urgency",
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

func TestNotificationRepository_Create(t *testing.T) {
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
}

func TestNotificationRepository_GetByID(t *testing.T) {
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

	var nonExistentNotification model.Notification
	err = repo.GetByID(999, &nonExistentNotification)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestNotificationRepository_GetPendingNotifications(t *testing.T) {
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

	notifications, err = repo.GetPendingNotifications(1)
	assert.NoError(t, err)
	assert.Len(t, notifications, 1)
	assert.Equal(t, model.NotificationPending, notifications[0].Status)
}

func TestNotificationRepository_GetByUrgencyID(t *testing.T) {
	db := setupNotificationTestDB(t)
	log := utils.NewTestLogger()
	repo := NewNotificationRepository(log, db)

	urgency1 := createTestUrgencyForNotification(t, db)
	urgency2 := &model.Urgency{
		Name:         "Test Urgency 2",
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

	notifications, err = repo.GetByUrgencyID(999)
	assert.NoError(t, err)
	assert.Len(t, notifications, 0)
}

func TestNotificationRepository_GetByEmployeeID(t *testing.T) {
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

	notifications, err = repo.GetByEmployeeID(999)
	assert.NoError(t, err)
	assert.Len(t, notifications, 0)
}

func TestNotificationRepository_Update(t *testing.T) {
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
}

func TestNotificationRepository_Delete(t *testing.T) {
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
}

func TestNotificationRepository_MarkAsSent(t *testing.T) {
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
}

func TestNotificationRepository_MarkAsFailed(t *testing.T) {
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
}

func TestNotificationRepository_IncrementAttempts(t *testing.T) {
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
}
