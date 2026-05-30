package email

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
	p.app.Singleton(NewEmail)
}

func (p *ServiceProvider) Boot() {}

// NewDefaultConfig returns an email Config populated from the application config.
func NewDefaultConfig() Config {
	return Config{
		Host:       config.String("mail.host"),
		Port:       config.Int("mail.port", 587),
		Username:   config.String("mail.username"),
		Password:   config.String("mail.password"),
		From:       config.String("mail.from"),
		Encryption: config.String("mail.encryption", "tls"),
	}
}
