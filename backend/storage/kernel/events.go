package kernel

import (
	"github.com/syncloud/platform/storage/model"
	"regexp"
	"strings"
)

var (
	ataPort   = regexp.MustCompile(`^(ata\d+)(\.\d+)?: (.*)$`)
	linkEvent = regexp.MustCompile(`hard resetting link|SError:|exception Emask|COMRESET failed|link is slow to respond|failed to recover|limiting SATA link speed`)
	ioDevice  = regexp.MustCompile(`(?:I/O error|critical medium error|critical target error), dev (\w+)|Buffer I/O error on dev (\w+)|bdev /dev/(\w+) errs:`)
)

type Events struct {
	lines []Line
}

func NewEvents(lines []Line) *Events {
	return &Events{lines: lines}
}

func (e *Events) Link(port string) model.KernelEvents {
	result := model.KernelEvents{}
	if port == "" {
		return result
	}
	for _, line := range e.lines {
		match := ataPort.FindStringSubmatch(line.Text)
		if match == nil || match[1] != port {
			continue
		}
		if linkEvent.MatchString(match[3]) {
			result.Count++
			result.Last = line.Time
		}
	}
	return result
}

func (e *Events) Io(device string) model.KernelEvents {
	result := model.KernelEvents{}
	name := strings.TrimPrefix(device, "/dev/")
	for _, line := range e.lines {
		match := ioDevice.FindStringSubmatch(line.Text)
		if match == nil {
			continue
		}
		if belongsTo(match[1]+match[2]+match[3], name) {
			result.Count++
			result.Last = line.Time
		}
	}
	return result
}

func belongsTo(name string, disk string) bool {
	if name == disk {
		return true
	}
	if !strings.HasPrefix(name, disk) {
		return false
	}
	suffix := strings.TrimPrefix(name, disk)
	if last := disk[len(disk)-1]; last >= '0' && last <= '9' {
		if !strings.HasPrefix(suffix, "p") {
			return false
		}
		suffix = strings.TrimPrefix(suffix, "p")
	}
	if suffix == "" {
		return false
	}
	for _, c := range suffix {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
