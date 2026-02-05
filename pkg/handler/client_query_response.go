package handler

import (
	"strings"

	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol/pdu"
	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// HandleClientQueryResponse 处理客户端查询响应包（$CR）
func HandleClientQueryResponse(conn gnet.Conn, p *pdu.ClientQueryResponse) error {
	logger.Debug("handling ClientQueryResponse",
		zap.String("from", p.From),
		zap.String("to", p.To),
		zap.String("type", p.Type),
		zap.Strings("payload", p.Payload),
		zap.String("remote", conn.RemoteAddr().String()),
	)

	// 验证登录状态和callsign
	_, err := ValidateLoginAndCallsign(conn, p.From)
	if err != nil {
		return err
	}

	to := strings.ToLower(p.To)

	// to 是 server，忽略
	if to == "server" {
		return nil
	}

	// to 以 @ 开头，转发给所有已登录用户
	if strings.HasPrefix(to, "@") {
		session.BroadcastToAll(p)
		return nil
	}

	// 默认：转发给目标
	session.SendToCallsign(p.To, p)

	return nil
}
