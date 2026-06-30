package ipban

// Action 封禁动作类型
type Action string

const (
	// ActionReject 直接拒绝连接：连接建立后立即断开
	ActionReject Action = "reject"
	// ActionSilent 静默处理：接受连接但不发送任何数据、丢弃所有输入，
	// 也不触发空闲/认证超时，让对方误以为服务器崩溃
	ActionSilent Action = "silent"
)

// Rule 单条封禁规则（JSON文件中的一项）
type Rule struct {
	CIDR   string `json:"cidr"`   // IP段（CIDR，如 192.0.2.0/24）或单个IP（如 192.0.2.7）
	Action Action `json:"action"` // 封禁动作：reject 或 silent
	Note   string `json:"note"`   // 备注，仅用于人工识别
}

// File 封禁规则JSON文件结构
type File struct {
	Rules []Rule `json:"rules"`
}
