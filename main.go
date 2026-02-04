package main

import (
	"advanced-flight-server/pkg/config"
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/server"

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

	// 获取配置
	appCfg := config.GetApp()
	serverCfg := config.GetServer()

	logger.Info("application started",
		zap.String("name", appCfg.Name),
		zap.String("version", appCfg.Version),
		zap.String("env", appCfg.Env),
	)

	// 启动 TCP 服务器
	srv := server.NewFlightServer(serverCfg.Host, serverCfg.Port)
	logger.Infof("starting TCP server on %s:%d", serverCfg.Host, serverCfg.Port)

	if err := srv.Run(); err != nil {
		logger.Error("server error", zap.Error(err))
		panic(err)
	}
}
