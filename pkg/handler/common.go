package handler

import (
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol/pdu"
	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// ValidateLogin 验证用户是否已登录，返回session和错误
// 如果验证失败，会自动发送错误并关闭连接
func ValidateLogin(conn gnet.Conn, callsign string) (*session.Session, error) {
	mgr := session.GetManager()
	sess := mgr.GetSession(conn)
	if sess == nil || !sess.IsLoggedIn() {
		logger.Warn("received packet from unauthenticated connection",
			zap.String("callsign", callsign),
		)
		return nil, session.SendErrorAndClose(conn, callsign, pdu.NetworkErrorInvalidLogon, "not logged in")
	}
	return sess, nil
}

// ValidateLoginAndCallsign 验证用户是否已登录且callsign匹配
// 如果验证失败，会自动发送错误并关闭连接
func ValidateLoginAndCallsign(conn gnet.Conn, callsign string) (*session.Session, error) {
	sess, err := ValidateLogin(conn, callsign)
	if err != nil {
		return nil, err
	}

	if sess.Callsign != callsign {
		logger.Warn("callsign mismatch",
			zap.String("session_callsign", sess.Callsign),
			zap.String("pdu_callsign", callsign),
		)
		return nil, session.SendErrorAndClose(conn, callsign, pdu.NetworkErrorInvalidControl, "callsign mismatch")
	}

	return sess, nil
}
