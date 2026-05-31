package email

import (
	"github.com/iVampireSP/beacon/config"
	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/support"
)

type EmailServiceProvider struct {
	app contracts.Application
}

func NewEmailServiceProvider(app contracts.Application) support.Provider {
	return &EmailServiceProvider{app: app}
}

func (p *EmailServiceProvider) Register() {
	p.app.Singleton(NewDefaultConfig)
	p.app.Singleton(NewEmail)
}

func (p *EmailServiceProvider) Boot() {}

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
