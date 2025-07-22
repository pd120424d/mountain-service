package repositories

//go:generate mockgen -source=urgency_repository.go -destination=urgency_repository_gomock.go -package=repositories mountain_service/urgency/internal/repositories -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"fmt"
	"maps"
	"slices"

	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/pd120424d/mountain-service/api/urgency/internal/model"

	"gorm.io/gorm"
)

type UrgencyRepository interface {
	Create(urgency *model.Urgency) error
	GetAll() ([]model.Urgency, error)
	GetByID(id uint, urgency *model.Urgency) error
	Update(urgency *model.Urgency) error
	Delete(urgencyID uint) error
	List(filters map[string]interface{}) ([]model.Urgency, error)
	ResetAllData() error
}

type urgencyRepository struct {
	log utils.Logger
	db  *gorm.DB
}

func NewUrgencyRepository(log utils.Logger, db *gorm.DB) UrgencyRepository {
	return &urgencyRepository{log: log.WithName("urgencyRepository"), db: db}
}

func (r *urgencyRepository) Create(urgency *model.Urgency) error {
	return r.db.Create(urgency).Error
}

func (r *urgencyRepository) GetAll() ([]model.Urgency, error) {
	var urgencies []model.Urgency
	err := r.db.Where("deleted_at IS NULL").Find(&urgencies).Error
	return urgencies, err
}

func (r *urgencyRepository) GetByID(id uint, urgency *model.Urgency) error {
	return r.db.First(urgency, "id = ?", id).Error
}

func (r *urgencyRepository) Update(urgency *model.Urgency) error {
	return r.db.Save(urgency).Error
}

func (r *urgencyRepository) Delete(urgencyID uint) error {
	return r.db.Delete(&model.Urgency{}, urgencyID).Error
}

func (r *urgencyRepository) List(filters map[string]interface{}) ([]model.Urgency, error) {
	allowedColumns := r.allowedColumns()
	var urgencies []model.Urgency
	query := r.db.Model(&model.Urgency{})

	filterKeys := slices.Collect(maps.Keys(filters))
	slices.Sort(filterKeys)

	for _, key := range filterKeys {
		if _, ok := allowedColumns[key]; !ok {
			return nil, fmt.Errorf("invalid filter key: %s", key)
		}

		value := filters[key]

		switch v := value.(type) {
		case string:
			// Use LIKE for string fields
			query = query.Where(fmt.Sprintf("%s LIKE ?", key), fmt.Sprintf("%%%s%%", v))
		case int, int32, int64, float32, float64, bool:
			// Use exact match for non-string types
			query = query.Where(fmt.Sprintf("%s = ?", key), v)
		default:
			return nil, fmt.Errorf("unsupported type for filter key: %s", key)
		}
	}

	err := query.Find(&urgencies).Error
	return urgencies, err
}

func (r *urgencyRepository) ResetAllData() error {
	return r.db.Unscoped().Delete(&model.Urgency{}, "1 = 1").Error
}

func (r *urgencyRepository) allowedColumns() map[string]bool {
	return map[string]bool{
		"name":          true,
		"email":         true,
		"contact_phone": true,
		"description":   true,
		"level":         true,
		"status":        true,
	}
}
