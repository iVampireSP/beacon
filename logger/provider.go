package logger

import (
	"github.com/iVampireSP/beacon/config"
	"github.com/iVampireSP/beacon/foundation"
)

type ServiceProvider struct {
	app *foundation.Application
}

func NewServiceProvider(app *foundation.Application) *ServiceProvider {
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
