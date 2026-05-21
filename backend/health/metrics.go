package health

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type CPU struct {
	User   uint64 `json:"user"`
	Nice   uint64 `json:"nice"`
	System uint64 `json:"system"`
	Idle   uint64 `json:"idle"`
	IOWait uint64 `json:"iowait"`
	IRQ    uint64 `json:"irq"`
	SoftIRQ uint64 `json:"softirq"`
	Steal  uint64 `json:"steal"`
}

func (c CPU) Total() uint64 {
	return c.User + c.Nice + c.System + c.Idle + c.IOWait + c.IRQ + c.SoftIRQ + c.Steal
}

func (c CPU) Busy() uint64 {
	return c.Total() - c.Idle - c.IOWait
}

type Memory struct {
	TotalKB     uint64 `json:"total_kb"`
	AvailableKB uint64 `json:"available_kb"`
	FreeKB      uint64 `json:"free_kb"`
	BuffersKB   uint64 `json:"buffers_kb"`
	CachedKB    uint64 `json:"cached_kb"`
	SwapTotalKB uint64 `json:"swap_total_kb"`
	SwapFreeKB  uint64 `json:"swap_free_kb"`
}

type Disk struct {
	Name        string `json:"name"`
	ReadsTotal  uint64 `json:"reads_total"`
	WritesTotal uint64 `json:"writes_total"`
	SectorsRead uint64 `json:"sectors_read"`
	SectorsWrt  uint64 `json:"sectors_written"`
}

type Mount struct {
	Path     string `json:"path"`
	TotalKB  uint64 `json:"total_kb"`
	UsedKB   uint64 `json:"used_kb"`
}

type Net struct {
	Name    string `json:"name"`
	RxBytes uint64 `json:"rx_bytes"`
	TxBytes uint64 `json:"tx_bytes"`
}

type Snapshot struct {
	CPU    CPU    `json:"cpu"`
	Memory Memory `json:"memory"`
	Disks  []Disk `json:"disks"`
	Mounts []Mount `json:"mounts"`
	Net    []Net  `json:"net"`
}

type Collector struct {
	procDir string
}

func NewCollector(procDir string) *Collector {
	return &Collector{procDir: procDir}
}

func (c *Collector) Snapshot() (Snapshot, error) {
	var s Snapshot
	cpu, err := c.readCPU()
	if err != nil {
		return s, err
	}
	s.CPU = cpu
	mem, err := c.readMemory()
	if err != nil {
		return s, err
	}
	s.Memory = mem
	s.Disks, _ = c.readDisks()
	s.Net, _ = c.readNet()
	s.Mounts = c.Mounts()
	return s, nil
}

func (c *Collector) readCPU() (CPU, error) {
	f, err := os.Open(filepath.Join(c.procDir, "stat"))
	if err != nil {
		return CPU{}, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if !strings.HasPrefix(line, "cpu ") {
			continue
		}
		fields := strings.Fields(line)
		nums := make([]uint64, 0, 8)
		for _, fld := range fields[1:] {
			n, _ := strconv.ParseUint(fld, 10, 64)
			nums = append(nums, n)
		}
		for len(nums) < 8 {
			nums = append(nums, 0)
		}
		return CPU{nums[0], nums[1], nums[2], nums[3], nums[4], nums[5], nums[6], nums[7]}, nil
	}
	return CPU{}, fmt.Errorf("cpu: 'cpu ' line missing")
}

func (c *Collector) readMemory() (Memory, error) {
	f, err := os.Open(filepath.Join(c.procDir, "meminfo"))
	if err != nil {
		return Memory{}, err
	}
	defer f.Close()
	m := Memory{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		idx := strings.IndexByte(line, ':')
		if idx < 0 {
			continue
		}
		key := line[:idx]
		rest := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(line[idx+1:]), " kB"))
		v, _ := strconv.ParseUint(rest, 10, 64)
		switch key {
		case "MemTotal":
			m.TotalKB = v
		case "MemAvailable":
			m.AvailableKB = v
		case "MemFree":
			m.FreeKB = v
		case "Buffers":
			m.BuffersKB = v
		case "Cached":
			m.CachedKB = v
		case "SwapTotal":
			m.SwapTotalKB = v
		case "SwapFree":
			m.SwapFreeKB = v
		}
	}
	return m, sc.Err()
}

func (c *Collector) readDisks() ([]Disk, error) {
	f, err := os.Open(filepath.Join(c.procDir, "diskstats"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []Disk
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) < 14 {
			continue
		}
		name := fields[2]
		if isPartition(name) {
			continue
		}
		reads, _ := strconv.ParseUint(fields[3], 10, 64)
		sectorsRead, _ := strconv.ParseUint(fields[5], 10, 64)
		writes, _ := strconv.ParseUint(fields[7], 10, 64)
		sectorsWrt, _ := strconv.ParseUint(fields[9], 10, 64)
		out = append(out, Disk{
			Name:        name,
			ReadsTotal:  reads,
			SectorsRead: sectorsRead,
			WritesTotal: writes,
			SectorsWrt:  sectorsWrt,
		})
	}
	return out, sc.Err()
}

func isPartition(name string) bool {
	if strings.HasPrefix(name, "loop") || strings.HasPrefix(name, "ram") || strings.HasPrefix(name, "dm-") {
		return true
	}
	if len(name) == 0 {
		return true
	}
	last := name[len(name)-1]
	if last < '0' || last > '9' {
		return false
	}
	if strings.HasPrefix(name, "mmcblk") || strings.HasPrefix(name, "nvme") {
		return strings.Contains(name, "p")
	}
	return true
}

func (c *Collector) readNet() ([]Net, error) {
	f, err := os.Open(filepath.Join(c.procDir, "net/dev"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []Net
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		idx := strings.IndexByte(line, ':')
		if idx < 0 {
			continue
		}
		name := strings.TrimSpace(line[:idx])
		if name == "lo" {
			continue
		}
		fields := strings.Fields(line[idx+1:])
		if len(fields) < 9 {
			continue
		}
		rx, _ := strconv.ParseUint(fields[0], 10, 64)
		tx, _ := strconv.ParseUint(fields[8], 10, 64)
		out = append(out, Net{Name: name, RxBytes: rx, TxBytes: tx})
	}
	return out, sc.Err()
}
