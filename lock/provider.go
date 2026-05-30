package lock

import "github.com/iVampireSP/beacon/foundation"

type ServiceProvider struct {
	app *foundation.Application
}

func NewServiceProvider(app *foundation.Application) *ServiceProvider {
	return &ServiceProvider{app: app}
}

func (p *ServiceProvider) Register() {
	p.app.Singleton(NewLocker)
}

func (p *ServiceProvider) Boot() {}
