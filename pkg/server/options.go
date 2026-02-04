package server

import (
	"fmt"
)

// Options 服务器配置选项
type Options struct {
	Host      string
	Port      int
	Multicore bool
	ReusePort bool
}

// DefaultOptions 返回默认配置
func DefaultOptions() *Options {
	return &Options{
		Host:      "0.0.0.0",
		Port:      6809,
		Multicore: true,
		ReusePort: true,
	}
}

// Address 返回监听地址
func (o *Options) Address() string {
	return fmt.Sprintf("tcp://%s:%d", o.Host, o.Port)
}
