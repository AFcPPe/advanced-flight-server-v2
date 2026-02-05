package protocol

// Packet 解析后的包结构
type Packet struct {
	Type    PacketType
	SubType string // 对于 # 和 $ 类型，存储后两位字母如 "AP", "CQ"
	RawData []byte
	Tokens  []string // 解析后的字段（按:分割）
}

// GetTypeName 获取包类型名称（用于日志）
func (p *Packet) GetTypeName() string {
	switch p.Type {
	case PacketTypeATCPosition:
		return "ATCPosition"
	case PacketTypePilotPosition:
		return "PilotPosition"
	case PacketTypeHash:
		return "Hash#" + p.SubType
	case PacketTypeDollar:
		return "Dollar$" + p.SubType
	default:
		return "Unknown"
	}
}
