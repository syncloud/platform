package stability

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type MemInfo struct {
	procDir string
}

type MemSnap struct {
	TotalKB     uint64
	AvailableKB uint64
}

func (s MemSnap) AvailableRatio() float64 {
	if s.TotalKB == 0 {
		return 1
	}
	return float64(s.AvailableKB) / float64(s.TotalKB)
}

func NewMemInfo(procDir string) *MemInfo {
	return &MemInfo{procDir: procDir}
}

func (m *MemInfo) Snapshot() (MemSnap, error) {
	f, err := os.Open(filepath.Join(m.procDir, "meminfo"))
	if err != nil {
		return MemSnap{}, err
	}
	defer f.Close()
	s := MemSnap{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		key, val := parseMemLine(line)
		switch key {
		case "MemTotal":
			s.TotalKB = val
		case "MemAvailable":
			s.AvailableKB = val
		}
	}
	if err := sc.Err(); err != nil {
		return MemSnap{}, err
	}
	if s.TotalKB == 0 {
		return MemSnap{}, fmt.Errorf("meminfo: MemTotal missing")
	}
	return s, nil
}

func parseMemLine(line string) (string, uint64) {
	idx := strings.IndexByte(line, ':')
	if idx < 0 {
		return "", 0
	}
	key := line[:idx]
	rest := strings.TrimSpace(line[idx+1:])
	rest = strings.TrimSuffix(rest, " kB")
	v, err := strconv.ParseUint(rest, 10, 64)
	if err != nil {
		return key, 0
	}
	return key, v
}

func (m *MemInfo) PSIAvailable() bool {
	_, err := os.Stat(filepath.Join(m.procDir, "pressure", "memory"))
	return err == nil
}

func (m *MemInfo) PSIMemoryAvg10() (float64, error) {
	f, err := os.Open(filepath.Join(m.procDir, "pressure", "memory"))
	if err != nil {
		return 0, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if !strings.HasPrefix(line, "some ") {
			continue
		}
		for _, field := range strings.Fields(line) {
			if v, ok := strings.CutPrefix(field, "avg10="); ok {
				return strconv.ParseFloat(v, 64)
			}
		}
	}
	return 0, fmt.Errorf("pressure/memory: 'some avg10' missing")
}
