package handler

import (
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol/pdu"
	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// HandlePilotPosition 处理飞行员位置更新
func HandlePilotPosition(conn gnet.Conn, p *pdu.PilotPosition) error {
	logger.Debug("handling PilotPosition",
		zap.String("callsign", p.From),
		zap.Int("squawk", p.SquawkCode),
		zap.Float64("lat", p.Lat),
		zap.Float64("lon", p.Lon),
		zap.Int("altitude", p.TrueAltitude),
		zap.String("remote", conn.RemoteAddr().String()),
	)

	mgr := session.GetManager()

	// 验证连接类型必须是Pilot
	if mgr.GetConnType(conn) != session.ConnectionTypePilot {
		logger.Warn("PilotPosition received from non-pilot connection",
			zap.String("callsign", p.From),
		)
		return session.SendErrorAndClose(conn, p.From, pdu.NetworkErrorInvalidControl, "invalid connection type for pilot position")
	}

	// 验证登录状态和callsign
	if _, err := ValidateLoginAndCallsign(conn, p.From); err != nil {
		return err
	}

	// 更新飞行员位置信息
	mgr.UpdatePilotPosition(conn, p.Lat, p.Lon, p.SquawkCode, p.SquawkingModeC, p.Identing,
		p.TrueAltitude, p.PressureAltitude, p.GroundSpeed, p.Pitch, p.Heading, p.Bank)

	// 范围广播原始PDU
	session.BroadcastInRange(conn, p)

	return nil
}
