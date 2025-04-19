package repositories

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/pd120424d/mountain-service/api/shared/utils"

	"github.com/pd120424d/mountain-service/api/employee/internal/model"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gorm.io/gorm/logger"

	_ "modernc.org/sqlite"
)

func setupSQLiteTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:testdb_%d?mode=memory&cache=shared", time.Now().UnixNano())

	dialector := sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        dsn,
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // Use stdout for logs
			logger.Config{
				LogLevel: logger.Info, // Show SQL + params
			},
		),
	})
	require.NoError(t, err, "failed to open sqlite in-memory db")

	// require.NoError(t, db.Migrator().DropTable(&model.EmployeeShift{}, &model.Shift{}, &model.Employee{}))
	require.NoError(t, db.AutoMigrate(&model.Employee{}, &model.Shift{}, &model.EmployeeShift{}))

	return db
}

func TestEmployeeRepository_Create(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB := setupSQLiteTestDB(t)
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

		log.Infof("type from test: %T\n", employee.ProfileType)
		log.Infof("type in SQL: %T\n", employee.ProfileType.String())

		err := repo.Create(employee)
		assert.NoError(t, err)
	})

	t.Run("it returns an error when password is too long and hashing failes", func(t *testing.T) {
		employee.Password = "verylongpasswordthatexceedstheallowedlength"
		err := repo.Create(employee)
		assert.Error(t, err)
	})
}

func TestEmployeeRepository_GetAll(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB := setupSQLiteTestDB(t)
	repo := NewEmployeeRepository(log, gormDB)

	t.Run("it returns all employees when they exist", func(t *testing.T) {
		gormDB.Create(&model.Employee{Username: "test-user", FirstName: "Bruce", Email: "test-user@example.com", LastName: "Lee", Password: "Pass123!"})
		gormDB.Create(&model.Employee{Username: "asmith", FirstName: "Alice", Email: "asmith@example.com", LastName: "Smith", Password: "Pass123!"})

		employees, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Len(t, employees, 2)
	})
}

func TestEmployeeRepository_GetEmployeeByID(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB := setupSQLiteTestDB(t)
	repo := NewEmployeeRepository(log, gormDB)

	t.Run("it returns an employee when it exists", func(t *testing.T) {
		gormDB.Create(&model.Employee{Username: "test-user", FirstName: "Bruce", LastName: "Lee", Password: "Pass123!"})

		var employee model.Employee
		err := repo.GetEmployeeByID("1", &employee)

		assert.NoError(t, err)
		assert.Equal(t, "test-user", employee.Username)
	})
}

func TestEmployeeRepository_GetEmployeeByUsername(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB := setupSQLiteTestDB(t)
	repo := NewEmployeeRepository(log, gormDB)

	t.Run("it returns an employee when it exists", func(t *testing.T) {
		gormDB.Create(&model.Employee{Username: "test-user", FirstName: "Bruce", LastName: "Lee", Password: "Pass123!"})

		employee, err := repo.GetEmployeeByUsername("test-user")

		assert.NoError(t, err)
		assert.Equal(t, "test-user", employee.Username)
	})
}

func TestEmployeeRepository_ListEmployees(t *testing.T) {
	log := utils.NewTestLogger()
	db := setupSQLiteTestDB(t)
	repo := NewEmployeeRepository(log, db)

	// Seed sample data
	tx := db.Create(&model.Employee{Username: "test-user", FirstName: "Bruce", LastName: "Lee", Email: "test-user@example.com", ProfileType: model.Medic})
	require.NoError(t, tx.Error)
	tx = db.Create(&model.Employee{Username: "jackiec", FirstName: "Jackie", LastName: "Chan", Email: "jackiec@example.com", ProfileType: model.Technical})
	require.NoError(t, tx.Error)
	tx = db.Create(&model.Employee{Username: "pd120424d", FirstName: "System", LastName: "Admin", Email: "admin@example.com", ProfileType: model.Administrator})
	require.NoError(t, tx.Error)

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
				{ID: 1, Username: "test-user", FirstName: "Bruce", LastName: "Lee", Email: "test-user@example.com", ProfileType: model.Medic},
			},
		},
		{
			name: "it returns employees filtered by username when they exist",
			filters: map[string]any{
				"username": "test-user",
			},
			expected: []model.Employee{
				{ID: 1, Username: "test-user", FirstName: "Bruce", LastName: "Lee", Email: "test-user@example.com", ProfileType: model.Medic},
			},
		},
		{
			name: "it returns employees filtered by both email and username when they exist",
			filters: map[string]any{
				"username": "test-user",
				"email":    "test-user@example.com",
			},
			expected: []model.Employee{
				{ID: 1, Username: "test-user", FirstName: "Bruce", LastName: "Lee", Email: "test-user@example.com", ProfileType: model.Medic},
			},
		},
		{
			name:    "it returns all employees when no filters are provided",
			filters: map[string]any{},
			expected: []model.Employee{
				{ID: 1, Username: "test-user", FirstName: "Bruce", LastName: "Lee", Email: "test-user@example.com", ProfileType: model.Medic},
				{ID: 2, Username: "jackiec", FirstName: "Jackie", LastName: "Chan", Email: "jackiec@example.com", ProfileType: model.Technical},
				{ID: 3, Username: "pd120424d", FirstName: "System", LastName: "Admin", Email: "admin@example.com", ProfileType: model.Administrator},
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

			employees, err := repo.ListEmployees(test.filters)

			assert.NoError(t, err)

			log.Infof("employees: %v", employees)

			assert.True(t, cmp.Equal(employees, test.expected, cmpopts.IgnoreFields(model.Employee{}, "CreatedAt", "UpdatedAt", "DeletedAt")))
		})
	}
}

func TestEmployeeRepository_UpdateEmployee(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB := setupSQLiteTestDB(t)
	repo := NewEmployeeRepository(log, gormDB)

	t.Run("it updates an employee when it exists", func(t *testing.T) {
		gormDB.Create(&model.Employee{Username: "test-user", FirstName: "Bruce", LastName: "Lee", Password: "Pass123!"})

		employee := &model.Employee{
			ID:        1,
			FirstName: "Bruce",
			LastName:  "Lee",
			Email:     "test-user@example.com",
			Username:  "test-user",
			Gender:    "M",
			Phone:     "123456789",
		}

		err := repo.UpdateEmployee(employee)
		assert.NoError(t, err)
	})
}

func TestEmployeeRepository_Delete(t *testing.T) {
	log := utils.NewTestLogger()

	gormDB := setupSQLiteTestDB(t)
	repo := NewEmployeeRepository(log, gormDB)

	t.Run("it deletes an employee when it exists", func(t *testing.T) {
		gormDB.Create(&model.Employee{Username: "test-user", FirstName: "Bruce", LastName: "Lee", Password: "Pass123!"})

		err := repo.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("it returns an error when employee Lees not exist", func(t *testing.T) {

		err := repo.Delete(999)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}
