package config

import (
	"fmt"
	"reflect"
	"strings"
	"super-web-server/internal/types"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Mode   types.ServerMode
	Server ServerConfig `mapstructure:"server"`
	DB     DBConfig     `mapstructure:"database"`
	Redis  RedisConfig  `mapstructure:"redis"`
	JWT    JWTConfig    `mapstructure:"jwt"`
	Log    LogConfig    `mapstructure:"log"`
}

var defaultConfig = &Config{
	Server: ServerConfig{
		Port:           8080,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		SnowflakeNode:  1,
	},
	Log: LogConfig{
		Level:      "debug",
		Path:       "./logs",
		Filename:   "app.log",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
		Stdout:     true,
	},
	DB: DBConfig{
		Host:            "localhost",
		Port:            3306,
		Username:        "root",
		Password:        "123456",
		Database:        "test",
		Charset:         "utf8mb4",
		Timezone:        "UTC",
		ParseTime:       true,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: 10 * time.Second,
		LogLevel:        "info",
		SlowThreshold:   1 * time.Second,
	},
	Redis: RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "123456",
		DB:       0,
	},
	JWT: JWTConfig{
		Secret: "123456",
		Expire: 1 * time.Hour,
		Issuer: "super-web-server",
	},
}

func LoadConfig(filePath string, serverMode types.ServerMode) (*Config, error) {
	var config = &Config{}

	v := viper.New()

	// 设置默认值 - 这里是关键
	setDefaults(v)

	var configFileName = fmt.Sprintf("config.%s", serverMode.String())

	v.AddConfigPath(filePath)
	v.SetConfigName(configFileName)
	v.SetConfigType("yml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config file failed: %v", err)
	}

	// read env config
	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unmarshal config file failed: %w", err)
	}

	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("validate config failed: %w", err)
	}

	config.Mode = serverMode

	return config, nil
}

// setDefaults 使用反射自动设置所有默认配置值
func setDefaults(v *viper.Viper) {
	// 设置顶级配置
	v.SetDefault("mode", defaultConfig.Mode)
	// 设置嵌套配置
	setDefaultsFromStruct(v, "server", defaultConfig.Server)
	setDefaultsFromStruct(v, "database", defaultConfig.DB)
	setDefaultsFromStruct(v, "redis", defaultConfig.Redis)
	setDefaultsFromStruct(v, "jwt", defaultConfig.JWT)
	setDefaultsFromStruct(v, "log", defaultConfig.Log)
}

// setDefaultsFromStruct 使用反射设置结构体的默认值
func setDefaultsFromStruct(v *viper.Viper, prefix string, structValue interface{}) {
	val := reflect.ValueOf(structValue)
	typ := reflect.TypeOf(structValue)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// 获取 mapstructure 标签作为配置键名
		tag := fieldType.Tag.Get("mapstructure")
		if tag == "" {
			// 如果没有 mapstructure 标签，使用字段名的小写形式
			tag = strings.ToLower(fieldType.Name)
		}

		configKey := fmt.Sprintf("%s.%s", prefix, tag)

		// 设置默认值
		if field.IsValid() && field.CanInterface() {
			v.SetDefault(configKey, field.Interface())
		}
	}
}

func validateConfig(config *Config) error {
	var validate = validator.New()

	if err := validate.Struct(config); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return formatValidationErrors(validationErrors)
		}
		return err
	}

	return nil
}

func formatValidationErrors(errors validator.ValidationErrors) error {
	var errMsgs []string
	for _, err := range errors {
		errMsgs = append(errMsgs, fmt.Sprintf(
			"%s",
			err.Error(),
		))
	}
	return fmt.Errorf("config validation error:\n%s", errMsgs)
}
