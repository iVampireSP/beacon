package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/iVampireSP/beacon/cache"
	"github.com/redis/go-redis/v9"
)

// newStore spins up an in-memory miniredis (no external Redis needed) and a Store
// over it. Single-node, so it exercises the facade logic + the lock Lua; the
// cluster-specific Flush fan-out and slot routing are go-redis concerns that need
// a real cluster.
func newStore(t *testing.T) *cache.Store {
	t.Helper()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = client.Close() })
	return cache.NewStore(client)
}

func TestPutGet(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	type user struct {
		Name string
		Age  int
	}
	if err := s.Put(ctx, "u", user{"alice", 30}, time.Minute); err != nil {
		t.Fatalf("put: %v", err)
	}
	got, found, err := cache.Get[user](ctx, s, "u")
	if err != nil || !found {
		t.Fatalf("get: err=%v found=%v", err, found)
	}
	if got.Name != "alice" || got.Age != 30 {
		t.Fatalf("got %+v", got)
	}
	if _, found, _ := cache.Get[user](ctx, s, "missing"); found {
		t.Fatal("expected miss for absent key")
	}
}

func TestAddOnlyIfAbsent(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	if ok, _ := s.Add(ctx, "k", "v1", time.Minute); !ok {
		t.Fatal("first Add should store")
	}
	if ok, _ := s.Add(ctx, "k", "v2", time.Minute); ok {
		t.Fatal("second Add must not overwrite")
	}
	if v, _, _ := cache.Get[string](ctx, s, "k"); v != "v1" {
		t.Fatalf("want v1, got %q", v)
	}
}

func TestHasForget(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	_ = s.Put(ctx, "k", 1, time.Minute)
	if has, _ := s.Has(ctx, "k"); !has {
		t.Fatal("Has should be true")
	}
	_ = s.Forget(ctx, "k")
	if has, _ := s.Has(ctx, "k"); has {
		t.Fatal("Has should be false after Forget")
	}
}

func TestIncrementDecrement(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	if n, _ := s.Increment(ctx, "c"); n != 1 {
		t.Fatalf("inc default = %d, want 1", n)
	}
	if n, _ := s.Increment(ctx, "c", 5); n != 6 {
		t.Fatalf("inc 5 = %d, want 6", n)
	}
	if n, _ := s.Decrement(ctx, "c", 2); n != 4 {
		t.Fatalf("dec 2 = %d, want 4", n)
	}
	// JSON-encoded integers stay INCR-compatible: Put(int) then Increment.
	_ = s.Put(ctx, "c2", 10, time.Minute)
	if n, _ := s.Increment(ctx, "c2"); n != 11 {
		t.Fatalf("put-then-inc = %d, want 11", n)
	}
	if v, _, _ := cache.Get[int](ctx, s, "c2"); v != 11 {
		t.Fatalf("get int = %d, want 11", v)
	}
}

func TestPullDeletes(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	_ = s.Put(ctx, "k", "v", time.Minute)
	v, found, err := cache.Pull[string](ctx, s, "k")
	if err != nil || !found || v != "v" {
		t.Fatalf("pull: err=%v found=%v v=%q", err, found, v)
	}
	if has, _ := s.Has(ctx, "k"); has {
		t.Fatal("Pull must delete the key")
	}
}

func TestRememberCachesAndComputesOnce(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	calls := 0
	fn := func() (int, error) { calls++; return 42, nil }

	if v, err := cache.Remember(ctx, s, "k", time.Minute, fn); err != nil || v != 42 {
		t.Fatalf("first remember: err=%v v=%d", err, v)
	}
	if v, err := cache.Remember(ctx, s, "k", time.Minute, fn); err != nil || v != 42 {
		t.Fatalf("second remember: err=%v v=%d", err, v)
	}
	if calls != 1 {
		t.Fatalf("compute fn called %d times, want 1 (second is a cache hit)", calls)
	}
}

func TestFlush(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	_ = s.Put(ctx, "a", 1, time.Minute)
	_ = s.Put(ctx, "b", 2, time.Minute)
	if err := s.Flush(ctx); err != nil {
		t.Fatalf("flush: %v", err)
	}
	if has, _ := s.Has(ctx, "a"); has {
		t.Fatal("Flush should clear all keys")
	}
}
