package scheduler

import (
	"github.com/go-co-op/gocron"
	"time"
)

type Master struct {
	cronScheduler *gocron.Scheduler
	providers     map[string]IProvider
}

func (m *Master) Start() {
	for _, p := range m.providers {
		m.cronScheduler.CronWithSeconds(p.GetCronExpression()).Do(p.Run)
	}

	m.cronScheduler.StartAsync()
}

func (m *Master) Stop() {
	m.cronScheduler.Stop()
}

func NewMaster(providers []IProvider) *Master {
	master := &Master{}
	master.cronScheduler = gocron.NewScheduler(time.Local)
	master.providers = make(map[string]IProvider, len(providers))

	for _, p := range providers {
		master.providers[p.GetName()] = p
	}

	return master
}
