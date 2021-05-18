package cron

import (
	"log"
	"os/exec"
)

func Job() {
	out, err := exec.Command("snap", "run", "platform.python", "/snap/platform/current/bin/cron.py").CombinedOutput()
	log.Printf("Cron: %s", out)
	if err != nil {
		log.Printf("Cron failed: %s", err)
	}
}
