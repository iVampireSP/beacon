package cron

import (
	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/support"
)

type ServiceProvider struct {
	app contracts.Application
}

func NewServiceProvider(app contracts.Application) support.Provider {
	return &ServiceProvider{app: app}
}

func (p *ServiceProvider) Register() {
	p.app.Singleton(NewCron)
}

func (p *ServiceProvider) Boot() {}
