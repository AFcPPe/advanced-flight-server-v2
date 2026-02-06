package handler

import (
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol/pdu"
	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// HandleHandOff 处理移交包（$HO）
// 验证来源后原封不动转发给目标
func HandleHandOff(conn gnet.Conn, p *pdu.HandOff) error {
	logger.Debug("handling HandOff",
		zap.String("from", p.From),
		zap.String("to", p.To),
		zap.String("target", p.Target),
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

// HandleHandOffAccept 处理移交接受包（$HA）
// 验证来源后原封不动转发给目标
func HandleHandOffAccept(conn gnet.Conn, p *pdu.HandOffAccept) error {
	logger.Debug("handling HandOffAccept",
		zap.String("from", p.From),
		zap.String("to", p.To),
		zap.String("target", p.Target),
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
