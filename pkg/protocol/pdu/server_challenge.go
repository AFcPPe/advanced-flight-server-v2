package pdu

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"
)

// ServerChallenge 服务器挑战PDU ($DI)
// 格式: $DISERVER:CLIENT:VATSIM FSD V3.41b:<ChallengeCode>
type ServerChallenge struct {
	From          string
	To            string
	ServerVersion string
	ChallengeCode string
}

// NewServerChallenge 创建一个新的ServerChallenge PDU，自动生成ChallengeCode
func NewServerChallenge() *ServerChallenge {
	return &ServerChallenge{
		From:          "SERVER",
		To:            "CLIENT",
		ServerVersion: "VATSIM FSD V3.41b",
		ChallengeCode: generateChallengeCode(),
	}
}

// generateChallengeCode 生成挑战码
func generateChallengeCode() string {
	data := []byte(strconv.FormatInt(time.Now().UnixNano(), 10))
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}

func (p *ServerChallenge) GetHeader() string {
	return "$DI"
}

func (p *ServerChallenge) ToTokens() []string {
	return []string{
		p.From,
		p.To,
		p.ServerVersion,
		p.ChallengeCode,
	}
}

// GetChallengeCode 返回当前的挑战码
func (p *ServerChallenge) GetChallengeCode() string {
	return p.ChallengeCode
}
