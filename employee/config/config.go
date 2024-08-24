package config

import (
	"fmt"

	"gorm.io/gorm"

	"mountain-service/shared/config"
	"mountain-service/shared/utils"
)

const (
	ServerPort = "8082"
	dbPort     = "5432"
	dbName     = "employee_service"
)

func GetEmployeeDB(log utils.Logger, connString string) *gorm.DB {
	return config.GetDbConnection(log, connString)
}

func CreateEmployeeDB(log utils.Logger, connString string) {
	db := config.GetDbConnection(log, connString)

	cmd := fmt.Sprintf("CREATE DATABASE %s", dbName)
	tx := db.Exec(cmd)
	errDbExist := fmt.Sprintf("ERROR: database \"%s\" already exists (SQLSTATE 42P04)", dbName)
	if tx.Error != nil && tx.Error.Error() != errDbExist {
		log.Fatalf("failed to create '%v' database: %s", dbName, tx.Error.Error())
	} else if tx.Error != nil {
		log.Warnf("'%s' database already exists", dbName)
	} else {
		log.Infof("'%v' database created successfully", dbName)
	}
}
