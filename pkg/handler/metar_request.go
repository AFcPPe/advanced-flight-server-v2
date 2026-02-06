package handler

import (
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/metar"
	"advanced-flight-server/pkg/protocol/pdu"
	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// HandleMetarRequest 处理METAR请求包（$AX）
func HandleMetarRequest(conn gnet.Conn, p *pdu.MetarRequest) error {
	logger.Debug("handling MetarRequest",
		zap.String("from", p.From),
		zap.String("to", p.To),
		zap.String("icao", p.ICAO),
		zap.String("remote", conn.RemoteAddr().String()),
	)

	// 验证登录状态
	_, err := ValidateLogin(conn, p.From)
	if err != nil {
		return err
	}

	// 从METAR存储中查找
	metarData, ok := metar.GetStore().Get(p.ICAO)
	if !ok {
		metarData = ""
	}

	// 回复METAR响应
	response := pdu.NewPDUMetarResponse("SERVER", p.From, p.Type, metarData)
	return session.Send(conn, response)
}
