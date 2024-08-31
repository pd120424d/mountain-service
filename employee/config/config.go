package config

import (
	"gorm.io/gorm"

	"mountain-service/shared/config"
	"mountain-service/shared/utils"
)

const (
	ServerPort = "8082"
)

func GetEmployeeDB(log utils.Logger, connString string) *gorm.DB {
	return config.GetDbConnection(log, connString)
}
