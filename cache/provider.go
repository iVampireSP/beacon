package cache

import (
	"github.com/iVampireSP/foundation/config"
	"github.com/iVampireSP/foundation/container"
	"github.com/redis/go-redis/v9"
)

type ServiceProvider struct {
	app *container.Application
}

func NewServiceProvider(app *container.Application) *ServiceProvider {
	return &ServiceProvider{app: app}
}

func (p *ServiceProvider) Register() {
	p.app.Singleton(NewDefaultRedisConfig)
	p.app.Singleton(NewCache)
	p.app.OnShutdown(func() error {
		if c, err := container.Make[redis.UniversalClient](p.app.Container); err == nil && c != nil {
			return c.Close()
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
