package config

import (
	"context"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/log"
	"time"

	"gorm.io/gorm/logger"
)

type MyCustomDatabaseLogger struct{}

func (d *MyCustomDatabaseLogger) LogMode(level logger.LogLevel) logger.Interface {
	return d
}
func (d *MyCustomDatabaseLogger) Info(ctx context.Context, s string, i ...interface{}) {
	log.WithLevel(
		constant.Info,
		ctx,
		s,
	)
}
func (d *MyCustomDatabaseLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	log.WithLevel(
		constant.Warn,
		ctx,
		s,
	)
}
func (d *MyCustomDatabaseLogger) Error(ctx context.Context, s string, i ...interface{}) {
	log.WithLevel(
		constant.Error,
		ctx,
		s,
	)
}
func (d *MyCustomDatabaseLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rowsAffected := fc()
	errorMessage := "no error"
	if err != nil {
		errorMessage = err.Error()
		log.WithLevel(
			constant.Error,
			ctx,
			"info:\n    - sql: %s\n    - rowsAffected: %v\n    - begin: %s\n    - error: %s",
			sql,
			rowsAffected,
			begin.Format(constant.YyyyMmDdHhMmSsFormat),
			errorMessage,
		)
		return
	}
	log.WithLevel(
		constant.Info,
		ctx,
		"info:\n    - sql: %s\n    - rowsAffected: %v\n    - begin: %s\n    - error: %s",
		sql,
		rowsAffected,
		begin.Format(constant.YyyyMmDdHhMmSsFormat),
		errorMessage,
	)
}
