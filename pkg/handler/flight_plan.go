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

	// 验证登录状态和callsign
	if _, err := ValidateLoginAndCallsign(conn, p.Callsign); err != nil {
		return err
	}

	// 检查飞行计划是否被ATC锁定
	locked, currentFP := mgr.GetFlightPlanLocked(conn)
	if locked && currentFP != nil {
		// 锁定状态下，只有DEP或ARR不同才允许提交
		if p.Dep == currentFP.Dep && p.Dest == currentFP.Dest {
			logger.Warn("FlightPlan rejected: plan is locked by ATC, DEP and ARR unchanged",
				zap.String("callsign", p.Callsign),
				zap.String("dep", p.Dep),
				zap.String("dest", p.Dest),
			)
			return nil
		}
		logger.Info("FlightPlan accepted: DEP or ARR changed, unlocking",
			zap.String("callsign", p.Callsign),
			zap.String("old_dep", currentFP.Dep),
			zap.String("new_dep", p.Dep),
			zap.String("old_dest", currentFP.Dest),
			zap.String("new_dest", p.Dest),
		)
	}

	// 存储飞行计划到Session
	fpData := session.NewFlightPlanData(p)
	mgr.UpdateFlightPlan(conn, fpData)

	// 广播给所有ATC
	session.BroadcastToATC(conn, p)

	return nil
}
