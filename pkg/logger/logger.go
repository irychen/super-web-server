package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*zap.Logger
}

var (
	globalLogger *Logger
)

type Format string

const (
	FormatJSON    Format = "json"
	FormatConsole Format = "console"
)

type Config struct {
	Level      zapcore.Level // logging level (debug, info, warn, error, fatal)
	Format     Format        // log format (json or console)
	Directory  string        // directory to store log files
	Filename   string        // filename for log file
	MaxSize    int           // max size in MB
	MaxBackups int           // max number of backups
	MaxAge     int           // max age in days
	Compress   bool          // compress old log files
	Stdout     bool          // output to stdout
}

func InitLogger(config Config) error {

	if err := os.MkdirAll(config.Directory, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	fmt.Println("config.Level", config.Level)

	logFile := filepath.Join(config.Directory, config.Filename)

	fileWriter := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder

	if config.Format == FormatJSON {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	var cores []zapcore.Core

	cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(fileWriter), config.Level))

	if config.Stdout {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), config.Level))
	}

	core := zapcore.NewTee(cores...)

	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	globalLogger = &Logger{Logger: logger}

	return nil
}

func GetModuleLogger(module string) *Logger {
	return &Logger{Logger: globalLogger.Named(module)}
}

func ParseStringLogLevel(level string) (zapcore.Level, error) {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	case "fatal":
		zapLevel = zapcore.FatalLevel
	case "panic":
		zapLevel = zapcore.PanicLevel
	default:
		return zapcore.InfoLevel, fmt.Errorf("invalid log level: %s", level)
	}
	return zapLevel, nil
}

func ParseStringFormat(format string) (Format, error) {
	var zapFormat Format
	switch format {
	case "json":
		zapFormat = FormatJSON
	case "console":
		zapFormat = FormatConsole
	default:
		return FormatConsole, fmt.Errorf("invalid log format: %s", format)
	}
	return zapFormat, nil
}

func Sync() error {
	return globalLogger.Sync()
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

func DebugF(msg string, args ...interface{}) {
	globalLogger.Debug(fmt.Sprintf(msg, args...))
}

func DebugE(msg string, err error, fields ...zap.Field) {
	globalLogger.Debug(msg, append([]zap.Field{zap.Error(err)}, fields...)...)
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

func InfoF(msg string, args ...interface{}) {
	globalLogger.Info(fmt.Sprintf(msg, args...))
}

func InfoE(msg string, err error, fields ...zap.Field) {
	globalLogger.Info(msg, append([]zap.Field{zap.Error(err)}, fields...)...)
}

func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

func WarnF(msg string, args ...interface{}) {
	globalLogger.Warn(fmt.Sprintf(msg, args...))
}

func WarnE(msg string, err error, fields ...zap.Field) {
	globalLogger.Warn(msg, append([]zap.Field{zap.Error(err)}, fields...)...)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

func ErrorF(msg string, args ...interface{}) {
	globalLogger.Error(fmt.Sprintf(msg, args...))
}

func ErrorE(msg string, err error, fields ...zap.Field) {
	globalLogger.Error(msg, append([]zap.Field{zap.Error(err)}, fields...)...)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

func FatalF(msg string, args ...interface{}) {
	globalLogger.Fatal(fmt.Sprintf(msg, args...))
}

func FatalE(msg string, err error, fields ...zap.Field) {
	globalLogger.Fatal(msg, append([]zap.Field{zap.Error(err)}, fields...)...)
}

func Panic(msg string, fields ...zap.Field) {
	globalLogger.Panic(msg, fields...)
}

func PanicF(msg string, args ...interface{}) {
	globalLogger.Panic(fmt.Sprintf(msg, args...))
}

func PanicE(msg string, err error, fields ...zap.Field) {
	globalLogger.Panic(msg, append([]zap.Field{zap.Error(err)}, fields...)...)
}
