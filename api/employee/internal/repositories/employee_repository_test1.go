package repositories

// import (
// 	"testing"

// 	"github.com/pd120424d/mountain-service/api/shared/utils"

// 	"github.com/pd120424d/mountain-service/api/employee/internal/model"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/stretchr/testify/assert"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// 	"gorm.io/gorm/logger"
// )

// func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)

// 	gormDB, err := gorm.Open(postgres.New(postgres.Config{
// 		Conn: db,
// 	}), &gorm.Config{
// 		Logger: logger.Default.LogMode(logger.Error),
// 	})
// 	assert.NoError(t, err)

// 	return gormDB, mock
// }

// func TestEmployeeRepository_Create(t *testing.T) {
// 	log := utils.NewTestLogger()

// 	gormDB, mock := setupMockDB(t)
// 	repo := NewEmployeeRepository(log, gormDB)

// 	employee := &model.Employee{
// 		Username:       "test-user",
// 		Password:       "Pass123!",
// 		FirstName:      "Bruce",
// 		LastName:       "Lee",
// 		Gender:         "M",
// 		Phone:          "123456789",
// 		Email:          "test-user@example.com",
// 		ProfilePicture: "https://example.com/profile.jpg",
// 		ProfileType:    model.Medic,
// 	}
// }

// func TestEmployeeRepository_UpdateEmployee(t *testing.T) {
// 	log := utils.NewTestLogger()

// 	gormDB, mock := setupMockDB(t)
// 	repo := NewEmployeeRepository(log, gormDB)

// }

// func TestEmployeeRepository_Delete(t *testing.T) {
// 	log := utils.NewTestLogger()

// 	gormDB, mock := setupMockDB(t)
// 	repo := NewEmployeeRepository(log, gormDB)

// 	t.Run("it returns an error when employee Lees not exist", func(t *testing.T) {
// 		mock.ExpectQuery(`SELECT \* FROM "employees"`).
// 			WithArgs(999, 1).
// 			WillReturnError(gorm.ErrRecordNotFound)

// 		err := repo.Delete(999)
// 		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
// 		assert.NoError(t, mock.ExpectationsWereMet())
// 	})

// 	t.Run("it returns an error when it fails to delete an employee", func(t *testing.T) {
// 		mock.ExpectQuery(`SELECT \* FROM "employees"`).
// 			WithArgs(1, 1).
// 			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "first_name", "last_name"}).
// 				AddRow(1, "test-user", "Bruce", "Lee"))

// 		mock.ExpectBegin()
// 		mock.ExpectExec(`DELETE FROM "employees"`).
// 			WithArgs(1).
// 			WillReturnError(sqlmock.ErrCancelled)
// 		mock.ExpectRollback()

// 		err := repo.Delete(1)
// 		assert.Error(t, err)
// 		assert.NoError(t, mock.ExpectationsWereMet())
// 	})
// }
