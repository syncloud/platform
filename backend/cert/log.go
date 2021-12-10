package cert

import "os/exec"

type Log struct{}

func NewLog() *Log {
	return &Log{}
}

func (l *Log) Load() string {
	output, err := exec.Command("journalctl", "-F", l.userConfDir, "-b", "cn=config", "-l", initScript).CombinedOutput()
	if err != nil {
		return err
	}

}
