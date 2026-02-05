package session

import (
	"github.com/panjf2000/gnet/v2"
)

// ConnectionType 连接类型
type ConnectionType int

const (
	ConnectionTypeUnknown ConnectionType = iota
	ConnectionTypePilot
	ConnectionTypeATC
)

const (
	// MaxBufferSize 每个连接的最大缓存大小 (500KB)
	MaxBufferSize = 500 * 1024
)

// Session 表示一个客户端连接的会话信息
type Session struct {
	Conn          gnet.Conn
	Callsign      string
	Buffer        []byte         // 用于处理粘包的缓存
	Authenticated bool           // 是否已通过$ID验证（首包必须是$ID）
	ConnType      ConnectionType // 连接类型：Pilot或ATC

	// 位置信息（Pilot和ATC共用）
	Lat             float64
	Lon             float64
	VisibilityRange int // Pilot固定为40，ATC由客户端设置

	// Pilot特有字段
	SquawkCode       int
	SquawkingModeC   bool
	Identing         bool
	TrueAltitude     int
	PressureAltitude int
	GroundSpeed      int
	Pitch            float64
	Heading          float64
	Bank             float64

	// ATC特有字段
	Frequencies []string
	Facility    int
}

// AppendBuffer 追加数据到缓存，如果超过限制则丢弃全部缓存
// 返回 true 表示追加成功，false 表示缓存已满并被清空
func (s *Session) AppendBuffer(data []byte) bool {
	if len(s.Buffer)+len(data) > MaxBufferSize {
		s.Buffer = nil
		return false
	}
	s.Buffer = append(s.Buffer, data...)
	return true
}

// ClearBuffer 清空缓存
func (s *Session) ClearBuffer() {
	s.Buffer = nil
}

// SetBuffer 设置缓存内容
func (s *Session) SetBuffer(data []byte) {
	s.Buffer = data
}

// IsLoggedIn 判断用户是否已认证且已登录
func (s *Session) IsLoggedIn() bool {
	return s.Authenticated && s.Callsign != ""
}
