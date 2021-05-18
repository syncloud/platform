package cron

import (
	"log"
	"os/exec"
)

func Job() error {
	out, err := exec.Command("snap", "run", "platform.python", "/snap/platform/current/bin/cron.py").CombinedOutput()
	log.Printf("Cron: %s", out)
	return err
}
