package config

import (
	"github.com/pd120424d/mountain-service/api/shared/config"
	"github.com/pd120424d/mountain-service/api/shared/utils"

	"gorm.io/gorm"
)

const (
	ServerPort = "8083"
)

func GetUrgencyDB(log utils.Logger, connString string) *gorm.DB {
	return config.GetDbConnection(log, connString)
}
