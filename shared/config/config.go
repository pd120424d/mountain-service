package config

import (
	"fmt"

	"mountain-service/shared/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	dbConnectionTemplate = "host=%v user=%v password=%v dbname=%v port=%v sslmode=disable"
	defaultHost          = "localhost"
	defaultPort          = "5432"
	defaultUser          = "postgres"
	defaultPassword      = "etf"
	defaultDbName        = "activity"
)

// GetDbConnection creates a connection to postgres db with provided host, port and name.
// If any of these is not provided (empty) it fall back to the default values.
//
//	default host = `localhost`
//	default port = `5432`
//	default db name = `activity`
func GetDbConnection(log utils.Logger, hostname, port, user, dbName string) *gorm.DB {
	if hostname == "" {
		log.Warnf("Hostname not provided, creating a connection with a default one (%v)", defaultHost)
		hostname = defaultHost
	}
	if port == "" {
		log.Warnf("Port not provided, creating a connection with a default one (%v)", defaultPort)
		port = defaultPort
	}
	if dbName == "" {
		log.Warnf("Database name not provided, creating a connection with a default one (%v)", defaultDbName)
		dbName = defaultDbName
	}
	if user == "" {
		log.Warnf("User not provided, creating a connection with a default one (%v)", defaultUser)
		user = defaultUser
	}

	dsn := fmt.Sprintf(dbConnectionTemplate, hostname, user, defaultPassword, dbName, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to '%v' database: %v", dbName, err)
	}
	return db
}
