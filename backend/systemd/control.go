package systemd

import (
	"fmt"
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
	output, err := exec.Command("systemctl", "reload", fmt.Sprintf("snap.%s", service)).CombinedOutput()
	log.Printf("systemctl output: %s", string(output))
	return err
}
