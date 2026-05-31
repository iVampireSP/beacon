package grpcclient

import (
	"sync/atomic"

	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/support"
)

// GRPCClientServiceProvider registers the gRPC client pool as a singleton and closes its
// connections on shutdown.
type GRPCClientServiceProvider struct {
	app contracts.Application
}

func NewGRPCClientServiceProvider(app contracts.Application) support.Provider {
	return &GRPCClientServiceProvider{app: app}
}

func (p *GRPCClientServiceProvider) Register() {
	var pool atomic.Pointer[GRPCClientPool]
	p.app.Singleton(func() *GRPCClientPool {
		c := NewGRPCClientPool()
		pool.Store(c)
		return c
	})
	p.app.OnShutdown(func() error {
		if c := pool.Load(); c != nil {
			return c.Close()
		}
		return nil
	})
}

func (p *GRPCClientServiceProvider) Boot() {}
