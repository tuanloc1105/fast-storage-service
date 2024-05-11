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
	usersOtpDataMigrateErr := db.AutoMigrate(&model.UsersOtpData{})
	if usersOtpDataMigrateErr != nil {
		panic(usersOtpDataMigrateErr)
	}
	userAuthenticationLogMigrateErr := db.AutoMigrate(&model.UserAuthenticationLog{})
	if userAuthenticationLogMigrateErr != nil {
		panic(userAuthenticationLogMigrateErr)
	}
}

func InsertData(db *gorm.DB) {
}
