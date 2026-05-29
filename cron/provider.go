package cron

import "github.com/iVampireSP/foundation/container"

type ServiceProvider struct {
	app *container.Application
}

func NewServiceProvider(app *container.Application) *ServiceProvider {
	return &ServiceProvider{app: app}
}

func (p *ServiceProvider) Register() {
	p.app.Singleton(NewCron)
}

func (p *ServiceProvider) Boot() {}
