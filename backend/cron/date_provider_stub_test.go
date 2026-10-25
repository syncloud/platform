package cron

import "time"

type DateProviderStub struct {
	now time.Time
}

func (d *DateProviderStub) Now() time.Time {
	return d.now
}
