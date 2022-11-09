package snap

import (
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/snap/model"
	"go.uber.org/zap"
)

type Cli struct {
	executor cli.Executor
	logger   *zap.Logger
}

func NewCli(executor cli.Executor, logger   *zap.Logger) *Cli {
	return &Cli{
		executor: executor,
 logger: logger,
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

func (s *Cli) RunCmdIfExists(snap model.Snap, name string) error {
	cmd := snap.FindCommand(name)
	if cmd != nil {
		err := s.Run(cmd.FullName())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Cli) run(command string, name string) error {
	_, err := s.executor.CombinedOutput("snap", command, name)
	if err != nil {
		s.logger.Error("snap failed", zap.String("command", command), zap.Error(err))
		return err
	}
	return nil
}
