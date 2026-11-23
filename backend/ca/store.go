package ca

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	cp "github.com/otiai10/copy"
	"github.com/syncloud/platform/cli"
	"go.uber.org/zap"
)

const prefix = "mozilla/"

type Store struct {
	shippedDir string
	systemDir  string
	confFile   string
	executor   cli.Executor
	logger     *zap.Logger
}

func NewStore(shippedDir, systemDir, confFile string, executor cli.Executor, logger *zap.Logger) *Store {
	return &Store{
		shippedDir: shippedDir,
		systemDir:  systemDir,
		confFile:   confFile,
		executor:   executor,
		logger:     logger,
	}
}

func (s *Store) Update() error {
	shipped, err := certificates(s.shippedDir)
	if err != nil {
		return err
	}
	if len(shipped) == 0 {
		return fmt.Errorf("no certificates in %s", s.shippedDir)
	}

	err = os.MkdirAll(s.systemDir, 0755)
	if err != nil {
		return err
	}
	installed, err := certificates(s.systemDir)
	if err != nil {
		return err
	}
	for _, name := range installed {
		if !slices.Contains(shipped, name) {
			s.logger.Info("removing certificate", zap.String("name", name))
			err = os.Remove(filepath.Join(s.systemDir, name))
			if err != nil {
				return err
			}
		}
	}

	err = cp.Copy(s.shippedDir, s.systemDir)
	if err != nil {
		return err
	}

	err = s.updateConf(shipped)
	if err != nil {
		return err
	}

	_, err = s.executor.CombinedOutput("update-ca-certificates", "--fresh")
	return err
}

func (s *Store) updateConf(shipped []string) error {
	content, err := os.ReadFile(s.confFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var lines []string
	listed := map[string]bool{}
	for _, line := range strings.Split(strings.TrimRight(string(content), "\n"), "\n") {
		name, ok := strings.CutPrefix(strings.TrimPrefix(line, "!"), prefix)
		if !ok {
			if line != "" {
				lines = append(lines, line)
			}
			continue
		}
		if !slices.Contains(shipped, name) || listed[name] {
			continue
		}
		listed[name] = true
		lines = append(lines, line)
	}
	for _, name := range shipped {
		if !listed[name] {
			lines = append(lines, prefix+name)
		}
	}

	return os.WriteFile(s.confFile, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

func certificates(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".crt") {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}
