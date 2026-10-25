package cron

import (
	"errors"
	"fmt"
	"time"

	"github.com/syncloud/platform/date"
	"go.uber.org/zap"
)

const TimeSyncMaxOffset = time.Minute

var TimeSyncServers = []string{"0.pool.ntp.org", "1.pool.ntp.org", "2.pool.ntp.org"}

type TimeSource interface {
	Offset(server string) (time.Duration, error)
}

type Clock interface {
	Set(now time.Time) error
}

type TimeSyncJob struct {
	source       TimeSource
	clock        Clock
	servers      []string
	dateProvider date.Provider
	logger       *zap.Logger
}

func NewTimeSyncJob(source TimeSource, clock Clock, servers []string, dateProvider date.Provider, logger *zap.Logger) *TimeSyncJob {
	return &TimeSyncJob{
		source:       source,
		clock:        clock,
		servers:      servers,
		dateProvider: dateProvider,
		logger:       logger,
	}
}

func (t *TimeSyncJob) Run() error {
	var errs []error
	for _, server := range t.servers {
		offset, err := t.source.Offset(server)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", server, err))
			continue
		}
		if offset.Abs() <= TimeSyncMaxOffset {
			return nil
		}
		now := t.dateProvider.Now().Add(offset)
		t.logger.Info("setting system clock", zap.String("server", server), zap.Duration("offset", offset), zap.Time("time", now))
		return t.clock.Set(now)
	}
	return fmt.Errorf("time sync failed: %w", errors.Join(errs...))
}
