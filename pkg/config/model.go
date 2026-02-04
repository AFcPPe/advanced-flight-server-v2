package config

// Config 全局配置结构
type Config struct {
	App    AppConfig    `mapstructure:"app"`
	Logger LoggerConfig `mapstructure:"logger"`
	Server ServerConfig `mapstructure:"server"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Env     string `mapstructure:"env"` // dev, test, prod
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `mapstructure:"level"`       // 日志级别: debug, info, warn, error
	Filename   string `mapstructure:"filename"`    // 日志文件路径
	MaxSize    int    `mapstructure:"max_size"`    // 单个日志文件最大大小(MB)
	MaxBackups int    `mapstructure:"max_backups"` // 保留的旧日志文件最大数量
	MaxAge     int    `mapstructure:"max_age"`     // 保留的旧日志文件最大天数
	Compress   bool   `mapstructure:"compress"`    // 是否压缩旧日志文件
	Console    bool   `mapstructure:"console"`     // 是否同时输出到控制台
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:    "advanced-flight-server",
			Version: "1.0.0",
			Env:     "dev",
		},
		Logger: LoggerConfig{
			Level:      "info",
			Filename:   "logs/app.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   true,
			Console:    true,
		},
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
	}
}
