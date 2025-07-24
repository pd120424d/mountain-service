package repositories

//go:generate mockgen -source=assignment_repository.go -destination=assignment_repository_gomock.go -package=repositories mountain_service/urgency/internal/repositories -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/pd120424d/mountain-service/api/urgency/internal/model"
	"gorm.io/gorm"
)

type AssignmentRepository interface {
	Create(assignment *model.EmergencyAssignment) error
	GetByID(id uint, assignment *model.EmergencyAssignment) error
	GetByUrgencyID(urgencyID uint) ([]model.EmergencyAssignment, error)
	GetByEmployeeID(employeeID uint) ([]model.EmergencyAssignment, error)
	GetPendingByEmployeeID(employeeID uint) ([]model.EmergencyAssignment, error)
	Update(assignment *model.EmergencyAssignment) error
	Delete(id uint) error
	GetByUrgencyAndEmployee(urgencyID, employeeID uint) (*model.EmergencyAssignment, error)
}

type assignmentRepository struct {
	log utils.Logger
	db  *gorm.DB
}

func NewAssignmentRepository(log utils.Logger, db *gorm.DB) AssignmentRepository {
	return &assignmentRepository{log: log.WithName("assignmentRepository"), db: db}
}

func (r *assignmentRepository) Create(assignment *model.EmergencyAssignment) error {
	r.log.Infof("Creating emergency assignment: urgencyID=%d, employeeID=%d", assignment.UrgencyID, assignment.EmployeeID)

	if err := r.db.Create(assignment).Error; err != nil {
		r.log.Errorf("Failed to create emergency assignment: %v", err)
		return err
	}

	r.log.Infof("Emergency assignment created successfully: id=%d", assignment.ID)
	return nil
}

func (r *assignmentRepository) GetByID(id uint, assignment *model.EmergencyAssignment) error {
	r.log.Infof("Getting emergency assignment by ID: %d", id)

	if err := r.db.Preload("Urgency").First(assignment, id).Error; err != nil {
		r.log.Errorf("Failed to get emergency assignment %d: %v", id, err)
		return err
	}

	return nil
}

func (r *assignmentRepository) GetByUrgencyID(urgencyID uint) ([]model.EmergencyAssignment, error) {
	r.log.Infof("Getting emergency assignments by urgency ID: %d", urgencyID)

	var assignments []model.EmergencyAssignment
	if err := r.db.Where("urgency_id = ?", urgencyID).Find(&assignments).Error; err != nil {
		r.log.Errorf("Failed to get emergency assignments by urgency ID %d: %v", urgencyID, err)
		return nil, err
	}

	return assignments, nil
}

func (r *assignmentRepository) GetByEmployeeID(employeeID uint) ([]model.EmergencyAssignment, error) {
	r.log.Infof("Getting emergency assignments by employee ID: %d", employeeID)

	var assignments []model.EmergencyAssignment
	if err := r.db.Preload("Urgency").Where("employee_id = ?", employeeID).Find(&assignments).Error; err != nil {
		r.log.Errorf("Failed to get emergency assignments by employee ID %d: %v", employeeID, err)
		return nil, err
	}

	return assignments, nil
}

func (r *assignmentRepository) GetPendingByEmployeeID(employeeID uint) ([]model.EmergencyAssignment, error) {
	r.log.Infof("Getting pending emergency assignments by employee ID: %d", employeeID)

	var assignments []model.EmergencyAssignment
	if err := r.db.Preload("Urgency").Where("employee_id = ? AND status = ?", employeeID, model.AssignmentPending).Find(&assignments).Error; err != nil {
		r.log.Errorf("Failed to get pending emergency assignments by employee ID %d: %v", employeeID, err)
		return nil, err
	}

	return assignments, nil
}

func (r *assignmentRepository) Update(assignment *model.EmergencyAssignment) error {
	r.log.Infof("Updating emergency assignment: %d", assignment.ID)

	if err := r.db.Save(assignment).Error; err != nil {
		r.log.Errorf("Failed to update emergency assignment %d: %v", assignment.ID, err)
		return err
	}

	r.log.Infof("Emergency assignment updated successfully: %d", assignment.ID)
	return nil
}

func (r *assignmentRepository) Delete(id uint) error {
	r.log.Infof("Deleting emergency assignment: %d", id)

	if err := r.db.Delete(&model.EmergencyAssignment{}, id).Error; err != nil {
		r.log.Errorf("Failed to delete emergency assignment %d: %v", id, err)
		return err
	}

	r.log.Infof("Emergency assignment deleted successfully: %d", id)
	return nil
}

func (r *assignmentRepository) GetByUrgencyAndEmployee(urgencyID, employeeID uint) (*model.EmergencyAssignment, error) {
	r.log.Infof("Getting emergency assignment by urgency %d and employee %d", urgencyID, employeeID)

	var assignment model.EmergencyAssignment
	if err := r.db.Where("urgency_id = ? AND employee_id = ?", urgencyID, employeeID).First(&assignment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.Errorf("Failed to get emergency assignment by urgency %d and employee %d: %v", urgencyID, employeeID, err)
		return nil, err
	}

	return &assignment, nil
}
