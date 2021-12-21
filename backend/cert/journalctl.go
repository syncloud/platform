package cert

import (
	"fmt"
	"os/exec"
)

type Ctl struct{}

type JournalCtl interface {
	Read(unit string) (string, error)
}

func NewJournalCtl() *Ctl {
	return &Ctl{}
}

func (j *Ctl) Read(unit string) (string, error) {
	output, err := exec.Command("journalctl", "-u", fmt.Sprintf("snap.%s", unit), "-n", "1000", "--no-pager").CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
