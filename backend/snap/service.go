package snap

import (
	"log"
	"os/exec"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Start(name string) error {
	return run("start", name)
}

func (s *Service) Stop(name string) error {
	return run("stop", name)
}

func run(command string, name string) error {
	_, err := exec.Command("snap", command, name).CombinedOutput()
	if err != nil {
		log.Printf("snap %s failed: %s", command, err)
		return err
	}
	return nil
}
