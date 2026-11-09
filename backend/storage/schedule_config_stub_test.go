package storage

import "github.com/syncloud/platform/storage/model"

type ScheduleConfigStub struct {
	schedule model.CheckSchedule
}

func (c *ScheduleConfigStub) StorageCheckSchedule() model.CheckSchedule {
	return c.schedule
}

func (c *ScheduleConfigStub) SetStorageCheckSchedule(schedule model.CheckSchedule) {
	c.schedule = schedule
}
