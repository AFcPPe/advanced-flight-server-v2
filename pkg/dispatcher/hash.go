package dispatcher

import (
	"advanced-flight-server/pkg/handler"
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol"
	"advanced-flight-server/pkg/protocol/pdu"
	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// closeOnParseError 解析失败时关闭连接，并标记会话为正在关闭，
// 使同一批粘包里后续的包被 Dispatch 丢弃，避免在关闭中的连接上继续处理。
func closeOnParseError(conn gnet.Conn) {
	if sess := session.GetManager().GetSession(conn); sess != nil {
		sess.Closing = true
	}
	_ = conn.Close()
}

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
	case "DP":
		return dispatchDeletePilot(conn, pkt)
	case "DA":
		return dispatchDeleteATC(conn, pkt)
	case "TM":
		return dispatchTextMessage(conn, pkt)
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
		closeOnParseError(conn)
		return err
	}
	return handler.HandleAddPilot(conn, p)
}

// dispatchAddATC 解析并分发管制员登录包
func dispatchAddATC(conn gnet.Conn, pkt *protocol.Packet) error {
	p, err := pdu.AddATCFromTokens(pkt.Tokens)
	if err != nil {
		logger.Error("failed to parse AddATC", zap.Error(err))
		closeOnParseError(conn)
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

// dispatchDeletePilot 解析并分发飞行员退出包
func dispatchDeletePilot(conn gnet.Conn, pkt *protocol.Packet) error {
	p, err := pdu.DeletePilotFromTokens(pkt.Tokens)
	if err != nil {
		logger.Error("failed to parse DeletePilot", zap.Error(err))
		return err
	}
	return handler.HandleDeletePilot(conn, p)
}

// dispatchDeleteATC 解析并分发管制员退出包
func dispatchDeleteATC(conn gnet.Conn, pkt *protocol.Packet) error {
	p, err := pdu.DeleteATCFromTokens(pkt.Tokens)
	if err != nil {
		logger.Error("failed to parse DeleteATC", zap.Error(err))
		return err
	}
	return handler.HandleDeleteATC(conn, p)
}

// dispatchTextMessage 解析并分发文本消息包
func dispatchTextMessage(conn gnet.Conn, pkt *protocol.Packet) error {
	p, err := pdu.TextMessageFromTokens(pkt.Tokens)
	if err != nil {
		logger.Error("failed to parse TextMessage", zap.Error(err))
		return err
	}
	return handler.HandleTextMessage(conn, p)
}
