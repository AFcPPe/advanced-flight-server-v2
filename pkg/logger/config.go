package logger

// Config 日志配置
type Config struct {
	Level        string // 日志级别: debug, info, warn, error
	Filename     string // 日志文件路径
	MaxSize      int    // 单个日志文件最大大小(MB)
	MaxBackups   int    // 保留的旧日志文件最大数量
	MaxAge       int    // 保留的旧日志文件最大天数
	Compress     bool   // 是否压缩旧日志文件
	Console      bool   // 是否同时输出到控制台
	RotateByDate bool   // 是否按日期自动分文件
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:        "info",
		Filename:     "logs/app.log",
		MaxSize:      100,
		MaxBackups:   3,
		MaxAge:       7,
		Compress:     true,
		Console:      true,
		RotateByDate: true,
	}
}
