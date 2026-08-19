package system

import (
	"time"

	"github.com/syncloud/platform/cli"
	"go.uber.org/zap"
)

const Delay = 2 * time.Second

type Power struct {
	executor cli.Executor
	delay    time.Duration
	logger   *zap.Logger
}

func NewPower(executor cli.Executor, logger *zap.Logger) *Power {
	return &Power{
		executor: executor,
		delay:    Delay,
		logger:   logger,
	}
}

func (p *Power) Restart() {
	p.schedule("-r", "now")
}

func (p *Power) Shutdown() {
	p.schedule("now")
}

func (p *Power) schedule(args ...string) {
	p.logger.Info("scheduling shutdown", zap.Strings("args", args), zap.Duration("delay", p.delay))
	go func() {
		time.Sleep(p.delay)
		output, err := p.executor.CombinedOutput("shutdown", args...)
		if err != nil {
			p.logger.Error("shutdown failed", zap.Error(err))
			return
		}
		p.logger.Info("shutdown", zap.String("output", string(output)))
	}()
}
