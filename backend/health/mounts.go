package health

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func (c *Collector) Mounts() []Mount {
	f, err := os.Open(filepath.Join(c.procDir, "mounts"))
	if err != nil {
		return nil
	}
	defer f.Close()
	var out []Mount
	seen := map[string]bool{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) < 3 {
			continue
		}
		dev, mount, fs := fields[0], fields[1], fields[2]
		if !strings.HasPrefix(dev, "/dev/") {
			continue
		}
		if fs == "squashfs" || fs == "tmpfs" || fs == "devtmpfs" {
			continue
		}
		if seen[mount] {
			continue
		}
		seen[mount] = true
		var st syscall.Statfs_t
		if err := syscall.Statfs(mount, &st); err != nil {
			continue
		}
		total := st.Blocks * uint64(st.Bsize) / 1024
		free := st.Bavail * uint64(st.Bsize) / 1024
		out = append(out, Mount{Path: mount, TotalKB: total, UsedKB: total - free})
	}
	return out
}
