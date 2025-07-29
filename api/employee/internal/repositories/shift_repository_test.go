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

		expectedSQL := `SELECT shifts.shift_date, shifts.shift_type, employees.profile_type AS employee_role, COUNT(*) AS count FROM "shifts" JOIN employee_shifts ON shifts.id = employee_shifts.shift_id JOIN employees ON employee_shifts.employee_id = employees.id WHERE shift_date >= $1 AND shift_date < $2 GROUP BY shifts.shift_date, shifts.shift_type, employees.profile_type`
		mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(
				time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 2, 10, 0, 0, 0, 0, time.UTC), // You had the same day twice before!
			).
			WillReturnError(sqlmock.ErrCancelled)

		_, err := repo.GetShiftAvailability(time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC), time.Date(2025, 2, 10, 0, 0, 0, 0, time.UTC))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "canceling query due to user request")
	})
}

func TestShiftRepositoryMockDB_RemoveEmployeeFromShiftByDetails(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewShiftRepository(log, gormDB)
	gormDB.Logger = gormDB.Logger.LogMode(logger.Info)

	t.Run("it fails when shift is not found", func(t *testing.T) {
		shiftDate := time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC)

		mock.ExpectQuery(`SELECT \* FROM "shifts"`).
			WithArgs(shiftDate, 1, 1).
			WillReturnError(sqlmock.ErrCancelled)

		err := repo.RemoveEmployeeFromShiftByDetails(1, shiftDate, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find shift")
	})

	t.Run("it fails when assignment is not found", func(t *testing.T) {
		shiftDate := time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC)

		mock.ExpectQuery(`SELECT \* FROM "shifts"`).
			WithArgs(shiftDate, 1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "shift_date", "shift_type", "created_at"}).
				AddRow(1, shiftDate, 1, time.Now()))

		mock.ExpectQuery(`SELECT \* FROM "employee_shifts"`).
			WithArgs(1, 1, 1).
			WillReturnError(sqlmock.ErrCancelled)

		err := repo.RemoveEmployeeFromShiftByDetails(1, shiftDate, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find assignment")
	})
	t.Run("it fails to remove an employee from a shift when the delete query fails", func(t *testing.T) {
		shiftDate := time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC)

		mock.ExpectQuery(`SELECT \* FROM "shifts"`).
			WithArgs(shiftDate, 1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "shift_date", "shift_type", "created_at"}).
				AddRow(1, shiftDate, 1, time.Now()))

		mock.ExpectQuery(`SELECT \* FROM "employee_shifts"`).
			WithArgs(1, 1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "employee_id", "shift_id"}).
				AddRow(1, 1, 1))

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "employee_shifts"`).
			WithArgs(1).
			WillReturnError(sqlmock.ErrCancelled)
		mock.ExpectRollback()

		err := repo.RemoveEmployeeFromShiftByDetails(1, shiftDate, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "canceling query due to user request")
	})
}

func TestShiftRepositoryMockDB_GetOnCallEmployees(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewShiftRepository(log, gormDB)
	gormDB.Logger = gormDB.Logger.LogMode(logger.Info)

	t.Run("it fails to get on-call employees when the query fails", func(t *testing.T) {
		mock.ExpectQuery(`SELECT DISTINCT employees\.\* FROM "employees" JOIN employee_shifts ON employees\.id = employee_shifts\.employee_id JOIN shifts ON employee_shifts\.shift_id = shifts\.id WHERE \(\(shifts\.shift_date = \$1 AND shifts\.shift_type = \$2\)\) AND "employees"\."deleted_at" IS NULL`).
			WillReturnError(sqlmock.ErrCancelled)

		testTime := time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC) // 10 AM, should be shift 1
		_, err := repo.GetOnCallEmployees(testTime, 0)            // No buffer to keep query simple
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get on-call employees: canceling query due to user request")
	})

	t.Run("it successfully returns on-call employees", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "password", "first_name", "last_name", "gender", "phone", "email", "profile_picture", "profile_type"}).
			AddRow(1, time.Now(), time.Now(), nil, "petar_petrovic", "hashed_password", "Petar", "Petrovic", "M", "123456789", "petar@example.com", "", "Medic").
			AddRow(2, time.Now(), time.Now(), nil, "marko_markovic", "hashed_password", "Marko", "Markovic", "F", "987654321", "marko@example.com", "", "Technical")

		mock.ExpectQuery(`SELECT DISTINCT employees\.\* FROM "employees" JOIN employee_shifts ON employees\.id = employee_shifts\.employee_id JOIN shifts ON employee_shifts\.shift_id = shifts\.id WHERE \(\(shifts\.shift_date = \$1 AND shifts\.shift_type = \$2\)\) AND "employees"\."deleted_at" IS NULL`).
			WillReturnRows(rows)

		testTime := time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC) // 10 AM, should be shift 1
		employees, err := repo.GetOnCallEmployees(testTime, 0)    // No buffer

		assert.NoError(t, err)
		assert.Len(t, employees, 2)
		assert.Equal(t, "petar_petrovic", employees[0].Username)
		assert.Equal(t, "marko_markovic", employees[1].Username)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("it includes next shift employees when within buffer", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "password", "first_name", "last_name", "gender", "phone", "email", "profile_picture", "profile_type"}).
			AddRow(1, time.Now(), time.Now(), nil, "current_shift", "hashed_password", "Current", "Shift", "M", "123456789", "current@example.com", "", "Medic").
			AddRow(2, time.Now(), time.Now(), nil, "next_shift", "hashed_password", "Next", "Shift", "F", "987654321", "next@example.com", "", "Technical")

		// Note: GORM adds extra parentheses around each condition in OR clauses
		mock.ExpectQuery(`SELECT DISTINCT employees\.\* FROM "employees" JOIN employee_shifts ON employees\.id = employee_shifts\.employee_id JOIN shifts ON employee_shifts\.shift_id = shifts\.id WHERE \(\(\(shifts\.shift_date = \$1 AND shifts\.shift_type = \$2\)\) OR \(\(shifts\.shift_date = \$3 AND shifts\.shift_type = \$4\)\)\) AND "employees"\."deleted_at" IS NULL`).
			WillReturnRows(rows)

		testTime := time.Date(2023, 1, 15, 13, 30, 0, 0, time.UTC)       // 1:30 PM, 30 min before shift 1 ends
		employees, err := repo.GetOnCallEmployees(testTime, 1*time.Hour) // 1 hour buffer

		assert.NoError(t, err)
		assert.Len(t, employees, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestShiftRepository_getShiftTypeForTime(t *testing.T) {
	t.Parallel()

	logger := utils.NewTestLogger()
	repo := &shiftRepository{log: logger}

	tests := []struct {
		name     string
		hour     int
		minute   int
		expected int
	}{
		// Shift 1: 6am-2pm
		{"Early morning - 6:00 AM", 6, 0, 1},
		{"Mid morning - 10:30 AM", 10, 30, 1},
		{"Just before shift end - 1:59 PM", 13, 59, 1},

		// Shift 2: 2pm-10pm
		{"Shift start - 2:00 PM", 14, 0, 2},
		{"Evening - 6:30 PM", 18, 30, 2},
		{"Just before shift end - 9:59 PM", 21, 59, 2},

		// Shift 3: 10pm-6am
		{"Night start - 10:00 PM", 22, 0, 3},
		{"Midnight", 0, 0, 3},
		{"Early morning - 3:30 AM", 3, 30, 3},
		{"Just before shift end - 5:59 AM", 5, 59, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTime := time.Date(2023, 1, 1, tt.hour, tt.minute, 0, 0, time.UTC)
			result := repo.getShiftTypeForTime(testTime)
			assert.Equal(t, tt.expected, result, "Expected shift type %d for time %s", tt.expected, testTime.Format("15:04"))
		})
	}

	t.Run("it works correctly when time is exactly at shift boundaries", func(t *testing.T) {
		boundaries := []struct {
			time      time.Time
			shiftType int
		}{
			{time.Date(2023, 1, 1, 6, 0, 0, 0, time.UTC), 1},  // Exactly 6:00 AM
			{time.Date(2023, 1, 1, 14, 0, 0, 0, time.UTC), 2}, // Exactly 2:00 PM
			{time.Date(2023, 1, 1, 22, 0, 0, 0, time.UTC), 3}, // Exactly 10:00 PM
		}

		for _, b := range boundaries {
			result := repo.getShiftTypeForTime(b.time)
			assert.Equal(t, b.shiftType, result, "Time %s should be shift %d", b.time.Format("15:04"), b.shiftType)
		}
	})
}

func TestShiftRepository_getTimeUntilShiftEnd(t *testing.T) {
	t.Parallel()

	logger := utils.NewTestLogger()
	repo := &shiftRepository{log: logger}

	tests := []struct {
		name          string
		hour          int
		minute        int
		second        int
		shiftType     int
		expectedHours int
		expectedMins  int
	}{
		// Shift 1: 6am-2pm (ends at 14:00)
		{"Shift 1 - Start of shift", 6, 0, 0, 1, 8, 0},
		{"Shift 1 - Mid shift", 10, 30, 0, 1, 3, 30},
		{"Shift 1 - 1 hour before end", 13, 0, 0, 1, 1, 0},
		{"Shift 1 - 30 minutes before end", 13, 30, 0, 1, 0, 30},
		{"Shift 1 - 1 minute before end", 13, 59, 0, 1, 0, 1},

		// Shift 2: 2pm-10pm (ends at 22:00)
		{"Shift 2 - Start of shift", 14, 0, 0, 2, 8, 0},
		{"Shift 2 - Mid shift", 18, 30, 0, 2, 3, 30},
		{"Shift 2 - 1 hour before end", 21, 0, 0, 2, 1, 0},
		{"Shift 2 - 30 minutes before end", 21, 30, 0, 2, 0, 30},

		// Shift 3: 10pm-6am (ends at 6:00 next day)
		{"Shift 3 - Start of shift (22:00)", 22, 0, 0, 3, 8, 0},
		{"Shift 3 - Midnight", 0, 0, 0, 3, 6, 0},
		{"Shift 3 - 3 AM", 3, 0, 0, 3, 3, 0},
		{"Shift 3 - 1 hour before end", 5, 0, 0, 3, 1, 0},
		{"Shift 3 - 30 minutes before end", 5, 30, 0, 3, 0, 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTime := time.Date(2023, 1, 1, tt.hour, tt.minute, tt.second, 0, time.UTC)
			result := repo.getTimeUntilShiftEnd(testTime, tt.shiftType)

			expectedDuration := time.Duration(tt.expectedHours)*time.Hour + time.Duration(tt.expectedMins)*time.Minute

			// small differences should be tolerated
			diff := result - expectedDuration
			if diff < 0 {
				diff = -diff
			}

			assert.True(t, diff < time.Minute,
				"Expected ~%v, got %v for time %s (diff: %v)",
				expectedDuration, result, testTime.Format("15:04:05"), diff)
		})
	}

	t.Run("Invalid shift type", func(t *testing.T) {
		testTime := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
		result := repo.getTimeUntilShiftEnd(testTime, 999) // Invalid shift type
		assert.Equal(t, time.Duration(0), result, "Invalid shift type should return 0 duration")
	})
}

func TestShiftRepository_getNextShift(t *testing.T) {
	t.Parallel()

	logger := utils.NewTestLogger()
	repo := &shiftRepository{log: logger}

	baseDate := time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC) // Sunday

	tests := []struct {
		name              string
		currentShiftType  int
		expectedShiftType int
		expectedDateDiff  int // days difference from base date
	}{
		{"From shift 1 to shift 2", 1, 2, 0}, // Same day
		{"From shift 2 to shift 3", 2, 3, 0}, // Same day
		{"From shift 3 to shift 1", 3, 1, 1}, // Next day
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextShiftType, nextDate := repo.getNextShift(tt.currentShiftType, baseDate)

			assert.Equal(t, tt.expectedShiftType, nextShiftType, "Expected next shift type %d", tt.expectedShiftType)

			expectedDate := baseDate.Add(time.Duration(tt.expectedDateDiff) * 24 * time.Hour)
			assert.Equal(t, expectedDate.Truncate(24*time.Hour), nextDate.Truncate(24*time.Hour),
				"Expected date %s, got %s", expectedDate.Format("2006-01-02"), nextDate.Format("2006-01-02"))
		})
	}
}

func TestShiftRepository_shiftBufferLogic(t *testing.T) {
	t.Parallel()

	logger := utils.NewTestLogger()
	repo := &shiftRepository{log: logger}

	tests := []struct {
		name              string
		currentHour       int
		currentMinute     int
		shiftBuffer       time.Duration
		shouldIncludeNext bool
		description       string
	}{
		// Shift 1 scenarios (ends at 14:00)
		{"Shift 1 - No buffer", 10, 0, 0, false, "Mid-shift with no buffer"},
		{"Shift 1 - Buffer but not close to end", 12, 0, 1 * time.Hour, false, "2 hours before end with 1h buffer"},
		{"Shift 1 - Within buffer window", 13, 30, 1 * time.Hour, true, "30 min before end with 1h buffer"},
		{"Shift 1 - Exactly at buffer threshold", 13, 0, 1 * time.Hour, true, "Exactly 1h before end with 1h buffer"},

		// Shift 2 scenarios (ends at 22:00)
		{"Shift 2 - Within buffer window", 21, 30, 1 * time.Hour, true, "30 min before end with 1h buffer"},
		{"Shift 2 - Outside buffer window", 20, 0, 1 * time.Hour, false, "2 hours before end with 1h buffer"},

		// Shift 3 scenarios (ends at 6:00 next day)
		{"Shift 3 - Late night within buffer", 23, 30, 1 * time.Hour, false, "23:30 with 1h buffer (6.5h remaining)"},
		{"Shift 3 - Early morning within buffer", 5, 30, 1 * time.Hour, true, "5:30 AM with 1h buffer (30min remaining)"},
		{"Shift 3 - Midnight outside buffer", 0, 0, 1 * time.Hour, false, "Midnight with 1h buffer (6h remaining)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTime := time.Date(2023, 1, 15, tt.currentHour, tt.currentMinute, 0, 0, time.UTC)
			currentShiftType := repo.getShiftTypeForTime(testTime)
			timeUntilEnd := repo.getTimeUntilShiftEnd(testTime, currentShiftType)

			shouldInclude := tt.shiftBuffer > 0 && timeUntilEnd <= tt.shiftBuffer

			assert.Equal(t, tt.shouldIncludeNext, shouldInclude,
				"%s: Time until end: %v, Buffer: %v, Should include next: %v",
				tt.description, timeUntilEnd, tt.shiftBuffer, tt.shouldIncludeNext)
		})
	}
}
