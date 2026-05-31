//go:build cluster

// Redis Cluster integration test, gated behind the `cluster` build tag so normal
// `go test` (miniredis) never needs a cluster. Run against a live cluster:
//
//	docker run -d --name rc -e IP=0.0.0.0 -p 7000-7005:7000-7005 grokzen/redis-cluster:7.0.10
//	go test -tags cluster ./cache/
//
// Override the seed addresses with BEACON_REDIS_CLUSTER (comma-separated).
//
// NOTE: run this on Linux/CI (host networking or a real cluster). macOS Docker
// Desktop's NAT can't reliably proxy a Redis Cluster client: the seed node
// connects, but cross-node routing times out regardless of port mapping — a
// known Docker-Desktop limitation, not a fault in this code.
package cache_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/iVampireSP/beacon/cache"
)

func newClusterStore(t *testing.T) *cache.CacheStore {
	t.Helper()
	addrs := os.Getenv("BEACON_REDIS_CLUSTER")
	if addrs == "" {
		addrs = "127.0.0.1:7000,127.0.0.1:7001,127.0.0.1:7002,127.0.0.1:7003,127.0.0.1:7004,127.0.0.1:7005"
	}
	client := cache.NewCache(cache.RedisConfig{ClusterAddrs: addrs})
	t.Cleanup(func() { _ = client.Close() })
	return cache.NewCacheStore(client)
}

// TestClusterCrossSlotReads is the core worry: in a cluster, keys hash to many
// different slots/nodes. Because every facade op is single-key, go-redis routes
// each to the owning node — so all reads succeed (no "can't read", no CROSSSLOT).
func TestClusterCrossSlotReads(t *testing.T) {
	s := newClusterStore(t)
	ctx := context.Background()

	const n = 20
	for i := 0; i < n; i++ {
		if err := s.Put(ctx, fmt.Sprintf("xslot:%d", i), i, time.Minute); err != nil {
			t.Fatalf("put %d: %v", i, err)
		}
	}
	for i := 0; i < n; i++ {
		v, found, err := cache.Get[int](ctx, s, fmt.Sprintf("xslot:%d", i))
		if err != nil || !found || v != i {
			t.Fatalf("get %d across slots: err=%v found=%v v=%d", i, err, found, v)
		}
	}
}

// TestClusterFlushFansOut verifies Flush clears keys on EVERY master (ForEachMaster),
// not just one node — keys are spread across slots, so a single-node FLUSHDB
// would leave most behind.
func TestClusterFlushFansOut(t *testing.T) {
	s := newClusterStore(t)
	ctx := context.Background()

	const n = 20
	for i := 0; i < n; i++ {
		_ = s.Put(ctx, fmt.Sprintf("flush:%d", i), i, time.Minute)
	}
	if err := s.Flush(ctx); err != nil {
		t.Fatalf("flush: %v", err)
	}
	for i := 0; i < n; i++ {
		if has, _ := s.Has(ctx, fmt.Sprintf("flush:%d", i)); has {
			t.Fatalf("key flush:%d survived Flush — fan-out missed its master", i)
		}
	}
}

// TestClusterRememberAndLock exercises Remember (cache key + a separate :lock key
// that may live on a different slot) and the distributed lock, on a real cluster.
func TestClusterRememberAndLock(t *testing.T) {
	s := newClusterStore(t)
	ctx := context.Background()

	calls := 0
	v, err := cache.Remember(ctx, s, "remember:key", time.Minute, func() (string, error) {
		calls++
		return "computed", nil
	})
	if err != nil || v != "computed" {
		t.Fatalf("remember: err=%v v=%q", err, v)
	}
	if v, _ := cache.Remember(ctx, s, "remember:key", time.Minute, func() (string, error) {
		calls++
		return "x", nil
	}); v != "computed" || calls != 1 {
		t.Fatalf("remember not cached: v=%q calls=%d", v, calls)
	}

	lk, err := s.Lock(ctx, "lock:cluster", time.Minute, nil)
	if err != nil {
		t.Fatalf("lock: %v", err)
	}
	if _, err := s.Lock(ctx, "lock:cluster", time.Minute, nil); err != cache.ErrNotObtained {
		t.Fatalf("contended lock: want ErrNotObtained, got %v", err)
	}
	if err := lk.Release(ctx); err != nil {
		t.Fatalf("release: %v", err)
	}
}
