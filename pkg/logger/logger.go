package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	globalLogger *zap.Logger
	sugarLogger  *zap.SugaredLogger
	once         sync.Once
)

// dateRotateWriter 按日期自动切换日志文件的 WriteSyncer
type dateRotateWriter struct {
	mu          sync.Mutex
	cfg         *Config
	currentDate string
	writer      *lumberjack.Logger
}

// newDateRotateWriter 创建按日期分文件的 writer
func newDateRotateWriter(cfg *Config) *dateRotateWriter {
	w := &dateRotateWriter{cfg: cfg}
	w.rotate()
	return w
}

// dateFilename 根据日期生成文件名，例如 logs/app-2026-02-11.log
func (w *dateRotateWriter) dateFilename(date string) string {
	ext := filepath.Ext(w.cfg.Filename)
	name := strings.TrimSuffix(w.cfg.Filename, ext)
	return fmt.Sprintf("%s-%s%s", name, date, ext)
}

// rotate 切换到当天日期对应的日志文件
func (w *dateRotateWriter) rotate() {
	now := time.Now().Format("2006-01-02")
	w.currentDate = now
	if w.writer != nil {
		_ = w.writer.Close()
	}
	w.writer = &lumberjack.Logger{
		Filename:   w.dateFilename(now),
		MaxSize:    w.cfg.MaxSize,
		MaxBackups: w.cfg.MaxBackups,
		MaxAge:     w.cfg.MaxAge,
		Compress:   w.cfg.Compress,
	}
}

func (w *dateRotateWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now().Format("2006-01-02")
	if now != w.currentDate {
		w.rotate()
	}
	return w.writer.Write(p)
}

func (w *dateRotateWriter) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.writer != nil {
		return w.writer.Close()
	}
	return nil
}

// Init 初始化全局 logger
func Init(cfg *Config) {
	once.Do(func() {
		if cfg == nil {
			cfg = DefaultConfig()
		}

		// 解析日志级别
		level := parseLevel(cfg.Level)

		// 编码器配置
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		var cores []zapcore.Core

		// 文件输出
		if cfg.Filename != "" {
			var fileWriteSyncer zapcore.WriteSyncer
			if cfg.RotateByDate {
				// 按日期自动分文件
				fileWriteSyncer = zapcore.AddSync(newDateRotateWriter(cfg))
			} else {
				fileWriteSyncer = zapcore.AddSync(&lumberjack.Logger{
					Filename:   cfg.Filename,
					MaxSize:    cfg.MaxSize,
					MaxBackups: cfg.MaxBackups,
					MaxAge:     cfg.MaxAge,
					Compress:   cfg.Compress,
				})
			}
			fileCore := zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig),
				fileWriteSyncer,
				level,
			)
			cores = append(cores, fileCore)
		}

		// 控制台输出
		if cfg.Console {
			consoleEncoderConfig := encoderConfig
			consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			consoleCore := zapcore.NewCore(
				zapcore.NewConsoleEncoder(consoleEncoderConfig),
				zapcore.AddSync(os.Stdout),
				level,
			)
			cores = append(cores, consoleCore)
		}

		// 创建 logger
		core := zapcore.NewTee(cores...)
		globalLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		sugarLogger = globalLogger.Sugar()
	})
}

// parseLevel 解析日志级别
func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// L 获取全局 Logger
func L() *zap.Logger {
	if globalLogger == nil {
		Init(nil)
	}
	return globalLogger
}

// S 获取全局 SugaredLogger
func S() *zap.SugaredLogger {
	if sugarLogger == nil {
		Init(nil)
	}
	return sugarLogger
}

// Sync 刷新日志缓冲
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}
