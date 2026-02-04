package pdu

import "fmt"

// TextMessage 文本消息PDU (#TM)
type TextMessage struct {
	From    string
	To      string
	Message string
}

// NewTextMessage 创建文本消息PDU
func NewTextMessage(from, to, message string) *TextMessage {
	return &TextMessage{
		From:    from,
		To:      to,
		Message: message,
	}
}

// NewServerTextMessage 创建服务器发送的文本消息PDU (便捷方法)
func NewServerTextMessage(to, message string) *TextMessage {
	return NewTextMessage("SERVER", to, message)
}

// TextMessageFromTokens 从tokens解析TextMessage
// 格式: From:To:Message (Message可能包含冒号)
func TextMessageFromTokens(tokens []string) (*TextMessage, error) {
	if len(tokens) < 3 {
		return nil, fmt.Errorf("TextMessage: invalid token count, got %d, need at least 3", len(tokens))
	}
	message := tokens[2]
	for i := 3; i < len(tokens); i++ {
		message += ":" + tokens[i]
	}
	return NewTextMessage(tokens[0], tokens[1], message), nil
}

func (p *TextMessage) GetHeader() string {
	return "#TM"
}

func (p *TextMessage) ToTokens() []string {
	return []string{p.From, p.To, p.Message}
}
