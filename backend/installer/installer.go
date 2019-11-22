package installer

import (
	"log"
	"os/exec"
)

type Installer struct {
}

const (
	INSTALLER_UPGRADE_CMD = "/snap/platform/current/bin/upgrade-snapd.sh"
)

func New() *Installer {
	return &Installer{}
}

func (installer *Installer) Upgrade() {
	log.Println("Running installer upgrade", INSTALLER_UPGRADE_CMD)
	out, err := exec.Command(INSTALLER_UPGRADE_CMD, "stable").CombinedOutput()
	log.Printf("Installer upgrade output %s", out)
	if err != nil {
		log.Printf("Installer upgrade failed: %v", err)
	} else {
		log.Printf("Installer upgrade completed")
	}
}
