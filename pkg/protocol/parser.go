package protocol

// ParsePacket 解析原始数据包
func ParsePacket(data []byte) *Packet {
	if len(data) == 0 {
		return &Packet{Type: PacketTypeUnknown, RawData: data}
	}

	packet := &Packet{
		RawData: data,
	}

	header := data[0]
	switch header {
	case '%':
		packet.Type = PacketTypeATCPosition
	case '@':
		packet.Type = PacketTypePilotPosition
	case '#':
		packet.Type = PacketTypeHash
		if len(data) >= 3 {
			packet.SubType = string(data[1:3])
		}
	case '$':
		packet.Type = PacketTypeDollar
		if len(data) >= 3 {
			packet.SubType = string(data[1:3])
		}
	default:
		packet.Type = PacketTypeUnknown
	}

	return packet
}
