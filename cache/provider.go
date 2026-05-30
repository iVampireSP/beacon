package cache

import (
	"sync/atomic"

	"github.com/iVampireSP/beacon/config"
	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/support"
	"github.com/redis/go-redis/v9"
)

type ServiceProvider struct {
	app contracts.Application
}

func NewServiceProvider(app contracts.Application) support.Provider {
	return &ServiceProvider{app: app}
}

func (p *ServiceProvider) Register() {
	p.app.Singleton(NewDefaultRedisConfig)

	// Capture the client when it is actually constructed and close it on
	// shutdown only if so — registered during Register for a controlled shutdown
	// order, and avoiding a forced build of an unused client.
	var client atomic.Value // redis.UniversalClient
	p.app.Singleton(func(cfg RedisConfig) redis.UniversalClient {
		c := NewCache(cfg)
		client.Store(c)
		return c
	})
	p.app.OnShutdown(func() error {
		if v := client.Load(); v != nil {
			return v.(redis.UniversalClient).Close()
		}
		return nil
	})
}

func (p *ServiceProvider) Boot() {}

// NewDefaultRedisConfig returns a RedisConfig populated from the application config.
func NewDefaultRedisConfig() RedisConfig {
	return RedisConfig{
		Host:         config.String("redis.host", "localhost"),
		Port:         config.Int("redis.port", 6379),
		Password:     config.String("redis.password"),
		DB:           config.Int("redis.db", 0),
		ClusterAddrs: config.String("redis.cluster_addrs"),
	}
}
