package hook

import (
	"errors"
	"fmt"
	cp "github.com/otiai10/copy"
	"github.com/syncloud/golib/linux"
	"github.com/syncloud/platform/auth"
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/storage"
	"go.uber.org/zap"
	"os"
	"path"
)

type Install struct {
	storageChecker storage.Checker
	storageLinker  DisksLinker
	systemConfig   SystemConfig
	certGenerator  CertificateGenerator
	ldap           Ldap
	nginx          Nginx
	web            auth.Web
	logDir         string
	logger         *zap.Logger
}

type SystemConfig interface {
	DiskRoot() string
	InternalDiskDir() string
	DiskLink() string
}

type CertificateGenerator interface {
	Generate() error
}

type DisksLinker interface {
	RelinkDisk(link string, target string) error
}

type Ldap interface {
	Init() error
}

type Nginx interface {
	InitConfig() error
}

const (
	App       = "platform"
	AppDir    = "/snap/platform/current"
	DataDir   = "/var/snap/platform/current"
	CommonDir = "/var/snap/platform/common"
)

func NewInstall(
	storageChecker storage.Checker,
	storageLinker DisksLinker,
	systemConfig SystemConfig,
	certGenerator CertificateGenerator,
	ldap Ldap,
	nginx Nginx,
	web auth.Web,
	logger *zap.Logger,
) *Install {
	return &Install{
		storageChecker: storageChecker,
		storageLinker:  storageLinker,
		systemConfig:   systemConfig,
		certGenerator:  certGenerator,
		ldap:           ldap,
		nginx:          nginx,
		web:            web,
		logDir:         path.Join(CommonDir, "log"),
		logger:         logger,
	}
}

func (i *Install) Install() error {
	err := linux.CreateUser(App)
	if err != nil {
		return err
	}

	err = i.InitConfigs()
	if err != nil {
		return err
	}

	err = i.InitDisk()
	if err != nil {
		return err
	}

	err = i.certGenerator.Generate()
	if err != nil {
		return err
	}
	err = i.ldap.Init()
	if err != nil {
		return err
	}
	err = i.nginx.InitConfig()
	if err != nil {
		return err
	}
	err = i.web.InitConfig()
	if err != nil {
		return err
	}
	return nil
}

func (i *Install) PostRefresh() error {

	err := i.InitConfigs()
	if err != nil {
		return err
	}

	err = i.InitDisk()
	if err != nil {
		return err
	}

	err = i.nginx.InitConfig()
	if err != nil {
		return err
	}

	err = i.web.InitConfig()
	if err != nil {
		return err
	}

	err = cli.Remove(fmt.Sprintf("%s/*.log", i.logDir))
	if err != nil {
		return err
	}

	// Clean up old manually-installed odroidhc4-display service
	// TODO: Remove this cleanup code after sufficient time has passed for devices to upgrade
	err = i.cleanupOldLcdService()
	if err != nil {
		i.logger.Warn("failed to cleanup old LCD service", zap.Error(err))
		// Don't fail the refresh if cleanup fails
	}

	return nil
}

func (i *Install) cleanupOldLcdService() error {
	servicePath := "/etc/systemd/system/odroidhc4-display.service"
	binaryPath := "/usr/bin/odroidhc4-display"

	// Check if the service file exists
	if _, err := os.Stat(servicePath); err == nil {
		i.logger.Info("found old odroidhc4-display service, removing it")

		// Stop and disable the service
		_ = linux.SystemCtl("stop", "odroidhc4-display")
		_ = linux.SystemCtl("disable", "odroidhc4-display")

		// Remove the service file
		if err := os.Remove(servicePath); err != nil {
			i.logger.Warn("failed to remove service file", zap.Error(err))
		}

		// Reload systemd daemon
		_ = linux.SystemCtl("daemon-reload")
	}

	// Check if the binary exists
	if _, err := os.Stat(binaryPath); err == nil {
		i.logger.Info("found old odroidhc4-display binary, removing it")
		if err := os.Remove(binaryPath); err != nil {
			i.logger.Warn("failed to remove binary", zap.Error(err))
		}
	}

	return nil
}

func (i *Install) InitConfigs() error {
	i.logger.Info("init configs")

	dataDirs := []string{
		i.logDir,
		path.Join(DataDir, "nginx"),
		path.Join(DataDir, "openldap"),
		path.Join(DataDir, "openldap-data"),
	}

	for _, dir := range dataDirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	err := cp.Copy(path.Join(AppDir, "certs"), "/usr/share/ca-certificates/mozilla")
	if err != nil {
		return err
	}

	return nil
}

func (i *Install) InitDisk() error {
	i.logger.Info("init disk")
	err := createDir(i.systemConfig.DiskRoot())
	if err != nil {
		return err
	}
	err = createDir(i.systemConfig.InternalDiskDir())
	if err != nil {
		return err
	}

	if !i.storageChecker.ExternalDiskLinkExists() {
		err = i.storageLinker.RelinkDisk(i.systemConfig.DiskLink(), i.systemConfig.InternalDiskDir())
		if err != nil {
			return err
		}
	}
	return nil
}

func createDir(dir string) error {
	_, err := os.Stat(dir)
	if errors.Is(err, os.ErrNotExist) {
		return os.MkdirAll(dir, 0755)
	}
	return err
}
