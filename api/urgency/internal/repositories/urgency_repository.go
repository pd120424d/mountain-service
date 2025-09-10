package repositories

//go:generate mockgen -source=urgency_repository.go -destination=urgency_repository_gomock.go -package=repositories mountain_service/urgency/internal/repositories -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"context"
	"database/sql"
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
	ListUnassignedIDs(ctx context.Context) ([]uint, error)
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
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(log, "UrgencyRepository.Create")()
	return r.dbWrite.WithContext(ctx).Create(urgency).Error
}

func (r *urgencyRepository) GetAll(ctx context.Context) ([]model.Urgency, error) {
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(log, "UrgencyRepository.GetAll")()
	var urgencies []model.Urgency
	err := r.withRead(ctx, func(db *gorm.DB) error {
		return db.Where("deleted_at IS NULL").Find(&urgencies).Error
	})
	return urgencies, err
}

func (r *urgencyRepository) GetByID(ctx context.Context, id uint, urgency *model.Urgency) error {
	return r.withRead(ctx, func(db *gorm.DB) error {
		return db.First(urgency, "id = ?", id).Error
	})
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
	if err := r.withRead(ctx, func(db *gorm.DB) error {
		query := db.Model(&model.Urgency{})
		filterKeys := slices.Collect(maps.Keys(filters))
		slices.Sort(filterKeys)
		for _, key := range filterKeys {
			if _, ok := allowedColumns[key]; !ok {
				return fmt.Errorf("invalid filter key: %s", key)
			}
			value := filters[key]
			switch v := value.(type) {
			case string:
				query = query.Where(fmt.Sprintf("%s LIKE ?", key), fmt.Sprintf("%%%s%%", v))
			case int, int32, int64, float32, float64, bool:
				query = query.Where(fmt.Sprintf("%s = ?", key), v)
			default:
				return fmt.Errorf("unsupported type for filter key: %s", key)
			}
		}
		return query.Find(&urgencies).Error
	}); err != nil {
		return nil, err
	}
	return urgencies, nil
}

func (r *urgencyRepository) ListPaginated(ctx context.Context, page int, pageSize int, assignedEmployeeID *uint) ([]model.Urgency, int64, error) {
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(log, "UrgencyRepository.ListPaginated")()
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

	if err := r.withRead(ctx, func(db *gorm.DB) error {
		q := db.Model(&model.Urgency{}).Where("deleted_at IS NULL")
		if assignedEmployeeID != nil {
			q = q.Where("assigned_employee_id = ?", *assignedEmployeeID)
		}
		if err := q.Count(&total).Error; err != nil {
			return err
		}
		orderExpr := "sort_priority ASC, created_at DESC"
		return q.Order(orderExpr).Limit(pageSize).Offset(offset).Find(&urgencies).Error
	}); err != nil {
		return nil, 0, err
	}

	return urgencies, total, nil
}

func (r *urgencyRepository) ListUnassignedIDs(ctx context.Context) ([]uint, error) {
	log := r.log.WithContext(ctx)
	defer utils.TimeOperation(log, "UrgencyRepository.ListUnassignedIDs")()
	var ids []uint
	if err := r.withRead(ctx, func(db *gorm.DB) error {
		return db.Model(&model.Urgency{}).
			Where("deleted_at IS NULL AND sort_priority = ?", 1).
			Pluck("id", &ids).Error
	}); err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *urgencyRepository) ResetAllData(ctx context.Context) error {
	return r.dbWrite.WithContext(ctx).Unscoped().Delete(&model.Urgency{}, "1 = 1").Error
}

func (r *urgencyRepository) getReadDB(ctx context.Context) *gorm.DB {
	if utils.IsFreshRequired(ctx) {
		// Read-Your-Writes: route to primary within fresh window
		r.log.WithContext(ctx).Debugf("RYW: using primary for read")
		return r.dbWrite.WithContext(ctx)
	}
	return r.dbRead.WithContext(ctx)
}

// withRead executes the provided function using the replica in a read-only transaction when applicable.
// If RYW fresh window is active or a replica is not configured, it executes against primary without a transaction.
func (r *urgencyRepository) withRead(ctx context.Context, fn func(db *gorm.DB) error) error {
	// If fresh is required or read pool is the same as write, avoid read-only tx
	if utils.IsFreshRequired(ctx) || r.dbRead == nil || r.dbRead == r.dbWrite {
		return fn(r.dbWrite.WithContext(ctx))
	}
	// Begin a read-only transaction on the read pool
	tx := r.dbRead.WithContext(ctx).Begin(&sql.TxOptions{ReadOnly: true})
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
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
