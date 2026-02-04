package session

import (
	"github.com/panjf2000/gnet/v2"
)

// GetSessionByConn 通过连接获取会话
func (m *Manager) GetSessionByConn(conn gnet.Conn) *Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.connToSession[conn]
}

// GetConnByCallsign 通过callsign获取连接
func (m *Manager) GetConnByCallsign(callsign string) gnet.Conn {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.callsignToConn[callsign]
}

// GetCallsignByConn 通过连接获取callsign
func (m *Manager) GetCallsignByConn(conn gnet.Conn) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.connToSession[conn]
	if !exists {
		return ""
	}
	return session.Callsign
}

// GetSessionByCallsign 通过callsign获取会话
func (m *Manager) GetSessionByCallsign(callsign string) *Session {
	m.mu.RLock()
	defer m.mu.RUnlock()

	conn, exists := m.callsignToConn[callsign]
	if !exists {
		return nil
	}
	return m.connToSession[conn]
}

// GetAllSessions 获取所有会话（用于广播等场景）
func (m *Manager) GetAllSessions() []*Session {
	m.mu.RLock()
	defer m.mu.RUnlock()

	sessions := make([]*Session, 0, len(m.connToSession))
	for _, session := range m.connToSession {
		sessions = append(sessions, session)
	}
	return sessions
}

// Count 返回当前连接数
func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.connToSession)
}
