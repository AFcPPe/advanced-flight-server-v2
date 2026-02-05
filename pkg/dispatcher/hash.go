package dispatcher

import (
	"advanced-flight-server/pkg/handler"
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol"
	"advanced-flight-server/pkg/protocol/pdu"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// DispatchHash 分发 # 开头的包
func DispatchHash(conn gnet.Conn, pkt *protocol.Packet) error {
	switch pkt.SubType {
	case "AP":
		return dispatchAddPilot(conn, pkt)
	case "AA":
		return dispatchAddATC(conn, pkt)
	case "PC":
		return dispatchControllerInfo(conn, pkt)
	case "SB":
		return dispatchPlaneInfo(conn, pkt)
	default:
		logger.Debug("unhandled hash packet subtype",
			zap.String("subtype", pkt.SubType),
		)
		return nil
	}
}

// dispatchAddPilot 解析并分发飞行员登录包
func dispatchAddPilot(conn gnet.Conn, pkt *protocol.Packet) error {
	p, err := pdu.AddPilotFromTokens(pkt.Tokens)
	if err != nil {
		logger.Error("failed to parse AddPilot", zap.Error(err))
		return err
	}
	return handler.HandleAddPilot(conn, p)
}

// dispatchAddATC 解析并分发管制员登录包
func dispatchAddATC(conn gnet.Conn, pkt *protocol.Packet) error {
	p, err := pdu.AddATCFromTokens(pkt.Tokens)
	if err != nil {
		logger.Error("failed to parse AddATC", zap.Error(err))
		return err
	}
	return handler.HandleAddATC(conn, p)
}

// dispatchControllerInfo 解析并分发管制员信息包
func dispatchControllerInfo(conn gnet.Conn, pkt *protocol.Packet) error {
	p, err := pdu.ControllerInfoFromTokens(pkt.Tokens)
	if err != nil {
		logger.Error("failed to parse ControllerInfo", zap.Error(err))
		return err
	}
	return handler.HandleControllerInfo(conn, p)
}

// dispatchPlaneInfo 解析并分发飞机信息包
func dispatchPlaneInfo(conn gnet.Conn, pkt *protocol.Packet) error {
	p, err := pdu.PlaneInfoFromTokens(pkt.Tokens)
	if err != nil {
		logger.Error("failed to parse PlaneInfo", zap.Error(err))
		return err
	}
	return handler.HandlePlaneInfo(conn, p)
}
