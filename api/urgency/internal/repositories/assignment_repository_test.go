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

func TestAssignmentRepository_Create(t *testing.T) {
	t.Parallel()

	t.Run("successfully creates assignment", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		urgency := createTestUrgency(t, db)

		assignment := &model.EmergencyAssignment{
			UrgencyID:  urgency.ID,
			EmployeeID: 1,
			Status:     model.AssignmentPending,
			AssignedAt: time.Now(),
		}

		err := repo.Create(assignment)
		assert.NoError(t, err)
		assert.NotZero(t, assignment.ID)

		var dbAssignment model.EmergencyAssignment
		err = db.First(&dbAssignment, assignment.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, urgency.ID, dbAssignment.UrgencyID)
		assert.Equal(t, uint(1), dbAssignment.EmployeeID)
		assert.Equal(t, model.AssignmentPending, dbAssignment.Status)
	})

	t.Run("creates assignment even when urgency does not exist in database", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		assignment := &model.EmergencyAssignment{
			UrgencyID:  999, // Non-existent urgency
			EmployeeID: 1,
			Status:     model.AssignmentPending,
			AssignedAt: time.Now(),
		}

		// GORM allows creating assignments with non-existent foreign keys in SQLite
		err := repo.Create(assignment)
		assert.NoError(t, err)
		assert.NotZero(t, assignment.ID)
	})

	t.Run("returns error when database operation fails", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		assignment := &model.EmergencyAssignment{
			UrgencyID:  1,
			EmployeeID: 1,
			Status:     model.AssignmentPending,
			AssignedAt: time.Now(),
		}

		err := repo.Create(assignment)
		assert.Error(t, err)
	})
}

func TestAssignmentRepository_GetByID(t *testing.T) {
	t.Parallel()

	t.Run("successfully retrieves assignment by ID", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		urgency := createTestUrgency(t, db)

		assignment := &model.EmergencyAssignment{
			UrgencyID:  urgency.ID,
			EmployeeID: 1,
			Status:     model.AssignmentPending,
			AssignedAt: time.Now(),
		}
		err := db.Create(assignment).Error
		require.NoError(t, err)

		var retrievedAssignment model.EmergencyAssignment
		err = repo.GetByID(assignment.ID, &retrievedAssignment)
		assert.NoError(t, err)
		assert.Equal(t, assignment.ID, retrievedAssignment.ID)
		assert.Equal(t, assignment.UrgencyID, retrievedAssignment.UrgencyID)
		assert.Equal(t, assignment.EmployeeID, retrievedAssignment.EmployeeID)
		assert.NotNil(t, retrievedAssignment.Urgency) // Should be preloaded
	})

	t.Run("returns error when assignment not found", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		var nonExistentAssignment model.EmergencyAssignment
		err := repo.GetByID(999, &nonExistentAssignment)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("returns error when database operation fails", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		var assignment model.EmergencyAssignment
		err := repo.GetByID(1, &assignment)
		assert.Error(t, err)
		assert.NotEqual(t, gorm.ErrRecordNotFound, err)
	})
}

func TestAssignmentRepository_GetByUrgencyID(t *testing.T) {
	t.Parallel()

	t.Run("successfully retrieves assignments by urgency ID", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		urgency := createTestUrgency(t, db)

		assignment1 := &model.EmergencyAssignment{
			UrgencyID:  urgency.ID,
			EmployeeID: 1,
			Status:     model.AssignmentPending,
			AssignedAt: time.Now(),
		}
		assignment2 := &model.EmergencyAssignment{
			UrgencyID:  urgency.ID,
			EmployeeID: 2,
			Status:     model.AssignmentAccepted,
			AssignedAt: time.Now(),
		}
		err := db.Create(assignment1).Error
		require.NoError(t, err)
		err = db.Create(assignment2).Error
		require.NoError(t, err)

		assignments, err := repo.GetByUrgencyID(urgency.ID)
		assert.NoError(t, err)
		assert.Len(t, assignments, 2)
	})

	t.Run("returns empty slice when urgency has no assignments", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		assignments, err := repo.GetByUrgencyID(999)
		assert.NoError(t, err)
		assert.Len(t, assignments, 0)
	})

	t.Run("returns error when database operation fails", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		assignments, err := repo.GetByUrgencyID(1)
		assert.Error(t, err)
		assert.Nil(t, assignments)
	})
}

