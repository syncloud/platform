package cron

import "github.com/syncloud/platform/storage/model"

type StorageChecksStub struct {
	schedule       model.CheckSchedule
	scrubAvailable bool
	scrubErr       error
	scrubs         int
	smartTests     int
}

func (s *StorageChecksStub) Schedule() model.CheckSchedule {
	return s.schedule
}

func (s *StorageChecksStub) ScrubAvailable() bool {
	return s.scrubAvailable
}

func (s *StorageChecksStub) StartScrub() error {
	if s.scrubErr != nil {
		return s.scrubErr
	}
	s.scrubs++
	return nil
}

func (s *StorageChecksStub) StartSmartTests() error {
	s.smartTests++
	return nil
}
