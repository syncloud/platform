package systemd

import (
	"log"
	"os/exec"
)

type Control struct {
}

func New() *Control {
	return &Control{}
}

func (c *Control) ReloadService(service string) error {

	log.Printf("reloading %s\n", service)
	_, err := exec.Command("systemctl", "reload", service).CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
