package ipban

import (
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestParseCIDR(t *testing.T) {
	cases := []struct {
		in    string
		valid bool
		match string // 一个应命中的IP
		miss  string // 一个不应命中的IP
	}{
		{"192.0.2.0/24", true, "192.0.2.55", "192.0.3.1"},
		{"198.51.100.7", true, "198.51.100.7", "198.51.100.8"},
		{"2001:db8::/32", true, "2001:db8::1", "2001:dead::1"},
		{"not-an-ip", false, "", ""},
		{"10.0.0.0/99", false, "", ""},
	}

	for _, c := range cases {
		ipNet, err := parseCIDR(c.in)
		if c.valid && err != nil {
			t.Errorf("parseCIDR(%q) unexpected error: %v", c.in, err)
			continue
		}
		if !c.valid {
			if err == nil {
				t.Errorf("parseCIDR(%q) expected error, got nil", c.in)
			}
			continue
		}
		if !ipNet.Contains(net.ParseIP(c.match)) {
			t.Errorf("parseCIDR(%q) should contain %s", c.in, c.match)
		}
		if ipNet.Contains(net.ParseIP(c.miss)) {
			t.Errorf("parseCIDR(%q) should not contain %s", c.in, c.miss)
		}
	}
}

func TestStoreLoadAndMatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ban.json")
	content := `{
		"rules": [
			{"cidr": "203.0.113.0/24", "action": "reject", "note": "block range"},
			{"cidr": "198.51.100.10", "action": "silent", "note": "silence one"},
			{"cidr": "bad-rule", "action": "reject", "note": "should be skipped"},
			{"cidr": "192.0.2.1", "action": "bogus", "note": "bad action skipped"}
		]
	}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	s := &Store{}
	s.entries.Store([]entry{})
	s.path.Store("")
	s.Load(path)

	if got := s.Count(); got != 2 {
		t.Fatalf("expected 2 valid rules, got %d", got)
	}

	if a, ok := s.Match(net.ParseIP("203.0.113.55")); !ok || a != ActionReject {
		t.Errorf("203.0.113.55 expected reject, got %v ok=%v", a, ok)
	}
	if a, ok := s.Match(net.ParseIP("198.51.100.10")); !ok || a != ActionSilent {
		t.Errorf("198.51.100.10 expected silent, got %v ok=%v", a, ok)
	}
	if _, ok := s.Match(net.ParseIP("8.8.8.8")); ok {
		t.Errorf("8.8.8.8 should not match any rule")
	}
}

func TestMatchAddr(t *testing.T) {
	s := &Store{}
	_, ipNet, _ := net.ParseCIDR("203.0.113.0/24")
	s.entries.Store([]entry{{net: ipNet, action: ActionReject}})
	s.path.Store("")

	if a, ok := s.MatchAddr("203.0.113.9:51234"); !ok || a != ActionReject {
		t.Errorf("host:port form expected reject, got %v ok=%v", a, ok)
	}
	if a, ok := s.MatchAddr("203.0.113.9"); !ok || a != ActionReject {
		t.Errorf("bare ip form expected reject, got %v ok=%v", a, ok)
	}
}

func TestLoadCreatesDefaultFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.json")

	s := &Store{}
	s.entries.Store([]entry{})
	s.path.Store("")
	s.Load(path)

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("default ban file should have been created: %v", err)
	}
}
