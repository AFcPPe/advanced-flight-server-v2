package pdu

// ClientType 客户端类型
type ClientType int

const (
	ClientTypePilot ClientType = iota // 飞行员
	ClientTypeATC                     // 管制员
)

func (c ClientType) String() string {
	switch c {
	case ClientTypePilot:
		return "Pilot"
	case ClientTypeATC:
		return "ATC"
	default:
		return "Unknown"
	}
}

// NetworkRating 网络等级
type NetworkRating int
