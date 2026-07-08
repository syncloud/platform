package installer

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
	AppsUrl        = "http://apps.syncloud.org/apps"
	WorkDir        = "/tmp"
	UpgradeScript  = "snapd/upgrade.sh"
	KeepArchives   = 2
	VerifyAttempts = 10
	VerifyDelay    = 3 * time.Second
)

var archByGoArch = map[string]string{
	"amd64": "amd64",
	"arm64": "arm64",
	"arm":   "armhf",
}

type SnapdVersion interface {
	InstalledVersion() (string, error)
}

type Installer struct {
	snapd  SnapdVersion
	sleep  func(time.Duration)
	logger *zap.Logger
}

func New(snapd SnapdVersion, logger *zap.Logger) *Installer {
	return &Installer{
		snapd:  snapd,
		sleep:  time.Sleep,
		logger: logger,
	}
}

func (i *Installer) Upgrade(version string) error {
	version = strings.TrimSpace(version)
	arch, err := i.arch()
	if err != nil {
		return err
	}
	name := fmt.Sprintf("snapd-%s-%s.tar.gz", version, arch)
	url := fmt.Sprintf("%s/%s", AppsUrl, name)
	archive := filepath.Join(WorkDir, name)
	i.logger.Info("downloading snapd", zap.String("url", url))
	err = i.download(url, archive)
	if err != nil {
		return err
	}
	i.prune(WorkDir, KeepArchives)

	err = os.RemoveAll(filepath.Join(WorkDir, "snapd"))
	if err != nil {
		return err
	}
	err = i.extract(archive, WorkDir)
	if err != nil {
		return err
	}

	cmd := exec.Command("./" + UpgradeScript)
	cmd.Dir = WorkDir
	i.logger.Info("running snapd upgrade", zap.String("script", UpgradeScript))
	out, err := cmd.CombinedOutput()
	i.logger.Info("snapd upgrade output", zap.String("output", string(out)))
	if err != nil {
		return err
	}

	return i.verify(version)
}

func (i *Installer) verify(expected string) error {
	var last string
	for attempt := 0; attempt < VerifyAttempts; attempt++ {
		if attempt > 0 {
			i.sleep(VerifyDelay)
		}
		installed, err := i.snapd.InstalledVersion()
		if err != nil {
			last = err.Error()
			continue
		}
		installed = strings.TrimSpace(installed)
		if installed == expected {
			i.logger.Info("snapd upgrade verified", zap.String("version", installed))
			return nil
		}
		last = fmt.Sprintf("installed %s, expected %s", installed, expected)
	}
	return fmt.Errorf("snapd upgrade verification failed: %s", last)
}

func (i *Installer) prune(dir string, keep int) {
	matches, err := filepath.Glob(filepath.Join(dir, "snapd-*.tar.gz"))
	if err != nil {
		i.logger.Warn("cannot list snapd archives", zap.Error(err))
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
		i.logger.Info("removing old snapd archive", zap.String("path", entries[idx].path))
		err = os.Remove(entries[idx].path)
		if err != nil {
			i.logger.Warn("cannot remove old snapd archive", zap.String("path", entries[idx].path), zap.Error(err))
		}
	}
}

func (i *Installer) arch() (string, error) {
	arch, ok := archByGoArch[runtime.GOARCH]
	if !ok {
		return "", fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
	}
	return arch, nil
}

func (i *Installer) download(url, dst string) error {
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

func (i *Installer) extract(archive, dst string) error {
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
