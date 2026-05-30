package httpserver

import (
	"time"

	"github.com/iVampireSP/beacon/config"
	"github.com/iVampireSP/beacon/container"
)

type ServiceProvider struct {
	app *container.Application
}

func NewServiceProvider(app *container.Application) *ServiceProvider {
	return &ServiceProvider{app: app}
}

func (p *ServiceProvider) Register() {
	p.app.Singleton(NewDefaultMetricsConfig)
}

func (p *ServiceProvider) Boot() {}

func NewDefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		Enabled:         config.Bool("metrics.enabled", true),
		Host:            config.String("metrics.host", "0.0.0.0"),
		Port:            config.Int("metrics.port", 9090),
		ShutdownTimeout: 30 * time.Second,
	}
}
