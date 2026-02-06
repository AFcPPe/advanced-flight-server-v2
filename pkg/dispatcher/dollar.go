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
	case "CQ":
		return dispatchClientQuery(conn, pkt)
	case "CR":
		return dispatchClientQueryResponse(conn, pkt)
	case "HO":
		return dispatchHandOff(conn, pkt)
	case "HA":
		return dispatchHandOffAccept(conn, pkt)
	case "AM":
		return dispatchTagModify(conn, pkt)
	case "AX":
		return dispatchMetarRequest(conn, pkt)
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

// dispatchClientQuery 解析并分发客户端查询包
func dispatchClientQuery(conn gnet.Conn, pkt *protocol.Packet) error {
	p, err := pdu.ClientQueryFromTokens(pkt.Tokens)
	if err != nil {
		logger.Error("failed to parse ClientQuery", zap.Error(err))
		return err
	}
	return handler.HandleClientQuery(conn, p)
}

// dispatchClientQueryResponse 解析并分发客户端查询响应包
func dispatchClientQueryResponse(conn gnet.Conn, pkt *protocol.Packet) error {
	p, err := pdu.ClientQueryResponseFromTokens(pkt.Tokens)
	if err != nil {
		logger.Error("failed to parse ClientQueryResponse", zap.Error(err))
		return err
	}
	return handler.HandleClientQueryResponse(conn, p)
}

// dispatchHandOff 解析并分发移交包
func dispatchHandOff(conn gnet.Conn, pkt *protocol.Packet) error {
	p, err := pdu.HandOffFromTokens(pkt.Tokens)
	if err != nil {
		logger.Error("failed to parse HandOff", zap.Error(err))
		return err
	}
	return handler.HandleHandOff(conn, p)
}

// dispatchHandOffAccept 解析并分发移交接受包
func dispatchHandOffAccept(conn gnet.Conn, pkt *protocol.Packet) error {
	p, err := pdu.HandOffAcceptFromTokens(pkt.Tokens)
	if err != nil {
		logger.Error("failed to parse HandOffAccept", zap.Error(err))
		return err
	}
	return handler.HandleHandOffAccept(conn, p)
}

// dispatchTagModify 解析并分发ATC修改飞行计划包
func dispatchTagModify(conn gnet.Conn, pkt *protocol.Packet) error {
	p, err := pdu.TagModifyFromTokens(pkt.Tokens)
	if err != nil {
		logger.Error("failed to parse TagModify", zap.Error(err))
		return err
	}
	return handler.HandleTagModify(conn, p)
}

// dispatchMetarRequest 解析并分发METAR请求包
func dispatchMetarRequest(conn gnet.Conn, pkt *protocol.Packet) error {
	p, err := pdu.MetarRequestFromTokens(pkt.Tokens)
	if err != nil {
		logger.Error("failed to parse MetarRequest", zap.Error(err))
		return err
	}
	return handler.HandleMetarRequest(conn, p)
}
