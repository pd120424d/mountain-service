package repositories

import (
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
		Username:       "jdoe",
		Password:       "Pass123!",
		FirstName:      "John",
		LastName:       "Doe",
		Gender:         "M",
		Phone:          "123456789",
		Email:          "jdoe@example.com",
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
				AddRow(1, "jdoe", "John", "Doe").
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
				AddRow(1, "jdoe", "John", "Doe"))

		var employee model.Employee
		err := repo.GetEmployeeByID("1", &employee)

		assert.NoError(t, err)
		assert.Equal(t, "jdoe", employee.Username)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestEmployeeRepository_GetEmployeeByUsername(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewEmployeeRepository(log, gormDB)

	t.Run("it returns an employee when it exists", func(t *testing.T) {
		mock.ExpectQuery(`FROM "employees"`).
			WithArgs("jdoe", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "first_name", "last_name"}).
				AddRow(1, "jdoe", "John", "Doe"))

		employee, err := repo.GetEmployeeByUsername("jdoe")

		assert.NoError(t, err)
		assert.Equal(t, "jdoe", employee.Username)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestEmployeeRepository_UpdateEmployee(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB, mock := setupMockDB(t)
	repo := NewEmployeeRepository(log, gormDB)

	t.Run("it updates an employee when it exists", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "employees"`).
			WithArgs(
				sqlmock.AnyArg(), // updated_at
				"John",
				"Doe",
				"jdoe@example.com",
				1,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		employee := &model.Employee{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "jdoe@example.com",
		}

		err := repo.UpdateEmployee(employee)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
