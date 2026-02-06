package handler

import (
	"time"

	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol"
	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// HandleClientIdentification 处理$ID包（客户端身份验证响应）
func HandleClientIdentification(conn gnet.Conn, pkt *protocol.Packet) error {
	logger.Debug("handling $ID packet",
		zap.String("remote", conn.RemoteAddr().String()),
		zap.ByteString("raw", pkt.RawData),
	)

	// 获取会话并设置已认证
	sess := session.GetManager().GetSessionByConn(conn)
	if sess == nil {
		return nil
	}

	// 标记为已认证
	sess.Authenticated = true
	sess.AuthenticatedTime = time.Now()

	// TODO: 解析$ID包内容并验证
	// TODO: 验证客户端版本等信息

	return nil
}
