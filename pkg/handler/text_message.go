package handler

import (
	"strings"

	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol/pdu"
	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// HandleTextMessage 处理文本消息包（#TM）
// 验证来源后根据目标转发
func HandleTextMessage(conn gnet.Conn, p *pdu.TextMessage) error {
	logger.Debug("handling TextMessage",
		zap.String("from", p.From),
		zap.String("to", p.To),
		zap.String("message", p.Message),
		zap.String("remote", conn.RemoteAddr().String()),
	)

	// 验证登录状态和callsign
	sess, err := ValidateLoginAndCallsign(conn, p.From)
	if err != nil {
		return err
	}

	to := strings.ToUpper(p.To)

	// *S：转发给所有Supervisor及以上（SUP、ADM）
	if to == "*S" {
		session.BroadcastToSupervisors(conn, p)
		return nil
	}

	// *：广播给所有用户，仅SUP及以上可用
	if to == "*" {
		if sess.Rating < int(pdu.RatingSUP) {
			return nil
		}
		session.BroadcastToAll(conn, p)
		return nil
	}

	// @开头：广播给所有已登录用户
	if strings.HasPrefix(to, "@") {
		session.BroadcastToAll(conn, p)
		return nil
	}

	// *A：广播给所有ATC
	if to == "*A" {
		session.BroadcastToATC(conn, p)
		return nil
	}

	// FP：忽略
	if to == "FP" {
		return nil
	}

	// 默认：转发给目标callsign
	session.SendToCallsign(p.To, p)

	return nil
}
