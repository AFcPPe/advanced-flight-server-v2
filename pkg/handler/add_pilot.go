package handler

import (
	"context"
	"strings"

	"advanced-flight-server/pkg/config"
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol/pdu"
	"advanced-flight-server/pkg/service"
	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// HandleAddPilot 处理飞行员登录
func HandleAddPilot(conn gnet.Conn, p *pdu.AddPilot) error {
	logger.Debug("handling AddPilot",
		zap.String("callsign", p.Callsign),
		zap.String("cid", p.Cid),
		zap.String("realname", p.RealName),
		zap.String("password", p.Password),
		zap.Int("rating", int(p.Rating)),
		zap.String("rating_str", p.Rating.String()),
		zap.Int("protocol_version", p.ProtocolVersion),
		zap.Int("simtype", p.SimType),
		zap.String("remote", conn.RemoteAddr().String()),
	)

	// 验证callsign长度不能超过12
	if len(p.Callsign) > 12 {
		logger.Warn("Pilot callsign too long",
			zap.String("callsign", p.Callsign),
			zap.Int("length", len(p.Callsign)),
		)
		return session.SendErrorAndClose(conn, p.Callsign, pdu.NetworkErrorCallsignInvalid, "callsign too long (max 12 characters)")
	}

	// 认证验证（密码 + rating）
	authService := service.NewAuthService()
	result := authService.ValidateLoginWithRating(context.Background(), p.Cid, p.Password, p.Rating)
	if !result.Success {
		logger.Warn("Pilot login failed",
			zap.String("cid", p.Cid),
			zap.String("callsign", p.Callsign),
			zap.Error(result.Error),
		)
		return session.SendErrorAndClose(conn, p.Callsign, pdu.NetworkErrorInvalidLogon, result.Error.Error())
	}

	logger.Info("Pilot login success",
		zap.String("cid", p.Cid),
		zap.String("callsign", p.Callsign),
		zap.Int("level", int(result.Level)),
	)

	// 创建session并设置callsign（在认证成功后检查重复，避免并发问题）
	mgr := session.GetManager()
	if existingConn := mgr.GetConnByCallsign(p.Callsign); existingConn != nil {
		logger.Warn("Pilot callsign already in use",
			zap.String("callsign", p.Callsign),
			zap.String("cid", p.Cid),
		)
		return session.SendErrorAndClose(conn, p.Callsign, pdu.NetworkErrorCallsignInUse, "callsign already in use")
	}
	mgr.SetCallsign(conn, p.Callsign)

	// 发送 motd
	if motd := config.GetServer().Motd; motd != "" {
		for _, line := range strings.Split(motd, "\n") {
			if line = strings.TrimSpace(line); line != "" {
				_ = session.Send(conn, pdu.NewServerTextMessage(p.Callsign, line))
			}
		}
	}

	return nil
}
