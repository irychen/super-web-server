package config

import "time"

type ServerConfig struct {
	Port           int           `mapstructure:"port"`
	ReadTimeout    time.Duration `mapstructure:"readTimeout"`    // 读取超时时间
	WriteTimeout   time.Duration `mapstructure:"writeTimeout"`   // 写入超时时间
	MaxHeaderBytes int           `mapstructure:"maxHeaderBytes"` // 最大头字节数
	SnowflakeNode  int64         `mapstructure:"snowflakeNode"`  // 雪花算法节点
}

type LogConfig struct {
	Level      string `mapstructure:"level" validate:"required,oneof=debug info warn error fatal panic"` // 日志级别
	Path       string `mapstructure:"path"`                                                              // 日志路径
	Filename   string `mapstructure:"filename"`                                                          // 日志文件名
	MaxSize    int    `mapstructure:"maxSize"`                                                           // 日志最大大小
	MaxBackups int    `mapstructure:"maxBackups"`                                                        // 日志最大备份数
	MaxAge     int    `mapstructure:"maxAge"`                                                            // 日志最大保存时间
	Compress   bool   `mapstructure:"compress"`                                                          // 日志是否压缩
	Stdout     bool   `mapstructure:"stdout"`                                                            // 日志是否输出到标准输出
}

type DBConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	Charset         string        `mapstructure:"charset"`
	ParseTime       bool          `mapstructure:"parseTime"`
	Timezone        string        `mapstructure:"timezone"`
	MaxIdleConns    int           `mapstructure:"maxIdleConns"`
	MaxOpenConns    int           `mapstructure:"maxOpenConns"`
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime"`
	LogLevel        string        `mapstructure:"logLevel" validate:"required,oneof=silent info warn error"` // 数据库日志级别
	SlowThreshold   time.Duration `mapstructure:"slowThreshold"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret string        `mapstructure:"secret"`
	Expire time.Duration `mapstructure:"expire"`
	Issuer string        `mapstructure:"issuer"`
}
