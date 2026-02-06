package cron

import (
	"fmt"
	"sync"

	"advanced-flight-server/pkg/config"
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

// Init 初始化全局调度器，注册所有定时任务并启动
func Init() {
	once.Do(func() {
		scheduler = &Scheduler{
			c: cron.New(),
		}

		// 注册所有定时任务
		registerMetarSync()
		registerSnapshotPublish()

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
