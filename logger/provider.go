package logger

import (
	"github.com/iVampireSP/beacon/config"
	"github.com/iVampireSP/beacon/contracts"
)

type ServiceProvider struct {
	app contracts.Application
}

func NewServiceProvider(app contracts.Application) *ServiceProvider {
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
