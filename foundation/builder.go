package foundation

import (
	"io/fs"

	"github.com/iVampireSP/beacon/config"
	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/foundation/bootstrap"
	fconsole "github.com/iVampireSP/beacon/foundation/console"
)

// Builder configures and creates an Application — the analog of Laravel's
// Illuminate\Foundation\Configuration\ApplicationBuilder, used as
// foundation.Configure().WithConfig(...).WithProviders(...).Create(), mirroring
// Application::configure(...)->withProviders(...)->create().
//
// Like Laravel's builder it configures BEHAVIOUR — providers and commands — not
// resource locations, with ONE exception: WithConfig. Config is privileged in
// every faithful framework (Laravel's LoadConfiguration bootstrapper; goravel's
// config base provider) because it must load before any provider reads it, so the
// builder takes the config source and registers it as the first provider. Every
// OTHER resource (translations, templates, migrations) is an ordinary provider in
// WithProviders — i18n.NewI18nServiceProvider(fs,dir), etc. — so foundation
// imports only config, not i18n/tmpl/db.
type Builder struct {
	providers    []contracts.ProviderConstructor
	commandCtors []any

	configFS  fs.FS
	configDir string
	hasConfig bool
}

// Configure starts building an application. The CLI root identity (use, short)
// is supplied later to HandleCommand.
func Configure() *Builder {
	return &Builder{}
}

// WithConfig sets the embedded filesystem the configuration is loaded from. It is
// registered as a PRIVILEGED provider — first, before the framework's default
// providers — so config is populated before anything reads it (the goravel model;
// Laravel's LoadConfiguration). This is the only resource the builder knows about
// directly, because only config must precede the providers.
func (b *Builder) WithConfig(fsys fs.FS, dir string) *Builder {
	b.configFS, b.configDir, b.hasConfig = fsys, dir, true
	return b
}

// WithProviders sets the application's service-provider list — the analog of
// Laravel's bootstrap/providers.php `return [...]`. Each entry is a constructor
// the RegisterProviders bootstrapper calls with the application.
func (b *Builder) WithProviders(providers []contracts.ProviderConstructor) *Builder {
	b.providers = providers
	return b
}

// WithCommands registers console-command constructors directly on the kernel, in
// addition to any a provider pushes via Add.
func (b *Builder) WithCommands(constructors ...any) *Builder {
	b.commandCtors = append(b.commandCtors, constructors...)
	return b
}

// Create assembles the application: it wires the console kernel, records the
// builder's commands, and builds the ordered bootstrapper list the application
// runs on HandleCommand. Providers are NOT registered here — RegisterProviders
// runs during bootstrap, mirroring Laravel.
func (b *Builder) Create() *Application {
	app := newApplication()
	app.consoleKernel = fconsole.NewKernel(app.Container)
	app.Add(b.commandCtors...)
	app.bootstrappers = b.bootstrappers()
	return app
}

// bootstrappers assembles the ordered boot sequence: register providers → boot
// providers (the fixed framework steps, like Laravel's Kernel::$bootstrappers).
// The provider list is ordered config (privileged, if set) → framework defaults →
// the app's own, so config is registered and loaded before any provider reads it.
func (b *Builder) bootstrappers() []bootstrap.Bootstrapper {
	var providers []contracts.ProviderConstructor
	if b.hasConfig {
		providers = append(providers, config.NewConfigServiceProvider(b.configFS, b.configDir))
	}
	// Built-in framework providers next, then the app's own — mirroring how
	// Laravel merges ServiceProvider::defaultProviders() with the app's list.
	providers = append(providers, DefaultProviders...)
	providers = append(providers, b.providers...)
	return []bootstrap.Bootstrapper{
		bootstrap.RegisterProviders{Providers: providers},
		bootstrap.BootProviders{},
	}
}