func TestAssignmentRepository_GetByEmployeeID(t *testing.T) {
	t.Parallel()

	t.Run("successfully retrieves assignments by employee ID", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		urgency1 := createTestUrgency(t, db)
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

		assignment1 := &model.EmergencyAssignment{
			UrgencyID:  urgency1.ID,
			EmployeeID: 1,
			Status:     model.AssignmentPending,
			AssignedAt: time.Now(),
		}
		assignment2 := &model.EmergencyAssignment{
			UrgencyID:  urgency2.ID,
			EmployeeID: 1,
			Status:     model.AssignmentAccepted,
			AssignedAt: time.Now(),
		}
		err = db.Create(assignment1).Error
		require.NoError(t, err)
		err = db.Create(assignment2).Error
		require.NoError(t, err)

		assignments, err := repo.GetByEmployeeID(1)
		assert.NoError(t, err)
		assert.Len(t, assignments, 2)
		for _, assignment := range assignments {
			assert.NotNil(t, assignment.Urgency)
		}
	})

	t.Run("returns empty slice when employee has no assignments", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		assignments, err := repo.GetByEmployeeID(999)
		assert.NoError(t, err)
		assert.Len(t, assignments, 0)
	})

	t.Run("returns error when database operation fails", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		assignments, err := repo.GetByEmployeeID(1)
		assert.Error(t, err)
		assert.Nil(t, assignments)
	})
}

func TestAssignmentRepository_GetPendingByEmployeeID(t *testing.T) {
	t.Parallel()

	t.Run("successfully retrieves pending assignments by employee ID", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		urgency := createTestUrgency(t, db)

		pendingAssignment := &model.EmergencyAssignment{
			UrgencyID:  urgency.ID,
			EmployeeID: 1,
			Status:     model.AssignmentPending,
			AssignedAt: time.Now(),
		}
		acceptedAssignment := &model.EmergencyAssignment{
			UrgencyID:  urgency.ID,
			EmployeeID: 1,
			Status:     model.AssignmentAccepted,
			AssignedAt: time.Now(),
		}
		err := db.Create(pendingAssignment).Error
		require.NoError(t, err)
		err = db.Create(acceptedAssignment).Error
		require.NoError(t, err)

		assignments, err := repo.GetPendingByEmployeeID(1)
		assert.NoError(t, err)
		assert.Len(t, assignments, 1)
		assert.Equal(t, model.AssignmentPending, assignments[0].Status)
		assert.NotNil(t, assignments[0].Urgency) // Should be preloaded
	})

	t.Run("returns empty slice when employee has no pending assignments", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		assignments, err := repo.GetPendingByEmployeeID(999)
		assert.NoError(t, err)
		assert.Len(t, assignments, 0)
	})

	t.Run("returns error when database operation fails", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		assignments, err := repo.GetPendingByEmployeeID(1)
		assert.Error(t, err)
		assert.Nil(t, assignments)
	})
}

