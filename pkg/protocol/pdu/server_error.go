package pdu

import (
	"fmt"
	"strconv"
)

type ServerError struct {
	Base
	ErrorCode NetworkError
	Msg       string
	Param     string
}

func NewServerError(from, to string, errorCode NetworkError, param, msg string) *ServerError {
	return &ServerError{
		Base: Base{
			From: from,
			To:   to,
		},
		ErrorCode: errorCode,
		Msg:       msg,
		Param:     param,
	}
}

func ServerErrorFromTokens(tokens []string) (*ServerError, error) {
	if len(tokens) < 5 {
		return nil, fmt.Errorf("server error: invalid token count, expected >= 5, got %d", len(tokens))
	}

	errorVal, err := strconv.Atoi(tokens[2])
	if err != nil {
		return nil, fmt.Errorf("server error: invalid error code: %w", err)
	}

	return NewServerError(tokens[0], tokens[1], NetworkError(errorVal), tokens[3], tokens[4]), nil
}

func (p *ServerError) ToTokens() []string {
	return []string{p.From, p.To, strconv.Itoa(int(p.ErrorCode)), p.Param, p.Msg}
}

func (p *ServerError) GetHeader() string {
	return "$ER"
}
