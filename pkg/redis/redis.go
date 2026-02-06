package redis

import (
	"context"
	"sync"
	"time"

	"advanced-flight-server/pkg/config"
	"advanced-flight-server/pkg/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	client *redis.Client
	once   sync.Once
)

// Init 初始化Redis连接
func Init() {
	once.Do(func() {
		cfg := config.GetRedis()
		client = redis.NewClient(&redis.Options{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Ping(ctx).Err(); err != nil {
			logger.Error("failed to connect to redis", zap.String("addr", cfg.Addr), zap.Error(err))
			return
		}

		logger.Info("redis connected", zap.String("addr", cfg.Addr), zap.Int("db", cfg.DB))
	})
}

// GetClient 获取Redis客户端
func GetClient() *redis.Client {
	return client
}

// Close 关闭Redis连接
func Close() error {
	if client == nil {
		return nil
	}
	return client.Close()
}
