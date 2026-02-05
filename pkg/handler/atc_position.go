package handler

import (
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol/pdu"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// HandleATCPosition 处理ATC位置更新
func HandleATCPosition(conn gnet.Conn, p *pdu.ATCPosition) error {
	logger.Debug("handling ATCPosition",
		zap.String("callsign", p.From),
		zap.Strings("frequencies", p.Frequencies),
		zap.Float64("lat", p.Lat),
		zap.Float64("lon", p.Lon),
		zap.String("remote", conn.RemoteAddr().String()),
	)

	// TODO: 更新ATC位置信息
	// TODO: 广播给相关客户端

	return nil
}
