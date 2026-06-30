package snapshot

import (
	"sync"
	"testing"

	"advanced-flight-server/pkg/session"

	"github.com/panjf2000/gnet/v2"
)

// fakeConn 是一个最小的 gnet.Conn 桩，仅用于把 Session 放进 Manager。
// Manager 只把 conn 当作 map 的 key，从不调用其方法，因此内嵌 nil 接口即可满足编译。
type fakeConn struct {
	gnet.Conn
	id int
}

// TestPublishNoRaceWithSessionMutation 在 -race 下验证：
// 快照发布与 event-loop 风格的会话字段并发读写不再产生数据竞争。
// 复现的是导致整服闪退的根因——快照在锁外读 *Session 的切片/指针字段。
func TestPublishNoRaceWithSessionMutation(t *testing.T) {
	mgr := session.GetManager()
	Init()
	pub := GetPublisher()

	const n = 16
	conns := make([]*fakeConn, n)
	for i := range n {
		c := &fakeConn{id: i}
		conns[i] = c
		mgr.AddConn(c)
		mgr.SetCallsign(c, callsignOf(i))
		mgr.SetConnType(c, session.ConnectionTypeATC)
		// 标记已认证，使 IsLoggedIn 为 true，从而进入快照采集
		if s := mgr.GetSession(c); s != nil {
			s.Authenticated = true
		}
	}

	var wg sync.WaitGroup
	stop := make(chan struct{})

	// writer：模拟 event-loop 持续更新会话字段（频率/ATIS 切片等）
	for i := range n {
		wg.Add(1)
		go func(c *fakeConn) {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				mgr.UpdateATCPosition(c, 30, 120, []string{"18100", "12150"}, 4, 50, 4)
				mgr.UpdateTextAtis(c, []string{"V"})
				mgr.UpdateTextAtis(c, []string{"T", "line one"})
			}
		}(conns[i])
	}

	// reader：快照发布（无 redis 时 Publish 内部会安全跳过写入）
	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				pub.Publish()
			}
		}()
	}

	// 跑若干轮后停止
	for range 200 {
		pub.Publish()
	}
	close(stop)
	wg.Wait()

	// 清理
	for _, c := range conns {
		mgr.RemoveConn(c)
	}
}

func callsignOf(i int) string {
	return "ATC" + string(rune('A'+i)) + "_CTR"
}
