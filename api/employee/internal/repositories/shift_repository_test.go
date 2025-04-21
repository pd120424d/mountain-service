package repositories

import (
	"regexp"
	"testing"
	"time"

	"github.com/pd120424d/mountain-service/api/shared/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/logger"
)

func TestShiftRepositoryMockDB_GetOrCreateShift(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewShiftRepository(log, gormDB)
	gormDB.Logger = gormDB.Logger.LogMode(logger.Info)

	t.Run("it fails to create a shift when the query fails", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "shifts"`).
			WithArgs(time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC), 1, 1).
			WillReturnError(sqlmock.ErrCancelled)

		_, err := repo.GetOrCreateShift(time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC), 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find or create shift: canceling query due to user request")
	})
}

func TestShiftRepositoryMockDB_AssignedToShift(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewShiftRepository(log, gormDB)
	gormDB.Logger = gormDB.Logger.LogMode(logger.Info)

	t.Run("it fails to check assignment when the query fails", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "employee_shifts"`).
			WithArgs(1, 1, 1).
			WillReturnError(sqlmock.ErrCancelled)

		_, err := repo.AssignedToShift(1, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to check assignment: canceling query due to user request")
	})
}

func TestShiftRepositoryMockDB_CountAssignmentsByProfile(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewShiftRepository(log, gormDB)
	gormDB.Logger = gormDB.Logger.LogMode(logger.Info)

	t.Run("it fails to count assignments when the query fails", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM "employee_shifts"`).
			WithArgs(1, "Medic").
			WillReturnError(sqlmock.ErrCancelled)

		_, err := repo.CountAssignmentsByProfile(1, "Medic")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to count assignments: canceling query due to user request")
	})
}

func TestShiftRepositoryMockDB_CreateAssignment(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewShiftRepository(log, gormDB)
	gormDB.Logger = gormDB.Logger.LogMode(logger.Info)

	t.Run("it fails to create an assignment when the query fails", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "employee_shifts"`).
			WithArgs(1, 2).
			WillReturnError(sqlmock.ErrCancelled)
		mock.ExpectRollback()

		_, err := repo.CreateAssignment(1, 2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create assignment: canceling query due to user request")
	})
}

func TestShiftRepositoryMockDB_GetShiftAvailability(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewShiftRepository(log, gormDB)
	gormDB.Logger = gormDB.Logger.LogMode(logger.Info)

	t.Run("it fails to get shifts availability when the query fails", func(t *testing.T) {

		expectedSQL := `SELECT shifts.shift_type, employees.profile_type AS employee_role, COUNT(*) AS count FROM "shifts" JOIN employee_shifts ON shifts.id = employee_shifts.shift_id JOIN employees ON employee_shifts.employee_id = employees.id WHERE shift_date >= $1 AND shift_date < $2 GROUP BY shifts.shift_type, employees.profile_type`
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(
				time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 2, 4, 0, 0, 0, 0, time.UTC), // You had the same day twice before!
			).
			WillReturnError(sqlmock.ErrCancelled)

		_, err := repo.GetShiftAvailability(time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "canceling query due to user request")
	})
}

func TestShiftRepositoryMockDB_RemoveEmployeeFromShift(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewShiftRepository(log, gormDB)
	gormDB.Logger = gormDB.Logger.LogMode(logger.Info)

	t.Run("it fails to remove an employee from a shift when the delete query fails", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "employee_shifts"`).
			WithArgs(int64(123), int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "employee_shifts"`).
			WithArgs(123).
			WillReturnError(sqlmock.ErrCancelled)
		mock.ExpectRollback()

		err := repo.RemoveEmployeeFromShift(123)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "canceling query due to user request")
	})
}
