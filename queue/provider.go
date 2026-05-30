package queue

import (
	"github.com/iVampireSP/foundation/container"
	"github.com/iVampireSP/foundation/support"
)

// QueueServiceProvider wires the job queue and owns the `worker` command,
// mirroring Laravel's Illuminate\Queue\QueueServiceProvider: the provider and
// its commands live together in the module root.
type QueueServiceProvider struct {
	support.ServiceProvider
	app *container.Application
}

func NewQueueServiceProvider(app *container.Application) *QueueServiceProvider {
	return &QueueServiceProvider{ServiceProvider: support.ServiceProvider{App: app}, app: app}
}

func (p *QueueServiceProvider) Register() {
	p.app.Singleton(NewDefaultConfig)
	p.app.Singleton(NewDefaultRedisConfig)
	p.app.Singleton(NewQueue)
	p.AddCommand(NewWorker)
}

func (p *QueueServiceProvider) Boot() {}
