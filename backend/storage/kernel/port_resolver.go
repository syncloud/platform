package kernel

import (
	"path/filepath"
	"regexp"
	"strings"
)

var ataSegment = regexp.MustCompile(`/(ata\d+)/`)
var usbSegment = regexp.MustCompile(`/(usb\d+)/`)

type PortResolver struct {
	sysRoot string
}

func NewPortResolver(sysRoot string) *PortResolver {
	return &PortResolver{sysRoot: sysRoot}
}

func (r *PortResolver) Port(device string) string {
	path, err := filepath.EvalSymlinks(filepath.Join(r.sysRoot, "block", strings.TrimPrefix(device, "/dev/")))
	if err != nil {
		return ""
	}
	if match := ataSegment.FindStringSubmatch(path); match != nil {
		return match[1]
	}
	if match := usbSegment.FindStringSubmatch(path); match != nil {
		return match[1]
	}
	return ""
}
