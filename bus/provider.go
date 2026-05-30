package bus

import (
	"github.com/iVampireSP/beacon/foundation"
	"github.com/iVampireSP/beacon/support"
)

// BusServiceProvider wires the event bus and owns the `eventbus` command,
// mirroring Laravel's provider layout: the provider and its commands live
// together in the module root.
type BusServiceProvider struct {
	support.ServiceProvider
	app *foundation.Application
}

func NewBusServiceProvider(app *foundation.Application) *BusServiceProvider {
	return &BusServiceProvider{ServiceProvider: support.ServiceProvider{Kernel: app}, app: app}
}

func (p *BusServiceProvider) Register() {
	p.app.Singleton(NewDefaultConfig)
	p.app.Singleton(NewBus)
	p.AddCommand(NewEventBus)
}
