package handler

import (
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol/pdu"
	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// HandleATCPosition 处理ATC位置更新
func HandleATCPosition(conn gnet.Conn, p *pdu.ATCPosition) error {
	logger.Debug("handling ATCPosition",
		zap.String("callsign", p.From),
		zap.Strings("frequencies", p.Frequencies),
		zap.Float64("lat", p.Lat),
		zap.Float64("lon", p.Lon),
		zap.String("remote", conn.RemoteAddr().String()),
	)

	mgr := session.GetManager()

	// 验证连接类型必须是ATC
	if mgr.GetConnType(conn) != session.ConnectionTypeATC {
		logger.Warn("ATCPosition received from non-ATC connection",
			zap.String("callsign", p.From),
		)
		return session.SendErrorAndClose(conn, p.From, pdu.NetworkErrorInvalidControl, "invalid connection type for ATC position")
	}

	// 验证callsign是否匹配
	if sess := mgr.GetSession(conn); sess == nil || sess.Callsign != p.From {
		logger.Warn("ATCPosition callsign mismatch",
			zap.String("pdu_from", p.From),
		)
		return session.SendErrorAndClose(conn, p.From, pdu.NetworkErrorInvalidControl, "callsign mismatch")
	}

	// 更新ATC位置信息
	mgr.UpdateATCPosition(conn, p.Lat, p.Lon, p.Frequencies, int(p.Facility), p.VisibilityRange)

	// 范围广播原始PDU
	session.BroadcastInRange(conn, p)

	return nil
}
