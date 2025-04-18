package repositories

import (
	"database/sql/driver"
	"fmt"
	"slices"
	"testing"

	"github.com/pd120424d/mountain-service/api/shared/utils"

	"github.com/pd120424d/mountain-service/api/employee/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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

func TestEmployeeRepository_Create(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewEmployeeRepository(log, gormDB)

	employee := &model.Employee{
		Username:       "test-user",
		Password:       "Pass123!",
		FirstName:      "Bruce",
		LastName:       "Lee",
		Gender:         "M",
		Phone:          "123456789",
		Email:          "test-user@example.com",
		ProfilePicture: "https://example.com/profile.jpg",
		ProfileType:    model.Medic,
	}

	t.Run("it creates an employee when data is valid", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "employees"`).
			WithArgs(
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
				sqlmock.AnyArg(), // deleted_at (nullable)
				employee.Username,
				sqlmock.AnyArg(), // password (hashed inside Create)
				employee.FirstName,
				employee.LastName,
				employee.Gender,
				employee.Phone,
				employee.Email,
				employee.ProfilePicture,
				employee.ProfileType.String(), // important: use .String()
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		log.Infof("type from test: %T\n", employee.ProfileType)
		log.Infof("type in SQL: %T\n", employee.ProfileType.String())

		err := repo.Create(employee)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("it returns an error when password is too long and hashing failes", func(t *testing.T) {
		employee.Password = "verylongpasswordthatexceedstheallowedlength"
		err := repo.Create(employee)
		assert.Error(t, err)
	})
}

func TestEmployeeRepository_GetAll(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewEmployeeRepository(log, gormDB)

	t.Run("it returns all employees when they exist", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "employees" WHERE deleted_at IS NULL`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "first_name", "last_name"}).
				AddRow(1, "test-user", "Bruce", "Lee").
				AddRow(2, "asmith", "Alice", "Smith"))

		employees, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Len(t, employees, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestEmployeeRepository_GetEmployeeByID(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewEmployeeRepository(log, gormDB)

	t.Run("it returns an employee when it exists", func(t *testing.T) {
		mock.ExpectQuery(`FROM "employees"`).
			WithArgs("1", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "first_name", "last_name"}).
				AddRow(1, "test-user", "Bruce", "Lee"))

		var employee model.Employee
		err := repo.GetEmployeeByID("1", &employee)

		assert.NoError(t, err)
		assert.Equal(t, "test-user", employee.Username)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestEmployeeRepository_GetEmployeeByUsername(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewEmployeeRepository(log, gormDB)

	t.Run("it returns an employee when it exists", func(t *testing.T) {
		mock.ExpectQuery(`FROM "employees"`).
			WithArgs("test-user", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "first_name", "last_name"}).
				AddRow(1, "test-user", "Bruce", "Lee"))

		employee, err := repo.GetEmployeeByUsername("test-user")

		assert.NoError(t, err)
		assert.Equal(t, "test-user", employee.Username)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestEmployeeRepository_ListEmployees(t *testing.T) {
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
				slices.Reverse(keys)

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

			employees, err := repo.ListEmployees(test.filters)

			assert.NoError(t, err)

			assert.Equal(t, test.expected, employees)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestEmployeeRepository_UpdateEmployee(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewEmployeeRepository(log, gormDB)

	t.Run("it updates an employee when it exists", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "employees"`).
			WithArgs(
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
				nil,              // deleted_at
				"test-user",      // username
				"",               // password
				"Bruce",
				"Lee",
				"", // gender
				"", // phone
				"test-user@example.com",
				"",
				"", // profile_type
				1,  // id
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		employee := &model.Employee{
			ID:        1,
			FirstName: "Bruce",
			LastName:  "Lee",
			Email:     "test-user@example.com",
			Username:  "test-user",
		}

		err := repo.UpdateEmployee(employee)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestEmployeeRepository_Delete(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewEmployeeRepository(log, gormDB)

	// Enable logging for debugging
	// gormDB.Logger = gormDB.Logger.LogMode(logger.Info)

	t.Run("it deletes an employee when it exists", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "employees"`).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "first_name", "last_name"}).
				AddRow(1, "test-user", "Bruce", "Lee"))

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "employees"`).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("it returns an error when employee Lees not exist", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "employees"`).
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		err := repo.Delete(999)
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

		err := repo.Delete(1)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
