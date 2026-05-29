package tracing

import (
	"time"

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
	p.app.Singleton(NewTracing)
}

func (p *ServiceProvider) Boot() {}

// NewDefaultConfig returns a Tracing Config populated from the application config.
func NewDefaultConfig() Config {
	return Config{
		Enabled:       config.Bool("tracing.enabled", false),
		Endpoint:      config.String("tracing.endpoint", "localhost:4317"),
		SampleRatio:   config.GetFloat("tracing.sample_ratio", 1.0),
		QueryURL:      config.String("tracing.query_url", "http://jaeger-query:16686"),
		QueryUsername: config.String("tracing.query_username"),
		QueryPassword: config.String("tracing.query_password"),
		QueryTimeout:  config.Duration("tracing.query_timeout", 15*time.Second),
	}
}
