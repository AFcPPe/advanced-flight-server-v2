package server

import (
	"bytes"
	"time"

	"advanced-flight-server/pkg/config"
	"advanced-flight-server/pkg/dispatcher"
	"advanced-flight-server/pkg/errs"
	"advanced-flight-server/pkg/ipban"
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol"
	"advanced-flight-server/pkg/protocol/pdu"
	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

const (
	// IdleTimeout 连接空闲超时时间
	IdleTimeout = 35 * time.Second
	// TickInterval OnTick检查间隔
	TickInterval = 5 * time.Second
)

// OnBoot 服务器启动时回调
func (s *FlightServer) OnBoot(eng gnet.Engine) gnet.Action {
	s.eng = eng
	logger.Info("flight server started", zap.String("addr", s.opts.Address()))
	return gnet.None
}

// OnShutdown 服务器关闭时回调
func (s *FlightServer) OnShutdown(eng gnet.Engine) {
	logger.Info("flight server shutdown")
}

// OnOpen 新连接建立时回调
func (s *FlightServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	remote := c.RemoteAddr().String()
	logger.Debug("OnOpen called",
		zap.String("remote", remote),
	)

	// IP封禁检查
	if config.GetIPBan().Enabled {
		if banAction, hit := ipban.GetStore().MatchAddr(remote); hit {
			switch banAction {
			case ipban.ActionReject:
				// 直接拒绝：连接建立后立即断开，不创建会话
				logger.Info("rejecting banned ip", zap.String("remote", remote))
				return nil, gnet.Close
			case ipban.ActionSilent:
				// 静默处理：接受连接但不回任何包，标记会话为静默
				logger.Info("silencing banned ip", zap.String("remote", remote))
				s.sessionMgr.AddConn(c)
				s.sessionMgr.SetSilenced(c)
				return nil, gnet.None
			}
		}
	}

	s.sessionMgr.AddConn(c)
	logger.Info("new connection",
		zap.String("remote", remote),
		zap.Int("total_connections", s.sessionMgr.Count()),
	)

	// 发送ServerChallenge包
	challenge := pdu.NewServerChallenge()
	return []byte(pdu.Serialize(challenge.GetHeader(), challenge.ToTokens())), gnet.None
}

// OnClose 连接关闭时回调
func (s *FlightServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	logger.Debug("OnClose called",
		zap.String("remote", c.RemoteAddr().String()),
		zap.Error(err),
	)

	// 在移除会话前，获取断链用户的信息
	sess := s.sessionMgr.GetSessionByConn(c)
	var callsign string
	var connType session.ConnectionType
	var cid string
	var authenticating bool
	if sess != nil {
		callsign = sess.Callsign
		connType = sess.ConnType
		cid = sess.Cid
		authenticating = sess.Authenticating
	}

	logger.Debug("OnClose session info before removal",
		zap.String("remote", c.RemoteAddr().String()),
		zap.String("callsign", callsign),
		zap.String("cid", cid),
		zap.Bool("authenticating", authenticating),
		zap.Bool("session_exists", sess != nil),
	)

	s.sessionMgr.RemoveConn(c)

	// 如果断链用户有callsign，广播 #DP 或 #DA 通知其他用户
	if callsign != "" {
		logger.Debug("broadcasting disconnect notification",
			zap.String("callsign", callsign),
			zap.Int("connType", int(connType)),
		)
		switch connType {
		case session.ConnectionTypePilot:
			dp := pdu.NewPDUDeletePilot(callsign, cid)
			session.BroadcastToAll(c, dp)
		case session.ConnectionTypeATC:
			da := pdu.NewPDUDeleteATC(callsign, cid)
			session.BroadcastToAll(c, da)
		}
	}

	logger.Info("connection closed, session cleaned up",
		zap.String("remote", c.RemoteAddr().String()),
		zap.String("callsign", callsign),
		zap.Error(err),
		zap.Int("total_connections", s.sessionMgr.Count()),
		zap.Int("total_callsigns", s.sessionMgr.CallsignCount()),
	)
	return gnet.None
}

// OnTraffic 收到数据时回调
func (s *FlightServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	// 兜底恢复：单个包/单条连接的处理 panic 不能掀翻整个 gnet 引擎。
	// gnet 在 event-loop goroutine panic 后会触发 OnShutdown 并强制关闭所有连接，
	// 表现为整服闪退。这里捕获 panic、打印真实堆栈，仅关闭出问题的连接。
	defer func() {
		if r := recover(); r != nil {
			var remote string
			if addr := c.RemoteAddr(); addr != nil {
				remote = addr.String()
			}
			logger.Error("PANIC recovered in OnTraffic",
				zap.Any("recover", r),
				zap.String("remote", remote),
				zap.Stack("stack"),
			)
			action = gnet.Close
		}
	}()

	// 读取所有可用数据
	data, err := c.Next(-1)
	if err != nil {
		logger.Error("failed to read data", zap.Error(err))
		return gnet.None
	}

	if len(data) == 0 {
		return gnet.None
	}

	// 静默连接：丢弃所有输入，不回任何包，伪装服务器崩溃
	sess := s.sessionMgr.GetSessionByConn(c)
	if sess != nil && sess.Silenced {
		return gnet.None
	}

	logger.Debug("OnTraffic received data",
		zap.String("remote", c.RemoteAddr().String()),
		zap.Int("data_len", len(data)),
	)

	// 更新最后活动时间
	s.sessionMgr.UpdateLastActivity(c)

	// 获取会话
	if sess == nil {
		logger.Error("session not found for connection", zap.String("remote", c.RemoteAddr().String()))
		return gnet.Close
	}

	// 追加数据到缓存，处理粘包
	if !sess.AppendBuffer(data) {
		logger.Warn("buffer overflow, discarded",
			zap.String("remote", c.RemoteAddr().String()),
			zap.String("callsign", sess.Callsign),
		)
		return gnet.None
	}

	// 按\r\n分割处理
	for {
		idx := bytes.Index(sess.Buffer, []byte("\r\n"))
		if idx == -1 {
			// 没有完整的包，等待更多数据
			break
		}

		// 提取一个完整的包
		line := sess.Buffer[:idx]
		sess.Buffer = sess.Buffer[idx+2:]

		if len(line) == 0 {
			continue
		}

		// 解析并分发包
		packet := protocol.ParsePacket(line)
		logger.Debug("dispatching packet from OnTraffic",
			zap.String("remote", c.RemoteAddr().String()),
			zap.String("callsign", sess.Callsign),
			zap.String("type", packet.GetTypeName()),
		)
		if err := dispatcher.Dispatch(c, packet); err != nil {
			logger.Error("failed to handle packet",
				zap.Error(err),
				zap.String("type", packet.GetTypeName()),
				zap.String("remote", c.RemoteAddr().String()),
				zap.String("callsign", sess.Callsign),
			)
			// 如果首包不是$ID，断开连接
			if err == errs.ErrNotAuthenticated {
				return gnet.Close
			}
		}
	}

	return gnet.None
}

// OnTick 定时检查空闲连接并断开
func (s *FlightServer) OnTick() (delay time.Duration, action gnet.Action) {
	logger.Debug("OnTick called",
		zap.Int("total_connections", s.sessionMgr.Count()),
		zap.Int("total_callsigns", s.sessionMgr.CallsignCount()),
	)

	// 检查空闲超时
	idleConns := s.sessionMgr.GetIdleConns(IdleTimeout)
	if len(idleConns) > 0 {
		logger.Debug("found idle connections", zap.Int("count", len(idleConns)))
	}
	for _, c := range idleConns {
		callsign := s.sessionMgr.GetCallsignByConn(c)
		logger.Warn("closing idle connection",
			zap.String("remote", c.RemoteAddr().String()),
			zap.String("callsign", callsign),
		)
		_ = c.Close()
	}

	// 检查认证/登录超时
	authTimeout := time.Duration(config.GetServer().AuthTimeout) * time.Second
	if authTimeout > 0 {
		authTimeoutConns := s.sessionMgr.GetAuthTimeoutConns(authTimeout)
		if len(authTimeoutConns) > 0 {
			logger.Debug("found auth timeout connections", zap.Int("count", len(authTimeoutConns)))
		}
		for _, c := range authTimeoutConns {
			logger.Warn("closing connection due to auth/login timeout",
				zap.String("remote", c.RemoteAddr().String()),
			)
			_ = c.Close()
		}
	}

	// 每30秒向每个ATC单独询问ATIS
	now := time.Now()
	sessions := s.sessionMgr.GetAllSessions()
	for _, sess := range sessions {
		if sess.IsLoggedIn() && sess.ConnType == session.ConnectionTypeATC {
			if now.Sub(sess.LastAtisQuery) >= 30*time.Second {
				sess.LastAtisQuery = now
				logger.Debug("querying ATIS from ATC",
					zap.String("callsign", sess.Callsign),
				)
				_ = session.Send(sess.Conn, pdu.NewPDUClientQuery("SERVER", sess.Callsign, "ATIS", nil))
			}
		}
	}

	logger.Debug("OnTick completed")
	return TickInterval, gnet.None
}
