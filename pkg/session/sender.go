package session

import (
	"strings"

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

// SendTo 向指定callsign发送PDU（针对性发送）
func SendTo(callsign string, p Sendable) error {
	mgr := GetManager()
	conn := mgr.GetConnByCallsign(callsign)
	if conn == nil {
		return nil // 目标不在线，静默忽略
	}
	return conn.AsyncWrite(Serialize(p), nil)
}

// SendToMultiple 向多个callsign发送PDU
func SendToMultiple(callsigns []string, p Sendable) {
	data := Serialize(p)
	mgr := GetManager()
	for _, callsign := range callsigns {
		conn := mgr.GetConnByCallsign(callsign)
		if conn != nil {
			_ = conn.AsyncWrite(data, nil)
		}
	}
}

// Broadcast 向所有连接广播PDU（广播性发送）
func Broadcast(p Sendable) {
	data := Serialize(p)
	mgr := GetManager()
	sessions := mgr.GetAllSessions()
	for _, s := range sessions {
		_ = s.Conn.AsyncWrite(data, nil)
	}
}

// BroadcastExcept 向除指定连接外的所有连接广播PDU
func BroadcastExcept(except gnet.Conn, p Sendable) {
	data := Serialize(p)
	mgr := GetManager()
	sessions := mgr.GetAllSessions()
	for _, s := range sessions {
		if s.Conn != except {
			_ = s.Conn.AsyncWrite(data, nil)
		}
	}
}

// BroadcastExceptCallsign 向除指定callsign外的所有连接广播PDU
func BroadcastExceptCallsign(exceptCallsign string, p Sendable) {
	data := Serialize(p)
	mgr := GetManager()
	sessions := mgr.GetAllSessions()
	for _, s := range sessions {
		if s.Callsign != exceptCallsign {
			_ = s.Conn.AsyncWrite(data, nil)
		}
	}
}
