package session

import (
	"sync"
	"time"

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

	now := time.Now()
	session := &Session{
		Conn:         conn,
		LastActivity: now,
		LogonTime:    now,
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

// SetCallsignIfNotExist 原子地为连接设置callsign，如果callsign已被占用则返回false
// 返回值: (success, callsignInUse)
//   - success=true, callsignInUse=false: 设置成功
//   - success=false, callsignInUse=true: callsign已被其他连接占用
//   - success=false, callsignInUse=false: 连接不存在
func (m *Manager) SetCallsignIfNotExist(conn gnet.Conn, callsign string) (success bool, callsignInUse bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.connToSession[conn]
	if !exists {
		return false, false
	}

	// 检查callsign是否已被其他连接占用
	if existingConn, occupied := m.callsignToConn[callsign]; occupied && existingConn != conn {
		// 检查占用callsign的连接是否还存在于connToSession中
		// 如果不存在，说明是残留的脏数据，清理后允许使用
		if _, stale := m.connToSession[existingConn]; !stale {
			delete(m.callsignToConn, callsign)
		} else {
			return false, true
		}
	}

	// 如果之前有callsign，先移除旧的映射
	if session.Callsign != "" {
		delete(m.callsignToConn, session.Callsign)
	}

	session.Callsign = callsign
	m.callsignToConn[callsign] = conn
	return true, false
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

// SetConnType 设置连接类型
func (m *Manager) SetConnType(conn gnet.Conn, connType ConnectionType) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, exists := m.connToSession[conn]; exists {
		session.ConnType = connType
	}
}

// GetSession 获取连接对应的Session
func (m *Manager) GetSession(conn gnet.Conn) *Session {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.connToSession[conn]
}

// GetConnType 获取连接类型
func (m *Manager) GetConnType(conn gnet.Conn) ConnectionType {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if session, exists := m.connToSession[conn]; exists {
		return session.ConnType
	}
	return ConnectionTypeUnknown
}

// UpdatePilotPosition 更新Pilot位置信息
func (m *Manager) UpdatePilotPosition(conn gnet.Conn, lat, lon float64, squawkCode int, squawkingModeC, identing bool, trueAlt, pressureAlt, gs int, pitch, heading, bank float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, exists := m.connToSession[conn]; exists {
		session.Lat = lat
		session.Lon = lon
		session.VisibilityRange = 40 // Pilot固定为40
		session.SquawkCode = squawkCode
		session.SquawkingModeC = squawkingModeC
		session.Identing = identing
		session.TrueAltitude = trueAlt
		session.PressureAltitude = pressureAlt
		session.GroundSpeed = gs
		session.Pitch = pitch
		session.Heading = heading
		session.Bank = bank
	}
}

// UpdateATCPosition 更新ATC位置信息
func (m *Manager) UpdateATCPosition(conn gnet.Conn, lat, lon float64, frequencies []string, facility, visRange, rating int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, exists := m.connToSession[conn]; exists {
		session.Lat = lat
		session.Lon = lon
		session.VisibilityRange = visRange
		session.Frequencies = frequencies
		session.Facility = facility
		session.Rating = rating
	}
}

// UpdateFlightPlan 更新Pilot的飞行计划（同时解除锁定）
func (m *Manager) UpdateFlightPlan(conn gnet.Conn, fp *FlightPlanData) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, exists := m.connToSession[conn]; exists {
		session.FlightPlan = fp
		session.FlightPlanLocked = false
	}
}

// UpdateFlightPlanByCallsign 通过callsign更新Pilot的飞行计划并锁定
func (m *Manager) UpdateFlightPlanByCallsign(callsign string, fp *FlightPlanData) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	conn, exists := m.callsignToConn[callsign]
	if !exists {
		return false
	}
	session, exists := m.connToSession[conn]
	if !exists {
		return false
	}
	session.FlightPlan = fp
	session.FlightPlanLocked = true
	return true
}

// GetFlightPlanLocked 获取指定连接的飞行计划锁定状态和当前飞行计划
func (m *Manager) GetFlightPlanLocked(conn gnet.Conn) (locked bool, fp *FlightPlanData) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if session, exists := m.connToSession[conn]; exists {
		return session.FlightPlanLocked, session.FlightPlan
	}
	return false, nil
}

// UpdateLastActivity 更新连接的最后活动时间
func (m *Manager) UpdateLastActivity(conn gnet.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, exists := m.connToSession[conn]; exists {
		session.LastActivity = time.Now()
	}
}

// GetIdleConns 获取超过指定时长未活动的连接列表
func (m *Manager) GetIdleConns(timeout time.Duration) []gnet.Conn {
	m.mu.RLock()
	defer m.mu.RUnlock()

	now := time.Now()
	var idle []gnet.Conn
	for conn, sess := range m.connToSession {
		if now.Sub(sess.LastActivity) > timeout {
			idle = append(idle, conn)
		}
	}
	return idle
}

// GetAuthTimeoutConns 获取认证/登录超时的连接列表
// 1. 连接后未在timeout内完成$ID认证的连接
// 2. $ID认证后未在timeout内完成登录（设置callsign）的连接
func (m *Manager) GetAuthTimeoutConns(timeout time.Duration) []gnet.Conn {
	m.mu.RLock()
	defer m.mu.RUnlock()

	now := time.Now()
	var result []gnet.Conn
	for conn, sess := range m.connToSession {
		// 已登录的连接不需要检查
		if sess.IsLoggedIn() {
			continue
		}
		// 正在关闭的连接不需要再踢
		if sess.Closing {
			continue
		}

		if !sess.Authenticated {
			// 未完成$ID认证，检查连接时间
			if now.Sub(sess.LogonTime) > timeout {
				result = append(result, conn)
			}
		} else {
			// 已完成$ID认证但未登录，检查认证完成时间
			if now.Sub(sess.AuthenticatedTime) > timeout {
				result = append(result, conn)
			}
		}
	}
	return result
}

// UpdateTextAtis 处理ATIS文本更新
func (m *Manager) UpdateTextAtis(conn gnet.Conn, payload []string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.connToSession[conn]
	if !exists {
		return
	}

	if len(payload) < 1 {
		return
	}

	switch payload[0] {
	case "T":
		if len(payload) >= 2 {
			session.TextAtis = append(session.TextAtis, payload[1])
		}
	case "V":
		session.TextAtis = make([]string, 0)
	case "Z":
		if len(payload) >= 2 {
			session.TextAtis = append(session.TextAtis, "Expected logoff time: "+payload[1])
		}
	}
}

// SetRating 设置用户登录等级
func (m *Manager) SetRating(conn gnet.Conn, rating int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, exists := m.connToSession[conn]; exists {
		session.Rating = rating
	}
}

// SetCid 设置用户CID
func (m *Manager) SetCid(conn gnet.Conn, cid string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, exists := m.connToSession[conn]; exists {
		session.Cid = cid
	}
}
