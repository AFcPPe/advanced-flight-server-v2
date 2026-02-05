package handler

import (
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol/pdu"
	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// HandleFlightPlan 处理飞行计划
func HandleFlightPlan(conn gnet.Conn, p *pdu.FlightPlan) error {
	logger.Debug("handling FlightPlan",
		zap.String("callsign", p.Callsign),
		zap.String("dep", p.Dep),
		zap.String("dest", p.Dest),
		zap.String("route", p.Route),
		zap.String("remote", conn.RemoteAddr().String()),
	)

	mgr := session.GetManager()

	// 验证连接类型必须是Pilot
	if mgr.GetConnType(conn) != session.ConnectionTypePilot {
		logger.Warn("FlightPlan received from non-pilot connection",
			zap.String("callsign", p.Callsign),
		)
		return session.SendErrorAndClose(conn, p.Callsign, pdu.NetworkErrorInvalidControl, "only pilots can send flight plans")
	}

	// 验证callsign是否匹配
	if sess := mgr.GetSession(conn); sess == nil || sess.Callsign != p.Callsign {
		logger.Warn("FlightPlan callsign mismatch",
			zap.String("pdu_callsign", p.Callsign),
		)
		return session.SendErrorAndClose(conn, p.Callsign, pdu.NetworkErrorInvalidControl, "callsign mismatch")
	}

	// 存储飞行计划到Session
	fpData := session.NewFlightPlanData(p)
	mgr.UpdateFlightPlan(conn, fpData)

	// 广播给所有ATC
	session.BroadcastToATC(p)

	return nil
}
