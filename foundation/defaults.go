package foundation

import (
	"github.com/iVampireSP/beacon/bus"
	"github.com/iVampireSP/beacon/cache"
	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/cron"
	"github.com/iVampireSP/beacon/logger"
	"github.com/iVampireSP/beacon/queue"
	"github.com/iVampireSP/beacon/schedule"
	"github.com/iVampireSP/beacon/tracing"
)

// DefaultProviders are the framework's built-in infrastructure service providers,
// registered automatically (before the app's own) by the builder — the analog of
// Laravel's ServiceProvider::defaultProviders(): the Illuminate providers every
// app gets for free. Each binds lazily, so an app that never runs the worker or
// event bus pays nothing. The database/ORM provider is deliberately NOT here —
// an app brings its own (ent-based) ORM provider.
//
// foundation can list these because they all depend only on contracts (not on
// foundation), so there is no import cycle.
var DefaultProviders = []contracts.ProviderConstructor{
	logger.NewServiceProvider,
	cache.NewServiceProvider, // also provides the distributed Locker
	cron.NewServiceProvider,
	tracing.NewServiceProvider,
	queue.NewQueueServiceProvider,
	bus.NewBusServiceProvider,
	schedule.NewScheduleServiceProvider,
}
