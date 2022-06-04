package systemd

import (
	"fmt"
	"github.com/syncloud/platform/cli"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

const Dir = "/lib/systemd/system"

type Control struct {
	config   ControlConfig
	executor cli.CommandExecutor
	logger   *zap.Logger
}

type ControlConfig interface {
	ExternalDiskDir() string
}

func New(executor cli.CommandExecutor, config ControlConfig, logger *zap.Logger) *Control {
	return &Control{executor: executor, config: config, logger: logger}
}

func (c *Control) ReloadService(service string) error {

	log.Printf("reloading %s\n", service)
	output, err := exec.Command("systemctl", "reload", fmt.Sprintf("snap.%s", service)).CombinedOutput()
	log.Printf("systemctl output: %s", string(output))
	return err
}

func (c *Control) RemoveMount() error {
	return c.remove(c.DirToSystemdMountFilename(c.config.ExternalDiskDir()))
}

func (c *Control) DirToSystemdMountFilename(directory string) string {
	directoryClean := strings.TrimPrefix(directory, "/")
	return strings.Join(strings.Split(directoryClean, "/"), "-") + ".mount"
}

func (c *Control) remove(filename string) error {

	status, err := c.stop(filename)
	if err != nil {
		return err
	}
	if slices.Contains([]string{"unknown", "inactive"}, *status) {
		return nil
	}
	output, err := c.executor.CommandOutput("systemctl", "disable", filename)
	if err != nil {
		c.logger.Error(string(output))
		return err
	}
	systemdFile := c.systemdFile(filename)
	if _, err = os.Stat(systemdFile); err == nil {
		err = os.Remove(systemdFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Control) systemdFile(filename string) string {
	return path.Join(Dir, filename)
}

func (c *Control) stop(service string) (*string, error) {

	c.logger.Info("checking", zap.String("service", service))
	//TODO: exit code 3 when inactive
	isAliveOutput, err := c.executor.CommandOutput("systemctl", "is-active", service)
	if err != nil {
		return nil, err
	}
	result := strings.TrimSpace(string(isAliveOutput))
	c.logger.Info("stopping", zap.String("service", service))
	stopOutput, err := c.executor.CommandOutput("systemctl", "stop", service)
	if err != nil {
		result = strings.TrimSpace(string(stopOutput))
	}

	c.logger.Info("result", zap.String(service, result))
	return &result, nil
}
