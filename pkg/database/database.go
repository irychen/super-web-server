package database

import (
	"fmt"
	"net/url"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host            string
	Port            int
	Username        string
	Password        string
	DatabaseName    string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	Timezone        string       // timezone configuration
	Charset         string       // character set (primarily for MySQL)
	ParseTime       bool         // parse time (for MySQL)
	GormConfig      *gorm.Config // additional GORM configuration options
}

func GetMySQLDNS(config Config) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.Username, config.Password, config.Host, config.Port, config.DatabaseName)
	parseTime := "True"
	if !config.ParseTime {
		parseTime = "False"
	}
	encodedTimezone := url.QueryEscape(config.Timezone)
	query := fmt.Sprintf("charset=%s&parseTime=%s&loc=%s", config.Charset, parseTime, encodedTimezone)
	return fmt.Sprintf("%s?%s", dsn, query)
}

type DB struct {
	*gorm.DB
}

func NewDB(config Config, logger logger.Interface) (*DB, error) {
	var dialector gorm.Dialector
	var dns = GetMySQLDNS(config)

	dialector = mysql.New(mysql.Config{
		DSN:               dns,
		DefaultStringSize: 256,
	})

	GConfig := config.GormConfig

	if GConfig == nil {
		GConfig = &gorm.Config{}
	}

	GConfig.Logger = logger

	db, err := gorm.Open(dialector, GConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w with dns: %s", err, dns)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql db: %w", err)
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	return &DB{db}, nil
}
