package pdu

// PDU 协议数据单元接口
type PDU interface {
	// GetHeader 获取协议头 (如 "#AP", "#AA", "%", "@")
	GetHeader() string
	// ToTokens 序列化为tokens
	ToTokens() []string
}
