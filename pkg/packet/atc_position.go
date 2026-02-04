package packet

import (
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// HandleATCPosition 处理 ATC 位置包 (% 开头)
func HandleATCPosition(conn gnet.Conn, pkt *protocol.Packet) error {
	logger.Debug("handling ATC position packet",
		zap.String("remote", conn.RemoteAddr().String()),
		zap.ByteString("data", pkt.RawData),
	)

	// TODO: 实现 ATC 位置包的具体处理逻辑
	// 1. 解析包内容
	// 2. 更新 ATC 位置信息
	// 3. 广播给相关客户端

	return nil
}
