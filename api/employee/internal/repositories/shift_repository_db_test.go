package repositories

import (
	"testing"
	"time"

	"github.com/pd120424d/mountain-service/api/employee/internal/model"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShiftRepository_GetOrCreateShift(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB := setupSQLiteTestDB(t)
	repo := NewShiftRepository(log, gormDB)

	t.Run("it creates a shift when it doesn't exist", func(t *testing.T) {
		shift, err := repo.GetOrCreateShift(time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC), 1)
		assert.NoError(t, err)
		assert.NotNil(t, shift)
	})

	t.Run("it retrieves a shift when it exists", func(t *testing.T) {
		gormDB.Create(&model.Shift{ShiftDate: time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC), ShiftType: 1})

		shift, err := repo.GetOrCreateShift(time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC), 1)
		assert.NoError(t, err)
		assert.NotNil(t, shift)
	})
}

func TestShiftRepository_AssignedToShift(t *testing.T) {
	log := utils.NewTestLogger()

	t.Run("it returns true when employee is already assigned to shift", func(t *testing.T) {
		gormDB := setupSQLiteTestDB(t)
		repo := NewShiftRepository(log, gormDB)
		gormDB.Create(&model.Employee{ID: 1, Username: "test-user", FirstName: "Bruce", LastName: "Lee", Email: "test-user@example.com", ProfileType: model.Medic})
		gormDB.Create(&model.Shift{ID: 1, ShiftDate: time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC), ShiftType: 1})
		gormDB.Create(&model.EmployeeShift{EmployeeID: 1, ShiftID: 1})

		assigned, err := repo.AssignedToShift(1, 1)
		assert.NoError(t, err)
		assert.True(t, assigned)
	})

	t.Run("it returns false when employee is not assigned to shift", func(t *testing.T) {
		gormDB := setupSQLiteTestDB(t)
		repo := NewShiftRepository(log, gormDB)
		gormDB.Create(&model.Employee{ID: 1, Username: "test-user", FirstName: "Bruce", LastName: "Lee", Email: "test-user@example.com", ProfileType: model.Medic})
		gormDB.Create(&model.Shift{ID: 1, ShiftDate: time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC), ShiftType: 1})

		assigned, err := repo.AssignedToShift(1, 1)
		assert.NoError(t, err)
		assert.False(t, assigned)
	})
}

func TestShiftRepository_CountAssignmentsByProfile(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB := setupSQLiteTestDB(t)
	repo := NewShiftRepository(log, gormDB)

	t.Run("it returns the count of assignments for a profile type", func(t *testing.T) {
		gormDB.Create(&model.Employee{ID: 1, Username: "test-user", FirstName: "Bruce", LastName: "Lee", Email: "test-user@example.com", ProfileType: model.Medic})
		gormDB.Create(&model.Shift{ID: 1, ShiftDate: time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC), ShiftType: 1})
		gormDB.Create(&model.EmployeeShift{EmployeeID: 1, ShiftID: 1})

		gormDB.Create(&model.Employee{ID: 2, Username: "test-user2", FirstName: "Bruce2", LastName: "Lee2", Email: "test-user2@example.com", ProfileType: model.Medic})
		gormDB.Create(&model.EmployeeShift{EmployeeID: 2, ShiftID: 1})

		count, err := repo.CountAssignmentsByProfile(1, "Medic")
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)
	})
}

func TestShiftRepository_CreateAssignment(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB := setupSQLiteTestDB(t)
	repo := NewShiftRepository(log, gormDB)

	t.Run("it creates an assignment when data is valid", func(t *testing.T) {

		assignmentID, err := repo.CreateAssignment(1, 1)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), assignmentID)
	})
}

func TestShiftRepository_GetShiftsByEmployeeID(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB := setupSQLiteTestDB(t)
	repo := NewShiftRepository(log, gormDB)

	t.Run("it returns shifts for an employee when they exist", func(t *testing.T) {
		gormDB.Create(&model.Employee{ID: 1, Username: "test-user", FirstName: "Bruce", LastName: "Lee", Email: "test-user@example.com", ProfileType: model.Medic})
		gormDB.Create(&model.Shift{ID: 1, ShiftDate: time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC), ShiftType: 1})
		gormDB.Create(&model.Shift{ID: 2, ShiftDate: time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC), ShiftType: 2})
		gormDB.Create(&model.EmployeeShift{EmployeeID: 1, ShiftID: 1})
		gormDB.Create(&model.EmployeeShift{EmployeeID: 1, ShiftID: 2})

		var shifts []model.Shift
		err := repo.GetShiftsByEmployeeID(1, &shifts)
		assert.NoError(t, err)
		assert.Len(t, shifts, 2)
	})
}

func TestShiftRepository_GetShiftAvailability(t *testing.T) {
	log := utils.NewTestLogger()

	t.Run("it returns shifts availability for a given date", func(t *testing.T) {
		gormDB := setupSQLiteTestDB(t)
		repo := NewShiftRepository(log, gormDB)

		availability, err := repo.GetShiftAvailability(time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC))
		assert.NoError(t, err)
		assert.NotNil(t, availability)
		assert.Equal(t, 3, len(availability))
		assert.Equal(t, 2, availability[1][model.Medic])
		assert.Equal(t, 4, availability[1][model.Technical])
	})

	t.Run("it returns shifts availability for a given date when there are assigned employees", func(t *testing.T) {
		gormDB := setupSQLiteTestDB(t)
		repo := NewShiftRepository(log, gormDB)

		tx := gormDB.Create(&model.Employee{ID: 1, Username: "test-user", FirstName: "Bruce", LastName: "Lee", Email: "test-user@example.com", ProfileType: model.Medic})
		require.NoError(t, tx.Error)
		tx = gormDB.Create(&model.Employee{ID: 2, Username: "jackiec", FirstName: "Jackie", LastName: "Chan", Email: "jackiec@example.com", ProfileType: model.Technical})
		require.NoError(t, tx.Error)
		tx = gormDB.Create(&model.Shift{ID: 1, ShiftDate: time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC), ShiftType: 1})
		require.NoError(t, tx.Error)
		tx = gormDB.Create(&model.EmployeeShift{EmployeeID: 1, ShiftID: 1})
		require.NoError(t, tx.Error)
		tx = gormDB.Create(&model.EmployeeShift{EmployeeID: 2, ShiftID: 1})
		require.NoError(t, tx.Error)

		availability, err := repo.GetShiftAvailability(time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC))
		assert.NoError(t, err)
		assert.NotNil(t, availability)
		assert.Equal(t, 3, len(availability))
		assert.Equal(t, 1, availability[1][model.Medic])
		assert.Equal(t, 3, availability[1][model.Technical])
	})
}
