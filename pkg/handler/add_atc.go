package handler

import (
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol/pdu"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// HandleAddATC 处理管制员登录
func HandleAddATC(conn gnet.Conn, p *pdu.AddATC) error {
	// 打印完整一点
	logger.Debug("handling AddATC",
		zap.String("callsign", p.Callsign),
		zap.String("cid", p.Cid),
		zap.String("realname", p.RealName),
		zap.String("remote", conn.RemoteAddr().String()),
		zap.String("password", p.Password),
		zap.Int("rating", int(p.Rating)),
		zap.Int("protocol_version", p.ProtocolVersion),
	)
	// TODO: 密码验证
	// TODO: 等级验证
	// TODO: 创建session
	// TODO: 广播上线消息

	return nil
}
