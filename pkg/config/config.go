package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
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
	// 使用反射自动将结构体转换为viper默认值
	if err := viper.MergeConfigMap(structToMap(defaults)); err != nil {
		// 静默处理，使用空配置
		return
	}
}

// structToMap 将结构体递归转换为map[string]interface{}
func structToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return result
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// 获取mapstructure tag作为key
		key := field.Tag.Get("mapstructure")
		if key == "" {
			key = strings.ToLower(field.Name)
		}

		result[key] = convertValue(value)
	}
	return result
}

// convertValue 转换reflect.Value为interface{}
func convertValue(v reflect.Value) interface{} {
	switch v.Kind() {
	case reflect.Struct:
		return structToMap(v.Interface())
	case reflect.Map:
		return mapToInterface(v)
	case reflect.Slice, reflect.Array:
		return sliceToInterface(v)
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		return convertValue(v.Elem())
	default:
		return v.Interface()
	}
}

// mapToInterface 将map转换为map[string]interface{}
func mapToInterface(v reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})
	for _, key := range v.MapKeys() {
		keyStr := fmt.Sprintf("%v", key.Interface())
		result[keyStr] = convertValue(v.MapIndex(key))
	}
	return result
}

// sliceToInterface 将slice转换为[]interface{}
func sliceToInterface(v reflect.Value) []interface{} {
	result := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		result[i] = convertValue(v.Index(i))
	}
	return result
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

// GetRedis 获取Redis配置
func GetRedis() *RedisConfig {
	return &Get().Redis
}

// GetIPBan 获取IP封禁配置
func GetIPBan() *IPBanConfig {
	return &Get().IPBan
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
