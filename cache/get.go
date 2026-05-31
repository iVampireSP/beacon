package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// These are free functions, not methods, because Go methods cannot carry type
// parameters. They JSON-decode the stored value into T.

// Get retrieves and decodes key into T. found is false on a cache miss.
func Get[T any](ctx context.Context, s *CacheStore, key string) (value T, found bool, err error) {
	data, err := s.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return value, false, nil
	}
	if err != nil {
		return value, false, err
	}
	if err = json.Unmarshal(data, &value); err != nil {
		return value, false, err
	}
	return value, true, nil
}

// Pull retrieves key and deletes it atomically via GETDEL (a single-key command,
// so it is cluster-safe). found is false on a cache miss.
func Pull[T any](ctx context.Context, s *CacheStore, key string) (value T, found bool, err error) {
	data, err := s.client.GetDel(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return value, false, nil
	}
	if err != nil {
		return value, false, err
	}
	if err = json.Unmarshal(data, &value); err != nil {
		return value, false, err
	}
	return value, true, nil
}

// Remember returns the cached value for key, or computes it with fn, stores it
// for ttl, and returns it. A per-key distributed lock guards the computation so
// concurrent callers don't all hit the source (stampede); losers wait then
// re-read. Lock and value are separate single-key operations — never one
// cross-slot atomic op — so this is safe under Redis Cluster.
func Remember[T any](ctx context.Context, s *CacheStore, key string, ttl time.Duration, fn func() (T, error)) (T, error) {
	if v, found, err := Get[T](ctx, s, key); err != nil {
		var zero T
		return zero, err
	} else if found {
		return v, nil
	}

	lk, err := s.locker.Obtain(ctx, key+":lock", 10*time.Second,
		&ObtainOptions{RetryCount: 50, RetryDelay: 100 * time.Millisecond})
	if err != nil {
		// Lock unavailable (contention/timeout): try one more read, else compute
		// directly rather than block forever.
		if v, found, gerr := Get[T](ctx, s, key); gerr == nil && found {
			return v, nil
		}
		return computeAndStore(ctx, s, key, ttl, fn)
	}
	defer func() { _ = lk.Release(context.Background()) }()

	// Double-check: another caller may have filled the cache before we locked.
	if v, found, err := Get[T](ctx, s, key); err != nil {
		var zero T
		return zero, err
	} else if found {
		return v, nil
	}

	return computeAndStore(ctx, s, key, ttl, fn)
}

// RememberForever is Remember with no expiry.
func RememberForever[T any](ctx context.Context, s *CacheStore, key string, fn func() (T, error)) (T, error) {
	return Remember(ctx, s, key, 0, fn)
}

func computeAndStore[T any](ctx context.Context, s *CacheStore, key string, ttl time.Duration, fn func() (T, error)) (T, error) {
	value, err := fn()
	if err != nil {
		return value, err
	}
	if err := s.Put(ctx, key, value, ttl); err != nil {
		return value, err
	}
	return value, nil
}
