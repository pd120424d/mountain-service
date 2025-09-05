package repositories

//go:generate mockgen -source=activity_repository.go -destination=activity_repository_gomock.go -package=repositories mountain_service/activity/internal/repositories -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pd120424d/mountain-service/api/activity/internal/model"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/models"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"gorm.io/gorm"
)

type ActivityRepository interface {
	Create(ctx context.Context, activity *model.Activity) error
	CreateWithOutbox(ctx context.Context, activity *model.Activity, event *models.OutboxEvent) error
	GetByID(ctx context.Context, id uint) (*model.Activity, error)
	List(ctx context.Context, filter *model.ActivityFilter) ([]model.Activity, int64, error)
	GetStats(ctx context.Context) (*model.ActivityStats, error)
	Delete(ctx context.Context, id uint) error
	ResetAllData(ctx context.Context) error
}

type activityRepository struct {
	log utils.Logger
	db  *gorm.DB
}

func NewActivityRepository(log utils.Logger, db *gorm.DB) ActivityRepository {
	return &activityRepository{log: log.WithName("activityRepository"), db: db}
}

func (r *activityRepository) Create(ctx context.Context, activity *model.Activity) error {
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(ctx, r.log, "ActivityRepository.Create")()
	log.Infof("Creating activity: employee_id=%d, urgency_id=%d", activity.EmployeeID, activity.UrgencyID)

	if err := r.db.WithContext(ctx).Create(activity).Error; err != nil {
		log.Errorf("Failed to create activity: %v", err)
		return fmt.Errorf("failed to create activity: %w", err)
	}

	log.Infof("Activity created successfully: id=%d", activity.ID)
	return nil
}

func (r *activityRepository) CreateWithOutbox(ctx context.Context, activity *model.Activity, event *models.OutboxEvent) error {
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(ctx, r.log, "ActivityRepository.CreateWithOutbox")()
	log.Infof("Creating activity with outbox event")
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(activity).Error; err != nil {
			log.Errorf("Failed to create activity: %v", err)
			return fmt.Errorf("failed to create activity: %w", err)
		}
		if event != nil {
			// Make sure that we have aggregate ID at this point since it won't be processed if it is 0
			event.AggregateID = fmt.Sprintf("activity-%d", activity.ID)
			// Try to update eventData.ActivityID if it's a valid ActivityEvent JSON
			var ev activityV1.ActivityEvent
			if json.Unmarshal([]byte(event.EventData), &ev) == nil {
				ev.ActivityID = activity.ID
				if b, mErr := json.Marshal(ev); mErr == nil {
					event.EventData = string(b)
				}
			}
			if err := tx.Create(event).Error; err != nil {
				log.Errorf("Failed to create outbox event: %v", err)
				return fmt.Errorf("failed to create outbox event: %w", err)
			}
		}
		return nil
	})
}

func (r *activityRepository) GetByID(ctx context.Context, id uint) (*model.Activity, error) {
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(ctx, r.log, "ActivityRepository.GetByID")()
	log.Infof("Getting activity by ID: %d", id)

	var activity model.Activity
	if err := r.db.WithContext(ctx).First(&activity, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warnf("Activity not found: %d", id)
			return nil, fmt.Errorf("activity not found")
		}
		log.Errorf("Failed to get activity %d: %v", id, err)
		return nil, fmt.Errorf("failed to get activity: %w", err)
	}

	log.Infof("Activity retrieved successfully: id=%d", id)
	return &activity, nil
}

