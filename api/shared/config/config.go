package config

import (
	"api/shared/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GetDbConnection creates a connection to postgres db with provided host, port and name.
// If any of these is not provided (empty) it fall back to the default values.
//
//	default host = `localhost`
//	default port = `5432`
//	default db name = `activity`
func GetDbConnection(log utils.Logger, connString string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to '%v': %v", connString, err)
	}
	return db
}
