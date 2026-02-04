package session

import (
	"github.com/panjf2000/gnet/v2"
)

// Session 表示一个客户端连接的会话信息
type Session struct {
	Conn     gnet.Conn
	Callsign string
}
