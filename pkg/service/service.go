package service

import (
	"sync"
)

// Container 服务容器，管理所有服务实例
type Container struct {
	authService AuthService
	mu          sync.RWMutex
}

var (
	container *Container
	once      sync.Once
)

// Init 初始化服务容器
func Init() {
	once.Do(func() {
		container = &Container{
			authService: NewAuthService(),
		}
	})
}

// GetContainer 获取服务容器
func GetContainer() *Container {
	return container
}

// Auth 获取认证服务
func (c *Container) Auth() AuthService {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.authService
}

// Auth 快捷方法：获取认证服务
func Auth() AuthService {
	return GetContainer().Auth()
}
