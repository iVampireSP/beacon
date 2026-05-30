package queue

import (
	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/support"
)

// QueueServiceProvider wires the job queue and owns the `worker` command,
// mirroring Laravel's Illuminate\Queue\QueueServiceProvider: the provider and
// its commands live together in the module root.
type QueueServiceProvider struct {
	support.ServiceProvider
	app contracts.Application
}

func NewQueueServiceProvider(app contracts.Application) support.Provider {
	return &QueueServiceProvider{ServiceProvider: support.ServiceProvider{Registry: app}, app: app}
}

func (p *QueueServiceProvider) Register() {
	p.app.Singleton(NewDefaultConfig)
	p.app.Singleton(NewDefaultRedisConfig)
	p.app.Singleton(NewQueue)
	p.Add(NewWorker)
}
