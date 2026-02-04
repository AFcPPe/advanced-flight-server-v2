package config

// Config 全局配置结构
type Config struct {
	App              AppConfig                        `mapstructure:"app"`
	Logger           LoggerConfig                     `mapstructure:"logger"`
	Server           ServerConfig                     `mapstructure:"server"`
	DatabaseAccounts map[string]DatabaseAccountConfig `mapstructure:"database_accounts"`
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
	Motd string `mapstructure:"motd"` // 登录成功后发送的欢迎消息
}

// DatabaseAccountConfig 数据库账户配置
type DatabaseAccountConfig struct {
	Driver          string `mapstructure:"driver"`             // 数据库驱动: mysql, postgres
	Host            string `mapstructure:"host"`               // 数据库主机
	Port            int    `mapstructure:"port"`               // 数据库端口
	Username        string `mapstructure:"username"`           // 用户名
	Password        string `mapstructure:"password"`           // 密码
	Database        string `mapstructure:"database"`           // 数据库名
	Charset         string `mapstructure:"charset"`            // 字符集
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`     // 最大空闲连接数
	MaxOpenConns    int    `mapstructure:"max_open_conns"`     // 最大打开连接数
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`  // 连接最大生命周期(秒)
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time"` // 连接最大空闲时间(秒)
	LogLevel        string `mapstructure:"log_level"`          // 日志级别: silent, error, warn, info
	SlowThreshold   int    `mapstructure:"slow_threshold"`     // 慢查询阈值(毫秒)
}

// DSN 生成数据库连接字符串
func (c *DatabaseAccountConfig) DSN() string {
	switch c.Driver {
	case "mysql":
		return c.mysqlDSN()
	case "postgres":
		return c.postgresDSN()
	default:
		return c.mysqlDSN()
	}
}

func (c *DatabaseAccountConfig) mysqlDSN() string {
	return c.Username + ":" + c.Password + "@tcp(" + c.Host + ":" + itoa(c.Port) + ")/" +
		c.Database + "?charset=" + c.Charset + "&parseTime=True&loc=Local"
}

func (c *DatabaseAccountConfig) postgresDSN() string {
	return "host=" + c.Host + " user=" + c.Username + " password=" + c.Password +
		" dbname=" + c.Database + " port=" + itoa(c.Port) + " sslmode=disable TimeZone=Asia/Shanghai"
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var b [20]byte
	pos := len(b)
	for i > 0 {
		pos--
		b[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(b[pos:])
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
			Port: 6809,
			Motd: "Welcome to Advanced Flight Server!",
		},
		DatabaseAccounts: map[string]DatabaseAccountConfig{
			"account": {
				Driver:          "mysql",
				Host:            "127.0.0.1",
				Port:            3306,
				Username:        "root",
				Password:        "",
				Database:        "flight_server",
				Charset:         "utf8mb4",
				MaxIdleConns:    10,
				MaxOpenConns:    100,
				ConnMaxLifetime: 3600,
				ConnMaxIdleTime: 600,
				LogLevel:        "warn",
				SlowThreshold:   200,
			},
		},
	}
}
