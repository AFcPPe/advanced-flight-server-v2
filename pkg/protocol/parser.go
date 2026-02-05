package protocol

import "strings"

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
		if len(data) > 1 {
			packet.Tokens = strings.Split(string(data[1:]), ":")
		}
	case '@':
		packet.Type = PacketTypePilotPosition
		if len(data) > 1 {
			packet.Tokens = strings.Split(string(data[1:]), ":")
		}
	case '#':
		packet.Type = PacketTypeHash
		if len(data) >= 3 {
			packet.SubType = string(data[1:3])
			if len(data) > 3 {
				packet.Tokens = strings.Split(string(data[3:]), ":")
			}
		}
	case '$':
		packet.Type = PacketTypeDollar
		if len(data) >= 3 {
			packet.SubType = string(data[1:3])
			if len(data) > 3 {
				packet.Tokens = strings.Split(string(data[3:]), ":")
			}
		}
	default:
		packet.Type = PacketTypeUnknown
	}

	return packet
}
