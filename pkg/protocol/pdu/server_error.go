package pdu

import (
	"fmt"
	"strconv"
)

// Error 错误PDU ($ER)
// 用于服务器向客户端发送错误消息，如登录失败原因
type Error struct {
	From  string       // 发送方
	To    string       // 目标客户端
	Error NetworkError // 错误码
	Param string       // 参数
	Msg   string       // 错误消息
}

// NewError 创建错误PDU
func NewError(from, to string, errorVal NetworkError, param, msg string) *Error {
	return &Error{
		From:  from,
		To:    to,
		Error: errorVal,
		Param: param,
		Msg:   msg,
	}
}

// NewServerError 创建服务器发送的错误PDU (便捷方法)
func NewServerError(to string, errorVal NetworkError, param, msg string) *Error {
	return NewError("SERVER", to, errorVal, param, msg)
}

// ErrorFromTokens 从tokens解析Error
// 格式: From:To:Error:Param:Msg
func ErrorFromTokens(tokens []string) (*Error, error) {
	if len(tokens) < 5 {
		return nil, fmt.Errorf("Error: invalid token count, got %d, need at least 5", len(tokens))
	}

	errorVal, err := strconv.Atoi(tokens[2])
	if err != nil {
		return nil, fmt.Errorf("Error: invalid error value: %w", err)
	}

	return NewError(tokens[0], tokens[1], NetworkError(errorVal), tokens[3], tokens[4]), nil
}

func (p *Error) GetHeader() string {
	return "$ER"
}

func (p *Error) ToTokens() []string {
	return []string{p.From, p.To, strconv.Itoa(int(p.Error)), p.Param, p.Msg}
}
