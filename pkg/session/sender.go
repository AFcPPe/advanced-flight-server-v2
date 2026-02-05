package session

import (
	"strings"

	"advanced-flight-server/pkg/util"

	"github.com/panjf2000/gnet/v2"
)

// Sendable PDU发送接口，避免循环引用pdu包
type Sendable interface {
	GetHeader() string
	ToTokens() []string
}

// Serialize 将Sendable序列化为可发送的字节数据（带\r\n结尾）
func Serialize(p Sendable) []byte {
	tokens := p.ToTokens()
	return []byte(p.GetHeader() + strings.Join(tokens, ":") + "\r\n")
}

// Send 向指定连接发送PDU（回复性发送）
func Send(conn gnet.Conn, p Sendable) error {
	return conn.AsyncWrite(Serialize(p), nil)
}

// BroadcastInRange 向指定连接范围内的所有用户广播PDU
// 规则：距离小于二者VisibilityRange之和
func BroadcastInRange(conn gnet.Conn, p Sendable) {
	data := Serialize(p)
	mgr := GetManager()

	// 获取发送者的session信息
	senderSession := mgr.GetSession(conn)
	if senderSession == nil {
		return
	}

	sessions := mgr.GetAllSessions()
	for _, s := range sessions {
		if s.Conn == conn {
			continue // 不发给自己
		}

		if !s.IsLoggedIn() {
			continue
		}

		// 计算距离
		distance := util.DistanceNM(senderSession.Lat, senderSession.Lon, s.Lat, s.Lon)

		// 距离小于二者VisibilityRange之和则发送
		if distance < float64(senderSession.VisibilityRange+s.VisibilityRange) {
			_ = s.Conn.AsyncWrite(data, nil)
		}
	}
}

// SendToCallsign 向指定callsign发送PDU
// 返回 true 表示目标存在并已发送，false 表示目标不存在
func SendToCallsign(callsign string, p Sendable) bool {
	mgr := GetManager()
	targetConn := mgr.GetConnByCallsign(callsign)
	if targetConn == nil {
		return false
	}
	_ = targetConn.AsyncWrite(Serialize(p), nil)
	return true
}

// BroadcastToATC 向所有ATC广播PDU
func BroadcastToATC(p Sendable) {
	data := Serialize(p)
	mgr := GetManager()

	sessions := mgr.GetAllSessions()
	for _, s := range sessions {
		if !s.IsLoggedIn() {
			continue
		}

		if s.ConnType == ConnectionTypeATC {
			_ = s.Conn.AsyncWrite(data, nil)
		}
	}
}

// BroadcastToAll 向所有已登录用户广播PDU
func BroadcastToAll(p Sendable) {
	data := Serialize(p)
	mgr := GetManager()

	sessions := mgr.GetAllSessions()
	for _, s := range sessions {
		if !s.IsLoggedIn() {
			continue
		}
		_ = s.Conn.AsyncWrite(data, nil)
	}
}
