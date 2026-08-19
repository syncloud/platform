package cron

import (
	"go.uber.org/zap"
)

type SnapshotsCleaner interface {
	DisableAutomatic() error
	ForgetAuto() error
}

type SnapshotsJob struct {
	cleaner SnapshotsCleaner
	logger  *zap.Logger
}

func NewSnapshotsJob(cleaner SnapshotsCleaner, logger *zap.Logger) *SnapshotsJob {
	return &SnapshotsJob{
		cleaner: cleaner,
		logger:  logger,
	}
}

func (j *SnapshotsJob) Run() error {
	err := j.cleaner.DisableAutomatic()
	if err != nil {
		return err
	}
	return j.cleaner.ForgetAuto()
}
