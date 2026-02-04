package main

import (
	"advanced-flight-server/pkg/config"
	"advanced-flight-server/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// 初始化配置（不存在则自动生成 config.yaml）
	if err := config.Init("config.yaml"); err != nil {
		panic(err)
	}

	// 用配置初始化 logger
	logCfg := config.GetLogger()
	logger.Init(&logger.Config{
		Level:      logCfg.Level,
		Filename:   logCfg.Filename,
		MaxSize:    logCfg.MaxSize,
		MaxBackups: logCfg.MaxBackups,
		MaxAge:     logCfg.MaxAge,
		Compress:   logCfg.Compress,
		Console:    logCfg.Console,
	})
	defer logger.Sync()

	// 测试日志输出
	appCfg := config.GetApp()
	serverCfg := config.GetServer()

	logger.Info("application started",
		zap.String("name", appCfg.Name),
		zap.String("version", appCfg.Version),
		zap.String("env", appCfg.Env),
	)

	logger.Infof("server listening on %s:%d", serverCfg.Host, serverCfg.Port)
	logger.Debug("this is a debug message")
	logger.Warn("this is a warning message")
}
