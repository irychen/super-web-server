package logger

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type GormLoggerConfig struct {
	LogLevel                  gormlogger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

type GormLogger struct {
	ZapLogger *Logger
	config    GormLoggerConfig
}

func NewGormLogger(zapLogger *Logger, config GormLoggerConfig) *GormLogger {
	return &GormLogger{
		ZapLogger: zapLogger,
		config:    config,
	}
}

func (l GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := l
	newlogger.config.LogLevel = level
	return &newlogger
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= gormlogger.Info {
		l.ZapLogger.Info(fmt.Sprintf(msg, data...))
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= gormlogger.Warn {
		l.ZapLogger.Warn(fmt.Sprintf(msg, data...))
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= gormlogger.Error {
		l.ZapLogger.Error(fmt.Sprintf(msg, data...))
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.config.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
	}

	if !l.config.SkipCallerLookup {
		fields = append(fields, zap.String("file", fileWithLineNum()))
	}

	switch {
	case err != nil && l.config.LogLevel >= gormlogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.config.IgnoreRecordNotFoundError):
		l.ZapLogger.Error("SQL Error", append(fields, zap.Error(err))...)
	case elapsed > l.config.SlowThreshold && l.config.SlowThreshold != 0 && l.config.LogLevel >= gormlogger.Warn:
		l.ZapLogger.Warn("Slow SQL", append(fields, zap.Duration("threshold", l.config.SlowThreshold))...)
	case l.config.LogLevel == gormlogger.Info:
		l.ZapLogger.Info("SQL", fields...)
	}
}

func fileWithLineNum() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && !strings.Contains(file, "gorm.io") {
			return fmt.Sprintf("%s:%d", filepath.Base(file), line)
		}
	}
	return ""
}

func ParseStringGormLogLevel(level string) (gormlogger.LogLevel, error) {
	switch level {
	case "silent":
		return gormlogger.Silent, nil
	case "info":
		return gormlogger.Info, nil
	case "warn":
		return gormlogger.Warn, nil
	case "error":
		return gormlogger.Error, nil
	default:
		return gormlogger.Info, fmt.Errorf("invalid log level: %s", level)
	}
}
