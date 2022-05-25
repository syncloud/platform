package snap

import (
	"github.com/syncloud/platform/executor"
	"log"
)

type Service struct {
	executor executor.Executor
}

func NewService(executor executor.Executor) *Service {
	return &Service{
		executor: executor,
	}
}

func (s *Service) Start(name string) error {
	return s.run("start", name)
}

func (s *Service) Stop(name string) error {
	return s.run("stop", name)
}

func (s *Service) run(command string, name string) error {
	_, err := s.executor.CommandOutput("snap", command, name)
	if err != nil {
		log.Printf("snap %s failed: %s", command, err)
		return err
	}
	return nil
}
