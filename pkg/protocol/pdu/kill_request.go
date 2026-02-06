package pdu

import "fmt"

type KillRequest struct {
	Base
	Reason string
}

func NewPDUKillRequest(from, target, reason string) *KillRequest {
	return &KillRequest{
		Base: Base{
			From: from,
			To:   target,
		},
		Reason: reason,
	}
}

func KillRequestFromTokens(tokens []string) (*KillRequest, error) {
	if len(tokens) < 3 {
		return nil, fmt.Errorf("kill request: invalid token count, expected >= 3, got %d", len(tokens))
	}
	reason := tokens[2]
	for i := 3; i < len(tokens); i++ {
		reason += ":" + tokens[i]
	}
	return NewPDUKillRequest(tokens[0], tokens[1], reason), nil
}

func (p *KillRequest) ToTokens() []string {
	return []string{p.From, p.To, p.Reason}
}

func (p *KillRequest) GetHeader() string {
	return "$!!"
}
