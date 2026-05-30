package bus

import (
	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/support"
)

// BusServiceProvider wires the event bus and owns the `eventbus` command,
// mirroring Laravel's provider layout: the provider and its commands live
// together in the module root.
type BusServiceProvider struct {
	support.ServiceProvider
	app contracts.Application
}

func NewBusServiceProvider(app contracts.Application) support.Provider {
	return &BusServiceProvider{ServiceProvider: support.ServiceProvider{Registry: app}, app: app}
}

func (p *BusServiceProvider) Register() {
	p.app.Singleton(NewDefaultConfig)
	p.app.Singleton(NewBus)
	p.Add(NewEventBus)
}
