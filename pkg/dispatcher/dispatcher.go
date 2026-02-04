package dispatcher

import (
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/packet"
	"advanced-flight-server/pkg/protocol"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// Dispatch 根据包类型分发到对应的处理函数
func Dispatch(conn gnet.Conn, pkt *protocol.Packet) error {
	logger.Debug("dispatching packet",
		zap.String("type", pkt.GetTypeName()),
		zap.String("remote", conn.RemoteAddr().String()),
	)

	switch pkt.Type {
	case protocol.PacketTypeATCPosition:
		return packet.HandleATCPosition(conn, pkt)
	case protocol.PacketTypePilotPosition:
		return packet.HandlePilotPosition(conn, pkt)
	case protocol.PacketTypeHash:
		return DispatchHash(conn, pkt)
	case protocol.PacketTypeDollar:
		return DispatchDollar(conn, pkt)
	default:
		logger.Warn("unknown packet type",
			zap.ByteString("raw", pkt.RawData),
		)
		return nil
	}
}
