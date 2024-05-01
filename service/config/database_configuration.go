package config

import (
	"context"
	"errors"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/log"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabaseConnection() (db *gorm.DB, err error) {
	databaseUsername, isDatabaseUsernameSet := os.LookupEnv("DATABASE_USERNAME")
	databasePassword, isDatabasePasswordSet := os.LookupEnv("DATABASE_PASSWORD")
	databaseHost, isDatabaseHostSet := os.LookupEnv("DATABASE_HOST")
	databasePort, isDatabasePortSet := os.LookupEnv("DATABASE_PORT")
	databaseName, isDatabaseNameSet := os.LookupEnv("DATABASE_NAME")

	if !isDatabaseUsernameSet || !isDatabasePasswordSet || !isDatabaseHostSet || !isDatabasePortSet || !isDatabaseNameSet {
		return nil, errors.New("database info was not full set in enviroment")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		databaseHost,
		databaseUsername,
		databasePassword,
		databaseName,
		databasePort,
	)
	var ctx = context.Background()
	log.WithLevel(
		constant.Info,
		ctx,
		"Database connect info:\n    - databaseHost: %s\n    - databaseUsername: %s\n    - databaseName: %s\n    - databasePort: %s",
		databaseHost,
		databaseUsername,
		databaseName,
		databasePort,
	)
	myLog := &MyCustomDatabaseLogger{}
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: myLog.LogMode(logger.Info),
	})
	if err != nil {
		log.WithLevel(
			constant.Error,
			ctx,
			"An error has been occurred when trying to connect to Database:\n\t- error: %v",
			err,
		)
		return db, err
	}

	sqlDB, getDbObjectError := db.DB()

	if getDbObjectError != nil {
		log.WithLevel(
			constant.Error,
			ctx,
			"An error has been occurred when trying to get Database Object:\n\t- error: %v",
			getDbObjectError,
		)
		return db, getDbObjectError
	}

	// sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	if isDatabaseMigration, isDatabaseMigrationExist := os.LookupEnv("DATABASE_MIGRATION"); isDatabaseMigration == "true" && isDatabaseMigrationExist {
		MigrationAndInsertDate(db)
	}
	if isDatabaseInitializationData, isDatabaseInitializationDataExist := os.LookupEnv("DATABASE_INITIALIZATION_DATA"); isDatabaseInitializationData == "true" && isDatabaseInitializationDataExist {
		InsertData(db)
	}

	return db, err
}
