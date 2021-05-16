package cron

import "time"

type Cron struct {
	job     func()
	started bool
	delay   time.Duration
}

func New(job func(), delay time.Duration) *Cron {
	return &Cron{job: job, delay: delay}
}

func (c *Cron) Start() {
	c.started = true
	go func() {
		for c.started {
			c.job()
			time.Sleep(c.delay)
		}
	}()
}

func (c *Cron) Stop() {
	c.started = false
}
