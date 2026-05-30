package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Store is the cache facade over Redis — the analog of Laravel's cache store,
// with distributed locking folded in (Cache::lock). Values are JSON-encoded.
// Every operation is single-key, so go-redis routes it to the right node under
// Redis Cluster — except Flush, which fans out to all masters. No multi-key
// command is issued, so there are no CROSSSLOT errors.
type Store struct {
	client redis.UniversalClient
	locker *Locker
}

// NewStore builds a cache Store (and its Locker) over the Redis client.
func NewStore(client redis.UniversalClient) *Store {
	return &Store{client: client, locker: NewLocker(client)}
}

// Client exposes the underlying Redis client for raw/advanced use.
func (s *Store) Client() redis.UniversalClient { return s.client }

// Locker exposes the distributed locker.
func (s *Store) Locker() *Locker { return s.locker }

// Put stores value (JSON-encoded) under key with ttl (0 = no expiry).
func (s *Store) Put(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, data, ttl).Err()
}

// Forever stores value with no expiry.
func (s *Store) Forever(ctx context.Context, key string, value any) error {
	return s.Put(ctx, key, value, 0)
}

// Add stores value only if key does not exist (SET NX); reports whether stored.
func (s *Store) Add(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	return s.client.SetNX(ctx, key, data, ttl).Result()
}

// Has reports whether key exists.
func (s *Store) Has(ctx context.Context, key string) (bool, error) {
	n, err := s.client.Exists(ctx, key).Result()
	return n > 0, err
}

// Forget deletes key.
func (s *Store) Forget(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

// Increment atomically increments key by n (default 1), returning the new value.
// JSON-encoded integers ("42") are valid INCR operands, so a value Put as an int
// can be incremented.
func (s *Store) Increment(ctx context.Context, key string, n ...int64) (int64, error) {
	return s.client.IncrBy(ctx, key, delta(n)).Result()
}

// Decrement atomically decrements key by n (default 1), returning the new value.
func (s *Store) Decrement(ctx context.Context, key string, n ...int64) (int64, error) {
	return s.client.DecrBy(ctx, key, delta(n)).Result()
}

func delta(n []int64) int64 {
	if len(n) > 0 {
		return n[0]
	}
	return 1
}

// Flush removes all keys. Under Redis Cluster, FLUSHDB only clears the node it is
// sent to, so fan out to every master; otherwise flush the single DB.
func (s *Store) Flush(ctx context.Context) error {
	if cc, ok := s.client.(*redis.ClusterClient); ok {
		return cc.ForEachMaster(ctx, func(ctx context.Context, m *redis.Client) error {
			return m.FlushDB(ctx).Err()
		})
	}
	return s.client.FlushDB(ctx).Err()
}

// Lock obtains a distributed lock named name with ttl — the analog of Laravel's
// Cache::lock(). Pass opts to retry/block; nil tries once.
func (s *Store) Lock(ctx context.Context, name string, ttl time.Duration, opts *ObtainOptions) (*Lock, error) {
	return s.locker.Obtain(ctx, name, ttl, opts)
}
