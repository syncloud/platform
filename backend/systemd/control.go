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
	"text/template"
)

const Dir = "/lib/systemd/system"

type Control struct {
	config   ControlConfig
	executor cli.Executor
	logger   *zap.Logger
}

type ControlConfig interface {
	ExternalDiskDir() string
	ConfigDir() string
}

func New(executor cli.Executor, config ControlConfig, logger *zap.Logger) *Control {
	return &Control{executor: executor, config: config, logger: logger}
}

func (c *Control) RestartService(service string) error {
	serviceName := c.serviceName(service)
	output, err := exec.Command("systemctl", "restart", serviceName).CombinedOutput()
	log.Printf("systemctl output: %s", string(output))
	return err
}

func (c *Control) ReloadService(service string) error {

	log.Printf("reloading %s\n", service)
	serviceName := c.serviceName(service)
	output, err := exec.Command("systemctl", "reload", serviceName).CombinedOutput()
	log.Printf("systemctl output: %s", string(output))
	return err
}

func (c *Control) serviceName(service string) string {
	return fmt.Sprintf("snap.%s", service)
}

func (c *Control) RemoveMount() error {
	return c.remove(c.DirToSystemdMountFilename(c.config.ExternalDiskDir()))
}

func (c *Control) AddMount(device string) error {
	c.logger.Info("adding mount", zap.String("device", device))
	mountTemplateFile := path.Join(c.config.ConfigDir(), "mount", "mount.template")
	mountDefinition, err := template.ParseFiles(mountTemplateFile)
	if err != nil {
		return err
	}
	mountFilename := c.DirToSystemdMountFilename(c.config.ExternalDiskDir())
	systemdFilename := c.systemdFile(mountFilename)
	f, err := os.Create(systemdFilename)
	if err != nil {
		return err
	}
	err = mountDefinition.Execute(f, &Mount{What: device, Where: c.config.ExternalDiskDir()})
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	c.logger.Info("enabling", zap.String("file", mountFilename))
	_, err = c.executor.CombinedOutput("systemctl", "enable", mountFilename)
	if err != nil {
		return err
	}
	return c.start(mountFilename)
}

func (c *Control) start(service string) error {
	c.logger.Info("starting", zap.String("service", service))
	output, err := c.executor.CombinedOutput("systemctl", "start", service)
	if err != nil {
		c.logger.Error("unable to start a service", zap.String("output", string(output)))
		logOutput, logErr := c.executor.CombinedOutput("journalctl", "-u", service)
		if logErr != nil {
			c.logger.Error("unable to get service log", zap.String("log output", string(logOutput)))
		}
		return err
	}
	return nil
}

func (c *Control) DirToSystemdMountFilename(directory string) string {
	directoryClean := strings.TrimPrefix(directory, "/")
	return strings.Join(strings.Split(directoryClean, "/"), "-") + ".mount"
}

func (c *Control) remove(filename string) error {

	status := c.stop(filename)
	if slices.Contains([]string{"unknown", "inactive", "failed"}, status) {
		return nil
	}
	output, err := c.executor.CombinedOutput("systemctl", "disable", filename)
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

func (c *Control) stop(service string) string {

	c.logger.Info("checking", zap.String("service", service))
	isAliveOutput, err := c.executor.CombinedOutput("systemctl", "is-active", service)
	isAliveResult := strings.TrimSpace(string(isAliveOutput))
	if err != nil {
		c.logger.Info("is-active", zap.String("output", string(isAliveOutput)))
		return isAliveResult
	}
	c.logger.Info("stopping", zap.String("service", service))
	stopOutput, err := c.executor.CombinedOutput("systemctl", "stop", service)
	if err != nil {
		resultResult := strings.TrimSpace(string(stopOutput))
		return resultResult
	}

	c.logger.Info("result", zap.String(service, isAliveResult))
	return isAliveResult
}
