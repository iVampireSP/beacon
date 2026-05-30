package httpserver

import (
	"time"

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
