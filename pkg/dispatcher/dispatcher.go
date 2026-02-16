package dispatcher

import (
	"errors"

	"advanced-flight-server/pkg/errs"
	"advanced-flight-server/pkg/handler"
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol"
	"advanced-flight-server/pkg/protocol/pdu"
	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// Dispatch 根据包类型分发到对应的处理函数
func Dispatch(conn gnet.Conn, pkt *protocol.Packet) error {
	logger.Debug("dispatching packet",
		zap.String("type", pkt.GetTypeName()),
		zap.String("remote", conn.RemoteAddr().String()),
	)

	// 获取会话，检查是否已认证
	sess := session.GetManager().GetSessionByConn(conn)
	if sess == nil {
		return errors.New("session not found")
	}

	// 连接正在关闭中，静默丢弃后续包
	if sess.Closing {
		return nil
	}

	// 正在进行异步登录认证，缓存后续包等认证完成后重放
	if sess.Authenticating {
		sess.PendingPackets = append(sess.PendingPackets, &session.PendingPacket{
			RawData: append([]byte(nil), pkt.RawData...),
		})
		// 设置重放回调（仅首次设置）
		if sess.ReplayFunc == nil {
			sess.ReplayFunc = func() {
				ReplayPendingPackets(conn, sess.PendingPackets)
				sess.PendingPackets = nil
			}
		}
		logger.Debug("buffered packet during authentication",
			zap.String("type", pkt.GetTypeName()),
			zap.String("remote", conn.RemoteAddr().String()),
			zap.Int("pending_count", len(sess.PendingPackets)),
		)
		return nil
	}

	// 如果未认证，首包必须是$ID
	if !sess.Authenticated {
		if pkt.Type != protocol.PacketTypeDollar || pkt.SubType != "ID" {
			logger.Warn("first packet is not $ID, closing connection",
				zap.String("remote", conn.RemoteAddr().String()),
				zap.String("type", pkt.GetTypeName()),
			)
			return errs.ErrNotAuthenticated
		}
	}

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

// ReplayPendingPackets 重放认证期间缓存的包
// 在异步认证完成后由goroutine调用
func ReplayPendingPackets(conn gnet.Conn, packets []*session.PendingPacket) {
	logger.Debug("ReplayPendingPackets starting",
		zap.String("remote", conn.RemoteAddr().String()),
		zap.Int("packet_count", len(packets)),
	)
	for i, pending := range packets {
		pkt := protocol.ParsePacket(pending.RawData)
		logger.Debug("replaying pending packet",
			zap.Int("index", i),
			zap.String("type", pkt.GetTypeName()),
			zap.String("remote", conn.RemoteAddr().String()),
		)
		if err := Dispatch(conn, pkt); err != nil {
			logger.Error("failed to replay pending packet",
				zap.Error(err),
				zap.String("type", pkt.GetTypeName()),
				zap.String("remote", conn.RemoteAddr().String()),
			)
		}
	}
	logger.Debug("ReplayPendingPackets completed",
		zap.String("remote", conn.RemoteAddr().String()),
	)
}

// dispatchATCPosition 解析并分发ATC位置包
func dispatchATCPosition(conn gnet.Conn, pkt *protocol.Packet) error {
	p, err := pdu.ATCPositionFromTokens(pkt.Tokens)
	if err != nil {
		logger.Error("failed to parse ATCPosition", zap.Error(err))
		return err
	}
	return handler.HandleATCPosition(conn, p)
}

// dispatchPilotPosition 解析并分发飞行员位置包
func dispatchPilotPosition(conn gnet.Conn, pkt *protocol.Packet) error {
	p, err := pdu.PilotPositionFromTokens(pkt.Tokens)
	if err != nil {
		logger.Error("failed to parse PilotPosition", zap.Error(err))
		return err
	}
	return handler.HandlePilotPosition(conn, p)
}
