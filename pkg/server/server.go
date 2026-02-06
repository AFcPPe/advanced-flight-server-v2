package server

import (
	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
)

// FlightServer gnet TCP 服务器
type FlightServer struct {
	gnet.BuiltinEventEngine
	eng        gnet.Engine
	opts       *Options
	sessionMgr *session.Manager
}

// NewFlightServer 创建新的服务器实例
func NewFlightServer(host string, port int) *FlightServer {
	opts := &Options{
		Host:      host,
		Port:      port,
		Multicore: true,
		ReusePort: true,
	}
	return &FlightServer{
		opts:       opts,
		sessionMgr: session.GetManager(),
	}
}

// NewFlightServerWithOptions 使用自定义选项创建服务器
func NewFlightServerWithOptions(opts *Options) *FlightServer {
	return &FlightServer{
		opts:       opts,
		sessionMgr: session.GetManager(),
	}
}

// Run 启动服务器
func (s *FlightServer) Run() error {
	return gnet.Run(s, s.opts.Address(),
		gnet.WithMulticore(s.opts.Multicore),
		gnet.WithReusePort(s.opts.ReusePort),
		gnet.WithTicker(true),
	)
}

// GetEngine 获取 gnet 引擎
func (s *FlightServer) GetEngine() gnet.Engine {
	return s.eng
}

// GetSessionManager 获取会话管理器
func (s *FlightServer) GetSessionManager() *session.Manager {
	return s.sessionMgr
}

// GetOptions 获取服务器配置
func (s *FlightServer) GetOptions() *Options {
	return s.opts
}
