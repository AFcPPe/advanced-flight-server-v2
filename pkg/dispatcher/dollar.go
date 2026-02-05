package dispatcher

import (
	"advanced-flight-server/pkg/handler"
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol"
	"advanced-flight-server/pkg/protocol/pdu"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// DispatchDollar 分发 $ 开头的包
func DispatchDollar(conn gnet.Conn, pkt *protocol.Packet) error {
	switch pkt.SubType {
	case "ID":
		return handler.HandleClientIdentification(conn, pkt)
	case "FP":
		return dispatchFlightPlan(conn, pkt)
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

// dispatchFlightPlan 解析并分发飞行计划包
func dispatchFlightPlan(conn gnet.Conn, pkt *protocol.Packet) error {
	p, err := pdu.FlightPlanFromTokens(pkt.Tokens)
	if err != nil {
		logger.Error("failed to parse FlightPlan", zap.Error(err))
		return err
	}
	return handler.HandleFlightPlan(conn, p)
}
