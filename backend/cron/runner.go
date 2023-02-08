package cron

import (
	"github.com/syncloud/platform/config"
	"log"
	"time"
)

type Cron struct {
	jobs       []Job
	delay      time.Duration
	userConfig *config.UserConfig
}

type Job interface {
	Run() error
}

func New(jobs []Job, delay time.Duration, userConfig *config.UserConfig) *Cron {
	return &Cron{jobs: jobs, delay: delay, userConfig: userConfig}
}

func (c *Cron) StartSingle() {
	if !c.userConfig.IsActivated() {
		log.Println("device is not activated yet, not running cron")
		return
	}
	for _, job := range c.jobs {
		err := job.Run()
		if err != nil {
			log.Printf("Cron job failed: %s", err)
		}
	}

}

func (c *Cron) Start() error {
	go func() {
		for {
			c.StartSingle()
			time.Sleep(c.delay)
		}
	}()
	return nil
}
