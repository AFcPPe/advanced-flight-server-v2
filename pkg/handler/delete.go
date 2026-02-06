package handler

import (
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol/pdu"
	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// HandleDeletePilot 处理飞行员退出
// 验证来源callsign匹配后，广播给所有用户，然后关闭连接
func HandleDeletePilot(conn gnet.Conn, p *pdu.DeletePilot) error {
	logger.Debug("handling DeletePilot",
		zap.String("callsign", p.From),
		zap.String("remote", conn.RemoteAddr().String()),
	)

	mgr := session.GetManager()

	// 验证连接类型必须是Pilot
	if mgr.GetConnType(conn) != session.ConnectionTypePilot {
		logger.Warn("DeletePilot received from non-pilot connection",
			zap.String("callsign", p.From),
		)
		return session.SendErrorAndClose(conn, p.From, pdu.NetworkErrorInvalidControl, "invalid connection type for delete pilot")
	}

	// 验证登录状态和callsign（防止伪造From踢掉别人）
	if _, err := ValidateLoginAndCallsign(conn, p.From); err != nil {
		return err
	}

	// 广播给所有已登录用户
	session.BroadcastToAll(p)

	// 关闭发送者的连接
	return conn.Close()
}

// HandleDeleteATC 处理管制员退出
// 验证来源callsign匹配后，广播给所有用户，然后关闭连接
func HandleDeleteATC(conn gnet.Conn, p *pdu.DeleteATC) error {
	logger.Debug("handling DeleteATC",
		zap.String("callsign", p.From),
		zap.String("remote", conn.RemoteAddr().String()),
	)

	mgr := session.GetManager()

	// 验证连接类型必须是ATC
	if mgr.GetConnType(conn) != session.ConnectionTypeATC {
		logger.Warn("DeleteATC received from non-ATC connection",
			zap.String("callsign", p.From),
		)
		return session.SendErrorAndClose(conn, p.From, pdu.NetworkErrorInvalidControl, "invalid connection type for delete ATC")
	}

	// 验证登录状态和callsign（防止伪造From踢掉别人）
	if _, err := ValidateLoginAndCallsign(conn, p.From); err != nil {
		return err
	}

	// 广播给所有已登录用户
	session.BroadcastToAll(p)

	// 关闭发送者的连接
	return conn.Close()
}
