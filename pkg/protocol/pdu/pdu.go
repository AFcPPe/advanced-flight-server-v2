package pdu

import "strings"

// PDU 协议数据单元接口
type PDU interface {
	// GetHeader 获取协议头 (如 "#AP", "#AA", "%", "@")
	GetHeader() string
	// ToTokens 序列化为tokens
	ToTokens() []string
}

// Serialize 将PDU序列化为可发送的字节数据（带\r\n结尾）
func Serialize(p PDU) []byte {
	tokens := p.ToTokens()
	return []byte(p.GetHeader() + strings.Join(tokens, ":") + "\r\n")
}
