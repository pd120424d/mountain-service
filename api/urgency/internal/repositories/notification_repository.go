package repositories

//go:generate mockgen -source=notification_repository.go -destination=notification_repository_gomock.go -package=repositories mountain_service/urgency/internal/repositories -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"context"
	"time"

	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/pd120424d/mountain-service/api/urgency/internal/model"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *model.Notification) error
	GetByID(ctx context.Context, id uint, notification *model.Notification) error
	GetPendingNotifications(ctx context.Context, limit int) ([]model.Notification, error)
	GetByUrgencyID(ctx context.Context, urgencyID uint) ([]model.Notification, error)
	GetByEmployeeID(ctx context.Context, employeeID uint) ([]model.Notification, error)
	Update(ctx context.Context, notification *model.Notification) error
	Delete(ctx context.Context, id uint) error
	MarkAsSent(ctx context.Context, id uint, sentAt time.Time) error
	MarkAsFailed(ctx context.Context, id uint, errorMessage string) error
	IncrementAttempts(ctx context.Context, id uint) error
}

type notificationRepository struct {
	log utils.Logger
	db  *gorm.DB
}

func NewNotificationRepository(log utils.Logger, db *gorm.DB) NotificationRepository {
	return &notificationRepository{log: log.WithName("notificationRepository"), db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notification *model.Notification) error {
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(log, "NotificationRepository.Create")()
	log.Infof("Creating notification: urgencyID=%d, employeeID=%d, type=%s", notification.UrgencyID, notification.EmployeeID, notification.NotificationType)

	if err := r.db.Create(notification).Error; err != nil {
		log.Errorf("Failed to create notification: %v", err)
		return err
	}

	log.Infof("Notification created successfully: id=%d", notification.ID)
	return nil
}

func (r *notificationRepository) GetByID(ctx context.Context, id uint, notification *model.Notification) error {
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(log, "NotificationRepository.GetByID")()
	log.Infof("Getting notification by ID: %d", id)

	if err := r.db.Preload("Urgency").First(notification, id).Error; err != nil {
		log.Errorf("Failed to get notification %d: %v", id, err)
		return err
	}

	return nil
}

func (r *notificationRepository) GetPendingNotifications(ctx context.Context, limit int) ([]model.Notification, error) {
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(log, "NotificationRepository.GetPendingNotifications")()
	log.Infof("Getting pending notifications with limit: %d", limit)

	var notifications []model.Notification
	query := r.db.Where("status = ?", model.NotificationPending).
		Order("created_at ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Preload("Urgency").Find(&notifications).Error; err != nil {
		log.Errorf("Failed to get pending notifications: %v", err)
		return nil, err
	}

	log.Infof("Retrieved pending notifications: count=%d", len(notifications))
	return notifications, nil
}

func (r *notificationRepository) GetByUrgencyID(ctx context.Context, urgencyID uint) ([]model.Notification, error) {
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(log, "NotificationRepository.GetByUrgencyID")()
	log.Infof("Getting notifications by urgency ID: %d", urgencyID)

	var notifications []model.Notification
	if err := r.db.Where("urgency_id = ?", urgencyID).Find(&notifications).Error; err != nil {
		log.Errorf("Failed to get notifications by urgency ID %d: %v", urgencyID, err)
		return nil, err
	}

	return notifications, nil
}

func (r *notificationRepository) GetByEmployeeID(ctx context.Context, employeeID uint) ([]model.Notification, error) {
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(log, "NotificationRepository.GetByEmployeeID")()
	log.Infof("Getting notifications by employee ID: %d", employeeID)

	var notifications []model.Notification
	if err := r.db.Preload("Urgency").Where("employee_id = ?", employeeID).Find(&notifications).Error; err != nil {
		log.Errorf("Failed to get notifications by employee ID %d: %v", employeeID, err)
		return nil, err
	}

	return notifications, nil
}

func (r *notificationRepository) Update(ctx context.Context, notification *model.Notification) error {
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(log, "NotificationRepository.Update")()
	log.Infof("Updating notification: %d", notification.ID)

	if err := r.db.Save(notification).Error; err != nil {
		log.Errorf("Failed to update notification %d: %v", notification.ID, err)
		return err
	}

	log.Infof("Notification updated successfully: %d", notification.ID)
	return nil
}

func (r *notificationRepository) Delete(ctx context.Context, id uint) error {
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(log, "NotificationRepository.Delete")()
	log.Infof("Deleting notification: %d", id)

	if err := r.db.Delete(&model.Notification{}, id).Error; err != nil {
		log.Errorf("Failed to delete notification %d: %v", id, err)
		return err
	}

	log.Infof("Notification deleted successfully: %d", id)
	return nil
}

func (r *notificationRepository) MarkAsSent(ctx context.Context, id uint, sentAt time.Time) error {
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(log, "NotificationRepository.MarkAsSent")()
	log.Infof("Marking notification as sent: %d", id)

	updates := map[string]interface{}{
		"status":  model.NotificationSent,
		"sent_at": sentAt,
	}

	if err := r.db.Model(&model.Notification{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Errorf("Failed to mark notification %d as sent: %v", id, err)
		return err
	}

	log.Infof("Notification marked as sent successfully: %d", id)
	return nil
}

func (r *notificationRepository) MarkAsFailed(ctx context.Context, id uint, errorMessage string) error {
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(log, "NotificationRepository.MarkAsFailed")()
	log.Infof("Marking notification as failed: %d", id)

	updates := map[string]interface{}{
		"status":          model.NotificationFailed,
		"error_message":   errorMessage,
		"last_attempt_at": time.Now(),
	}

	if err := r.db.Model(&model.Notification{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Errorf("Failed to mark notification %d as failed: %v", id, err)
		return err
	}

	log.Infof("Notification marked as failed successfully: %d", id)
	return nil
}

func (r *notificationRepository) IncrementAttempts(ctx context.Context, id uint) error {
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(log, "NotificationRepository.IncrementAttempts")()
	log.Infof("Incrementing notification attempts: %d", id)

	if err := r.db.Model(&model.Notification{}).Where("id = ?", id).Updates(map[string]interface{}{
		"attempts":        gorm.Expr("attempts + 1"),
		"last_attempt_at": time.Now(),
	}).Error; err != nil {
		log.Errorf("Failed to increment notification attempts for %d: %v", id, err)
		return err
	}

	log.Infof("Notification attempts incremented successfully: %d", id)
	return nil
}
