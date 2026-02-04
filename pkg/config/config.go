package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	globalConfig *Config
	once         sync.Once
)

// Init 初始化配置
// configPath: 配置文件路径，如 "config.yaml" 或 "configs/config.yaml"
func Init(configPath string) error {
	var initErr error
	once.Do(func() {
		if configPath == "" {
			configPath = "config.yaml"
		}

		// 检查配置文件是否存在，不存在则生成默认配置文件
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			if err := generateDefaultConfig(configPath); err != nil {
				initErr = fmt.Errorf("failed to generate default config: %w", err)
				return
			}
			fmt.Printf("配置文件不存在，已生成默认配置文件: %s\n", configPath)
		}

		viper.SetConfigFile(configPath)
		viper.SetConfigType("yaml")

		// 设置默认值
		setDefaults()

		// 读取配置文件
		if err := viper.ReadInConfig(); err != nil {
			initErr = fmt.Errorf("failed to read config file: %w", err)
			return
		}

		// 解析到结构体
		globalConfig = &Config{}
		if err := viper.Unmarshal(globalConfig); err != nil {
			initErr = fmt.Errorf("failed to unmarshal config: %w", err)
			return
		}
	})
	return initErr
}

// setDefaults 设置默认值
func setDefaults() {
	defaults := DefaultConfig()

	// App
	viper.SetDefault("app.name", defaults.App.Name)
	viper.SetDefault("app.version", defaults.App.Version)
	viper.SetDefault("app.env", defaults.App.Env)

	// Logger
	viper.SetDefault("logger.level", defaults.Logger.Level)
	viper.SetDefault("logger.filename", defaults.Logger.Filename)
	viper.SetDefault("logger.max_size", defaults.Logger.MaxSize)
	viper.SetDefault("logger.max_backups", defaults.Logger.MaxBackups)
	viper.SetDefault("logger.max_age", defaults.Logger.MaxAge)
	viper.SetDefault("logger.compress", defaults.Logger.Compress)
	viper.SetDefault("logger.console", defaults.Logger.Console)

	// Server
	viper.SetDefault("server.host", defaults.Server.Host)
	viper.SetDefault("server.port", defaults.Server.Port)
}

// generateDefaultConfig 使用 viper 生成默认配置文件
func generateDefaultConfig(configPath string) error {
	setDefaults()
	return viper.SafeWriteConfigAs(configPath)
}

// Get 获取全局配置
func Get() *Config {
	if globalConfig == nil {
		// 如果未初始化，使用默认配置
		return DefaultConfig()
	}
	return globalConfig
}

// GetApp 获取应用配置
func GetApp() *AppConfig {
	return &Get().App
}

// GetLogger 获取日志配置
func GetLogger() *LoggerConfig {
	return &Get().Logger
}

// GetServer 获取服务器配置
func GetServer() *ServerConfig {
	return &Get().Server
}

// Reload 重新加载配置
func Reload() error {
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}
	if err := viper.Unmarshal(globalConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return nil
}

// WatchConfig 监听配置文件变化
func WatchConfig(onChange func()) {
	viper.OnConfigChange(func(e fsnotify.Event) {
		_ = viper.Unmarshal(globalConfig)
		if onChange != nil {
			onChange()
		}
	})
	viper.WatchConfig()
}
