package repositories

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pd120424d/mountain-service/api/employee/internal/model"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
)

func TestShiftRepository_GetOrCreateShift(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewShiftRepository(log, gormDB)

	t.Run("it creates a shift when it doesn't exist", func(t *testing.T) {
		shiftDate := time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC)
		shiftType := 1

		mock.ExpectQuery(`SELECT \* FROM "shifts" WHERE "shifts"\."shift_date" = \$1 AND "shifts"\."shift_type" = \$2 ORDER BY "shifts"\."id" LIMIT \$3`).
			WithArgs(shiftDate, shiftType, 1).
			WillReturnRows(sqlmock.NewRows([]string{}))

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "shifts"`).
			WithArgs(shiftDate, shiftType, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		shift, err := repo.GetOrCreateShift(time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC), 1)
		assert.NoError(t, err)
		assert.NotNil(t, shift)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("it returns an error when it fails to create a shift", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "shifts" WHERE "shifts"\."shift_date" = \$1 AND "shifts"\."shift_type" = \$2 ORDER BY "shifts"\."id" LIMIT \$3`).
			WithArgs(time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC), 1, 1).
			WillReturnRows(sqlmock.NewRows([]string{}))

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "shifts"`).
			WithArgs(time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC), 1, sqlmock.AnyArg()).
			WillReturnError(sqlmock.ErrCancelled)
		mock.ExpectRollback()

		_, err := repo.GetOrCreateShift(time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC), 1)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("it retrieves a shift when it exists", func(t *testing.T) {
		shiftDate := time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC)
		shiftType := 1

		mock.ExpectQuery(`SELECT \* FROM "shifts" WHERE "shifts"\."shift_date" = \$1 AND "shifts"\."shift_type" = \$2 ORDER BY "shifts"\."id" LIMIT \$3`).
			WithArgs(shiftDate, shiftType, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		shift, err := repo.GetOrCreateShift(time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC), 1)
		assert.NoError(t, err)
		assert.NotNil(t, shift)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestShiftRepository_AssignedToShift(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewShiftRepository(log, gormDB)

	t.Run("it returns true when employee is already assigned to shift", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "employee_shifts" WHERE employee_id = \$1 AND shift_id = \$2 ORDER BY "employee_shifts"\."id" LIMIT \$3`).
			WithArgs(1, 1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		assigned, err := repo.AssignedToShift(1, 1)
		assert.NoError(t, err)
		assert.True(t, assigned)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("it returns false when employee is not assigned to shift", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "employee_shifts" WHERE employee_id = \$1 AND shift_id = \$2 ORDER BY "employee_shifts"\."id" LIMIT \$3`).
			WithArgs(1, 1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		assigned, err := repo.AssignedToShift(1, 1)
		assert.NoError(t, err)
		assert.False(t, assigned)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("it returns an error when it fails to check assignment", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "employee_shifts" WHERE employee_id = \$1 AND shift_id = \$2 ORDER BY "employee_shifts"\."id" LIMIT \$3`).
			WithArgs(1, 1, 1).
			WillReturnError(sqlmock.ErrCancelled)

		assigned, err := repo.AssignedToShift(1, 1)
		assert.Error(t, err)
		assert.False(t, assigned)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestShiftRepository_CountAssignmentsByProfile(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewShiftRepository(log, gormDB)

	t.Run("it returns the count of assignments for a profile type", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM "employee_shifts" WHERE shift_id = \$1 AND profile_type = \$2`).
			WithArgs(1, "Medic").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		count, err := repo.CountAssignmentsByProfile(1, "Medic")
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("it returns an error when it fails to count assignments", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM "employee_shifts" WHERE shift_id = \$1 AND profile_type = \$2`).
			WithArgs(1, "Medic").
			WillReturnError(sqlmock.ErrCancelled)

		count, err := repo.CountAssignmentsByProfile(1, "Medic")
		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestShiftRepository_CreateAssignment(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewShiftRepository(log, gormDB)

	t.Run("it creates an assignment when data is valid", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "employee_shifts"`).
			WithArgs(1, 1, "Medic").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		assignmentID, err := repo.CreateAssignment(1, 1, "Medic")
		assert.NoError(t, err)
		assert.Equal(t, uint(1), assignmentID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("it returns an error when it fails to create an assignment", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "employee_shifts"`).
			WithArgs(1, 1, "Medic").
			WillReturnError(sqlmock.ErrCancelled)
		mock.ExpectRollback()

		_, err := repo.CreateAssignment(1, 1, "Medic")
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestShiftRepository_GetShiftsByEmployeeID(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewShiftRepository(log, gormDB)

	t.Run("it returns shifts for an employee when they exist", func(t *testing.T) {
		mock.ExpectQuery(`SELECT employee_shifts.id, shifts.shift_date, shifts.shift_type, employee_shifts.profile_type FROM "employee_shifts" JOIN shifts ON employee_shifts.shift_id = shifts.id WHERE employee_shifts.employee_id = \$1 ORDER BY shifts.shift_date ASC, shifts.shift_type ASC`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "shift_date", "shift_type", "profile_type"}).
				AddRow(1, time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC), 1, "Medic").
				AddRow(2, time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC), 2, "Technical"))

		var shifts []model.Shift
		err := repo.GetShiftsByEmployeeID(1, &shifts)
		assert.NoError(t, err)
		assert.Len(t, shifts, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("it returns an error when it fails to retrieve shifts", func(t *testing.T) {
		mock.ExpectQuery(`SELECT employee_shifts.id, shifts.shift_date, shifts.shift_type, employee_shifts.profile_type FROM "employee_shifts" JOIN shifts ON employee_shifts.shift_id = shifts.id WHERE employee_shifts.employee_id = \$1 ORDER BY shifts.shift_date ASC, shifts.shift_type ASC`).
			WithArgs(1).
			WillReturnError(sqlmock.ErrCancelled)

		var shifts []model.Shift
		err := repo.GetShiftsByEmployeeID(1, &shifts)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
