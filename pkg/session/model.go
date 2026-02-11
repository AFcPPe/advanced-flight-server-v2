package session

import (
	"time"

	"advanced-flight-server/pkg/protocol/pdu"

	"github.com/panjf2000/gnet/v2"
)

// PendingPacket 认证期间缓存的待处理包（原始字节）
type PendingPacket struct {
	RawData []byte
}

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
	Conn              gnet.Conn
	Callsign          string
	Cid               string           // 用户CID
	RealName          string           // 用户真实姓名
	Buffer            []byte           // 用于处理粘包的缓存
	Authenticated     bool             // 是否已通过$ID验证（首包必须是$ID）
	AuthenticatedTime time.Time        // $ID认证完成的时间
	Authenticating    bool             // 是否正在进行异步登录认证（DB查询中）
	PendingPackets    []*PendingPacket // 认证期间缓存的待处理包
	ReplayFunc        func()           // 认证完成后重放缓存包的回调函数
	Closing           bool             // 是否正在关闭连接（发送错误后等待断开）
	ConnType          ConnectionType   // 连接类型：Pilot或ATC
	LastActivity      time.Time        // 最后一次收到数据的时间
	LogonTime         time.Time        // 连接建立的时间

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
	FlightPlan       *FlightPlanData // 飞行计划
	FlightPlanLocked bool            // 飞行计划是否被ATC锁定（被TagModify修改后锁定）
	Rating           int             // 登录时请求的等级（Pilot和ATC共用）

	// ATC特有字段
	Frequencies   []string
	Facility      int
	TextAtis      []string  // ATIS文本信息
	LastAtisQuery time.Time // 上次向该ATC询问ATIS的时间
}

// FinishAuthenticating 完成异步认证，清除认证状态并重放缓存包
// loginSuccess: 登录是否成功，只有成功时才重放缓存包
func (s *Session) FinishAuthenticating(loginSuccess bool) {
	s.Authenticating = false
	if loginSuccess && s.ReplayFunc != nil {
		s.ReplayFunc()
	}
	s.PendingPackets = nil
	s.ReplayFunc = nil
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

// FlightPlanData 存储飞行计划数据
type FlightPlanData struct {
	FlightRules   pdu.FlightRule
	Type          string
	TAS           string
	Dep           string
	DepTime       string
	ActualDepTime string
	CruiseAlt     string
	Dest          string
	EnrouteHour   string
	EnrouteMin    string
	FobHour       string
	FobMin        string
	AlterDest     string
	Remark        string
	Route         string
}

// NewFlightPlanData 从FlightPlan PDU创建FlightPlanData
func NewFlightPlanData(fp *pdu.FlightPlan) *FlightPlanData {
	return &FlightPlanData{
		FlightRules:   fp.FlightRules,
		Type:          fp.Type,
		TAS:           fp.TAS,
		Dep:           fp.Dep,
		DepTime:       fp.DepTime,
		ActualDepTime: fp.ActualDepTime,
		CruiseAlt:     fp.CruiseAlt,
		Dest:          fp.Dest,
		EnrouteHour:   fp.EnrouteHour,
		EnrouteMin:    fp.EnrouteMin,
		FobHour:       fp.FobHour,
		FobMin:        fp.FobMin,
		AlterDest:     fp.AlterDest,
		Remark:        fp.Remark,
		Route:         fp.Route,
	}
}

// NewFlightPlanDataFromTagModify 从TagModify PDU创建FlightPlanData
func NewFlightPlanDataFromTagModify(tm *pdu.TagModify) *FlightPlanData {
	return &FlightPlanData{
		FlightRules:   tm.FlightRules,
		Type:          tm.Type,
		TAS:           tm.TAS,
		Dep:           tm.Dep,
		DepTime:       tm.DepTime,
		ActualDepTime: tm.ActualDepTime,
		CruiseAlt:     tm.CruiseAlt,
		Dest:          tm.Dest,
		EnrouteHour:   tm.EnrouteHour,
		EnrouteMin:    tm.EnrouteMin,
		FobHour:       tm.FobHour,
		FobMin:        tm.FobMin,
		AlterDest:     tm.AlterDest,
		Remark:        tm.Remark,
		Route:         tm.Route,
	}
}
