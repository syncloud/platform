package storage

import "github.com/syncloud/platform/storage/model"

type ScrubStub struct {
	status  model.ScrubStatus
	started bool
}

func (s *ScrubStub) Start() error {
	s.started = true
	return nil
}

func (s *ScrubStub) Status() (*model.ScrubStatus, error) {
	return &s.status, nil
}
