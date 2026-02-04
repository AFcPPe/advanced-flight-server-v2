package dispatcher

import (
	"advanced-flight-server/pkg/handler"
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol"
	"advanced-flight-server/pkg/protocol/pdu"

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
		return dispatchATCPosition(conn, pkt)
	case protocol.PacketTypePilotPosition:
		return dispatchPilotPosition(conn, pkt)
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

// dispatchATCPosition 解析并分发ATC位置包
func dispatchATCPosition(conn gnet.Conn, pkt *protocol.Packet) error {
	p, err := pdu.ATCPositionFromRaw(pkt.RawData)
	if err != nil {
		logger.Error("failed to parse ATCPosition", zap.Error(err))
		return err
	}
	return handler.HandleATCPosition(conn, p)
}

// dispatchPilotPosition 解析并分发飞行员位置包
func dispatchPilotPosition(conn gnet.Conn, pkt *protocol.Packet) error {
	p, err := pdu.PilotPositionFromRaw(pkt.RawData)
	if err != nil {
		logger.Error("failed to parse PilotPosition", zap.Error(err))
		return err
	}
	return handler.HandlePilotPosition(conn, p)
}
