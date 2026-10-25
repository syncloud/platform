package cron

import (
	"fmt"
	"time"
)

type TimeSourceStub struct {
	offsets map[string]time.Duration
	asked   []string
}

func (s *TimeSourceStub) Offset(server string) (time.Duration, error) {
	s.asked = append(s.asked, server)
	offset, ok := s.offsets[server]
	if !ok {
		return 0, fmt.Errorf("i/o timeout")
	}
	return offset, nil
}
