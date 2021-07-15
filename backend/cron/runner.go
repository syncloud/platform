package cron

import (
	"log"
	"time"
)

type Cron struct {
	job     func() error
	started bool
	delay   time.Duration
}

func New(job func() error, delay time.Duration) *Cron {
	return &Cron{job: job, delay: delay}
}

func (c *Cron) Start() {
	c.started = true
	go func() {
		for c.started {
			err := c.job()
			if err != nil {
				log.Printf("Cron job failed: %s", err)
			}
			time.Sleep(c.delay)
		}
	}()
}

func (c *Cron) Stop() {
	c.started = false
}
