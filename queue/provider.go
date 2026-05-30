package queue

import (
	"github.com/iVampireSP/beacon/foundation"
	"github.com/iVampireSP/beacon/support"
)

// QueueServiceProvider wires the job queue and owns the `worker` command,
// mirroring Laravel's Illuminate\Queue\QueueServiceProvider: the provider and
// its commands live together in the module root.
type QueueServiceProvider struct {
	support.ServiceProvider
	app *foundation.Application
}

func NewQueueServiceProvider(app *foundation.Application) *QueueServiceProvider {
	return &QueueServiceProvider{ServiceProvider: support.ServiceProvider{Kernel: app}, app: app}
}

func (p *QueueServiceProvider) Register() {
	p.app.Singleton(NewDefaultConfig)
	p.app.Singleton(NewDefaultRedisConfig)
	p.app.Singleton(NewQueue)
	p.AddCommand(NewWorker)
}
