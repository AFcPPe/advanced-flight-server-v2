package handler

import (
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol/pdu"
	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// HandleControllerInfo 处理管制员信息包（#PC）
// 验证来源后原封不动转发给目标
func HandleControllerInfo(conn gnet.Conn, p *pdu.ControllerInfo) error {
	logger.Debug("handling ControllerInfo",
		zap.String("from", p.From),
		zap.String("to", p.To),
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
