package server

import (
	"bytes"
	"time"

	"advanced-flight-server/pkg/config"
	"advanced-flight-server/pkg/dispatcher"
	"advanced-flight-server/pkg/errs"
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
	s.sessionMgr.AddConn(c)
	logger.Info("new connection",
		zap.String("remote", c.RemoteAddr().String()),
		zap.Int("total_connections", s.sessionMgr.Count()),
	)

	// 发送ServerChallenge包
	challenge := pdu.NewServerChallenge()
	return []byte(pdu.Serialize(challenge.GetHeader(), challenge.ToTokens())), gnet.None
}

// OnClose 连接关闭时回调
func (s *FlightServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	// 在移除会话前，获取断链用户的信息
	sess := s.sessionMgr.GetSessionByConn(c)
	var callsign string
	var connType session.ConnectionType
	var cid string
	if sess != nil {
		callsign = sess.Callsign
		connType = sess.ConnType
		cid = sess.Cid
	}

	s.sessionMgr.RemoveConn(c)

	// 如果断链用户有callsign，广播 #DP 或 #DA 通知其他用户
	if callsign != "" {
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
	// 读取所有可用数据
	data, err := c.Next(-1)
	if err != nil {
		logger.Error("failed to read data", zap.Error(err))
		return gnet.None
	}

	if len(data) == 0 {
		return gnet.None
	}

	// 更新最后活动时间
	s.sessionMgr.UpdateLastActivity(c)

	// 获取会话
	session := s.sessionMgr.GetSessionByConn(c)
	if session == nil {
		logger.Error("session not found for connection", zap.String("remote", c.RemoteAddr().String()))
		return gnet.Close
	}

	// 追加数据到缓存，处理粘包
	if !session.AppendBuffer(data) {
		logger.Warn("buffer overflow, discarded",
			zap.String("remote", c.RemoteAddr().String()),
			zap.String("callsign", session.Callsign),
		)
		return gnet.None
	}

	// 按\r\n分割处理
	for {
		idx := bytes.Index(session.Buffer, []byte("\r\n"))
		if idx == -1 {
			// 没有完整的包，等待更多数据
			break
		}

		// 提取一个完整的包
		line := session.Buffer[:idx]
		session.Buffer = session.Buffer[idx+2:]

		if len(line) == 0 {
			continue
		}

		// 解析并分发包
		packet := protocol.ParsePacket(line)
		if err := dispatcher.Dispatch(c, packet); err != nil {
			logger.Error("failed to handle packet",
				zap.Error(err),
				zap.String("type", packet.GetTypeName()),
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
	// 检查空闲超时
	idleConns := s.sessionMgr.GetIdleConns(IdleTimeout)
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
		for _, c := range authTimeoutConns {
			logger.Warn("closing connection due to auth/login timeout",
				zap.String("remote", c.RemoteAddr().String()),
			)
			_ = c.Close()
		}
	}

	return TickInterval, gnet.None
}