func (r *activityRepository) List(ctx context.Context, filter *model.ActivityFilter) ([]model.Activity, int64, error) {
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(ctx, r.log, "ActivityRepository.List")()
	log.Infof("Listing activities with filter: page=%d, pageSize=%d", filter.Page, filter.PageSize)

	// Validation of filter just sets defaults if needed and doesn't return error
	_ = filter.Validate()

	query := r.db.WithContext(ctx).Model(&model.Activity{})

	// Apply filters
	if filter.EmployeeID != nil {
		query = query.Where("employee_id = ?", *filter.EmployeeID)
	}
	if filter.UrgencyID != nil {
		query = query.Where("urgency_id = ?", *filter.UrgencyID)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Errorf("Failed to count activities: %v", err)
		return nil, 0, fmt.Errorf("failed to count activities: %w", err)
	}

	// Get paginated results
	var activities []model.Activity
	if err := query.
		Order("created_at DESC").
		Offset(filter.GetOffset()).
		Limit(filter.GetLimit()).
		Find(&activities).Error; err != nil {
		log.Errorf("Failed to list activities: %v", err)
		return nil, 0, fmt.Errorf("failed to list activities: %w", err)
	}

	log.Infof("Activities listed successfully: count=%d, total=%d", len(activities), total)
	return activities, total, nil
}

func (r *activityRepository) GetStats(ctx context.Context) (*model.ActivityStats, error) {
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(ctx, r.log, "ActivityRepository.GetStats")()
	log.Info("Getting activity statistics")

	stats := &model.ActivityStats{}

	// Get total count
	if err := r.db.WithContext(ctx).Model(&model.Activity{}).Count(&stats.TotalActivities).Error; err != nil {
		log.Errorf("Failed to get total activities count: %v", err)
		return nil, fmt.Errorf("failed to get total activities count: %w", err)
	}

	// Level statistics removed - activities don't have levels

	// Get recent activities (last 10)
	if err := r.db.WithContext(ctx).Order("created_at DESC").Limit(10).Find(&stats.RecentActivities).Error; err != nil {
		log.Errorf("Failed to get recent activities: %v", err)
		return nil, fmt.Errorf("failed to get recent activities: %w", err)
	}

	// Get activities for different time periods
	now := time.Now()

	// Last 24 hours
	if err := r.db.WithContext(ctx).Model(&model.Activity{}).
		Where("created_at >= ?", now.Add(-24*time.Hour)).
		Count(&stats.ActivitiesLast24h).Error; err != nil {
		log.Errorf("Failed to get activities last 24h: %v", err)
		return nil, fmt.Errorf("failed to get activities last 24h: %w", err)
	}

	// Last 7 days
	if err := r.db.WithContext(ctx).Model(&model.Activity{}).
		Where("created_at >= ?", now.Add(-7*24*time.Hour)).
		Count(&stats.ActivitiesLast7Days).Error; err != nil {
		log.Errorf("Failed to get activities last 7 days: %v", err)
		return nil, fmt.Errorf("failed to get activities last 7 days: %w", err)
	}

	// Last 30 days
	if err := r.db.WithContext(ctx).Model(&model.Activity{}).
		Where("created_at >= ?", now.Add(-30*24*time.Hour)).
		Count(&stats.ActivitiesLast30Days).Error; err != nil {
		log.Errorf("Failed to get activities last 30 days: %v", err)
		return nil, fmt.Errorf("failed to get activities last 30 days: %w", err)
	}

	log.Info("Activity statistics retrieved successfully")
	return stats, nil
}

func (r *activityRepository) Delete(ctx context.Context, id uint) error {
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(ctx, r.log, "ActivityRepository.Delete")()
	log.Infof("Deleting activity: %d", id)

	result := r.db.WithContext(ctx).Delete(&model.Activity{}, id)
	if result.Error != nil {
		log.Errorf("Failed to delete activity %d: %v", id, result.Error)
		return fmt.Errorf("failed to delete activity: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Warnf("Activity not found for deletion: %d", id)
		return fmt.Errorf("activity not found")
	}

	log.Infof("Activity deleted successfully: %d", id)
	return nil
}

func (r *activityRepository) ResetAllData(ctx context.Context) error {
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(ctx, r.log, "ActivityRepository.ResetAllData")()
	log.Warn("Resetting all activity data")

	if err := r.db.WithContext(ctx).Exec("DELETE FROM activities").Error; err != nil {
		log.Errorf("Failed to reset activity data: %v", err)
		return fmt.Errorf("failed to reset activity data: %w", err)
	}

	log.Info("All activity data reset successfully")
	return nil
}
