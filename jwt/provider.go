package jwt

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
	p.app.Singleton(NewJWT)
}

func (p *ServiceProvider) Boot() {}

// NewDefaultConfig returns a JWT Config populated from the application config.
func NewDefaultConfig() Config {
	return Config{
		KeyName:   config.String("jwt.key", "rsa"),
		Issuer:    config.String("discovery.issuer", "http://localhost"),
		ExpiresIn: config.Int("jwt.expires_in", 86400),
	}
}
