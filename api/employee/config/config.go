package config

import (
	"api/shared/config"
	"api/shared/utils"
	"fmt"

	"gorm.io/gorm"
)

const (
	ServerPort = "8082"
)

func GetEmployeeDB(log utils.Logger, connString string) *gorm.DB {
	return config.GetDbConnection(log, connString)
}
