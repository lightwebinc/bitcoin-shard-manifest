package sender

import (
	"testing"

	"github.com/lightwebinc/shard-manifest/config"
)

func TestResolveGroups_JoinAll(t *testing.T) {
	c := &config.Config{ShardBits: 3, JoinAll: true}
	groups, claim := resolveGroups(c)
	if !claim {
		t.Fatal("JoinAll must produce a claim")
	}
	if len(groups) != 8 {
		t.Fatalf("len = %d, want 8", len(groups))
	}
	for i, g := range groups {
		if int(g) != i {
			t.Errorf("groups[%d] = %d", i, g)
		}
	}
}

func TestResolveGroups_Explicit(t *testing.T) {
	c := &config.Config{ShardBits: 3, JoinedGroups: []uint16{1, 4}}
	groups, claim := resolveGroups(c)
	if !claim || len(groups) != 2 || groups[0] != 1 || groups[1] != 4 {
		t.Fatalf("explicit groups = %v claim=%v", groups, claim)
	}
	// Returned slice must be a copy, not the config's backing array.
	groups[0] = 99
	if c.JoinedGroups[0] != 1 {
		t.Error("resolveGroups must not alias the config slice")
	}
}

func TestResolveGroups_NoClaim(t *testing.T) {
	c := &config.Config{ShardBits: 3}
	if groups, claim := resolveGroups(c); claim || groups != nil {
		t.Errorf("empty config should make no claim: groups=%v claim=%v", groups, claim)
	}
}

func TestJoinCount(t *testing.T) {
	if got := joinCount(&config.Config{ShardBits: 4, JoinAll: true}); got != 16 {
		t.Errorf("JoinAll count = %d, want 16", got)
	}
	if got := joinCount(&config.Config{JoinedGroups: []uint16{1, 2, 3}}); got != 3 {
		t.Errorf("explicit count = %d, want 3", got)
	}
}

func TestCurrentSources(t *testing.T) {
	s := &Sender{}
	if got := s.currentSources(); got != nil {
		t.Errorf("empty sources should return nil, got %v", got)
	}
	s.sources = [][16]byte{{0x01}, {0x02}}
	got := s.currentSources()
	if len(got) != 2 {
		t.Fatalf("len = %d", len(got))
	}
	// Mutating the returned slice must not affect the cached set.
	got[0] = [16]byte{0xFF}
	if s.sources[0] != [16]byte{0x01} {
		t.Error("currentSources must return a copy")
	}
}
