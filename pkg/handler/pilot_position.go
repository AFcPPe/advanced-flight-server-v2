package handler

import (
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol/pdu"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// HandlePilotPosition 处理飞行员位置更新
func HandlePilotPosition(conn gnet.Conn, p *pdu.PilotPosition) error {
	logger.Debug("handling PilotPosition",
		zap.String("callsign", p.From),
		zap.Int("squawk", p.SquawkCode),
		zap.Float64("lat", p.Lat),
		zap.Float64("lon", p.Lon),
		zap.Int("altitude", p.TrueAltitude),
		zap.String("remote", conn.RemoteAddr().String()),
	)

	// TODO: 更新飞行员位置信息
	// TODO: 广播给相关客户端

	return nil
}
