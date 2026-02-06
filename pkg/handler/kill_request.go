package handler

import (
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol/pdu"
	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// HandleKillRequest 处理踢人请求
// 仅SUP及以上等级可以踢人，且不能踢登录等级大于或等于自己的人
func HandleKillRequest(conn gnet.Conn, p *pdu.KillRequest) error {
	logger.Debug("handling KillRequest",
		zap.String("from", p.From),
		zap.String("target", p.To),
		zap.String("reason", p.Reason),
		zap.String("remote", conn.RemoteAddr().String()),
	)

	// 验证登录状态和callsign
	sess, err := ValidateLoginAndCallsign(conn, p.From)
	if err != nil {
		return err
	}

	// 验证发送者等级必须是SUP及以上
	if sess.Rating < int(pdu.RatingSUP) {
		logger.Warn("KillRequest from non-supervisor",
			zap.String("callsign", p.From),
			zap.Int("rating", sess.Rating),
		)
		return session.SendErrorAndClose(conn, p.From, pdu.NetworkErrorInvalidControl, "insufficient rating to kill")
	}

	// 查找目标会话
	mgr := session.GetManager()
	targetSess := mgr.GetSessionByCallsign(p.To)
	if targetSess == nil {
		logger.Warn("KillRequest target not found",
			zap.String("from", p.From),
			zap.String("target", p.To),
		)
		return nil
	}

	// 不能踢登录等级大于或等于自己的人
	if targetSess.Rating >= sess.Rating {
		logger.Warn("KillRequest target rating too high",
			zap.String("from", p.From),
			zap.Int("from_rating", sess.Rating),
			zap.String("target", p.To),
			zap.Int("target_rating", targetSess.Rating),
		)
		_ = session.Send(conn, pdu.NewPDUTextMessage("SERVER", p.From, "cannot kill user with equal or higher rating"))
		return nil
	}

	logger.Info("killing user",
		zap.String("by", p.From),
		zap.String("target", p.To),
		zap.String("reason", p.Reason),
	)

	// 通知目标被踢原因
	_ = session.Send(targetSess.Conn, pdu.NewPDUTextMessage("SERVER", p.To, "you have been kicked by "+p.From+": "+p.Reason))

	// 向目标发送踢人PDU，然后断开目标连接
	_ = session.Send(targetSess.Conn, p)

	// 根据目标连接类型广播对应的Delete包，通知其他客户端
	if targetSess.ConnType == session.ConnectionTypePilot {
		deletePdu := pdu.NewPDUDeletePilot(p.To, targetSess.Cid)
		session.BroadcastToAll(targetSess.Conn, deletePdu)
	} else if targetSess.ConnType == session.ConnectionTypeATC {
		deletePdu := pdu.NewPDUDeleteATC(p.To, targetSess.Cid)
		session.BroadcastToAll(targetSess.Conn, deletePdu)
	}

	// 关闭目标连接
	return targetSess.Conn.Close()
}
