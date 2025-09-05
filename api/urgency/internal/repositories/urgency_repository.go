package repositories

//go:generate mockgen -source=urgency_repository.go -destination=urgency_repository_gomock.go -package=repositories mountain_service/urgency/internal/repositories -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/pd120424d/mountain-service/api/urgency/internal/model"

	"gorm.io/gorm"
)

type UrgencyRepository interface {
	Create(ctx context.Context, urgency *model.Urgency) error
	GetAll(ctx context.Context) ([]model.Urgency, error)
	GetByID(ctx context.Context, id uint, urgency *model.Urgency) error
	GetByIDPrimary(ctx context.Context, id uint, urgency *model.Urgency) error
	Update(ctx context.Context, urgency *model.Urgency) error
	Delete(ctx context.Context, urgencyID uint) error
	ListPaginated(ctx context.Context, page int, pageSize int, assignedEmployeeID *uint) ([]model.Urgency, int64, error)

	List(ctx context.Context, filters map[string]interface{}) ([]model.Urgency, error)
	ResetAllData(ctx context.Context) error
}

type urgencyRepository struct {
	log     utils.Logger
	dbWrite *gorm.DB
	dbRead  *gorm.DB
}

func NewUrgencyRepository(log utils.Logger, db *gorm.DB) UrgencyRepository {
	return &urgencyRepository{log: log.WithName("urgencyRepository"), dbWrite: db, dbRead: db}
}

func NewUrgencyRepositoryRW(log utils.Logger, writeDB *gorm.DB, readDB *gorm.DB) UrgencyRepository {
	if readDB == nil {
		readDB = writeDB
	}
	return &urgencyRepository{log: log.WithName("urgencyRepository"), dbWrite: writeDB, dbRead: readDB}
}

func (r *urgencyRepository) Create(ctx context.Context, urgency *model.Urgency) error {
	defer utils.TimeOperation(ctx, r.log, "UrgencyRepository.Create")()
	return r.dbWrite.WithContext(ctx).Create(urgency).Error
}

func (r *urgencyRepository) GetAll(ctx context.Context) ([]model.Urgency, error) {
	var urgencies []model.Urgency
	err := r.dbRead.WithContext(ctx).Where("deleted_at IS NULL").Find(&urgencies).Error
	return urgencies, err
}

func (r *urgencyRepository) GetByID(ctx context.Context, id uint, urgency *model.Urgency) error {
	return r.dbRead.WithContext(ctx).First(urgency, "id = ?", id).Error
}

func (r *urgencyRepository) GetByIDPrimary(ctx context.Context, id uint, urgency *model.Urgency) error {
	return r.dbWrite.WithContext(ctx).First(urgency, "id = ?", id).Error
}

func (r *urgencyRepository) Update(ctx context.Context, urgency *model.Urgency) error {
	return r.dbWrite.WithContext(ctx).Save(urgency).Error
}

func (r *urgencyRepository) Delete(ctx context.Context, urgencyID uint) error {
	return r.dbWrite.WithContext(ctx).Delete(&model.Urgency{}, urgencyID).Error
}

func (r *urgencyRepository) List(ctx context.Context, filters map[string]interface{}) ([]model.Urgency, error) {
	allowedColumns := r.allowedColumns()
	var urgencies []model.Urgency
	query := r.dbRead.WithContext(ctx).Model(&model.Urgency{})

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

func (r *urgencyRepository) ListPaginated(ctx context.Context, page int, pageSize int, assignedEmployeeID *uint) ([]model.Urgency, int64, error) {
	defer utils.TimeOperation(ctx, r.log, "UrgencyRepository.ListPaginated")()
	var urgencies []model.Urgency
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 1000 {
		pageSize = 1000
	}

	offset := (page - 1) * pageSize

	q := r.dbRead.WithContext(ctx).Model(&model.Urgency{}).Where("deleted_at IS NULL")
	if assignedEmployeeID != nil {
		q = q.Where("assigned_employee_id = ?", *assignedEmployeeID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	orderExpr := "CASE " +
		"WHEN status = 'open' AND assigned_employee_id IS NULL THEN 0 " +
		"WHEN status = 'open' AND assigned_employee_id IS NOT NULL THEN 1 " +
		"WHEN status = 'in_progress' THEN 2 " +
		"WHEN status = 'resolved' THEN 3 " +
		"WHEN status = 'closed' THEN 4 " +
		"ELSE 5 END ASC, created_at DESC"
	if err := q.Order(orderExpr).Limit(pageSize).Offset(offset).Find(&urgencies).Error; err != nil {
		return nil, 0, err
	}
	return urgencies, total, nil
}

func (r *urgencyRepository) ResetAllData(ctx context.Context) error {
	return r.dbWrite.WithContext(ctx).Unscoped().Delete(&model.Urgency{}, "1 = 1").Error
}

func (r *urgencyRepository) allowedColumns() map[string]bool {
	return map[string]bool{
		"first_name":    true,
		"last_name":     true,
		"email":         true,
		"contact_phone": true,
		"description":   true,
		"level":         true,
		"status":        true,
	}
}
