package snapshot

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
	"sync"
	"time"

	"advanced-flight-server/pkg/logger"
	"advanced-flight-server/pkg/redis"
	"advanced-flight-server/pkg/session"

	"go.uber.org/zap"
)

const (
	// KeyUsers Redis键：在线用户快照
	KeyUsers = "afs:snapshot:users"
	// TTL 快照过期时间，兜底防止服务异常退出后残留脏数据
	TTL = 10 * time.Second
	// Interval 快照推送间隔（秒）
	Interval = 3
	// atisSuffix ATIS类型callsign后缀
	atisSuffix = "_ATIS"
)

var (
	publisher *Publisher
	once      sync.Once
)

// Publisher 快照发布器
type Publisher struct {
	sessionMgr *session.Manager
}

// Init 初始化快照发布器
func Init() {
	once.Do(func() {
		publisher = &Publisher{
			sessionMgr: session.GetManager(),
		}
	})
}

// GetPublisher 获取快照发布器
func GetPublisher() *Publisher {
	return publisher
}

// Publish 采集当前所有在线用户信息，按类型分类后推送到Redis
func (p *Publisher) Publish() {
	var pilots []*Pilot
	var atcs []*ATC
	var atiss []*ATIS

	// 在持有会话管理器读锁的前提下完成采集与拷贝，
	// 避免与 event-loop/auth goroutine 并发读写 Session 字段产生数据竞争。
	p.sessionMgr.ForEachLoggedInSession(func(s *session.Session) {
		switch s.ConnType {
		case session.ConnectionTypePilot:
			pilots = append(pilots, buildPilot(s))
		case session.ConnectionTypeATC:
			if isATIS(s.Callsign) {
				atiss = append(atiss, buildATIS(s))
			} else {
				atcs = append(atcs, buildATC(s))
			}
		}
	})

	snap := &Snapshot{
		Timestamp: time.Now().Unix(),
		Pilots:    pilots,
		ATCs:      atcs,
		ATISs:     atiss,
	}

	data, err := json.Marshal(snap)
	if err != nil {
		logger.Error("snapshot: failed to marshal", zap.Error(err))
		return
	}

	client := redis.GetClient()
	if client == nil {
		logger.Warn("snapshot: redis client not available, skipping publish")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Set(ctx, KeyUsers, data, TTL).Err(); err != nil {
		logger.Error("snapshot: failed to publish to redis", zap.Error(err))
	}
}

// isATIS 判断callsign是否为ATIS类型（以_ATIS结尾）
func isATIS(callsign string) bool {
	return strings.HasSuffix(callsign, atisSuffix)
}

// buildPilot 从Session构建Pilot快照
func buildPilot(s *session.Session) *Pilot {
	p := &Pilot{
		CID:              s.Cid,
		Callsign:         s.Callsign,
		Name:             s.RealName,
		Lat:              s.Lat,
		Lon:              s.Lon,
		TrueAltitude:     s.TrueAltitude,
		PressureAltitude: s.PressureAltitude,
		GroundSpeed:      s.GroundSpeed,
		Pitch:            s.Pitch,
		Heading:          s.Heading,
		Bank:             s.Bank,
		Squawk:           s.SquawkCode,
		SquawkingModeC:   s.SquawkingModeC,
		Identing:         s.Identing,
		VisibilityRange:  s.VisibilityRange,
		FlightPlanLocked: s.FlightPlanLocked,
		Rating:           s.Rating,
		LogonTime:        s.LogonTime.UTC().Format(time.RFC3339),
	}
	if s.FlightPlan != nil {
		p.FlightPlan = buildFlightPlan(s.FlightPlan)
	}
	return p
}

// buildATC 从Session构建ATC快照
func buildATC(s *session.Session) *ATC {
	return &ATC{
		CID:             s.Cid,
		Callsign:        s.Callsign,
		Name:            s.RealName,
		Lat:             s.Lat,
		Lon:             s.Lon,
		Frequencies:     formatFrequencies(s.Frequencies),
		Facility:        s.Facility,
		Rating:          s.Rating,
		VisibilityRange: s.VisibilityRange,
		TextAtis:        slices.Clone(s.TextAtis),
		LogonTime:       s.LogonTime.UTC().Format(time.RFC3339),
	}
}

// buildATIS 从Session构建ATIS快照
func buildATIS(s *session.Session) *ATIS {
	return &ATIS{
		CID:             s.Cid,
		Callsign:        s.Callsign,
		Name:            s.RealName,
		Lat:             s.Lat,
		Lon:             s.Lon,
		Frequencies:     formatFrequencies(s.Frequencies),
		Facility:        s.Facility,
		Rating:          s.Rating,
		VisibilityRange: s.VisibilityRange,
		TextAtis:        slices.Clone(s.TextAtis),
		LogonTime:       s.LogonTime.UTC().Format(time.RFC3339),
	}
}

// formatFrequency 将原始频率字符串转换为标准格式（如 "18100" → "118.100"）
func formatFrequency(raw string) string {
	if len(raw) < 4 {
		return raw
	}
	intPart := "1" + raw[:2]
	decPart := raw[2:]
	return intPart + "." + decPart
}

// formatFrequencies 批量转换频率列表
func formatFrequencies(raw []string) []string {
	if len(raw) == 0 {
		return raw
	}
	result := make([]string, len(raw))
	for i, f := range raw {
		result[i] = formatFrequency(f)
	}
	return result
}

// buildFlightPlan 从FlightPlanData构建飞行计划快照
func buildFlightPlan(fp *session.FlightPlanData) *FlightPlan {
	return &FlightPlan{
		FlightRules:   fp.FlightRules.String(),
		AircraftType:  fp.Type,
		TAS:           fp.TAS,
		Departure:     fp.Dep,
		DepTime:       fp.DepTime,
		ActualDepTime: fp.ActualDepTime,
		CruiseAlt:     fp.CruiseAlt,
		Arrival:       fp.Dest,
		EnrouteHour:   fp.EnrouteHour,
		EnrouteMin:    fp.EnrouteMin,
		FobHour:       fp.FobHour,
		FobMin:        fp.FobMin,
		Alternate:     fp.AlterDest,
		Route:         fp.Route,
		Remark:        fp.Remark,
	}
}
