package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/pd120424d/mountain-service/api/shared/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host   string
	Port   string
	Name   string
	Models []interface{}
}

const (
	EmployeeServiceName = "employee-service"
	UrgencyServiceName  = "urgency-service"
	ActivityServiceName = "activity-service"

	EmployeeDBName = "employee_service_db"
	UrgencyDBName  = "urgency_service_db"
	ActivityDBName = "activity_service_db"

	EmployeeServicePort = "8082"
	UrgencyServicePort  = "8083"
	ActivityServicePort = "8084"
)

const (
	DB_USER     = "DB_USER"
	DB_PASSWORD = "DB_PASSWORD"
)

func InitDb(log utils.Logger, serviceName string, dbConfig DatabaseConfig) *gorm.DB {
	log.Infof("Setting up database for %s...", serviceName)

	svcNameAllCaps := strings.ToUpper(serviceName)
	svcDbUser := fmt.Sprintf("%s_DB_USER", svcNameAllCaps)
	svcDbPassword := fmt.Sprintf("%s_DB_PASSWORD", svcNameAllCaps)
	svcDbUserFile := fmt.Sprintf("%s_DB_USER_FILE", svcNameAllCaps)
	svcDbPasswordFile := fmt.Sprintf("%s_DB_PASSWORD_FILE", svcNameAllCaps)

	dbUser := os.Getenv(DB_USER)
	if dbUser == "" {
		log.Infof("DB_USER is empty, checking %s", svcDbUserFile)
		userFile := os.Getenv(svcDbUserFile)
		if userFile != "" && userFile != " " {
			var err error
			dbUser, err = ReadSecret(userFile)
			if err != nil {
				log.Fatalf("Failed to read %s from file %s: %v", svcDbUser, userFile, err)
			}
		} else {
			log.Fatalf("Neither DB_USER environment variable nor %s is set. DB_USER='%s', %s='%s'", svcDbUser, dbUser, svcDbUserFile, userFile)
		}
	}

	dbPassword := os.Getenv(DB_PASSWORD)
	if dbPassword == "" {
		log.Infof("DB_PASSWORD is empty, checking %s", svcDbPasswordFile)
		passwordFile := os.Getenv(svcDbPasswordFile)
		if passwordFile != "" && passwordFile != " " {
			var err error
			dbPassword, err = ReadSecret(passwordFile)
			if err != nil {
				log.Fatalf("Failed to read %s from file %s: %v", svcDbPassword, passwordFile, err)
			}
		} else {
			log.Fatalf("Neither DB_PASSWORD environment variable nor %s is set. DB_PASSWORD='%s', %s='%s'", svcDbPassword, dbPassword, svcDbPasswordFile, passwordFile)
		}
	}

	log.Infof("Connecting to database at %s:%s as user %s", dbConfig.Host, dbConfig.Port, dbUser)
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbUser, dbPassword, dbConfig.Name)

	// Create the {serviceName} database if it doesn't exist
	db := getDbConnection(log, connectionString)

	// Auto migrate the model
	err := db.AutoMigrate(dbConfig.Models...)
	if err != nil {
		log.Fatalf("failed to migrate %s models: %v", serviceName, err)
	}
	log.Infof("Successfully migrated %s models", serviceName)

	log.Infof("Database setup finished succsesfully for %s", serviceName)
	return db
}

func ReadSecret(filePath string) (string, error) {
	secret, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(secret)), nil
}

func getDbConnection(log utils.Logger, connString string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to '%v': %v", connString, err)
	}
	return db
}
