package support

// Provider is the service provider contract: the application registers and
// boots these. Register() binds services into the container; Boot() runs after
// all providers have been registered. Mirrors the registrable role of Laravel's
// Illuminate\Support\ServiceProvider.
type Provider interface {
	Register()
	Boot()
}

// Application is the slice of the container that a ServiceProvider needs. It is
// declared here, rather than imported from the container package, to invert the
// dependency (container imports support, never the reverse) — mirroring how
// Laravel types $app as the Contracts\Foundation\Application interface.
// *container.Application satisfies it.
type Application interface {
	RegisterCommands(constructors ...any)
}

// ServiceProvider is the base that application providers embed, mirroring
// Illuminate\Support\ServiceProvider. It holds the application and registers
// console commands. Concrete providers embed it and override Register/Boot.
type ServiceProvider struct {
	App Application
}

// Register binds services into the container. The embedded default is a no-op;
// concrete providers override it.
func (p *ServiceProvider) Register() {}

// Boot runs after all providers are registered. The embedded default is a
// no-op; concrete providers override it.
func (p *ServiceProvider) Boot() {}

// AddCommand registers console commands by their constructors. The container
// resolves each constructor's dependencies (constructor injection) when the
// command is built. This is the analog of Laravel's ServiceProvider::commands(),
// named AddCommand to avoid ambiguity with the command-collection capability.
func (p *ServiceProvider) AddCommand(commands ...any) {
	p.App.RegisterCommands(commands...)
}
