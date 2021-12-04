package cron

import (
	"github.com/syncloud/platform/config"
	"log"
	"os/exec"
)

type PortsJob struct {
	userConfig *config.UserConfig
}

func NewPortsJob(userConfig *config.UserConfig) *PortsJob {
	return &PortsJob{
		userConfig: userConfig,
	}
}

func (j *PortsJob) Run() error {
	if j.userConfig.IsRedirectEnabled() {
		out, err := exec.Command("snap", "run", "platform.python", "/snap/platform/current/bin/ports_job.py").CombinedOutput()
		log.Printf("Cron: %s", out)
		return err
	}
	return nil
}
