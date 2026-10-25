package cron

import "time"

type ClockStub struct {
	set []time.Time
}

func (c *ClockStub) Set(now time.Time) error {
	c.set = append(c.set, now)
	return nil
}
