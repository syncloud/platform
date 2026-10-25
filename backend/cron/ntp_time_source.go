package cron

import (
	"time"

	"github.com/beevik/ntp"
)

type NtpTimeSource struct {
	timeout time.Duration
}

func NewNtpTimeSource(timeout time.Duration) *NtpTimeSource {
	return &NtpTimeSource{timeout: timeout}
}

func (s *NtpTimeSource) Offset(server string) (time.Duration, error) {
	response, err := ntp.QueryWithOptions(server, ntp.QueryOptions{Timeout: s.timeout})
	if err != nil {
		return 0, err
	}
	err = response.Validate()
	if err != nil {
		return 0, err
	}
	return response.ClockOffset, nil
}
