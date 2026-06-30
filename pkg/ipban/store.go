package ipban

import (
	"encoding/json"
	"net"
	"os"
	"sync"
	"sync/atomic"

	"advanced-flight-server/pkg/logger"

	"go.uber.org/zap"
)

var (
	store     *Store
	storeOnce sync.Once
)

// entry 已解析的封禁规则
type entry struct {
	net    *net.IPNet
	action Action
	note   string
}

// Store IP封禁规则存储，线程安全
// 使用 atomic.Value 实现无锁读，loadMu 防止并发加载
type Store struct {
	entries atomic.Value // 存储 []entry
	loadMu  sync.Mutex   // 防止 Load 并发重入
	path    atomic.Value // 存储 string，当前规则文件路径
}

// GetStore 获取全局IP封禁存储实例
func GetStore() *Store {
	storeOnce.Do(func() {
		store = &Store{}
		store.entries.Store([]entry{})
		store.path.Store("")
	})
	return store
}

// Match 判断给定IP是否命中封禁规则，命中则返回对应动作和true
// 按规则在文件中的顺序匹配，第一条命中的规则生效（无锁读）
func (s *Store) Match(ip net.IP) (Action, bool) {
	if ip == nil {
		return "", false
	}
	entries := s.entries.Load().([]entry)
	for _, e := range entries {
		if e.net.Contains(ip) {
			return e.action, true
		}
	}
	return "", false
}

// MatchAddr 解析 host:port 形式的地址并匹配封禁规则
func (s *Store) MatchAddr(remoteAddr string) (Action, bool) {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		// 没有端口，按纯IP处理
		host = remoteAddr
	}
	return s.Match(net.ParseIP(host))
}

// Count 返回当前已加载的封禁规则条数（无锁读）
func (s *Store) Count() int {
	return len(s.entries.Load().([]entry))
}

// Load 从指定路径读取封禁规则文件并替换当前规则
// 文件不存在时自动创建一份带示例的默认文件
// 使用 loadMu 防止并发重入
func (s *Store) Load(path string) {
	if !s.loadMu.TryLock() {
		logger.Info("ip ban reload already in progress, skipping")
		return
	}
	defer s.loadMu.Unlock()

	s.path.Store(path)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Info("ip ban file not found, creating default", zap.String("path", path))
			if werr := writeDefaultFile(path); werr != nil {
				logger.Error("failed to create default ip ban file", zap.String("path", path), zap.Error(werr))
				return
			}
			// 默认文件没有真实生效的规则，直接置空
			s.entries.Store([]entry{})
			return
		}
		logger.Error("failed to read ip ban file", zap.String("path", path), zap.Error(err))
		return
	}

	var f File
	if err := json.Unmarshal(data, &f); err != nil {
		logger.Error("failed to parse ip ban file, keeping old rules", zap.String("path", path), zap.Error(err))
		return
	}

	newEntries := make([]entry, 0, len(f.Rules))
	for _, r := range f.Rules {
		ipNet, perr := parseCIDR(r.CIDR)
		if perr != nil {
			logger.Warn("skipping invalid ip ban rule", zap.String("cidr", r.CIDR), zap.Error(perr))
			continue
		}
		action := r.Action
		if action != ActionReject && action != ActionSilent {
			logger.Warn("skipping ip ban rule with unknown action",
				zap.String("cidr", r.CIDR), zap.String("action", string(r.Action)))
			continue
		}
		newEntries = append(newEntries, entry{net: ipNet, action: action, note: r.Note})
	}

	s.entries.Store(newEntries)
	logger.Info("ip ban rules loaded", zap.String("path", path), zap.Int("count", len(newEntries)))
}

// Reload 使用上一次 Load 的路径重新加载（供定时任务调用）
func (s *Store) Reload() {
	path, _ := s.path.Load().(string)
	if path == "" {
		return
	}
	s.Load(path)
}

// parseCIDR 解析CIDR或单个IP，单个IP按 /32（IPv4）或 /128（IPv6）处理
func parseCIDR(s string) (*net.IPNet, error) {
	if _, ipNet, err := net.ParseCIDR(s); err == nil {
		return ipNet, nil
	}
	// 尝试按单个IP解析
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, &net.ParseError{Type: "IP/CIDR address", Text: s}
	}
	if v4 := ip.To4(); v4 != nil {
		return &net.IPNet{IP: v4, Mask: net.CIDRMask(32, 32)}, nil
	}
	return &net.IPNet{IP: ip, Mask: net.CIDRMask(128, 128)}, nil
}

// writeDefaultFile 写出一份带示例规则的默认封禁文件
// 示例使用 RFC5737 文档保留地址，不会影响真实流量
func writeDefaultFile(path string) error {
	def := File{
		Rules: []Rule{
			{CIDR: "192.0.2.0/24", Action: ActionReject, Note: "示例：整段直接拒绝连接（建立即断开）"},
			{CIDR: "198.51.100.7", Action: ActionSilent, Note: "示例：单个IP静默处理（伪装服务器崩溃）"},
		},
	}
	data, err := json.MarshalIndent(def, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
