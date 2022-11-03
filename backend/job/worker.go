package job

import (
	"go.uber.org/zap"
	"time"
)

type Master interface {
	Take() (func() error, error)
	Complete() error
}

type Worker struct {
	master Master
	logger *zap.Logger
}

func NewWorker(master Master, logger *zap.Logger) *Worker {
	return &Worker{master, logger}
}

func (w *Worker) Start() {
	for {
		if !w.Do() {
			time.Sleep(1 * time.Second)
		}
	}
}

func (w *Worker) Do() bool {
	job, err := w.master.Take()
	if err != nil {
		w.logger.Error("cannot take task", zap.Error(err))
		return false
	}
	err = job()
	if err != nil {
		w.logger.Error("error in the task", zap.Error(err))
	}
	err = w.master.Complete()
	if err != nil {
		w.logger.Error("cannot complete task", zap.Error(err))
	}
	return true
}
