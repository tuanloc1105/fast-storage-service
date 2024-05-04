package config

import (
	"fast-storage-go-service/model"

	"gorm.io/gorm"
)

func MigrationAndInsertDate(db *gorm.DB) {
	userAccountActivationLogMigrateErr := db.AutoMigrate(&model.UserAccountActivationLog{})
	if userAccountActivationLogMigrateErr != nil {
		panic(userAccountActivationLogMigrateErr)
	}
}

func InsertData(db *gorm.DB) {
}
