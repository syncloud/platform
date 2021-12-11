package cert

import (
	"os/exec"
	"strings"
)

type Reader struct{}

func NewReader() *Reader {
	return &Reader{}
}

func (l *Reader) Read() []string {

	output, err := exec.Command("journalctl", "-u", "snap.platform.backend", "-n", "1000", "--no-pager").CombinedOutput()
	if err != nil {
		return []string{err.Error()}
	}

	return strings.Split(string(output), "\n")
}
