package bus

import (
	"github.com/iVampireSP/beacon/container"
	"github.com/iVampireSP/beacon/support"
)

// BusServiceProvider wires the event bus and owns the `eventbus` command,
// mirroring Laravel's provider layout: the provider and its commands live
// together in the module root.
type BusServiceProvider struct {
	support.ServiceProvider
	app *container.Application
}

func NewBusServiceProvider(app *container.Application, kernel support.Kernel) *BusServiceProvider {
	return &BusServiceProvider{ServiceProvider: support.ServiceProvider{Kernel: kernel}, app: app}
}

func (p *BusServiceProvider) Register() {
	p.app.Singleton(NewDefaultConfig)
	p.app.Singleton(NewBus)
	p.AddCommand(NewEventBus)
}
