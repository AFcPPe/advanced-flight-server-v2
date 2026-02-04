package pdu

import (
	"fmt"
	"strconv"
)

// AddPilot 飞行员登录PDU (#AP)
type AddPilot struct {
	Callsign        string
	Cid             string
	RealName        string
	Password        string        // TODO: 密码验证
	Rating          NetworkRating // TODO: 等级验证
	ProtocolVersion int
	SimType         int
}

// FromTokens 从tokens解析AddPilot
// 格式: From:To:Cid:Password:Rating:ProtocolVersion:SimType:RealName
func AddPilotFromTokens(tokens []string) (*AddPilot, error) {
	if len(tokens) < 8 {
		return nil, fmt.Errorf("AddPilot: invalid token count, got %d, need 8", len(tokens))
	}

	rating, err := strconv.Atoi(tokens[4])
	if err != nil {
		return nil, fmt.Errorf("AddPilot: invalid rating: %w", err)
	}

	protocolVersion, err := strconv.Atoi(tokens[5])
	if err != nil {
		return nil, fmt.Errorf("AddPilot: invalid protocol version: %w", err)
	}

	simType, err := strconv.Atoi(tokens[6])
	if err != nil {
		return nil, fmt.Errorf("AddPilot: invalid sim type: %w", err)
	}

	return &AddPilot{
		Callsign:        tokens[0],
		Cid:             tokens[2],
		RealName:        tokens[7],
		Password:        tokens[3],
		Rating:          NetworkRating(rating),
		ProtocolVersion: protocolVersion,
		SimType:         simType,
	}, nil
}

func (p *AddPilot) GetHeader() string {
	return "#AP"
}

func (p *AddPilot) ToTokens() []string {
	return []string{
		p.Callsign,
		"SERVER",
		p.Cid,
		p.Password,
		strconv.Itoa(int(p.Rating)),
		strconv.Itoa(p.ProtocolVersion),
		strconv.Itoa(p.SimType),
		p.RealName,
	}
}

// GetClientType 返回客户端类型
func (p *AddPilot) GetClientType() ClientType {
	return ClientTypePilot
}
