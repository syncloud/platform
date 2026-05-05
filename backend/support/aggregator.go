package support

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"path/filepath"
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
	log += a.cmd("sh", "-c", "du -sh /var/snap/*/current/database 2>/dev/null; du -sh /var/snap/*/current/database.dump 2>/dev/null")
	log += a.cmd("uptime")
	log += a.cmd("snap", "version")
	log += a.cmd("snap", "list")
	log += a.cmd("snap", "list", "--all")
	log += a.cmd("snap", "changes")
	log += a.cmd("sh", "-c", "for id in $(snap changes 2>/dev/null | awk 'NR>1 && $2!~/^Done$/ {print $1}' | head -5; snap changes 2>/dev/null | awk 'NR>1 {print $1}' | head -5); do echo === snap change $id ===; snap change \"$id\" 2>&1; done")
	log += a.cmd("snap", "services")
	log += a.cmd("snap", "run", "platform.cli", "ipv4", "public")
	log += a.cmd("journalctl", "--since=-3h", "--no-pager")
	log += a.cmd("dmesg", "-T")
	log += a.cmd("sh", "-c", "dmesg -T 2>&1 | grep -iE 'I/O error|EXT4-fs error|usb-storage|disconnect' | tail -50")
	log += a.cmd("sh", "-c", "for d in /dev/sd? /dev/nvme?n? /dev/mmcblk?; do [ -b \"$d\" ] && smartctl -a \"$d\" 2>&1; done")
	log += a.cmd("cat", "/proc/diskstats")
	return log
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