func TestAssignmentRepository_Update(t *testing.T) {
	t.Parallel()

	t.Run("successfully updates assignment", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		urgency := createTestUrgency(t, db)

		assignment := &model.EmergencyAssignment{
			UrgencyID:  urgency.ID,
			EmployeeID: 1,
			Status:     model.AssignmentPending,
			AssignedAt: time.Now(),
		}
		err := db.Create(assignment).Error
		require.NoError(t, err)

		assignment.Status = model.AssignmentAccepted

		err = repo.Update(assignment)
		assert.NoError(t, err)

		var updatedAssignment model.EmergencyAssignment
		err = db.First(&updatedAssignment, assignment.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, model.AssignmentAccepted, updatedAssignment.Status)
	})

	t.Run("creates new assignment when ID does not exist", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		assignment := &model.EmergencyAssignment{
			ID:         999, // Non-existent assignment
			UrgencyID:  1,
			EmployeeID: 1,
			Status:     model.AssignmentAccepted,
			AssignedAt: time.Now(),
		}

		// GORM's Save() will create a new record if ID doesn't exist
		err := repo.Update(assignment)
		assert.NoError(t, err)

		// Verify the assignment was created with the specified ID
		var createdAssignment model.EmergencyAssignment
		err = db.First(&createdAssignment, 999).Error
		assert.NoError(t, err)
		assert.Equal(t, uint(999), createdAssignment.ID)
	})

	t.Run("returns error when database operation fails", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		assignment := &model.EmergencyAssignment{
			ID:         1,
			UrgencyID:  1,
			EmployeeID: 1,
			Status:     model.AssignmentAccepted,
			AssignedAt: time.Now(),
		}

		err := repo.Update(assignment)
		assert.Error(t, err)
	})
}

func TestAssignmentRepository_Delete(t *testing.T) {
	t.Parallel()

	t.Run("successfully deletes assignment", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		urgency := createTestUrgency(t, db)

		assignment := &model.EmergencyAssignment{
			UrgencyID:  urgency.ID,
			EmployeeID: 1,
			Status:     model.AssignmentPending,
			AssignedAt: time.Now(),
		}
		err := db.Create(assignment).Error
		require.NoError(t, err)

		err = repo.Delete(assignment.ID)
		assert.NoError(t, err)

		var deletedAssignment model.EmergencyAssignment
		err = db.First(&deletedAssignment, assignment.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("succeeds even when assignment does not exist", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		// GORM's Delete() succeeds even if no records are deleted
		err := repo.Delete(999) // Non-existent assignment
		assert.NoError(t, err)
	})

	t.Run("returns error when database operation fails", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		err := repo.Delete(1)
		assert.Error(t, err)
	})
}

func TestAssignmentRepository_GetByUrgencyAndEmployee(t *testing.T) {
	t.Parallel()

	t.Run("successfully retrieves assignment by urgency and employee", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		urgency := createTestUrgency(t, db)

		assignment := &model.EmergencyAssignment{
			UrgencyID:  urgency.ID,
			EmployeeID: 1,
			Status:     model.AssignmentPending,
			AssignedAt: time.Now(),
		}
		err := db.Create(assignment).Error
		require.NoError(t, err)

		retrievedAssignment, err := repo.GetByUrgencyAndEmployee(urgency.ID, 1)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedAssignment)
		assert.Equal(t, assignment.ID, retrievedAssignment.ID)
		assert.Equal(t, urgency.ID, retrievedAssignment.UrgencyID)
		assert.Equal(t, uint(1), retrievedAssignment.EmployeeID)
	})

	t.Run("returns nil when assignment not found by employee", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		urgency := createTestUrgency(t, db)

		retrievedAssignment, err := repo.GetByUrgencyAndEmployee(urgency.ID, 999)
		assert.NoError(t, err)
		assert.Nil(t, retrievedAssignment)
	})

	t.Run("returns nil when assignment not found by urgency", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		retrievedAssignment, err := repo.GetByUrgencyAndEmployee(999, 1)
		assert.NoError(t, err)
		assert.Nil(t, retrievedAssignment)
	})

	t.Run("returns error when database operation fails", func(t *testing.T) {
		db := setupAssignmentTestDB(t)
		log := utils.NewTestLogger()
		repo := NewAssignmentRepository(log, db)

		// Close the database to simulate a database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		retrievedAssignment, err := repo.GetByUrgencyAndEmployee(1, 1)
		assert.Error(t, err)
		assert.Nil(t, retrievedAssignment)
	})
}

func setupAssignmentTestDB(t *testing.T) *gorm.DB {
	sqlDB, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	db, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Urgency{}, &model.EmergencyAssignment{})
	require.NoError(t, err)

	return db
}

func createTestUrgency(t *testing.T, db *gorm.DB) *model.Urgency {
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
