package session

import (
	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/protocol/pdu"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

// SendErrorAndClose 发送错误消息给客户端并断开连接
// 确保消息发送后再关闭连接
func SendErrorAndClose(conn gnet.Conn, to string, errorCode pdu.NetworkError, message string) error {
	errPdu := pdu.NewServerError(to, errorCode, "", message)
	data := pdu.Serialize(errPdu)

	logger.Debug("sending error and closing connection",
		zap.String("to", to),
		zap.Int("errorCode", int(errorCode)),
		zap.String("message", message),
		zap.String("remote", conn.RemoteAddr().String()),
	)

	// 异步写入，写入完成后关闭连接
	return conn.AsyncWrite(data, func(c gnet.Conn, err error) error {
		if err != nil {
			logger.Error("failed to send error message",
				zap.String("to", to),
				zap.Error(err),
			)
		}
		// 无论发送成功与否，都关闭连接
		return c.Close()
	})
}
