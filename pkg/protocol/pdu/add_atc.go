package pdu

import (
	"fmt"
	"strconv"
)

// AddATC 管制员登录PDU (#AA)
type AddATC struct {
	Callsign        string
	Cid             string
	RealName        string
	Password        string        // TODO: 密码验证
	Rating          NetworkRating // TODO: 等级验证
	ProtocolVersion int
}

// FromTokens 从tokens解析AddATC
// 格式: Callsign:To:Name:Cid:Password:Rating:ProtocolRevision
func AddATCFromTokens(tokens []string) (*AddATC, error) {
	if len(tokens) < 7 {
		return nil, fmt.Errorf("AddATC: invalid token count, got %d, need 7", len(tokens))
	}

	rating, err := strconv.Atoi(tokens[5])
	if err != nil {
		return nil, fmt.Errorf("AddATC: invalid rating: %w", err)
	}

	protocolVersion, err := strconv.Atoi(tokens[6])
	if err != nil {
		return nil, fmt.Errorf("AddATC: invalid protocol version: %w", err)
	}

	return &AddATC{
		Callsign:        tokens[0],
		Cid:             tokens[3],
		RealName:        tokens[2],
		Password:        tokens[4],
		Rating:          NetworkRating(rating),
		ProtocolVersion: protocolVersion,
	}, nil
}

func (p *AddATC) GetHeader() string {
	return "#AA"
}

func (p *AddATC) ToTokens() []string {
	return []string{
		p.Callsign,
		"SERVER",
		p.RealName,
		p.Cid,
		p.Password,
		strconv.Itoa(int(p.Rating)),
		strconv.Itoa(p.ProtocolVersion),
	}
}

// GetClientType 返回客户端类型
func (p *AddATC) GetClientType() ClientType {
	return ClientTypeATC
}
