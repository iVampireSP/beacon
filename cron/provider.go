package cron

import "github.com/iVampireSP/beacon/foundation"

type ServiceProvider struct {
	app *foundation.Application
}

func NewServiceProvider(app *foundation.Application) *ServiceProvider {
	return &ServiceProvider{app: app}
}

func (p *ServiceProvider) Register() {
	p.app.Singleton(NewCron)
}

func (p *ServiceProvider) Boot() {}
