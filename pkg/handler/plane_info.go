package handler

import (
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol/pdu"
	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// HandlePlaneInfo 处理飞机信息包（#SB）
// 验证来源后原封不动转发给目标
func HandlePlaneInfo(conn gnet.Conn, p *pdu.PlaneInfo) error {
	logger.Debug("handling PlaneInfo",
		zap.String("from", p.From),
		zap.String("to", p.To),
		zap.String("type", p.Type),
		zap.String("remote", conn.RemoteAddr().String()),
	)

	// 验证登录状态和callsign
	if _, err := ValidateLoginAndCallsign(conn, p.From); err != nil {
		return err
	}

	// 转发给目标
	session.SendToCallsign(p.To, p)

	return nil
}
