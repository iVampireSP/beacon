package command

import (
	"github.com/iVampireSP/foundation/console"
	"github.com/iVampireSP/foundation/container"
	jobqueue "github.com/iVampireSP/foundation/queue"
)

// QueueServiceProvider wires the job queue and owns the `worker` command.
// It lives alongside the command (rather than in the queue core package) so it
// can declare its command without a queue -> queue/command import cycle.
type QueueServiceProvider struct {
	app *container.Application
}

func NewQueueServiceProvider(app *container.Application) *QueueServiceProvider {
	return &QueueServiceProvider{app: app}
}

func (p *QueueServiceProvider) Register() {
	p.app.Singleton(jobqueue.NewDefaultConfig)
	p.app.Singleton(jobqueue.NewDefaultRedisConfig)
	p.app.Singleton(jobqueue.NewQueue)
}

func (p *QueueServiceProvider) Boot() {}

// Commands implements console.CommandProvider.
func (p *QueueServiceProvider) Commands() []console.ConsoleCommand {
	return []console.ConsoleCommand{NewWorker(p.app)}
}
