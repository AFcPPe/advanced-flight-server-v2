package handler

import (
	"context"
	"strings"
	"time"

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
		zap.String("realname", p.Name),
		zap.String("password", p.Password),
		zap.Int("rating", int(p.Rating)),
		zap.Int("protocol_version", p.ProtocolRevision),
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

	// 标记会话正在进行异步认证，防止被auth timeout踢掉
	mgr := session.GetManager()
	if sess := mgr.GetSession(conn); sess != nil {
		sess.Authenticating = true
	}

	// 将阻塞的DB认证移到goroutine中，避免阻塞gnet事件循环
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("PANIC in HandleAddATC goroutine",
					zap.Any("recover", r),
					zap.String("callsign", p.Callsign),
					zap.String("cid", p.Cid),
					zap.Stack("stack"),
				)
			}
		}()

		logger.Debug("ATC auth goroutine started",
			zap.String("callsign", p.Callsign),
			zap.String("cid", p.Cid),
			zap.String("remote", conn.RemoteAddr().String()),
		)

		sess := mgr.GetSession(conn)

		// 加超时保护，防止DB无限阻塞
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		authService := service.NewAuthService()
		result := authService.ValidateLoginWithRating(ctx, p.Cid, p.Password, p.Rating)
		logger.Debug("ATC auth DB query completed",
			zap.String("callsign", p.Callsign),
			zap.String("cid", p.Cid),
			zap.Bool("success", result.Success),
		)
		if !result.Success {
			logger.Warn("ATC login failed",
				zap.String("cid", p.Cid),
				zap.String("callsign", p.Callsign),
				zap.Error(result.Error),
			)
			if sess != nil {
				sess.FinishAuthenticating(false)
			}
			_ = session.SendErrorAndClose(conn, p.Callsign, pdu.NetworkErrorInvalidLogon, result.Error.Error())
			return
		}

		logger.Info("ATC login success",
			zap.String("cid", p.Cid),
			zap.String("callsign", p.Callsign),
			zap.Int("level", int(result.Level)),
		)

		// 检查连接是否仍然有效（可能在认证期间已断开）
		if mgr.GetSession(conn) == nil {
			logger.Warn("ATC connection already closed before callsign assignment",
				zap.String("callsign", p.Callsign),
				zap.String("cid", p.Cid),
			)
			return
		}

		// 原子地设置callsign，避免并发竞态条件
		success, callsignInUse := mgr.SetCallsignIfNotExist(conn, p.Callsign)
		logger.Debug("ATC SetCallsignIfNotExist result",
			zap.String("callsign", p.Callsign),
			zap.String("cid", p.Cid),
			zap.Bool("success", success),
			zap.Bool("callsignInUse", callsignInUse),
		)
		if !success {
			if callsignInUse {
				logger.Warn("ATC callsign already in use",
					zap.String("callsign", p.Callsign),
					zap.String("cid", p.Cid),
				)
				if sess != nil {
					sess.FinishAuthenticating(false)
				}
				_ = session.SendErrorAndClose(conn, p.Callsign, pdu.NetworkErrorCallsignInUse, "callsign already in use")
			} else {
				logger.Warn("ATC session not found",
					zap.String("callsign", p.Callsign),
					zap.String("cid", p.Cid),
				)
				if sess != nil {
					sess.FinishAuthenticating(false)
				}
				_ = session.SendErrorAndClose(conn, p.Callsign, pdu.NetworkErrorInvalidLogon, "session not found")
			}
			return
		}

		// 再次检查连接是否仍然有效
		if mgr.GetSession(conn) == nil {
			logger.Warn("ATC connection closed after callsign assignment, cleaning up",
				zap.String("callsign", p.Callsign),
				zap.String("cid", p.Cid),
			)
			return
		}

		// 设置连接类型为ATC，并保存CID和登录等级
		mgr.SetConnType(conn, session.ConnectionTypeATC)
		mgr.SetCid(conn, p.Cid)
		mgr.SetRating(conn, int(p.Rating))
		mgr.SetRealName(conn, p.Name)

		logger.Debug("ATC session fully configured",
			zap.String("callsign", p.Callsign),
			zap.String("cid", p.Cid),
		)

		// 发送 motd
		if motd := config.GetServer().Motd; motd != "" {
			for _, line := range strings.Split(motd, "\n") {
				if line = strings.TrimSpace(line); line != "" {
					_ = session.Send(conn, pdu.NewPDUTextMessage("SERVER", p.Callsign, line))
				}
			}
		}

		// 询问 ATC 的 ATIS 信息
		_ = session.Send(conn, pdu.NewPDUClientQuery("SERVER", p.Callsign, "ATIS", nil))

		// 认证完成，重放缓存的包
		sess = mgr.GetSession(conn) // 重新获取最新的session
		if sess != nil {
			logger.Debug("ATC finishing authentication, replaying pending packets",
				zap.String("callsign", p.Callsign),
				zap.String("cid", p.Cid),
				zap.Int("pending_count", len(sess.PendingPackets)),
			)
			sess.FinishAuthenticating(true)
		} else {
			logger.Warn("ATC session gone before FinishAuthenticating",
				zap.String("callsign", p.Callsign),
				zap.String("cid", p.Cid),
			)
		}
	}()

	return nil
}
