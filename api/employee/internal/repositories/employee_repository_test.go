package repositories

import (
	"testing"
	"time"

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
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestEmployeeRepository_Create(t *testing.T) {
	log := utils.NewNamedLogger("testLogger")

	gormDB, mock := setupMockDB(t)
	repo := NewEmployeeRepository(log, gormDB)

	employee := &model.Employee{
		Username:  "jdoe",
		Password:  "Password123!",
		FirstName: "John",
		LastName:  "Doe",
	}

	t.Run("it does create an employee when data is valid", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO \"employees\"").
			WithArgs(employee.Username, employee.Password, employee.FirstName, employee.LastName).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(employee)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestEmployeeRepository_Delete(t *testing.T) {
	log := utils.NewNamedLogger("testLogger")

	gormDB, mock := setupMockDB(t)
	repo := NewEmployeeRepository(log, gormDB)

	employeeID := uint(1)

	t.Run("it returns an error when the employee is already deleted", func(t *testing.T) {
		mock.ExpectQuery("SELECT \"deleted_at\" FROM \"employees\" WHERE id = ?").
			WithArgs(employeeID).
			WillReturnRows(sqlmock.NewRows([]string{"deleted_at"}).AddRow(time.Now()))

		err := repo.Delete(employeeID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	//t.Run("it does soft delete the employee when the employee is not deleted", func(t *testing.T) {
	//	mock.ExpectQuery("SELECT \"deleted_at\" FROM \"employees\" WHERE id = ?").
	//		WithArgs(employeeID).
	//		WillReturnRows(sqlmock.NewRows([]string{"deleted_at"}).AddRow(nil))
	//
	//	mock.ExpectBegin()
	//	mock.ExpectExec(`UPDATE \"employees\" SET \"deleted_at\"`).
	//		WithArgs(sqlmock.AnyArg(), employeeID).
	//		WillReturnResult(sqlmock.NewResult(1, 1))
	//	mock.ExpectCommit()
	//
	//	err := repo.Delete(employeeID)
	//	assert.NoError(t, err)
	//	assert.NoError(t, mock.ExpectationsWereMet())
	//})
}
