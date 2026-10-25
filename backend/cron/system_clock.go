package cron

import (
	"time"

	"golang.org/x/sys/unix"
)

type SystemClock struct {
}

func NewSystemClock() *SystemClock {
	return &SystemClock{}
}

func (c *SystemClock) Set(now time.Time) error {
	tv := unix.NsecToTimeval(now.UnixNano())
	return unix.Settimeofday(&tv)
}
