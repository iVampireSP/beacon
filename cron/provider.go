package cron

import (
	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/support"
)

type CronServiceProvider struct {
	app contracts.Application
}

func NewCronServiceProvider(app contracts.Application) support.Provider {
	return &CronServiceProvider{app: app}
}

func (p *CronServiceProvider) Register() {
	p.app.Singleton(NewCron)
}

func (p *CronServiceProvider) Boot() {}
