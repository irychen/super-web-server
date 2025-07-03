package app

import (
	"fmt"
	"super-web-server/internal/model"
	"super-web-server/internal/seed"
	"super-web-server/pkg/database"
	"super-web-server/pkg/logger"

	"go.uber.org/zap"
)

func (a *App) InitDatabase() error {
	var dbConfig = a.config.DB

	gormLogLevel, err := logger.ParseStringGormLogLevel(dbConfig.LogLevel)
	if err != nil {
		return fmt.Errorf("parse db log level failed %w", err)
	}

	gormLogger := logger.NewGormLogger(logger.GetModuleLogger("gorm"), logger.GormLoggerConfig{
		LogLevel:                  gormLogLevel,
		SlowThreshold:             dbConfig.SlowThreshold,
		SkipCallerLookup:          true,
		IgnoreRecordNotFoundError: true,
	})

	config := database.Config{
		Host:            dbConfig.Host,
		Port:            dbConfig.Port,
		Username:        dbConfig.Username,
		Password:        dbConfig.Password,
		DatabaseName:    dbConfig.Database,
		MaxIdleConns:    dbConfig.MaxIdleConns,
		MaxOpenConns:    dbConfig.MaxOpenConns,
		ConnMaxLifetime: dbConfig.ConnMaxLifetime,
		Timezone:        dbConfig.Timezone,
		Charset:         dbConfig.Charset,
		ParseTime:       dbConfig.ParseTime,
	}

	db, err := database.NewDB(config, gormLogger)

	if err != nil {
		return err
	}

	a.db = db

	logger.Info("database initialized successfully")

	err = db.AutoMigrate(
		&model.User{},
		&model.UserRole{},
	)

	if err != nil {
		logger.Error("database migrate failed", zap.Error(err))
		return err
	}

	logger.Info("database migrate successfully")

	if err := seed.Seed(db, a.snowflake); err != nil {
		logger.Error("database seed failed", zap.Error(err))
		return err
	}

	logger.Info("database seed successfully")

	return nil
}
