package schedule

import (
	"github.com/iVampireSP/foundation/container"
	"github.com/iVampireSP/foundation/support"
)

// ScheduleServiceProvider owns the `scheduler` command. It binds no singletons
// of its own (the scheduler resolves cron, lock and queue from their providers).
// The provider and its command live together in the module root.
type ScheduleServiceProvider struct {
	support.ServiceProvider
}

func NewScheduleServiceProvider(app *container.Application) *ScheduleServiceProvider {
	return &ScheduleServiceProvider{ServiceProvider: support.ServiceProvider{App: app}}
}

func (p *ScheduleServiceProvider) Register() {
	p.AddCommand(NewSchedulerCommand)
}

func (p *ScheduleServiceProvider) Boot() {}
