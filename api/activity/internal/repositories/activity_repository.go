package repositories

//go:generate mockgen -source=activity_repository.go -destination=activity_repository_gomock.go -package=repositories mountain_service/activity/internal/repositories -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"fmt"
	"time"

	"github.com/pd120424d/mountain-service/api/activity/internal/model"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"gorm.io/gorm"
)

type ActivityRepository interface {
	Create(activity *model.Activity) error
	GetByID(id uint) (*model.Activity, error)
	List(filter *model.ActivityFilter) ([]model.Activity, int64, error)
	GetStats() (*model.ActivityStats, error)
	Delete(id uint) error
	ResetAllData() error
}

type activityRepository struct {
	log utils.Logger
	db  *gorm.DB
}

func NewActivityRepository(log utils.Logger, db *gorm.DB) ActivityRepository {
	return &activityRepository{log: log.WithName("activityRepository"), db: db}
}

func (r *activityRepository) Create(activity *model.Activity) error {
	r.log.Infof("Creating activity: type=%s, level=%s, title=%s", activity.Type, activity.Level, activity.Title)

	if err := r.db.Create(activity).Error; err != nil {
		r.log.Errorf("Failed to create activity: %v", err)
		return fmt.Errorf("failed to create activity: %w", err)
	}

	r.log.Infof("Activity created successfully: id=%d", activity.ID)
	return nil
}

func (r *activityRepository) GetByID(id uint) (*model.Activity, error) {
	r.log.Infof("Getting activity by ID: %d", id)

	var activity model.Activity
	if err := r.db.First(&activity, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.Warnf("Activity not found: %d", id)
			return nil, fmt.Errorf("activity not found")
		}
		r.log.Errorf("Failed to get activity %d: %v", id, err)
		return nil, fmt.Errorf("failed to get activity: %w", err)
	}

	r.log.Infof("Activity retrieved successfully: id=%d", id)
	return &activity, nil
}

func (r *activityRepository) List(filter *model.ActivityFilter) ([]model.Activity, int64, error) {
	r.log.Infof("Listing activities with filter: page=%d, pageSize=%d", filter.Page, filter.PageSize)

	// Validate filter
	if err := filter.Validate(); err != nil {
		return nil, 0, fmt.Errorf("invalid filter: %w", err)
	}

	query := r.db.Model(&model.Activity{})

	// Apply filters
	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}
	if filter.Level != nil {
		query = query.Where("level = ?", *filter.Level)
	}
	if filter.ActorID != nil {
		query = query.Where("actor_id = ?", *filter.ActorID)
	}
	if filter.TargetID != nil {
		query = query.Where("target_id = ?", *filter.TargetID)
	}
	if filter.TargetType != nil {
		query = query.Where("target_type = ?", *filter.TargetType)
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
		r.log.Errorf("Failed to count activities: %v", err)
		return nil, 0, fmt.Errorf("failed to count activities: %w", err)
	}

	// Get paginated results
	var activities []model.Activity
	if err := query.
		Order("created_at DESC").
		Offset(filter.GetOffset()).
		Limit(filter.GetLimit()).
		Find(&activities).Error; err != nil {
		r.log.Errorf("Failed to list activities: %v", err)
		return nil, 0, fmt.Errorf("failed to list activities: %w", err)
	}

	r.log.Infof("Activities listed successfully: count=%d, total=%d", len(activities), total)
	return activities, total, nil
}

func (r *activityRepository) GetStats() (*model.ActivityStats, error) {
	r.log.Info("Getting activity statistics")

	stats := &model.ActivityStats{
		ActivitiesByType:  make(map[activityV1.ActivityType]int64),
		ActivitiesByLevel: make(map[activityV1.ActivityLevel]int64),
	}

	// Get total count
	if err := r.db.Model(&model.Activity{}).Count(&stats.TotalActivities).Error; err != nil {
		r.log.Errorf("Failed to get total activities count: %v", err)
		return nil, fmt.Errorf("failed to get total activities count: %w", err)
	}

	// Get activities by type
	var typeStats []struct {
		Type  activityV1.ActivityType `json:"type"`
		Count int64                   `json:"count"`
	}
	if err := r.db.Model(&model.Activity{}).
		Select("type, COUNT(*) as count").
		Group("type").
		Find(&typeStats).Error; err != nil {
		r.log.Errorf("Failed to get activities by type: %v", err)
		return nil, fmt.Errorf("failed to get activities by type: %w", err)
	}
	for _, stat := range typeStats {
		stats.ActivitiesByType[stat.Type] = stat.Count
	}

	// Get activities by level
	var levelStats []struct {
		Level activityV1.ActivityLevel `json:"level"`
		Count int64                    `json:"count"`
	}
	if err := r.db.Model(&model.Activity{}).
		Select("level, COUNT(*) as count").
		Group("level").
		Find(&levelStats).Error; err != nil {
		r.log.Errorf("Failed to get activities by level: %v", err)
		return nil, fmt.Errorf("failed to get activities by level: %w", err)
	}
	for _, stat := range levelStats {
		stats.ActivitiesByLevel[stat.Level] = stat.Count
	}

	// Get recent activities (last 10)
	if err := r.db.Order("created_at DESC").Limit(10).Find(&stats.RecentActivities).Error; err != nil {
		r.log.Errorf("Failed to get recent activities: %v", err)
		return nil, fmt.Errorf("failed to get recent activities: %w", err)
	}

	// Get activities for different time periods
	now := time.Now()

	// Last 24 hours
	if err := r.db.Model(&model.Activity{}).
		Where("created_at >= ?", now.Add(-24*time.Hour)).
		Count(&stats.ActivitiesLast24h).Error; err != nil {
		r.log.Errorf("Failed to get activities last 24h: %v", err)
		return nil, fmt.Errorf("failed to get activities last 24h: %w", err)
	}

	// Last 7 days
	if err := r.db.Model(&model.Activity{}).
		Where("created_at >= ?", now.Add(-7*24*time.Hour)).
		Count(&stats.ActivitiesLast7Days).Error; err != nil {
		r.log.Errorf("Failed to get activities last 7 days: %v", err)
		return nil, fmt.Errorf("failed to get activities last 7 days: %w", err)
	}

	// Last 30 days
	if err := r.db.Model(&model.Activity{}).
		Where("created_at >= ?", now.Add(-30*24*time.Hour)).
		Count(&stats.ActivitiesLast30Days).Error; err != nil {
		r.log.Errorf("Failed to get activities last 30 days: %v", err)
		return nil, fmt.Errorf("failed to get activities last 30 days: %w", err)
	}

	r.log.Info("Activity statistics retrieved successfully")
	return stats, nil
}

func (r *activityRepository) Delete(id uint) error {
	r.log.Infof("Deleting activity: %d", id)

	result := r.db.Delete(&model.Activity{}, id)
	if result.Error != nil {
		r.log.Errorf("Failed to delete activity %d: %v", id, result.Error)
		return fmt.Errorf("failed to delete activity: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		r.log.Warnf("Activity not found for deletion: %d", id)
		return fmt.Errorf("activity not found")
	}

	r.log.Infof("Activity deleted successfully: %d", id)
	return nil
}

func (r *activityRepository) ResetAllData() error {
	r.log.Warn("Resetting all activity data")

	if err := r.db.Exec("DELETE FROM activities").Error; err != nil {
		r.log.Errorf("Failed to reset activity data: %v", err)
		return fmt.Errorf("failed to reset activity data: %w", err)
	}

	r.log.Info("All activity data reset successfully")
	return nil
}
