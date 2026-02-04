package database

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"advanced-flight-server/pkg/config"
	"advanced-flight-server/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// 数据库实例名称常量
const (
	DBAccount = "account" // 账户数据库
)

var (
	instances = make(map[string]*gorm.DB)
	mu        sync.RWMutex
)

// ErrInstanceNotFound 数据库实例未找到
var ErrInstanceNotFound = errors.New("database instance not found")

// Register 注册数据库实例
func Register(name string, cfg *config.DatabaseAccountConfig) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := instances[name]; exists {
		return fmt.Errorf("database instance '%s' already registered", name)
	}

	conn, err := connect(cfg)
	if err != nil {
		return err
	}

	instances[name] = conn
	logger.Info("database registered",
		zap.String("name", name),
		zap.String("driver", cfg.Driver),
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.Database),
	)
	return nil
}

// connect 创建数据库连接
func connect(cfg *config.DatabaseAccountConfig) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.DSN())
	case "postgres":
		dialector = postgres.Open(cfg.DSN())
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	gormConfig := &gorm.Config{
		Logger:                                   newGormLogger(cfg),
		DisableForeignKeyConstraintWhenMigrating: true,
		PrepareStmt:                              true,
	}

	conn, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := conn.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Second)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return conn, nil
}

// Get 获取指定名称的数据库实例
func Get(name string) (*gorm.DB, error) {
	mu.RLock()
	defer mu.RUnlock()

	db, exists := instances[name]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrInstanceNotFound, name)
	}
	return db, nil
}

// Account 获取账户数据库实例的快捷方法
func Account() (*gorm.DB, error) {
	return Get(DBAccount)
}

// Close 关闭指定数据库连接
func Close(name string) error {
	mu.Lock()
	defer mu.Unlock()

	db, exists := instances[name]
	if !exists {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	delete(instances, name)
	return sqlDB.Close()
}

// CloseAll 关闭所有数据库连接
func CloseAll() error {
	mu.Lock()
	defer mu.Unlock()

	var lastErr error
	for name, db := range instances {
		sqlDB, err := db.DB()
		if err != nil {
			lastErr = err
			continue
		}
		if err := sqlDB.Close(); err != nil {
			lastErr = err
		}
		delete(instances, name)
		logger.Info("database closed", zap.String("name", name))
	}
	return lastErr
}

// Transaction 在指定数据库上执行事务
func Transaction(name string, fn func(tx *gorm.DB) error) error {
	db, err := Get(name)
	if err != nil {
		return err
	}
	return db.Transaction(fn)
}

// gormLogger 适配器
type gormLoggerAdapter struct {
	slowThreshold time.Duration
	logLevel      gormlogger.LogLevel
}

func newGormLogger(cfg *config.DatabaseAccountConfig) gormlogger.Interface {
	var level gormlogger.LogLevel
	switch cfg.LogLevel {
	case "silent":
		level = gormlogger.Silent
	case "error":
		level = gormlogger.Error
	case "warn":
		level = gormlogger.Warn
	case "info":
		level = gormlogger.Info
	default:
		level = gormlogger.Warn
	}

	return &gormLoggerAdapter{
		slowThreshold: time.Duration(cfg.SlowThreshold) * time.Millisecond,
		logLevel:      level,
	}
}

func (l *gormLoggerAdapter) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

func (l *gormLoggerAdapter) Info(_ context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Info {
		logger.Infof(msg, data...)
	}
}

func (l *gormLoggerAdapter) Warn(_ context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Warn {
		logger.Warnf(msg, data...)
	}
}

func (l *gormLoggerAdapter) Error(_ context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Error {
		logger.Errorf(msg, data...)
	}
}

func (l *gormLoggerAdapter) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.logLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && l.logLevel >= gormlogger.Error:
		logger.Error("gorm trace",
			zap.Error(err),
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= gormlogger.Warn:
		logger.Warn("slow sql",
			zap.Duration("elapsed", elapsed),
			zap.Duration("threshold", l.slowThreshold),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case l.logLevel >= gormlogger.Info:
		logger.Debug("gorm trace",
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	}
}
