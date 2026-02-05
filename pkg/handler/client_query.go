package handler

import (
	"strings"

	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol/pdu"
	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// ServerCaps 服务器支持的能力
var ServerCaps = []string{"ATCINFO=1", "ICAOEQ=1", "FASTPOS=1"}

// HandleClientQuery 处理客户端查询包（$CQ）
func HandleClientQuery(conn gnet.Conn, p *pdu.ClientQuery) error {
	logger.Debug("handling ClientQuery",
		zap.String("from", p.From),
		zap.String("to", p.To),
		zap.String("type", p.Type),
		zap.Strings("payload", p.Payload),
		zap.String("remote", conn.RemoteAddr().String()),
	)

	// 验证登录状态和callsign
	sess, err := ValidateLoginAndCallsign(conn, p.From)
	if err != nil {
		return err
	}

	to := strings.ToLower(p.To)

	// 特殊处理：to 是 server
	if to == "server" {
		return handleServerQuery(conn, sess, p)
	}

	// // 特殊处理：to 是 @94835，转发给所有ATC
	// if to == "@94835" {
	// 	session.BroadcastToATC(p)
	// 	return nil
	// }

	// 特殊处理：to 以 @ 开头，转发给所有已登录用户
	if strings.HasPrefix(to, "@") {
		session.BroadcastToAll(p)
		return nil
	}

	// 默认：转发给目标
	session.SendToCallsign(p.To, p)

	return nil
}

// handleServerQuery 处理对服务器的查询
func handleServerQuery(conn gnet.Conn, sess *session.Session, p *pdu.ClientQuery) error {
	switch p.Type {
	case "ATC":
		return handleATCQuery(conn, sess, p)
	case "CAPS":
		return handleCAPSQuery(conn, p)
	case "IP":
		return handleIPQuery(conn, p)
	case "FP":
		return handleFPQuery(conn, sess, p)
	default:
		// 未知类型，忽略或记录日志
		logger.Debug("unknown server query type",
			zap.String("type", p.Type),
			zap.String("from", p.From),
		)
	}
	return nil
}

// handleATCQuery 处理 ATC 查询
// 检查连接类型是否为ATC，等级不是OBS，设施不是Observer
func handleATCQuery(conn gnet.Conn, sess *session.Session, p *pdu.ClientQuery) error {
	// 检查是否是ATC连接
	if sess.ConnType != session.ConnectionTypeATC {
		return nil
	}

	// 检查等级不是OBS (RatingOBS = 1)
	if sess.Rating == int(pdu.RatingOBS) {
		return nil
	}

	// 检查设施不是Observer (Observer = 0)
	if sess.Facility == int(pdu.Observer) {
		return nil
	}

	// 回复：From SERVER, type ATC, Payload: ["Y", callsign]
	response := pdu.NewPDUClientQueryResponse("SERVER", p.From, "ATC", []string{"Y", p.From})
	return session.Send(conn, response)
}

// handleCAPSQuery 处理 CAPS 查询
func handleCAPSQuery(conn gnet.Conn, p *pdu.ClientQuery) error {
	// 回复：From SERVER, type CAPS, Payload: ServerCaps
	response := pdu.NewPDUClientQueryResponse("SERVER", p.From, "CAPS", ServerCaps)
	return session.Send(conn, response)
}

// handleIPQuery 处理 IP 查询
func handleIPQuery(conn gnet.Conn, p *pdu.ClientQuery) error {
	// 回复：From SERVER, type IP, Payload: ["unknown"]
	response := pdu.NewPDUClientQueryResponse("SERVER", p.From, "IP", []string{"unknown"})
	return session.Send(conn, response)
}

// handleFPQuery 处理 FP 查询（查询指定用户的飞行计划）
func handleFPQuery(conn gnet.Conn, sess *session.Session, p *pdu.ClientQuery) error {
	// payload[0] 是要查询的目标用户 callsign
	if len(p.Payload) == 0 {
		return nil
	}

	targetCallsign := p.Payload[0]

	// 查找目标用户的会话
	targetSess := session.GetManager().GetSessionByCallsign(targetCallsign)
	if targetSess == nil {
		return nil
	}

	// 检查目标用户是否有飞行计划
	if targetSess.FlightPlan == nil {
		return nil
	}

	fp := targetSess.FlightPlan

	// 组装 FlightPlan PDU 回复
	flightPlanPDU := &pdu.FlightPlan{
		Base: pdu.Base{
			From: "SERVER",
			To:   p.From,
		},
		Callsign:      targetCallsign,
		FlightRules:   fp.FlightRules,
		Type:          fp.Type,
		TAS:           fp.TAS,
		Dep:           fp.Dep,
		DepTime:       fp.DepTime,
		ActualDepTime: fp.ActualDepTime,
		CruiseAlt:     fp.CruiseAlt,
		Dest:          fp.Dest,
		EnrouteHour:   fp.EnrouteHour,
		EnrouteMin:    fp.EnrouteMin,
		FobHour:       fp.FobHour,
		FobMin:        fp.FobMin,
		AlterDest:     fp.AlterDest,
		Remark:        fp.Remark,
		Route:         fp.Route,
	}

	return session.Send(conn, flightPlanPDU)
}
