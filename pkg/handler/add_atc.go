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

// HandleAddATC 处理管制员登录
func HandleAddATC(conn gnet.Conn, p *pdu.AddATC) error {
	logger.Debug("handling AddATC",
		zap.String("callsign", p.Callsign),
		zap.String("cid", p.Cid),
		zap.String("realname", p.RealName),
		zap.String("password", p.Password),
		zap.Int("rating", int(p.Rating)),
		zap.String("rating_str", p.Rating.String()),
		zap.Int("protocol_version", p.ProtocolVersion),
		zap.String("remote", conn.RemoteAddr().String()),
	)

	// 验证callsign长度不能超过12
	if len(p.Callsign) > 12 {
		logger.Warn("ATC callsign too long",
			zap.String("callsign", p.Callsign),
			zap.Int("length", len(p.Callsign)),
		)
		return session.SendErrorAndClose(conn, p.Callsign, pdu.NetworkErrorCallsignInvalid, "callsign too long (max 12 characters)")
	}

	// 认证验证（密码 + rating）
	authService := service.NewAuthService()
	result := authService.ValidateLoginWithRating(context.Background(), p.Cid, p.Password, p.Rating)
	if !result.Success {
		logger.Warn("ATC login failed",
			zap.String("cid", p.Cid),
			zap.String("callsign", p.Callsign),
			zap.Error(result.Error),
		)
		return session.SendErrorAndClose(conn, p.Callsign, pdu.NetworkErrorInvalidLogon, result.Error.Error())
	}

	logger.Info("ATC login success",
		zap.String("cid", p.Cid),
		zap.String("callsign", p.Callsign),
		zap.Int("level", int(result.Level)),
	)

	// 原子地设置callsign，避免并发竞态条件
	mgr := session.GetManager()
	success, callsignInUse := mgr.SetCallsignIfNotExist(conn, p.Callsign)
	if !success {
		if callsignInUse {
			logger.Warn("ATC callsign already in use",
				zap.String("callsign", p.Callsign),
				zap.String("cid", p.Cid),
			)
			return session.SendErrorAndClose(conn, p.Callsign, pdu.NetworkErrorCallsignInUse, "callsign already in use")
		}
		logger.Warn("ATC session not found",
			zap.String("callsign", p.Callsign),
			zap.String("cid", p.Cid),
		)
		return session.SendErrorAndClose(conn, p.Callsign, pdu.NetworkErrorInvalidLogon, "session not found")
	}

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
