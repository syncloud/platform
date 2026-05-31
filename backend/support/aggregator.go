package support

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	Pattern   = "/var/snap/*/common/log/*.log"
	Separator = "\n----------------------\n"
)

type LogAggregator struct {
	logger *zap.Logger
}

func NewAggregator(logger *zap.Logger) *LogAggregator {
	return &LogAggregator{
		logger: logger,
	}
}

func (a *LogAggregator) GetLogs() string {
	log := a.fileLogs()
	log += a.cmd("date")
	log += a.cmd("mount")
	log += a.cmd("systemctl", "status", "--state=failed", "snap.*")
	log += a.cmd("top", "-n", "1", "-bc")
	log += a.cmd("ping", "google.com", "-c", "5")
	log += a.cmd("uname", "-a")
	log += a.cmd("cat", "/etc/debian_version")
	log += a.cmd("free", "-h")
	log += a.cmd("df", "-h")
	log += a.cmd("lsblk", "-o", "+UUID")
	log += a.cmd("lsblk", "-Pp", "-o", "NAME,SIZE,TYPE,MOUNTPOINT,PARTTYPE,FSTYPE,MODEL")
	log += a.cmd("ls", "-la", "/data")
	log += a.cmd("uptime")
	log += a.cmd("snap", "version")
	log += a.cmd("snap", "list")
	log += a.cmd("snap", "list", "--all")
	log += a.cmd("snap", "changes")
	log += a.snapChangesDetail()
	log += a.cmd("snap", "services")
	log += a.cmd("snap", "run", "platform.cli", "ipv4", "public")
	log += a.journal()
	log += a.previousBootTail()
	log += a.cmd("dmesg", "-T")
	log += a.dmesgErrors()
	log += a.cmd("cat", "/proc/diskstats")
	return log
}

var noise = regexp.MustCompile(`dhclient\[|]: (XMT|RCV|PRC):`)

func (a *LogAggregator) journal() string {
	command := exec.Command("journalctl", "-n", "5000", "--no-pager")
	out, err := command.CombinedOutput()
	if err != nil {
		a.logger.Warn("failed", zap.Error(err))
	}
	signal, noisy := splitNoise(string(out))
	return command.String() + "\n\n" + tail(signal, 1000) + Separator +
		"journal noise tail (filtered out above)\n\n" + tail(noisy, 100) + Separator
}

func splitNoise(input string) (string, string) {
	var signal, noisy []string
	for _, line := range strings.Split(input, "\n") {
		if noise.MatchString(line) {
			noisy = append(noisy, line)
		} else {
			signal = append(signal, line)
		}
	}
	return strings.Join(signal, "\n"), strings.Join(noisy, "\n")
}

func (a *LogAggregator) previousBootTail() string {
	command := exec.Command("journalctl", "-b", "-1", "-n", "100", "--no-pager")
	out, err := command.CombinedOutput()
	if err != nil {
		a.logger.Warn("failed", zap.Error(err))
	}
	return command.String() + " (tail before last reboot)\n\n" + string(out) + Separator
}

func (a *LogAggregator) snapChangesDetail() string {
	result := "snap changes detail\n\n"
	out, err := exec.Command("snap", "changes").CombinedOutput()
	if err != nil {
		a.logger.Warn("failed", zap.Error(err))
		return result + string(out) + Separator
	}
	var pending, recent []string
	for i, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if i == 0 {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		id, status := fields[0], fields[1]
		if status != "Done" && len(pending) < 5 {
			pending = append(pending, id)
		}
		if len(recent) < 5 {
			recent = append(recent, id)
		}
	}
	seen := map[string]bool{}
	for _, id := range append(pending, recent...) {
		if seen[id] {
			continue
		}
		seen[id] = true
		result += fmt.Sprintf("=== snap change %s ===\n", id)
		out, err := exec.Command("snap", "change", id).CombinedOutput()
		if err != nil {
			a.logger.Warn("failed", zap.Error(err))
		}
		result += tail(string(out), 100) + "\n"
	}
	return result + Separator
}

func (a *LogAggregator) dmesgErrors() string {
	result := "dmesg errors\n\n"
	out, err := exec.Command("dmesg", "-T").CombinedOutput()
	if err != nil {
		a.logger.Warn("failed", zap.Error(err))
	}
	re := regexp.MustCompile(`(?i)I/O error|EXT4-fs error|usb-storage|disconnect`)
	var matches []string
	for _, line := range strings.Split(string(out), "\n") {
		if re.MatchString(line) {
			matches = append(matches, line)
		}
	}
	if len(matches) > 50 {
		matches = matches[len(matches)-50:]
	}
	result += strings.Join(matches, "\n") + "\n"
	return result + Separator
}

func tail(s string, n int) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, "\n")
}

func (a *LogAggregator) cmd(app string, args ...string) string {
	command := exec.Command(app, args...)
	if app == "top" {
		command.Env = append(os.Environ(), "COLUMNS=1000")
	}
	result := command.String() + "\n\n"
	output, err := command.CombinedOutput()
	if err != nil {
		a.logger.Info(string(output))
		a.logger.Warn("failed", zap.Error(err))
	}

	result += string(output) + Separator
	return result
}

func (a *LogAggregator) fileLogs() string {

	matches, err := filepath.Glob(Pattern)
	if err != nil {
		a.logger.Error("failed", zap.Error(err))
		return ""
	}

	log := ""
	for _, file := range matches {
		log += fmt.Sprintf("file: %s\n\n", file)
		output, err := exec.Command("tail", "-100", file).CombinedOutput()
		if err != nil {
			a.logger.Error("failed", zap.Error(err))
			return ""
		}
		log += string(output)
		log += Separator
	}
	return log
}
