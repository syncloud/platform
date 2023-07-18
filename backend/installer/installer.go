package installer

import (
	"log"
	"os/exec"
)

type Installer struct {
}

const (
	UpgradeCmd = "/snap/platform/current/bin/upgrade-snapd.sh"
)

func New() *Installer {
	return &Installer{}
}

func (installer *Installer) Upgrade() error {
	log.Println("Running installer upgrade", UpgradeCmd)
	out, err := exec.Command(UpgradeCmd).CombinedOutput()
	log.Printf("Installer upgrade output %s", out)
	if err != nil {
		return err
	}
	return nil
}
