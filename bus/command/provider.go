package command

import (
	"github.com/iVampireSP/foundation/bus"
	"github.com/iVampireSP/foundation/console"
	"github.com/iVampireSP/foundation/container"
)

// BusServiceProvider wires the event bus and owns the `eventbus` command.
// It lives alongside the command (not in the bus core package) so it can
// declare its command without a bus -> bus/command import cycle.
type BusServiceProvider struct {
	app *container.Application
}

func NewBusServiceProvider(app *container.Application) *BusServiceProvider {
	return &BusServiceProvider{app: app}
}

func (p *BusServiceProvider) Register() {
	p.app.Singleton(bus.NewDefaultConfig)
	p.app.Singleton(bus.NewBus)
}

func (p *BusServiceProvider) Boot() {}

// Commands implements console.CommandProvider.
func (p *BusServiceProvider) Commands() []console.ConsoleCommand {
	return []console.ConsoleCommand{NewEventBus(p.app)}
}
