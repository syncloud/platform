package stability

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Victim struct {
	PID    int
	Comm   string
	Cgroup string
	RSSkB  uint64
	OOMAdj int
	Score  float64
}

type ProcScanner struct {
	procDir string
}

func NewProcScanner(procDir string) *ProcScanner {
	return &ProcScanner{procDir: procDir}
}

func (s *ProcScanner) Candidates(protect Protect, selfPID int) ([]Victim, error) {
	entries, err := os.ReadDir(s.procDir)
	if err != nil {
		return nil, err
	}
	var out []Victim
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue
		}
		if pid <= 1 || pid == selfPID {
			continue
		}
		v, err := s.readVictim(pid)
		if err != nil {
			continue
		}
		if v.RSSkB == 0 {
			continue
		}
		if protect.IsProtected(v) {
			continue
		}
		v.Score = score(v.RSSkB, v.OOMAdj)
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out, nil
}

func (s *ProcScanner) readVictim(pid int) (Victim, error) {
	base := filepath.Join(s.procDir, strconv.Itoa(pid))
	v := Victim{PID: pid}
	if err := readStatus(filepath.Join(base, "status"), &v); err != nil {
		return Victim{}, err
	}
	if adj, err := readInt(filepath.Join(base, "oom_score_adj")); err == nil {
		v.OOMAdj = adj
	}
	if cg, err := os.ReadFile(filepath.Join(base, "cgroup")); err == nil {
		v.Cgroup = strings.TrimSpace(string(cg))
	}
	return v, nil
}

func readStatus(path string, v *Victim) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		switch {
		case strings.HasPrefix(line, "Name:"):
			v.Comm = strings.TrimSpace(strings.TrimPrefix(line, "Name:"))
		case strings.HasPrefix(line, "VmRSS:"):
			rest := strings.TrimSpace(strings.TrimPrefix(line, "VmRSS:"))
			rest = strings.TrimSuffix(rest, " kB")
			v.RSSkB, _ = strconv.ParseUint(rest, 10, 64)
		}
	}
	return sc.Err()
}

func readInt(path string) (int, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	s := strings.TrimSpace(string(b))
	if s == "" {
		return 0, errors.New("empty")
	}
	return strconv.Atoi(s)
}

func score(rssKB uint64, adj int) float64 {
	mult := 1.0
	if adj > 0 {
		mult += float64(adj) / 1000.0
	}
	return float64(rssKB) * mult
}

var ErrNoVictim = fmt.Errorf("stability: no eligible victim")
