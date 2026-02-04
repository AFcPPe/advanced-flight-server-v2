package dispatcher

import (
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// DispatchHash 分发 # 开头的包
func DispatchHash(conn gnet.Conn, pkt *protocol.Packet) error {
	switch pkt.SubType {
	// TODO: 根据子类型分发到具体处理函数
	// case "AP":
	//     return packet.HandleHashAP(conn, pkt)
	// case "AA":
	//     return packet.HandleHashAA(conn, pkt)
	default:
		logger.Debug("unhandled hash packet subtype",
			zap.String("subtype", pkt.SubType),
		)
		return nil
	}
}
