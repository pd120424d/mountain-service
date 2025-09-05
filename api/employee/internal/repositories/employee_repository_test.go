package repositories

import (
	"context"
	"database/sql/driver"
	"fmt"
	"slices"
	"testing"

	"github.com/pd120424d/mountain-service/api/employee/internal/model"
	"github.com/pd120424d/mountain-service/api/shared/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestEmployeeRepositoryMockDB_ListEmployees(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewEmployeeRepository(log, gormDB)
	gormDB.Logger = gormDB.Logger.LogMode(logger.Info)

	tests := []struct {
		name     string
		filters  map[string]any
		expected []model.Employee
	}{
		{
			name: "it returns employees filtered by email when they exist",
			filters: map[string]any{
				"email": "test-user@example.com",
			},
			expected: []model.Employee{
				{ID: 1, Username: "test-user", FirstName: "Bruce", LastName: "Lee", Email: "test-user@example.com"},
			},
		},
		{
			name: "it returns employees filtered by username when they exist",
			filters: map[string]any{
				"username": "test-user",
			},
			expected: []model.Employee{
				{ID: 1, Username: "test-user", FirstName: "Bruce", LastName: "Lee", Email: "test-user@example.com"},
			},
		},
		{
			name: "it returns employees filtered by both email and username when they exist",
			filters: map[string]any{
				"username": "test-user",
				"email":    "test-user@example.com",
			},
			expected: []model.Employee{
				{ID: 1, Username: "test-user", FirstName: "Bruce", LastName: "Lee", Email: "test-user@example.com"},
			},
		},
		{
			name:    "it returns all employees when no filters are provided",
			filters: map[string]any{},
			expected: []model.Employee{
				{ID: 1, Username: "test-user", FirstName: "Bruce", LastName: "Lee", Email: "test-user@example.com"},
			},
		},
		{
			name: "it returns an empty list when no employees match the filter",
			filters: map[string]any{
				"username": "nonexistent-user",
			},
			expected: []model.Employee{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			var rowsToReturn = sqlmock.NewRows([]string{})
			if len(test.expected) > 0 {
				rowsToReturn = sqlmock.NewRows([]string{"id", "username", "first_name", "last_name", "email"}).
					AddRow(1, "test-user", "Bruce", "Lee", "test-user@example.com")
			}

			expectedArgs := []driver.Value{}
			expectedQuery := `SELECT \* FROM "employees" WHERE "employees"\."deleted_at" IS NULL`
			if len(test.filters) > 0 {
				keys := []string{}
				for k := range test.filters {
					keys = append(keys, k)
				}
				slices.Sort(keys)

				queryFilter := ""
				i := 1
				for _, key := range keys {
					queryFilter += fmt.Sprintf(`%s LIKE \$%d AND `, key, i)
					expectedArgs = append(expectedArgs, "%"+test.filters[key].(string)+"%")
					i++
				}

				expectedQuery = fmt.Sprintf(`SELECT \* FROM "employees" WHERE %v"employees"\."deleted_at" IS NULL`, queryFilter)
			}

			mock.ExpectQuery(expectedQuery).
				WithArgs(expectedArgs...).
				WillReturnRows(rowsToReturn)

			employees, err := repo.ListEmployees(context.Background(), test.filters)

			assert.NoError(t, err)

			assert.Equal(t, test.expected, employees)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

	t.Run("it returns an error when filter is invalid", func(t *testing.T) {
		employees, err := repo.ListEmployees(context.Background(), map[string]any{"invalid": "invalid"})

		assert.Error(t, err)
		assert.Equal(t, "invalid filter key: invalid", err.Error())
		assert.Nil(t, employees)
	})

	t.Run("it returns an error when it fails to list employees", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "employees"`).
			WillReturnError(sqlmock.ErrCancelled)

		employees, err := repo.ListEmployees(context.Background(), map[string]any{})

		assert.Error(t, err)
		assert.Nil(t, employees)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestEmployeeRepositoryMockDB_Delete(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewEmployeeRepository(log, gormDB)

	t.Run("it returns an error when employee Lees not exist", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "employees"`).
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		err := repo.Delete(context.Background(), 999)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("it returns an error when it fails to delete an employee", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "employees"`).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "first_name", "last_name"}).
				AddRow(1, "test-user", "Bruce", "Lee"))

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "employees"`).
			WithArgs(1).
			WillReturnError(sqlmock.ErrCancelled)
		mock.ExpectRollback()

		err := repo.Delete(context.Background(), 1)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestEmployeeRepositoryMockDB_ResetAllData(t *testing.T) {
	t.Parallel()

	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewEmployeeRepository(log, gormDB)

	t.Run("it returns an error when it fails to delete employee-shift associations", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "employee_shifts" WHERE 1=1`).
			WillReturnError(sqlmock.ErrCancelled)
		mock.ExpectRollback()

		err := repo.ResetAllData(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "canceling query due to user request")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("it returns an error when it fails to delete shifts", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "employee_shifts" WHERE 1=1`).
			WillReturnResult(sqlmock.NewResult(0, 5))
		mock.ExpectCommit()

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "shifts" WHERE 1=1`).
			WillReturnError(sqlmock.ErrCancelled)
		mock.ExpectRollback()

		err := repo.ResetAllData(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "canceling query due to user request")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("it returns an error when it fails to delete employees", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "employee_shifts" WHERE 1=1`).
			WillReturnResult(sqlmock.NewResult(0, 5))
		mock.ExpectCommit()

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "shifts" WHERE 1=1`).
			WillReturnResult(sqlmock.NewResult(0, 3))
		mock.ExpectCommit()

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "employees" WHERE 1=1`).
			WillReturnError(sqlmock.ErrCancelled)
		mock.ExpectRollback()

		err := repo.ResetAllData(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "canceling query due to user request")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("it successfully resets all data", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "employee_shifts" WHERE 1=1`).
			WillReturnResult(sqlmock.NewResult(0, 5))
		mock.ExpectCommit()

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "shifts" WHERE 1=1`).
			WillReturnResult(sqlmock.NewResult(0, 3))
		mock.ExpectCommit()

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "employees" WHERE 1=1`).
			WillReturnResult(sqlmock.NewResult(0, 2))
		mock.ExpectCommit()

		err := repo.ResetAllData(context.Background())
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	assert.NoError(t, err)

	return gormDB, mock
}
