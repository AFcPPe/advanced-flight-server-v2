package handler

import (
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol/pdu"
	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// HandleTagModify 处理ATC修改Pilot飞行计划（$AM）
func HandleTagModify(conn gnet.Conn, p *pdu.TagModify) error {
	logger.Debug("handling TagModify",
		zap.String("from", p.From),
		zap.String("to", p.To),
		zap.String("callsign", p.Callsign),
		zap.String("dep", p.Dep),
		zap.String("dest", p.Dest),
		zap.String("remote", conn.RemoteAddr().String()),
	)

	mgr := session.GetManager()

	// 验证连接类型必须是ATC
	if mgr.GetConnType(conn) != session.ConnectionTypeATC {
		logger.Warn("TagModify received from non-ATC connection",
			zap.String("from", p.From),
		)
		return session.SendErrorAndClose(conn, p.From, pdu.NetworkErrorInvalidControl, "only ATC can send tag modify")
	}

	// 验证登录状态和callsign
	if _, err := ValidateLoginAndCallsign(conn, p.From); err != nil {
		return err
	}

	// 更新目标Pilot的飞行计划并锁定
	fpData := session.NewFlightPlanDataFromTagModify(p)
	if !mgr.UpdateFlightPlanByCallsign(p.Callsign, fpData) {
		logger.Warn("TagModify target pilot not found",
			zap.String("callsign", p.Callsign),
		)
	}

	// 广播给所有ATC
	session.BroadcastToATC(conn, p)

	return nil
}
