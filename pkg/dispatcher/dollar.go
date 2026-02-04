package dispatcher

import (
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// DispatchDollar 分发 $ 开头的包
func DispatchDollar(conn gnet.Conn, pkt *protocol.Packet) error {
	switch pkt.SubType {
	// TODO: 根据子类型分发到具体处理函数
	// case "CQ":
	//     return packet.HandleDollarCQ(conn, pkt)
	// case "CR":
	//     return packet.HandleDollarCR(conn, pkt)
	default:
		logger.Debug("unhandled dollar packet subtype",
			zap.String("subtype", pkt.SubType),
		)
		return nil
	}
}
