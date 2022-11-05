package snap

import (
	"github.com/syncloud/platform/cli"
	"go.uber.org/zap"
)

type Cli struct {
	executor cli.CommandExecutor
	logger   *zap.Logger
}

func NewCli(executor cli.CommandExecutor) *Cli {
	return &Cli{
		executor: executor,
	}
}

func (s *Cli) Start(name string) error {
	return s.run("start", name)
}

func (s *Cli) Stop(name string) error {
	return s.run("stop", name)
}
func (s *Cli) Run(name string) error {
	return s.run("run", name)
}

func (s *Cli) run(command string, name string) error {
	_, err := s.executor.CommandOutput("snap", command, name)
	if err != nil {
		s.logger.Error("snap failed", zap.String("command", command), zap.Error(err))
		return err
	}
	return nil
}
