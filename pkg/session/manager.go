package session

import (
	"sync"

	"github.com/panjf2000/gnet/v2"
)

// Manager 管理所有连接的会话，支持双向查找
type Manager struct {
	mu             sync.RWMutex
	connToSession  map[gnet.Conn]*Session // 通过连接找会话
	callsignToConn map[string]gnet.Conn   // 通过callsign找连接
}

var (
	instance *Manager
	once     sync.Once
)

// GetManager 获取全局会话管理器单例
func GetManager() *Manager {
	once.Do(func() {
		instance = &Manager{
			connToSession:  make(map[gnet.Conn]*Session),
			callsignToConn: make(map[string]gnet.Conn),
		}
	})
	return instance
}

// AddConn 添加新连接（此时还没有callsign）
func (m *Manager) AddConn(conn gnet.Conn) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	session := &Session{
		Conn: conn,
	}
	m.connToSession[conn] = session
	return session
}

// SetCallsign 为连接设置callsign
func (m *Manager) SetCallsign(conn gnet.Conn, callsign string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.connToSession[conn]
	if !exists {
		return false
	}

	// 如果之前有callsign，先移除旧的映射
	if session.Callsign != "" {
		delete(m.callsignToConn, session.Callsign)
	}

	session.Callsign = callsign
	m.callsignToConn[callsign] = conn
	return true
}

// RemoveConn 移除连接
func (m *Manager) RemoveConn(conn gnet.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.connToSession[conn]
	if !exists {
		return
	}

	if session.Callsign != "" {
		delete(m.callsignToConn, session.Callsign)
	}
	delete(m.connToSession, conn)
}
