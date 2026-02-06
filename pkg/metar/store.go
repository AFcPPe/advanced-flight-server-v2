package metar

import (
	"bufio"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"advanced-flight-server/pkg/config"
	"advanced-flight-server/pkg/logger"

	"go.uber.org/zap"
)

var (
	store     *Store
	storeOnce sync.Once
)

// Store METAR数据存储，线程安全
// 使用 atomic.Value 实现无锁读，fetchMu 防止并发拉取
type Store struct {
	data    atomic.Value // 存储 map[string]string
	fetchMu sync.Mutex   // 防止 Fetch 并发重入
}

// GetStore 获取全局METAR存储实例
func GetStore() *Store {
	storeOnce.Do(func() {
		store = &Store{}
		store.data.Store(make(map[string]string))
	})
	return store
}

// Get 根据ICAO获取METAR数据（无锁读）
func (s *Store) Get(icao string) (string, bool) {
	data := s.data.Load().(map[string]string)
	metar, ok := data[strings.ToUpper(icao)]
	return metar, ok
}

// Count 返回当前存储的METAR条目数（无锁读）
func (s *Store) Count() int {
	data := s.data.Load().(map[string]string)
	return len(data)
}

// Fetch 从配置的URL拉取METAR数据并更新存储
// 使用 fetchMu 防止并发重入，同一时刻只允许一个 Fetch 执行
func (s *Store) Fetch() {
	if !s.fetchMu.TryLock() {
		logger.Info("METAR fetch already in progress, skipping")
		return
	}
	defer s.fetchMu.Unlock()

	cfg := config.Get().Metar
	url := cfg.URL
	if url == "" {
		url = "http://metar.vatsim.net/metar.php?id=ALL"
	}

	logger.Info("fetching METAR data", zap.String("url", url))

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		logger.Error("failed to fetch METAR data", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("METAR fetch returned non-200 status", zap.Int("status", resp.StatusCode))
		return
	}

	newData := make(map[string]string)
	scanner := bufio.NewScanner(resp.Body)
	// 增大缓冲区以处理可能的长行
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) < 4 {
			continue
		}
		// 前4位是ICAO
		icao := strings.ToUpper(line[:4])
		// 简单校验：ICAO应为字母数字
		if !isValidICAO(icao) {
			continue
		}
		newData[icao] = line
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		logger.Error("error reading METAR response body", zap.Error(err))
		return
	}

	if len(newData) == 0 {
		logger.Warn("fetched METAR data is empty, keeping old data")
		return
	}

	s.data.Store(newData)

	logger.Info("METAR data updated", zap.Int("count", len(newData)))
}

// isValidICAO 简单校验ICAO是否为4位字母/数字
func isValidICAO(icao string) bool {
	if len(icao) != 4 {
		return false
	}
	for _, c := range icao {
		if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			return false
		}
	}
	return true
}
