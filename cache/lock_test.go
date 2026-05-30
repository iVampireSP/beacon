package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/iVampireSP/beacon/cache"
)

func TestLockContention(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()

	lk, err := s.Lock(ctx, "lock:x", time.Minute, nil)
	if err != nil {
		t.Fatalf("obtain: %v", err)
	}
	// A second acquire while held must fail (single try → ErrNotObtained).
	if _, err := s.Lock(ctx, "lock:x", time.Minute, nil); err != cache.ErrNotObtained {
		t.Fatalf("second acquire: want ErrNotObtained, got %v", err)
	}
	// After release it is obtainable again.
	if err := lk.Release(ctx); err != nil {
		t.Fatalf("release: %v", err)
	}
	lk2, err := s.Lock(ctx, "lock:x", time.Minute, nil)
	if err != nil {
		t.Fatalf("re-acquire after release: %v", err)
	}
	_ = lk2.Release(ctx)
}

func TestLockDoubleReleaseIsNotHeld(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	lk, err := s.Lock(ctx, "lock:y", time.Minute, nil)
	if err != nil {
		t.Fatalf("obtain: %v", err)
	}
	if err := lk.Release(ctx); err != nil {
		t.Fatalf("first release: %v", err)
	}
	if err := lk.Release(ctx); err != cache.ErrLockNotHeld {
		t.Fatalf("double release: want ErrLockNotHeld, got %v", err)
	}
}

func TestLockRefreshExtendsTTL(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	lk, err := s.Lock(ctx, "lock:z", time.Second, nil)
	if err != nil {
		t.Fatalf("obtain: %v", err)
	}
	if err := lk.Refresh(ctx, time.Minute); err != nil {
		t.Fatalf("refresh: %v", err)
	}
	ttl, err := lk.TTL(ctx)
	if err != nil {
		t.Fatalf("ttl: %v", err)
	}
	if ttl < 30*time.Second {
		t.Fatalf("ttl after refresh = %v, want ~1m", ttl)
	}
}

func TestAcquireHelperReleases(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	release, err := cache.Acquire(ctx, s.Locker(), "lock:acq", time.Minute)
	if err != nil {
		t.Fatalf("acquire: %v", err)
	}
	// Held: a direct acquire fails.
	if _, err := s.Lock(ctx, "lock:acq", time.Minute, nil); err != cache.ErrNotObtained {
		t.Fatalf("want held, got %v", err)
	}
	release()
	// Released: now obtainable.
	lk, err := s.Lock(ctx, "lock:acq", time.Minute, nil)
	if err != nil {
		t.Fatalf("after release(): %v", err)
	}
	_ = lk.Release(ctx)
}
