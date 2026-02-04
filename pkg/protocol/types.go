package protocol

// PacketType 包类型
type PacketType int

const (
	PacketTypeUnknown PacketType = iota
	PacketTypeATCPosition   // % 开头
	PacketTypePilotPosition // @ 开头
	PacketTypeHash          // # 开头，需要进一步解析
	PacketTypeDollar        // $ 开头，需要进一步解析
)

// HashSubType # 开头的子类型
type HashSubType string

// DollarSubType $ 开头的子类型
type DollarSubType string
