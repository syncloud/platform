package system

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const UptimeFile = "/proc/uptime"

type Uptime struct {
	file string
}

func NewUptime() *Uptime {
	return &Uptime{file: UptimeFile}
}

func (u *Uptime) Seconds() (int, error) {
	content, err := os.ReadFile(u.file)
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(string(content))
	if len(fields) == 0 {
		return 0, fmt.Errorf("cannot parse %s: %s", u.file, string(content))
	}
	seconds, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse %s: %w", u.file, err)
	}
	return int(seconds), nil
}
