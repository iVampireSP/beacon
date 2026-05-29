package command

import (
	"github.com/iVampireSP/foundation/console"
	"github.com/iVampireSP/foundation/container"
)

// ScheduleServiceProvider owns the `scheduler` command. It binds no singletons
// of its own (the scheduler resolves cron, lock and queue from their providers);
// it lives alongside the command to avoid a schedule -> schedule/command cycle.
type ScheduleServiceProvider struct {
	app *container.Application
}

func NewScheduleServiceProvider(app *container.Application) *ScheduleServiceProvider {
	return &ScheduleServiceProvider{app: app}
}

func (p *ScheduleServiceProvider) Register() {}

func (p *ScheduleServiceProvider) Boot() {}

// Commands implements console.CommandProvider.
func (p *ScheduleServiceProvider) Commands() []console.ConsoleCommand {
	return []console.ConsoleCommand{NewScheduler(p.app)}
}
