package lock

import "github.com/iVampireSP/beacon/contracts"

type ServiceProvider struct {
	app contracts.Application
}

func NewServiceProvider(app contracts.Application) *ServiceProvider {
	return &ServiceProvider{app: app}
}

func (p *ServiceProvider) Register() {
	p.app.Singleton(NewLocker)
}

func (p *ServiceProvider) Boot() {}
