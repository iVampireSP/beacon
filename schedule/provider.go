package schedule

import (
	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/support"
)

// ScheduleServiceProvider owns the `scheduler` command. It binds no singletons
// of its own (the scheduler resolves cron, lock and queue from their providers).
// The provider and its command live together in the module root.
type ScheduleServiceProvider struct {
	support.ServiceProvider
}

func NewScheduleServiceProvider(app contracts.Application) support.Provider {
	return &ScheduleServiceProvider{ServiceProvider: support.ServiceProvider{Registry: app}}
}

func (p *ScheduleServiceProvider) Register() {
	p.Add(NewSchedulerCommand)
}
