package handler

import (
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol/pdu"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// HandleAddPilot 处理飞行员登录
func HandleAddPilot(conn gnet.Conn, p *pdu.AddPilot) error {
	logger.Debug("handling AddPilot",
		zap.String("callsign", p.Callsign),
		zap.String("cid", p.Cid),
		zap.String("realname", p.RealName),
		zap.Int("simtype", p.SimType),
		zap.String("remote", conn.RemoteAddr().String()),
	)

	// TODO: 密码验证
	// TODO: 等级验证
	// TODO: 创建session
	// TODO: 广播上线消息

	return nil
}
