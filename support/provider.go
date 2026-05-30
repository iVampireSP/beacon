package support

// Provider is the service provider contract: the application registers and
// boots these. Register() binds services into the container; Boot() runs after
// all providers have been registered. Mirrors the registrable role of Laravel's
// Illuminate\Support\ServiceProvider.
type Provider interface {
	Register()
	Boot()
}

// Registry is the application's single contribution sink. A ServiceProvider
// pushes ALL its contributions — console commands, queue jobs, event listeners,
// cron jobs — through the one Add method; there is no per-kind registration. The
// items are `any` so this stays a pure leaf: each runtime claims what it
// understands (the console kernel builds the command constructors; the
// worker/eventbus/scheduler type-assert queue/job.Handler, bus.Listener,
// schedule.CronJob). *foundation.Application satisfies it.
type Registry interface {
	Add(contributions ...any)
}

// ServiceProvider is the embeddable base, mirroring Illuminate\Support\ServiceProvider.
// It holds the Registry so a provider declares everything with one Add — the
// generalized analog of Laravel's $this->commands([...]). Concrete providers
// embed it and override Register/Boot as needed.
type ServiceProvider struct {
	Registry Registry
}

// Register binds services into the container. The embedded default is a no-op;
// concrete providers override it.
func (ServiceProvider) Register() {}

// Boot runs after all providers are registered. The embedded default is a
// no-op; concrete providers override it.
func (ServiceProvider) Boot() {}

// Add declares the provider's contributions — command constructors, job
// handlers, event listeners, cron jobs — through the single registry. Each
// runtime picks up the kinds it understands.
func (p ServiceProvider) Add(contributions ...any) {
	p.Registry.Add(contributions...)
}
