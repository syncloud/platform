package cron

import "time"

type SimpleScheduler struct {
}

func NewSimpleScheduler() *SimpleScheduler {
	return &SimpleScheduler{}
}
func (s *SimpleScheduler) ShouldRun(day int, hour int, now time.Time, last time.Time) bool {
	if now.Truncate(time.Hour) == last.Truncate(time.Hour) {
		return false
	}
	if day == 0 {
		return now.Hour() == hour
	} else {
		if s.weekDay(now) == day {
			return now.Hour() == hour
		}
		return false
	}
}

func (s *SimpleScheduler) weekDay(now time.Time) int {
	weekday := now.Weekday()
	if weekday == time.Sunday {
		return 7
	}
	return int(weekday)
}
