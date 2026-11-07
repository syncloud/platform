package storage

import (
	"fmt"
	"github.com/syncloud/platform/storage/model"
)

type SmartStub struct {
	health  map[string]model.SmartHealth
	fail    map[string]bool
	started []string
}

func (s *SmartStub) Health(device string) model.SmartHealth {
	return s.health[device]
}

func (s *SmartStub) StartShortTest(device string) error {
	if s.fail[device] {
		return fmt.Errorf("test refused")
	}
	s.started = append(s.started, device)
	return nil
}
