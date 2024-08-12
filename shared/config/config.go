package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"mountain-service/shared/utils"
)

const (
	dbConnectionTemplate = "host=%v user=%v password=%v dbname=%v port=%v sslmode=disable"
	defaultHost          = "localhost"
	defaultPort          = "5432"
	defaultUser          = "dpavlovic"
	defaultPassword      = "etf"
	defaultDbName        = "activity"
)

// GetDbConnection creates a connection to postgres db with provided host, port and name.
// If any of these is not provided (empty) it fall back to the default values.
//
//	default host = `defaultHost`
//	default port = `5432`
//	default db name = `activity`
func GetDbConnection(log utils.Logger, hostname, port, dbName string) *gorm.DB {
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
		dbName = defaultHost
	}

	dsn := fmt.Sprintf(dbConnectionTemplate, defaultHost, defaultUser, defaultPassword, dbName, defaultPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to employee database: %v", err)
	}
	return db
}
