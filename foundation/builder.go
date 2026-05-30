package foundation

import (
	"io/fs"

	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/foundation/bootstrap"
	fconsole "github.com/iVampireSP/beacon/foundation/console"
	"github.com/iVampireSP/beacon/support"
)

// Builder configures and creates an Application — the analog of Laravel's
// Illuminate\Foundation\Configuration\ApplicationBuilder, used as
// foundation.Configure().WithConfig(...).WithProviders(...).Create(), mirroring
// Application::configure(...)->withProviders(...)->create(). The configured
// resources and providers become the ordered bootstrappers the application runs
// on HandleCommand.
type Builder struct {
	providerFactory func(contracts.Application) []support.Provider
	commandCtors    []any

	configFS  fs.FS
	configDir string
	langFS    fs.FS
	langDir   string

	templatesFS  fs.FS
	templatesDir string
	hasTemplates bool
}

// Configure starts building an application. The CLI root identity (use, short)
// is supplied later to HandleCommand.
func Configure() *Builder {
	return &Builder{}
}

// WithConfig sets the embedded filesystem the LoadConfiguration bootstrapper
// loads application config from.
func (b *Builder) WithConfig(fsys fs.FS, dir string) *Builder {
	b.configFS, b.configDir = fsys, dir
	return b
}

// WithLocale sets the embedded filesystem the LoadTranslations bootstrapper loads
// i18n catalogs from.
func (b *Builder) WithLocale(fsys fs.FS, dir string) *Builder {
	b.langFS, b.langDir = fsys, dir
	return b
}

// WithTemplates sets the filesystem and subdir the LoadTemplates bootstrapper
// parses templates from.
func (b *Builder) WithTemplates(fsys fs.FS, dir string) *Builder {
	b.templatesFS, b.templatesDir, b.hasTemplates = fsys, dir, true
	return b
}

// WithProviders sets the factory that builds the service providers. The factory
// receives the application so each provider can take it (new Provider($app)).
func (b *Builder) WithProviders(factory func(contracts.Application) []support.Provider) *Builder {
	b.providerFactory = factory
	return b
}

// WithCommands registers console-command constructors directly on the kernel, in
// addition to any a provider pushes via AddCommand.
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
	app.RegisterCommands(b.commandCtors...)
	app.bootstrappers = b.bootstrappers()
	return app
}

// bootstrappers assembles the ordered boot sequence from the configured sources:
// load config → load translations → load templates → register providers → boot
// providers. Only the resource loaders that were configured are included.
func (b *Builder) bootstrappers() []bootstrap.Bootstrapper {
	var list []bootstrap.Bootstrapper
	if b.configFS != nil {
		list = append(list, bootstrap.LoadConfiguration{FS: b.configFS, Dir: b.configDir})
	}
	if b.langFS != nil {
		list = append(list, bootstrap.LoadTranslations{FS: b.langFS, Dir: b.langDir})
	}
	if b.hasTemplates {
		list = append(list, bootstrap.LoadTemplates{FS: b.templatesFS, Dir: b.templatesDir})
	}
	list = append(list,
		bootstrap.RegisterProviders{Factory: b.providerFactory},
		bootstrap.BootProviders{},
	)
	return list
}
