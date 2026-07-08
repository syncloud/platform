package snap

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"
)

const (
	SnapdAppsUrl        = "http://apps.syncloud.org/apps"
	SnapdWorkDir        = "/tmp"
	SnapdUpgradeScript  = "snapd/upgrade.sh"
	SnapdKeepArchives   = 2
	SnapdVerifyAttempts = 10
	SnapdVerifyDelay    = 3 * time.Second
)

var snapdArchByGoArch = map[string]string{
	"amd64": "amd64",
	"arm64": "arm64",
	"arm":   "armhf",
}

type InstalledVersionProvider interface {
	InstalledVersion() (string, error)
}

type Snapd struct {
	version InstalledVersionProvider
	sleep   func(time.Duration)
	logger  *zap.Logger
}

func NewSnapd(version InstalledVersionProvider, logger *zap.Logger) *Snapd {
	return &Snapd{
		version: version,
		sleep:   time.Sleep,
		logger:  logger,
	}
}

func (s *Snapd) Upgrade(version string) error {
	version = strings.TrimSpace(version)
	arch, err := s.arch()
	if err != nil {
		return err
	}
	name := fmt.Sprintf("snapd-%s-%s.tar.gz", version, arch)
	url := fmt.Sprintf("%s/%s", SnapdAppsUrl, name)
	archive := filepath.Join(SnapdWorkDir, name)
	s.logger.Info("downloading snapd", zap.String("url", url))
	err = s.download(url, archive)
	if err != nil {
		return err
	}
	s.prune(SnapdWorkDir, SnapdKeepArchives)

	err = os.RemoveAll(filepath.Join(SnapdWorkDir, "snapd"))
	if err != nil {
		return err
	}
	err = s.extract(archive, SnapdWorkDir)
	if err != nil {
		return err
	}

	cmd := exec.Command("./" + SnapdUpgradeScript)
	cmd.Dir = SnapdWorkDir
	s.logger.Info("running snapd upgrade", zap.String("script", SnapdUpgradeScript))
	out, err := cmd.CombinedOutput()
	s.logger.Info("snapd upgrade output", zap.String("output", string(out)))
	if err != nil {
		return err
	}

	return s.verify(version)
}

func (s *Snapd) verify(expected string) error {
	var last string
	for attempt := 0; attempt < SnapdVerifyAttempts; attempt++ {
		if attempt > 0 {
			s.sleep(SnapdVerifyDelay)
		}
		installed, err := s.version.InstalledVersion()
		if err != nil {
			last = err.Error()
			continue
		}
		installed = strings.TrimSpace(installed)
		if installed == expected {
			s.logger.Info("snapd upgrade verified", zap.String("version", installed))
			return nil
		}
		last = fmt.Sprintf("installed %s, expected %s", installed, expected)
	}
	return fmt.Errorf("snapd upgrade verification failed: %s", last)
}

func (s *Snapd) prune(dir string, keep int) {
	matches, err := filepath.Glob(filepath.Join(dir, "snapd-*.tar.gz"))
	if err != nil {
		s.logger.Warn("cannot list snapd archives", zap.Error(err))
		return
	}
	type entry struct {
		path string
		mod  time.Time
	}
	var entries []entry
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			continue
		}
		entries = append(entries, entry{path: match, mod: info.ModTime()})
	}
	sort.Slice(entries, func(a, b int) bool { return entries[a].mod.After(entries[b].mod) })
	for idx := keep; idx < len(entries); idx++ {
		s.logger.Info("removing old snapd archive", zap.String("path", entries[idx].path))
		err = os.Remove(entries[idx].path)
		if err != nil {
			s.logger.Warn("cannot remove old snapd archive", zap.String("path", entries[idx].path), zap.Error(err))
		}
	}
}

func (s *Snapd) arch() (string, error) {
	arch, ok := snapdArchByGoArch[runtime.GOARCH]
	if !ok {
		return "", fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
	}
	return arch, nil
}

func (s *Snapd) download(url, dst string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed %s: %s", url, resp.Status)
	}
	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	return err
}

func (s *Snapd) extract(archive, dst string) error {
	file, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer file.Close()
	gr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(dst, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}
			out.Close()
		}
	}
	return nil
}
