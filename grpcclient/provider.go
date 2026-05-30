package grpcclient

import (
	"sync/atomic"

	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/support"
)

// ServiceProvider registers the gRPC client pool as a singleton and closes its
// connections on shutdown.
type ServiceProvider struct {
	app contracts.Application
}

func NewServiceProvider(app contracts.Application) support.Provider {
	return &ServiceProvider{app: app}
}

func (p *ServiceProvider) Register() {
	var pool atomic.Pointer[Clients]
	p.app.Singleton(func() *Clients {
		c := NewClients()
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

func (p *ServiceProvider) Boot() {}
