package logger

import (
	"github.com/iVampireSP/foundation/config"
	"github.com/iVampireSP/foundation/container"
)

type ServiceProvider struct {
	app *container.Application
}

func NewServiceProvider(app *container.Application) *ServiceProvider {
	return &ServiceProvider{app: app}
}

func (p *ServiceProvider) Register() {
	p.app.Singleton(NewDefaultConfig)
	p.app.Singleton(New)
}

func (p *ServiceProvider) Boot() {}

// NewDefaultConfig returns a logger Config populated from the application config.
func NewDefaultConfig() Config {
	return Config{
		Level: config.String("log.level", "info"),
		Debug: config.Bool("app.debug", false),
	}
}
