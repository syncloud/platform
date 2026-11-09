package cron

import "time"

type StorageCheckConfigStub struct {
	times map[string]time.Time
}

func (c *StorageCheckConfigStub) GetStorageCheckTime(check string) time.Time {
	return c.times[check]
}

func (c *StorageCheckConfigStub) SetStorageCheckTime(check string, time time.Time) {
	c.times[check] = time
}
