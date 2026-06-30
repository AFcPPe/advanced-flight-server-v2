package cron

import (
	"fmt"
	"sync"

	"advanced-flight-server/pkg/config"
	"advanced-flight-server/pkg/ipban"
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/metar"
	"advanced-flight-server/pkg/snapshot"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var (
	scheduler *Scheduler
	once      sync.Once
)

// Scheduler 全局定时任务调度器
type Scheduler struct {
	c *cron.Cron
}

// cronLogger 将 cron 的日志接口适配到项目的 zap logger。
// 主要用于 cron.Recover：定时任务 panic 时记录堆栈而非崩溃整个进程。
type cronLogger struct{}

func (cronLogger) Info(msg string, keysAndValues ...any) {
	logger.S().Infow(msg, keysAndValues...)
}

func (cronLogger) Error(err error, msg string, keysAndValues ...any) {
	kv := append([]any{"error", err}, keysAndValues...)
	logger.S().Errorw(msg, kv...)
}

// Init 初始化全局调度器，注册所有定时任务并启动
func Init() {
	once.Do(func() {
		scheduler = &Scheduler{
			// 挂上 Recover 中间件：任一定时任务（如快照发布）panic 时
			// 仅记录堆栈并跳过本次执行，不再让整个进程闪退。
			c: cron.New(cron.WithChain(cron.Recover(cronLogger{}))),
		}

		// 注册所有定时任务
		registerMetarSync()
		registerSnapshotPublish()
		registerIPBanReload()

		// 启动调度器
		scheduler.Start()
	})
}

// registerMetarSync 注册METAR天气同步任务：立即拉取一次，然后按配置间隔定时拉取
func registerMetarSync() {
	store := metar.GetStore()
	store.Fetch()

	cfg := config.Get().Metar
	interval := cfg.Interval
	if interval <= 0 {
		interval = 10
	}
	_ = scheduler.AddEveryMinutes("metar_sync", interval, store.Fetch)
}

// registerSnapshotPublish 注册快照发布任务：每3秒将在线用户数据推送到Redis
func registerSnapshotPublish() {
	snapshot.Init()
	pub := snapshot.GetPublisher()
	_ = scheduler.AddEverySeconds("snapshot_publish", snapshot.Interval, pub.Publish)
}

// registerIPBanReload 注册IP封禁规则重载任务：启动时立即加载一次，然后按配置间隔定时重载
func registerIPBanReload() {
	cfg := config.GetIPBan()
	if !cfg.Enabled {
		logger.Info("ip ban disabled, skip reload job")
		return
	}

	path := cfg.File
	if path == "" {
		path = "ip_ban.json"
	}

	store := ipban.GetStore()
	store.Load(path)

	interval := cfg.Interval
	if interval <= 0 {
		interval = 1
	}
	_ = scheduler.AddEveryMinutes("ip_ban_reload", interval, store.Reload)
}

// GetScheduler 获取全局调度器
func GetScheduler() *Scheduler {
	return scheduler
}

// AddEveryMinutes 添加一个按分钟间隔执行的任务
func (s *Scheduler) AddEveryMinutes(name string, minutes int, fn func()) error {
	spec := fmt.Sprintf("@every %dm", minutes)
	_, err := s.c.AddFunc(spec, fn)
	if err != nil {
		logger.Error("failed to add cron job", zap.String("name", name), zap.Error(err))
		return err
	}
	logger.Info("cron job registered", zap.String("name", name), zap.String("spec", spec))
	return nil
}

// AddEverySeconds 添加一个按秒间隔执行的任务
func (s *Scheduler) AddEverySeconds(name string, seconds int, fn func()) error {
	spec := fmt.Sprintf("@every %ds", seconds)
	_, err := s.c.AddFunc(spec, fn)
	if err != nil {
		logger.Error("failed to add cron job", zap.String("name", name), zap.Error(err))
		return err
	}
	logger.Info("cron job registered", zap.String("name", name), zap.String("spec", spec))
	return nil
}

// AddSpec 添加一个自定义cron表达式的任务
func (s *Scheduler) AddSpec(name string, spec string, fn func()) error {
	_, err := s.c.AddFunc(spec, fn)
	if err != nil {
		logger.Error("failed to add cron job", zap.String("name", name), zap.Error(err))
		return err
	}
	logger.Info("cron job registered", zap.String("name", name), zap.String("spec", spec))
	return nil
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.c.Start()
	logger.Info("cron scheduler started")
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.c.Stop()
	logger.Info("cron scheduler stopped")
}
